package entrypoints

import framework "golang.org/x/tools/go/analysis"

var Late = alias{Run: late}

type alias = framework.Analyzer

type support struct{}

func (*support) Method() {}

func late(*framework.Pass) (any, error) { // want `analysis entrypoints must appear before methods`
	return nil, nil
}

func supportingFunction() {}

func lateCrossFile(*framework.Pass) (any, error) { // want `analysis entrypoints must appear before functions`
	return nil, nil
}
