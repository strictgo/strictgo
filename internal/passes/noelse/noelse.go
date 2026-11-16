// Package noelse flags else branches made unnecessary by an early return.
package noelse

import (
	"bytes"
	"go/ast"
	"go/scanner"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noelse",
	Doc:  "remove unnecessary else branches after returning branches",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		checkFile(pass, file)
	}
	return nil, nil
}

func checkFile(pass *analysis.Pass, file *ast.File) {
	source, err := pass.ReadFile(pass.Fset.File(file.Pos()).Name())
	if err != nil {
		source = nil
	}
	introduced := make(map[*types.Scope]map[string]bool)
	var flattened []ast.Node
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
		conditional, ok := node.(*ast.IfStmt)
		if !ok || conditional.Else == nil || !returns(conditional.Body) {
			return true
		}
		position := elsePosition(pass.Fset, conditional, source)
		diagnostic := analysis.Diagnostic{Pos: conditional.Else.Pos(), Message: "unnecessary else after return"}
		if position.IsValid() {
			diagnostic.Pos = position
			diagnostic.End = position + token.Pos(len("else"))
		}
		if position.IsValid() && canFlatten(pass, file, conditional, parent, introduced, flattened) {
			edits := []analysis.TextEdit{{Pos: position, End: position + token.Pos(len("else")), NewText: []byte("\n")}}
			if block, ok := conditional.Else.(*ast.BlockStmt); ok {
				tokenFile := pass.Fset.File(block.Pos())
				edits = append(edits, braceEdit(tokenFile, source, block.Lbrace), braceEdit(tokenFile, source, block.Rbrace))
				reserveNames(pass.TypesInfo, conditional, block, introduced)
			}
			diagnostic.SuggestedFixes = []analysis.SuggestedFix{{Message: "Remove unnecessary else", TextEdits: edits}}
			flattened = append(flattened, conditional.Else)
		}
		pass.Report(diagnostic)
		return true
	})
}

func returns(block *ast.BlockStmt) bool {
	if block == nil || len(block.List) == 0 {
		return false
	}
	switch last := block.List[len(block.List)-1].(type) {
	case *ast.ReturnStmt:
		return true
	case *ast.BlockStmt:
		return returns(last)
	case *ast.IfStmt:
		return returns(last.Body) && elseReturns(last.Else)
	}
	return false
}

func elseReturns(statement ast.Stmt) bool {
	switch statement := statement.(type) {
	case *ast.BlockStmt:
		return returns(statement)
	case *ast.IfStmt:
		return returns(statement.Body) && elseReturns(statement.Else)
	}
	return false
}

func canFlatten(pass *analysis.Pass, file *ast.File, conditional *ast.IfStmt, parent ast.Node, introduced map[*types.Scope]map[string]bool, flattened []ast.Node) bool {
	if conditional.Init != nil {
		return false
	}
	// Flattening a nested else-if on its own would move code outside the outer
	// conditional. Let a later analyzer run handle it after the outer fix.
	if outer, ok := parent.(*ast.IfStmt); ok && outer.Else == conditional {
		return false
	}
	// Do not combine nested scope-changing fixes. Each proposed fix must remain
	// safe when applied alone as well as alongside the other suggestions.
	for _, scope := range flattened {
		if conditional.Pos() >= scope.Pos() && conditional.End() <= scope.End() {
			return false
		}
	}
	block, ok := conditional.Else.(*ast.BlockStmt)
	if !ok {
		return true
	}
	scope := pass.TypesInfo.Scopes[block]
	ifScope := pass.TypesInfo.Scopes[conditional]
	if scope == nil || ifScope == nil || ifScope.Parent() == nil {
		return false
	}
	parentScope := ifScope.Parent()
	if len(scope.Names()) > 0 && hasGoto(file) {
		return false
	}
	for _, name := range scope.Names() {
		if parentScope.Lookup(name) != nil || introduced[parentScope][name] {
			return false
		}
		// A new binding must not change how subsequent code resolves names from
		// outer scopes, including captures in subsequently declared closures.
		for identifier := range pass.TypesInfo.Uses {
			if identifier.Name == name && identifier.Pos() > block.End() && identifier.Pos() < parentScope.End() {
				return false
			}
		}
	}
	return true
}

func reserveNames(info *types.Info, conditional *ast.IfStmt, block *ast.BlockStmt, introduced map[*types.Scope]map[string]bool) {
	parent := info.Scopes[conditional].Parent()
	if introduced[parent] == nil {
		introduced[parent] = make(map[string]bool)
	}
	for _, name := range info.Scopes[block].Names() {
		introduced[parent][name] = true
	}
}

func hasGoto(file *ast.File) bool {
	var found bool
	ast.Inspect(file, func(node ast.Node) bool {
		if branch, ok := node.(*ast.BranchStmt); ok && branch.Tok == token.GOTO {
			found = true
		}
		return !found
	})
	return found
}

// Remove an otherwise empty brace line too, avoiding spurious blank lines.
// Inline braces and braces with comments are removed without touching comments.
func braceEdit(file *token.File, source []byte, brace token.Pos) analysis.TextEdit {
	offset := file.Offset(brace)
	start := bytes.LastIndexByte(source[:offset], '\n') + 1
	end := len(source)
	if newline := bytes.IndexByte(source[offset+1:], '\n'); newline >= 0 {
		end = offset + 1 + newline
	}
	if len(bytes.TrimSpace(source[start:offset])) == 0 && len(bytes.TrimSpace(source[offset+1:end])) == 0 {
		if end < len(source) {
			end++
		}
		return analysis.TextEdit{Pos: file.Pos(start), End: file.Pos(end)}
	}
	return analysis.TextEdit{Pos: brace, End: brace + 1}
}

// Scan just the gap between branches so comments containing "else" are not
// mistaken for the keyword. All edits preserve surrounding comments verbatim.
func elsePosition(files *token.FileSet, conditional *ast.IfStmt, source []byte) token.Pos {
	if source == nil {
		return token.NoPos
	}
	file := files.File(conditional.Pos())
	start := file.Offset(conditional.Body.End())
	end := file.Offset(conditional.Else.Pos())
	fragment := source[start:end]
	fragmentFile := token.NewFileSet().AddFile("", -1, len(fragment))
	var lexer scanner.Scanner
	lexer.Init(fragmentFile, fragment, nil, 0)
	for {
		position, kind, _ := lexer.Scan()
		if kind == token.EOF {
			return token.NoPos
		}
		if kind == token.ELSE {
			return file.Pos(start + fragmentFile.Offset(position))
		}
	}
}
