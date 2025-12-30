package domain

import (
	"testing"
	"time"
)

func TestNewJournalEntry(t *testing.T) {
	date := time.Date(2025, 12, 30, 0, 0, 0, 0, time.UTC)
	templateName := "daily-human-dev"

	entry := NewJournalEntry(date, templateName)

	if entry == nil {
		t.Fatal("NewJournalEntry() returned nil")
	}

	if !entry.Date.Equal(date) {
		t.Errorf("Date = %v, want %v", entry.Date, date)
	}

	if entry.Template != templateName {
		t.Errorf("Template = %v, want %v", entry.Template, templateName)
	}

	if entry.Todos == nil {
		t.Error("Todos should be initialized, got nil")
	}

	if entry.Backlog == nil {
		t.Error("Backlog should be initialized, got nil")
	}

	if entry.Questions == nil {
		t.Error("Questions should be initialized, got nil")
	}

	if len(entry.Todos) != 0 {
		t.Errorf("Todos should be empty, got %d items", len(entry.Todos))
	}

	if len(entry.Backlog) != 0 {
		t.Errorf("Backlog should be empty, got %d items", len(entry.Backlog))
	}

	if len(entry.Questions) != 0 {
		t.Errorf("Questions should be empty, got %d items", len(entry.Questions))
	}
}

func TestTodoStruct(t *testing.T) {
	tests := []struct {
		name string
		todo Todo
		want Todo
	}{
		{
			name: "unchecked todo",
			todo: Todo{Text: "Write tests", Done: false},
			want: Todo{Text: "Write tests", Done: false},
		},
		{
			name: "checked todo",
			todo: Todo{Text: "Build project", Done: true},
			want: Todo{Text: "Build project", Done: true},
		},
		{
			name: "empty todo",
			todo: Todo{Text: "", Done: false},
			want: Todo{Text: "", Done: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.todo.Text != tt.want.Text {
				t.Errorf("Text = %v, want %v", tt.todo.Text, tt.want.Text)
			}
			if tt.todo.Done != tt.want.Done {
				t.Errorf("Done = %v, want %v", tt.todo.Done, tt.want.Done)
			}
		})
	}
}

func TestJournalEntryFields(t *testing.T) {
	date := time.Date(2025, 12, 30, 0, 0, 0, 0, time.UTC)
	entry := NewJournalEntry(date, "test-template")

	// Test setting fields
	entry.Mood = "Happy"
	entry.Energy = "High"
	entry.Highlight = "Completed the project"

	if entry.Mood != "Happy" {
		t.Errorf("Mood = %v, want Happy", entry.Mood)
	}

	if entry.Energy != "High" {
		t.Errorf("Energy = %v, want High", entry.Energy)
	}

	if entry.Highlight != "Completed the project" {
		t.Errorf("Highlight = %v, want 'Completed the project'", entry.Highlight)
	}
}

func TestJournalEntryTodos(t *testing.T) {
	date := time.Date(2025, 12, 30, 0, 0, 0, 0, time.UTC)
	entry := NewJournalEntry(date, "test-template")

	// Add todos
	entry.Todos = append(entry.Todos, Todo{Text: "Task 1", Done: false})
	entry.Todos = append(entry.Todos, Todo{Text: "Task 2", Done: true})

	if len(entry.Todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(entry.Todos))
	}

	if entry.Todos[0].Text != "Task 1" {
		t.Errorf("First todo text = %v, want 'Task 1'", entry.Todos[0].Text)
	}

	if entry.Todos[1].Done != true {
		t.Errorf("Second todo should be done")
	}
}

func TestJournalEntryBacklog(t *testing.T) {
	date := time.Date(2025, 12, 30, 0, 0, 0, 0, time.UTC)
	entry := NewJournalEntry(date, "test-template")

	// Add backlog items
	entry.Backlog = append(entry.Backlog, Todo{Text: "Backlog 1", Done: false})
	entry.Backlog = append(entry.Backlog, Todo{Text: "Backlog 2", Done: false})

	if len(entry.Backlog) != 2 {
		t.Errorf("Expected 2 backlog items, got %d", len(entry.Backlog))
	}
}

func TestJournalEntryQuestions(t *testing.T) {
	date := time.Date(2025, 12, 30, 0, 0, 0, 0, time.UTC)
	entry := NewJournalEntry(date, "test-template")

	// Add questions and answers
	entry.Questions["How am I feeling?"] = "Great!"
	entry.Questions["What did I learn?"] = "Testing is important"

	if len(entry.Questions) != 2 {
		t.Errorf("Expected 2 questions, got %d", len(entry.Questions))
	}

	if entry.Questions["How am I feeling?"] != "Great!" {
		t.Errorf("Question answer mismatch")
	}
}
