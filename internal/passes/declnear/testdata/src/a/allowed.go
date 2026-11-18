package a

type Public struct{}

type withMethods struct{}

type (
	grouped struct{}
)

type aliased = struct{}

type number int

type registered struct{}

var registration = registered{}

type shared struct{}

type prerequisite struct{}

type dependent struct{ value prerequisite }

var dependency = dependent{}

type unused struct{}

func unrelated() {}

func (withMethods) Method() {}

func useAllowed(public Public, method withMethods, group grouped, alias aliased, numeric number, registration registered, value shared) {
	_, _, _, _, _, _, _ = public, method, group, alias, numeric, registration, value
}

func useLate(value late) { _ = value }

type late struct{}

// No unrelated function separates an adjacent type from its first reference.
type adjacent struct{}

func parseAdjacent() adjacent { return adjacent{} }

func local() {
	type localOptions struct{}
	println("local declarations are outside the initial rule")
	_ = localOptions{}
}
