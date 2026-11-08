// Package badname defines an Analyzer that reports badly named variables.
package badname

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "badname",
	Doc:      "flags badly named variables",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		// explicit check: should never be false due to inspect.Analyzer dependency
		return nil, fmt.Errorf("unexpected type: %T", pass.ResultOf[inspect.Analyzer])
	}
	receivers := make(map[types.Object]bool)
	filter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.Ident)(nil),
		(*ast.ValueSpec)(nil),
	}
	inspector.Preorder(filter, func(node ast.Node) {
		switch node := node.(type) {
		case *ast.FuncDecl:
			if node.Recv != nil {
				for _, field := range node.Recv.List {
					for _, name := range field.Names {
						receivers[pass.TypesInfo.Defs[name]] = true
					}
				}
			}
		case *ast.Ident:
			handleSingleLetter(pass, node, receivers)
		case *ast.ValueSpec:
			handleMaps(pass, node)
			handleStrings(pass, node)
		}
	})
	return nil, nil
}

// maxShortScopeLines is the inclusive line span from a variable's declaration
// to the end of its lexical scope. Using the scope rather than the last use
// prevents unused or rarely used variables from appearing artificially local.
const maxShortScopeLines = 7

func handleSingleLetter(pass *analysis.Pass, name *ast.Ident, receivers map[types.Object]bool) {
	letter, size := utf8.DecodeRuneInString(name.Name)
	if size != len(name.Name) || !unicode.IsLetter(letter) {
		return
	}
	variable, ok := pass.TypesInfo.Defs[name].(*types.Var)
	if !ok || variable.IsField() || receivers[variable] || conventionalTestingName(variable) {
		return
	}
	scope := variable.Parent()
	if scope == nil {
		return
	}
	if scope != pass.Pkg.Scope() {
		start := pass.Fset.PositionFor(name.Pos(), false).Line
		end := pass.Fset.PositionFor(scope.End(), false).Line
		if end-start+1 <= maxShortScopeLines {
			return
		}
	}
	pass.Reportf(name.Pos(), "variable %q has a single-letter name outside a short scope", name.Name)
}

// Match type identity rather than source spelling so aliases work, while
// similarly named types from other packages do not receive an exemption.
func conventionalTestingName(variable *types.Var) bool {
	typ := types.Unalias(variable.Type())
	pointer, isPointer := typ.(*types.Pointer)
	if isPointer {
		typ = types.Unalias(pointer.Elem())
	}
	named, ok := typ.(*types.Named)
	if !ok || named.Obj().Pkg() == nil || named.Obj().Pkg().Path() != "testing" {
		return false
	}
	name := variable.Name()
	if !isPointer {
		return named.Obj().Name() == "TB" && (name == "t" || name == "b" || name == "f")
	}
	switch named.Obj().Name() {
	case "T":
		return name == "t"
	case "B":
		return name == "b"
	case "F":
		return name == "f"
	case "M":
		return name == "m"
	}
	return false
}

func handleStrings(pass *analysis.Pass, spec *ast.ValueSpec) {
	ident, ok := spec.Type.(*ast.Ident)
	if !ok || ident.Name != "string" {
		return
	}
	for _, name := range spec.Names {
		if strings.HasSuffix(name.Name, "Str") {
			pass.Reportf(
				name.Pos(),
				"variable %q unnecessarily contains type information with suffix %q",
				name.Name,
				"Str",
			)
		}
	}
}

func handleMaps(pass *analysis.Pass, spec *ast.ValueSpec) {
	if _, ok := spec.Type.(*ast.MapType); !ok {
		return
	}
	for _, name := range spec.Names {
		if strings.HasSuffix(name.Name, "Map") {
			pass.Reportf(
				name.Pos(),
				"variable %q unnecessarily contains type information with suffix %q",
				name.Name,
				"Map",
			)
		}
	}
}
