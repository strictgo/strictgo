// Package scopesprawl estimates how many distinct concepts a file contains.
package scopesprawl

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const defaultLimit = 5

var maxGroups = defaultLimit

var Analyzer = &analysis.Analyzer{
	Name: "scopesprawl",
	Doc:  "limit declaration groups per file using package-wide ownership",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	if maxGroups < 1 {
		return nil, fmt.Errorf("scopesprawl.limit must be at least 1")
	}
	graph := inferOwnership(pass)
	for _, file := range pass.Files {
		if ast.IsGenerated(file) || testFile(pass, file) {
			continue
		}
		names := fileGroups(graph, file)
		if len(names) > maxGroups {
			pass.Reportf(file.Name.Pos(), "%d declaration groups (%s); limit %d", len(names), strings.Join(names, ", "), maxGroups)
		}
	}
	return nil, nil
}

func init() {
	Analyzer.Flags.IntVar(&maxGroups, "limit", defaultLimit, "maximum declaration groups per non-test file")
}

func inferOwnership(pass *analysis.Pass) *ownership {
	graph := buildOwnership(pass)
	foldPrivateDeclarations(graph)
	foldOptionConfigurations(graph)
	foldPrivateDeclarations(graph)
	return graph
}

// A shared dependency is not a reason to merge its consumers. Fold only a
// private helper/data/state group with exactly one distinct incoming owner group.
// Recompute after every fold so chains of exclusive helpers can follow owners.
func foldPrivateDeclarations(graph *ownership) {
	for {
		consumers := consumerGroups(graph)
		var changed bool
		for index, declaration := range graph.declarations {
			if root(graph, index) != index || declaration.anchor || len(consumers[index]) != 1 {
				continue
			}
			for owner := range consumers[index] {
				graph.parents[index] = owner
			}
			changed = true
			break
		}
		if !changed {
			return
		}
	}
}

func consumerGroups(graph *ownership) map[int]map[int]bool {
	consumers := make(map[int]map[int]bool)
	for source, targets := range graph.references {
		owner := root(graph, source)
		for target := range targets {
			dependency := root(graph, target)
			if dependency == owner {
				continue
			}
			if consumers[dependency] == nil {
				consumers[dependency] = make(map[int]bool)
			}
			consumers[dependency][owner] = true
		}
	}
	return consumers
}

func root(graph *ownership, index int) int {
	if graph.parents[index] != index {
		graph.parents[index] = root(graph, graph.parents[index])
	}
	return graph.parents[index]
}

func fileGroups(graph *ownership, file *ast.File) []string {
	first := make(map[int]int)
	consumers := consumerGroups(graph)
	for index, declaration := range graph.declarations {
		if declaration.file != file {
			continue
		}
		owner := root(graph, index)
		// Shared private constants are background dependencies, not another
		// concept. Keep their consumers separate; mutable/public state counts.
		if group, ok := graph.declarations[owner].node.(*ast.GenDecl); ok && group.Tok == token.CONST && !graph.declarations[owner].anchor && len(consumers[owner]) > 1 {
			continue
		}
		if previous, ok := first[owner]; !ok || declaration.node.Pos() < graph.declarations[previous].node.Pos() {
			first[owner] = index
		}
	}
	var owners []int
	for owner := range first {
		owners = append(owners, owner)
	}
	slices.SortFunc(owners, func(left, right int) int {
		return cmp.Compare(graph.declarations[first[left]].node.Pos(), graph.declarations[first[right]].node.Pos())
	})
	var names []string
	for _, owner := range owners {
		names = append(names, graph.declarations[owner].name)
	}
	return names
}

func testFile(pass *analysis.Pass, file *ast.File) bool {
	return strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go")
}

// Generic instantiations refer to the same declaration group as their origins.
func declarationIndex(graph *ownership, object types.Object) (int, bool) {
	if function, ok := object.(*types.Func); ok {
		object = function.Origin()
	}
	index, ok := graph.objects[object]
	return index, ok
}
