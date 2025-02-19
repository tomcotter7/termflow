package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomcotter7/termflow/internal/storage"
)

func (m model) inputModeView() string {
	var b strings.Builder
	for i := range m.createTaskInput.textInputs.ti {
		b.WriteString(m.createTaskInput.textInputs.ti[i].View())
		if i < len(m.createTaskInput.textInputs.ti)-1 {
			b.WriteRune('\n')
		}
	}
	button := &blurredButton
	if m.createTaskInput.textInputs.onSubmitButton() {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	content := b.String()

	contentHeight := strings.Count(content, "\n") + 1
	topPadding := (m.height - contentHeight) / 8

	style := lipgloss.NewStyle().
		Width(m.width).
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
			m.createTaskInput.textInputs.resetTextInputs()
			return m, nil
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.createTaskInput.textInputs.onSubmitButton() {
				switch strings.ToLower(m.createTaskInput.textInputs.ti[2].Value()) {
				case "today":
					m.createTaskInput.textInputs.ti[2].SetValue(time.Now().Format("2006-01-02"))
				case "tomorrow":
					m.createTaskInput.textInputs.ti[2].SetValue(time.Now().Add(24 * time.Hour).Format("2006-01-02"))
				}

				dd, err := time.Parse("2006-01-02", m.createTaskInput.textInputs.ti[2].Value())
				if err != nil {
					m.createTaskInput.textInputs.focusTextInput(2)
					m.createTaskInput.textInputs.focusedIdx = 2
					return m, nil
				} else {
					if len(m.createTaskInput.inputTaskId) == 0 {
						m.createTaskInput.inputTaskId = randomId()
					}
					newTask := storage.Task{
						Status:   columnNames[m.cursor.col],
						Desc:     m.createTaskInput.textInputs.ti[0].Value(),
						FullDesc: m.createTaskInput.textInputs.ti[1].Value(),
						Created:  time.Now().Format("2006-01-02"),
						Due:      dd.Format("2006-01-02"),
					}
					m.structuredTasks[m.createTaskInput.inputTaskId] = newTask
					m.handler.SaveTasks(m.project+".json", m.structuredTasks)
					m.formattedTasks = formatTasks(m.structuredTasks)
					m.mode = NormalMode

					m.createTaskInput.textInputs.resetTextInputs()
					m.createTaskInput.inputTaskId = ""
					return m, nil
				}
			}

			if s == "up" || s == "shift+tab" {
				m.createTaskInput.textInputs.decreaseFocusedIndex()
			} else {
				m.createTaskInput.textInputs.increaseFocusedIndex()
			}

			for i := 0; i <= len(m.createTaskInput.textInputs.ti)-1; i++ {
				if i == m.createTaskInput.textInputs.focusedIdx {
					m.createTaskInput.textInputs.focusTextInput(i)
					continue
				}

				m.createTaskInput.textInputs.deFocusTextInput(i)
			}
		}
	}

	cmd := m.createTaskInput.textInputs.updateTextInputs(msg)
	return m, cmd
}
