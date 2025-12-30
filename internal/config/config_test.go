package config

import (
    "os"
    "path/filepath"
    "testing"

    "gopkg.in/yaml.v3"
)

func TestLoadConfigDefaultWhenMissing(t *testing.T) {
    tmp := t.TempDir()
    // Set XDG_CONFIG_HOME for deterministic UserConfigDir resolution
    prev := os.Getenv("XDG_CONFIG_HOME")
    os.Setenv("XDG_CONFIG_HOME", tmp)
    defer os.Setenv("XDG_CONFIG_HOME", prev)

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
    prev := os.Getenv("XDG_CONFIG_HOME")
    os.Setenv("XDG_CONFIG_HOME", tmp)
    defer os.Setenv("XDG_CONFIG_HOME", prev)

    cfgDir := filepath.Join(tmp, "journal-cli")
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
