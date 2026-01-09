package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ObsidianVault   string        `yaml:"obsidian_vault"`
	JournalDir      string        `yaml:"journal_dir"` // Relative to ObsidianVault
	DefaultTemplate string        `yaml:"default_template"`
	Notification    *Notification `yaml:"notification,omitempty"`
}

// Notification holds notification settings
type Notification struct {
	Enabled   bool   `yaml:"enabled"`
	Frequency string `yaml:"frequency"` // "daily" for now
	Time      string `yaml:"time"`      // HH:MM format (24-hour)
}

func LoadConfig() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "journal-cli", "config.yaml")
	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		// Return default config if not found
		return &Config{
			ObsidianVault: "", // User must set this
			JournalDir:    filepath.Join("Journal", "Daily"),
		}, nil
	}
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(cfg *Config) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	journalConfigDir := filepath.Join(configDir, "journal-cli")
	if err := os.MkdirAll(journalConfigDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(journalConfigDir, "config.yaml")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "journal-cli", "config.yaml"), nil
}
