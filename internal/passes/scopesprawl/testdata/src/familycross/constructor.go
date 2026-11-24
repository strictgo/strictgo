package familycross

func New(opts ...Option) *Client {
	var settings config
	for _, option := range opts {
		_ = option(&settings)
	}
	return &Client{}
}
