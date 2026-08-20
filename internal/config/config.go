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
	Host     string     `yaml:"host"`
	Port     int        `yaml:"port"`
	Name     string     `yaml:"name"`
	User     string     `yaml:"user"`
	Password string     `yaml:"password"`
	SSL      SSLConfig  `yaml:"ssl"`
	Pool     PoolConfig `yaml:"pool"`
}

type SSLConfig struct {
	Mode     string `yaml:"mode"`
	RootCert string `yaml:"root_cert"`
	Cert     string `yaml:"cert"`
	Key      string `yaml:"key"`
}

type PoolConfig struct {
	MaxConns          int    `yaml:"max_conns"`
	MinConns          int    `yaml:"min_conns"`
	MaxConnLifetime   string `yaml:"max_conn_lifetime"`
	MaxConnIdleTime   string `yaml:"max_conn_idle_time"`
	HealthCheckPeriod string `yaml:"health_check_period"`
}
