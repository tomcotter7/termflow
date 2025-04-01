package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tomcotter7/termflow/internal/storage"
)

func sortTasks(l *[3][]string) {
	for i := range l {
		sort.Strings(l[i])
	}
}

func transpose(l [3][]string) [][3]string {
	max_len := max(len(l[0]), len(l[1]), len(l[2]))
	l_t := make([][3]string, max_len)
	for i := range max_len {
		l_t[i] = [3]string{}
		for j := range 3 {
			if i < len(l[j]) {
				l_t[i][j] = l[j][i]
			}
		}
	}

	return l_t
}

func addPadding(ipt string, space int, title bool) string {
	diff := max(space, len(ipt)) - len(ipt)

	if title {
		lpadding := (diff / 2)
		rpadding := max(space-len(ipt)-lpadding, 0)
		return strings.Repeat(" ", lpadding) + ipt + strings.Repeat(" ", rpadding)
	}

	return " " + ipt + strings.Repeat(" ", (diff-1))
}

func maxTaskLength(tasks map[string]storage.Task) int {
	maxLength := 0
	for _, v := range tasks {
		maxLength = max(maxLength, len(v.Desc))
	}

	return maxLength
}

func (m *model) saveAndUpdateTasks() {
	m.formattedTasks = formatTasks(m.tasks)
	filename := m.activeProject + ".json"
	m.handler.SaveTasks(filename, m.tasks)
}

func (m model) handleNormalModelUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
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

			if task, exists := m.tasks[item]; exists && m.cursor.col < 2 {
				task.Status = columnNames[m.cursor.col+1]
				m.tasks[item] = task
				m.saveAndUpdateTasks()
			}

			m.cursor.IncCol(m.formattedTasks)

		case "r":
			if len(m.formattedTasks[m.cursor.col]) <= m.cursor.row {
				return m, nil
			}

			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.tasks[item]; exists && m.cursor.col > 0 {
				task.Status = columnNames[m.cursor.col-1]
				m.tasks[item] = task
				m.saveAndUpdateTasks()
			}
			m.cursor.DecCol(m.formattedTasks)
		case "b":
			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.tasks[item]; exists {
				task.Blocked = !task.Blocked
				m.tasks[item] = task
				m.saveAndUpdateTasks()
			}

		case "a":
			m.mode = InputMode
			m.createTaskForm.textInputs.focusTextInput(0)
		case "e":

			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.tasks[item]; exists {
				m.mode = InputMode
				m.createTaskForm.textInputs.focusTextInput(0)
				m.createTaskForm.inputTaskId = item
				m.createTaskForm.textInputs.ti[0].SetValue(task.Desc)
				m.createTaskForm.textInputs.ti[1].SetValue(task.FullDesc)
				m.createTaskForm.textInputs.ti[2].SetValue(task.Due)
			}
		case "s", "enter":
			m.mode = ShowMode
		case "t":
			if len(m.formattedTasks[m.cursor.col]) == 0 {
				return m, nil
			}
			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.tasks[item]; exists {

				today := time.Now().Format("2006-01-02")

				if task.Due != today {
					task.Due = today
				} else {
					task.Due = ""
				}

				m.tasks[item] = task
				m.saveAndUpdateTasks()
			}
		case "d":
			if len(m.formattedTasks[m.cursor.col]) == 0 {
				return m, nil
			}
			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			delete(m.tasks, item)
			m.saveAndUpdateTasks()
		case "?":
			m.showHelp = !m.showHelp
		case ":":
			m.mode = CommandMode
			m.commands.SetSize(m.termWidth-2, m.termHeight-2)
		}
	}

	return m, nil
}

func (m model) normalModeView() string {
	ml := maxTaskLength(m.tasks)

	tTitle, iTitle, dTitle := "todo", "inprogress", "done"

	minPadding := 2
	space := max(ml, len(tTitle), len(iTitle), len(dTitle)) + minPadding + 2

	switch m.cursor.col {
	case 0:
		tTitle = "* " + tTitle
	case 1:
		iTitle = "* " + iTitle
	case 2:
		dTitle = "* " + dTitle
	}

	tTitle, iTitle, dTitle = addPadding(tTitle, space, true), addPadding(iTitle, space, true), addPadding(dTitle, space, true)

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
			taskData := m.tasks[task]

			tasks[j] = taskData.Desc
			if m.cursor.row == i && m.cursor.col == j {
				if len(tasks[j]) > 0 {
					tasks[j] = "> " + tasks[j]
				} else {
					tasks[j] = "--+--"
				}
			}
			tasks[j] = addPadding(tasks[j], space, false)

			if j < 2 {
				if taskData.Due == time.Now().Format("2006-01-02") {
					if !taskData.Blocked {
						tasks[j] = blueText.Render(tasks[j])
					} else {
						tasks[j] = blueTextRedBackground.Render(tasks[j])
					}
				} else {
					if taskData.Blocked {
						tasks[j] = redBackground.Render(tasks[j])
					}
				}
			} else {
				tasks[j] = blurredStyle.Render(tasks[j])
			}
		}

		tTask, iTask, dTask := tasks[0], tasks[1], tasks[2]
		s.WriteString(fmt.Sprintf("║%s║%s║%s║\n", tTask, iTask, dTask))
	}

	s.WriteString("╠" + strings.Repeat("═", space) + "╩" + strings.Repeat("═", space) + "╩" + strings.Repeat("═", space) + "╣\n")

	projectPadding := (space * 3) + 2 - len(m.activeProject) - 1

	s.WriteString("║" + " " + m.activeProject + strings.Repeat(" ", projectPadding) + "║\n")
	if m.err != nil {
		errorPadding := (space * 3) + 2 - len(m.err.Error()) - 1
		s.WriteString("║" + " " + m.err.Error() + strings.Repeat(" ", errorPadding) + "║\n")
	}
	s.WriteString("╚" + strings.Repeat("═", (space*3)+2) + "╝\n")

	if m.showHelp {
		s.WriteString(helpStyle.Render("\nCommands:\n"))
		s.WriteString(helpStyle.Render("\na: (a)dd • p: (p)romote • r: (r)egress • d: (d)elete • e: (e)dit • s: (s)how • \nt: (t)oday • b: (b)locked • q: (q)uit • ':': command-mode • ?: hide\n"))
	} else {
		s.WriteString(helpStyle.Render("\n?: help\n"))
	}

	content := s.String()
	contentHeight := strings.Count(content, "\n") + 1
	topPadding := (m.termHeight - 4 - contentHeight) / 8
	style := lipgloss.NewStyle().
		Width(m.termWidth - 4).
		Align(lipgloss.Center).
		PaddingTop(topPadding)

	return style.Render(content)
}
