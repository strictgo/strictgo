package notanalysis

// A similarly shaped type from another package is not go/analysis registration.
var Custom = Analyzer{Run: run}

type (
	Analyzer struct{ Run func() }
	support  struct{}
)

func run() {}

func (*support) Method() {} // want `methods must appear before functions`
