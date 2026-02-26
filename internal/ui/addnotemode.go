package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomcotter7/termflow/internal/storage"
)

func (m model) addNoteModeView() string {
	content := m.addNoteForm.inputs.buildFormView()
	return m.centeredView(content)
}

func (m model) handleAddNoteModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

		m.addNoteForm.inputs.ta[0].SetWidth(m.termWidth / 2)
		m.addNoteForm.inputs.ta[0].SetHeight(m.termHeight / 2)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.addNoteForm.inputs.onSubmitButton() {

				id := m.addNoteForm.prevID
				var created string

				if id != "" {
					created = m.notes[id].Created
				} else {
					var err error
					id, err = randomId()
					if err != nil {
						m.err = err
						m.mode = ErrorMode
						return m, nil
					}
					created = time.Now().Format("2006-01-02")
				}

				newNote := storage.Note{
					ID:      id,
					Created: created,
					Content: m.addNoteForm.inputs.ta[0].Value(),
				}
				m.notes[id] = newNote
				m.notesList = createNotesListModel(m.notes)
				m.handler.SaveNotes(m.activeProject+"_notes.json", m.notes)
				m.mode = NormalMode
				m.addNoteForm.inputs.reset()
				m.addNoteForm.prevID = ""
				return m, nil
			}
		case "esc":
			m.mode = CommandMode
			return m, nil
		case "tab":
			m.addNoteForm.inputs.increaseFocusedIndex()
			return m, nil
		case "shift+tab":
			m.addNoteForm.inputs.decreaseFocusedIndex()
			return m, nil
		}
	}

	cmd := m.addNoteForm.inputs.updateInputs(msg)
	return m, cmd
}
