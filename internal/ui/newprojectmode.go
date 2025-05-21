package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) newProjectModeView() string {
	var b strings.Builder
	for i := range m.createProjectForm.textInputs.ti {
		b.WriteString(m.createProjectForm.textInputs.ti[i].View())
		if i < len(m.createProjectForm.textInputs.ti)-1 {
			b.WriteRune('\n')
		}
	}
	button := &blurredButton
	if m.createProjectForm.textInputs.onSubmitButton() {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	content := b.String()

	contentHeight := strings.Count(content, "\n") + 1
	topPadding := (m.termHeight - contentHeight) / 8

	style := lipgloss.NewStyle().
		Width(m.termWidth).
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
			m.createProjectForm.textInputs.resetTextInputs()
			return m, nil
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			if s == "enter" && m.createProjectForm.textInputs.onSubmitButton() {
				newProject := m.createProjectForm.textInputs.ti[0].Value()
				if len(newProject) == 0 {
					m.createProjectForm.textInputs.resetTextInputs()
					m.createProjectForm.textInputs.focusTextInput(0)
				}

				m.activeProject = newProject
				m.handler.SaveCurrent(newProject)
				m.mode = NormalMode
				m.createProjectForm.textInputs.resetTextInputs()
				sts, err := m.handler.LoadTasks(m.activeProject + ".json")
				if err != nil {
					m.err = err
					m.mode = ErrorMode
					return m, nil
				}
				m.tasks = sts
				m.formattedTasks = formatTasks(m.tasks)
				m.projects, err = newProjectListModel(m.handler)
				if err != nil {
					m.err = err
					m.mode = ErrorMode
					return m, nil
				}
				return m, nil
			}
			if s == "up" || s == "shift+tab" {
				m.createProjectForm.textInputs.decreaseFocusedIndex()
			} else {
				m.createProjectForm.textInputs.increaseFocusedIndex()
			}
		}
	}

	cmd := m.createProjectForm.textInputs.updateTextInputs(msg)
	return m, cmd
}
