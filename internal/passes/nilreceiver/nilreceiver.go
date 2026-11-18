// Package nilreceiver flags nil checks on method receivers.
package nilreceiver

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "nilreceiver",
	Doc:  "flag nil checks on method receivers",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		for _, declaration := range file.Decls {
			method, ok := declaration.(*ast.FuncDecl)
			if !ok || method.Recv == nil || method.Body == nil {
				continue
			}
			for _, field := range method.Recv.List {
				for _, name := range field.Names {
					receiver, ok := pass.TypesInfo.Defs[name].(*types.Var)
					if ok && name.Name != "_" {
						checkReceiver(pass, method.Body, receiver)
					}
				}
			}
		}
	}
	return nil, nil
}

func checkReceiver(pass *analysis.Pass, body *ast.BlockStmt, receiver *types.Var) {
	ast.Inspect(body, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.BinaryExpr:
			if node.Op != token.EQL && node.Op != token.NEQ {
				return true
			}
			if (isReceiver(pass.TypesInfo, node.X, receiver) && isNil(pass.TypesInfo, node.Y)) ||
				(isNil(pass.TypesInfo, node.X) && isReceiver(pass.TypesInfo, node.Y, receiver)) {
				pass.Reportf(node.Pos(), "do not check method receiver %q for nil", receiver.Name())
			}
		case *ast.SwitchStmt:
			if node.Tag == nil {
				return true
			}
			if !isReceiver(pass.TypesInfo, node.Tag, receiver) {
				return true
			}
			for _, statement := range node.Body.List {
				clause := statement.(*ast.CaseClause)
				for _, expression := range clause.List {
					if isNil(pass.TypesInfo, expression) {
						pass.Reportf(expression.Pos(), "do not check method receiver %q for nil", receiver.Name())
					}
				}
			}
		}
		return true
	})
}

func isReceiver(info *types.Info, expression ast.Expr, receiver *types.Var) bool {
	name, ok := ast.Unparen(expression).(*ast.Ident)
	return ok && info.ObjectOf(name) == receiver
}

func isNil(info *types.Info, expression ast.Expr) bool {
	return info.Types[ast.Unparen(expression)].IsNil()
}
