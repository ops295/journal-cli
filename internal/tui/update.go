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
			// We'll treat backlog items and added todos as a single linear selectable list.
			// Backlog items come first, then Entry.Todos. The BacklogCursor indexes into that combined list.
			totalSelectable := len(m.Entry.Backlog) + len(m.Entry.Todos)

			switch msg := msg.(type) {
			case tea.KeyMsg:
				// Use Tab/Shift+Tab to switch focus between input and the selectable list.
				if msg.Type == tea.KeyTab && m.TodoInput.Focused() && totalSelectable > 0 {
					m.TodoInput.Blur()
					// focus on first selectable item
					m.BacklogCursor = 0
					return m, nil
				}

				if msg.Type == tea.KeyShiftTab && !m.TodoInput.Focused() {
					// shift+tab moves focus back to input
					m.TodoInput.Focus()
					return m, nil
				}

				// When not focused, Up/Down (or k/j) navigate the linear selection; Tab/Down past end focuses input
				if !m.TodoInput.Focused() {
					if msg.String() == "up" || msg.String() == "k" {
						if m.BacklogCursor > 0 {
							m.BacklogCursor--
						}
						return m, nil
					}
					if msg.String() == "down" || msg.String() == "j" {
						if m.BacklogCursor < totalSelectable-1 {
							m.BacklogCursor++
						} else {
							m.TodoInput.Focus()
						}
						return m, nil
					}

					// Space toggles backlog selection only (backlog indices are 0..len(backlog)-1)
					if msg.String() == " " {
						if m.BacklogCursor < len(m.Entry.Backlog) {
							if m.SelectedBacklog[m.BacklogCursor] {
								delete(m.SelectedBacklog, m.BacklogCursor)
							} else {
								m.SelectedBacklog[m.BacklogCursor] = true
							}
						}
						return m, nil
					}

					// Enter on an added todo opens it for editing
					if msg.String() == "enter" {
						if m.BacklogCursor >= len(m.Entry.Backlog) {
							idx := m.BacklogCursor - len(m.Entry.Backlog)
							if idx >= 0 && idx < len(m.Entry.Todos) {
								// Load the todo into input for editing, remove from list temporarily
								val := m.Entry.Todos[idx].Text
								// remove from slice
								m.Entry.Todos = append(m.Entry.Todos[:idx], m.Entry.Todos[idx+1:]...)
								m.TodoInput.SetValue(val)
								m.TodoInput.Focus()
								return m, nil
							}
						}
					}
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
						// Add todo (or re-add edited todo)
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
			// Shift+Tab => Next question (save current)
			if msg.Type == tea.KeyShiftTab && !msg.Alt {
				currentTemplate := m.Templates[m.TemplateCursor]
				question := currentTemplate.Questions[m.QuestionIndex].Title
				m.Entry.Questions[question] = m.QuestionInput.Value()
				m.QuestionInput.Reset()
				m.QuestionIndex++
				if m.QuestionIndex >= len(currentTemplate.Questions) {
					m.CurrentStep = StepDone
					return m, tea.Quit
				}
				nextQ := currentTemplate.Questions[m.QuestionIndex].Title
				m.QuestionInput.SetValue(m.Entry.Questions[nextQ])
				m.QuestionInput.Focus()
				return m, nil
			}

			// Shift+Left => Previous question
			if msg.String() == "shift+left" {
				if m.QuestionIndex > 0 {
					currentTemplate := m.Templates[m.TemplateCursor]
					question := currentTemplate.Questions[m.QuestionIndex].Title
					m.Entry.Questions[question] = m.QuestionInput.Value()
					m.QuestionIndex--
					prevQ := currentTemplate.Questions[m.QuestionIndex].Title
					m.QuestionInput.SetValue(m.Entry.Questions[prevQ])
					m.QuestionInput.Focus()
				}
				return m, nil
			}

			// Shift+Right => Next question
			if msg.String() == "shift+right" {
				currentTemplate := m.Templates[m.TemplateCursor]
				question := currentTemplate.Questions[m.QuestionIndex].Title
				m.Entry.Questions[question] = m.QuestionInput.Value()
				m.QuestionInput.Reset()
				m.QuestionIndex++
				if m.QuestionIndex >= len(currentTemplate.Questions) {
					m.CurrentStep = StepDone
					return m, tea.Quit
				}
				nextQ := currentTemplate.Questions[m.QuestionIndex].Title
				m.QuestionInput.SetValue(m.Entry.Questions[nextQ])
				m.QuestionInput.Focus()
				return m, nil
			}

			// Enter => save current and advance
			if msg.Type == tea.KeyEnter {
				currentTemplate := m.Templates[m.TemplateCursor]
				question := currentTemplate.Questions[m.QuestionIndex].Title
				m.Entry.Questions[question] = m.QuestionInput.Value()
				m.QuestionInput.Reset()
				m.QuestionIndex++
				if m.QuestionIndex >= len(currentTemplate.Questions) {
					m.CurrentStep = StepDone
					return m, tea.Quit
				}
				nextQ := currentTemplate.Questions[m.QuestionIndex].Title
				m.QuestionInput.SetValue(m.Entry.Questions[nextQ])
				m.QuestionInput.Focus()
				return m, nil
			}

			// Backwards compatible: Ctrl+S or Ctrl+N to submit answer
			if msg.Type == tea.KeyCtrlS || msg.Type == tea.KeyCtrlN {
				currentTemplate := m.Templates[m.TemplateCursor]
				question := currentTemplate.Questions[m.QuestionIndex].Title
				m.Entry.Questions[question] = m.QuestionInput.Value()
				m.QuestionInput.Reset()
				m.QuestionIndex++
				if m.QuestionIndex >= len(currentTemplate.Questions) {
					m.CurrentStep = StepDone
					return m, tea.Quit
				}
				nextQ := currentTemplate.Questions[m.QuestionIndex].Title
				m.QuestionInput.SetValue(m.Entry.Questions[nextQ])
				m.QuestionInput.Focus()
				return m, nil
			}
		}
		return m, cmd
	}

	return m, nil
}
