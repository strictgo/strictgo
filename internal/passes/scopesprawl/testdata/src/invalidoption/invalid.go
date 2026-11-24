package invalidoption // want `3 declaration groups \(config, Option, WithValue\); limit 2`

type config struct{}

// This is a callback, not a recognized functional option contract.
type Option func(*config) bool

func WithValue() Option { return func(settings *config) bool { return true } }
