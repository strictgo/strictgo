package testingtb

import (
	"go/ast"
	"go/types"
)

type parameterUse struct {
	info    *types.Info
	parents map[ast.Node]ast.Node
	tb      *types.Named
	methods *types.MethodSet
}

func (usage parameterUse) permitted(name *ast.Ident) bool {
	use := expressionUse(name, usage.parents)
	expression := use.expression
	switch parent := use.parent.(type) {
	case *ast.SelectorExpr:
		selection := usage.info.Selections[parent]
		if parent.X != expression || selection == nil || selection.Kind() != types.MethodVal {
			return false
		}
		method := usage.methods.Lookup(nil, parent.Sel.Name)
		return method != nil && types.Identical(method.Obj().Type(), selection.Obj().Type())
	case *ast.CallExpr:
		return usage.callArgument(parent, expression)
	case *ast.ReturnStmt:
		return usage.returnValue(parent, expression)
	case *ast.AssignStmt:
		if len(parent.Lhs) != len(parent.Rhs) {
			return false
		}
		for index, value := range parent.Rhs {
			if value == expression {
				if blank, ok := parent.Lhs[index].(*ast.Ident); ok && blank.Name == "_" {
					return true
				}
				return usage.assignableTo(usage.info.TypeOf(parent.Lhs[index]))
			}
		}
	case *ast.ValueSpec:
		if len(parent.Names) != len(parent.Values) {
			return false
		}
		for index, value := range parent.Values {
			if value == expression {
				if parent.Names[index].Name == "_" {
					return true
				}
				return usage.assignableTo(usage.info.Defs[parent.Names[index]].Type())
			}
		}
	}
	// In particular, nil comparisons, dereferences, address-taking, boxing
	// conversions/type assertions, and inferred *T aliases are not proven safe.
	return false
}

func (usage parameterUse) callArgument(call *ast.CallExpr, expression ast.Expr) bool {
	if usage.info.Types[ast.Unparen(call.Fun)].IsType() {
		return false
	}
	typ := usage.info.TypeOf(call.Fun)
	if typ == nil {
		return false
	}
	signature, ok := typ.Underlying().(*types.Signature)
	if !ok {
		return false
	}
	for index, argument := range call.Args {
		if argument != expression {
			continue
		}
		parameters := signature.Params()
		parameterIndex := index
		if signature.Variadic() && index >= parameters.Len()-1 {
			parameterIndex = parameters.Len() - 1
		}
		if parameterIndex < 0 || parameterIndex >= parameters.Len() {
			return false
		}
		parameter := parameters.At(parameterIndex).Type()
		if signature.Variadic() && index >= parameters.Len()-1 && !call.Ellipsis.IsValid() {
			parameter = parameter.(*types.Slice).Elem()
		}
		return usage.assignableTo(parameter)
	}
	return false
}

func (usage parameterUse) returnValue(returned *ast.ReturnStmt, expression ast.Expr) bool {
	signature := enclosingSignature(usage.info, usage.parents, returned)
	if signature == nil || signature.Results().Len() != len(returned.Results) {
		return false
	}
	for index, value := range returned.Results {
		if value == expression {
			return usage.assignableTo(signature.Results().At(index).Type())
		}
	}
	return false
}

func (usage parameterUse) assignableTo(typ types.Type) bool {
	return typ != nil && types.AssignableTo(usage.tb, typ)
}

func enclosingSignature(info *types.Info, parents map[ast.Node]ast.Node, node ast.Node) *types.Signature {
	for node != nil {
		switch function := node.(type) {
		case *ast.FuncDecl:
			if object, ok := info.Defs[function.Name].(*types.Func); ok {
				return object.Type().(*types.Signature)
			}
		case *ast.FuncLit:
			if signature, ok := info.TypeOf(function).(*types.Signature); ok {
				return signature
			}
		}
		node = parents[node]
	}
	return nil
}
