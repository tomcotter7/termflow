package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomcotter7/termflow/internal/storage"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Cursor struct {
	row int
	col int
}

func (c *Cursor) AdjustRow(formattedTasks [3][]string) {
	newColHeight := max(len(formattedTasks[c.col])-1, 0)
	if c.row > newColHeight {
		c.row = newColHeight
	}
}

func (c *Cursor) IncCol(formattedTasks [3][]string) {
	if c.col < 2 {

		c.col++
		c.AdjustRow(formattedTasks)
	}
}

func (c *Cursor) DecCol(formattedTasks [3][]string) {
	if c.col > 0 {
		c.col--
		c.AdjustRow(formattedTasks)
	}
}

func (c *Cursor) IncRow(maxLen int) {
	if c.row < maxLen {
		c.row++
	}
}

func (c *Cursor) DecRow() {
	if c.row > 0 {
		c.row--
	}
}

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
	ti.focusedIdx = max(0, ti.focusedIdx-1)
}

func (ti *TextInputs) increaseFocusedIndex() {
	ti.focusedIdx = min(ti.focusedIdx+1, len(ti.ti))
}

func (ti *TextInputs) focusTextInput(idx int) {
	ti.ti[idx].Focus()
	ti.ti[idx].PromptStyle = focusedStyle
	ti.ti[idx].TextStyle = focusedStyle
}

func (ti *TextInputs) deFocusTextInput(idx int) {
	ti.ti[idx].Blur()
	ti.ti[idx].PromptStyle = noStyle
	ti.ti[idx].TextStyle = noStyle
}

func (ti *TextInputs) updateTextInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(ti.ti))

	for i := range ti.ti {
		ti.ti[i], cmds[i] = ti.ti[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

type CreateProjectInput struct {
	textInputs TextInputs
}

type CreateTaskInput struct {
	textInputs  TextInputs
	inputTaskId string
}

type Mode int

const (
	NormalMode Mode = iota
	InputMode
	ShowMode
	CommandMode
	NewProjectMode
	SwitchProjectMode
	ErrorMode
)

type model struct {
	handler         *storage.Handler
	structuredTasks map[string]storage.Task
	formattedTasks  [3][]string
	cursor          Cursor
	mode            Mode
	help            bool
	commands        list.Model
	projects        list.Model
	height          int
	width           int
	project         string
	error           error

	createTaskInput    CreateTaskInput
	createProjectInput CreateProjectInput
}

func formatTasks(tasks map[string]storage.Task) [3][]string {
	taskIds := [3][]string{{}, {}, {}}
	for id, task := range tasks {
		switch task.Status {
		case "todo":
			taskIds[0] = append(taskIds[0], id)
		case "inprogress":
			taskIds[1] = append(taskIds[1], id)
		case "done":
			taskIds[2] = append(taskIds[2], id)
		}
	}

	return taskIds
}

func NewModel() model {
	h, err := storage.New()
	if err != nil {
		log.Fatal("Unable to create Handler object:", err)
	}

	structuredTasks, err := h.LoadTasks("default.json")
	if err != nil {
		log.Fatal("Unable to load initial model:", err)
	}

	task_inputs := make([]textinput.Model, 3)
	for i := range task_inputs {
		t := textinput.New()
		switch i {
		case 0:
			t.Placeholder = "Short Description"
		case 1:
			t.Placeholder = "Full Description"
		case 2:
			t.Placeholder = "Due Date"
		}

		task_inputs[i] = t
	}
	cti := CreateTaskInput{textInputs: TextInputs{ti: task_inputs}}

	project_inputs := make([]textinput.Model, 1)
	t := textinput.New()
	t.Placeholder = "Project Name"
	project_inputs[0] = t

	pti := CreateProjectInput{textInputs: TextInputs{ti: project_inputs}}

	commandItems := []list.Item{
		item{title: "Print", desc: "Produce a Carmack-like .plan file with all done tasks."},
		item{title: "Clear", desc: "Delete all done tasks."},
		item{title: "Create Project", desc: "Create a new project & switch to it."},
		item{title: "Switch to Project", desc: "Switch to a different project."},
	}

	commands := list.New(commandItems, list.NewDefaultDelegate(), 0, 0)
	commands.Title = "Available Commands"

	projectNames, err := h.ListAllProjects()
	projectItems := make([]list.Item, len(projectNames))

	for i, project := range projectNames {
		projectItems[i] = item{title: project, desc: ""}
	}

	projects := list.New(projectItems, list.NewDefaultDelegate(), 0, 0)
	projects.Title = "Available Projects"

	return model{
		handler:            h,
		structuredTasks:    structuredTasks,
		formattedTasks:     formatTasks(structuredTasks),
		createTaskInput:    cti,
		createProjectInput: pti,
		mode:               NormalMode,
		help:               false,
		commands:           commands,
		project:            "default",
		projects:           projects,
	}
}
