package screen

import "github.com/agiles231/gotui/terminal"

// Cell represents a single character cell on the screen
type Cell struct {
	Rune  rune
	Style terminal.Style
}

// NewCell creates a new cell with the given rune and style
func NewCell(r rune, style terminal.Style) Cell {
	return Cell{
		Rune:  r,
		Style: style,
	}
}

// EmptyCell returns an empty cell with default style
func EmptyCell() Cell {
	return Cell{
		Rune:  ' ',
		Style: terminal.DefaultStyle(),
	}
}

// Equals checks if two cells are identical
func (c Cell) Equals(other Cell) bool {
	return c.Rune == other.Rune && c.Style.Equals(other.Style)
}

// IsEmpty returns true if the cell is a space with default style
func (c Cell) IsEmpty() bool {
	return c.Rune == ' ' && c.Style.Equals(terminal.DefaultStyle())
}

// WithRune returns a copy of the cell with a different rune
func (c Cell) WithRune(r rune) Cell {
	c.Rune = r
	return c
}

// WithStyle returns a copy of the cell with a different style
func (c Cell) WithStyle(s terminal.Style) Cell {
	c.Style = s
	return c
}

// WithFG returns a copy of the cell with a different foreground color
func (c Cell) WithFG(color terminal.Color) Cell {
	c.Style = c.Style.WithFG(color)
	return c
}

// WithBG returns a copy of the cell with a different background color
func (c Cell) WithBG(color terminal.Color) Cell {
	c.Style = c.Style.WithBG(color)
	return c
}

