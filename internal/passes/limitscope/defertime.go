package limitscope

import (
	"bytes"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

func checkDeferredTimes(pass *analysis.Pass, list []ast.Stmt) {
	for index := range len(list) - 1 {
		assignment, ok := list[index].(*ast.AssignStmt)
		if !ok || assignment.Tok != token.DEFINE || len(assignment.Lhs) != 1 || len(assignment.Rhs) != 1 {
			continue
		}
		name, ok := assignment.Lhs[0].(*ast.Ident)
		if !ok || name.Name == "_" {
			continue
		}
		call, ok := ast.Unparen(assignment.Rhs[0]).(*ast.CallExpr)
		if !ok {
			continue
		}
		callee := typeutil.StaticCallee(pass.TypesInfo, call)
		if callee == nil || callee.Pkg() == nil || callee.Pkg().Path() != "time" || callee.Name() != "Now" {
			continue
		}
		deferred, ok := list[index+1].(*ast.DeferStmt)
		if !ok || len(deferred.Call.Args) != 0 {
			continue
		}
		callback, ok := ast.Unparen(deferred.Call.Fun).(*ast.FuncLit)
		if !ok || len(callback.Type.Params.List) != 0 {
			continue
		}
		variable := pass.TypesInfo.Defs[name]
		if !usedOnlyIn(pass.TypesInfo, variable, callback.Body) {
			continue
		}
		// Parameters and direct body declarations share a scope. Adding a
		// parameter must not invalidate an existing short declaration or type.
		if scope := pass.TypesInfo.Scopes[callback.Type]; scope != nil && scope.Lookup(name.Name) != nil {
			continue
		}
		pass.Report(analysis.Diagnostic{
			Pos: assignment.Pos(), End: assignment.End(),
			Message:        "pass " + name.Name + " as a deferred function argument to reduce scope",
			SuggestedFixes: deferredTimeFix(pass, assignment, deferred, callback, name.Name, call),
		})
	}
}

func usedOnlyIn(info *types.Info, variable types.Object, body *ast.BlockStmt) bool {
	if variable == nil {
		return false
	}
	var used bool
	for name, object := range info.Uses {
		if object != variable {
			continue
		}
		if name.Pos() < body.Pos() || name.End() > body.End() {
			return false
		}
		used = true
	}
	return used
}

func deferredTimeFix(pass *analysis.Pass, assignment *ast.AssignStmt, deferred *ast.DeferStmt, callback *ast.FuncLit, name string, call *ast.CallExpr) []analysis.SuggestedFix {
	typeName, ok := deferredTimeType(pass, call, deferred.Pos())
	if !ok {
		return nil
	}
	// Comments outside the declaration survive its removal. Leave declarations
	// with internal comments untouched rather than silently discarding comments.
	for _, file := range pass.Files {
		for _, group := range file.Comments {
			if group.Pos() < assignment.End() && group.End() > assignment.Pos() {
				return nil
			}
		}
	}
	file := pass.Fset.File(assignment.Pos())
	source, err := pass.ReadFile(file.Name())
	if err != nil {
		return nil
	}
	value := assignment.Rhs[0]
	argument := source[file.Offset(value.Pos()):file.Offset(value.End())]
	end := assignment.End()
	tail := source[file.Offset(end):]
	if newline := bytes.IndexByte(tail, '\n'); newline >= 0 && len(bytes.TrimSpace(tail[:newline])) == 0 {
		end = file.Pos(file.Offset(end) + newline + 1)
	}
	return []analysis.SuggestedFix{{
		Message: "Pass time.Now() directly to the deferred function",
		TextEdits: []analysis.TextEdit{
			{Pos: assignment.Pos(), End: end},
			{Pos: callback.Type.Params.Closing, End: callback.Type.Params.Closing, NewText: []byte(name + " " + typeName)},
			{Pos: deferred.Call.Rparen, End: deferred.Call.Rparen, NewText: argument},
		},
	}}
}

// Reuse the original time import spelling, including aliases and dot imports.
// A qualifier shadowed by the removed variable is not usable in the signature.
func deferredTimeType(pass *analysis.Pass, call *ast.CallExpr, position token.Pos) (string, bool) {
	scope := pass.Pkg.Scope().Innermost(position)
	if scope == nil {
		return "", false
	}
	switch function := ast.Unparen(call.Fun).(type) {
	case *ast.SelectorExpr:
		qualifier, ok := function.X.(*ast.Ident)
		if !ok {
			return "", false
		}
		_, object := scope.LookupParent(qualifier.Name, position)
		if object == pass.TypesInfo.ObjectOf(qualifier) {
			return qualifier.Name + ".Time", true
		}
	case *ast.Ident:
		_, object := scope.LookupParent("Time", position)
		named, ok := pass.TypesInfo.TypeOf(call).(*types.Named)
		if ok && object == named.Obj() {
			return "Time", true
		}
	}
	return "", false
}
