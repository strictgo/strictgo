// Package singleusestring flags single-use unexported, untyped string constants.
package singleusestring

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "singleusestring",
	Doc:  "inline single-use unexported, untyped, single-line string constants",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	references := make(map[types.Object][]*ast.Ident)
	for name, object := range pass.TypesInfo.Uses {
		references[object] = append(references[object], name)
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			declaration, ok := node.(*ast.GenDecl)
			if !ok || declaration.Tok != token.CONST || declaration.Lparen.IsValid() || len(declaration.Specs) != 1 {
				return true
			}
			spec := declaration.Specs[0].(*ast.ValueSpec)
			if len(spec.Names) != 1 || len(spec.Values) != 1 || spec.Names[0].Name == "_" || ast.IsExported(spec.Names[0].Name) {
				return true
			}
			object, ok := pass.TypesInfo.Defs[spec.Names[0]].(*types.Const)
			if !ok || !types.Identical(object.Type(), types.Typ[types.UntypedString]) || len(references[object]) != 1 {
				return true
			}
			value := constant.StringVal(object.Val())
			if strings.ContainsAny(value, "\r\n") || pass.Fset.File(spec.Pos()).Line(spec.Values[0].Pos()) != pass.Fset.File(spec.Pos()).Line(spec.Values[0].End()) {
				return true
			}
			diagnostic := analysis.Diagnostic{
				Pos: spec.Names[0].Pos(), End: spec.Names[0].End(),
				Message: "string constant " + strconv.Quote(object.Name()) + " is used only once; inline its value",
			}
			// Literal-only fixes cannot delete another candidate's reference edit.
			// Nonliteral initializers can still be reported without overlapping fixes.
			literal, ok := ast.Unparen(spec.Values[0]).(*ast.BasicLit)
			if ok && literal.Kind == token.STRING && !commented(pass.Fset, file, declaration, spec) {
				reference := references[object][0]
				diagnostic.SuggestedFixes = []analysis.SuggestedFix{{
					Message: "Inline string constant",
					TextEdits: []analysis.TextEdit{
						{Pos: declaration.Pos(), End: declaration.End()},
						{Pos: reference.Pos(), End: reference.End(), NewText: []byte(strconv.Quote(value))},
					},
				}}
			}
			pass.Report(diagnostic)
			return true
		})
	}
	return nil, nil
}

func commented(files *token.FileSet, file *ast.File, declaration *ast.GenDecl, spec *ast.ValueSpec) bool {
	if declaration.Doc != nil || spec.Doc != nil || spec.Comment != nil {
		return true
	}
	physical := files.File(declaration.Pos())
	start := physical.Line(declaration.Pos())
	end := physical.Line(declaration.End())
	for _, comment := range file.Comments {
		if physical.Line(comment.End()) >= start && physical.Line(comment.Pos()) <= end {
			return true
		}
	}
	return false
}
