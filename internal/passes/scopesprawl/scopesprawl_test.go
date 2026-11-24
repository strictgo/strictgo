package scopesprawl

import (
	"path/filepath"
	"reflect"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	previous := maxGroups
	maxGroups = 2
	t.Cleanup(func() { maxGroups = previous })
	results := analysistest.Run(t, analysistest.TestData(), Analyzer, "owned", "shared", "concepts", "functions", "state", "crossfile", "parameters", "represented", "generic", "shadowing", "generated", "tests", "options", "genericoptions", "sharedoptions", "optionconsumer", "invalidoption", "support", "privatedata", "sharedconstants", "sharedstate", "familycross")
	want := map[string][]string{
		"familycross/types.go":         {"Client", "Option"},
		"familycross/constructor.go":   {"Client"},
		"familycross/options.go":       {"Option"},
		"owned/client.go":              {"Client"},
		"owned/methods.go":             {"Client"},
		"shared/shared.go":             {"Client", "ParseConfig", "normalize"},
		"concepts/concepts.go":         {"Client", "ParseConfig", "RetryPolicy"},
		"functions/functions.go":       {"Encode", "Decode", "Inspect"},
		"functions/private.go":         {"first", "second", "normalize"},
		"state/state.go":               {"Input", "Output", "Mode"},
		"state/blocks.go":              {"First", "Current"},
		"crossfile/client.go":          {"Client", "privateHelper"},
		"crossfile/parser.go":          {"ParseConfig"},
		"parameters/config.go":         {"Config", "ParseConfig", "ValidateConfig"},
		"represented/types.go":         {"Client"},
		"represented/methods.go":       {"Client"},
		"represented/caller.go":        {"ParseConfig"},
		"generic/generic.go":           {"Box", "Operate"},
		"shadowing/shadowing.go":       {"Client", "ParseConfig"},
		"generated/source.go":          {"Client", "normalize"},
		"tests/source.go":              {"Client"},
		"options/family.go":            {"Client", "Option", "WithNumber"},
		"genericoptions/generic.go":    {"Client", "Option"},
		"sharedoptions/shared.go":      {"config", "ReadOption", "WriteOption"},
		"optionconsumer/consumer.go":   {"config", "Option", "Audit"},
		"invalidoption/invalid.go":     {"config", "Option", "WithValue"},
		"support/support.go":           {"Client", "policy", "recipe"},
		"privatedata/private.go":       {"config", "ParseConfig", "ValidateConfig"},
		"sharedconstants/constants.go": {"Client", "ParseConfig"},
		"sharedconstants/public.go":    {"Limit", "First", "Second"},
		"sharedstate/state.go":         {"count", "Client", "ParseConfig"},
	}
	for _, result := range results {
		pkg := result.Action.Package
		pass := &analysis.Pass{Files: pkg.Syntax, TypesInfo: pkg.TypesInfo, Pkg: pkg.Types, Fset: pkg.Fset}
		graph := inferOwnership(pass)
		for _, file := range pass.Files {
			key := pkg.PkgPath + "/" + filepath.Base(pass.Fset.File(file.Pos()).Name())
			if expected, ok := want[key]; ok {
				if got := fileGroups(graph, file); !reflect.DeepEqual(got, expected) {
					t.Errorf("%s groups = %v, want %v", key, got, expected)
				}
			}
		}
		for _, diagnostic := range result.Action.Diagnostics {
			if len(diagnostic.SuggestedFixes) != 0 {
				t.Error("concept boundaries cannot be fixed automatically")
			}
		}
	}
}

func TestInvalidLimit(t *testing.T) {
	previous := maxGroups
	maxGroups = 0
	t.Cleanup(func() { maxGroups = previous })
	if _, err := run(&analysis.Pass{}); err != nil {
		return
	}
	t.Fatal("accepted a non-positive declaration-group limit")
}
