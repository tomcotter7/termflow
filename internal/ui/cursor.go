package ui

import "github.com/tomcotter7/termflow/internal/storage"

type Cursor struct {
	row int
	col int
}

func (c *Cursor) MoveToFirstRow() {
	c.row = 0
}

func (c *Cursor) MoveToLastRow(maxLen int) {
	c.row = maxLen
}

func (c *Cursor) MoveToFirstCol(formattedTasks [4][]storage.Task) {
	if c.col != 0 {
		c.col = 0
		c.AdjustRow(formattedTasks)
	}
}

func (c *Cursor) MoveToLastCol(formattedTasks [4][]storage.Task) {
	if c.col < 3 {
		c.col = 3
		c.AdjustRow(formattedTasks)
	}
}

func (c *Cursor) AdjustRow(formattedTasks [4][]storage.Task) {
	newColHeight := max(len(formattedTasks[c.col])-1, 0)
	if c.row > newColHeight {
		c.row = newColHeight
	}
}

func (c *Cursor) IncCol(formattedTasks [4][]storage.Task) {
	if c.col < 3 {

		c.col++
		c.AdjustRow(formattedTasks)
	}
}

func (c *Cursor) DecCol(formattedTasks [4][]storage.Task) {
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
