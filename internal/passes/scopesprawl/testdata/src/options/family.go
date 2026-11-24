package options // want `3 declaration groups \(Client, Option, WithNumber\); limit 2`

type (
	Client      struct{}
	config      struct{ enabled bool }
	Option      func(*config) error
	OptionAlias = Option
)

func New(opts ...Option) (*Client, error) {
	var settings config
	for _, option := range opts {
		if err := option(&settings); err != nil {
			return nil, err
		}
	}
	return &Client{}, nil
}

// Factory membership comes from the contract, not the function name.
func Configure() Option {
	return func(settings *config) error { settings.enabled = true; return nil }
}

func WithAlias() OptionAlias     { return Configure() }
func WithError() (Option, error) { return Configure(), nil }
func WithNumber() int            { return 1 }
