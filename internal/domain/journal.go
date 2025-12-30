package domain

import "time"

type Todo struct {
	Text string
	Done bool
}

type JournalEntry struct {
	Date      time.Time
	Template  string
	Mood      string
	Energy    string
	Highlight string
	Todos     []Todo
	Backlog   []Todo
	Questions map[string]string // Question -> Answer
}

func NewJournalEntry(date time.Time, templateName string) *JournalEntry {
	return &JournalEntry{
		Date:      date,
		Template:  templateName,
		Todos:     make([]Todo, 0),
		Backlog:   make([]Todo, 0),
		Questions: make(map[string]string),
	}
}
