package ui

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomcotter7/termflow/internal/storage"
)

func randomId() string {
	b := make([]byte, 16)
	rand.Read(b)

	return hex.EncodeToString(b)
}

func (m model) inputModeView() string {
	var b strings.Builder
	for i := range m.createTaskForm.inputs.ti {
		b.WriteString(m.createTaskForm.inputs.ti[i].View())
		b.WriteRune('\n')
	}
	b.WriteRune('\n')
	for i := range m.createTaskForm.inputs.ta {
		b.WriteString(m.createTaskForm.inputs.ta[i].View())
		if i < len(m.createTaskForm.inputs.ta)-1 {
			b.WriteRune('\n')
		}
	}
	button := &blurredButton
	if m.createTaskForm.inputs.onSubmitButton() {
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

func (m model) handleInputModelUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
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

			if len(m.createTaskForm.inputTaskId) == 0 {
				m.createTaskForm.inputTaskId = randomId()
			}
			newTask := storage.Task{
				ID:       m.createTaskForm.inputTaskId,
				Status:   columnNames[m.cursor.col],
				Desc:     m.createTaskForm.inputs.ti[0].Value(),
				FullDesc: m.createTaskForm.inputs.ta[0].Value(),
				Priority: priority,
				Created:  time.Now().Format("2006-01-02"),
				Due:      dd,
			}
			m.tasks[m.createTaskForm.inputTaskId] = newTask
			m.handler.SaveTasks(m.activeProject+".json", m.tasks)
			m.formattedTasks = formatTasks(m.tasks)
			m.mode = NormalMode

			m.createTaskForm.inputs.reset()
			m.createTaskForm.inputTaskId = ""
			return m, nil
		}

		switch msg.String() {
		case "esc":
			m.mode = NormalMode
			m.createTaskForm.inputs.reset()
			return m, nil
		case "tab", "shift+tab":
			s := msg.String()
			if s == "shift+tab" {
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
