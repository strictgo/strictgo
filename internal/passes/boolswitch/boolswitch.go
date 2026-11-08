// Package boolswitch flags boolean switches in favor of if statements.
package boolswitch

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "boolswitch",
	Doc:  "replace boolean switches with if statements and early returns",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		var source []byte
		ast.Inspect(file, func(node ast.Node) bool {
			sw, ok := node.(*ast.SwitchStmt)
			if !ok {
				return true
			}
			if sw.Tag != nil {
				typ := pass.TypesInfo.TypeOf(sw.Tag)
				if typ == nil {
					return true
				}
				basic, ok := typ.Underlying().(*types.Basic)
				if !ok || basic.Info()&types.IsBoolean == 0 {
					return true
				}
			}
			diagnostic := analysis.Diagnostic{Pos: sw.Switch, End: sw.End(), Message: "replace boolean switch with if statements and early returns"}
			if source == nil {
				var err error
				source, err = pass.ReadFile(pass.Fset.Position(file.Pos()).Filename)
				if err != nil {
					pass.Report(diagnostic)
					return true
				}
			}
			if replacement, ok := rewrite(pass, sw, source); ok {
				diagnostic.SuggestedFixes = []analysis.SuggestedFix{{Message: "Replace with if statements", TextEdits: []analysis.TextEdit{{Pos: sw.Pos(), End: sw.End(), NewText: replacement}}}}
			}
			pass.Report(diagnostic)
			return true
		})
	}
	return nil, nil
}

func rewrite(pass *analysis.Pass, sw *ast.SwitchStmt, source []byte) ([]byte, bool) {
	safe := true
	ast.Inspect(sw.Body, func(node ast.Node) bool {
		switch node.(type) {
		case *ast.BranchStmt, *ast.LabeledStmt:
			safe = false
		case *ast.FuncLit:
			return false
		}
		return true
	})
	if !safe {
		return nil, false
	}
	text := func(start, end token.Pos) string {
		return string(source[pass.Fset.Position(start).Offset:pass.Fset.Position(end).Offset])
	}
	// Choose a fresh identifier even when a case body introduces its own names.
	name := "strictgoSwitch"
	var suffix int
	for {
		var used bool
		for id := range pass.TypesInfo.Defs {
			if id.Name == name {
				used = true
			}
		}
		for id := range pass.TypesInfo.Uses {
			if id.Name == name {
				used = true
			}
		}
		if !used {
			break
		}
		suffix++
		name = fmt.Sprintf("strictgoSwitch%d", suffix)
	}
	var prefix strings.Builder
	for _, file := range pass.Files {
		for _, group := range file.Comments {
			for _, comment := range group.List {
				if comment.Pos() < sw.Pos() || comment.End() > sw.End() {
					continue
				}
				var inBody bool
				for index, stmt := range sw.Body.List {
					clause := stmt.(*ast.CaseClause)
					end := sw.Body.Rbrace
					if index+1 < len(sw.Body.List) {
						end = sw.Body.List[index+1].Pos()
					}
					if comment.Pos() > clause.Colon && comment.End() <= end {
						inBody = true
					}
				}
				if !inBody {
					prefix.WriteString(comment.Text + "\n")
				}
			}
		}
	}
	if sw.Init != nil {
		prefix.WriteString(text(sw.Init.Pos(), sw.Init.End()) + "\n")
	}
	if sw.Tag != nil {
		prefix.WriteString(name + " := " + text(sw.Tag.Pos(), sw.Tag.End()) + "\n")
	}
	type branch struct {
		condition, body string
		returns         bool
	}
	var branches []branch
	var fallback string
	for index, stmt := range sw.Body.List {
		clause := stmt.(*ast.CaseClause)
		end := sw.Body.Rbrace
		if index+1 < len(sw.Body.List) {
			end = sw.Body.List[index+1].Pos()
		}
		body := strings.TrimSpace(text(clause.Colon+1, end))
		if clause.List == nil {
			fallback = "{\n" + body + "\n}\n"
			continue
		}
		var conditions []string
		for _, expr := range clause.List {
			condition := "(" + text(expr.Pos(), expr.End()) + ")"
			if sw.Tag != nil {
				condition = name + " == " + condition
			}
			conditions = append(conditions, condition)
		}
		branches = append(branches, branch{strings.Join(conditions, " || "), body, returns(clause.Body)})
	}
	if sw.Tag != nil && len(branches) == 0 {
		prefix.WriteString("_ = " + name + "\n")
	}
	result := fallback
	for _, branch := range slices.Backward(branches) {
		head := "if " + branch.condition + " {\n" + branch.body + "\n}"
		if result == "" {
			result = head + "\n"
			continue
		}
		if branch.returns {
			result = head + "\n" + result
			continue
		}
		result = head + " else {\n" + strings.TrimSpace(result) + "\n}\n"
	}
	// A block preserves initializer scope and case-local declarations.
	replacement := "{\n" + strings.TrimSpace(prefix.String()+result) + "\n}"
	formatted, err := format.Source([]byte(replacement))
	return formatted, err == nil
}

func returns(body []ast.Stmt) bool {
	if len(body) == 0 {
		return false
	}
	switch last := body[len(body)-1].(type) {
	case *ast.ReturnStmt:
		return true
	case *ast.BlockStmt:
		return returns(last.List)
	case *ast.IfStmt:
		if last.Else == nil {
			return false
		}
		return returns(last.Body.List) && returns([]ast.Stmt{last.Else})
	}
	return false
}
