package errcheck

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	results := analysistest.Run(t, analysistest.TestData(), Analyzer, "a")
	for _, result := range results {
		for _, diagnostic := range result.Action.Diagnostics {
			if len(diagnostic.SuggestedFixes) != 0 {
				t.Errorf("error handling must not be changed automatically: %s", diagnostic.Message)
			}
		}
	}
}
