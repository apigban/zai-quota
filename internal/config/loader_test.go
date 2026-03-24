package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig_EnvVars tests loading configuration from environment variables
func TestLoadConfig_EnvVars(t *testing.T) {
	// Save and restore original env vars
	oldAPIKey := os.Getenv("ZAI_API_KEY")
	oldEndpoint := os.Getenv("ZAI_ENDPOINT")
	oldTimeout := os.Getenv("ZAI_TIMEOUT_SECONDS")
	defer func() {
		if oldAPIKey != "" {
			os.Setenv("ZAI_API_KEY", oldAPIKey)
		} else {
			os.Unsetenv("ZAI_API_KEY")
		}
		if oldEndpoint != "" {
			os.Setenv("ZAI_ENDPOINT", oldEndpoint)
		} else {
			os.Unsetenv("ZAI_ENDPOINT")
		}
		if oldTimeout != "" {
			os.Setenv("ZAI_TIMEOUT_SECONDS", oldTimeout)
		} else {
			os.Unsetenv("ZAI_TIMEOUT_SECONDS")
		}
	}()

	// Set env vars
	os.Setenv("ZAI_API_KEY", "test-env-key")
	os.Setenv("ZAI_ENDPOINT", "https://env.example.com")
	os.Setenv("ZAI_TIMEOUT_SECONDS", "10")

	// Create a minimal Cobra command
	cmd := &cobra.Command{
		Use: "test",
	}

	cfg, err := LoadConfig(cmd)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify env vars were loaded (they override defaults)
	assert.Equal(t, "test-env-key", cfg.APIKey)
	assert.Equal(t, "https://env.example.com", cfg.Endpoint)
	assert.Equal(t, 10, cfg.TimeoutSeconds)
}

// TestLoadConfig_ConfigFile tests loading configuration from YAML file
func TestLoadConfig_ConfigFile(t *testing.T) {
	// Save and restore original env vars
	oldAPIKey := os.Getenv("ZAI_API_KEY")
	defer func() {
		if oldAPIKey != "" {
			os.Setenv("ZAI_API_KEY", oldAPIKey)
		} else {
			os.Unsetenv("ZAI_API_KEY")
		}
	}()

	// Ensure no env var interference
	os.Unsetenv("ZAI_API_KEY")

	// Create a temporary config file
	homeDir := t.TempDir()
	configPath := filepath.Join(homeDir, ".zai-quota.yaml")

	configContent := `
api_key: "test-file-key"
endpoint: "https://file.example.com"
timeout_seconds: 15
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Override home directory for testing
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer func() {
		os.Setenv("HOME", oldHome)
	}()

	cmd := &cobra.Command{
		Use: "test",
	}

	cfg, err := LoadConfig(cmd)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify config file was loaded
	assert.Equal(t, "test-file-key", cfg.APIKey)
	assert.Equal(t, "https://file.example.com", cfg.Endpoint)
	assert.Equal(t, 15, cfg.TimeoutSeconds)
}

// TestLoadConfig_Precedence tests that env vars override config file
func TestLoadConfig_Precedence(t *testing.T) {
	// Save and restore original env vars
	oldAPIKey := os.Getenv("ZAI_API_KEY")
	oldEndpoint := os.Getenv("ZAI_ENDPOINT")
	defer func() {
		if oldAPIKey != "" {
			os.Setenv("ZAI_API_KEY", oldAPIKey)
		} else {
			os.Unsetenv("ZAI_API_KEY")
		}
		if oldEndpoint != "" {
			os.Setenv("ZAI_ENDPOINT", oldEndpoint)
		} else {
			os.Unsetenv("ZAI_ENDPOINT")
		}
	}()

	// Create a config file with values
	homeDir := t.TempDir()
	configPath := filepath.Join(homeDir, ".zai-quota.yaml")

	configContent := `
api_key: "file-key"
endpoint: "https://file.example.com"
timeout_seconds: 20
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Override home directory for testing
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer func() {
		os.Setenv("HOME", oldHome)
	}()

	// Set env vars (these should override file)
	os.Setenv("ZAI_API_KEY", "env-override-key")
	os.Setenv("ZAI_ENDPOINT", "https://env-override.example.com")

	cmd := &cobra.Command{
		Use: "test",
	}

	cfg, err := LoadConfig(cmd)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify env vars take precedence over file
	assert.Equal(t, "env-override-key", cfg.APIKey, "env should override file")
	assert.Equal(t, "https://env-override.example.com", cfg.Endpoint, "env should override file")
	assert.Equal(t, 20, cfg.TimeoutSeconds, "file value should be used when env not set")
}

// TestLoadConfig_MissingConfigFile tests that missing config file uses defaults
func TestLoadConfig_MissingConfigFile(t *testing.T) {
	// Save and restore original env vars
	oldAPIKey := os.Getenv("ZAI_API_KEY")
	oldEndpoint := os.Getenv("ZAI_ENDPOINT")
	oldTimeout := os.Getenv("ZAI_TIMEOUT_SECONDS")
	defer func() {
		if oldAPIKey != "" {
			os.Setenv("ZAI_API_KEY", oldAPIKey)
		} else {
			os.Unsetenv("ZAI_API_KEY")
		}
		if oldEndpoint != "" {
			os.Setenv("ZAI_ENDPOINT", oldEndpoint)
		} else {
			os.Unsetenv("ZAI_ENDPOINT")
		}
		if oldTimeout != "" {
			os.Setenv("ZAI_TIMEOUT_SECONDS", oldTimeout)
		} else {
			os.Unsetenv("ZAI_TIMEOUT_SECONDS")
		}
	}()

	// Ensure no config file and no env vars
	os.Unsetenv("ZAI_API_KEY")
	os.Unsetenv("ZAI_ENDPOINT")
	os.Unsetenv("ZAI_TIMEOUT_SECONDS")

	// Use a temp directory with no config file
	homeDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer func() {
		os.Setenv("HOME", oldHome)
	}()

	cmd := &cobra.Command{
		Use: "test",
	}

	cfg, err := LoadConfig(cmd)
	require.NoError(t, err, "should not error with missing config file")
	require.NotNil(t, cfg)

	// Verify defaults are used
	assert.Equal(t, "", cfg.APIKey, "APIKey should be empty default")
	assert.Equal(t, "https://api.z.ai/api/monitor/usage/quota/limit", cfg.Endpoint, "should use default endpoint")
	assert.Equal(t, 5, cfg.TimeoutSeconds, "should use default timeout")
}

// TestLoadConfig_FlagBinding tests that Cobra flags are properly bound
func TestLoadConfig_FlagBinding(t *testing.T) {
	// Save and restore env vars
	oldAPIKey := os.Getenv("ZAI_API_KEY")
	defer func() {
		if oldAPIKey != "" {
			os.Setenv("ZAI_API_KEY", oldAPIKey)
		} else {
			os.Unsetenv("ZAI_API_KEY")
		}
	}()

	// Use temp directory with no config file
	homeDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer func() {
		os.Setenv("HOME", oldHome)
	}()

	// Create command with a flag
	var jsonOutput bool
	cmd := &cobra.Command{
		Use: "test",
	}
	cmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	cfg, err := LoadConfig(cmd)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Config should load successfully even with flags
	assert.NotNil(t, cfg)
}

func TestSaveConfig(t *testing.T) {
	homeDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", oldHome)

	cfg := &Config{
		APIKey:         "test-api-key-123",
		Endpoint:       "https://test.example.com",
		TimeoutSeconds: 30,
	}

	err := SaveConfig(cfg)
	require.NoError(t, err)

	configPath := filepath.Join(homeDir, ".zai-quota.yaml")
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	assert.Contains(t, string(data), "api_key: test-api-key-123")
	assert.Contains(t, string(data), "endpoint: https://test.example.com")
	assert.Contains(t, string(data), "timeout_seconds: 30")
}

func TestSaveConfig_FilePermissions(t *testing.T) {
	homeDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", oldHome)

	cfg := &Config{
		APIKey:         "secret-key",
		Endpoint:       "https://api.z.ai",
		TimeoutSeconds: 5,
	}

	err := SaveConfig(cfg)
	require.NoError(t, err)

	configPath := filepath.Join(homeDir, ".zai-quota.yaml")
	info, err := os.Stat(configPath)
	require.NoError(t, err)

	assert.Equal(t, os.FileMode(0600), info.Mode().Perm(), "config file should have 0600 permissions")
}

func TestConfigFileExists(t *testing.T) {
	homeDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", oldHome)

	exists, err := ConfigFileExists()
	require.NoError(t, err)
	assert.False(t, exists)

	configPath := filepath.Join(homeDir, ".zai-quota.yaml")
	err = os.WriteFile(configPath, []byte("api_key: test"), 0600)
	require.NoError(t, err)

	exists, err = ConfigFileExists()
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestLoadConfigFromFile(t *testing.T) {
	homeDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", oldHome)

	configContent := `
api_key: file-key-xyz
endpoint: https://file.example.com
timeout_seconds: 25
`
	configPath := filepath.Join(homeDir, ".zai-quota.yaml")
	err := os.WriteFile(configPath, []byte(configContent), 0600)
	require.NoError(t, err)

	cfg, err := LoadConfigFromFile()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "file-key-xyz", cfg.APIKey)
	assert.Equal(t, "https://file.example.com", cfg.Endpoint)
	assert.Equal(t, 25, cfg.TimeoutSeconds)
}
