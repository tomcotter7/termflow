package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func (m *model) showModeView() string {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 20
		height = 10
	}
	if m.cursor.row < len(m.formattedTasks[m.cursor.col]) {
		item := m.formattedTasks[m.cursor.col][m.cursor.row]
		if task, exists := m.structuredTasks[item]; exists {

			var s strings.Builder

			s.WriteString("Title: " + task.Desc + "\n")
			s.WriteString("Description: " + task.FullDesc + "\n")
			s.WriteString("Created: " + task.Created + "\n")
			s.WriteString("Due: " + task.Due + "\n")
			s.WriteString("Blocked: " + fmt.Sprintf("%v", task.Blocked) + "\n")

			content := s.String()
			contentHeight := strings.Count(content, "\n") + 1
			topPadding := (height - contentHeight) / 8
			style := lipgloss.NewStyle().
				Width(width).
				Align(lipgloss.Center).
				PaddingTop(topPadding)

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

func (m *model) handleShowModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.mode = "normal"
			m.resetInputs()
		}
	}
	return m, nil
}
