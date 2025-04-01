package ui

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
