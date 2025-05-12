package ui

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomcotter7/termflow/internal/storage"
)

func randomId() string {
	b := make([]byte, 5)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (m model) inputModeView() string {
	var b strings.Builder
	for i := range m.createTaskForm.textInputs.ti {
		b.WriteString(m.createTaskForm.textInputs.ti[i].View())
		if i < len(m.createTaskForm.textInputs.ti)-1 {
			b.WriteRune('\n')
		}
	}
	button := &blurredButton
	if m.createTaskForm.textInputs.onSubmitButton() {
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
		switch msg.String() {
		case "esc":
			m.mode = NormalMode
			m.createTaskForm.textInputs.resetTextInputs()
			return m, nil
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.createTaskForm.textInputs.onSubmitButton() {
				switch strings.ToLower(m.createTaskForm.textInputs.ti[2].Value()) {
				case "today":
					m.createTaskForm.textInputs.ti[2].SetValue(time.Now().Format("2006-01-02"))
				case "tomorrow":
					m.createTaskForm.textInputs.ti[2].SetValue(time.Now().Add(24 * time.Hour).Format("2006-01-02"))
				case "end-of-week":
					today := time.Now().Weekday()
					friday := time.Friday
					diff := int(friday - today)
					if diff <= 0 {
						diff += 7
					}
					m.createTaskForm.textInputs.ti[2].SetValue(time.Now().Add(time.Hour * 24 * time.Duration(diff)).Format("2006-01-02"))
				}

				dd, err := time.Parse("2006-01-02", m.createTaskForm.textInputs.ti[2].Value())
				if err != nil {
					m.createTaskForm.textInputs.focusTextInput(2)
					return m, nil
				} else {
					if len(m.createTaskForm.inputTaskId) == 0 {
						m.createTaskForm.inputTaskId = randomId()
					}
					newTask := storage.Task{
						Status:   columnNames[m.cursor.col],
						Desc:     m.createTaskForm.textInputs.ti[0].Value(),
						FullDesc: m.createTaskForm.textInputs.ti[1].Value(),
						Created:  time.Now().Format("2006-01-02"),
						Due:      dd.Format("2006-01-02"),
					}
					m.tasks[m.createTaskForm.inputTaskId] = newTask
					m.handler.SaveTasks(m.activeProject+".json", m.tasks)
					m.formattedTasks = formatTasks(m.tasks)
					m.mode = NormalMode

					m.createTaskForm.textInputs.resetTextInputs()
					m.createTaskForm.inputTaskId = ""
					return m, nil
				}
			}

			if s == "up" || s == "shift+tab" {
				m.createTaskForm.textInputs.decreaseFocusedIndex()
			} else {
				m.createTaskForm.textInputs.increaseFocusedIndex()
			}
		}
	}

	cmd := m.createTaskForm.textInputs.updateTextInputs(msg)
	return m, cmd
}
