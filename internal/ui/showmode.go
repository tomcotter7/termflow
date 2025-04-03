package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var showModeFocusedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Width(15).Align(lipgloss.Left)

func (m model) showModeView() string {
	if m.cursor.row < len(m.formattedTasks[m.cursor.col]) {
		item := m.formattedTasks[m.cursor.col][m.cursor.row]
		if task, exists := m.tasks[item]; exists {

			var s strings.Builder

			s.WriteString(fmt.Sprintf("%-12s %s\n", showModeFocusedStyle.Render("Title:"), task.Desc))
			s.WriteString(fmt.Sprintf("%-12s %s\n", showModeFocusedStyle.Render("Description:"), task.FullDesc))
			s.WriteString(fmt.Sprintf("%-12s %s\n", showModeFocusedStyle.Render("Created:"), task.Created))
			s.WriteString(fmt.Sprintf("%-12s %s\n", showModeFocusedStyle.Render("Due:"), task.Due))
			s.WriteString(fmt.Sprintf("%-12s %s\n", showModeFocusedStyle.Render("Blocked:"), fmt.Sprintf("%v", task.Blocked)))

			contentWidth := max(len(task.Desc), len(task.FullDesc), len(task.Created), len(task.Due))

			content := s.String()
			contentHeight := strings.Count(content, "\n") + 1
			topPadding := (m.termHeight - contentHeight) / 8
			leftPadding := (m.termWidth - contentWidth) / 2
			style := lipgloss.NewStyle().
				Width(m.termWidth).
				Align(lipgloss.Left).
				PaddingTop(topPadding).
				PaddingLeft(leftPadding)

			return style.Render(content)

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
