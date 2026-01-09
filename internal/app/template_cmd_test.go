package app

import (
	"os"
	"path/filepath"
	"testing"

	"journal-cli/internal/config"
	"journal-cli/internal/template"

	"gopkg.in/yaml.v3"
)

func TestSetDefaultTemplate(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("APPDATA", tmp)

	// Ensure templates are loaded (creates defaults)
	_, err := template.LoadTemplates()
	if err != nil {
		t.Fatalf("LoadTemplates error: %v", err)
	}

	// Test setting valid template
	err = SetDefaultTemplate("daily-human-dev")
	if err != nil {
		t.Fatalf("SetDefaultTemplate error: %v", err)
	}

	// Verify it was saved
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if cfg.DefaultTemplate != "daily-human-dev" {
		t.Fatalf("default template mismatch: got %s want daily-human-dev", cfg.DefaultTemplate)
	}
}

func TestSetDefaultTemplateInvalid(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("APPDATA", tmp)

	// Ensure templates are loaded
	_, err := template.LoadTemplates()
	if err != nil {
		t.Fatalf("LoadTemplates error: %v", err)
	}

	// Test setting invalid template
	err = SetDefaultTemplate("nonexistent-template")
	if err == nil {
		t.Fatalf("expected error for nonexistent template, got nil")
	}
}

func TestGetDefaultTemplate(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("APPDATA", tmp)

	// Create config with default template
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("UserConfigDir error: %v", err)
	}

	cfgDir := filepath.Join(userConfigDir, "journal-cli")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	cfg := config.Config{
		ObsidianVault:   "/tmp/vault",
		JournalDir:      "Journal/Daily",
		DefaultTemplate: "gentle-day",
	}
	data, _ := yaml.Marshal(cfg)
	if err := os.WriteFile(filepath.Join(cfgDir, "config.yaml"), data, 0644); err != nil {
		t.Fatalf("write config failed: %v", err)
	}

	// Test getting default template
	defaultTmpl, err := GetDefaultTemplate()
	if err != nil {
		t.Fatalf("GetDefaultTemplate error: %v", err)
	}

	if defaultTmpl != "gentle-day" {
		t.Fatalf("default template mismatch: got %s want gentle-day", defaultTmpl)
	}
}

func TestListTemplatesCommand(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("APPDATA", tmp)

	// Ensure templates are loaded
	_, err := template.LoadTemplates()
	if err != nil {
		t.Fatalf("LoadTemplates error: %v", err)
	}

	// Test listing templates (should not error)
	err = ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates error: %v", err)
	}
}
