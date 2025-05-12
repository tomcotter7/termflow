package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) errorModeView() string {
	error := "So sorry for inconvenience! (⌒_⌒;) termflow-san is experiencing difficulties. Please submit a issue on `github.com/tomcotter7/termflow`. Thank you for your most honorable patience!\n"

	var s strings.Builder

	s.WriteString(redText.Render(error))
	s.WriteString(fmt.Sprintf("%s \n", redText.Render(m.err.Error())))
	s.WriteString("Press (q) to re-enter normal mode\n")

	content := s.String()

	contentHeight := strings.Count(content, "\n")

	topPadding := (m.termHeight - contentHeight) / 8

	style := lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Padding(topPadding).Bold(true)

	return style.Render(content)
}

func (m model) handleErrorModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.mode = NormalMode
			m.err = nil
		}
	}

	return m, nil
}
