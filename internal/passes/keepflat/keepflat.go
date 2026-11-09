// Package keepflat defines an Analyzer that flags deeply nested blocks.
package keepflat

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const maxDepth = 5

var Analyzer = &analysis.Analyzer{
	Name: "keepflat",
	Doc:  "flags code nested more than five blocks deep within a function",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Walk(visitor{pass: pass}, file)
	}
	return nil, nil
}

// A value receiver keeps each child's depth local to its branch of the AST.
type visitor struct {
	pass         *analysis.Pass
	depth        int
	functionBody *ast.BlockStmt
}

func (v visitor) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case nil:
		return nil
	case *ast.FuncDecl:
		v.depth = 0
		v.functionBody = node.Body
	case *ast.FuncLit:
		v.depth = 0
		v.functionBody = node.Body
	case *ast.BlockStmt:
		if node != v.functionBody {
			v.depth++
			if v.depth == maxDepth+1 {
				v.pass.Reportf(node.Lbrace, "nesting depth %d exceeds maximum of %d", v.depth, maxDepth)
			}
		}
	}
	return v
}
