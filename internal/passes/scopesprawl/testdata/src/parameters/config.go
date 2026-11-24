package parameters // want `3 declaration groups \(Config, ParseConfig, ValidateConfig\); limit 2`

type Config struct{}

func ParseConfig(config Config)     { _ = config }
func ValidateConfig(config *Config) { _ = config }
