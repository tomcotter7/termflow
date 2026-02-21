package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) notesViewModeView() string {
	return listContainerStyle.Render(m.notesList.View())
}

func (m model) handleNotesViewModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		m.notesList.SetSize(m.termWidth, m.termHeight)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q", "esc":
			m.mode = CommandMode
			return m, nil
		case "enter":
			selected := m.notesList.SelectedItem().(item)
			id := selected.id
			desc := selected.desc
			delete(m.notes, id)
			m.notesList = createNotesListModel(m.notes)
			m.handler.SaveNotes(m.activeProject+"_notes.json", m.notes)

			m.mode = EditMode
			m.createTaskForm.inputs.focusInput(0)
			m.createTaskForm.inputs.ta[0].SetValue(desc)
			m.createTaskForm.inputTaskId = ""
			m.createTaskForm.inputs.ta[0].SetWidth(m.termWidth / 2)
			m.createTaskForm.inputs.ta[0].SetHeight(m.termHeight / 4)
			m.createTaskForm.inputs.ta[1].SetWidth(m.termWidth / 2)
			m.createTaskForm.inputs.ta[1].SetHeight(m.termHeight / 4)

		}

	}

	var cmd tea.Cmd
	m.notesList, cmd = m.notesList.Update(msg)
	return m, cmd
}
