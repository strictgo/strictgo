package noelse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestSuggestedFixes(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), Analyzer, "a")
}

func TestElsePosition(t *testing.T) {
	source := []byte("package test; func run() { if true { return } /* else in a comment */ else {} }")
	files := token.NewFileSet()
	file, err := parser.ParseFile(files, "test.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	function := file.Decls[0].(*ast.FuncDecl)
	conditional := function.Body.List[0].(*ast.IfStmt)
	position := elsePosition(files, conditional, source)
	if !position.IsValid() {
		t.Fatal("else keyword was not found")
	}
	got := files.File(position).Offset(position)
	want := strings.LastIndex(string(source), "else")
	if got != want {
		t.Errorf("else offset = %d, want %d", got, want)
	}
	if elsePosition(files, conditional, nil).IsValid() {
		t.Error("cannot locate the keyword without source")
	}
}

func TestUnsafeFixes(t *testing.T) {
	results := analysistest.Run(t, analysistest.TestData(), Analyzer, "nofix")
	for _, result := range results {
		for _, diagnostic := range result.Action.Diagnostics {
			if len(diagnostic.SuggestedFixes) != 0 {
				t.Errorf("unexpected fix for %s", diagnostic.Message)
			}
		}
	}
}
