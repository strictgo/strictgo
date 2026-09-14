// Package badname defines an Analyzer that reports badly named variables.
package badname

import (
	"fmt"
	"go/ast"
	"strings"

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
	filter := []ast.Node{
		(*ast.ValueSpec)(nil),
	}
	inspector.Preorder(filter, func(n ast.Node) {
		spec := n.(*ast.ValueSpec)
		handleMaps(pass, spec)
		handleStrings(pass, spec)
	})
	return nil, nil
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
