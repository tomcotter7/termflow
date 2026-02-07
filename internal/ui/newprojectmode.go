package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) newProjectModeView() string {
	var b strings.Builder
	for i := range m.createProjectForm.inputs.ti {
		b.WriteString(m.createProjectForm.inputs.ti[i].View())
		if i < len(m.createProjectForm.inputs.ti)-1 {
			b.WriteRune('\n')
		}
	}
	button := &blurredButton
	if m.createProjectForm.inputs.onSubmitButton() {
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
			m.createProjectForm.inputs.reset()
			return m, nil
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			if s == "enter" && m.createProjectForm.inputs.onSubmitButton() {
				newProject := m.createProjectForm.inputs.ti[0].Value()
				if len(newProject) == 0 {
					m.createProjectForm.inputs.reset()
					m.createProjectForm.inputs.focusInput(0)
				}

				m.activeProject = newProject
				m.handler.SaveCurrent(newProject)
				m.mode = NormalMode
				m.createProjectForm.inputs.reset()
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
				m.createProjectForm.inputs.decreaseFocusedIndex()
			} else {
				m.createProjectForm.inputs.increaseFocusedIndex()
			}
		}
	}

	cmd := m.createProjectForm.inputs.updateInputs(msg)
	return m, cmd
}
