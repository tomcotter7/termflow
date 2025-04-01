package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) switchProjectModeView() string {
	return docStyle.Render(m.projects.View())
}

func (m model) handleSwitchProjectModeUpdates(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			m.mode = NormalMode
			return m, nil
		case "enter":
			m.activeProject = m.projects.SelectedItem().FilterValue()
			m.handler.SaveCurrent(m.activeProject)
			sts, err := m.handler.LoadTasks(m.activeProject + ".json")
			if err != nil {
				m.err = err
				m.mode = NormalMode
				return m, nil
			}
			m.tasks = sts
			m.formattedTasks = formatTasks(m.tasks)
			m.mode = NormalMode
			m.cursor.row = 0
			m.cursor.col = 0
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.projects.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.projects, cmd = m.projects.Update(msg)
	return m, cmd
}
