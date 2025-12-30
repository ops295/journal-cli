package todo

import (
	"os"
	"time"

	"journal-cli/internal/domain"
	"journal-cli/internal/fs"
	"journal-cli/internal/markdown"
)

// GetBacklog reads the journal entry from the given path and returns unchecked todos.
func GetBacklog(path string) ([]domain.Todo, error) {
	if !fs.Exists(path) {
		return []domain.Todo{}, nil
	}

	content, err := fs.ReadFile(path)
	if err != nil {
		return nil, err
	}

	entry, err := markdown.ParseMarkdown(content)
	if err != nil {
		return nil, err
	}

	var backlog []domain.Todo

	// Collect unchecked todos from "Todos" section
	for _, todo := range entry.Todos {
		if !todo.Done {
			backlog = append(backlog, todo)
		}
	}

	// Collect unchecked todos from "Backlog" section (recursive backlog)
	for _, todo := range entry.Backlog {
		if !todo.Done {
			backlog = append(backlog, todo)
		}
	}

	return backlog, nil
}

// GetPreviousJournalPath calculates the path for the previous day's journal.
// This is a helper, but the actual path construction depends on config.
// So maybe we just pass the date and let the caller construct the path?
// Or we pass the base dir.
func GetPreviousJournalPath(baseDir string, date time.Time) string {
	prevDate := date.AddDate(0, 0, -1)
	filename := prevDate.Format("2006-01-02") + ".md"
	// Assuming flat structure as per prompt: <ObsidianVault>/Journal/Daily/YYYY-MM-DD.md
	// But the prompt says "Stored in <ObsidianVault>/Journal/Daily/"
	// So we just join baseDir with filename.
	return baseDir + string(os.PathSeparator) + filename
}
