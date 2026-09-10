package limitscope

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	path := analysistest.TestData()
	analysistest.Run(t, path, Analyzer, "a")
}
