package errcheck

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkSwallowedErrors(pass *analysis.Pass, body *ast.BlockStmt, signature *types.Signature) {
	if body == nil {
		return
	}
	ast.Inspect(body, func(node ast.Node) bool {
		if _, ok := node.(*ast.FuncLit); ok {
			// Each literal is checked separately with its own return signature.
			return false
		}
		conditional, ok := node.(*ast.IfStmt)
		if !ok {
			return true
		}
		checked := checkedError(pass.TypesInfo, conditional.Cond)
		if checked == nil {
			return true
		}
		for index, statement := range conditional.Body.List {
			returned, ok := statement.(*ast.ReturnStmt)
			if !ok || !nilErrorReturn(pass.TypesInfo, signature, returned) {
				continue
			}
			if errorReassigned(pass.TypesInfo, conditional.Body.List[:index], checked) {
				continue
			}
			pass.Reportf(returned.Pos(), "do not swallow a non-nil error by returning a nil error")
		}
		return true
	})
}

// Only a direct non-nil check establishes the initial narrow rule. A shadowed
// nil identifier is not a nil literal, and aliases of error retain its identity.
func checkedError(info *types.Info, condition ast.Expr) *types.Var {
	comparison, ok := ast.Unparen(condition).(*ast.BinaryExpr)
	if !ok || comparison.Op != token.NEQ {
		return nil
	}
	var expression ast.Expr
	if info.Types[ast.Unparen(comparison.Y)].IsNil() {
		expression = comparison.X
	} else if info.Types[ast.Unparen(comparison.X)].IsNil() {
		expression = comparison.Y
	} else {
		return nil
	}
	name, ok := ast.Unparen(expression).(*ast.Ident)
	if !ok {
		return nil
	}
	variable, ok := info.ObjectOf(name).(*types.Var)
	if ok && types.Identical(variable.Type(), types.Universe.Lookup("error").Type()) {
		return variable
	}
	return nil
}

func nilErrorReturn(info *types.Info, signature *types.Signature, returned *ast.ReturnStmt) bool {
	results := signature.Results()
	if len(returned.Results) != results.Len() {
		return false
	}
	var found bool
	for index := range results.Len() {
		if !types.Identical(results.At(index).Type(), types.Universe.Lookup("error").Type()) {
			continue
		}
		// Multiple error results must all be nil: returning an error in another
		// position is propagation rather than swallowing.
		if !info.Types[ast.Unparen(returned.Results[index])].IsNil() {
			return false
		}
		found = true
	}
	return found
}

// Reassigning the checked error (or passing its address) may handle or replace
// it before the return. Avoid diagnosing that without control-flow analysis.
func errorReassigned(info *types.Info, statements []ast.Stmt, checked *types.Var) bool {
	var changed bool
	for _, statement := range statements {
		ast.Inspect(statement, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.AssignStmt:
				for _, expression := range node.Lhs {
					if name, ok := ast.Unparen(expression).(*ast.Ident); ok && info.ObjectOf(name) == checked {
						changed = true
					}
				}
			case *ast.UnaryExpr:
				if node.Op == token.AND {
					if name, ok := ast.Unparen(node.X).(*ast.Ident); ok && info.ObjectOf(name) == checked {
						changed = true
					}
				}
			}
			return !changed
		})
		if changed {
			return true
		}
	}
	return false
}
