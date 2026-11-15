package entrypoints

import framework "golang.org/x/tools/go/analysis"

var Analyzer = &framework.Analyzer{Run: analyze}

// Registration by identity also recognizes entrypoints not named run.
func analyze(*framework.Pass) (any, error) { return nil, nil }

type visitor struct{}

func (*visitor) Visit() {}

func helper() {}
