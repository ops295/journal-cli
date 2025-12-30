package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"journal-cli/internal/config"
	"journal-cli/internal/domain"
	"journal-cli/internal/fs"
	"journal-cli/internal/markdown"
	"journal-cli/internal/template"
	"journal-cli/internal/todo"
	"journal-cli/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// 2. Load Templates
	templates, err := template.LoadTemplates()
	if err != nil {
		fmt.Printf("Error loading templates: %v\n", err)
		os.Exit(1)
	}

	// 3. Setup Date and Paths
	now := time.Now()
	journalDir := cfg.JournalDir
	if cfg.ObsidianVault != "" {
		journalDir = filepath.Join(cfg.ObsidianVault, cfg.JournalDir)
	} else {
		// Fallback if not set, though LoadConfig sets a default relative path
		// But we need a base. Let's assume current directory if Vault is empty?
		// Or better, warn user.
		// For MVP, if Vault is empty, we might just use current directory or a default in UserHome.
		// Let's use a default in UserHome/Documents/Journal if not set?
		// Or just fail?
		// The prompt says "Stored in <ObsidianVault>/Journal/Daily/"
		// If ObsidianVault is not configured, we can't really guess.
		// But for now, let's just use a local "journal" folder if not set.
		if journalDir == "Journal/Daily" { // Default value from config.go
			// Try to find a good place
			home, _ := os.UserHomeDir()
			journalDir = filepath.Join(home, "Documents", "Journal", "Daily")
		}
	}

	if err := fs.EnsureDir(journalDir); err != nil {
		fmt.Printf("Error ensuring journal directory: %v\n", err)
		os.Exit(1)
	}

	todayFile := filepath.Join(journalDir, now.Format("2006-01-02")+".md")

	// 4. Load Backlog
	yesterdayFile := todo.GetPreviousJournalPath(journalDir, now)
	backlog, err := todo.GetBacklog(yesterdayFile)
	if err != nil {
		// Non-fatal, just log or ignore
		fmt.Printf("Warning: could not load backlog: %v\n", err)
	}

	// 5. Initialize Entry
	entry := domain.NewJournalEntry(now, "")
	entry.Backlog = backlog

	// 6. Initialize TUI
	model := tui.NewModel(cfg, templates, entry)
	p := tea.NewProgram(model)

	// 7. Run TUI
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}

	m, ok := finalModel.(tui.Model)
	if !ok {
		fmt.Printf("Error: unexpected model type\n")
		os.Exit(1)
	}

	if m.CurrentStep != tui.StepDone {
		fmt.Println("Journaling cancelled.")
		return
	}

	// Merge selected backlog items into Todos
	// And keep unselected ones in Backlog (or remove them from Backlog if we want to drop them?
	// Usually backlog persists until done or deleted.
	// For this MVP, let's say:
	// - Selected items move to Today's Todos (and are removed from Backlog for *this* entry, effectively)
	// - Unselected items remain in Backlog for *this* entry.
	// When generating markdown:
	// - Todos section has new items + selected backlog items.
	// - Backlog section has unselected backlog items.
	
	var newTodos []domain.Todo
	var remainingBacklog []domain.Todo

	// Add selected backlog items first
	for i, t := range m.Entry.Backlog {
		if m.SelectedBacklog[i] {
			newTodos = append(newTodos, t)
		} else {
			remainingBacklog = append(remainingBacklog, t)
		}
	}
	// Append newly added todos
	newTodos = append(newTodos, m.Entry.Todos...)
	
	m.Entry.Todos = newTodos
	m.Entry.Backlog = remainingBacklog

	// 8. Save to Disk
	content, err := markdown.GenerateMarkdown(m.Entry)
	if err != nil {
		fmt.Printf("Error generating markdown: %v\n", err)
		os.Exit(1)
	}

	if err := fs.WriteFile(todayFile, content); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Journal entry saved to %s\n", todayFile)
}
