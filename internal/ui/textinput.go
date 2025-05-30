package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInputs struct {
	ti         []textinput.Model
	focusedIdx int
}

func (ti *TextInputs) onSubmitButton() bool {
	return ti.focusedIdx == len(ti.ti)
}

func (ti *TextInputs) resetTextInputs() {
	ti.focusedIdx = 0
	for i := range ti.ti {
		ti.ti[i].Reset()
		ti.ti[i].Blur()
	}
}

func (ti *TextInputs) decreaseFocusedIndex() {
	ti.deFocusTextInput(ti.focusedIdx)
	ti.focusedIdx = max(0, ti.focusedIdx-1)
	ti.focusTextInput(ti.focusedIdx)
}

func (ti *TextInputs) increaseFocusedIndex() {
	ti.deFocusTextInput(ti.focusedIdx)
	ti.focusedIdx = min(ti.focusedIdx+1, len(ti.ti))

	if ti.focusedIdx < len(ti.ti) {
		ti.focusTextInput(ti.focusedIdx)
	}
}

func (ti *TextInputs) focusTextInput(idx int) {
	if idx < len(ti.ti) {
		ti.ti[idx].Focus()
		ti.ti[idx].PromptStyle = focusedStyle
		ti.ti[idx].TextStyle = focusedStyle
	}
}

func (ti *TextInputs) deFocusTextInput(idx int) {
	if idx < len(ti.ti) {
		ti.ti[idx].Blur()
		ti.ti[idx].PromptStyle = noStyle
		ti.ti[idx].TextStyle = noStyle
	}
}

func (ti *TextInputs) updateTextInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(ti.ti))

	for i := range ti.ti {
		ti.ti[i], cmds[i] = ti.ti[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
