package main

import (
	"github.com/strictgo/strictgo/internal/passes/badname"
	"github.com/strictgo/strictgo/internal/passes/boolswitch"
	"github.com/strictgo/strictgo/internal/passes/errcheck"
	"github.com/strictgo/strictgo/internal/passes/forcederef"
	"github.com/strictgo/strictgo/internal/passes/forint"
	"github.com/strictgo/strictgo/internal/passes/keepflat"
	"github.com/strictgo/strictgo/internal/passes/limitscope"
	"github.com/strictgo/strictgo/internal/passes/nilresult"
	"github.com/strictgo/strictgo/internal/passes/okcheck"
	"github.com/strictgo/strictgo/internal/passes/returnstruct"
	"github.com/strictgo/strictgo/internal/passes/singleassign"
	"github.com/strictgo/strictgo/internal/passes/zeroinit"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		badname.Analyzer,
		boolswitch.Analyzer,
		errcheck.Analyzer,
		forcederef.Analyzer,
		forint.Analyzer,
		keepflat.Analyzer,
		limitscope.Analyzer,
		nilresult.Analyzer,
		okcheck.Analyzer,
		returnstruct.Analyzer,
		singleassign.Analyzer,
		zeroinit.Analyzer,
	)
}
