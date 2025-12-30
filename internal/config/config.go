package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ObsidianVault string `yaml:"obsidian_vault"`
	JournalDir    string `yaml:"journal_dir"` // Relative to ObsidianVault
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
