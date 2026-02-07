package ui

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomcotter7/termflow/internal/storage"
)

func randomId() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (m model) editModeView() string {
	content := m.createTaskForm.inputs.buildFormView()

	contentHeight := strings.Count(content, "\n") + 1
	topPadding := (m.termHeight - contentHeight) / 8

	style := lipgloss.NewStyle().
		Width(m.termWidth).
		Align(lipgloss.Center).
		PaddingTop(topPadding)

	return style.Render(content)
}

func (m model) handleEditModelUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if msg.String() == "enter" && m.createTaskForm.inputs.onSubmitButton() {
			var dd string
			switch strings.ToLower(m.createTaskForm.inputs.ti[1].Value()) {
			case "today":
				dd = time.Now().Format("2006-01-02")
				m.createTaskForm.inputs.ti[1].SetValue(dd)
			case "tomorrow":
				dd = time.Now().Add(24 * time.Hour).Format("2006-01-02")
				m.createTaskForm.inputs.ti[1].SetValue(dd)
			case "end-of-week":
				today := time.Now().Weekday()
				friday := time.Friday
				diff := int(friday - today)
				if diff <= 0 {
					diff += 7
				}
				dd = time.Now().Add(time.Hour * 24 * time.Duration(diff)).Format("2006-01-02")
				m.createTaskForm.inputs.ti[1].SetValue(dd)
			case "none":
				dd = "none"
			default:
				date, err := time.Parse("2006-01-02", m.createTaskForm.inputs.ti[1].Value())
				if err != nil {
					m.createTaskForm.inputs.focusInput(1)
					return m, nil
				}
				dd = date.Format("2006-01-02")
			}

			priority, err := strconv.Atoi(m.createTaskForm.inputs.ti[2].Value())
			if err != nil {
				m.createTaskForm.inputs.focusInput(2)
				return m, nil
			}

			created := time.Now().Format("2006-01-02")
			blocked := false
			id, err := randomId()
			if err != nil {
				m.err = err
				m.mode = ErrorMode
				return m, nil
			}

			if len(m.createTaskForm.inputTaskId) > 0 {
				created = m.tasks[m.createTaskForm.inputTaskId].Created
				blocked = m.tasks[m.createTaskForm.inputTaskId].Blocked
				id = m.createTaskForm.inputTaskId
			}

			newTask := storage.Task{
				ID:       id,
				Status:   columnNames[m.cursor.col],
				Desc:     m.createTaskForm.inputs.ti[0].Value(),
				FullDesc: m.createTaskForm.inputs.ta[0].Value(),
				Priority: priority,
				Due:      dd,
				Created:  created,
				Blocked:  blocked,
				Result:   m.createTaskForm.inputs.ta[1].Value(),
			}
			m.tasks[id] = newTask
			m.handler.SaveTasks(m.activeProject+".json", m.tasks)
			m.formattedTasks = formatTasks(m.tasks)
			m.mode = NormalMode

			m.createTaskForm.inputs.reset()
			m.createTaskForm.inputTaskId = ""
			return m, nil
		}

		switch k := msg.String(); k {
		case "esc":
			m.mode = NormalMode
			m.createTaskForm.inputs.reset()
			m.createTaskForm.inputTaskId = ""
			return m, nil
		case "tab", "shift+tab":
			if k == "shift+tab" {
				m.createTaskForm.inputs.decreaseFocusedIndex()
			} else {
				m.createTaskForm.inputs.increaseFocusedIndex()
			}
			return m, nil
		}
	}

	cmd := m.createTaskForm.inputs.updateInputs(msg)
	return m, cmd
}
