package blankparam

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestSuggestedFixes(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), Analyzer, "a")
}

func TestCollisions(t *testing.T) {
	results := analysistest.Run(t, analysistest.TestData(), Analyzer, "nofix")
	for _, result := range results {
		for _, diagnostic := range result.Action.Diagnostics {
			if len(diagnostic.SuggestedFixes) != 0 {
				t.Errorf("unexpected fix for %s", diagnostic.Message)
			}
		}
	}
}

func TestLowerCamel(t *testing.T) {
	cases := map[string]string{
		"Request":    "request",
		"HTTPClient": "httpClient",
		"ID":         "id",
		"config":     "config",
	}
	for name, want := range cases {
		if got := lowerCamel(name); got != want {
			t.Errorf("lowerCamel(%q) = %q, want %q", name, got, want)
		}
	}
}
