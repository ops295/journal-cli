package tui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	var s strings.Builder

	if m.Err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.Err))
	}

	switch m.CurrentStep {
	case StepSelectTemplate:
		s.WriteString(titleStyle.Render("Select Template"))
		s.WriteString("\n\n")
		for i, t := range m.Templates {
			cursor := " "
			style := itemStyle
			if m.TemplateCursor == i {
				cursor = ">"
				style = selectedItemStyle
			}
			s.WriteString(style.Render(fmt.Sprintf("%s %s", cursor, t.Name)) + "\n")
			if m.TemplateCursor == i && t.Description != "" {
				s.WriteString(itemStyle.Render(fmt.Sprintf("    %s", t.Description)) + "\n")
			}
		}
		s.WriteString("\n(Use arrow keys to select, Enter to confirm)")

	case StepMood:
		s.WriteString(titleStyle.Render("How are you feeling?"))
		s.WriteString("\n\n")
		s.WriteString(m.MoodInput.View())
		s.WriteString("\n\n(Enter to continue)")

	case StepEnergy:
		s.WriteString(titleStyle.Render("How is your energy?"))
		s.WriteString("\n\n")
		s.WriteString(m.EnergyInput.View())
		s.WriteString("\n\n(Enter to continue)")

	case StepHighlight:
		s.WriteString(titleStyle.Render("Daily Highlight"))
		s.WriteString("\n\n")
		s.WriteString(m.HighlightInput.View())
		s.WriteString("\n\n(Enter to continue)")

	case StepTodos:
		s.WriteString(titleStyle.Render("Today's Todos"))
		s.WriteString("\n\n")
		
		// Backlog
		if len(m.Entry.Backlog) > 0 {
			s.WriteString("🔁 Backlog (Up/Down to select, Space to toggle):\n")
			for i, t := range m.Entry.Backlog {
				cursor := " "
				if !m.TodoInput.Focused() && m.BacklogCursor == i {
					cursor = ">"
				}
				checked := "[ ]"
				if m.SelectedBacklog[i] {
					checked = "[x]"
				}
				s.WriteString(fmt.Sprintf("%s %s %s\n", cursor, checked, t.Text))
			}
			s.WriteString("\n")
		}

		// Show added todos (including selected backlog items preview?)
		// For now just show newly added ones
		if len(m.Entry.Todos) > 0 {
			s.WriteString("Added:\n")
			for _, t := range m.Entry.Todos {
				s.WriteString(fmt.Sprintf("- %s\n", t.Text))
			}
			s.WriteString("\n")
		}

		s.WriteString(m.TodoInput.View())
		s.WriteString("\n\n(Enter to add, Empty Enter to finish)")

	case StepQuestions:
		currentTemplate := m.Templates[m.TemplateCursor]
		if m.QuestionIndex < len(currentTemplate.Questions) {
			q := currentTemplate.Questions[m.QuestionIndex].Title
			s.WriteString(titleStyle.Render(q))
			s.WriteString("\n\n")
			s.WriteString(m.QuestionInput.View())
			s.WriteString("\n\n(Ctrl+S or Ctrl+N to next)")
		}

	case StepDone:
		s.WriteString(titleStyle.Render("All done!"))
		s.WriteString("\n\nSaving journal entry...")
	}

	return s.String()
}
