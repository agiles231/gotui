package screen

import "github.com/agiles231/gotui/terminal"

// Buffer represents a 3D grid of cells with z-ordering for layered rendering
type Buffer struct {
	cells  [][][]Cell // indexed as cells[z][y][x]
	width  int
	height int
	depth  int
}

// NewBuffer creates a new buffer with the specified dimensions
func NewBuffer(width, height, depth int) *Buffer {
	b := &Buffer{
		width:  width,
		height: height,
		depth:  depth,
	}
	b.cells = make([][][]Cell, depth)
	for z := range b.cells {
		b.cells[z] = make([][]Cell, height)
		for y := range b.cells[z] {
			b.cells[z][y] = make([]Cell, width)
			for x := range b.cells[z][y] {
				b.cells[z][y][x] = EmptyCell()
			}
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

// Depth returns the buffer depth (number of z-layers)
func (b *Buffer) Depth() int {
	return b.depth
}

// Get returns the cell at the given position
func (b *Buffer) Get(x, y, z int) Cell {
	if x < 0 || x >= b.width || y < 0 || y >= b.height || z < 0 || z >= b.depth {
		return EmptyCell()
	}
	return b.cells[z][y][x]
}

// Set sets the cell at the given position
func (b *Buffer) Set(x, y, z int, cell Cell) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height || z < 0 || z >= b.depth {
		return
	}
	b.cells[z][y][x] = cell
}

// SetRune sets just the rune at the given position
func (b *Buffer) SetRune(x, y, z int, r rune) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height || z < 0 || z >= b.depth {
		return
	}
	b.cells[z][y][x].Rune = r
}

// SetStyle sets just the style at the given position
func (b *Buffer) SetStyle(x, y, z int, style terminal.Style) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height || z < 0 || z >= b.depth {
		return
	}
	b.cells[z][y][x].Style = style
}

// Fill fills the entire buffer with the given cell (all z-layers)
func (b *Buffer) Fill(cell Cell) {
	for z := range b.cells {
		for y := range b.cells[z] {
			for x := range b.cells[z][y] {
				b.cells[z][y][x] = cell
			}
		}
	}
}

// Clear clears the buffer with empty cells (all z-layers)
func (b *Buffer) Clear() {
	b.Fill(EmptyCell())
}

// FillRect fills a rectangular region with the given cell at the specified z-layer
func (b *Buffer) FillRect(x, y, z, width, height int, cell Cell) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			b.Set(x+dx, y+dy, z, cell)
		}
	}
}

// DrawString draws a string at the given position with the given style
func (b *Buffer) DrawString(x, y, z int, s string, style terminal.Style) {
	// TODO: handle multi rune characters
	for i, r := range s {
		b.Set(x+i, y, z, NewCell(r, style))
	}
}

// DrawStringClipped draws a string clipped to a maximum width
func (b *Buffer) DrawStringClipped(x, y, z int, s string, style terminal.Style, maxWidth int) {
	i := 0
	for _, r := range s {
		if i >= maxWidth {
			break
		}
		b.Set(x+i, y, z, NewCell(r, style))
		i++
	}
}

// DrawHLine draws a horizontal line
func (b *Buffer) DrawHLine(x, y, z, width int, r rune, style terminal.Style) {
	cell := NewCell(r, style)
	for i := 0; i < width; i++ {
		b.Set(x+i, y, z, cell)
	}
}

// DrawVLine draws a vertical line
func (b *Buffer) DrawVLine(x, y, z, height int, r rune, style terminal.Style) {
	cell := NewCell(r, style)
	for i := 0; i < height; i++ {
		b.Set(x, y+i, z, cell)
	}
}

// DrawBox draws a box with the specified box-drawing characters
func (b *Buffer) DrawBox(x, y, z, width, height int, style terminal.Style) {
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
	b.Set(x, y, z, NewCell(topLeft, style))
	b.DrawHLine(x+1, y, z, width-2, horizontal, style)
	b.Set(x+width-1, y, z, NewCell(topRight, style))

	// Side borders
	b.DrawVLine(x, y+1, z, height-2, vertical, style)
	b.DrawVLine(x+width-1, y+1, z, height-2, vertical, style)

	// Bottom border
	b.Set(x, y+height-1, z, NewCell(bottomLeft, style))
	b.DrawHLine(x+1, y+height-1, z, width-2, horizontal, style)
	b.Set(x+width-1, y+height-1, z, NewCell(bottomRight, style))
}

// DrawDoubleBox draws a box with double-line characters
func (b *Buffer) DrawDoubleBox(x, y, z, width, height int, style terminal.Style) {
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
	b.Set(x, y, z, NewCell(topLeft, style))
	b.DrawHLine(x+1, y, z, width-2, horizontal, style)
	b.Set(x+width-1, y, z, NewCell(topRight, style))

	// Side borders
	b.DrawVLine(x, y+1, z, height-2, vertical, style)
	b.DrawVLine(x+width-1, y+1, z, height-2, vertical, style)

	// Bottom border
	b.Set(x, y+height-1, z, NewCell(bottomLeft, style))
	b.DrawHLine(x+1, y+height-1, z, width-2, horizontal, style)
	b.Set(x+width-1, y+height-1, z, NewCell(bottomRight, style))
}

// Flatten composites all z-layers into a 2D slice for rendering
// Higher z-values overwrite lower z-values, empty cells are skipped
func (b *Buffer) Flatten() [][]Cell {
	result := make([][]Cell, b.height)
	for y := 0; y < b.height; y++ {
		result[y] = make([]Cell, b.width)
		for x := 0; x < b.width; x++ {
			// Start with empty cell
			result[y][x] = EmptyCell()
			// Composite from z=0 to z=depth-1
			for z := 0; z < b.depth; z++ {
				cell := b.cells[z][y][x]
				if !cell.IsEmpty() {
					result[y][x] = cell
				}
			}
		}
	}
	return result
}

// Resize creates a new buffer with the given dimensions, copying existing content
func (b *Buffer) Resize(width, height, depth int) *Buffer {
	newBuf := NewBuffer(width, height, depth)

	// Copy existing content
	for z := 0; z < min(b.depth, depth); z++ {
		for y := 0; y < min(b.height, height); y++ {
			for x := 0; x < min(b.width, width); x++ {
				newBuf.cells[z][y][x] = b.cells[z][y][x]
			}
		}
	}

	return newBuf
}

// Clone creates a deep copy of the buffer
func (b *Buffer) Clone() *Buffer {
	clone := NewBuffer(b.width, b.height, b.depth)
	for z := range b.cells {
		for y := range b.cells[z] {
			copy(clone.cells[z][y], b.cells[z][y])
		}
	}
	return clone
}

// Blit copies a region from another buffer onto this buffer at a specific z-layer
func (b *Buffer) Blit(src *Buffer, srcX, srcY, srcZ, dstX, dstY, dstZ, width, height int) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			cell := src.Get(srcX+dx, srcY+dy, srcZ)
			b.Set(dstX+dx, dstY+dy, dstZ, cell)
		}
	}
}

// BlitBuffer copies the entire source buffer onto this buffer at the given position
func (b *Buffer) BlitBuffer(src *Buffer, dstX, dstY, dstZ int) {
	for z := 0; z < src.depth; z++ {
		targetZ := dstZ + z
		if targetZ >= b.depth {
			break
		}
		for dy := 0; dy < src.height; dy++ {
			for dx := 0; dx < src.width; dx++ {
				cell := src.Get(dx, dy, z)
				b.Set(dstX+dx, dstY+dy, targetZ, cell)
			}
		}
	}
}
