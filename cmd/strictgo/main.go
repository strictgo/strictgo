package main

import (
	"os"

	"github.com/strictgo/strictgo/internal/diff"
	"github.com/strictgo/strictgo/internal/passes/badname"
	"github.com/strictgo/strictgo/internal/passes/blankparam"
	"github.com/strictgo/strictgo/internal/passes/boolswitch"
	"github.com/strictgo/strictgo/internal/passes/declnear"
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
	if len(os.Args) > 1 && os.Args[1] == "diff" {
		os.Exit(diff.Run(os.Args[2:], os.Stdout, os.Stderr))
	}
	multichecker.Main(
		badname.Analyzer,
		blankparam.Analyzer,
		boolswitch.Analyzer,
		declnear.Analyzer,
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
