package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"journal-cli/internal/config"
	"journal-cli/internal/fs"
	"journal-cli/internal/preview"
)

// PreviewEntry handles the logic for finding a journal file and opening it in the browser
func PreviewEntry(dateStr string) error {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// 2. Determine Journal Directory
	journalDir := cfg.JournalDir
	if cfg.ObsidianVault != "" {
		journalDir = filepath.Join(cfg.ObsidianVault, cfg.JournalDir)
	} else if journalDir == "Journal/Daily" {
		// Fallback similar to app.Run
		home, _ := os.UserHomeDir()
		journalDir = filepath.Join(home, "Documents", "Journal", "Daily")
	}

	// 3. Determine Date
	targetDate := time.Now()
	if dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return fmt.Errorf("invalid date format (use YYYY-MM-DD): %w", err)
		}
		targetDate = parsed
	}

	dateFmt := targetDate.Format("2006-01-02")
	targetFile := filepath.Join(journalDir, dateFmt+".md")

	if !fs.Exists(targetFile) {
		return fmt.Errorf("journal entry for %s not found at %s", dateFmt, targetFile)
	}

	// 4. Read File
	content, err := fs.ReadFile(targetFile)
	if err != nil {
		return fmt.Errorf("reading journal file: %w", err)
	}

	// 5. Render HTML
	htmlContent, err := preview.Render(content, "Journal - "+dateFmt)
	if err != nil {
		return fmt.Errorf("rendering html: %w", err)
	}

	// 6. Open in Browser
	fmt.Printf("🔍 Opening preview for %s...\n", dateFmt)
	if err := preview.OpenInBrowser(htmlContent); err != nil {
		return fmt.Errorf("opening browser: %w", err)
	}

	return nil
}
