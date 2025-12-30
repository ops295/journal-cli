package app

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"journal-cli/internal/config"
	"journal-cli/internal/domain"
	"journal-cli/internal/fs"
	"journal-cli/internal/markdown"
)

// UpdateTodos loads the journal file for the given date (empty = today)
// and prompts the user for each todo: complete (c), partial (p), not yet (n).
// 'not yet' items are moved to the Backlog section so they'll be carried
// forward when the next day's journal is opened.
func UpdateTodos(dateStr string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Determine journal dir similar to app.Run
	journalDir := cfg.JournalDir
	if cfg.ObsidianVault != "" {
		journalDir = filepath.Join(cfg.ObsidianVault, cfg.JournalDir)
	} else {
		if journalDir == "Journal/Daily" {
			home, _ := os.UserHomeDir()
			journalDir = filepath.Join(home, "Documents", "Journal", "Daily")
		}
	}

	var date time.Time
	if dateStr == "" {
		date = time.Now()
	} else {
		d, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return fmt.Errorf("invalid date format (use YYYY-MM-DD): %w", err)
		}
		date = d
	}

	file := filepath.Join(journalDir, date.Format("2006-01-02")+".md")
	if !fs.Exists(file) {
		return fmt.Errorf("journal file not found: %s", file)
	}

	data, err := fs.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	entry, err := markdown.ParseMarkdown(data)
	if err != nil {
		return fmt.Errorf("parse markdown: %w", err)
	}

	if len(entry.Todos) == 0 {
		fmt.Println("No todos found in the entry.")
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Updating todos in %s\n", file)

	// iterate over todos, allow removing while iterating
	for i := 0; i < len(entry.Todos); i++ {
		t := entry.Todos[i]
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}
		fmt.Printf("%d) %s %s\n", i+1, status, t.Text)
		fmt.Printf("(c)omplete, (p)artial, (n)ot yet -> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		switch input {
		case "c", "complete":
			entry.Todos[i].Done = true
		case "p", "partial":
			// mark as partial: keep in todos but append marker
			if !strings.Contains(entry.Todos[i].Text, "(partial)") {
				entry.Todos[i].Text = entry.Todos[i].Text + " (partial)"
			}
		case "n", "not":
			// move to backlog: append to Backlog and remove from Todos
			entry.Backlog = append(entry.Backlog, domain.Todo{Text: entry.Todos[i].Text, Done: false})
			// remove this todo
			entry.Todos = append(entry.Todos[:i], entry.Todos[i+1:]...)
			i-- // stay at same index
		default:
			// treat as skip/no change
			fmt.Println("skipped")
		}
	}

	// Generate markdown and write back
	out, err := markdown.GenerateMarkdown(entry)
	if err != nil {
		return fmt.Errorf("generate markdown: %w", err)
	}

	if err := fs.WriteFile(file, out); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	fmt.Printf("Updated file: %s\n", file)
	return nil
}
