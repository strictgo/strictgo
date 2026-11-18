// Package declorder enforces a consistent order for top-level declarations.
package declorder

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

const (
	constantDeclarations = iota
	variableDeclarations
	mainEntrypoints
	testEntrypoints
	typeDeclarations
	constructors
	analysisEntrypoints
	methods
	functions
)

var Analyzer = &analysis.Analyzer{
	Name: "declorder",
	Doc:  "order file preambles and type blocks, with test entrypoints before test helpers",
	Run:  run,
}

var categoryNames = [...]string{
	"const declarations",
	"var declarations",
	"main entrypoints",
	"test entrypoints",
	"type declarations",
	"constructors",
	"analysis entrypoints",
	"methods",
	"functions",
}

func run(pass *analysis.Pass) (any, error) {
	entrypoints := registeredEntrypoints(pass)
	for _, file := range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		if isTestFile(pass, file) {
			checkTestOrder(pass, file)
		}
		checkFile(pass, file, entrypoints)
	}
	return nil, nil
}

// Preamble declarations are ordered against the whole file. Each constructed
// type has independent ordering state, so a later type block can follow an
// earlier block's methods or free functions. Free helpers belong to the most
// recently encountered block; returning to an earlier type does not reset it.
func checkFile(pass *analysis.Pass, file *ast.File, entrypoints map[types.Object]bool) {
	mainTypes := constructedTypes(pass, file)
	testFile := isTestFile(pass, file)
	groups := make(map[types.Object]int)
	var active types.Object
	highest := -1
	for _, declaration := range file.Decls {
		if function, ok := declaration.(*ast.FuncDecl); testFile && ok && testEntrypoint(pass.TypesInfo, function) {
			// Tests lead their own section rather than belonging to a helper
			// type block or pushing subsequent helper methods below functions.
			highest = max(highest, testEntrypoints)
			continue
		}
		category, ok := declarationCategory(declaration, pass.TypesInfo, mainTypes, entrypoints)
		if !ok {
			continue
		}
		owners := declarationOwners(declaration, pass.TypesInfo, mainTypes)
		previous := highest
		if len(owners) > 0 {
			previous = -1
			for _, owner := range owners {
				previous = max(previous, groups[owner])
			}
			active = owners[len(owners)-1]
		} else if category > mainEntrypoints && active != nil {
			owners = []types.Object{active}
			previous = groups[active]
		}
		if category < previous {
			pass.Reportf(declaration.Pos(), "%s must appear before %s", categoryNames[category], categoryNames[previous])
		}
		for _, owner := range owners {
			groups[owner] = max(groups[owner], category)
		}
		highest = max(highest, category)
	}
}

// Resolve the actual Analyzer type and Run function rather than assuming every
// function named run is an entrypoint. Scan the package so registration and
// implementation may live in different files, and import/type aliases work.
func registeredEntrypoints(pass *analysis.Pass) map[types.Object]bool {
	entrypoints := make(map[types.Object]bool)
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			literal, ok := node.(*ast.CompositeLit)
			if !ok {
				return true
			}
			named, ok := types.Unalias(pass.TypesInfo.TypeOf(literal)).(*types.Named)
			if !ok || named.Obj().Name() != "Analyzer" || named.Obj().Pkg() == nil || named.Obj().Pkg().Path() != "golang.org/x/tools/go/analysis" {
				return true
			}
			for _, element := range literal.Elts {
				field, ok := element.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, ok := field.Key.(*ast.Ident)
				if !ok || key.Name != "Run" {
					continue
				}
				var name *ast.Ident
				switch value := ast.Unparen(field.Value).(type) {
				case *ast.Ident:
					name = value
				case *ast.SelectorExpr:
					name = value.Sel
				}
				if name != nil {
					function, ok := pass.TypesInfo.ObjectOf(name).(*types.Func)
					if ok && function.Pkg() == pass.Pkg {
						entrypoints[function] = true
					}
				}
			}
			return true
		})
	}
	return entrypoints
}

// A main type has a conventionally named constructor in this file returning
// that type (or a pointer to it). Supporting types do not establish an ordering
// boundary, even when a different file supplies a constructor for them.
func constructedTypes(pass *analysis.Pass, file *ast.File) map[types.Object]bool {
	mainTypes := make(map[types.Object]bool)
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if named := constructorType(pass.TypesInfo, function); named != nil && named.Obj().Pkg() == pass.Pkg {
			mainTypes[named.Obj()] = true
		}
	}
	return mainTypes
}

func declarationOwners(declaration ast.Decl, info *types.Info, mainTypes map[types.Object]bool) []types.Object {
	switch declaration := declaration.(type) {
	case *ast.GenDecl:
		if declaration.Tok == token.TYPE {
			var owners []types.Object
			for _, specification := range declaration.Specs {
				typ := specification.(*ast.TypeSpec)
				if object := info.Defs[typ.Name]; mainTypes[object] {
					owners = append(owners, object)
				}
			}
			return owners
		}
	case *ast.FuncDecl:
		var named *types.Named
		if function, ok := info.Defs[declaration.Name].(*types.Func); ok {
			signature := function.Type().(*types.Signature)
			if receiver := signature.Recv(); receiver != nil {
				named = namedType(receiver.Type())
			} else {
				named = constructorType(info, declaration)
			}
		}
		if named != nil && mainTypes[named.Obj()] {
			return []types.Object{named.Obj()}
		}
	}
	return nil
}

func constructorType(info *types.Info, function *ast.FuncDecl) *types.Named {
	if function.Recv != nil || !constructorName(function.Name.Name) {
		return nil
	}
	object, ok := info.Defs[function.Name].(*types.Func)
	if !ok {
		return nil
	}
	results := object.Type().(*types.Signature).Results()
	if results.Len() == 0 {
		return nil
	}
	return namedType(results.At(0).Type())
}

func namedType(typ types.Type) *types.Named {
	typ = types.Unalias(typ)
	if pointer, ok := typ.(*types.Pointer); ok {
		typ = types.Unalias(pointer.Elem())
	}
	if named, ok := typ.(*types.Named); ok {
		return named.Origin()
	}
	return nil
}

func declarationCategory(declaration ast.Decl, info *types.Info, mainTypes, entrypoints map[types.Object]bool) (int, bool) {
	switch declaration := declaration.(type) {
	case *ast.GenDecl:
		switch declaration.Tok {
		case token.CONST:
			return constantDeclarations, true
		case token.VAR:
			return variableDeclarations, true
		case token.TYPE:
			for _, specification := range declaration.Specs {
				typ := specification.(*ast.TypeSpec)
				if mainTypes[info.Defs[typ.Name]] {
					return typeDeclarations, true
				}
			}
		}
	case *ast.FuncDecl:
		function, ok := info.Defs[declaration.Name].(*types.Func)
		if ok && declaration.Recv == nil && function.Name() == "main" && function.Pkg() != nil && function.Pkg().Name() == "main" {
			return mainEntrypoints, true
		}
		if entrypoints[info.Defs[declaration.Name]] {
			return analysisEntrypoints, true
		}
		if declaration.Recv != nil {
			return methods, true
		}
		if constructorName(declaration.Name.Name) {
			return constructors, true
		}
		return functions, true
	}
	return 0, false
}

// Use naming conventions for this initial rule, not a guess about what a
// function constructs. New, NewClient, and newClient are constructors;
// Newline and newer are ordinary functions. Return types identify which type
// declarations are main types, but do not affect the function's own category.
func constructorName(name string) bool {
	for _, prefix := range []string{"New", "new"} {
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		suffix := strings.TrimPrefix(name, prefix)
		if suffix == "" {
			return true
		}
		letter, _ := utf8.DecodeRuneInString(suffix)
		return unicode.IsUpper(letter)
	}
	return false
}
