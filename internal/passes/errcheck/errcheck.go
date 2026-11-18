// Package errcheck enforces consistent error checks and flags swallowed errors.
package errcheck

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "require err != nil comparisons and flag nil error returns in non-nil error branches",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	errorType := types.Universe.Lookup("error").Type()
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.FuncDecl:
				if function, ok := pass.TypesInfo.Defs[node.Name].(*types.Func); ok {
					checkSwallowedErrors(pass, node.Body, function.Type().(*types.Signature))
				}
			case *ast.FuncLit:
				if signature, ok := pass.TypesInfo.TypeOf(node).(*types.Signature); ok {
					checkSwallowedErrors(pass, node.Body, signature)
				}
			}
			comparison, ok := node.(*ast.BinaryExpr)
			if !ok || (comparison.Op != token.EQL && comparison.Op != token.NEQ) {
				return true
			}
			info := pass.TypesInfo
			errorOnLeft := types.Identical(info.TypeOf(comparison.X), errorType) && info.Types[comparison.Y].IsNil()
			errorOnRight := info.Types[comparison.X].IsNil() && types.Identical(info.TypeOf(comparison.Y), errorType)
			if (errorOnLeft || errorOnRight) && (!errorOnLeft || comparison.Op != token.NEQ) {
				pass.Reportf(comparison.Pos(), "check errors using err != nil")
			}
			return true
		})
	}
	return nil, nil
}
