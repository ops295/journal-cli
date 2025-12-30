package markdown

import (
    "strings"
    "testing"
    "time"

    "journal-cli/internal/domain"
)

func TestGenerateAndParseRoundtrip(t *testing.T) {
    date := time.Date(2025, 12, 30, 0, 0, 0, 0, time.UTC)
    entry := domain.NewJournalEntry(date, "daily-human-dev")
    entry.Mood = "Calm"
    entry.Energy = "Medium"
    entry.Highlight = "Wrote tests"
    entry.Todos = append(entry.Todos, domain.Todo{Text: "Do thing", Done: false})
    entry.Backlog = append(entry.Backlog, domain.Todo{Text: "Carryover", Done: false})
    entry.Questions["What did I learn?"] = "Testing roundtrip"

    md, err := GenerateMarkdown(entry)
    if err != nil {
        t.Fatalf("GenerateMarkdown error: %v", err)
    }

    parsed, err := ParseMarkdown(md)
    if err != nil {
        t.Fatalf("ParseMarkdown error: %v", err)
    }

    if !parsed.Date.Equal(entry.Date) {
        t.Fatalf("date mismatch: got %v want %v", parsed.Date, entry.Date)
    }

    if parsed.Template != entry.Template {
        t.Fatalf("template mismatch: got %s want %s", parsed.Template, entry.Template)
    }

    if len(parsed.Todos) != 1 || parsed.Todos[0].Text != "Do thing" {
        t.Fatalf("todos mismatch: %v", parsed.Todos)
    }

    if len(parsed.Backlog) != 1 || parsed.Backlog[0].Text != "Carryover" {
        t.Fatalf("backlog mismatch: %v", parsed.Backlog)
    }

    // Keys may include emoji prefixes; find answer by substring match
    found := false
    for k, v := range parsed.Questions {
        if strings.Contains(k, "What did I learn") && v == "Testing roundtrip" {
            found = true
            break
        }
    }
    if !found {
        t.Fatalf("questions mismatch: %v", parsed.Questions)
    }
}

func TestParseMalformed(t *testing.T) {
    bad := []byte("no-frontmatter-here")
    if _, err := ParseMarkdown(bad); err == nil {
        t.Fatalf("expected error parsing malformed markdown")
    }
}
