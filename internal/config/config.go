package config

// Config holds the configuration for the z.ai quota tool.
type Config struct {
	APIKey         string `mapstructure:"api_key" yaml:"api_key"`
	Endpoint       string `mapstructure:"endpoint" yaml:"endpoint"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds" yaml:"timeout_seconds"`
}

// NewConfig returns a new Config with default values.
func NewConfig() *Config {
	return &Config{
		Endpoint:       "https://api.z.ai/api/monitor/usage/quota/limit",
		TimeoutSeconds: 5,
	}
}
