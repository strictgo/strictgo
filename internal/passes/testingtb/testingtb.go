// Package testingtb recommends testing.TB for reusable test helpers.
package testingtb

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "testingtb",
	Doc:  "prefer testing.TB for helpers that do not require *testing.T",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	parents := syntaxParents(pass.Files)
	references := make(map[types.Object][]*ast.Ident)
	for name, object := range pass.TypesInfo.Uses {
		references[object] = append(references[object], name)
	}
	for _, file := range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			var function *ast.FuncType
			switch node := node.(type) {
			case *ast.FuncDecl:
				// A method signature may be fixed by an interface contract.
				if node.Recv != nil || node.Body == nil || testEntrypoint(pass, file, node) {
					return true
				}
				if !directCallsOnly(references, parents, pass.TypesInfo.Defs[node.Name]) {
					return true
				}
				function = node.Type
			case *ast.FuncLit:
				if !helperLiteral(pass.TypesInfo, references, parents, node) {
					return true
				}
				function = node.Type
			default:
				return true
			}
			checkParameters(pass, function, references, parents)
			return true
		})
	}
	return nil, nil
}

func checkParameters(pass *analysis.Pass, function *ast.FuncType, references map[types.Object][]*ast.Ident, parents map[ast.Node]ast.Node) {
	for _, field := range function.Params.List {
		for _, name := range field.Names {
			variable, ok := pass.TypesInfo.Defs[name].(*types.Var)
			if !ok {
				continue
			}
			tb := testingTBFor(variable.Type())
			if tb == nil {
				continue
			}
			usage := parameterUse{info: pass.TypesInfo, parents: parents, tb: tb, methods: types.NewMethodSet(tb)}
			var requiresT bool
			for _, reference := range references[variable] {
				if !usage.permitted(reference) {
					requiresT = true
					break
				}
			}
			if !requiresT {
				pass.Reportf(name.Pos(), "helper parameter %q can use testing.TB instead of *testing.T; name it tb", name.Name)
			}
		}
	}
}

func testingTBFor(typ types.Type) *types.Named {
	pointer, ok := types.Unalias(typ).(*types.Pointer)
	if !ok {
		return nil
	}
	named, ok := types.Unalias(pointer.Elem()).(*types.Named)
	if !ok || named.Obj().Name() != "T" || named.Obj().Pkg() == nil || named.Obj().Pkg().Path() != "testing" {
		return nil
	}
	object := named.Obj().Pkg().Scope().Lookup("TB")
	if object == nil {
		return nil
	}
	tb, ok := types.Unalias(object.Type()).(*types.Named)
	if !ok {
		return nil
	}
	return tb
}

func testEntrypoint(pass *analysis.Pass, file *ast.File, function *ast.FuncDecl) bool {
	if !strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go") || !strings.HasPrefix(function.Name.Name, "Test") {
		return false
	}
	if suffix := strings.TrimPrefix(function.Name.Name, "Test"); suffix != "" {
		letter, _ := utf8.DecodeRuneInString(suffix)
		if unicode.IsLower(letter) {
			return false
		}
	}
	object, ok := pass.TypesInfo.Defs[function.Name].(*types.Func)
	if !ok {
		return false
	}
	signature := object.Type().(*types.Signature)
	return signature.Params().Len() == 1 && signature.Results().Len() == 0 && signature.TypeParams().Len() == 0 && !signature.Variadic() && testingTBFor(signature.Params().At(0).Type()) != nil
}

func syntaxParents(files []*ast.File) map[ast.Node]ast.Node {
	parents := make(map[ast.Node]ast.Node)
	for _, file := range files {
		var stack []ast.Node
		ast.Inspect(file, func(node ast.Node) bool {
			if node == nil {
				stack = stack[:len(stack)-1]
				return true
			}
			if len(stack) > 0 {
				parents[node] = stack[len(stack)-1]
			}
			stack = append(stack, node)
			return true
		})
	}
	return parents
}

type expressionContext struct {
	expression ast.Expr
	parent     ast.Node
}

func expressionUse(name ast.Expr, parents map[ast.Node]ast.Node) expressionContext {
	for {
		parent, ok := parents[name].(*ast.ParenExpr)
		if !ok {
			return expressionContext{expression: name, parent: parents[name]}
		}
		name = parent
	}
}

// Widening parameters is not valid for function values assigned to a concrete
// callback type. Keep those, including named t.Run callbacks, unchanged.
func directCallsOnly(references map[types.Object][]*ast.Ident, parents map[ast.Node]ast.Node, function types.Object) bool {
	if function == nil {
		return false
	}
	for _, name := range references[function] {
		use := expressionUse(name, parents)
		switch instance := use.parent.(type) {
		case *ast.IndexExpr:
			if instance.X == use.expression {
				use = expressionUse(instance, parents)
			}
		case *ast.IndexListExpr:
			if instance.X == use.expression {
				use = expressionUse(instance, parents)
			}
		}
		call, ok := use.parent.(*ast.CallExpr)
		if !ok || call.Fun != use.expression {
			return false
		}
	}
	return len(references[function]) > 0
}

func helperLiteral(info *types.Info, references map[types.Object][]*ast.Ident, parents map[ast.Node]ast.Node, literal *ast.FuncLit) bool {
	use := expressionUse(literal, parents)
	switch parent := use.parent.(type) {
	case *ast.AssignStmt:
		if parent.Tok == token.DEFINE && len(parent.Lhs) == 1 && len(parent.Rhs) == 1 {
			if name, ok := parent.Lhs[0].(*ast.Ident); ok {
				return directCallsOnly(references, parents, info.Defs[name])
			}
		}
	case *ast.ValueSpec:
		if parent.Type == nil && len(parent.Names) == 1 && len(parent.Values) == 1 {
			return directCallsOnly(references, parents, info.Defs[parent.Names[0]])
		}
	case *ast.CallExpr:
		return parent.Fun == use.expression
	}
	return false
}
