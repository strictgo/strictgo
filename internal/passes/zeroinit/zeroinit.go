// Package zeroinit flags redundant zero-value initializers for basic types.
package zeroinit

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "zeroinit",
	Doc:  "use var declarations instead of explicit basic-type zero initializers",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		// These statement positions cannot be replaced with a var declaration.
		noFix := make(map[ast.Stmt]bool)
		ast.Inspect(file, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.IfStmt:
				noFix[node.Init] = true
			case *ast.ForStmt:
				noFix[node.Init] = true
				noFix[node.Post] = true
			case *ast.SwitchStmt:
				noFix[node.Init] = true
			case *ast.TypeSwitchStmt:
				noFix[node.Init] = true
			case *ast.LabeledStmt:
				noFix[node.Stmt] = true
			}
			return true
		})

		qualifier := func(pkg *types.Package) string {
			if pkg == pass.Pkg {
				return ""
			}
			for _, imp := range file.Imports {
				path, _ := strconv.Unquote(imp.Path.Value)
				if path == pkg.Path() {
					if imp.Name != nil {
						if imp.Name.Name == "." {
							return ""
						}
						return imp.Name.Name
					}
					return pkg.Name()
				}
			}
			return pkg.Name()
		}

		report := func(name *ast.Ident, value ast.Expr, start, end token.Pos, fix bool, prefix string) {
			obj := pass.TypesInfo.Defs[name]
			if name.Name == "_" || obj == nil || !isZero(pass.TypesInfo.Types[value]) {
				return
			}
			if _, ok := obj.Type().Underlying().(*types.Basic); !ok {
				return
			}
			diagnostic := analysis.Diagnostic{
				Pos: value.Pos(), End: value.End(),
				Message: "omit explicit zero-value initializer for " + name.Name + "; use a var declaration",
			}
			typeName := types.TypeString(obj.Type(), qualifier)
			// A predeclared type or import qualifier may be shadowed here.
			resolved, err := types.Eval(pass.Fset, pass.Pkg, start, typeName)
			if err != nil {
				pass.Report(diagnostic)
				return
			}
			if fix && resolved.IsType() && types.Identical(resolved.Type, obj.Type()) {
				replacement := prefix + name.Name + " " + typeName
				diagnostic.SuggestedFixes = []analysis.SuggestedFix{{
					Message:   "Use a var declaration",
					TextEdits: []analysis.TextEdit{{Pos: start, End: end, NewText: []byte(replacement)}},
				}}
			}
			pass.Report(diagnostic)
		}

		ast.Inspect(file, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.AssignStmt:
				if node.Tok != token.DEFINE || len(node.Lhs) != len(node.Rhs) {
					break
				}
				for i, lhs := range node.Lhs {
					if name, ok := lhs.(*ast.Ident); ok {
						report(name, node.Rhs[i], node.Pos(), node.End(), len(node.Lhs) == 1 && !noFix[node], "var ")
					}
				}
			case *ast.GenDecl:
				if node.Tok != token.VAR {
					break
				}
				for _, spec := range node.Specs {
					value := spec.(*ast.ValueSpec)
					if len(value.Names) != len(value.Values) {
						continue
					}
					for i, name := range value.Names {
						report(name, value.Values[i], value.Pos(), value.End(), len(value.Names) == 1, "")
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func isZero(value types.TypeAndValue) bool {
	if value.Value == nil {
		return false
	}
	switch value.Value.Kind() {
	case constant.Bool:
		return !constant.BoolVal(value.Value)
	case constant.String:
		return constant.StringVal(value.Value) == ""
	case constant.Int, constant.Float, constant.Complex:
		return constant.Compare(value.Value, token.EQL, constant.MakeInt64(0))
	}
	return false
}
