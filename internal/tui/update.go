package tui

import (
	"journal-cli/internal/domain"


	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	switch m.CurrentStep {
	case StepSelectTemplate:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if m.TemplateCursor > 0 {
					m.TemplateCursor--
				}
			case "down", "j":
				if m.TemplateCursor < len(m.Templates)-1 {
					m.TemplateCursor++
				}
			case "enter":
				m.Entry.Template = m.Templates[m.TemplateCursor].Name
				m.CurrentStep = StepMood
				m.MoodInput.Focus()
				return m, nil
			}
		}

	case StepMood:
		m.MoodInput, cmd = m.MoodInput.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter {
				m.Entry.Mood = m.MoodInput.Value()
				m.CurrentStep = StepEnergy
				m.EnergyInput.Focus()
				return m, nil
			}
		}
		return m, cmd

	case StepEnergy:
		m.EnergyInput, cmd = m.EnergyInput.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter {
				m.Entry.Energy = m.EnergyInput.Value()
				m.CurrentStep = StepHighlight
				m.HighlightInput.Focus()
				return m, nil
			}
		}
		return m, cmd

	case StepHighlight:
		m.HighlightInput, cmd = m.HighlightInput.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter {
				m.Entry.Highlight = m.HighlightInput.Value()
				m.CurrentStep = StepTodos
				// If we have backlog, start there? Or just focus input?
				// Let's focus input but allow navigating up to backlog.
				m.TodoInput.Focus()
				return m, nil
			}
		}
		return m, cmd

	case StepTodos:
		// If TodoInput is focused, handle that.
		// But we also want to navigate backlog.
		// Let's say: Up/Down navigates backlog if focused?
		// Or maybe: Tab switches focus between Input and Backlog?
		// Simpler: Up arrow from Input goes to Backlog. Down from Backlog goes to Input.
		
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "up" && m.TodoInput.Focused() && len(m.Entry.Backlog) > 0 {
				m.TodoInput.Blur()
				m.BacklogCursor = len(m.Entry.Backlog) - 1
				return m, nil
			}
			if msg.String() == "down" && !m.TodoInput.Focused() {
				if m.BacklogCursor < len(m.Entry.Backlog)-1 {
					m.BacklogCursor++
				} else {
					m.TodoInput.Focus()
				}
				return m, nil
			}
			if msg.String() == "up" && !m.TodoInput.Focused() {
				if m.BacklogCursor > 0 {
					m.BacklogCursor--
				}
				return m, nil
			}
			if msg.String() == " " && !m.TodoInput.Focused() {
				// Toggle selection
				if m.SelectedBacklog[m.BacklogCursor] {
					delete(m.SelectedBacklog, m.BacklogCursor)
				} else {
					m.SelectedBacklog[m.BacklogCursor] = true
				}
				return m, nil
			}
		}

		if m.TodoInput.Focused() {
			m.TodoInput, cmd = m.TodoInput.Update(msg)
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if msg.Type == tea.KeyEnter {
					val := m.TodoInput.Value()
					if val == "" {
						// Empty line means we are done with todos
						m.CurrentStep = StepQuestions
						m.QuestionIndex = 0
						m.QuestionInput.Focus()
						return m, nil
					}
					// Add todo
					m.Entry.Todos = append(m.Entry.Todos, domain.Todo{Text: val, Done: false})
					m.TodoInput.Reset()
				}
			}
			return m, cmd
		}

	case StepQuestions:
		m.QuestionInput, cmd = m.QuestionInput.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			// Use Ctrl+S or Ctrl+N to submit answer
			if msg.Type == tea.KeyCtrlS || msg.Type == tea.KeyCtrlN {
				// Save answer
				currentTemplate := m.Templates[m.TemplateCursor]
				question := currentTemplate.Questions[m.QuestionIndex].Title
				m.Entry.Questions[question] = m.QuestionInput.Value()
				
				m.QuestionInput.Reset()
				m.QuestionIndex++

				if m.QuestionIndex >= len(currentTemplate.Questions) {
					m.CurrentStep = StepDone
					return m, tea.Quit
				}
			}
		}
		return m, cmd
	}

	return m, nil
}
