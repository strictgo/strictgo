package scopesprawl

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

type declaration struct {
	file   *ast.File
	node   ast.Node
	name   string
	anchor bool
}

type ownership struct {
	declarations         []declaration
	objects              map[types.Object]int
	parents              []int
	references           []map[int]bool
	optionConfigurations map[int]map[int]bool
	optionConstructors   map[int]map[int]bool
}

func buildOwnership(pass *analysis.Pass) *ownership {
	graph := &ownership{
		objects:              make(map[types.Object]int),
		optionConfigurations: make(map[int]map[int]bool),
		optionConstructors:   make(map[int]map[int]bool),
	}
	for _, file := range pass.Files {
		// Tests exercise production concepts rather than own them. Including
		// test callers would make ownership vary between ordinary/test builds.
		if testFile(pass, file) {
			continue
		}
		for _, node := range file.Decls {
			switch node := node.(type) {
			case *ast.GenDecl:
				collectDeclarations(pass.TypesInfo, graph, file, node)
			case *ast.FuncDecl:
				function, ok := pass.TypesInfo.Defs[node.Name].(*types.Func)
				if ok {
					addDeclaration(graph, file, node, function.Name(), node.Recv == nil && function.Exported(), []types.Object{function})
				}
			}
		}
	}
	groupOptionTypes(pass.TypesInfo, graph)
	groupMethodsAndConstructors(pass.TypesInfo, graph)
	for index, declaration := range graph.declarations {
		ast.Inspect(declaration.node, func(node ast.Node) bool {
			name, ok := node.(*ast.Ident)
			if !ok {
				return true
			}
			if target, ok := declarationIndex(graph, pass.TypesInfo.Uses[name]); ok {
				graph.references[index][target] = true
			}
			return true
		})
	}
	return graph
}

func collectDeclarations(info *types.Info, graph *ownership, file *ast.File, group *ast.GenDecl) {
	if group.Tok == token.TYPE {
		for _, specification := range group.Specs {
			typ := specification.(*ast.TypeSpec)
			if object := info.Defs[typ.Name]; object != nil && object.Name() != "_" {
				// Private data and protocol types may belong to a sole consumer.
				// Public types and types with methods retain independent identity.
				anchor := object.Exported()
				if named, ok := types.Unalias(object.Type()).(*types.Named); ok {
					anchor = anchor || named.NumMethods() > 0
				}
				addDeclaration(graph, file, typ, object.Name(), anchor, []types.Object{object})
			}
		}
		return
	}
	if group.Tok != token.CONST && group.Tok != token.VAR {
		return
	}
	var objects []types.Object
	var exported bool
	for _, specification := range group.Specs {
		for _, name := range specification.(*ast.ValueSpec).Names {
			object := info.Defs[name]
			if object == nil || name.Name == "_" {
				continue
			}
			objects = append(objects, object)
			exported = exported || object.Exported()
		}
	}
	if len(objects) > 0 {
		// Explicit const/var blocks stay together, labelled by their first name.
		addDeclaration(graph, file, group, objects[0].Name(), exported, objects)
	}
}

func addDeclaration(graph *ownership, file *ast.File, node ast.Node, name string, anchor bool, objects []types.Object) {
	index := len(graph.declarations)
	graph.declarations = append(graph.declarations, declaration{file: file, node: node, name: name, anchor: anchor})
	graph.parents = append(graph.parents, index)
	graph.references = append(graph.references, make(map[int]bool))
	for _, object := range objects {
		if function, ok := object.(*types.Func); ok {
			object = function.Origin()
		}
		graph.objects[object] = index
	}
}

func groupMethodsAndConstructors(info *types.Info, graph *ownership) {
	for index, declaration := range graph.declarations {
		function, ok := declaration.node.(*ast.FuncDecl)
		if !ok {
			continue
		}
		object, ok := info.Defs[function.Name].(*types.Func)
		if !ok {
			continue
		}
		signature := object.Type().(*types.Signature)
		var owner *types.TypeName
		if receiver := signature.Recv(); receiver != nil {
			owner = namedType(receiver.Type())
		} else if target, ok := optionFactoryGroup(graph, signature); ok {
			// Returning the resolved option contract is stronger evidence than
			// a With... name, and also covers aliases and generic factories.
			graph.parents[index] = target
			continue
		} else if constructorName(function.Name.Name) && signature.Results().Len() > 0 {
			owner = namedType(signature.Results().At(0).Type())
		}
		if owner == nil {
			continue
		}
		if target, ok := graph.objects[owner]; ok {
			graph.parents[index] = target
			// Even a private type has its own purpose when it has a constructor
			// or methods. Do not fold domain types into their only caller.
			graph.declarations[target].anchor = true
			if signature.Recv() == nil {
				recordConstructorOptions(graph, signature, index)
			}
		}
	}
}

func namedType(typ types.Type) *types.TypeName {
	typ = types.Unalias(typ)
	if pointer, ok := typ.(*types.Pointer); ok {
		typ = types.Unalias(pointer.Elem())
	}
	if named, ok := typ.(*types.Named); ok {
		return named.Origin().Obj()
	}
	return nil
}

func constructorName(name string) bool {
	for _, prefix := range []string{"New", "new"} {
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		suffix := strings.TrimPrefix(name, prefix)
		if suffix == "" {
			return true
		}
		letter, _ := utf8.DecodeRuneInString(suffix)
		return unicode.IsUpper(letter)
	}
	return false
}
