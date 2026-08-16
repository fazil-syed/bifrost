package config

type Config struct {
	Logging  LoggingConfig  `yaml:"logging"`
	Bifrost  BifrostConfig  `yaml:"bifrost"`
	Database DatabaseConfig `yaml:"database"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

type BifrostConfig struct {
	Name string `yaml:"name"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslmode"`
}
