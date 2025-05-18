package ui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const labelWidth = 25

var showModeFocusedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Align(lipgloss.Left).Width(labelWidth)

func trimToLength(s string, maxLength int) string {
	if len(s) > maxLength {
		s = s[:maxLength-3] + "..."
	}

	return s
}

func (m model) showModeView() string {
	if m.cursor.row < len(m.formattedTasks[m.cursor.col]) {
		item := m.formattedTasks[m.cursor.col][m.cursor.row]
		if task, exists := m.tasks[item.ID]; exists {

			var s strings.Builder

			s.WriteString(fmt.Sprintf("%s %s\n", showModeFocusedStyle.Render("Title:"), trimToLength(task.Desc, m.termWidth-labelWidth)))
			s.WriteString(fmt.Sprintf("%s %s\n", showModeFocusedStyle.Render("Description:"), trimToLength(task.FullDesc, m.termWidth-labelWidth)))
			s.WriteString(fmt.Sprintf("%s %s\n", showModeFocusedStyle.Render("Created:"), trimToLength(task.Created, m.termWidth-labelWidth)))
			s.WriteString(fmt.Sprintf("%s %s\n", showModeFocusedStyle.Render("Due:"), trimToLength(task.Due, m.termWidth-labelWidth)))
			s.WriteString(fmt.Sprintf("%s %s\n", showModeFocusedStyle.Render("Blocked:"), trimToLength(strconv.FormatBool(task.Blocked), m.termWidth-labelWidth)))
			s.WriteString(fmt.Sprintf("%s %s\n", showModeFocusedStyle.Render("Priority"), trimToLength(strconv.Itoa(task.Priority), m.termWidth-labelWidth)))
			s.WriteString(fmt.Sprintf("%s %s\n", showModeFocusedStyle.Render("Ignore from .plan:"), trimToLength(strconv.FormatBool(task.IgnoreFromPlan), m.termWidth-labelWidth)))

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
