package layout

// Direction represents layout direction
type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

// FlexChild represents a child in a flex container
type FlexChild struct {
	// Fixed size (if > 0, takes precedence over Flex)
	Fixed int
	// Flex factor for proportional sizing
	Flex int
	// Minimum size
	Min int
	// Maximum size
	Max int
}

// NewFixedChild creates a child with fixed size
func NewFixedChild(size int) FlexChild {
	return FlexChild{Fixed: size}
}

// NewFlexChild creates a child with flex factor
func NewFlexChild(flex int) FlexChild {
	return FlexChild{Flex: flex}
}

// NewFlexChildWithBounds creates a child with flex and bounds
func NewFlexChildWithBounds(flex, minSize, maxSize int) FlexChild {
	return FlexChild{Flex: flex, Min: minSize, Max: maxSize}
}

// Flex performs flex layout on children within a container
type Flex struct {
	Direction Direction
	Gap       int // Gap between children
}

// NewFlex creates a new flex layout
func NewFlex(direction Direction) *Flex {
	return &Flex{Direction: direction}
}

// NewHFlex creates a horizontal flex layout
func NewHFlex() *Flex {
	return NewFlex(Horizontal)
}

// NewVFlex creates a vertical flex layout
func NewVFlex() *Flex {
	return NewFlex(Vertical)
}

// WithGap sets the gap between children
func (f *Flex) WithGap(gap int) *Flex {
	f.Gap = gap
	return f
}

// Layout calculates rectangles for each child within the container
func (f *Flex) Layout(container Rect, children []FlexChild) []Rect {
	if len(children) == 0 {
		return nil
	}

	// Calculate total available space
	var totalSpace int
	if f.Direction == Horizontal {
		totalSpace = container.Width
	} else {
		totalSpace = container.Height
	}

	// Subtract gaps
	totalGaps := f.Gap * (len(children) - 1)
	availableSpace := totalSpace - totalGaps

	// First pass: allocate fixed sizes and calculate flex total
	flexTotal := 0
	allocated := make([]int, len(children))
	remaining := availableSpace

	for i, child := range children {
		if child.Fixed > 0 {
			size := child.Fixed
			if child.Max > 0 && size > child.Max {
				size = child.Max
			}
			if size > remaining {
				size = remaining
			}
			allocated[i] = size
			remaining -= size
		} else {
			flexTotal += max(1, child.Flex)
		}
	}

	// Second pass: allocate flex space
	if flexTotal > 0 && remaining > 0 {
		for i, child := range children {
			if child.Fixed > 0 {
				continue
			}

			flex := child.Flex
			if flex == 0 {
				flex = 1
			}

			// Calculate proportional size
			size := (remaining * flex) / flexTotal

			// Apply bounds
			if child.Min > 0 && size < child.Min {
				size = child.Min
			}
			if child.Max > 0 && size > child.Max {
				size = child.Max
			}

			allocated[i] = size
		}
	}

	// Generate result rectangles
	results := make([]Rect, len(children))
	offset := 0

	for i, size := range allocated {
		if f.Direction == Horizontal {
			results[i] = NewRect(
				container.X+offset,
				container.Y,
				size,
				container.Height,
			)
		} else {
			results[i] = NewRect(
				container.X,
				container.Y+offset,
				container.Width,
				size,
			)
		}
		offset += size + f.Gap
	}

	return results
}

// Split evenly divides a container into n equal parts
func Split(container Rect, n int, direction Direction, gap int) []Rect {
	if n <= 0 {
		return nil
	}

	children := make([]FlexChild, n)
	for i := range children {
		children[i] = NewFlexChild(1)
	}

	flex := &Flex{Direction: direction, Gap: gap}
	return flex.Layout(container, children)
}

// SplitHorizontal splits a container into n horizontal parts
func SplitHorizontal(container Rect, n int) []Rect {
	return Split(container, n, Horizontal, 0)
}

// SplitVertical splits a container into n vertical parts
func SplitVertical(container Rect, n int) []Rect {
	return Split(container, n, Vertical, 0)
}

// Grid creates a grid layout
type Grid struct {
	Rows    int
	Cols    int
	RowGap  int
	ColGap  int
	Padding int
}

// NewGrid creates a new grid layout
func NewGrid(rows, cols int) *Grid {
	return &Grid{Rows: rows, Cols: cols}
}

// WithGaps sets row and column gaps
func (g *Grid) WithGaps(rowGap, colGap int) *Grid {
	g.RowGap = rowGap
	g.ColGap = colGap
	return g
}

// WithPadding sets internal padding
func (g *Grid) WithPadding(padding int) *Grid {
	g.Padding = padding
	return g
}

// Layout calculates cell rectangles for a grid
func (g *Grid) Layout(container Rect) [][]Rect {
	if g.Rows <= 0 || g.Cols <= 0 {
		return nil
	}

	// Apply padding
	inner := container.InsetAll(g.Padding)

	// Calculate cell sizes
	totalColGaps := g.ColGap * (g.Cols - 1)
	totalRowGaps := g.RowGap * (g.Rows - 1)

	cellWidth := (inner.Width - totalColGaps) / g.Cols
	cellHeight := (inner.Height - totalRowGaps) / g.Rows

	// Generate cells
	cells := make([][]Rect, g.Rows)
	for row := 0; row < g.Rows; row++ {
		cells[row] = make([]Rect, g.Cols)
		for col := 0; col < g.Cols; col++ {
			cells[row][col] = NewRect(
				inner.X+col*(cellWidth+g.ColGap),
				inner.Y+row*(cellHeight+g.RowGap),
				cellWidth,
				cellHeight,
			)
		}
	}

	return cells
}

// CellAt returns the rectangle for a specific cell
func (g *Grid) CellAt(container Rect, row, col int) Rect {
	cells := g.Layout(container)
	if row < 0 || row >= len(cells) || col < 0 || col >= len(cells[row]) {
		return Zero()
	}
	return cells[row][col]
}

