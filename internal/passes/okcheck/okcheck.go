// Package okcheck flags discarded boolean results from calls and type assertions.
package okcheck

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "okcheck",
	Doc:  "require preserving boolean results from (value, bool) calls and type assertions",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			var result *ast.Ident
			var values []ast.Expr
			switch node := node.(type) {
			case *ast.AssignStmt:
				if len(node.Lhs) != 2 {
					return true
				}
				identifier, ok := node.Lhs[1].(*ast.Ident)
				if !ok {
					return true
				}
				result = identifier
				values = node.Rhs
			case *ast.ValueSpec:
				if len(node.Names) != 2 {
					return true
				}
				result = node.Names[1]
				values = node.Values
			default:
				return true
			}
			if result == nil || result.Name != "_" || len(values) != 1 {
				return true
			}
			expression := ast.Unparen(values[0])
			if _, ok := expression.(*ast.TypeAssertExpr); ok {
				pass.Reportf(result.Pos(), "do not discard the boolean result of a type assertion")
				return true
			}
			call, ok := expression.(*ast.CallExpr)
			if !ok {
				return true
			}
			results, ok := pass.TypesInfo.TypeOf(call).(*types.Tuple)
			if !ok || results.Len() != 2 {
				return true
			}
			boolean, ok := results.At(1).Type().Underlying().(*types.Basic)
			if ok && boolean.Kind() == types.Bool {
				pass.Reportf(result.Pos(), "do not discard the boolean result of a function returning (value, bool)")
			}
			return true
		})
	}
	return nil, nil
}
