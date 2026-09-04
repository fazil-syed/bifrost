package config

import "fmt"

func addDefaultIfNotExist[T comparable](field *T, defaultValue T) {
	var zeroVal T
	if *field == zeroVal {
		*field = defaultValue
	}
}

func Validate(cfg *Config) error {
	addDefaultIfNotExist(&cfg.Logging.Level, "info")
	addDefaultIfNotExist(&cfg.Bifrost.Name, "Bifrost")

	if err := validateAerospike(cfg.Aerospike); err != nil {
		return err
	}

	if cfg.Session.Lifetime <= 0 {
		return fmt.Errorf("session lifetime must be greater than 0")
	}

	return nil
}
