package ui

import (
	"log"

	"github.com/tomcotter7/termflow/internal/storage"

	"github.com/charmbracelet/bubbles/textinput"
)

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

type model struct {
	handler         *storage.Handler
	structuredTasks map[string]storage.Task
	formattedTasks  [3][]string
	cursor          Cursor
	textInputs      []textinput.Model
	focusedIndex    int
	inputTaskId     string
	mode            string
	err             error
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
	ti := make([]textinput.Model, 3)
	h, err := storage.New()
	if err != nil {
		log.Fatal("Unable to create Handler object:", err)
	}

	structuredTasks, err := h.LoadTasks("default.json")
	if err != nil {
		log.Fatal("Unable to load initial model:", err)
	}

	for i := range ti {
		t := textinput.New()
		switch i {
		case 0:
			t.Placeholder = "Short Description"
		case 1:
			t.Placeholder = "Full Description"
		case 2:
			t.Placeholder = "Due Date"
		}

		ti[i] = t
	}

	return model{
		handler:         h,
		structuredTasks: structuredTasks,
		formattedTasks:  formatTasks(structuredTasks),
		textInputs:      ti,
		mode:            "normal",
	}
}
