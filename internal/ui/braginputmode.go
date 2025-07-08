package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

func (m model) getRecentlyCompletedTasks(lookbackDays int) string {
	today := time.Now()
	completed_tasks := make([]string, lookbackDays)

	for i := lookbackDays; i >= 0; i-- {
		filename := (today.AddDate(0, 0, -i).Format("2006-01-02")) + m.activeProject + ".plan"
		data, err := m.handler.ReadPlanFile(filename)
		if err != nil {
			continue
		}

		completed_tasks = append(completed_tasks, string(data))
	}

	return strings.TrimSpace(strings.Join(completed_tasks, "\n"))
}

func (m model) addBragModeView() string {
	m.addBragForm.tasksPager.Width = m.termWidth / 4
	m.addBragForm.tasksPager.Height = m.termHeight / 2
	pager := fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.addBragForm.tasksPager.View(), m.footerView())
	form := m.addBragForm.inputs.buildFormView()

	if m.addBragForm.focusOnPager {
		pager = lipgloss.NewStyle().BorderStyle(lipgloss.ThickBorder()).BorderForeground(lipgloss.Color("63")).Render(pager)
		form = lipgloss.NewStyle().Height(m.termHeight / 2).BorderStyle(lipgloss.ThickBorder()).BorderForeground(lipgloss.Color("0")).Render(form)
	} else {
		pager = lipgloss.NewStyle().BorderStyle(lipgloss.ThickBorder()).BorderForeground(lipgloss.Color("0")).Render(pager)
		form = lipgloss.NewStyle().Height(m.termHeight / 2).BorderStyle(lipgloss.ThickBorder()).BorderForeground(lipgloss.Color("63")).Render(form)
	}

	content := lipgloss.JoinHorizontal(lipgloss.Top, form, pager)

	return lipgloss.NewStyle().Width(m.termWidth - 4).Align(lipgloss.Center).Render(content)
}

func (m model) headerView() string {
	title := titleStyle.Render(fmt.Sprintf("Tasks completed in the last %d days", m.addBragForm.taskLookbackDays))
	line := strings.Repeat("─", max(0, m.addBragForm.tasksPager.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.addBragForm.tasksPager.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.addBragForm.tasksPager.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m model) handleAddBragModeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termHeight = msg.Height
		m.termWidth = msg.Width
	case tea.KeyMsg:
		switch k := msg.String(); k {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.mode = NormalMode
			return m, nil
		case "right":
			m.addBragForm.focusOnPager = true
		case "left":
			m.addBragForm.focusOnPager = false
		case "enter":
			if !m.addBragForm.focusOnPager && m.addBragForm.inputs.onSubmitButton() {
				brag := m.addBragForm.inputs.ta[0].Value()
				m.handler.SaveBragFile(brag)
				m.mode = NormalMode
				m.addBragForm.inputs.reset()
				return m, nil
			}
		case "tab", "shift+tab":
			if !m.addBragForm.focusOnPager {
				if k == "shift+tab" {
					m.addBragForm.inputs.decreaseFocusedIndex()
				} else {
					m.addBragForm.inputs.increaseFocusedIndex()
				}
			}
		}
	}

	if m.addBragForm.focusOnPager {
		m.addBragForm.tasksPager, cmd = m.addBragForm.tasksPager.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	} else {
		cmd = m.addBragForm.inputs.updateInputs(msg)
		return m, cmd
	}
}
