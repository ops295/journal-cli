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
	t.Setenv("APPDATA", tmp)

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
	t.Setenv("APPDATA", tmp)

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

func TestSaveConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("APPDATA", tmp)

	cfg := &Config{
		ObsidianVault:   "/tmp/vault",
		JournalDir:      "Journal/Daily",
		DefaultTemplate: "daily-human-dev",
	}

	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig error: %v", err)
	}

	// Load it back
	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if loaded.DefaultTemplate != cfg.DefaultTemplate {
		t.Fatalf("default template mismatch: got %s want %s", loaded.DefaultTemplate, cfg.DefaultTemplate)
	}
}

func TestLoadConfigWithDefaultTemplate(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("APPDATA", tmp)

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("UserConfigDir error: %v", err)
	}

	cfgDir := filepath.Join(userConfigDir, "journal-cli")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	c := Config{
		ObsidianVault:   "/tmp/vault",
		JournalDir:      "Journal/Daily",
		DefaultTemplate: "gentle-day",
	}
	data, _ := yaml.Marshal(c)
	if err := os.WriteFile(filepath.Join(cfgDir, "config.yaml"), data, 0644); err != nil {
		t.Fatalf("write config failed: %v", err)
	}

	got, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if got.DefaultTemplate != c.DefaultTemplate {
		t.Fatalf("default template mismatch: got %s want %s", got.DefaultTemplate, c.DefaultTemplate)
	}
}
