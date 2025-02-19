package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) newProjectModeView() string {
	var b strings.Builder
	for i := range m.createProjectInput.textInputs.ti {
		b.WriteString(m.createProjectInput.textInputs.ti[i].View())
		if i < len(m.createProjectInput.textInputs.ti)-1 {
			b.WriteRune('\n')
		}
	}
	button := &blurredButton
	if m.createProjectInput.textInputs.onSubmitButton() {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	content := b.String()

	contentHeight := strings.Count(content, "\n") + 1
	topPadding := (m.height - contentHeight) / 8

	style := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		PaddingTop(topPadding)

	return style.Render(content)
}

func (m model) handleNewProjectModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mode = NormalMode
			m.createProjectInput.textInputs.resetTextInputs()
			return m, nil
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			if s == "enter" && m.createProjectInput.textInputs.onSubmitButton() {
				m.project = m.createProjectInput.textInputs.ti[0].Value()
				m.mode = NormalMode
				m.createProjectInput.textInputs.resetTextInputs()
				sts, err := m.handler.LoadTasks(m.project + ".json")
				m.projects = formatProjects(m.handler)
				if err != nil {
					m.error = err
				}
				m.structuredTasks = sts
				return m, nil
			}
			if s == "up" || s == "shift+tab" {
				m.createProjectInput.textInputs.decreaseFocusedIndex()
			} else {
				m.createProjectInput.textInputs.increaseFocusedIndex()
			}
			for i := 0; i <= len(m.createProjectInput.textInputs.ti)-1; i++ {
				if i == m.createTaskInput.textInputs.focusedIdx {
					m.createProjectInput.textInputs.focusTextInput(i)
					continue
				}

				m.createProjectInput.textInputs.deFocusTextInput(i)
			}
		}
	}

	cmd := m.createProjectInput.textInputs.updateTextInputs(msg)
	return m, cmd
}
