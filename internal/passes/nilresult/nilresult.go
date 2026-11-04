// Package nilresult checks that pointer results and errors are mutually exclusive.
package nilresult

import (
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

var Analyzer = &analysis.Analyzer{
	Name:     "nilresult",
	Doc:      "require a nil result or nil error when returning (*T, error)",
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
		if _, ok := results.At(0).Type().Underlying().(*types.Pointer); !ok {
			continue
		}
		for _, block := range fn.Blocks {
			for _, instruction := range block.Instrs {
				ret, ok := instruction.(*ssa.Return)
				if !ok {
					continue
				}
				if isNil(ret.Results[0], block, make(map[visit]bool)) || isNil(ret.Results[1], block, make(map[visit]bool)) {
					continue
				}
				pos := ret.Pos()
				if !pos.IsValid() {
					pos = fn.Pos()
				}
				pass.Reportf(pos, "return a nil result or a nil error: both returned values may be non-nil")
			}
		}
	}
	return nil, nil
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
		edge := 0
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
