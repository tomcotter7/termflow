package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var showModetitleStyle = focusedStyle.Copy().Bold(true).PaddingBottom(1)

func trimToLength(s string, maxLength int) string {
	if len(s) > maxLength {
		s = s[:maxLength-3] + "..."
	}

	return s
}

func createRenderedAttribute(attributeTitle string, attributeDescription string, titleStyle lipgloss.Style, boxWidth int) string {
	titleStyle = showModetitleStyle.Copy().Width(boxWidth - 2).Align(lipgloss.Center)
	attributeTitle = titleStyle.Render(attributeTitle)
	attributeDescription = trimToLength(attributeDescription, boxWidth-2)

	contentStr := fmt.Sprintf("%s\n%s", attributeTitle, attributeDescription)

	boxStyle := borderStyle.Copy().Width(boxWidth).Align(lipgloss.Center)
	return boxStyle.Render(contentStr)
}

func (m model) showModeView() string {
	if m.cursor.row < len(m.formattedTasks[m.cursor.col]) {
		item := m.formattedTasks[m.cursor.col][m.cursor.row]
		if task, exists := m.tasks[item.ID]; exists {

			boxWidth := int(float64(m.termWidth) * 0.8)
			style := showModetitleStyle.Copy().Width(boxWidth - 2).Align(lipgloss.Center)

			title := createRenderedAttribute("Title", task.Desc, style, boxWidth)
			description := createRenderedAttribute("Description", task.FullDesc, style, boxWidth)
			dates := createRenderedAttribute("Created | Due", task.Created+" | "+task.Due, style, boxWidth)
			other := createRenderedAttribute("Other Attributes", fmt.Sprintf("%s: %t\n%s: %d\n%s: %t", boldStyle.Render("Blocked"), task.Blocked, boldStyle.Render("Priority"), task.Priority, boldStyle.Render("Ignore from .plan"), task.IgnoreFromPlan), style, boxWidth)

			var s strings.Builder

			s.WriteString(title + "\n")
			s.WriteString(description + "\n")
			s.WriteString(dates + "\n")
			s.WriteString(other)

			content := s.String()
			contentHeight := strings.Count(content, "\n") + 1
			topPadding := (m.termHeight - contentHeight) / 8
			rootStyle := lipgloss.NewStyle().
				Width(m.termWidth).
				Align(lipgloss.Center).
				PaddingTop(topPadding)

			return rootStyle.Render(content)

		}
	}
	content := "Nothing to see here!"
	contentHeight := strings.Count(content, "\n") + 1
	topPadding := (m.termHeight - contentHeight) / 8
	style := lipgloss.NewStyle().
		Width(m.termWidth).
		Align(lipgloss.Center).
		PaddingTop(topPadding)
	return style.Render(content)
}

func (m model) handleShowModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.mode = NormalMode
		}
	}
	return m, nil
}
