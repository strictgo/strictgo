package badname

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

// Recognize named functional-option contracts rather than exempting every
// callback or every struct called Options. Aliases retain the contract identity.
func optionTarget(typ types.Type) *types.Pointer {
	named, ok := types.Unalias(typ).(*types.Named)
	if !ok || !strings.HasSuffix(named.Obj().Name(), "Option") {
		return nil
	}
	signature, ok := named.Underlying().(*types.Signature)
	if !ok || !optionSignature(signature) {
		return nil
	}
	pointer, ok := types.Unalias(signature.Params().At(0).Type()).(*types.Pointer)
	if !ok {
		return nil
	}
	if _, ok := pointer.Elem().Underlying().(*types.Struct); !ok {
		return nil
	}
	return pointer
}

func optionSignature(signature *types.Signature) bool {
	if signature.Variadic() || signature.Params().Len() != 1 {
		return false
	}
	results := signature.Results()
	return results.Len() == 0 || (results.Len() == 1 && types.Identical(results.At(0).Type(), types.Universe.Lookup("error").Type()))
}

func exemptFunctionalOptions(pass *analysis.Pass, tree *inspector.Inspector, exempt map[types.Object]bool) {
	info := pass.TypesInfo
	var targets []*types.Pointer
	for _, value := range info.Types {
		if target := optionTarget(value.Type); target != nil {
			targets = append(targets, target)
		}
	}
	// Run before checking declarations: an option application may appear after
	// the configuration variable was declared.
	tree.Preorder([]ast.Node{(*ast.Ident)(nil), (*ast.CallExpr)(nil), (*ast.FuncLit)(nil), (*ast.FuncDecl)(nil)}, func(node ast.Node) {
		switch node := node.(type) {
		case *ast.Ident:
			variable, ok := info.Defs[node].(*types.Var)
			if ok && optionTarget(variable.Type()) != nil {
				exempt[variable] = true
			}
		case *ast.CallExpr:
			if optionTarget(info.TypeOf(node.Fun)) == nil || len(node.Args) != 1 {
				return
			}
			argument := ast.Unparen(node.Args[0])
			if address, ok := argument.(*ast.UnaryExpr); ok && address.Op == token.AND {
				argument = ast.Unparen(address.X)
			}
			if name, ok := argument.(*ast.Ident); ok {
				exempt[info.ObjectOf(name)] = true
			}
		case *ast.FuncLit:
			signature, ok := info.TypeOf(node).(*types.Signature)
			if ok {
				exemptOptionParameter(signature, targets, exempt)
			}
		case *ast.FuncDecl:
			function, ok := info.Defs[node.Name].(*types.Func)
			if ok {
				exemptOptionParameter(function.Type().(*types.Signature), targets, exempt)
			}
		}
	})
}

func exemptOptionParameter(signature *types.Signature, targets []*types.Pointer, exempt map[types.Object]bool) {
	if !optionSignature(signature) {
		return
	}
	parameter := signature.Params().At(0)
	for _, target := range targets {
		if types.Identical(parameter.Type(), target) {
			exempt[parameter] = true
			return
		}
	}
}
