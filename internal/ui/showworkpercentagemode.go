package ui

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getCellPos() string {
	day := time.Now().Day()
	col := ""
	if day >= 26 {
		col += "A"
		day -= 26
	}

	col += string(rune(day + 65))

	return col + "34"
}

func (m *model) getCurrentWorkPercentage() string {
	sheetName := time.Now().Format("2006January")
	cellPos := getCellPos()

	spreadsheetID := "1g9gV_a_j03qAw5ZepbY9RfnButjeXeFbB_cKq7BlTp8"

	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=csv&sheet=%s&range=%s", spreadsheetID, sheetName, cellPos)

	resp, err := http.Get(url)
	if err != nil {
		m.err = err
		return err.Error()
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		m.err = err
		m.mode = ErrorMode
		return err.Error()
	}

	workPercentage := records[0][0]

	return workPercentage
}

func (m model) getDoneFromAllProjects() string {
	projects, _ := m.handler.ListAllProjects()

	totalDone := 0
	for i := range projects {
		project := projects[i]
		tasks, _ := m.handler.LoadTasks(project + ".json")
		fTasks := formatTasks(tasks)
		totalDone += len(fTasks[2])
	}

	return strconv.Itoa(totalDone)
}

func (m model) showWPModeView() string {
	var s strings.Builder
	s.WriteString("Amount of day spent 🔒in: " + blueText.Render(m.getCurrentWorkPercentage()) + "\n")
	s.WriteString("Tasks completed: " + blueText.Render(m.getDoneFromAllProjects()))

	content := s.String()

	contentHeight := strings.Count(content, "\n")

	topPadding := (m.termHeight - contentHeight) / 8

	style := lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Padding(topPadding).Bold(true)

	return style.Render(content)
}

func (m model) handleWPModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.mode = NormalMode
		}
	}
	return m, nil
}
