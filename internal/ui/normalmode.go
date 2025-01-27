package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func (m model) handleNormalModelUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
				m.formattedTasks = formatTasks(m.structuredTasks)
				m.handler.SaveTasks("default.json", m.structuredTasks)
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
				m.formattedTasks = formatTasks(m.structuredTasks)
				m.handler.SaveTasks("default.json", m.structuredTasks)
			}
			m.cursor.DecCol(m.formattedTasks)
		case "b":
			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.structuredTasks[item]; exists {
				task.Blocked = !task.Blocked
				m.structuredTasks[item] = task
				m.formattedTasks = formatTasks(m.structuredTasks)
				m.handler.SaveTasks("default.json", m.structuredTasks)
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
		case "s":
			m.mode = "show"
			// TODO: add 'show' mode
		case "d":
			if len(m.formattedTasks[m.cursor.col]) == 0 {
				return m, nil
			}
			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			delete(m.structuredTasks, item)
			m.formattedTasks = formatTasks(m.structuredTasks)
			m.handler.SaveTasks("default.json", m.structuredTasks)
		}
	}

	return m, nil
}

func (m *model) normalModeView() string {
	ml := maxTaskLength(m.structuredTasks)

	tTitle, iTitle, dTitle := "todo", "inprogress", "done"

	minPadding := 1
	space := max(ml, len(tTitle), len(iTitle), len(dTitle)) + minPadding + 2

	if ml == 0 {
		switch m.cursor.col {
		case 0:
			tTitle = "* " + tTitle
		case 1:
			iTitle = "* " + iTitle
		case 2:
			dTitle = "* " + dTitle
		}
	}

	tTitle, iTitle, dTitle = addPadding(tTitle, space), addPadding(iTitle, space), addPadding(dTitle, space)

	tTitle = redText.Render(tTitle)
	iTitle = yellowText.Render(iTitle)
	dTitle = greenText.Render(dTitle)

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 20
		height = 10
	}

	var s strings.Builder

	s.WriteString("╔" + strings.Repeat("═", space) + "╦" + strings.Repeat("═", space) + "╦" + strings.Repeat("═", space) + "╗\n")

	s.WriteString(fmt.Sprintf("║%s║%s║%s║", tTitle, iTitle, dTitle) + "\n")

	s.WriteString("╠" + strings.Repeat("═", space) + "╬" + strings.Repeat("═", space) + "╬" + strings.Repeat("═", space) + "╣\n")

	sortTasks(&m.formattedTasks)
	tt := transpose(m.formattedTasks)
	for i := range tt {
		tTask, iTask, dTask := tt[i][0], tt[i][1], tt[i][2]

		tDueDate := m.structuredTasks[tTask].Due
		iDueDate := m.structuredTasks[iTask].Due
		dDueDate := m.structuredTasks[dTask].Due

		tBlocked := m.structuredTasks[tTask].Blocked
		iBlocked := m.structuredTasks[iTask].Blocked
		dBlocked := m.structuredTasks[dTask].Blocked

		tTask, iTask, dTask = m.structuredTasks[tTask].Desc, m.structuredTasks[iTask].Desc, m.structuredTasks[dTask].Desc

		if m.cursor.row == i {
			switch m.cursor.col {
			case 0:
				tTask = "> " + tTask
			case 1:
				iTask = "> " + iTask
			case 2:
				dTask = "> " + dTask
			}
		}
		tTask = addPadding(tTask, space)
		iTask = addPadding(iTask, space)
		dTask = addPadding(dTask, space)

		if tDueDate == time.Now().Format("2006-01-02") && !tBlocked {
			tTask = blueText.Render(tTask)
		}

		if iDueDate == time.Now().Format("2006-01-02") && !iBlocked {
			iTask = blueText.Render(iTask)
		}

		if dDueDate == time.Now().Format("2006-01-02") && !dBlocked {
			dTask = blueText.Render(dTask)
		}

		if tBlocked {
			tTask = redText.Render(tTask)
		}
		if iBlocked {
			iTask = redText.Render(iTask)
		}
		if dBlocked {
			dTask = redText.Render(dTask)
		}

		s.WriteString(fmt.Sprintf("║%s║%s║%s║\n", tTask, iTask, dTask))
	}

	s.WriteString("╚" + strings.Repeat("═", space) + "╩" + strings.Repeat("═", space) + "╩" + strings.Repeat("═", space) + "╝\n")

	s.WriteString(helpStyle.Render("\na: (a)dd • p: (p)romote • r: (r)egress • d: (d)elete • e: (e)dit • s: (s)how • q: (q)uit\n"))

	content := s.String()
	contentHeight := strings.Count(content, "\n") + 1
	topPadding := (height - contentHeight) / 8
	style := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		PaddingTop(topPadding)

	return style.Render(content)
}
