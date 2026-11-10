package optionnames

type config struct{ Enabled bool }

type (
	Option         func(*config)
	FallibleOption func(*config) error
	Alias          = Option
)

type GenericOption[T any] func(*T)

func WithEnabled() Option {
	return func(o *config) {
		x := true // want `variable "x" has a single-letter name outside a short scope`
		o.Enabled = x
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
	}
}

func WithError() FallibleOption {
	return func(o *config) error {
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		return nil
	}
}

func apply(options ...Option) {
	o := config{}
	for _, f := range options {
		f(&o)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
	}
	println(o.Enabled)
}

func evaluate(f Alias, o *config) {
	f(o)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
}

func generic(f GenericOption[config], o *config) {
	f(o)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
}

type (
	Options       struct{ Enabled bool }
	Callback      func(*Options)
	InvalidOption func(int)
)

func unrelated(o *Options, f Callback, g InvalidOption) { // want `variable "o" has a single-letter name outside a short scope` `variable "f" has a single-letter name outside a short scope` `variable "g" has a single-letter name outside a short scope`
	f(o)
	g(1)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
	println(o.Enabled)
}

func errorNames() FallibleOption {
	return func(o *config) (e error) { // want `error variable "e" should be named err`
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		println(o.Enabled)
		return nil
	}
}
