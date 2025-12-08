package screen

import "github.com/agiles231/gotui/terminal"

// Buffer represents a 2D grid of cells
type Buffer struct {
	cells  [][]Cell
	width  int
	height int
}

// NewBuffer creates a new buffer with the specified dimensions
func NewBuffer(width, height int) *Buffer {
	b := &Buffer{
		width:  width,
		height: height,
	}
	b.cells = make([][]Cell, height)
	for y := range b.cells {
		b.cells[y] = make([]Cell, width)
		for x := range b.cells[y] {
			b.cells[y][x] = EmptyCell()
		}
	}
	return b
}

// Width returns the buffer width
func (b *Buffer) Width() int {
	return b.width
}

// Height returns the buffer height
func (b *Buffer) Height() int {
	return b.height
}

// Get returns the cell at the given position
func (b *Buffer) Get(x, y int) Cell {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return EmptyCell()
	}
	return b.cells[y][x]
}

// Set sets the cell at the given position
func (b *Buffer) Set(x, y int, cell Cell) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return
	}
	b.cells[y][x] = cell
}

// SetRune sets just the rune at the given position
func (b *Buffer) SetRune(x, y int, r rune) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return
	}
	b.cells[y][x].Rune = r
}

// SetStyle sets just the style at the given position
func (b *Buffer) SetStyle(x, y int, style terminal.Style) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return
	}
	b.cells[y][x].Style = style
}

// Fill fills the entire buffer with the given cell
func (b *Buffer) Fill(cell Cell) {
	for y := range b.cells {
		for x := range b.cells[y] {
			b.cells[y][x] = cell
		}
	}
}

// Clear clears the buffer with empty cells
func (b *Buffer) Clear() {
	b.Fill(EmptyCell())
}

// FillRect fills a rectangular region with the given cell
func (b *Buffer) FillRect(x, y, width, height int, cell Cell) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			b.Set(x+dx, y+dy, cell)
		}
	}
}

// DrawString draws a string at the given position with the given style
func (b *Buffer) DrawString(x, y int, s string, style terminal.Style) {
	for i, r := range s {
		b.Set(x+i, y, NewCell(r, style))
	}
}

// DrawStringClipped draws a string clipped to a maximum width
func (b *Buffer) DrawStringClipped(x, y int, s string, style terminal.Style, maxWidth int) {
	i := 0
	for _, r := range s {
		if i >= maxWidth {
			break
		}
		b.Set(x+i, y, NewCell(r, style))
		i++
	}
}

// DrawHLine draws a horizontal line
func (b *Buffer) DrawHLine(x, y, width int, r rune, style terminal.Style) {
	cell := NewCell(r, style)
	for i := 0; i < width; i++ {
		b.Set(x+i, y, cell)
	}
}

// DrawVLine draws a vertical line
func (b *Buffer) DrawVLine(x, y, height int, r rune, style terminal.Style) {
	cell := NewCell(r, style)
	for i := 0; i < height; i++ {
		b.Set(x, y+i, cell)
	}
}

// DrawBox draws a box with the specified box-drawing characters
func (b *Buffer) DrawBox(x, y, width, height int, style terminal.Style) {
	if width < 2 || height < 2 {
		return
	}

	// Box drawing characters
	topLeft := '┌'
	topRight := '┐'
	bottomLeft := '└'
	bottomRight := '┘'
	horizontal := '─'
	vertical := '│'

	// Top border
	b.Set(x, y, NewCell(topLeft, style))
	b.DrawHLine(x+1, y, width-2, horizontal, style)
	b.Set(x+width-1, y, NewCell(topRight, style))

	// Side borders
	b.DrawVLine(x, y+1, height-2, vertical, style)
	b.DrawVLine(x+width-1, y+1, height-2, vertical, style)

	// Bottom border
	b.Set(x, y+height-1, NewCell(bottomLeft, style))
	b.DrawHLine(x+1, y+height-1, width-2, horizontal, style)
	b.Set(x+width-1, y+height-1, NewCell(bottomRight, style))
}

// DrawDoubleBox draws a box with double-line characters
func (b *Buffer) DrawDoubleBox(x, y, width, height int, style terminal.Style) {
	if width < 2 || height < 2 {
		return
	}

	// Double box drawing characters
	topLeft := '╔'
	topRight := '╗'
	bottomLeft := '╚'
	bottomRight := '╝'
	horizontal := '═'
	vertical := '║'

	// Top border
	b.Set(x, y, NewCell(topLeft, style))
	b.DrawHLine(x+1, y, width-2, horizontal, style)
	b.Set(x+width-1, y, NewCell(topRight, style))

	// Side borders
	b.DrawVLine(x, y+1, height-2, vertical, style)
	b.DrawVLine(x+width-1, y+1, height-2, vertical, style)

	// Bottom border
	b.Set(x, y+height-1, NewCell(bottomLeft, style))
	b.DrawHLine(x+1, y+height-1, width-2, horizontal, style)
	b.Set(x+width-1, y+height-1, NewCell(bottomRight, style))
}

// Resize creates a new buffer with the given dimensions, copying existing content
func (b *Buffer) Resize(width, height int) *Buffer {
	newBuf := NewBuffer(width, height)

	// Copy existing content
	for y := 0; y < min(b.height, height); y++ {
		for x := 0; x < min(b.width, width); x++ {
			newBuf.cells[y][x] = b.cells[y][x]
		}
	}

	return newBuf
}

// Clone creates a deep copy of the buffer
func (b *Buffer) Clone() *Buffer {
	clone := NewBuffer(b.width, b.height)
	for y := range b.cells {
		copy(clone.cells[y], b.cells[y])
	}
	return clone
}

// Blit copies a region from another buffer onto this buffer
func (b *Buffer) Blit(src *Buffer, srcX, srcY, dstX, dstY, width, height int) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			cell := src.Get(srcX+dx, srcY+dy)
			b.Set(dstX+dx, dstY+dy, cell)
		}
	}
}

// BlitBuffer copies the entire source buffer onto this buffer at the given position
func (b *Buffer) BlitBuffer(src *Buffer, dstX, dstY int) {
	b.Blit(src, 0, 0, dstX, dstY, src.width, src.height)
}

