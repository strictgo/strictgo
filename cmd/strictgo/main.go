package main

import (
	"github.com/konradreiche/strictgo/internal/passes/badname"
	"github.com/konradreiche/strictgo/internal/passes/boolswitch"
	"github.com/konradreiche/strictgo/internal/passes/forint"
	"github.com/konradreiche/strictgo/internal/passes/keepflat"
	"github.com/konradreiche/strictgo/internal/passes/limitscope"
	"github.com/konradreiche/strictgo/internal/passes/nilresult"
	"github.com/konradreiche/strictgo/internal/passes/returnstruct"
	"github.com/konradreiche/strictgo/internal/passes/zeroinit"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		badname.Analyzer,
		boolswitch.Analyzer,
		forint.Analyzer,
		keepflat.Analyzer,
		limitscope.Analyzer,
		nilresult.Analyzer,
		returnstruct.Analyzer,
		zeroinit.Analyzer,
	)
}
