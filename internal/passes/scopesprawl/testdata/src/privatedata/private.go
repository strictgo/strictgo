package privatedata // want `3 declaration groups \(config, ParseConfig, ValidateConfig\); limit 2`

type config struct{}

func ParseConfig(value *config)    { _ = value }
func ValidateConfig(value *config) { _ = value }
