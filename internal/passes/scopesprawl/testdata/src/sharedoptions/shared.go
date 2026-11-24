package sharedoptions // want `3 declaration groups \(config, ReadOption, WriteOption\); limit 2`

type (
	config      struct{}
	ReadOption  func(*config) error
	WriteOption func(*config) error
)

func Read() ReadOption   { return func(settings *config) error { return nil } }
func Write() WriteOption { return func(settings *config) error { return nil } }
