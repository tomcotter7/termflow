package ui

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"

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
	return m.centeredView(content)
}

func get_date_diff(dow time.Weekday) string {
	today := time.Now().Weekday()
	diff := int(dow - today)
	if diff <= 0 {
		diff += 7
	}
	dd := time.Now().Add(time.Hour * 24 * time.Duration(diff)).Format("2006-01-02")
	return dd
}

func get_true_datetime(datestring string) (string, error) {
	var dd string
	switch datestring {
	case "today":
		dd = time.Now().Format("2006-01-02")
	case "tomorrow":
		dd = time.Now().Add(24 * time.Hour).Format("2006-01-02")
	case "monday", "mon":
		dd = get_date_diff(time.Monday)
	case "tuesday", "tues":
		dd = get_date_diff(time.Tuesday)
	case "wednesday", "wed":
		dd = get_date_diff(time.Wednesday)
	case "thursday", "thurs":
		dd = get_date_diff(time.Thursday)
	case "end-of-week", "friday", "fri":
		dd = get_date_diff(time.Friday)
	case "none":
		dd = "none"
	default:
		date, err := time.Parse("2006-01-02", datestring)
		if err != nil {
			return "", err
		}
		dd = date.Format("2006-01-02")
	}
	return dd, nil
}

func (m model) handleEditModelUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		m.createTaskForm.inputs.ta[0].SetWidth(m.termWidth / 2)
		m.createTaskForm.inputs.ta[0].SetHeight(m.termHeight / 4)
		m.createTaskForm.inputs.ta[1].SetWidth(m.termWidth / 2)
		m.createTaskForm.inputs.ta[1].SetHeight(m.termHeight / 4)
	case tea.KeyMsg:

		if msg.String() == "enter" && m.createTaskForm.inputs.onSubmitButton() {
			dd, err := get_true_datetime(m.createTaskForm.inputs.ti[1].Value())
			if err != nil {
				m.createTaskForm.inputs.focusInput(1)
				return m, nil
			}
			m.createTaskForm.inputs.ti[1].SetValue(dd)

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
			m.saveAndUpdateTasks()
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
		case "tab":
			m.createTaskForm.inputs.increaseFocusedIndex()
			return m, nil
		case "shift+tab":
			m.createTaskForm.inputs.decreaseFocusedIndex()
			return m, nil
		}
	}

	cmd := m.createTaskForm.inputs.updateInputs(msg)
	return m, cmd
}
