package badname

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "a", "globals", "testingnames", "sortnames", "errornames", "optionnames", "loopnames", "filenames", "packagenames/...")
}

func TestVagueName(t *testing.T) {
	for _, name := range []string{"util", "utils", "stringutil", "httputils", "helper", "helpers", "testhelpers", "misc", "appmisc", "Utils"} {
		if !vagueName(name) {
			t.Errorf("vagueName(%q) = false, want true", name)
		}
	}
	for _, name := range []string{"http", "strings", "utility", "utilization", "helperService", "miscellaneous", "civil"} {
		if vagueName(name) {
			t.Errorf("vagueName(%q) = true, want false", name)
		}
	}
}
