package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Form struct {
	ti         []textinput.Model
	ta         []textarea.Model
	focusedIdx int
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
	if idx < len(f.ti) {
		f.ti[idx].Focus()
		f.ti[idx].PromptStyle = focusedStyle
		f.ti[idx].TextStyle = focusedStyle
	} else if idx < len(f.ti)+len(f.ta) {
		idx -= len(f.ti)
		f.ta[idx].Focus()
	}
}

func (f *Form) decreaseFocusedIndex() {
	f.deFocusInput(f.focusedIdx)
	f.focusedIdx = max(0, f.focusedIdx-1)
	f.focusInput(f.focusedIdx)
}

func (f *Form) increaseFocusedIndex() {
	f.deFocusInput(f.focusedIdx)
	f.focusedIdx = min(f.focusedIdx+1, len(f.ti)+len(f.ta))

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
		f.ta[i], cmds[i] = f.ta[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
