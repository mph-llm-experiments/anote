package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the application configuration.
type Config struct {
	IdeasDirectory string `toml:"ideas_directory"`
	Editor         string `toml:"editor"`
}

// DefaultConfig returns default configuration.
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		IdeasDirectory: filepath.Join(homeDir, "ideas"),
		Editor:         "vim",
	}
}

// Load reads configuration from a file. If path is empty, searches standard locations.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		path = findConfigFile()
	}

	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	cfg.IdeasDirectory = expandHome(cfg.IdeasDirectory)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.IdeasDirectory == "" {
		return fmt.Errorf("ideas_directory cannot be empty")
	}

	if info, err := os.Stat(c.IdeasDirectory); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("ideas_directory does not exist: %s", c.IdeasDirectory)
		}
		return fmt.Errorf("failed to check ideas_directory: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("ideas_directory is not a directory: %s", c.IdeasDirectory)
	}

	return nil
}

// findConfigFile looks for a config file in standard locations.
func findConfigFile() string {
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		path := filepath.Join(xdgConfig, "anote", "config.toml")
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	if homeDir, err := os.UserHomeDir(); err == nil {
		path := filepath.Join(homeDir, ".config", "anote", "config.toml")
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// expandHome expands ~ to the user's home directory.
func expandHome(path string) string {
	if path == "" || path[0] != '~' {
		return path
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	return filepath.Join(homeDir, path[1:])
}

// ConfigPath returns the default config file path.
func ConfigPath() string {
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "anote", "config.toml")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".config", "anote", "config.toml")
}
