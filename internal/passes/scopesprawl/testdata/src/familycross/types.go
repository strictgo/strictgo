package familycross

type (
	Client struct{}
	config struct{ enabled bool }
	Option func(*config) error
)
