package screen

import (
	"fmt"
	"os"
	"strings"

	"github.com/agiles231/gotui/terminal"
)

// Screen manages terminal rendering with double-buffering
type Screen struct {
	terminal *terminal.Terminal
	front    *Buffer // What's currently displayed
	back     *Buffer // What we're drawing to
	width    int
	height   int
	output   *strings.Builder
}

// NewScreen creates a new screen instance
func NewScreen(term *terminal.Terminal) (*Screen, error) {
	width, height, err := term.Size()
	if err != nil {
		return nil, err
	}

	return &Screen{
		terminal: term,
		front:    NewBuffer(width, height),
		back:     NewBuffer(width, height),
		width:    width,
		height:   height,
		output:   &strings.Builder{},
	}, nil
}

// Width returns the screen width
func (s *Screen) Width() int {
	return s.width
}

// Height returns the screen height
func (s *Screen) Height() int {
	return s.height
}

// Buffer returns the back buffer for drawing
func (s *Screen) Buffer() *Buffer {
	return s.back
}

// Resize resizes the screen buffers
func (s *Screen) Resize(width, height int) {
	s.width = width
	s.height = height
	s.front = NewBuffer(width, height)
	s.back = NewBuffer(width, height)
}

// Clear clears the back buffer
func (s *Screen) Clear() {
	s.back.Clear()
}

// Render renders the back buffer to the terminal using diff-based updates
func (s *Screen) Render() {
	s.output.Reset()

	var lastStyle terminal.Style
	styleSet := false
	lastX, lastY := -1, -1

	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			backCell := s.back.Get(x, y)
			frontCell := s.front.Get(x, y)

			// Skip if cell hasn't changed
			if backCell.Equals(frontCell) {
				continue
			}

			// Move cursor if not consecutive
			if x != lastX+1 || y != lastY {
				s.output.WriteString(terminal.CursorMove(x+1, y+1))
			}

			// Update style if changed
			if !styleSet || !backCell.Style.Equals(lastStyle) {
				s.output.WriteString(backCell.Style.Sequence())
				lastStyle = backCell.Style
				styleSet = true
			}

			// Write the character
			s.output.WriteRune(backCell.Rune)

			lastX = x
			lastY = y
		}
	}

	// Reset style at end
	if styleSet {
		s.output.WriteString(terminal.StyleReset)
	}

	// Write to terminal
	if s.output.Len() > 0 {
		os.Stdout.WriteString(s.output.String())
	}

	// Swap buffers by copying back to front
	for y := 0; y < s.height; y++ {
		copy(s.front.cells[y], s.back.cells[y])
	}
}

// ForceRender renders the entire back buffer regardless of changes
func (s *Screen) ForceRender() {
	s.output.Reset()

	var lastStyle terminal.Style
	styleSet := false

	// Move to home
	s.output.WriteString(terminal.CursorHome)

	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			cell := s.back.Get(x, y)

			// Update style if changed
			if !styleSet || !cell.Style.Equals(lastStyle) {
				s.output.WriteString(cell.Style.Sequence())
				lastStyle = cell.Style
				styleSet = true
			}

			s.output.WriteRune(cell.Rune)
		}

		// Don't add newline on last row
		if y < s.height-1 {
			s.output.WriteString("\r\n")
		}
	}

	// Reset style at end
	s.output.WriteString(terminal.StyleReset)

	// Write to terminal
	os.Stdout.WriteString(s.output.String())

	// Copy back to front
	for y := 0; y < s.height; y++ {
		copy(s.front.cells[y], s.back.cells[y])
	}
}

// Flush ensures all output is written
func (s *Screen) Flush() {
	os.Stdout.Sync()
}

// SetCell sets a cell in the back buffer
func (s *Screen) SetCell(x, y int, cell Cell) {
	s.back.Set(x, y, cell)
}

// DrawString draws a string in the back buffer
func (s *Screen) DrawString(x, y int, str string, style terminal.Style) {
	s.back.DrawString(x, y, str, style)
}

// DrawBox draws a box in the back buffer
func (s *Screen) DrawBox(x, y, width, height int, style terminal.Style) {
	s.back.DrawBox(x, y, width, height, style)
}

// Fill fills the back buffer with a cell
func (s *Screen) Fill(cell Cell) {
	s.back.Fill(cell)
}

// ShowCursor moves the cursor to the specified position and shows it
func (s *Screen) ShowCursor(x, y int) {
	fmt.Print(terminal.CursorMove(x+1, y+1))
	fmt.Print(terminal.CursorShow)
}

// HideCursor hides the cursor
func (s *Screen) HideCursor() {
	fmt.Print(terminal.CursorHide)
}

