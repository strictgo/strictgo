package reportingrules

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"reflect"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestGeneratedReportsPreserveAnalysis(t *testing.T) {
	source := parseInputs(t)
	dependency := &analysis.Analyzer{Name: "dependency"}
	result := &sampleFact{}
	analyzer := &analysis.Analyzer{
		Name: "sample", Requires: []*analysis.Analyzer{dependency},
		ResultType: reflect.TypeOf(result), FactTypes: []analysis.Fact{result},
	}
	analyzer.Flags.Bool("sample", false, "sample flag")
	var reported []analysis.Diagnostic
	var exported analysis.Fact
	pass := &analysis.Pass{
		Analyzer: analyzer, Fset: source.fset, Files: source.files,
		Pkg: types.NewPackage("example.com/sample", "sample"), TypesInfo: &types.Info{},
		ResultOf:          map[*analysis.Analyzer]any{dependency: result},
		ExportPackageFact: func(fact analysis.Fact) { exported = fact },
		Report:            func(diagnostic analysis.Diagnostic) { reported = append(reported, diagnostic) },
	}
	analyzer.Run = func(filtered *analysis.Pass) (any, error) {
		if filtered == pass || filtered.Fset != pass.Fset || filtered.Pkg != pass.Pkg || filtered.TypesInfo != pass.TypesInfo || filtered.Analyzer != analyzer {
			t.Error("only the reporting callback should change")
		}
		if len(filtered.Files) != 2 || filtered.Files[1] != source.files[1] || filtered.ResultOf[dependency] != result {
			t.Error("generated analysis input or dependency results were lost")
		}
		filtered.ExportPackageFact(result)
		filtered.Reportf(source.files[0].Package, "handwritten")
		filtered.Reportf(source.files[1].Package, "generated")
		return result, nil
	}
	registered := IgnoreGenerated(analyzer, dependency)
	if registered[0] != analyzer || analyzer.Requires[0] != dependency || analyzer.ResultType != reflect.TypeOf(result) || analyzer.FactTypes[0] != result || analyzer.Flags.Lookup("sample") == nil {
		t.Fatal("registration changed analyzer identity or metadata")
	}
	got, err := analyzer.Run(pass)
	if err != nil {
		t.Fatal(err)
	}
	if got != result || exported != result || len(reported) != 1 || reported[0].Message != "handwritten" {
		t.Fatalf("result, facts, or reports changed: %v, %v, %v", got, exported, reported)
	}
	// The original pass is not mutated, and its callback still works unchanged.
	pass.Report(analysis.Diagnostic{Pos: source.files[1].Package})
	if len(reported) != 2 {
		t.Fatal("the original reporting callback was replaced")
	}
}

func TestGeneratedFixes(t *testing.T) {
	source := parseInputs(t)
	handwritten := source.files[0].Package
	generated := source.files[1].Package
	unknown := source.fset.AddFile("unknown.go", -1, 10).Pos(0)
	original := analysis.Diagnostic{
		Pos: handwritten, Message: "keep this diagnostic",
		Related: []analysis.RelatedInformation{{Pos: generated, Message: "generated context"}},
		SuggestedFixes: []analysis.SuggestedFix{
			{Message: "safe", TextEdits: []analysis.TextEdit{{Pos: handwritten, End: handwritten + 1, NewText: []byte("x")}}},
			{Message: "safe insertion", TextEdits: []analysis.TextEdit{{Pos: handwritten, NewText: []byte("x")}}},
			{Message: "generated", TextEdits: []analysis.TextEdit{{Pos: generated, End: generated + 1}}},
			{Message: "mixed", TextEdits: []analysis.TextEdit{{Pos: handwritten, End: handwritten + 1}, {Pos: generated, End: generated + 1}}},
			{Message: "generated insertion", TextEdits: []analysis.TextEdit{{Pos: generated}}},
			{Message: "unknown file", TextEdits: []analysis.TextEdit{{Pos: unknown}}},
		},
	}
	var reported []analysis.Diagnostic
	analyzer := &analysis.Analyzer{Run: func(pass *analysis.Pass) (any, error) {
		pass.Report(original)
		pass.Report(analysis.Diagnostic{Message: "no source position"})
		return nil, nil
	}}
	IgnoreGenerated(analyzer)
	_, err := analyzer.Run(&analysis.Pass{
		Fset: source.fset, Files: source.files,
		Report: func(diagnostic analysis.Diagnostic) { reported = append(reported, diagnostic) },
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(reported) != 2 || reported[0].Message != original.Message || len(reported[0].Related) != 1 {
		t.Fatalf("nongenerated diagnostics or related context lost: %v", reported)
	}
	fixes := reported[0].SuggestedFixes
	if len(fixes) != 2 || fixes[0].Message != "safe" || fixes[1].Message != "safe insertion" {
		t.Fatalf("unsafe fixes retained or safe fixes lost: %v", fixes)
	}
	if len(original.SuggestedFixes) != 6 || original.SuggestedFixes[2].Message != "generated" {
		t.Fatal("the original diagnostic was mutated")
	}
}

func TestGeneratedDependenciesAndFailures(t *testing.T) {
	source := parseInputs(t)
	var calls int
	dependency := &analysis.Analyzer{Run: func(pass *analysis.Pass) (any, error) {
		calls++
		pass.Reportf(source.files[1].Package, "generated")
		return nil, io.EOF
	}}
	root := &analysis.Analyzer{Requires: []*analysis.Analyzer{dependency}}
	IgnoreGenerated(root, dependency)
	_, err := dependency.Run(&analysis.Pass{
		Fset: source.fset, Files: source.files[1:],
		Report: func(diagnostic analysis.Diagnostic) { t.Error("generated diagnostic escaped", diagnostic) },
	})
	if err != io.EOF || calls != 1 {
		t.Fatalf("an all-generated action must still run and preserve errors: %d, %v", calls, err)
	}
}

type sampleFact struct{}

func (*sampleFact) AFact() {}

type testInputs struct {
	fset  *token.FileSet
	files []*ast.File
}

func parseInputs(tb testing.TB) testInputs {
	tb.Helper()
	source := testInputs{fset: token.NewFileSet()}
	for _, file := range []struct{ name, text string }{
		{"handwritten.go", "package sample\nvar value int\n"},
		// Test physical file identity despite a misleading logical filename.
		{"generated.go", "// Code generated by rules tests. DO NOT EDIT.\n//line handwritten.go:1\npackage sample\nvar generated int\n"},
	} {
		parsed, err := parser.ParseFile(source.fset, file.name, file.text, parser.ParseComments)
		if err != nil {
			tb.Fatal(err)
		}
		source.files = append(source.files, parsed)
	}
	if source.fset.Position(source.files[1].Package).Filename != "handwritten.go" {
		tb.Fatal("generated fixture must have a misleading logical filename")
	}
	return source
}
