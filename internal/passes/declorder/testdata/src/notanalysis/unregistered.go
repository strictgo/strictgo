package notanalysis

import "golang.org/x/tools/go/analysis"

var unusedRegistration = analysis.Analyzer{}

type otherSupport struct{}

func unregisteredRun(*analysis.Pass) (any, error) { return nil, nil }

func (*otherSupport) Method() {} // want `methods must appear before functions`
