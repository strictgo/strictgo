package badname

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"

	"golang.org/x/tools/go/ast/inspector"
)

// Discover counters before checking declarations: a while-style loop often
// increments a variable declared earlier in the enclosing function's scope.
func exemptLoopCounters(info *types.Info, tree *inspector.Inspector, exempt map[types.Object]bool) {
	tree.Preorder([]ast.Node{(*ast.ForStmt)(nil), (*ast.RangeStmt)(nil)}, func(loop ast.Node) {
		ast.Inspect(loop, func(node ast.Node) bool {
			if _, ok := node.(*ast.FuncLit); ok {
				// A callback's increment is not part of the enclosing loop.
				return false
			}
			if statement, ok := node.(ast.Stmt); ok {
				if variable := incrementedVariable(info, statement); variable != nil && variable.Name() == "i" {
					exempt[variable] = true
				}
			}
			return true
		})
	})

	tree.WithStack([]ast.Node{(*ast.ForStmt)(nil), (*ast.RangeStmt)(nil)}, func(loop ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return true
		}
		index := scopedLoopIndex(info, loop)
		if index == nil {
			return true
		}
		if index.Name() == "i" || (index.Name() == "j" && enclosingIndexI(info, stack[:len(stack)-1])) {
			exempt[index] = true
		}
		return true
	})
}

func incrementedVariable(info *types.Info, statement ast.Stmt) *types.Var {
	var expression ast.Expr
	switch statement := statement.(type) {
	case *ast.IncDecStmt:
		if statement.Tok == token.INC {
			expression = statement.X
		}
	case *ast.AssignStmt:
		if statement.Tok == token.ADD_ASSIGN && len(statement.Lhs) == 1 {
			expression = statement.Lhs[0]
		}
	}
	if expression == nil {
		return nil
	}
	name, ok := ast.Unparen(expression).(*ast.Ident)
	if !ok {
		return nil
	}
	variable, ok := info.ObjectOf(name).(*types.Var)
	if !ok || variable.IsField() || variable.Parent() == nil || (variable.Pkg() != nil && variable.Parent() == variable.Pkg().Scope()) {
		return nil
	}
	basic, ok := variable.Type().Underlying().(*types.Basic)
	if !ok || basic.Info()&types.IsInteger == 0 {
		return nil
	}
	return variable
}

// Restrict the new exemptions to indexes declared by the loop itself. A map
// key, channel/iterator value, or unrelated variable assigned by a loop is not
// necessarily an ordinal index, even when its type is integer.
func scopedLoopIndex(info *types.Info, loop ast.Node) *types.Var {
	switch loop := loop.(type) {
	case *ast.RangeStmt:
		if loop.Tok != token.DEFINE || !indexedRange(info.TypeOf(loop.X)) {
			return nil
		}
		if name, ok := loop.Key.(*ast.Ident); ok {
			if variable, ok := info.Defs[name].(*types.Var); ok {
				return variable
			}
		}
	case *ast.ForStmt:
		initialization, ok := loop.Init.(*ast.AssignStmt)
		if !ok || initialization.Tok != token.DEFINE {
			return nil
		}
		variable := incrementedVariable(info, loop.Post)
		if variable == nil {
			return nil
		}
		for _, expression := range initialization.Lhs {
			if name, ok := expression.(*ast.Ident); ok && info.Defs[name] == variable {
				return variable
			}
		}
	}
	return nil
}

func indexedRange(typ types.Type) bool {
	if typ == nil {
		return false
	}
	switch typ := typ.Underlying().(type) {
	case *types.Array, *types.Slice:
		return true
	case *types.Pointer:
		_, ok := typ.Elem().Underlying().(*types.Array)
		return ok
	case *types.Basic:
		return typ.Info()&(types.IsInteger|types.IsString) != 0
	}
	return false
}

func enclosingIndexI(info *types.Info, stack []ast.Node) bool {
	for _, ancestor := range slices.Backward(stack) {
		switch ancestor.(type) {
		case *ast.FuncLit, *ast.FuncDecl:
			return false
		}
		if variable := scopedLoopIndex(info, ancestor); variable != nil && variable.Name() == "i" {
			return true
		}
	}
	return false
}
