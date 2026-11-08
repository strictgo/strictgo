package forint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(test *testing.T) {
	results := analysistest.Run(test, analysistest.TestData(), Analyzer, "a")
	for _, result := range results {
		for _, diagnostic := range result.Diagnostics {
			if len(diagnostic.SuggestedFixes) != 0 {
				test.Errorf("unexpected suggested fix: %s", diagnostic.Message)
			}
		}
	}
}
