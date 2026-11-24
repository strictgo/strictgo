package genericoptions

type (
	Client        struct{}
	config        struct{ enabled bool }
	Option[T any] func(*T) error
	Alias         = Option[config]
)

func New(opts ...Option[config]) *Client {
	var settings config
	for _, option := range opts {
		_ = option(&settings)
	}
	return &Client{}
}

func Configure() Alias {
	return func(settings *config) error { settings.enabled = true; return nil }
}
