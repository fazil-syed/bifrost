package config

func Validate(cfg *Config) error {

	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
	if cfg.Bifrost.Name == "" {
		cfg.Bifrost.Name = "bifrost"
	}

	return nil
}
