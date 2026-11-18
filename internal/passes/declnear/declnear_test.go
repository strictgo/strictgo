package declnear

import (
	"testing"

	"github.com/strictgo/strictgo/internal/passes/declorder"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	results := analysistest.Run(t, analysistest.TestData(), Analyzer, "a")
	for _, result := range results {
		for _, diagnostic := range result.Action.Diagnostics {
			if len(diagnostic.SuggestedFixes) != 0 {
				t.Errorf("declaration moves are diagnostic-only initially: %s", diagnostic.Message)
			}
		}
	}
}

func TestCompatibleWithDeclorder(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "compatible")
	analysistest.Run(t, analysistest.TestData(), declorder.Analyzer, "compatible")
}
