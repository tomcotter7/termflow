package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) newProjectModeView() string {
	content := m.createProjectForm.inputs.buildFormView()
	return m.centeredView(content)
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
