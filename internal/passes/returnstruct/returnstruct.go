package returnstruct

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "returnstructs",
	Doc:  "finds functions and methods that return interfaces",
	Run:  run,
}

var (
	errorType = types.Universe.Lookup("error").Type()
	anyType   = types.Universe.Lookup("any").Type()
)

func run(pass *analysis.Pass) (any, error) {
	interfaces := externalInterfaces(pass.Pkg)
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Type.Results == nil {
				continue
			}

			method, ok := pass.TypesInfo.Defs[fn.Name].(*types.Func)
			if !ok {
				continue
			}
			if implementsExternalMethod(method, interfaces) {
				continue
			}

			for _, result := range fn.Type.Results.List {
				typ := pass.TypesInfo.TypeOf(result.Type)
				if typ == nil {
					continue
				}

				if _, ok := typ.Underlying().(*types.Interface); ok {
					if types.Identical(typ, errorType) {
						continue
					}
					if types.Identical(typ, anyType) {
						continue
					}
					pass.Reportf(
						result.Type.Pos(),
						"%s returns an interface",
						fn.Name.Name,
					)
				}
			}
		}
	}

	return nil, nil
}

// externalInterfaces collects exported interfaces from dependencies, including
// transitive imports so callers need not import an interface's defining package.
func externalInterfaces(pkg *types.Package) []*types.Interface {
	var interfaces []*types.Interface
	seen := make(map[*types.Package]bool)
	var visit func(*types.Package)
	visit = func(dependency *types.Package) {
		if seen[dependency] {
			return
		}
		seen[dependency] = true
		scope := dependency.Scope()
		for _, name := range scope.Names() {
			obj, ok := scope.Lookup(name).(*types.TypeName)
			if !ok || !obj.Exported() {
				continue
			}
			iface, ok := obj.Type().Underlying().(*types.Interface)
			if !ok || !iface.IsMethodSet() || iface.NumMethods() == 0 {
				continue
			}
			interfaces = append(interfaces, iface)
		}
		for _, imported := range dependency.Imports() {
			visit(imported)
		}
	}
	for _, imported := range pkg.Imports() {
		visit(imported)
	}
	return interfaces
}

// A matching method name alone is insufficient: the receiver must implement
// the complete interface, and this method must belong to that contract.
func implementsExternalMethod(method *types.Func, interfaces []*types.Interface) bool {
	if method == nil {
		return false
	}
	signature, ok := method.Type().(*types.Signature)
	if !ok || signature.Recv() == nil {
		return false
	}
	receiver := signature.Recv().Type()
	candidates := []types.Type{receiver}
	if _, pointer := receiver.(*types.Pointer); !pointer {
		candidates = append(candidates, types.NewPointer(receiver))
	}
	for _, iface := range interfaces {
		var required bool
		for i := 0; i < iface.NumMethods(); i++ {
			if iface.Method(i).Id() == method.Id() {
				required = true
				break
			}
		}
		if !required {
			continue
		}
		for _, candidate := range candidates {
			if types.Implements(candidate, iface) {
				return true
			}
		}
	}
	return false
}
