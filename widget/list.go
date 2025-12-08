package widget

import (
	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
	"github.com/agiles231/gotui/terminal"
)

// ListItem represents an item in a list
type ListItem struct {
	Text  string
	Value interface{}
}

// List is a scrollable list widget
type List struct {
	BaseWidget
	items         []ListItem
	selected      int
	offset        int // Scroll offset
	style         terminal.Style
	selectedStyle terminal.Style
	height        int
	showBorder    bool
	onSelect      func(index int, item ListItem)
	onChange      func(index int, item ListItem)
}

// NewList creates a new list widget
func NewList() *List {
	l := &List{
		BaseWidget:    NewBaseWidget(),
		style:         terminal.DefaultStyle(),
		selectedStyle: terminal.DefaultStyle().WithReverse(),
		height:        10,
	}
	l.SetInteractive(true)
	return l
}

// SetItems sets the list items
func (l *List) SetItems(items []ListItem) *List {
	l.items = items
	if l.selected >= len(items) {
		l.selected = len(items) - 1
	}
	if l.selected < 0 {
		l.selected = 0
	}
	l.ensureVisible()
	return l
}

// SetStrings sets items from a string slice
func (l *List) SetStrings(strings []string) *List {
	items := make([]ListItem, len(strings))
	for i, s := range strings {
		items[i] = ListItem{Text: s, Value: s}
	}
	return l.SetItems(items)
}

// Items returns the list items
func (l *List) Items() []ListItem {
	return l.items
}

// Selected returns the selected index
func (l *List) Selected() int {
	return l.selected
}

// SelectedItem returns the selected item
func (l *List) SelectedItem() *ListItem {
	if l.selected >= 0 && l.selected < len(l.items) {
		return &l.items[l.selected]
	}
	return nil
}

// Select sets the selected index
func (l *List) Select(index int) *List {
	if index < 0 {
		index = 0
	}
	if index >= len(l.items) {
		index = len(l.items) - 1
	}
	l.selected = index
	l.ensureVisible()
	return l
}

// SetHeight sets the visible height
func (l *List) SetHeight(height int) *List {
	l.height = height
	return l
}

// SetStyle sets the normal style
func (l *List) SetStyle(style terminal.Style) *List {
	l.style = style
	return l
}

// SetSelectedStyle sets the selected item style
func (l *List) SetSelectedStyle(style terminal.Style) *List {
	l.selectedStyle = style
	return l
}

// SetShowBorder enables or disables the border
func (l *List) SetShowBorder(show bool) *List {
	l.showBorder = show
	return l
}

// OnSelect sets the callback for Enter key
func (l *List) OnSelect(fn func(index int, item ListItem)) *List {
	l.onSelect = fn
	return l
}

// OnChange sets the callback for selection change
func (l *List) OnChange(fn func(index int, item ListItem)) *List {
	l.onChange = fn
	return l
}

// Render draws the list
func (l *List) Render(buf *screen.Buffer, bounds layout.Rect) {
	if !l.visible {
		return
	}

	innerBounds := bounds
	if l.showBorder {
		buf.DrawBox(bounds.X, bounds.Y, bounds.Width, bounds.Height, l.style)
		innerBounds = bounds.Inset(1, 1, 1, 1)
	}

	visibleHeight := innerBounds.Height
	if visibleHeight > l.height {
		visibleHeight = l.height
	}

	for i := 0; i < visibleHeight; i++ {
		itemIndex := l.offset + i
		if itemIndex >= len(l.items) {
			break
		}

		item := l.items[itemIndex]
		style := l.style
		if itemIndex == l.selected && l.focused {
			style = l.selectedStyle
		}

		// Clear line
		for x := 0; x < innerBounds.Width; x++ {
			buf.Set(innerBounds.X+x, innerBounds.Y+i, screen.NewCell(' ', style))
		}

		// Draw item text
		buf.DrawStringClipped(innerBounds.X, innerBounds.Y+i, item.Text, style, innerBounds.Width)
	}

	// Draw scrollbar if needed
	if len(l.items) > visibleHeight {
		l.drawScrollbar(buf, innerBounds, visibleHeight)
	}
}

// drawScrollbar draws a scrollbar
func (l *List) drawScrollbar(buf *screen.Buffer, bounds layout.Rect, visibleHeight int) {
	if visibleHeight <= 0 || len(l.items) <= visibleHeight {
		return
	}

	scrollX := bounds.X + bounds.Width - 1
	
	// Calculate thumb position and size
	thumbSize := max(1, (visibleHeight*visibleHeight)/len(l.items))
	thumbPos := (l.offset * (visibleHeight - thumbSize)) / (len(l.items) - visibleHeight)

	scrollStyle := terminal.DefaultStyle().WithDim()
	thumbStyle := terminal.DefaultStyle().WithReverse()

	for y := 0; y < visibleHeight; y++ {
		style := scrollStyle
		char := '│'
		if y >= thumbPos && y < thumbPos+thumbSize {
			style = thumbStyle
			char = '█'
		}
		buf.Set(scrollX, bounds.Y+y, screen.NewCell(char, style))
	}
}

// HandleEvent handles input events
func (l *List) HandleEvent(event input.Event) bool {
	if !l.focused {
		return false
	}

	keyEvent, ok := event.(input.KeyEvent)
	if !ok {
		return false
	}

	switch keyEvent.Key {
	case input.KeyUp:
		l.moveUp()
		return true
	case input.KeyDown:
		l.moveDown()
		return true
	case input.KeyPageUp:
		l.pageUp()
		return true
	case input.KeyPageDown:
		l.pageDown()
		return true
	case input.KeyHome:
		l.Select(0)
		l.notifyChange()
		return true
	case input.KeyEnd:
		l.Select(len(l.items) - 1)
		l.notifyChange()
		return true
	case input.KeyEnter:
		if l.onSelect != nil && l.selected < len(l.items) {
			l.onSelect(l.selected, l.items[l.selected])
		}
		return true
	}

	return false
}

func (l *List) moveUp() {
	if l.selected > 0 {
		l.selected--
		l.ensureVisible()
		l.notifyChange()
	}
}

func (l *List) moveDown() {
	if l.selected < len(l.items)-1 {
		l.selected++
		l.ensureVisible()
		l.notifyChange()
	}
}

func (l *List) pageUp() {
	l.selected -= l.height
	if l.selected < 0 {
		l.selected = 0
	}
	l.ensureVisible()
	l.notifyChange()
}

func (l *List) pageDown() {
	l.selected += l.height
	if l.selected >= len(l.items) {
		l.selected = len(l.items) - 1
	}
	l.ensureVisible()
	l.notifyChange()
}

func (l *List) ensureVisible() {
	if l.selected < l.offset {
		l.offset = l.selected
	}
	if l.selected >= l.offset+l.height {
		l.offset = l.selected - l.height + 1
	}
}

func (l *List) notifyChange() {
	if l.onChange != nil && l.selected < len(l.items) {
		l.onChange(l.selected, l.items[l.selected])
	}
}

// Size returns the preferred size
func (l *List) Size() layout.Size {
	width := 0
	for _, item := range l.items {
		if len(item.Text) > width {
			width = len(item.Text)
		}
	}
	height := len(l.items)
	if height > l.height {
		height = l.height
	}
	if l.showBorder {
		width += 2
		height += 2
	}
	return layout.NewSize(width, height)
}

// MinSize returns the minimum size
func (l *List) MinSize() layout.Size {
	if l.showBorder {
		return layout.NewSize(5, 4)
	}
	return layout.NewSize(3, 1)
}

