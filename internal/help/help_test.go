package help

import (
	"strings"
	"testing"
)

func TestLoadHelp(t *testing.T) {
	doc, err := LoadHelp()
	if err != nil {
		t.Fatalf("LoadHelp error: %v", err)
	}

	if doc.App.Name != "journal-cli" {
		t.Errorf("App name mismatch: got %s want journal-cli", doc.App.Name)
	}

	if len(doc.Commands) == 0 {
		t.Error("Expected commands to be loaded")
	}

	if doc.Configuration.Title == "" {
		t.Error("Expected configuration section to be loaded")
	}
}

func TestRender(t *testing.T) {
	doc, err := LoadHelp()
	if err != nil {
		t.Fatalf("LoadHelp error: %v", err)
	}

	output := doc.Render("journal")

	// Check for key sections
	expectedStrings := []string{
		"Usage: journal [options]",
		"Options:",
		"--help",
		"Configuration:",
		"Templates:",
		"Template Management:",
	}

	for _, s := range expectedStrings {
		if !strings.Contains(output, s) {
			t.Errorf("Expected output to contain %q", s)
		}
	}
}
