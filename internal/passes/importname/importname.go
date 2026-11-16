// Package importname checks import bindings and unnecessary import aliases.
package importname

import (
	"cmp"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "importname",
	Doc:  "require lowercase import names without underscores and reject redundant aliases",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		for _, specification := range file.Imports {
			path, err := strconv.Unquote(specification.Path.Value)
			if err != nil || path == "C" {
				continue
			}
			position := specification.Path.Pos()
			end := specification.Path.End()
			if specification.Name != nil {
				position = specification.Name.Pos()
				end = specification.Name.End()
			}
			binding := importBinding(pass.TypesInfo, specification)
			if binding == nil {
				continue
			}
			name := binding.Name()
			redundant := specification.Name != nil && name == binding.Imported().Name()
			var message string
			if !lowercaseName(name) {
				message = "import name " + strconv.Quote(name) + " must use lowercase letters and digits, without underscores"
				if redundant {
					message += "; the alias also repeats the package name"
				}
			} else if redundant {
				message = "import alias " + strconv.Quote(name) + " is redundant; the package already has that name"
			} else {
				continue
			}
			candidate := strings.ReplaceAll(strings.ToLower(name), "_", "")
			pass.Report(analysis.Diagnostic{
				Pos: position, End: end, Message: message,
				SuggestedFixes: importFix(pass, file, specification, binding, candidate),
			})
		}
	}
	return nil, nil
}

func importBinding(info *types.Info, specification *ast.ImportSpec) *types.PkgName {
	object := info.Implicits[specification]
	if specification.Name != nil {
		object = info.Defs[specification.Name]
	}
	binding, ok := object.(*types.PkgName)
	if !ok {
		return nil
	}
	return binding
}

func lowercaseName(name string) bool {
	if name == "" || token.Lookup(name).IsKeyword() {
		return false
	}
	for index, letter := range name {
		if !unicode.IsLower(letter) && !(index > 0 && unicode.IsDigit(letter)) {
			return false
		}
	}
	return true
}

func importFix(pass *analysis.Pass, file *ast.File, specification *ast.ImportSpec, binding *types.PkgName, candidate string) []analysis.SuggestedFix {
	if !lowercaseName(candidate) {
		return nil
	}
	if candidate != binding.Name() && !canRename(pass, file, binding, candidate) {
		return nil
	}
	var edits []analysis.TextEdit
	if specification.Name == nil {
		edits = append(edits, analysis.TextEdit{Pos: specification.Path.Pos(), End: specification.Path.Pos(), NewText: []byte(candidate + " ")})
	} else {
		var replacement []byte
		if candidate != binding.Imported().Name() {
			replacement = []byte(candidate)
		}
		edits = append(edits, analysis.TextEdit{Pos: specification.Name.Pos(), End: specification.Name.End(), NewText: replacement})
	}
	if candidate != binding.Name() {
		for name, object := range pass.TypesInfo.Uses {
			if object == binding {
				edits = append(edits, analysis.TextEdit{Pos: name.Pos(), End: name.End(), NewText: []byte(candidate)})
			}
		}
	}
	slices.SortFunc(edits, func(left, right analysis.TextEdit) int { return cmp.Compare(left.Pos, right.Pos) })
	return []analysis.SuggestedFix{{Message: "Use import name " + strconv.Quote(candidate), TextEdits: edits}}
}

// Renaming a file-scoped import must not collide with another import, a package
// declaration, or a predeclared identifier. Local bindings must not shadow the
// proposed name at any of the package's existing references either.
func canRename(pass *analysis.Pass, file *ast.File, binding *types.PkgName, candidate string) bool {
	fileScope := pass.TypesInfo.Scopes[file]
	if fileScope == nil {
		return false
	}
	if _, object := fileScope.LookupParent(candidate, token.NoPos); object != nil {
		return false
	}
	// Fixes for two invalid aliases must not introduce the same new name when
	// the driver applies all suggested fixes together.
	for _, specification := range file.Imports {
		other := importBinding(pass.TypesInfo, specification)
		if other != nil && other != binding && strings.ReplaceAll(strings.ToLower(other.Name()), "_", "") == candidate {
			return false
		}
	}
	for name, object := range pass.TypesInfo.Uses {
		if object != binding {
			continue
		}
		scope := fileScope.Innermost(name.Pos())
		if scope == nil {
			return false
		}
		if _, shadow := scope.LookupParent(candidate, name.Pos()); shadow != nil {
			return false
		}
	}
	return true
}
