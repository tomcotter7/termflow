package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tomcotter7/termflow/internal/storage"
)

func transpose(l [4][]storage.Task) [][4]storage.Task {
	max_len := max(len(l[0]), len(l[1]), len(l[2]))
	l_t := make([][4]storage.Task, max_len)
	for i := range max_len {
		l_t[i] = [4]storage.Task{}
		for j := range 4 {
			if i < len(l[j]) {
				l_t[i][j] = l[j][i]
			}
		}
	}

	return l_t
}

func addPadding(ipt string, space int, title bool) string {
	diff := space - len(ipt)

	if diff <= 0 {
		newS := " " + ipt[:(space-4)] + "..."
		return newS
	}

	if title {
		lpadding := (diff / 2)
		rpadding := max(space-len(ipt)-lpadding, 0)
		return strings.Repeat(" ", lpadding) + ipt + strings.Repeat(" ", rpadding)
	}

	return " " + ipt + strings.Repeat(" ", (diff-1))
}

func getLongestTaskLength(tasks map[string]storage.Task) int {
	maxLength := 0
	for _, v := range tasks {
		maxLength = max(maxLength, len(v.Desc))
	}

	return maxLength
}

func (m *model) switchToEditMode(task storage.Task, focusIdx int) {
	m.mode = EditMode
	m.createTaskForm.inputTaskId = task.ID
	m.createTaskForm.inputs.ti[0].SetValue(task.Desc)
	m.createTaskForm.inputs.ti[1].SetValue(task.Due)
	m.createTaskForm.inputs.ti[2].SetValue(strconv.Itoa(task.Priority))
	m.createTaskForm.inputs.ta[0].SetValue(task.FullDesc)
	m.createTaskForm.inputs.ta[1].SetValue(task.Result)
	m.createTaskForm.inputs.focusInput(focusIdx)
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
		case "0":
			m.cursor.MoveToFirstCol(m.formattedTasks)
		case "$":
			m.cursor.MoveToLastCol(m.formattedTasks)
		case "p":

			if len(m.formattedTasks[m.cursor.col]) <= m.cursor.row {
				return m, nil
			}

			item := m.formattedTasks[m.cursor.col][m.cursor.row]

			if task, exists := m.tasks[item.ID]; exists && m.cursor.col < len(columnNames)-1 {
				task.Status = columnNames[m.cursor.col+1]
				m.tasks[item.ID] = task
				m.saveAndUpdateTasks()
				m.cursor.IncCol(m.formattedTasks)
				if m.cursor.col == len(columnNames)-1 {
					m.switchToEditMode(task, 4)
				}
			}

		case "r":
			if len(m.formattedTasks[m.cursor.col]) <= m.cursor.row {
				return m, nil
			}
			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.tasks[item.ID]; exists && m.cursor.col > 0 {
				task.Status = columnNames[m.cursor.col-1]
				m.tasks[item.ID] = task
				m.saveAndUpdateTasks()
				m.cursor.DecCol(m.formattedTasks)
			}
		case "b":
			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.tasks[item.ID]; exists {
				task.Blocked = !task.Blocked
				m.tasks[item.ID] = task
				m.saveAndUpdateTasks()
			}

		case "a":
			m.mode = EditMode
			m.createTaskForm.inputs.focusInput(0)
			m.createTaskForm.inputTaskId = ""
		case "e":
			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.tasks[item.ID]; exists {
				m.switchToEditMode(task, 0)
			}
		case "s", "enter":
			m.mode = ShowMode
		case "t":
			if len(m.formattedTasks[m.cursor.col]) == 0 {
				return m, nil
			}
			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.tasks[item.ID]; exists {

				today := time.Now().Format("2006-01-02")

				if task.Due != today {
					task.Due = today
				} else {
					task.Due = "none"
				}

				m.tasks[item.ID] = task
				m.saveAndUpdateTasks()
			}
		case "d":
			if len(m.formattedTasks[m.cursor.col]) == 0 {
				return m, nil
			}
			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			delete(m.tasks, item.ID)
			m.saveAndUpdateTasks()
		case "i":
			if len(m.formattedTasks[m.cursor.col]) == 0 {
				return m, nil
			}

			item := m.formattedTasks[m.cursor.col][m.cursor.row]
			if task, exists := m.tasks[item.ID]; exists {
				task.IgnoreFromPlan = !task.IgnoreFromPlan
				m.tasks[item.ID] = task
				m.saveAndUpdateTasks()
			}
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
	numColumns := len(columnNames)
	ll := getLongestTaskLength(m.tasks)

	tTitle := fmt.Sprintf("%s (%d)", storage.StatusTodo, len(m.formattedTasks[0]))
	iTitle := fmt.Sprintf("%s (%d)", storage.StatusInProgress, len(m.formattedTasks[1]))
	rTitle := fmt.Sprintf("%s (%d)", storage.StatusInReview, len(m.formattedTasks[2]))
	dTitle := fmt.Sprintf("%s (%d)", storage.StatusDone, len(m.formattedTasks[3]))

	maxTaskLength := (m.termWidth - 8) / numColumns

	minPadding := 0
	space := max(ll, len(tTitle), len(iTitle), len(rTitle), len(dTitle)) + minPadding

	if maxTaskLength > 0 {
		space = min(space, maxTaskLength)
	}

	switch m.cursor.col {
	case 0:
		tTitle = "* " + tTitle
	case 1:
		iTitle = "* " + iTitle
	case 2:
		rTitle = "* " + rTitle
	case 3:
		dTitle = "* " + dTitle
	}

	tTitle, iTitle, rTitle, dTitle = addPadding(tTitle, space, true), addPadding(iTitle, space, true), addPadding(rTitle, space, true), addPadding(dTitle, space, true)

	tTitle = redText.Render(tTitle)
	iTitle = yellowText.Render(iTitle)
	rTitle = blueText.Render(rTitle)
	dTitle = greenText.Render(dTitle)

	var s strings.Builder

	cells := make([]string, numColumns)
	for i := range cells {
		cells[i] = strings.Repeat("═", space)
	}

	s.WriteString("╔" + strings.Join(cells, "╦") + "╗\n")
	s.WriteString(fmt.Sprintf("║%s║%s║%s║%s║", tTitle, iTitle, rTitle, dTitle) + "\n")
	s.WriteString("╠" + strings.Join(cells, "╬") + "╣\n")

	tt := transpose(m.formattedTasks)

	for i := range tt {
		tasks := make([]string, numColumns)
		for j := range numColumns {
			taskData := tt[i][j]

			if len(taskData.Desc) > 0 {
				priorityIndicator := strings.Repeat("!", max((3-(taskData.Priority-1)), 0))
				tasks[j] = priorityIndicator + " " + taskData.Desc
			} else {
				tasks[j] = taskData.Desc
			}

			if m.cursor.row == i && m.cursor.col == j {
				if len(tasks[j]) > 0 {
					tasks[j] = "> " + tasks[j]
				} else {
					tasks[j] = "--+--"
				}
			}
			tasks[j] = addPadding(tasks[j], space, false)

			if j < numColumns-1 {
				if taskData.Due == time.Now().Format("2006-01-02") {
					if !taskData.Blocked {
						tasks[j] = blueText.Render(tasks[j])
					} else {
						tasks[j] = blueTextRedBackground.Render(tasks[j])
					}
				} else if taskData.Status != storage.StatusDone && taskData.Due < time.Now().Format("2006-01-02") {
					if !taskData.Blocked {
						tasks[j] = orangeText.Render(tasks[j])
					} else {
						tasks[j] = orangeTextRedBackground.Render(tasks[j])
					}
				} else {
					if taskData.Blocked {
						tasks[j] = redBackground.Render(tasks[j])
					}
				}
			} else {
				if taskData.IgnoreFromPlan {
					tasks[j] = excludedDoneStyle.Render(tasks[j])
				} else {
					tasks[j] = doneStyle.Render(tasks[j])
				}
			}
		}

		// Build format string and arguments dynamically
		formatStr := "║" + strings.Repeat("%s║", numColumns) + "\n"
		args := make([]interface{}, numColumns)
		for j := range numColumns {
			args[j] = tasks[j]
		}
		s.WriteString(fmt.Sprintf(formatStr, args...))
	}

	s.WriteString("╠" + strings.Join(cells, "╩") + "╣\n")

	projectPadding := (space * numColumns) + 2 - len(m.activeProject)

	s.WriteString("║" + " " + m.activeProject + strings.Repeat(" ", projectPadding) + "║\n")
	if m.err != nil {
		errorPadding := (space * numColumns) + 2 - len(m.err.Error())
		s.WriteString("║" + " " + redBackground.Render(m.err.Error()) + strings.Repeat(" ", errorPadding) + "║\n")
	}
	wpPadding := (space * numColumns) + 2 - len(m.wp)
	s.WriteString("║" + " " + blueBackground.Render(m.wp) + strings.Repeat(" ", wpPadding) + "║\n")

	s.WriteString("╚" + strings.Repeat("═", (space*numColumns)+3) + "╝\n")

	if m.showHelp {
		s.WriteString(helpStyle.Render("\nCommands:\n"))
		s.WriteString(helpStyle.Render("\na: (a)dd • p: (p)romote • r: (r)egress • d: (d)elete • e: (e)dit • s: (s)how • \nt: (t)oday • b: (b)locked • i: (i)gnore from .plan • q: (q)uit • ':': command-mode • ?: hide\n"))
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
