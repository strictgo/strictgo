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
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Type.Results == nil {
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
