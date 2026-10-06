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
	if cfg.Token.AccessLifetime <= 0 {
		return fmt.Errorf("access token lifetime must be greater than 0")
	}
	if cfg.Token.RefreshLifetime <= 0 {
		return fmt.Errorf("refresh token lifetime must be greater than 0")
	}

	if err := validateHTTP(cfg.HTTP); err != nil {
		return err
	}

	return nil
}

func validateHTTP(cfg HTTPConfig) error {
	if cfg.Host == "" {
		return fmt.Errorf("http host must not be empty")
	}

	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("http port must be between 1 and 65535")
	}

	if cfg.ReadTimeout <= 0 {
		return fmt.Errorf("http read timeout must be greater than 0")
	}
	if cfg.WriteTimeout <= 0 {
		return fmt.Errorf("http write timeout must be greater than 0")
	}
	if cfg.IdleTimeout <= 0 {
		return fmt.Errorf("http idle timeout must be greater than 0")
	}
	if cfg.ShutdownTimeout <= 0 {
		return fmt.Errorf("http shutdown timeout must be greater than 0")
	}
	return nil
}
