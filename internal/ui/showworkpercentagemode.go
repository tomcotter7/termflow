package ui

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getCurrentWorkPercentage() string {
	sheetName := time.Now().Format("2006January")
	cellPos := string(rune(time.Now().Day()+65)) + "34"

	spreadsheetID := "1g9gV_a_j03qAw5ZepbY9RfnButjeXeFbB_cKq7BlTp8"

	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=csv&sheet=%s&range=%s", spreadsheetID, sheetName, cellPos)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	workPercentage := records[0][0]

	return workPercentage
}

func (m model) showWPModeView() string {
	wp := "Amount of day spent 🔒in: " + blueText.Render(getCurrentWorkPercentage())

	contentHeight := 1

	topPadding := (m.termHeight - contentHeight) / 8

	style := lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Padding(topPadding).Bold(true)

	return style.Render(wp)
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
