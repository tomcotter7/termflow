package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) saveAndUpdateTasks(filename string) {
	m.formattedTasks = formatTasks(m.structuredTasks)
	m.handler.SaveTasks(filename, m.structuredTasks)
}

func (m model) handleNormalModelUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "l", "right":
			m.cursor.IncCol(m.formattedTasks)
		case "h", "left":
			m.cursor.DecCol(m.formattedTasks)
		case "j", "down":
			m.cursor.IncRow(len(m.formattedTasks[m.cursor.col]) - 1)
		case "k", "up":
			m.cursor.DecRow()
		case "p":

			if len(m.formattedTasks[m.cursor.col]) <= m.cursor.row {
				return m, nil
			}

			item := m.formattedTasks[m.cursor.col][m.cursor.row]

			if task, exists := m.structuredTasks[item]; exists && m.cursor.col < 2 {
				task.Status = columnNames[m.cursor.col+1]
				m.structuredTasks[item] = task
				m.saveAndUpdateTasks("default.json")
			}

			m.cursor.IncCol(m.formattedTasks)

		case "r":
			if len(m.formattedTasks[m.cursor.col]) <= m.cursor.row {
				return m, nil
			}

			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.structuredTasks[item]; exists && m.cursor.col > 0 {
				task.Status = columnNames[m.cursor.col-1]
				m.structuredTasks[item] = task
				m.saveAndUpdateTasks("default.json")
			}
			m.cursor.DecCol(m.formattedTasks)
		case "b":
			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.structuredTasks[item]; exists {
				task.Blocked = !task.Blocked
				m.structuredTasks[item] = task
				m.saveAndUpdateTasks("default.json")
			}

		case "a":
			m.mode = "input"
			m.textInputs[0].Focus()
			m.textInputs[0].PromptStyle = focusedStyle
			m.textInputs[0].TextStyle = focusedStyle
		case "e":

			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.structuredTasks[item]; exists {
				m.mode = "input"
				m.textInputs[0].Focus()
				m.textInputs[0].PromptStyle = focusedStyle
				m.textInputs[0].TextStyle = focusedStyle
				m.inputTaskId = item
				m.textInputs[0].SetValue(task.Desc)
				m.textInputs[1].SetValue(task.FullDesc)
				m.textInputs[2].SetValue(task.Due)
			}
		case "s", "enter":
			m.mode = "show"
		case "t":
			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.structuredTasks[item]; exists {
				task.Due = time.Now().Format("2006-01-02")

				m.structuredTasks[item] = task
				m.saveAndUpdateTasks("default.json")
			}
		case "d":
			if len(m.formattedTasks[m.cursor.col]) == 0 {
				return m, nil
			}
			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			delete(m.structuredTasks, item)
			m.saveAndUpdateTasks("default.json")
		case "?":
			m.help = !m.help
		case ":":
			m.mode = "command"
			m.commands.SetSize(m.width-2, m.height-2)
		}
	}

	return m, nil
}

func (m model) normalModeView() string {
	ml := maxTaskLength(m.structuredTasks)

	tTitle, iTitle, dTitle := "todo", "inprogress", "done"

	minPadding := 3
	space := max(ml, len(tTitle), len(iTitle), len(dTitle)) + minPadding + 2

	switch m.cursor.col {
	case 0:
		tTitle = "* " + tTitle
	case 1:
		iTitle = "* " + iTitle
	case 2:
		dTitle = "* " + dTitle
	}

	tTitle, iTitle, dTitle = addPadding(tTitle, space), addPadding(iTitle, space), addPadding(dTitle, space)

	tTitle = redText.Render(tTitle)
	iTitle = yellowText.Render(iTitle)
	dTitle = greenText.Render(dTitle)

	var s strings.Builder

	s.WriteString("╔" + strings.Repeat("═", space) + "╦" + strings.Repeat("═", space) + "╦" + strings.Repeat("═", space) + "╗\n")

	s.WriteString(fmt.Sprintf("║%s║%s║%s║", tTitle, iTitle, dTitle) + "\n")

	s.WriteString("╠" + strings.Repeat("═", space) + "╬" + strings.Repeat("═", space) + "╬" + strings.Repeat("═", space) + "╣\n")

	sortTasks(&m.formattedTasks)
	tt := transpose(m.formattedTasks)
	for i := range tt {
		tasks := make([]string, 3)
		for j := range 3 {
			task := tt[i][j]
			taskData := m.structuredTasks[task]

			tasks[j] = taskData.Desc
			if m.cursor.row == i && m.cursor.col == j {
				if len(tasks[j]) > 0 {
					tasks[j] = "> " + tasks[j]
				} else {
					tasks[j] = "--+--"
				}
			}
			tasks[j] = addPadding(tasks[j], space)

			if j < 2 {

				if taskData.Due == time.Now().Format("2006-01-02") && !taskData.Blocked {
					tasks[j] = blueText.Render(tasks[j])
				}

				if taskData.Blocked {
					tasks[j] = redText.Render(tasks[j])
				}

			} else {
				tasks[j] = blurredStyle.Render(tasks[j])
			}
		}

		tTask, iTask, dTask := tasks[0], tasks[1], tasks[2]
		s.WriteString(fmt.Sprintf("║%s║%s║%s║\n", tTask, iTask, dTask))
	}

	s.WriteString("╚" + strings.Repeat("═", space) + "╩" + strings.Repeat("═", space) + "╩" + strings.Repeat("═", space) + "╝\n")

	if m.help {
		s.WriteString(helpStyle.Render("\nCommands:\n"))
		s.WriteString(helpStyle.Render("\na: (a)dd • p: (p)romote • r: (r)egress • d: (d)elete • e: (e)dit • s: (s)how • \nt: (t)oday • b: (b)locked • q: (q)uit • ':': command-mode • ?: hide\n"))
	} else {
		s.WriteString(helpStyle.Render("\n?: help\n"))
	}

	content := s.String()
	contentHeight := strings.Count(content, "\n") + 1
	topPadding := (m.height - 4 - contentHeight) / 8
	style := lipgloss.NewStyle().
		Width(m.width - 4).
		Align(lipgloss.Center).
		PaddingTop(topPadding)

	return style.Render(content)
}
