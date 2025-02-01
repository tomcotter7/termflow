package ui

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tomcotter7/termflow/internal/storage"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	redText             = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#B33B3B"))
	yellowText          = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#EABD30"))
	greenText           = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#30EA40"))
	blueText            = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2563BE"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

var columnNames = map[int]string{
	0: "todo",
	1: "inprogress",
	2: "done",
}

type (
	errMsg error
)

func randomId() string {
	b := make([]byte, 5)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func sortTasks(l *[3][]string) {
	for i := range l {
		sort.Strings(l[i])
	}
}

func transpose(l [3][]string) [][3]string {
	max_len := max(len(l[0]), len(l[1]), len(l[2]))
	l_t := make([][3]string, max_len)

	for i := 0; i < max_len; i++ {
		l_t[i] = [3]string{}
		for j := 0; j < 3; j++ {
			if i < len(l[j]) {
				l_t[i][j] = l[j][i]
			}
		}
	}

	return l_t
}

func maxTaskLength(tasks map[string]storage.Task) int {
	maxLength := 0
	for _, v := range tasks {
		maxLength = max(maxLength, len(v.Desc))
	}

	return maxLength
}

func addPadding(ipt string, space int) string {
	diff := max(space, len(ipt)) - len(ipt)

	lpadding := (diff / 2)
	rpadding := max(space-len(ipt)-lpadding, 0)

	return strings.Repeat(" ", lpadding) + ipt + strings.Repeat(" ", rpadding)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.mode == "input" {
		return m.handleInputModelUpdate(msg)
	} else if m.mode == "normal" {
		return m.handleNormalModelUpdate(msg)
	} else if m.mode == "show" {
		return m.handleShowModeUpdate(msg)
	}

	return m, nil
}

func (m model) View() string {
	if m.mode == "input" {
		return m.inputModeView()
	} else if m.mode == "normal" {
		return m.normalModeView()
	} else if m.mode == "show" {
		return m.showModeView()
	}

	return ""
}
