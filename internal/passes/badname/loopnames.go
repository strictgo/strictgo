package badname

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ast/inspector"
)

// Discover counters before checking declarations: a while-style loop often
// increments a variable declared earlier in the enclosing function's scope.
func exemptLoopCounters(info *types.Info, tree *inspector.Inspector, exempt map[types.Object]bool) {
	tree.Preorder([]ast.Node{(*ast.ForStmt)(nil), (*ast.RangeStmt)(nil)}, func(loop ast.Node) {
		ast.Inspect(loop, func(node ast.Node) bool {
			var expression ast.Expr
			switch node := node.(type) {
			case *ast.FuncLit:
				// A callback's increment is not part of the enclosing loop.
				return false
			case *ast.IncDecStmt:
				if node.Tok == token.INC {
					expression = node.X
				}
			case *ast.AssignStmt:
				if node.Tok == token.ADD_ASSIGN && len(node.Lhs) == 1 {
					expression = node.Lhs[0]
				}
			}
			if expression == nil {
				return true
			}
			name, ok := ast.Unparen(expression).(*ast.Ident)
			if !ok || name.Name != "i" {
				return true
			}
			variable, ok := info.ObjectOf(name).(*types.Var)
			if !ok || variable.IsField() || variable.Parent() == nil || (variable.Pkg() != nil && variable.Parent() == variable.Pkg().Scope()) {
				return true
			}
			basic, ok := variable.Type().Underlying().(*types.Basic)
			if ok && basic.Info()&types.IsInteger != 0 {
				exempt[variable] = true
			}
			return true
		})
	})
}
