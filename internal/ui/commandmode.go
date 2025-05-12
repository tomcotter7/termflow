package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tomcotter7/termflow/internal/storage"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m model) writeToPlanFile(tasks map[string]storage.Task) error {
	var s strings.Builder

	for k, v := range tasks {
		if v.Status == "done" && !v.IgnoreFromPlan {
			s.WriteString("Task ID: " + k + "\n")
			s.WriteString("Title: " + v.Desc + "\n")
			if len(v.FullDesc) > 1 {
				s.WriteString("Full Description: " + v.FullDesc + "\n")
			}
			s.WriteString("Created on: " + v.Created + "\n")
			s.WriteString("---\n")
		}
	}

	content := s.String()
	today := time.Now().Format("2006-01-02")
	filename := today + m.activeProject + ".plan"
	err := m.handler.SavePlanFile(filename, content)
	return err
}

func (m *model) executeCommand(command string) {
	switch strings.ToLower(command) {
	case "clear":
		newTasks := make(map[string]storage.Task)
		for k, v := range m.tasks {
			if v.Status != "done" {
				newTasks[k] = v
			}
		}
		m.tasks = newTasks
		m.saveAndUpdateTasks()
		m.mode = NormalMode
	case "print":
		err := m.writeToPlanFile(m.tasks)
		if err != nil {
			m.mode = ErrorMode
			m.err = err
		}
		m.mode = NormalMode
	case "create project":
		m.createProjectForm.textInputs.focusTextInput(0)
		m.mode = NewProjectMode
	case "switch to project":
		m.projects.SetSize(m.termWidth-2, m.termHeight-2)
		m.mode = SwitchProjectMode
	case "show daily work %":
		m.mode = ShowWorkPercentageMode
	}
}

func (m model) commandModeView() string {
	return docStyle.Render(m.commands.View())
}

func (m model) handleCommandModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			m.mode = NormalMode
			return m, nil
		case "enter":
			command := m.commands.SelectedItem().FilterValue()
			m.executeCommand(command)
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.commands.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.commands, cmd = m.commands.Update(msg)
	return m, cmd
}
