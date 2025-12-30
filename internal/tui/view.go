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

		// Display Stats
		s.WriteString(subtle.Render(fmt.Sprintf(" 📝 Total Entries: %d", m.Stats.TotalEntries)))
		if !m.Stats.LastMissed.IsZero() {
			s.WriteString(subtle.Render(fmt.Sprintf(" | 🗓️  Last Missed: %s", m.Stats.LastMissed.Format("Monday, 02 Jan"))))
		}
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

		// If Todos menu active, show options to avoid navigation deadlocks
		if m.TodosMenuActive {
			opts := []string{"Edit Todos", "Manage Backlog", "Start Fresh", "Continue"}
			s.WriteString("Choose how to manage today's todos:\n")
			for i, o := range opts {
				cursor := " "
				style := itemStyle
				if m.TodosMenuCursor == i {
					cursor = ">"
					style = selectedItemStyle
				}
				s.WriteString(style.Render(fmt.Sprintf("%s %s\n", cursor, o)))
			}
			s.WriteString("\n(Use Up/Down to select, Enter to confirm)")
			return s.String()
		}

		// Backlog + Added todos rendered as a single linear selectable list.
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

		if len(m.Entry.Todos) > 0 {
			s.WriteString("Added (Enter to edit):\n")
			// offset index for added todos is len(Backlog)
			off := len(m.Entry.Backlog)
			for j, t := range m.Entry.Todos {
				cursor := " "
				if !m.TodoInput.Focused() && m.BacklogCursor == off+j {
					cursor = ">"
				}
				s.WriteString(fmt.Sprintf("%s - %s\n", cursor, t.Text))
			}
			s.WriteString("\n")
		}

		s.WriteString(m.TodoInput.View())
		s.WriteString("\n\n(Enter to add, Empty Enter to finish; Tab/Shift+Tab to switch focus; Up/Down to navigate; Enter on added todo to edit)")

	case StepQuestions:
		currentTemplate := m.Templates[m.TemplateCursor]
		if m.QuestionIndex < len(currentTemplate.Questions) {
			q := currentTemplate.Questions[m.QuestionIndex].Title
			s.WriteString(titleStyle.Render(q))
			s.WriteString("\n\n")
			s.WriteString(m.QuestionInput.View())
			s.WriteString("\n\n(Enter to save+next; Shift+Right to next; Shift+Left to previous; Ctrl+S or Ctrl+N still advances)")
		}

	case StepDone:
		s.WriteString(titleStyle.Render("All done!"))
		s.WriteString("\n\nSaving journal entry...")
	}

	return s.String()
}
