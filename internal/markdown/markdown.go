package markdown

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"time"

	"journal-cli/internal/domain"

	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	Date      string `yaml:"date"`
	Template  string `yaml:"template"`
	Mood      string `yaml:"mood"`
	Energy    string `yaml:"energy"`
	Highlight string `yaml:"highlight"`
}

func GenerateMarkdown(entry *domain.JournalEntry) ([]byte, error) {
	// Frontmatter
	fm := FrontMatter{
		Date:      entry.Date.Format("2006-01-02"),
		Template:  entry.Template,
		Mood:      entry.Mood,
		Energy:    entry.Energy,
		Highlight: entry.Highlight,
	}
	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return nil, err
	}

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.Write(fmBytes)
	sb.WriteString("---\n\n")

	sb.WriteString(fmt.Sprintf("# Daily Journal – %s\n\n", entry.Date.Format("2006-01-02")))

	if entry.Highlight != "" {
		sb.WriteString(fmt.Sprintf("## ⭐️ Daily Highlight\n%s\n\n", entry.Highlight))
	}

	sb.WriteString("## ✅ Todos – Today\n")
	for _, todo := range entry.Todos {
		check := " "
		if todo.Done {
			check = "x"
		}
		sb.WriteString(fmt.Sprintf("- [%s] %s\n", check, todo.Text))
	}
	sb.WriteString("\n")

	if len(entry.Backlog) > 0 {
		// Backlog usually comes from previous day, so we might want to label it differently
		// But for the current day's file, it's just backlog items that were carried over
		// The prompt says "Backlog (from YYYY-MM-DD)" but that date changes.
		// For simplicity, we'll just list them.
		// Or maybe we should track where they came from?
		// The prompt says: "## 🔁 Backlog (from YYYY-MM-DD)"
		// We'll assume for now we just list them under a generic Backlog or we need to know the source date.
		// Let's just use "Backlog" for now or "Backlog (Previous)"
		sb.WriteString("## 🔁 Backlog\n")
		for _, todo := range entry.Backlog {
			check := " "
			if todo.Done {
				check = "x"
			}
			sb.WriteString(fmt.Sprintf("- [%s] %s\n", check, todo.Text))
		}
		sb.WriteString("\n")
	}

	for q, a := range entry.Questions {
		sb.WriteString(fmt.Sprintf("## 🧠 %s\n", q))
		sb.WriteString(fmt.Sprintf("%s\n\n", a))
	}

	return []byte(sb.String()), nil
}

func ParseMarkdown(content []byte) (*domain.JournalEntry, error) {
	// Split frontmatter and body
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid markdown format: missing frontmatter")
	}

	// Parse Frontmatter
	var fm FrontMatter
	if err := yaml.Unmarshal(parts[1], &fm); err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02", fm.Date)
	if err != nil {
		return nil, err
	}

	entry := domain.NewJournalEntry(date, fm.Template)
	entry.Mood = fm.Mood
	entry.Energy = fm.Energy
	entry.Highlight = fm.Highlight

	// Parse Body
	scanner := bufio.NewScanner(bytes.NewReader(parts[2]))
	var currentSection string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "## ") {
			currentSection = strings.TrimPrefix(line, "## ")
			continue
		}

		if strings.HasPrefix(line, "- [") {
			// Todo item
			done := strings.HasPrefix(line, "- [x]") || strings.HasPrefix(line, "- [X]")
			text := strings.TrimSpace(line[5:])
			todo := domain.Todo{Text: text, Done: done}

			if strings.Contains(currentSection, "Todos") {
				entry.Todos = append(entry.Todos, todo)
			} else if strings.Contains(currentSection, "Backlog") {
				entry.Backlog = append(entry.Backlog, todo)
			}
		} else {
			// Probably an answer to a question
			// This is a simple parser, assuming questions are headers and answers follow
			// We need to match questions from the template or just store them as found
			// For now, let's assume any other section is a question
			if currentSection != "" && !strings.Contains(currentSection, "Todos") && !strings.Contains(currentSection, "Backlog") {
				// Clean up the section name (remove emoji if present)
				// This is a bit hacky, but works for the MVP
				q := currentSection
				if idx := strings.Index(q, " "); idx != -1 {
					// q = q[idx+1:] // Remove emoji? Maybe not, keep it for now
				}
				// Append to existing answer if multi-line?
				if val, ok := entry.Questions[q]; ok {
					entry.Questions[q] = val + "\n" + line
				} else {
					entry.Questions[q] = line
				}
			}
		}
	}

	return entry, nil
}
