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
	"journal-cli/internal/stats"
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
	var entry *domain.JournalEntry

	// If today's journal file exists, load it and start in edit mode
	editFields := false
	if fs.Exists(todayFile) {
		data, err := fs.ReadFile(todayFile)
		if err != nil {
			fmt.Printf("Warning: could not read today's file: %v\n", err)
			entry = domain.NewJournalEntry(now, "")
			entry.Backlog = backlog
		} else {
			parsed, err := markdown.ParseMarkdown(data)
			if err != nil {
				fmt.Printf("Warning: could not parse today's file, starting fresh: %v\n", err)
				entry = domain.NewJournalEntry(now, "")
				entry.Backlog = backlog
			} else {
				// Use parsed entry as starting point
				entry = parsed
				// Ensure Backlog from yesterday is present too (merge if missing)
				if len(entry.Backlog) == 0 {
					entry.Backlog = backlog
				}
			}
		}
	} else {
		entry = domain.NewJournalEntry(now, "")
		entry.Backlog = backlog
	}

	// 6. Stats
	s, err := stats.GetStats(journalDir)
	if err != nil {
		fmt.Printf("Warning: could not calculate stats: %v\n", err)
	}

	// 7. Initialize TUI
	// If today's file existed and was parsed, prompt the user whether to edit it
	// or start fresh. Offer an option to edit fields (mood/energy/highlight) directly.
	if fs.Exists(todayFile) {
		// Prompt the user
		fmt.Printf("Today's journal exists at %s.\n", todayFile)
		fmt.Printf("[Enter] Edit full entry  |  f Edit mood/energy/highlight  |  n Start fresh\n")
		fmt.Printf("Choose an option: ")
		var resp string
		_, err := fmt.Scanln(&resp)
		if err != nil {
			// Treat empty input (enter) as default edit
			resp = ""
		}

		switch resp {
		case "n", "N":
			// Start fresh: override parsed entry with a new one but keep backlog
			entry = domain.NewJournalEntry(now, "")
			entry.Backlog = backlog
		case "f", "F":
			// Edit fields: ensure entry is used but start at Mood input
			editFields = true
		default:
			// Default: edit full entry (do nothing)
		}
	}

	model := tui.NewModel(cfg, templates, entry, s)

	// If we loaded an existing entry (from today's file), initialize the UI
	// so user can edit rather than starting a fresh flow.
	if entry != nil && entry.Template != "" {
		// select template in list
		for i, t := range templates {
			if t.Name == entry.Template {
				model.TemplateCursor = i
				break
			}
		}

		// populate inputs with existing values
		model.MoodInput.SetValue(entry.Mood)
		model.EnergyInput.SetValue(entry.Energy)
		model.HighlightInput.SetValue(entry.Highlight)

		// If user chose to edit fields, force start at Mood
		if editFields {
			model.CurrentStep = tui.StepMood
			model.MoodInput.Focus()
		} else {
			// Determine starting step: find first missing field
			switch {
			case entry.Template == "":
				model.CurrentStep = tui.StepSelectTemplate
				// focus handled by NewModel defaults
			case entry.Mood == "":
				model.CurrentStep = tui.StepMood
				model.MoodInput.Focus()
			case entry.Energy == "":
				model.CurrentStep = tui.StepEnergy
				model.EnergyInput.Focus()
			case entry.Highlight == "":
				model.CurrentStep = tui.StepHighlight
				model.HighlightInput.Focus()
			default:
				// If todos/questions have missing data, jump accordingly
				if len(entry.Todos) == 0 && len(entry.Backlog) >= 0 {
					model.CurrentStep = tui.StepTodos
					// Start with an explicit Todos menu to avoid navigation deadlocks
					model.TodosMenuActive = true
					model.TodoInput.Blur()
				} else {
					// find first unanswered question
					if len(templates) > 0 {
						tmpl := templates[model.TemplateCursor]
						qi := 0
						foundUnanswered := false
						for i, q := range tmpl.Questions {
							if _, ok := entry.Questions[q.Title]; !ok || entry.Questions[q.Title] == "" {
								qi = i
								foundUnanswered = true
								break
							}
						}
						if foundUnanswered {
							model.CurrentStep = tui.StepQuestions
							model.QuestionIndex = qi
							q := tmpl.Questions[qi].Title
							model.QuestionInput.SetValue(entry.Questions[q])
							model.QuestionInput.Focus()
						} else {
							model.CurrentStep = tui.StepDone
						}
					} else {
						model.CurrentStep = tui.StepTodos
						// Start with Todos menu active when opening existing entry
						model.TodosMenuActive = true
						model.TodoInput.Blur()
					}
				}
			}
		}
	}

	p := tea.NewProgram(model)

	// 8. Run TUI
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

	fmt.Printf("Journal entry saved to: %s\n", todayFile)
	fmt.Printf("To view:  cat \"%s\"\n", todayFile)
	fmt.Printf("To edit:  nano \"%s\"\n", todayFile)
}
