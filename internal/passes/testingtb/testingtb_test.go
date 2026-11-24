package testingtb

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	results := analysistest.Run(t, analysistest.TestData(), Analyzer, "a")
	for _, result := range results {
		for _, diagnostic := range result.Action.Diagnostics {
			if len(diagnostic.SuggestedFixes) != 0 {
				t.Error("widening helper APIs is diagnostic-only initially")
			}
		}
	}
}
