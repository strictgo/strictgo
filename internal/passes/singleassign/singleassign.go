// Package singleassign flags assignments that combine multiple value sources.
package singleassign

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "singleassign",
	Doc:  "require separate assignments for multiple sources, except swaps and rotations",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		var stack []ast.Node
		ast.Inspect(file, func(node ast.Node) bool {
			if node == nil {
				stack = stack[:len(stack)-1]
				return true
			}
			var parent ast.Node
			if len(stack) > 0 {
				parent = stack[len(stack)-1]
			}
			stack = append(stack, node)
			assignment, ok := node.(*ast.AssignStmt)
			if !ok || len(assignment.Rhs) < 2 || len(assignment.Lhs) != len(assignment.Rhs) || permutation(assignment) {
				return true
			}
			diagnostic := analysis.Diagnostic{
				Pos:     assignment.Pos(),
				End:     assignment.End(),
				Message: "use separate assignments for multiple value sources",
			}
			if replacement, ok := split(pass, file, parent, assignment); ok {
				diagnostic.SuggestedFixes = []analysis.SuggestedFix{{
					Message: "Split into separate assignments",
					TextEdits: []analysis.TextEdit{{
						Pos: assignment.Pos(), End: assignment.End(), NewText: []byte(replacement),
					}},
				}}
			}
			pass.Report(diagnostic)
			return true
		})
	}
	return nil, nil
}

// A permutation must use every distinct destination exactly once as a source.
// Calls and receives are excluded: matching text does not mean matching values.
func permutation(assignment *ast.AssignStmt) bool {
	if assignment.Tok != token.ASSIGN {
		return false
	}
	destinations := make(map[string]bool)
	for _, expression := range assignment.Lhs {
		key := types.ExprString(ast.Unparen(expression))
		if key == "_" || destinations[key] || !stable(expression) {
			return false
		}
		destinations[key] = true
	}
	var moved bool
	for index, expression := range assignment.Rhs {
		key := types.ExprString(ast.Unparen(expression))
		if !destinations[key] || !stable(expression) {
			return false
		}
		delete(destinations, key)
		if key != types.ExprString(ast.Unparen(assignment.Lhs[index])) {
			moved = true
		}
	}
	return moved
}

func stable(expression ast.Expr) bool {
	valid := true
	ast.Inspect(expression, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.CallExpr:
			valid = false
		case *ast.UnaryExpr:
			if node.Op == token.ARROW {
				valid = false
			}
		}
		return valid
	})
	return valid
}

func split(pass *analysis.Pass, file *ast.File, parent ast.Node, assignment *ast.AssignStmt) (string, bool) {
	// Header and labeled statements cannot be replaced by several statements.
	switch parent.(type) {
	case *ast.BlockStmt, *ast.CaseClause, *ast.CommClause:
	default:
		return "", false
	}
	if assignment.Tok != token.ASSIGN && assignment.Tok != token.DEFINE {
		return "", false
	}
	for _, group := range file.Comments {
		if group.Pos() < assignment.End() && group.End() > assignment.Pos() {
			return "", false
		}
	}
	names := make(map[string]bool)
	allNew := assignment.Tok == token.DEFINE
	var fields bool
	for _, expression := range assignment.Lhs {
		switch expression := ast.Unparen(expression).(type) {
		case *ast.Ident:
			if expression.Name != "_" {
				names[expression.Name] = true
				if pass.TypesInfo.Defs[expression] == nil {
					allNew = false
				}
			}
		case *ast.SelectorExpr:
			// Only a single field access: nested pointer traversal may panic
			// before the original statement performs any stores.
			if _, ok := ast.Unparen(expression.X).(*ast.Ident); !ok {
				return "", false
			}
			selection := pass.TypesInfo.Selections[expression]
			if selection == nil || selection.Kind() != types.FieldVal || len(selection.Index()) != 1 {
				return "", false
			}
			fields = true
			allNew = false
		default:
			return "", false
		}
	}
	for _, expression := range assignment.Lhs {
		if field, ok := ast.Unparen(expression).(*ast.SelectorExpr); ok && references(field.X, names) {
			return "", false
		}
	}
	for _, expression := range assignment.Rhs {
		if references(expression, names) {
			return "", false
		}
		// Newly declared variables cannot be observed by subsequent sources.
		// Existing destinations need sources that neither call nor panic. Field
		// stores also exclude memory reads that could alias an earlier store.
		if !allNew && !safeSource(pass.TypesInfo, expression, !fields) {
			return "", false
		}
	}
	source, err := pass.ReadFile(pass.Fset.Position(file.Pos()).Filename)
	if err != nil {
		return "", false
	}
	var result strings.Builder
	for index, expression := range assignment.Lhs {
		operator := "="
		if name, ok := expression.(*ast.Ident); ok && assignment.Tok == token.DEFINE && name.Name != "_" && pass.TypesInfo.Defs[name] != nil {
			operator = ":="
		}
		if index > 0 {
			result.WriteByte('\n')
		}
		value := assignment.Rhs[index]
		start := pass.Fset.Position(value.Pos()).Offset
		end := pass.Fset.Position(value.End()).Offset
		result.WriteString(types.ExprString(expression) + " " + operator + " " + string(source[start:end]))
	}
	return result.String(), true
}

// Selector names are not lexical references: x.Field must not be confused with
// a local variable named Field introduced by the split declaration.
func references(expression ast.Expr, names map[string]bool) bool {
	var found bool
	ast.Inspect(expression, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.SelectorExpr:
			if references(node.X, names) {
				found = true
			}
			return false
		case *ast.Ident:
			if names[node.Name] {
				found = true
			}
		}
		return !found
	})
	return found
}

func safeSource(info *types.Info, expression ast.Expr, memoryReads bool) bool {
	expression = ast.Unparen(expression)
	if info.Types[expression].Value != nil {
		return true
	}
	switch expression := expression.(type) {
	case *ast.Ident:
		return memoryReads || !aggregate(info.TypeOf(expression))
	case *ast.BasicLit:
		return true
	case *ast.SelectorExpr:
		selection := info.Selections[expression]
		if selection == nil {
			return memoryReads || !aggregate(info.TypeOf(expression))
		}
		return memoryReads && selection.Kind() == types.FieldVal && !selection.Indirect() && safeSource(info, expression.X, memoryReads)
	case *ast.UnaryExpr:
		if expression.Op == token.AND {
			_, ok := ast.Unparen(expression.X).(*ast.Ident)
			return ok
		}
	}
	return false
}

func aggregate(typ types.Type) bool {
	switch typ.Underlying().(type) {
	case *types.Struct, *types.Array:
		return true
	}
	_, parameter := typ.(*types.TypeParam)
	return parameter
}
