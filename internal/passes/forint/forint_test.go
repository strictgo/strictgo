package forint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	results := analysistest.Run(t, analysistest.TestData(), Analyzer, "a")
	for _, result := range results {
		for _, diagnostic := range result.Diagnostics {
			if len(diagnostic.SuggestedFixes) != 0 {
				t.Errorf("unexpected suggested fix: %s", diagnostic.Message)
			}
		}
	}
}
