package config

func addDefaultIfNotExist[T comparable](field *T, defaultValue T) {
	var zeroVal T
	if field == &zeroVal {
		field = &defaultValue
	}
}

func Validate(cfg *Config) error {
	addDefaultIfNotExist(&cfg.Logging.Level, "info")
	addDefaultIfNotExist(&cfg.Bifrost.Name, "Bifrost")

	return nil
}
