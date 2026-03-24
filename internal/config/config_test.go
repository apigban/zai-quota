package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "default values",
			want: &Config{
				APIKey:         "",
				Endpoint:       "https://api.z.ai/api/monitor/usage/quota/limit",
				TimeoutSeconds: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); *got != *tt.want {
				t.Errorf("NewConfig() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestConfig_YAMLUnmarshaling(t *testing.T) {
	yamlData := `
api_key: test_key
endpoint: https://custom.endpoint.com
timeout_seconds: 10
`
	expected := &Config{
		APIKey:         "test_key",
		Endpoint:       "https://custom.endpoint.com",
		TimeoutSeconds: 10,
	}

	var cfg Config
	if err := yaml.Unmarshal([]byte(yamlData), &cfg); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if cfg != *expected {
		t.Errorf("UnmarshalYAML() = %v, want %v", cfg, *expected)
	}
}
