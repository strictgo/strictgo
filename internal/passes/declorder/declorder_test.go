package declorder

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "a", "blocks", "blockerrors", "entrypoints", "notanalysis", "mainbefore", "mainafter", "mainlate", "testorder")
}

func TestConstructorName(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{"New", true},
		{"new", true},
		{"NewClient", true},
		{"newClient", true},
		{"NewÉlément", true},
		{"Newline", false},
		{"newer", false},
		{"Renew", false},
		{"New_client", false},
		{"New2", false},
		{"run", false},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if got := constructorName(test.name); got != test.want {
				t.Errorf("constructorName(%q) = %v, want %v", test.name, got, test.want)
			}
		})
	}
}
