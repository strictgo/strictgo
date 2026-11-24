// Package blankparam requires descriptive names instead of blank parameters.
package blankparam

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "blankparam",
	Doc:  "flag blank function and method parameters and suggest descriptive names",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		owners := functionOwners(file)
		ast.Inspect(file, func(node ast.Node) bool {
			function, ok := node.(*ast.FuncType)
			if ok && function.Params != nil {
				checkParameters(pass, function, owners[function])
			}
			return true
		})
	}
	return nil, nil
}

func functionOwners(file *ast.File) map[*ast.FuncType]ast.Node {
	owners := make(map[*ast.FuncType]ast.Node)
	ast.Inspect(file, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.FuncDecl:
			owners[node.Type] = node
		case *ast.FuncLit:
			owners[node.Type] = node
		}
		return true
	})
	return owners
}

func checkParameters(pass *analysis.Pass, function *ast.FuncType, owner ast.Node) {
	candidates := make(map[*ast.Ident]string)
	counts := make(map[string]int)
	for _, field := range function.Params.List {
		for _, name := range field.Names {
			if name.Name != "_" {
				continue
			}
			variable, ok := pass.TypesInfo.Defs[name].(*types.Var)
			if ok {
				candidate := parameterName(variable.Type())
				candidates[name] = candidate
				counts[candidate]++
			}
		}
	}
	for _, field := range function.Params.List {
		for _, name := range field.Names {
			if name.Name != "_" {
				continue
			}
			diagnostic := analysis.Diagnostic{
				Pos: name.Pos(), End: name.End(),
				Message: "blank parameter obscures its purpose; give it a descriptive name",
			}
			candidate := candidates[name]
			if candidate != "" && counts[candidate] == 1 && canNameParameter(pass.TypesInfo, function, owner, candidate) {
				diagnostic.SuggestedFixes = []analysis.SuggestedFix{{
					Message:   "Name parameter " + candidate,
					TextEdits: []analysis.TextEdit{{Pos: name.Pos(), End: name.End(), NewText: []byte(candidate)}},
				}}
			}
			pass.Report(diagnostic)
		}
	}
}

func parameterName(typ types.Type) string {
	typ = types.Unalias(typ)
	var pointer bool
	for {
		indirect, ok := typ.(*types.Pointer)
		if !ok {
			break
		}
		pointer = true
		typ = types.Unalias(indirect.Elem())
	}
	if types.Identical(typ, types.Universe.Lookup("error").Type()) {
		return "err"
	}
	if named, ok := typ.(*types.Named); ok {
		object := named.Obj()
		if object.Pkg() != nil {
			if object.Pkg().Path() == "context" && object.Name() == "Context" {
				return "ctx"
			}
			if object.Pkg().Path() == "testing" {
				if pointer && (object.Name() == "T" || object.Name() == "B" || object.Name() == "F" || object.Name() == "M") {
					return strings.ToLower(object.Name())
				}
				if object.Name() == "TB" {
					return "tb"
				}
			}
		}
		name := lowerCamel(strings.TrimLeft(object.Name(), "_"))
		if len([]rune(name)) > 1 {
			return name
		}
		return "value"
	}
	switch typ := typ.(type) {
	case *types.Basic:
		if typ.Kind() == types.String {
			return "text"
		}
		if typ.Kind() == types.Bool {
			return "enabled"
		}
	case *types.Slice, *types.Array:
		return "values"
	case *types.Map:
		return "entries"
	case *types.Chan:
		return "channel"
	case *types.Signature:
		return "callback"
	}
	return "value"
}

func lowerCamel(name string) string {
	letters := []rune(name)
	for index, letter := range letters {
		if !unicode.IsUpper(letter) {
			break
		}
		if index > 0 && index+1 < len(letters) && unicode.IsLower(letters[index+1]) {
			break
		}
		letters[index] = unicode.ToLower(letter)
	}
	return string(letters)
}

func canNameParameter(info *types.Info, function *ast.FuncType, owner ast.Node, candidate string) bool {
	if token.Lookup(candidate).IsKeyword() {
		return false
	}
	if scope := info.Scopes[function]; scope != nil && scope.Lookup(candidate) != nil {
		return false
	}
	for _, fields := range []*ast.FieldList{function.Params, function.Results} {
		if containsName(fields, candidate) {
			return false
		}
	}
	var body *ast.BlockStmt
	switch owner := owner.(type) {
	case *ast.FuncDecl:
		if containsName(owner.Recv, candidate) {
			return false
		}
		body = owner.Body
	case *ast.FuncLit:
		body = owner.Body
	}
	if body == nil {
		return true
	}
	var collision bool
	ast.Inspect(body, func(node ast.Node) bool {
		name, ok := node.(*ast.Ident)
		if !ok || name.Name != candidate {
			return true
		}
		object := info.ObjectOf(name)
		if field, ok := object.(*types.Var); ok && field.IsField() {
			return true
		}
		if _, ok := object.(*types.Label); ok {
			return true
		}
		if method, ok := object.(*types.Func); ok && method.Type().(*types.Signature).Recv() != nil {
			return true
		}
		if object != nil {
			collision = true
		}
		return !collision
	})
	return !collision
}

func containsName(fields *ast.FieldList, candidate string) bool {
	if fields == nil {
		return false
	}
	for _, field := range fields.List {
		for _, name := range field.Names {
			if name.Name == candidate {
				return true
			}
		}
	}
	return false
}
