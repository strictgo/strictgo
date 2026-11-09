// Package forcederef flags dereferences used to satisfy local call signatures.
package forcederef

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

var Analyzer = &analysis.Analyzer{
	Name:      "forcederef",
	Doc:       "flag explicit argument dereferences for functions and methods in the same module",
	Run:       run,
	FactTypes: []analysis.Fact{new(moduleFact)},
}

// Module facts distinguish same-module packages from dependencies, including
// dependencies whose import paths share the current module's prefix.
type moduleFact struct {
	Path string
}

func (*moduleFact) AFact() {}

func (fact *moduleFact) String() string {
	return "module " + fact.Path
}

func run(pass *analysis.Pass) (any, error) {
	if pass.Module != nil {
		pass.ExportPackageFact(&moduleFact{Path: pass.Module.Path})
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			callee, ok := typeutil.Callee(pass.TypesInfo, call).(*types.Func)
			if !ok || !localPackage(pass, callee.Pkg()) {
				return true
			}
			signature, ok := pass.TypesInfo.TypeOf(call.Fun).(*types.Signature)
			if !ok {
				return true
			}
			for index, argument := range call.Args {
				deref, ok := ast.Unparen(argument).(*ast.StarExpr)
				if !ok {
					continue
				}
				parameters := signature.Params()
				parameterIndex := index
				if signature.Variadic() && index >= parameters.Len()-1 {
					parameterIndex = parameters.Len() - 1
				}
				if parameterIndex >= parameters.Len() {
					continue
				}
				parameter := parameters.At(parameterIndex).Type()
				if signature.Variadic() && index >= parameters.Len()-1 && !call.Ellipsis.IsValid() {
					parameter = parameter.(*types.Slice).Elem()
				}
				if types.AssignableTo(pass.TypesInfo.TypeOf(deref.X), parameter) {
					continue
				}
				pass.Reportf(deref.Pos(), "explicit dereference to call %s; align pointer/value types across local calls", callee.Name())
			}
			return true
		})
	}
	return nil, nil
}

func localPackage(pass *analysis.Pass, pkg *types.Package) bool {
	if pkg == pass.Pkg {
		return true
	}
	if pkg == nil || pass.Module == nil {
		return false
	}
	var fact moduleFact
	return pass.ImportPackageFact(pkg, &fact) && fact.Path == pass.Module.Path
}
