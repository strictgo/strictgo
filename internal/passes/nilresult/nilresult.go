// Package nilresult checks that error returns have nil pointer or zero struct results.
package nilresult

import (
	"go/constant"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

var Analyzer = &analysis.Analyzer{
	Name:     "nilresult",
	Doc:      "require a nil pointer or zero struct result when returning an error",
	Requires: []*analysis.Analyzer{buildssa.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	errorType := types.Universe.Lookup("error").Type()
	built := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	functions := append([]*ssa.Function(nil), built.SrcFuncs...)
	functions = append(functions, built.Pkg.Func("init"))
	seenFunctions := make(map[*ssa.Function]bool)
	for len(functions) > 0 {
		fn := functions[0]
		functions = functions[1:]
		if fn == nil || seenFunctions[fn] {
			continue
		}
		seenFunctions[fn] = true
		functions = append(functions, fn.AnonFuncs...)
		results := fn.Signature.Results()
		if results.Len() != 2 || !types.Identical(results.At(1).Type(), errorType) {
			continue
		}
		var structResult bool
		switch results.At(0).Type().Underlying().(type) {
		case *types.Pointer:
		case *types.Struct:
			structResult = true
		default:
			continue
		}
		for _, block := range fn.Blocks {
			// Unnamed result slots cannot be modified by deferred code. A
			// recovered panic either returns their initial zeros or the values
			// from an explicit return, which we check at its own return site.
			if block == fn.Recover && results.At(0).Name() == "" && results.At(1).Name() == "" {
				continue
			}
			for _, instruction := range block.Instrs {
				ret, ok := instruction.(*ssa.Return)
				if !ok {
					continue
				}
				if isNil(ret.Results[1], block, make(map[visit]bool)) {
					continue
				}
				if structResult {
					if isZero(ret.Results[0], block, make(map[visit]bool)) {
						continue
					}
				} else if isNil(ret.Results[0], block, make(map[visit]bool)) {
					continue
				}
				pos := ret.Pos()
				if !pos.IsValid() {
					pos = fn.Pos()
				}
				if structResult {
					pass.Reportf(pos, "return a zero struct result or a nil error: use %s{}, err in error cases", results.At(0).Type())
				} else {
					pass.Reportf(pos, "return a nil result or a nil error: both returned values may be non-nil")
				}
			}
		}
	}
	return nil, nil
}

// isZero recognizes zero constants, conversions, merges, and zero-initialized
// aggregate storage. Unknown values are deliberately not assumed to be zero.
func isZero(value ssa.Value, block *ssa.BasicBlock, seen map[visit]bool) bool {
	key := visit{value, block}
	if seen[key] {
		return false
	}
	seen[key] = true
	defer delete(seen, key)
	if c, ok := value.(*ssa.Const); ok {
		// SSA also represents zero structs and arrays with a nil Value.
		if c.Value == nil {
			return true
		}
		switch c.Value.Kind() {
		case constant.Bool:
			return !constant.BoolVal(c.Value)
		case constant.String:
			return constant.StringVal(c.Value) == ""
		case constant.Int, constant.Float, constant.Complex:
			return constant.Compare(c.Value, token.EQL, constant.MakeInt64(0))
		}
	}
	switch value := value.(type) {
	case *ssa.ChangeType:
		return isZero(value.X, block, seen)
	case *ssa.Convert:
		return isZero(value.X, block, seen)
	case *ssa.Phi:
		for i, edge := range value.Edges {
			if !isZero(edge, value.Block().Preds[i], seen) {
				return false
			}
		}
		return len(value.Edges) > 0
	case *ssa.UnOp:
		if value.Op == token.MUL {
			if alloc, ok := value.X.(*ssa.Alloc); ok {
				return loadedValueIs(value, isZero, seen) || zeroStorage(alloc, seen)
			}
		}
	}
	return false
}

// zeroStorage proves that a fresh allocation only receives zero writes and
// never escapes to code that could mutate it. Reads do not affect this proof.
func zeroStorage(address ssa.Value, seen map[visit]bool) bool {
	refs := address.Referrers()
	if refs == nil {
		return false
	}
	for _, ref := range *refs {
		switch ref := ref.(type) {
		case *ssa.Store:
			if ref.Addr != address || !isZero(ref.Val, ref.Block(), seen) {
				return false
			}
		case *ssa.FieldAddr:
			if !zeroStorage(ref, seen) {
				return false
			}
		case *ssa.IndexAddr:
			if !zeroStorage(ref, seen) {
				return false
			}
		case *ssa.UnOp:
			if ref.Op != token.MUL {
				return false
			}
		case *ssa.DebugRef:
		default:
			return false
		}
	}
	return true
}

// loadedValueIs follows whole-value stores reaching a load, rather than all
// writes during an allocation's lifetime. SSA uses shared result slots in
// functions with defer, even for explicit return expressions.
func loadedValueIs(load *ssa.UnOp, prove func(ssa.Value, *ssa.BasicBlock, map[visit]bool) bool, seen map[visit]bool) bool {
	alloc, ok := load.X.(*ssa.Alloc)
	if !ok || load.Op != token.MUL {
		return false
	}
	// Only private storage accessed by whole-value stores and loads is safe
	// here. In particular, deferred closures may mutate captured result slots.
	for _, ref := range *alloc.Referrers() {
		switch ref := ref.(type) {
		case *ssa.Store:
			if ref.Addr != alloc {
				return false
			}
		case *ssa.UnOp:
			if ref.Op != token.MUL {
				return false
			}
		case *ssa.DebugRef:
		default:
			return false
		}
	}
	visiting := make(map[*ssa.BasicBlock]bool)
	var reaching func(*ssa.BasicBlock, int) bool
	reaching = func(block *ssa.BasicBlock, end int) bool {
		if visiting[block] {
			return false
		}
		visiting[block] = true
		defer delete(visiting, block)
		for i := end - 1; i >= 0; i-- {
			instruction := block.Instrs[i]
			if store, ok := instruction.(*ssa.Store); ok && store.Addr == alloc {
				return prove(store.Val, block, seen)
			}
			if instruction == alloc {
				return true // Fresh allocations contain the zero value.
			}
		}
		for _, pred := range block.Preds {
			if !reaching(pred, len(pred.Instrs)) {
				return false
			}
		}
		return len(block.Preds) > 0
	}
	for i, instruction := range load.Block().Instrs {
		if instruction == load {
			return reaching(load.Block(), i)
		}
	}
	return false
}

type visit struct {
	value ssa.Value
	block *ssa.BasicBlock
}

// isNil accepts only values proven nil, including facts established by guards.
// A pointer wrapped in an error interface is not a nil interface.
func isNil(value ssa.Value, block *ssa.BasicBlock, seen map[visit]bool) bool {
	key := visit{value, block}
	if seen[key] {
		return false
	}
	seen[key] = true
	defer delete(seen, key)
	if constant, ok := value.(*ssa.Const); ok {
		return constant.IsNil()
	}
	for ancestor := block; ancestor != nil; ancestor = ancestor.Idom() {
		if len(ancestor.Instrs) == 0 {
			continue
		}
		branch, ok := ancestor.Instrs[len(ancestor.Instrs)-1].(*ssa.If)
		if !ok {
			continue
		}
		comparison, ok := branch.Cond.(*ssa.BinOp)
		if !ok || (comparison.Op != token.EQL && comparison.Op != token.NEQ) {
			continue
		}
		var other ssa.Value
		if comparison.X == value {
			other = comparison.Y
		}
		if comparison.Y == value {
			other = comparison.X
		}
		constant, ok := other.(*ssa.Const)
		if !ok || !constant.IsNil() {
			continue
		}
		var edge int
		if comparison.Op == token.NEQ {
			edge = 1
		}
		if ancestor.Succs[edge].Dominates(block) {
			return true
		}
	}
	switch value := value.(type) {
	case *ssa.ChangeType:
		return isNil(value.X, block, seen)
	case *ssa.Convert:
		return isNil(value.X, block, seen)
	case *ssa.ChangeInterface:
		return isNil(value.X, block, seen)
	case *ssa.UnOp:
		return loadedValueIs(value, isNil, seen)
	case *ssa.Phi:
		for i, edge := range value.Edges {
			if !isNil(edge, value.Block().Preds[i], seen) {
				return false
			}
		}
		return len(value.Edges) > 0
	}
	return false
}
