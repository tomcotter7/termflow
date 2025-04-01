package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/tomcotter7/termflow/internal/storage"
)

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

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type CreateProjectForm struct {
	textInputs TextInputs
}

type CreateTaskForm struct {
	textInputs  TextInputs
	inputTaskId string
}

type model struct {
	handler        *storage.Handler
	tasks          map[string]storage.Task
	formattedTasks [3][]string
	cursor         Cursor
	mode           Mode
	showHelp       bool
	commands       list.Model
	projects       list.Model
	termHeight     int
	termWidth      int
	activeProject  string
	err            error

	createTaskForm    CreateTaskForm
	createProjectForm CreateProjectForm
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

func newProjectListModel(h *storage.Handler) (list.Model, error) {
	projectNames, err := h.ListAllProjects()
	if err != nil {
		return list.Model{}, err
	}
	projectItems := make([]list.Item, len(projectNames))

	for i, project := range projectNames {
		projectItems[i] = item{title: project, desc: ""}
	}

	projects := list.New(projectItems, list.NewDefaultDelegate(), 0, 0)
	projects.Title = "Available Projects"

	return projects, nil
}

func newCommandsListModel() list.Model {
	commandItems := []list.Item{
		item{title: "Print", desc: "Produce a Carmack-like .plan file with all done tasks."},
		item{title: "Clear", desc: "Delete all done tasks."},
		item{title: "Create Project", desc: "Create a new project & switch to it."},
		item{title: "Switch to Project", desc: "Switch to a different project."},
	}

	commands := list.New(commandItems, list.NewDefaultDelegate(), 0, 0)
	commands.Title = "Available Commands"
	return commands
}

func newCreateTaskForm() CreateTaskForm {
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
	cti := CreateTaskForm{textInputs: TextInputs{ti: task_inputs}}
	return cti
}

func newCreateProjectForm() CreateProjectForm {
	project_inputs := make([]textinput.Model, 1)
	t := textinput.New()
	t.Placeholder = "Project Name"
	project_inputs[0] = t

	pti := CreateProjectForm{textInputs: TextInputs{ti: project_inputs}}
	return pti
}

func NewModel() model {
	h, err := storage.New()
	if err != nil {
		log.Fatal("Unable to create Handler object:", err)
	}

	currentProject := h.GetCurrent()
	structuredTasks, err := h.LoadTasks(currentProject + ".json")
	if err != nil {
		log.Fatal("Unable to load initial model:", err)
	}

	commands := newCommandsListModel()
	projects, err := newProjectListModel(h)
	if err != nil {
		log.Printf("Warning: Unable to load list of projects: %v", err)
		projects = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
		projects.Title = "Available Projects"
	}

	return model{
		handler:           h,
		tasks:             structuredTasks,
		formattedTasks:    formatTasks(structuredTasks),
		createTaskForm:    newCreateTaskForm(),
		createProjectForm: newCreateProjectForm(),
		mode:              NormalMode,
		showHelp:          false,
		commands:          commands,
		activeProject:     currentProject,
		projects:          projects,
	}
}
