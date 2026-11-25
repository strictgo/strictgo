package singleusestring

import (
	"strings"
	"testing"

	"github.com/strictgo/strictgo/internal/reportingrules"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestSuggestedFixes(t *testing.T) {
	analyzer := *Analyzer
	reportingrules.IgnoreGenerated(&analyzer)
	results := analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), &analyzer, "a")
	for _, result := range results {
		for _, diagnostic := range result.Action.Diagnostics {
			if strings.Contains(diagnostic.Message, `"generatedUse"`) && len(diagnostic.SuggestedFixes) != 0 {
				t.Error("a fix must not edit generated consumers")
			}
		}
	}
}

func TestUnsafeFixes(t *testing.T) {
	analyzer := *Analyzer
	reportingrules.IgnoreGenerated(&analyzer)
	results := analysistest.Run(t, analysistest.TestData(), &analyzer, "nofix")
	for _, result := range results {
		for _, diagnostic := range result.Action.Diagnostics {
			if len(diagnostic.SuggestedFixes) != 0 {
				t.Errorf("unexpected fix for %s", diagnostic.Message)
			}
		}
	}
}
