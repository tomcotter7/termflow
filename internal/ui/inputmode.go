package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
	"github.com/tomcotter7/termflow/internal/storage"

	"golang.org/x/term"
)

func (m model) inputModeView() string {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 20
		height = 10
	}

	var b strings.Builder
	for i := range m.textInputs {
		b.WriteString(m.textInputs[i].View())
		if i < len(m.textInputs)-1 {
			b.WriteRune('\n')
		}
	}
	button := &blurredButton
	if m.focusedIndex == len(m.textInputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	content := b.String()

	contentHeight := strings.Count(content, "\n") + 1
	topPadding := (height - contentHeight) / 8

	style := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		PaddingTop(topPadding)

	return style.Render(content)
}

func (m model) handleInputModelUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mode = "normal"
			m.resetInputs()
			return m, nil
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusedIndex == len(m.textInputs) {
				dd, err := time.Parse("2006-01-02", m.textInputs[2].Value())
				if err != nil {
					m.textInputs[2].Reset()
					m.textInputs[2].Focus()
					m.focusedIndex = 2
				} else {
					if m.inputTaskId == "" {
						m.inputTaskId = randomId()
					}
					newTask := storage.Task{
						Status:   columnNames[m.cursor.col],
						Desc:     m.textInputs[0].Value(),
						FullDesc: m.textInputs[1].Value(),
						Created:  time.Now().Format("2006-01-02"),
						Due:      dd.Format("2006-01-02"),
					}
					m.structuredTasks[m.inputTaskId] = newTask
					m.handler.SaveTasks("default.json", m.structuredTasks)
					m.formattedTasks = formatTasks(m.structuredTasks)
					m.mode = "normal"

					m.resetInputs()
					return m, nil
				}
			}

			if s == "up" || s == "shift+tab" {
				m.focusedIndex = max(0, m.focusedIndex-1)
			} else {
				m.focusedIndex = min(m.focusedIndex+1, len(m.textInputs))
			}

			cmds := make([]tea.Cmd, len(m.textInputs))
			for i := 0; i <= len(m.textInputs)-1; i++ {
				if i == m.focusedIndex {
					cmds[i] = m.textInputs[i].Focus()
					m.textInputs[i].PromptStyle = focusedStyle
					m.textInputs[i].TextStyle = focusedStyle
					continue
				}

				m.textInputs[i].Blur()
				m.textInputs[i].PromptStyle = noStyle
				m.textInputs[i].TextStyle = noStyle
			}
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}
