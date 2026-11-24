package optionconsumer // want `3 declaration groups \(config, Option, Audit\); limit 2`

type (
	config struct{}
	Option func(*config) error
)

func Configure() Option { return func(settings *config) error { return nil } }

// An unrelated consumer is not a constructor accepting the option contract.
func Audit(settings *config) { _ = settings }
