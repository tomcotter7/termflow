package ui

import (
	"log"
	"slices"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/tomcotter7/termflow/internal/storage"
)

type Mode int

const (
	NormalMode Mode = iota
	EditMode
	ShowMode
	CommandMode
	NewProjectMode
	SwitchProjectMode
	AddBragMode
	AddNoteMode
	NotesViewMode
	ParseNoteMode
	ErrorMode
)

type item struct {
	id          string
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title + " " + i.desc }

type CreateProjectForm struct {
	inputs Form
}

type CreateTaskForm struct {
	inputs      Form
	inputTaskId string
}

func (f *CreateTaskForm) PopulateFromTask(task storage.Task) {
	f.inputTaskId = task.ID
	f.inputs.ti[0].SetValue(task.Desc)
	f.inputs.ti[1].SetValue(task.Due)
	f.inputs.ti[2].SetValue(strconv.Itoa(task.Priority))
	f.inputs.ta[0].SetValue(task.FullDesc)
	f.inputs.ta[1].SetValue(task.Result)
}

type AddBragForm struct {
	inputs           Form
	tasksPager       viewport.Model
	taskLookbackDays int
	focusOnPager     bool
}

type AddNoteForm struct {
	inputs Form
	prevID string
}

type model struct {
	handler        *storage.Handler
	tasks          map[string]storage.Task
	formattedTasks [4][]storage.Task
	cursor         Cursor
	notes          map[string]storage.Note
	notesList      list.Model
	mode           Mode
	showHelp       bool
	commands       list.Model
	projects       list.Model
	termHeight     int
	termWidth      int
	activeProject  string
	err            error
	wp             string
	lastKey        string

	createTaskForm    CreateTaskForm
	createProjectForm CreateProjectForm
	addBragForm       AddBragForm
	addNoteForm       AddNoteForm
}

func priorityOrdering(t_a storage.Task, t_b storage.Task) int {
	today := time.Now().Format("2006-01-02")
	if (t_a.Due <= today) && (t_b.Due != today) {
		return -1
	} else if (t_a.Due != today) && (t_b.Due <= today) {
		return 1
	}

	if t_a.Priority < t_b.Priority {
		return -1
	} else if t_a.Priority > t_b.Priority {
		return 1
	}

	if t_a.ID < t_b.ID {
		return -1
	} else if t_a.ID > t_b.ID {
		return 1
	}

	return 0
}

func formatTasks(tasks map[string]storage.Task) [4][]storage.Task {
	fTasks := [4][]storage.Task{{}, {}, {}, {}}

	todoTasks := []storage.Task{}
	ipTasks := []storage.Task{}
	irTasks := []storage.Task{}
	doneTasks := []storage.Task{}

	for _, task := range tasks {
		switch task.Status {
		case storage.StatusTodo:
			todoTasks = append(todoTasks, task)
		case storage.StatusInProgress:
			ipTasks = append(ipTasks, task)
		case storage.StatusInReview:
			irTasks = append(irTasks, task)
		case storage.StatusDone:
			doneTasks = append(doneTasks, task)
		}
	}

	slices.SortFunc(todoTasks, priorityOrdering)
	slices.SortFunc(ipTasks, priorityOrdering)
	slices.SortFunc(irTasks, priorityOrdering)
	slices.SortFunc(doneTasks, priorityOrdering)

	fTasks[0] = todoTasks
	fTasks[1] = ipTasks
	fTasks[2] = irTasks
	fTasks[3] = doneTasks

	return fTasks
}

func createNotesListModel(notes map[string]storage.Note) list.Model {
	notesItems := make([]list.Item, len(notes))

	idx := 0
	for id, note := range notes {
		notesItems[idx] = item{id: id, title: note.Created, desc: note.Content}
		idx++
	}

	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(16)
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(lipgloss.Color("#555555"))

	notesList := list.New(notesItems, delegate, 0, 0)
	notesList.Title = "Notes"
	return notesList
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
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(lipgloss.Color("#555555"))

	projects := list.New(projectItems, delegate, 0, 0)
	projects.Title = "Available Projects"

	return projects, nil
}

func newCommandsListModel() list.Model {
	commandItems := []list.Item{
		item{title: "New Note", desc: "Add a WIP note to be later parsed into todo item."},
		item{title: "View Notes", desc: "View all current notes."},
		item{title: "Print", desc: "Produce a Carmack-like .plan file with all done tasks."},
		item{title: "Clear", desc: "Delete all done tasks."},
		item{title: "Create Project", desc: "Create a new project & switch to it."},
		item{title: "Switch to Project", desc: "Switch to a different project."},
		item{title: "Brag", desc: "Add an item to your brag document."},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(lipgloss.Color("#555555"))

	commands := list.New(commandItems, delegate, 0, 0)
	commands.Title = "Available Commands"
	return commands
}

func newAddBragForm() AddBragForm {
	text_inputs := make([]textinput.Model, 1)
	ti := textinput.New()
	ti.Placeholder = "Brag Title"
	text_inputs[0] = ti
	text_areas := make([]textarea.Model, 1)
	ta := textarea.New()
	ta.Placeholder = "Brag Content"
	ta.CharLimit = 0
	text_areas[0] = ta

	abf := AddBragForm{
		inputs:           Form{ti: text_inputs, ta: text_areas},
		tasksPager:       viewport.New(10, 10),
		taskLookbackDays: 7,
	}
	return abf
}

func newAddNoteForm() AddNoteForm {
	text_areas := make([]textarea.Model, 1)
	ta := textarea.New()
	ta.Placeholder = "Note Content..."
	ta.CharLimit = 0
	text_areas[0] = ta

	anf := AddNoteForm{
		inputs: Form{ta: text_areas},
	}

	return anf
}

func newCreateTaskForm() CreateTaskForm {
	text_inputs := make([]textinput.Model, 3)
	for i := range text_inputs {
		t := textinput.New()
		switch i {
		case 0:
			t.Placeholder = "Short Description"
		case 1:
			t.Placeholder = "Due Date"
		case 2:
			t.Placeholder = "Priority"
		}
		text_inputs[i] = t
	}

	text_areas := make([]textarea.Model, 2)
	t := textarea.New()
	t.Placeholder = "Full Description"
	t.CharLimit = 0
	text_areas[0] = t

	t2 := textarea.New()
	t2.Placeholder = "Results"
	t2.CharLimit = 0
	text_areas[1] = t2

	cti := CreateTaskForm{inputs: Form{ti: text_inputs, ta: text_areas}}
	return cti
}

func newCreateProjectForm() CreateProjectForm {
	project_inputs := make([]textinput.Model, 1)
	t := textinput.New()
	t.Placeholder = "Project Name"
	project_inputs[0] = t

	pti := CreateProjectForm{inputs: Form{ti: project_inputs}}
	return pti
}

func (m model) setError(err error) model {
	m.err = err
	m.mode = ErrorMode
	return m
}

func NewModel() model {
	h, err := storage.New()
	if err != nil {
		log.Fatal("Unable to create Handler object:", err)
	}

	currentProject := h.GetCurrent()
	structuredTasks, err := h.LoadTasks(currentProject + ".json")
	if err != nil {
		log.Fatal("Unable to load initial tasks:", err)
	}

	notes, err := h.LoadNotes(currentProject + "_notes.json")
	if err != nil {
		log.Fatal("Unable to load initial notes:", err)
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
		notes:             notes,
		notesList:         createNotesListModel(notes),
		createTaskForm:    newCreateTaskForm(),
		createProjectForm: newCreateProjectForm(),
		addBragForm:       newAddBragForm(),
		addNoteForm:       newAddNoteForm(),
		mode:              NormalMode,
		showHelp:          false,
		commands:          commands,
		activeProject:     currentProject,
		projects:          projects,
	}
}
