package main

import (
	"github.com/konradreiche/strictgo/internal/passes/badname"
	"github.com/konradreiche/strictgo/internal/passes/limitscope"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		badname.Analyzer,
		limitscope.Analyzer,
	)
}
