package entrypoints

import "golang.org/x/tools/go/analysis"

var CrossFile = analysis.Analyzer{Run: lateCrossFile}
