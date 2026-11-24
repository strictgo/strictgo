package scopesprawl

import (
	"go/ast"
	"go/types"
	"strings"
)

// A functional option is a named ...Option function taking one pointer to a
// struct (or generic target), returning nothing or error. A With... spelling
// alone is insufficient, and ordinary callback/result types are not families.
func optionContract(typ types.Type) (*types.TypeName, *types.TypeName) {
	named, ok := types.Unalias(typ).(*types.Named)
	if !ok || !strings.HasSuffix(strings.ToLower(named.Obj().Name()), "option") {
		return nil, nil
	}
	signature, ok := named.Underlying().(*types.Signature)
	if !ok || signature.Variadic() || signature.Params().Len() != 1 {
		return nil, nil
	}
	results := signature.Results()
	if results.Len() > 1 || (results.Len() == 1 && !types.Identical(results.At(0).Type(), types.Universe.Lookup("error").Type())) {
		return nil, nil
	}
	pointer, ok := types.Unalias(signature.Params().At(0).Type()).(*types.Pointer)
	if !ok {
		return nil, nil
	}
	target := types.Unalias(pointer.Elem())
	if _, parameter := target.(*types.TypeParam); !parameter {
		if _, structure := target.Underlying().(*types.Struct); !structure {
			return nil, nil
		}
	}
	return named.Origin().Obj(), namedType(target)
}

func groupOptionTypes(info *types.Info, graph *ownership) {
	for index, declaration := range graph.declarations {
		specification, ok := declaration.node.(*ast.TypeSpec)
		if !ok {
			continue
		}
		object := info.Defs[specification.Name]
		if object == nil {
			continue
		}
		if target, ok := recordOptionContract(graph, object.Type()); ok {
			// An alias denotes the same option family, not another concept.
			graph.parents[index] = target
		}
	}
}

func recordOptionContract(graph *ownership, typ types.Type) (int, bool) {
	option, configuration := optionContract(typ)
	if option == nil {
		return 0, false
	}
	target, ok := graph.objects[option]
	if !ok {
		return 0, false
	}
	graph.declarations[target].anchor = true
	if config, ok := graph.objects[configuration]; ok {
		if graph.optionConfigurations[config] == nil {
			graph.optionConfigurations[config] = make(map[int]bool)
		}
		graph.optionConfigurations[config][target] = true
	}
	return target, true
}

func optionFactoryGroup(graph *ownership, signature *types.Signature) (int, bool) {
	results := signature.Results()
	if results.Len() == 0 || results.Len() > 2 || (results.Len() == 2 && !types.Identical(results.At(1).Type(), types.Universe.Lookup("error").Type())) {
		return 0, false
	}
	return recordOptionContract(graph, results.At(0).Type())
}

func recordConstructorOptions(graph *ownership, signature *types.Signature, constructor int) {
	for parameter := range signature.Params().Variables() {
		typ := parameter.Type()
		if slice, ok := types.Unalias(typ).Underlying().(*types.Slice); ok {
			typ = slice.Elem()
		}
		if option, ok := recordOptionContract(graph, typ); ok {
			if graph.optionConstructors[option] == nil {
				graph.optionConstructors[option] = make(map[int]bool)
			}
			graph.optionConstructors[option][constructor] = true
		}
	}
}

// The config belongs to its option contract even when a constructor accepting
// that contract also instantiates it. Other consumers, independent behavior,
// exported configs, and configs shared by different contracts remain separate.
func foldOptionConfigurations(graph *ownership) {
	for index := range graph.declarations {
		contracts := graph.optionConfigurations[index]
		if len(contracts) == 0 {
			continue
		}
		config := root(graph, index)
		if graph.declarations[config].anchor {
			continue
		}
		owners := make(map[int]bool)
		for contract := range contracts {
			owners[root(graph, contract)] = true
		}
		if len(owners) != 1 {
			continue
		}
		for owner := range owners {
			if config != owner && optionOwnsConfiguration(graph, config, owner) {
				graph.parents[config] = owner
			}
		}
	}
}

func optionOwnsConfiguration(graph *ownership, config, option int) bool {
	allowed := map[int]bool{option: true}
	for contract, constructors := range graph.optionConstructors {
		if root(graph, contract) != option {
			continue
		}
		for constructor := range constructors {
			allowed[root(graph, constructor)] = true
		}
	}
	consumers := consumerGroups(graph)
	if !consumers[config][option] {
		return false
	}
	for consumer := range consumers[config] {
		if !allowed[consumer] {
			return false
		}
	}
	return true
}
