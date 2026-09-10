// Package limitscope defines an Analyzer that flags opportunities to reduce
// scope.
package limitscope

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "limitscope",
	Doc:      "hints and fixes to reduce scope",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		// explicit check: should never be false due to inspect.Analyzer dependency
		return nil, fmt.Errorf("unexpected type: %T", pass.ResultOf[inspect.Analyzer])
	}
	inspector.Preorder([]ast.Node{(*ast.BlockStmt)(nil)}, func(n ast.Node) {
		block := n.(*ast.BlockStmt)
		check(pass, block.List)
	})
	return nil, nil
}

func check(pass *analysis.Pass, list []ast.Stmt) {
	for i := range len(list) - 1 {
		assign, ok := list[i].(*ast.AssignStmt)
		if !ok {
			continue
		}
		if assign.Tok != token.DEFINE {
			continue
		}
		if len(assign.Lhs) != 1 {
			continue
		}
		if id, ok := assign.Lhs[0].(*ast.Ident); !ok || id.Name != "err" {
			continue
		}

		ifStmt, ok := list[i+1].(*ast.IfStmt)
		if !ok {
			continue
		}

		if ifStmt.Init != nil {
			continue
		}

		binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
		if !ok {
			continue
		}
		rightOp, ok := binExpr.Y.(*ast.Ident)
		if !ok {
			continue
		}
		if binExpr.Op != token.NEQ || rightOp.Name != "nil" {
			continue
		}

		pass.Report(analysis.Diagnostic{
			Pos:     assign.Pos(),
			End:     ifStmt.Cond.Pos(),
			Message: "inline err assignment and check to reduce scope",
		})
		return
	}
}
