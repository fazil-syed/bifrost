package config

type Config struct {
	Logging LoggingConfig `yaml:"logging"`
	Bifrost BifrostConfig `yaml:"bifrost"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

type BifrostConfig struct {
	Name string `yaml:"name"`
}
