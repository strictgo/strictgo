package a

type options struct { // want `move supporting struct "options" directly before parseOptions, its first explicit use`
	value int
}

// Inferred result types and unrelated local types are not explicit references
// to the package-level options declaration.
func Run() {
	type options struct{}
	_ = options{}
	_ = parseOptions()
}

func report() {}

func parseOptions() options { return options{} }

func check(settings options) { _ = settings }
