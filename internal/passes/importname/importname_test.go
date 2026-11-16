package importname

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestSuggestedFixes(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), Analyzer, "a")
}

func TestUnsafeRenames(t *testing.T) {
	results := analysistest.Run(t, analysistest.TestData(), Analyzer, "nofix")
	for _, result := range results {
		for _, diagnostic := range result.Action.Diagnostics {
			if len(diagnostic.SuggestedFixes) != 0 {
				t.Errorf("unexpected fix for %s", diagnostic.Message)
			}
		}
	}
}

func TestLowercaseName(t *testing.T) {
	for _, name := range []string{"fmt", "jsonv2", "http2", "λ"} {
		if !lowercaseName(name) {
			t.Errorf("lowercaseName(%q) = false, want true", name)
		}
	}
	for _, name := range []string{"", "_", ".", "Fmt", "foo_bar", "2json", "type", "中文"} {
		if lowercaseName(name) {
			t.Errorf("lowercaseName(%q) = true, want false", name)
		}
	}
}
