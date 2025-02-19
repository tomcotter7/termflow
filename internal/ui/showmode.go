package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var showModeFocusedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Width(15).Align(lipgloss.Left)

func (m model) showModeView() string {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 20
		height = 10
	}
	if m.cursor.row < len(m.formattedTasks[m.cursor.col]) {
		item := m.formattedTasks[m.cursor.col][m.cursor.row]
		if task, exists := m.structuredTasks[item]; exists {

			var s strings.Builder

			s.WriteString(fmt.Sprintf("%-12s %s\n", showModeFocusedStyle.Render("Title:"), task.Desc))
			s.WriteString(fmt.Sprintf("%-12s %s\n", showModeFocusedStyle.Render("Description:"), task.FullDesc))
			s.WriteString(fmt.Sprintf("%-12s %s\n", showModeFocusedStyle.Render("Created:"), task.Created))
			s.WriteString(fmt.Sprintf("%-12s %s\n", showModeFocusedStyle.Render("Due:"), task.Due))
			s.WriteString(fmt.Sprintf("%-12s %s\n", showModeFocusedStyle.Render("Blocked:"), fmt.Sprintf("%v", task.Blocked)))

			contentWidth := max(len(task.Desc), len(task.FullDesc), len(task.Created), len(task.Due))

			content := s.String()
			contentHeight := strings.Count(content, "\n") + 1
			topPadding := (height - contentHeight) / 8
			leftPadding := (width - contentWidth) / 2
			style := lipgloss.NewStyle().
				Width(width).
				Align(lipgloss.Left).
				PaddingTop(topPadding).
				PaddingLeft(leftPadding)

			return style.Render(content)

		}
	}
	content := "Nothing to see here!"
	contentHeight := strings.Count(content, "\n") + 1
	topPadding := (height - contentHeight) / 8
	style := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		PaddingTop(topPadding)
	return style.Render(content)
}

func (m model) handleShowModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.mode = NormalMode
		}
	}
	return m, nil
}
