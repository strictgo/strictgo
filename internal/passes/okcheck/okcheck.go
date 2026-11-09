// Package okcheck flags discarded boolean results from two-result calls.
package okcheck

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "okcheck",
	Doc:  "require preserving the boolean result of functions returning (value, bool)",
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
				result, _ = node.Lhs[1].(*ast.Ident)
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
			call, ok := ast.Unparen(values[0]).(*ast.CallExpr)
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
