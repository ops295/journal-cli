package tui

import (
	"journal-cli/internal/config"
	"journal-cli/internal/domain"
	"journal-cli/internal/stats"
	"journal-cli/internal/template"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Step int

const (
	StepSelectTemplate Step = iota
	StepMood
	StepEnergy
	StepHighlight
	StepTodos
	StepQuestions
	StepDone
)

type Model struct {
	Config    *config.Config
	Templates []template.Template
	Entry     *domain.JournalEntry
	Stats     stats.Stats

	CurrentStep    Step
	TemplateCursor int
	QuestionIndex  int
	TodoInput      textinput.Model
	QuestionInput  textarea.Model
	MoodInput      textinput.Model
	EnergyInput    textinput.Model
	HighlightInput textinput.Model

	// For Todos
	BacklogCursor   int
	SelectedBacklog map[int]bool // Index in Entry.Backlog -> true if selected
	TodoMode        bool         // true if adding a todo, false if reviewing backlog

	// Todos menu when opening an existing entry to avoid navigation deadlocks
	TodosMenuActive bool
	TodosMenuCursor int

	Err error
}

func NewModel(cfg *config.Config, templates []template.Template, entry *domain.JournalEntry, s stats.Stats) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter a task..."
	ti.Focus()

	ta := textarea.New()
	ta.Placeholder = "Write your answer..."
	ta.SetHeight(5)

	mi := textinput.New()
	mi.Placeholder = "How are you feeling?"
	mi.Focus()

	ei := textinput.New()
	ei.Placeholder = "How is your energy?"
	ei.Focus()

	hi := textinput.New()
	hi.Placeholder = "What is your main focus today?"
	hi.Focus()

	return Model{
		Config:          cfg,
		Templates:       templates,
		Entry:           entry,
		Stats:           s,
		CurrentStep:     StepSelectTemplate,
		TodoInput:       ti,
		QuestionInput:   ta,
		MoodInput:       mi,
		EnergyInput:     ei,
		HighlightInput:  hi,
		SelectedBacklog: make(map[int]bool),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
