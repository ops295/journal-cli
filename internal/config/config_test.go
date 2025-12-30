package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadConfigDefaultWhenMissing(t *testing.T) {
	tmp := t.TempDir()
	// Set HOME/XDG_CONFIG_HOME for deterministic UserConfigDir resolution
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if cfg.JournalDir == "" {
		t.Fatalf("expected default JournalDir, got empty")
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)

	// Determine where the system expects the config to be
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("UserConfigDir error: %v", err)
	}

	cfgDir := filepath.Join(userConfigDir, "journal-cli")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	c := Config{ObsidianVault: "/tmp/vault", JournalDir: "Journal/Daily"}
	data, _ := yaml.Marshal(c)
	if err := os.WriteFile(filepath.Join(cfgDir, "config.yaml"), data, 0644); err != nil {
		t.Fatalf("write config failed: %v", err)
	}

	got, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if got.ObsidianVault != c.ObsidianVault {
		t.Fatalf("vault mismatch: got %s want %s", got.ObsidianVault, c.ObsidianVault)
	}
}
