package familycross

func Configure() Option {
	return func(settings *config) error { settings.enabled = true; return nil }
}
