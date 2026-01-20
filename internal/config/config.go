package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configDirName  = ".k8ctl"
	configFileName = "config.yaml"
)

type Config struct {
	CurrentContext string            `mapstructure:"current_context"`
	CurrentNS      string            `mapstructure:"current_namespace"`
	Aliases        map[string]string `mapstructure:"aliases"`
	Colors         ColorConfig       `mapstructure:"colors"`
	Output         OutputConfig      `mapstructure:"output"`
}

type ColorConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type OutputConfig struct {
	Format string `mapstructure:"format"` // table, json, yaml
}

var (
	configDir  string
	configPath string
	cfg        *Config
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get home directory: %v", err))
	}

	configDir = filepath.Join(homeDir, configDirName)
	configPath = filepath.Join(configDir, configFileName)
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() string {
	return configDir
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() string {
	return configPath
}

// Load loads the configuration from file
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// Set defaults
	viper.SetDefault("current_context", "")
	viper.SetDefault("current_namespace", "")
	viper.SetDefault("aliases", map[string]string{
		"g":  "get",
		"d":  "describe",
		"l":  "logs",
		"w":  "watch",
		"pf": "port-forward",
		"h":  "health",
	})
	viper.SetDefault("colors.enabled", true)
	viper.SetDefault("output.format", "table")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create default one
			if err := Save(&Config{
				CurrentContext: "",
				CurrentNS:      "",
				Aliases: map[string]string{
					"g":  "get",
					"d":  "describe",
					"l":  "logs",
					"w":  "watch",
					"pf": "port-forward",
					"h":  "health",
				},
				Colors: ColorConfig{Enabled: true},
				Output: OutputConfig{Format: "table"},
			}); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// Save saves the configuration to file
func Save(c *Config) error {
	viper.Set("current_context", c.CurrentContext)
	viper.Set("current_namespace", c.CurrentNS)
	viper.Set("aliases", c.Aliases)
	viper.Set("colors.enabled", c.Colors.Enabled)
	viper.Set("output.format", c.Output.Format)

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	cfg = c
	return nil
}

// Get returns the current configuration
func Get() *Config {
	if cfg == nil {
		var err error
		cfg, err = Load()
		if err != nil || cfg == nil {
			// If loading fails, use default config
			cfg = &Config{
				Aliases: make(map[string]string),
				Colors:  ColorConfig{Enabled: true},
				Output:  OutputConfig{Format: "table"},
			}
		}
	}
	return cfg
}

// SetCurrentContext sets the current context
func SetCurrentContext(ctx string) error {
	cfg := Get()
	cfg.CurrentContext = ctx
	return Save(cfg)
}

// SetCurrentNamespace sets the current namespace
func SetCurrentNamespace(ns string) error {
	cfg := Get()
	cfg.CurrentNS = ns
	return Save(cfg)
}

// GetCurrentContext returns the current context
func GetCurrentContext() string {
	return Get().CurrentContext
}

// GetCurrentNamespace returns the current namespace
func GetCurrentNamespace() string {
	return Get().CurrentNS
}
