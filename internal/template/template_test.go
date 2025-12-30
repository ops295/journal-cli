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
