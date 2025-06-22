package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle      = lipgloss.NewStyle()
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Align(lipgloss.Center)

	doneStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Align(lipgloss.Center)
	excludedDoneStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Align(lipgloss.Center).Strikethrough(true)

	redText    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#B33B3B"))
	yellowText = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#EABD30"))
	orangeText = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff3c00"))
	greenText  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#30EA40"))
	blueText   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2563BE"))

	blueTextRedBackground   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2563BE")).Background(lipgloss.Color("#B33B3B"))
	orangeTextRedBackground = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff3c00")).Background(lipgloss.Color("#B33B3B"))
	redBackground           = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#B33B3B"))
	blueBackground          = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#2563BE"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = blurredStyle.Render("[ Submit ]")
)

var columnNames = map[int]string{
	0: "todo",
	1: "inprogress",
	2: "done",
}

type (
	errMsg error
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil || m.mode == ErrorMode {
		return m.handleErrorModeUpdate(msg)
	}

	switch m.mode {
	case NormalMode:
		return m.handleNormalModelUpdate(msg)
	case InputMode:
		return m.handleInputModelUpdate(msg)
	case ShowMode:
		return m.handleShowModeUpdate(msg)
	case CommandMode:
		return m.handleCommandModeUpdate(msg)
	case NewProjectMode:
		return m.handleNewProjectModeUpdate(msg)
	case SwitchProjectMode:
		return m.handleSwitchProjectModeUpdates(msg)
	case ShowWorkPercentageMode:
		return m.handleWPModeUpdate(msg)
	case AddBragMode:
		return m.handleAddBragModeUpdate(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.mode = NormalMode
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil || m.mode == ErrorMode {
		return m.errorModeView()
	}

	switch m.mode {
	case NormalMode:
		return m.normalModeView()
	case InputMode:
		return m.inputModeView()
	case ShowMode:
		return m.showModeView()
	case CommandMode:
		return m.commandModeView()
	case NewProjectMode:
		return m.newProjectModeView()
	case SwitchProjectMode:
		return m.switchProjectModeView()
	case ShowWorkPercentageMode:
		return m.showWPModeView()
	case AddBragMode:
		return m.addBragModeView()
	}

	return "Something has gone wrong!\n\nReport at bug at https://github.com/tomcotter7/termflow/issues\n\n Press (q) to go back."
}
