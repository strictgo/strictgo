package main

import (
	"github.com/strictgo/strictgo/internal/passes/badname"
	"github.com/strictgo/strictgo/internal/passes/blankparam"
	"github.com/strictgo/strictgo/internal/passes/boolswitch"
	"github.com/strictgo/strictgo/internal/passes/declorder"
	"github.com/strictgo/strictgo/internal/passes/errcheck"
	"github.com/strictgo/strictgo/internal/passes/forcederef"
	"github.com/strictgo/strictgo/internal/passes/forint"
	"github.com/strictgo/strictgo/internal/passes/importname"
	"github.com/strictgo/strictgo/internal/passes/keepflat"
	"github.com/strictgo/strictgo/internal/passes/limitscope"
	"github.com/strictgo/strictgo/internal/passes/nilreceiver"
	"github.com/strictgo/strictgo/internal/passes/nilresult"
	"github.com/strictgo/strictgo/internal/passes/noelse"
	"github.com/strictgo/strictgo/internal/passes/okcheck"
	"github.com/strictgo/strictgo/internal/passes/returnstruct"
	"github.com/strictgo/strictgo/internal/passes/singleassign"
	"github.com/strictgo/strictgo/internal/passes/zeroinit"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		badname.Analyzer,
		blankparam.Analyzer,
		boolswitch.Analyzer,
		declorder.Analyzer,
		errcheck.Analyzer,
		forcederef.Analyzer,
		forint.Analyzer,
		importname.Analyzer,
		keepflat.Analyzer,
		limitscope.Analyzer,
		nilreceiver.Analyzer,
		nilresult.Analyzer,
		noelse.Analyzer,
		okcheck.Analyzer,
		returnstruct.Analyzer,
		singleassign.Analyzer,
		zeroinit.Analyzer,
	)
}
