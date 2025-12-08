package widget

import (
	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
	"github.com/agiles231/gotui/terminal"
)

// TableColumn defines a column in the table
type TableColumn struct {
	Title string
	Width int
	Flex  int // Flex factor for auto-sizing
	Align layout.Alignment
}

// Table is a table widget with columns and rows
type Table struct {
	BaseWidget
	columns        []TableColumn
	rows           [][]string
	selectedRow    int
	offset         int
	height         int
	showHeader     bool
	showBorder     bool
	style          terminal.Style
	headerStyle    terminal.Style
	selectedStyle  terminal.Style
	onSelect       func(row int)
	onChange       func(row int)
}

// NewTable creates a new table widget
func NewTable() *Table {
	t := &Table{
		BaseWidget:    NewBaseWidget(),
		showHeader:    true,
		height:        10,
		style:         terminal.DefaultStyle(),
		headerStyle:   terminal.DefaultStyle().WithBold(),
		selectedStyle: terminal.DefaultStyle().WithReverse(),
	}
	t.SetInteractive(true)
	return t
}

// SetColumns sets the table columns
func (t *Table) SetColumns(columns []TableColumn) *Table {
	t.columns = columns
	return t
}

// SetRows sets the table rows
func (t *Table) SetRows(rows [][]string) *Table {
	t.rows = rows
	if t.selectedRow >= len(rows) {
		t.selectedRow = len(rows) - 1
	}
	if t.selectedRow < 0 {
		t.selectedRow = 0
	}
	t.ensureVisible()
	return t
}

// Rows returns the table rows
func (t *Table) Rows() [][]string {
	return t.rows
}

// SelectedRow returns the selected row index
func (t *Table) SelectedRow() int {
	return t.selectedRow
}

// SelectRow sets the selected row
func (t *Table) SelectRow(row int) *Table {
	if row < 0 {
		row = 0
	}
	if row >= len(t.rows) {
		row = len(t.rows) - 1
	}
	t.selectedRow = row
	t.ensureVisible()
	return t
}

// SetHeight sets the visible height
func (t *Table) SetHeight(height int) *Table {
	t.height = height
	return t
}

// SetShowHeader enables or disables the header row
func (t *Table) SetShowHeader(show bool) *Table {
	t.showHeader = show
	return t
}

// SetShowBorder enables or disables the border
func (t *Table) SetShowBorder(show bool) *Table {
	t.showBorder = show
	return t
}

// SetStyle sets the normal style
func (t *Table) SetStyle(style terminal.Style) *Table {
	t.style = style
	return t
}

// SetHeaderStyle sets the header style
func (t *Table) SetHeaderStyle(style terminal.Style) *Table {
	t.headerStyle = style
	return t
}

// SetSelectedStyle sets the selected row style
func (t *Table) SetSelectedStyle(style terminal.Style) *Table {
	t.selectedStyle = style
	return t
}

// OnSelect sets the callback for Enter key
func (t *Table) OnSelect(fn func(row int)) *Table {
	t.onSelect = fn
	return t
}

// OnChange sets the callback for selection change
func (t *Table) OnChange(fn func(row int)) *Table {
	t.onChange = fn
	return t
}

// Render draws the table
func (t *Table) Render(buf *screen.Buffer, bounds layout.Rect) {
	if !t.visible || len(t.columns) == 0 {
		return
	}

	innerBounds := bounds
	if t.showBorder {
		buf.DrawBox(bounds.X, bounds.Y, bounds.Width, bounds.Height, t.style)
		innerBounds = bounds.Inset(1, 1, 1, 1)
	}

	// Calculate column widths
	colWidths := t.calculateColumnWidths(innerBounds.Width)

	y := innerBounds.Y

	// Draw header
	if t.showHeader {
		t.drawRow(buf, innerBounds.X, y, colWidths, t.getColumnTitles(), t.headerStyle)
		y++
		// Draw separator
		t.drawSeparator(buf, innerBounds.X, y, colWidths)
		y++
	}

	// Draw rows
	visibleHeight := innerBounds.Height
	if t.showHeader {
		visibleHeight -= 2
	}

	for i := 0; i < visibleHeight && i+t.offset < len(t.rows); i++ {
		rowIndex := i + t.offset
		style := t.style
		if rowIndex == t.selectedRow && t.focused {
			style = t.selectedStyle
		}
		t.drawRow(buf, innerBounds.X, y+i, colWidths, t.rows[rowIndex], style)
	}
}

func (t *Table) getColumnTitles() []string {
	titles := make([]string, len(t.columns))
	for i, col := range t.columns {
		titles[i] = col.Title
	}
	return titles
}

func (t *Table) calculateColumnWidths(totalWidth int) []int {
	widths := make([]int, len(t.columns))
	flexTotal := 0
	fixedWidth := 0

	// First pass: calculate fixed widths and flex total
	for i, col := range t.columns {
		if col.Width > 0 {
			widths[i] = col.Width
			fixedWidth += col.Width
		} else {
			flex := col.Flex
			if flex == 0 {
				flex = 1
			}
			flexTotal += flex
		}
	}

	// Add separators
	separatorWidth := len(t.columns) - 1
	remaining := totalWidth - fixedWidth - separatorWidth

	// Second pass: distribute remaining width
	if flexTotal > 0 && remaining > 0 {
		for i, col := range t.columns {
			if col.Width == 0 {
				flex := col.Flex
				if flex == 0 {
					flex = 1
				}
				widths[i] = (remaining * flex) / flexTotal
			}
		}
	}

	return widths
}

func (t *Table) drawRow(buf *screen.Buffer, x, y int, widths []int, cells []string, style terminal.Style) {
	currentX := x
	for i, width := range widths {
		// Clear cell
		for dx := 0; dx < width; dx++ {
			buf.Set(currentX+dx, y, screen.NewCell(' ', style))
		}

		// Draw cell content
		if i < len(cells) {
			text := cells[i]
			if len(text) > width {
				text = text[:width]
			}

			// Apply alignment
			offset := 0
			if i < len(t.columns) {
				offset = layout.Align(len(text), width, t.columns[i].Align)
			}

			buf.DrawString(currentX+offset, y, text, style)
		}

		currentX += width

		// Draw separator
		if i < len(widths)-1 {
			buf.Set(currentX, y, screen.NewCell('│', style))
			currentX++
		}
	}
}

func (t *Table) drawSeparator(buf *screen.Buffer, x, y int, widths []int) {
	currentX := x
	for i, width := range widths {
		for dx := 0; dx < width; dx++ {
			buf.Set(currentX+dx, y, screen.NewCell('─', t.style))
		}
		currentX += width

		if i < len(widths)-1 {
			buf.Set(currentX, y, screen.NewCell('┼', t.style))
			currentX++
		}
	}
}

// HandleEvent handles input events
func (t *Table) HandleEvent(event input.Event) bool {
	if !t.focused {
		return false
	}

	keyEvent, ok := event.(input.KeyEvent)
	if !ok {
		return false
	}

	switch keyEvent.Key {
	case input.KeyUp:
		t.moveUp()
		return true
	case input.KeyDown:
		t.moveDown()
		return true
	case input.KeyPageUp:
		t.pageUp()
		return true
	case input.KeyPageDown:
		t.pageDown()
		return true
	case input.KeyHome:
		t.SelectRow(0)
		t.notifyChange()
		return true
	case input.KeyEnd:
		t.SelectRow(len(t.rows) - 1)
		t.notifyChange()
		return true
	case input.KeyEnter:
		if t.onSelect != nil {
			t.onSelect(t.selectedRow)
		}
		return true
	}

	return false
}

func (t *Table) moveUp() {
	if t.selectedRow > 0 {
		t.selectedRow--
		t.ensureVisible()
		t.notifyChange()
	}
}

func (t *Table) moveDown() {
	if t.selectedRow < len(t.rows)-1 {
		t.selectedRow++
		t.ensureVisible()
		t.notifyChange()
	}
}

func (t *Table) pageUp() {
	t.selectedRow -= t.height
	if t.selectedRow < 0 {
		t.selectedRow = 0
	}
	t.ensureVisible()
	t.notifyChange()
}

func (t *Table) pageDown() {
	t.selectedRow += t.height
	if t.selectedRow >= len(t.rows) {
		t.selectedRow = len(t.rows) - 1
	}
	t.ensureVisible()
	t.notifyChange()
}

func (t *Table) ensureVisible() {
	visibleRows := t.height
	if t.showHeader {
		visibleRows -= 2
	}
	if t.selectedRow < t.offset {
		t.offset = t.selectedRow
	}
	if t.selectedRow >= t.offset+visibleRows {
		t.offset = t.selectedRow - visibleRows + 1
	}
}

func (t *Table) notifyChange() {
	if t.onChange != nil {
		t.onChange(t.selectedRow)
	}
}

// Size returns the preferred size
func (t *Table) Size() layout.Size {
	width := 0
	for _, col := range t.columns {
		if col.Width > 0 {
			width += col.Width
		} else {
			width += 10 // Default width
		}
	}
	width += len(t.columns) - 1 // Separators

	height := len(t.rows)
	if t.showHeader {
		height += 2
	}
	if height > t.height {
		height = t.height
	}
	if t.showBorder {
		width += 2
		height += 2
	}

	return layout.NewSize(width, height)
}

// MinSize returns the minimum size
func (t *Table) MinSize() layout.Size {
	width := len(t.columns)*3 + len(t.columns) - 1
	if t.showBorder {
		width += 2
	}
	return layout.NewSize(width, 3)
}

