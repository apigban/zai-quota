package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	configFileName = ".zai-quota"
	configFileExt  = "yaml"
	envPrefix      = "ZAI"
)

// LoadConfig loads configuration from multiple sources with precedence:
// flags > env vars > config file > defaults
// Returns error only if there's a critical failure (not for missing file)
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	v := viper.New()

	// Set defaults first
	setDefaults(v)

	// Bind Cobra flags to Viper
	if err := bindFlags(v, cmd); err != nil {
		return nil, fmt.Errorf("failed to bind flags: %w", err)
	}

	// Configure environment variables
	setupEnvVars(v)

	// Load from config file if it exists (gracefully ignore if missing)
	if err := loadConfigFile(v); err != nil {
		// Only log, don't return error - missing file is OK
		_ = err
	}

	// Unmarshal into Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	v.SetDefault("endpoint", "https://api.z.ai/api/monitor/usage/quota/limit")
	v.SetDefault("timeout_seconds", 5)
}

// bindFlags binds Cobra flags to Viper keys
func bindFlags(v *viper.Viper, cmd *cobra.Command) error {
	// Bind flags to Viper using flag names as keys
	// Example: --json flag could be bound if needed
	if err := v.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	if err := v.BindPFlags(cmd.PersistentFlags()); err != nil {
		return err
	}
	return nil
}

// setupEnvVars configures environment variable binding
func setupEnvVars(v *viper.Viper) {
	// Set env prefix (ZAI_)
	v.SetEnvPrefix(envPrefix)

	// Automatically map env vars (e.g., ZAI_API_KEY -> api_key)
	v.AutomaticEnv()

	// Use replacer to convert uppercase env var suffixes to lowercase keys
	// Example: ZAI_API_KEY -> prefix="ZAI", key="API_KEY" -> "api_key"
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Explicitly bind env vars to ensure proper mapping
	v.BindEnv("api_key")
	v.BindEnv("endpoint")
	v.BindEnv("timeout_seconds")
}

// SaveConfig saves the configuration to ~/.zai-quota.yaml with 0600 permissions
func SaveConfig(cfg *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, configFileName+"."+configFileExt)

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ConfigFileExists checks if the config file exists
func ConfigFileExists() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, configFileName+"."+configFileExt)
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// LoadConfigFromFile loads configuration from file only (no flags/env)
// Used for merging with existing config during setup
func LoadConfigFromFile() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, configFileName+"."+configFileExt)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// loadConfigFile attempts to load config from ~/.zai-quota.yaml
// Returns error only if file exists but can't be read (not for missing file)
func loadConfigFile(v *viper.Viper) error {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Construct config file path
	configPath := filepath.Join(homeDir, configFileName+"."+configFileExt)

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// File doesn't exist - this is OK, use defaults
		return nil
	}

	// Configure Viper to read from file
	v.SetConfigName(configFileName)
	v.SetConfigType(configFileExt)
	v.AddConfigPath(homeDir)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	return nil
}
