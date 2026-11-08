// Package forint flags traditional for loops with integer counters.
package forint

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "forint",
	Doc:  "flag traditional for loops with integer counters",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			loop, ok := n.(*ast.ForStmt)
			if !ok {
				return true
			}
			if integerControl(pass.TypesInfo, loop.Init) || integerControl(pass.TypesInfo, loop.Post) {
				pass.Reportf(loop.For, "avoid for loops with integer counters")
			}
			return true
		})
	}
	return nil, nil
}

func integerControl(info *types.Info, stmt ast.Stmt) bool {
	var expressions []ast.Expr
	switch stmt := stmt.(type) {
	case *ast.AssignStmt:
		expressions = stmt.Lhs
	case *ast.IncDecStmt:
		expressions = []ast.Expr{stmt.X}
	}
	for _, expr := range expressions {
		if id, ok := expr.(*ast.Ident); ok && id.Name == "_" {
			continue
		}
		if typ := info.TypeOf(expr); typ != nil {
			if basic, ok := typ.Underlying().(*types.Basic); ok && basic.Info()&types.IsInteger != 0 {
				return true
			}
		}
	}
	return false
}
