package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Form struct {
	ti         []textinput.Model
	ta         []textarea.Model
	taHidden   []bool
	focusedIdx int
}

func (f *Form) updateTextAreas(termWidth int, termHeight int) {
	for i := range f.ta {
		f.ta[i].SetWidth(termWidth / 2)
		f.ta[i].SetHeight(termHeight / 4)
	}
}

func (f Form) getTotalHiddenTextAreas() int {
	totalHidden := 0
	for i := range f.taHidden {
		if f.taHidden[i] {
			totalHidden++
		}
	}
	return totalHidden
}

func (f Form) buildFormView() string {
	var b strings.Builder

	for i := range f.ti {
		b.WriteString(f.ti[i].View())
		b.WriteRune('\n')
	}
	b.WriteRune('\n')
	for i := range f.ta {
		if f.taHidden[i] {
			continue
		}
		b.WriteString(f.ta[i].View())
		b.WriteRune('\n')
		b.WriteRune('\n')
	}

	button := &blurredButton
	if f.onSubmitButton() {
		button = &focusedButton
	}

	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	content := b.String()
	return lipgloss.NewStyle().Align(lipgloss.Center).Render(content)
}

func (f *Form) onSubmitButton() bool {
	return f.focusedIdx == (len(f.ti) + len(f.ta))
}

func (f *Form) reset() {
	f.focusedIdx = 0
	for i := range f.ti {
		f.ti[i].Reset()
		f.ti[i].Blur()
	}

	for i := range f.ta {
		f.ta[i].Reset()
		f.ta[i].Blur()
		f.taHidden[i] = false
	}
}

func (f *Form) deFocusInput(idx int) {
	if idx < len(f.ti) {
		f.ti[idx].Blur()
		f.ti[idx].PromptStyle = noStyle
		f.ti[idx].TextStyle = noStyle
	} else if idx < len(f.ti)+len(f.ta) {
		idx -= len(f.ti)
		f.ta[idx].Blur()
	}
}

func (f *Form) focusInput(idx int) {
	f.focusedIdx = idx
	if idx < len(f.ti) {
		f.ti[idx].Focus()
		f.ti[idx].PromptStyle = focusedStyle
		f.ti[idx].TextStyle = focusedStyle
	} else if idx < len(f.ti)+len(f.ta) {
		adjustedIdx := idx - len(f.ti)
		if f.taHidden[adjustedIdx] {
			return
		}
		f.ta[adjustedIdx].Focus()
	}
}

func (f *Form) decreaseFocusedIndex() {
	f.deFocusInput(f.focusedIdx)
	f.focusedIdx = max(0, f.focusedIdx-1)
	for f.focusedIdx >= len(f.ti) && f.focusedIdx < len(f.ti)+len(f.ta) && f.taHidden[f.focusedIdx-len(f.ti)] {
		f.focusedIdx--
	}
	f.focusInput(f.focusedIdx)
}

func (f *Form) increaseFocusedIndex() {
	f.deFocusInput(f.focusedIdx)
	f.focusedIdx = min(f.focusedIdx+1, len(f.ti)+len(f.ta))

	for f.focusedIdx >= len(f.ti) && f.focusedIdx < len(f.ti)+len(f.ta) && f.taHidden[f.focusedIdx-len(f.ti)] {
		f.focusedIdx++
	}
	if f.focusedIdx < len(f.ti)+len(f.ta) {
		f.focusInput(f.focusedIdx)
	}
}

func (f *Form) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.ti)+len(f.ta))

	for i := range f.ti {
		f.ti[i], cmds[i] = f.ti[i].Update(msg)
	}
	for i := range f.ta {
		if f.taHidden[i] {
			continue
		}
		f.ta[i], cmds[len(f.ti)+i] = f.ta[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
