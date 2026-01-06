package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTemplatesDefaults(t *testing.T) {
	tmp := t.TempDir()
	prev := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmp)
	defer os.Setenv("XDG_CONFIG_HOME", prev)

	_, err := LoadTemplates()
	if err != nil {
		t.Fatalf("LoadTemplates error: %v", err)
	}

	// It's acceptable (depending on environment) for templates to be empty.
}

func TestLoadTemplatesWithMalformedFile(t *testing.T) {
	tmp := t.TempDir()
	prev := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmp)
	defer os.Setenv("XDG_CONFIG_HOME", prev)

	templatesDir := filepath.Join(tmp, "journal-cli", "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	// Write a malformed template file
	bad := "not: : valid: yaml: :::"
	if err := os.WriteFile(filepath.Join(templatesDir, "bad.yaml"), []byte(bad), 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	_, err := LoadTemplates()
	if err != nil {
		t.Fatalf("LoadTemplates error: %v", err)
	}

	// The loader should succeed even if templates are malformed; no further guarantees.
}

func TestListTemplateNames(t *testing.T) {
	tmp := t.TempDir()
	prev := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmp)
	defer os.Setenv("XDG_CONFIG_HOME", prev)

	names, err := ListTemplateNames()
	if err != nil {
		t.Fatalf("ListTemplateNames error: %v", err)
	}

	// Should have default templates
	if len(names) == 0 {
		t.Fatalf("expected at least one template, got none")
	}
}

func TestTemplateExists(t *testing.T) {
	tmp := t.TempDir()
	prev := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmp)
	defer os.Setenv("XDG_CONFIG_HOME", prev)

	// Load templates first to ensure defaults are created
	_, err := LoadTemplates()
	if err != nil {
		t.Fatalf("LoadTemplates error: %v", err)
	}

	// Test existing template
	exists, err := TemplateExists("daily-human-dev")
	if err != nil {
		t.Fatalf("TemplateExists error: %v", err)
	}
	if !exists {
		t.Fatalf("expected daily-human-dev to exist")
	}

	// Test non-existing template
	exists, err = TemplateExists("nonexistent-template")
	if err != nil {
		t.Fatalf("TemplateExists error: %v", err)
	}
	if exists {
		t.Fatalf("expected nonexistent-template to not exist")
	}
}

func TestGetTemplateByName(t *testing.T) {
	tmp := t.TempDir()
	prev := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmp)
	defer os.Setenv("XDG_CONFIG_HOME", prev)

	// Load templates first
	_, err := LoadTemplates()
	if err != nil {
		t.Fatalf("LoadTemplates error: %v", err)
	}

	// Test getting existing template
	tmpl, err := GetTemplateByName("daily-human-dev")
	if err != nil {
		t.Fatalf("GetTemplateByName error: %v", err)
	}
	if tmpl == nil {
		t.Fatalf("expected template, got nil")
	}
	if tmpl.Name != "daily-human-dev" {
		t.Fatalf("template name mismatch: got %s want daily-human-dev", tmpl.Name)
	}

	// Test getting non-existing template
	tmpl, err = GetTemplateByName("nonexistent")
	if err != nil {
		t.Fatalf("GetTemplateByName error: %v", err)
	}
	if tmpl != nil {
		t.Fatalf("expected nil for nonexistent template, got %v", tmpl)
	}
}
