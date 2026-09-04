package config

type AerospikeConfig struct {
	Namespace      string                    `yaml:"namespace"`
	Hosts          []AerospikeHostConfig     `yaml:"hosts"`
	Authentication AerospikeAuthConfig       `yaml:"authentication"`
	TLS            AerospikeTLSConfig        `yaml:"tls"`
	Connection     AerospikeConnectionConfig `yaml:"connection"`
	Timeout        AerospikeTimeoutConfig    `yaml:"timeout"`
}

type AerospikeHostConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type AerospikeAuthConfig struct {
	Mode     string `yaml:"mode"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type AerospikeTLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Name     string `yaml:"name"`
	CAFile   string `yaml:"ca_file"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type AerospikeConnectionConfig struct {
	MinConnectionsPerNode   int  `yaml:"min_connections_per_node"`
	ConnectionQueueSize     int  `yaml:"connection_queue_size"`
	LimitConnectionsToQueue bool `yaml:"limit_connections_to_queue"`
}

type AerospikeTimeoutConfig struct {
	Socket              string `yaml:"socket"`
	Total               string `yaml:"total"`
	MaxRetries          int    `yaml:"max_retries"`
	SleepBetweenRetries string `yaml:"sleep_between_retries"`
	TimeoutDelay        string `yaml:"timeout_delay"`
}
