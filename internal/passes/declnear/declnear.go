// Package declnear keeps supporting declarations near their explicit uses.
package declnear

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "declnear",
	Doc:  "place unexported method-free structs near their first explicit function reference",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		// Test entrypoints must precede their helpers, even when they use those
		// helpers. Do not impose a competing first-use ordering in test files.
		if ast.IsGenerated(file) || strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go") {
			continue
		}
		for _, declaration := range file.Decls {
			group, ok := declaration.(*ast.GenDecl)
			if !ok || group.Tok != token.TYPE || group.Lparen.IsValid() || len(group.Specs) != 1 {
				continue
			}
			specification := group.Specs[0].(*ast.TypeSpec)
			if specification.Name.IsExported() || specification.Assign.IsValid() {
				continue
			}
			if _, ok := specification.Type.(*ast.StructType); !ok {
				continue
			}
			object, ok := pass.TypesInfo.Defs[specification.Name].(*types.TypeName)
			if !ok {
				continue
			}
			named, ok := object.Type().(*types.Named)
			if !ok || named.NumMethods() != 0 {
				continue
			}
			first := firstReference(pass.TypesInfo, file, specification, object)
			if first == nil || !interveningFunction(file, group.End(), first.Pos()) {
				continue
			}
			pass.Reportf(specification.Name.Pos(), "move supporting struct %q directly before %s, its first explicit use", object.Name(), first.Name.Name)
		}
	}
	return nil, nil
}

// This first rule considers only types used explicitly by functions in this
// file. Package-level registrations, dependencies between types, references in
// other files, and uses preceding the declaration make placement ambiguous.
func firstReference(info *types.Info, file *ast.File, specification *ast.TypeSpec, object *types.TypeName) *ast.FuncDecl {
	var first *ast.FuncDecl
	for name, used := range info.Uses {
		if used != object || (name.Pos() >= specification.Pos() && name.End() <= specification.End()) {
			continue
		}
		if name.Pos() < file.Pos() || name.End() > file.End() {
			return nil
		}
		function := referencingFunction(file, name.Pos())
		if function == nil || function.Pos() < specification.End() {
			return nil
		}
		if first == nil || function.Pos() < first.Pos() {
			first = function
		}
	}
	return first
}

func referencingFunction(file *ast.File, position token.Pos) *ast.FuncDecl {
	for _, declaration := range file.Decls {
		if declaration.Pos() <= position && position < declaration.End() {
			function, ok := declaration.(*ast.FuncDecl)
			if ok {
				return function
			}
			return nil
		}
	}
	return nil
}

func interveningFunction(file *ast.File, start, end token.Pos) bool {
	for _, declaration := range file.Decls {
		if declaration.Pos() >= start && declaration.End() <= end {
			if _, ok := declaration.(*ast.FuncDecl); ok {
				return true
			}
		}
	}
	return false
}
