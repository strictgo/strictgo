// Package errcheck requires error nil comparisons to use the form err != nil.
package errcheck

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "require error nil comparisons to use err != nil",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	errorType := types.Universe.Lookup("error").Type()
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			comparison, ok := node.(*ast.BinaryExpr)
			if !ok || (comparison.Op != token.EQL && comparison.Op != token.NEQ) {
				return true
			}
			info := pass.TypesInfo
			errorOnLeft := types.Identical(info.TypeOf(comparison.X), errorType) && info.Types[comparison.Y].IsNil()
			errorOnRight := info.Types[comparison.X].IsNil() && types.Identical(info.TypeOf(comparison.Y), errorType)
			if (errorOnLeft || errorOnRight) && !(errorOnLeft && comparison.Op == token.NEQ) {
				pass.Reportf(comparison.Pos(), "check errors using err != nil")
			}
			return true
		})
	}
	return nil, nil
}
