package main

import (
	"github.com/konradreiche/strictgo/internal/passes/badname"
	"github.com/konradreiche/strictgo/internal/passes/keepflat"
	"github.com/konradreiche/strictgo/internal/passes/limitscope"
	"github.com/konradreiche/strictgo/internal/passes/returnstruct"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		badname.Analyzer,
		keepflat.Analyzer,
		limitscope.Analyzer,
		returnstruct.Analyzer,
	)
}
