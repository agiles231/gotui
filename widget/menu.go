package widget

import (
	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
	"github.com/agiles231/gotui/terminal"
)

// MenuItem represents an item in a menu
type MenuItem struct {
	Label    string
	Shortcut string
	Action   func()
	Disabled bool
	Children []*MenuItem // For submenus
}

// Menu is a menu widget
type Menu struct {
	BaseWidget
	items         []*MenuItem
	selected      int
	style         terminal.Style
	selectedStyle terminal.Style
	disabledStyle terminal.Style
	showBorder    bool
	width         int
	onSelect      func(index int, item *MenuItem)
}

// NewMenu creates a new menu widget
func NewMenu() *Menu {
	m := &Menu{
		BaseWidget:    NewBaseWidget(),
		style:         terminal.DefaultStyle(),
		selectedStyle: terminal.DefaultStyle().WithReverse(),
		disabledStyle: terminal.DefaultStyle().WithDim(),
		showBorder:    true,
	}
	m.SetInteractive(true)
	return m
}

// SetItems sets the menu items
func (m *Menu) SetItems(items []*MenuItem) *Menu {
	m.items = items
	if m.selected >= len(items) {
		m.selected = len(items) - 1
	}
	if m.selected < 0 {
		m.selected = 0
	}
	m.skipDisabled(1)
	return m
}

// Items returns the menu items
func (m *Menu) Items() []*MenuItem {
	return m.items
}

// Selected returns the selected index
func (m *Menu) Selected() int {
	return m.selected
}

// SelectedItem returns the selected item
func (m *Menu) SelectedItem() *MenuItem {
	if m.selected >= 0 && m.selected < len(m.items) {
		return m.items[m.selected]
	}
	return nil
}

// Select sets the selected index
func (m *Menu) Select(index int) *Menu {
	if index < 0 {
		index = 0
	}
	if index >= len(m.items) {
		index = len(m.items) - 1
	}
	m.selected = index
	return m
}

// SetWidth sets the menu width
func (m *Menu) SetWidth(width int) *Menu {
	m.width = width
	return m
}

// SetShowBorder enables or disables the border
func (m *Menu) SetShowBorder(show bool) *Menu {
	m.showBorder = show
	return m
}

// SetStyle sets the normal style
func (m *Menu) SetStyle(style terminal.Style) *Menu {
	m.style = style
	return m
}

// SetSelectedStyle sets the selected item style
func (m *Menu) SetSelectedStyle(style terminal.Style) *Menu {
	m.selectedStyle = style
	return m
}

// SetDisabledStyle sets the disabled item style
func (m *Menu) SetDisabledStyle(style terminal.Style) *Menu {
	m.disabledStyle = style
	return m
}

// OnSelect sets the callback for item selection
func (m *Menu) OnSelect(fn func(index int, item *MenuItem)) *Menu {
	m.onSelect = fn
	return m
}

// Render draws the menu
func (m *Menu) Render(buf *screen.Buffer, bounds layout.Rect) {
	if !m.visible {
		return
	}

	width := m.calculateWidth()
	if m.width > 0 {
		width = m.width
	}
	if width > bounds.Width {
		width = bounds.Width
	}

	height := len(m.items)
	if m.showBorder {
		height += 2
		width += 2
	}

	innerBounds := bounds
	if m.showBorder {
		buf.DrawBox(bounds.X, bounds.Y, width, height, m.style)
		innerBounds = layout.NewRect(bounds.X+1, bounds.Y+1, width-2, len(m.items))
	}

	for i, item := range m.items {
		if i >= innerBounds.Height {
			break
		}

		style := m.style
		if item.Disabled {
			style = m.disabledStyle
		} else if i == m.selected && m.focused {
			style = m.selectedStyle
		}

		// Clear line
		for x := 0; x < innerBounds.Width; x++ {
			buf.Set(innerBounds.X+x, innerBounds.Y+i, screen.NewCell(' ', style))
		}

		// Draw label
		label := item.Label
		if len(label) > innerBounds.Width {
			label = label[:innerBounds.Width]
		}
		buf.DrawString(innerBounds.X, innerBounds.Y+i, label, style)

		// Draw shortcut if present
		if item.Shortcut != "" && innerBounds.Width > len(label)+len(item.Shortcut)+2 {
			shortcutX := innerBounds.X + innerBounds.Width - len(item.Shortcut)
			buf.DrawString(shortcutX, innerBounds.Y+i, item.Shortcut, style.WithDim())
		}

		// Draw submenu indicator
		if len(item.Children) > 0 {
			buf.Set(innerBounds.X+innerBounds.Width-1, innerBounds.Y+i, screen.NewCell('▶', style))
		}
	}
}

func (m *Menu) calculateWidth() int {
	width := 0
	for _, item := range m.items {
		itemWidth := len(item.Label)
		if item.Shortcut != "" {
			itemWidth += 2 + len(item.Shortcut)
		}
		if len(item.Children) > 0 {
			itemWidth += 2
		}
		if itemWidth > width {
			width = itemWidth
		}
	}
	return width
}

// HandleEvent handles input events
func (m *Menu) HandleEvent(event input.Event) bool {
	if !m.focused {
		return false
	}

	keyEvent, ok := event.(input.KeyEvent)
	if !ok {
		return false
	}

	switch keyEvent.Key {
	case input.KeyUp:
		m.moveUp()
		return true
	case input.KeyDown:
		m.moveDown()
		return true
	case input.KeyEnter:
		m.activate()
		return true
	case input.KeyEscape:
		return true
	}

	return false
}

func (m *Menu) moveUp() {
	m.selected--
	if m.selected < 0 {
		m.selected = len(m.items) - 1
	}
	m.skipDisabled(-1)
}

func (m *Menu) moveDown() {
	m.selected++
	if m.selected >= len(m.items) {
		m.selected = 0
	}
	m.skipDisabled(1)
}

func (m *Menu) skipDisabled(direction int) {
	// Skip disabled items
	for i := 0; i < len(m.items); i++ {
		if m.selected >= 0 && m.selected < len(m.items) && !m.items[m.selected].Disabled {
			return
		}
		m.selected += direction
		if m.selected < 0 {
			m.selected = len(m.items) - 1
		}
		if m.selected >= len(m.items) {
			m.selected = 0
		}
	}
}

func (m *Menu) activate() {
	if m.selected < 0 || m.selected >= len(m.items) {
		return
	}

	item := m.items[m.selected]
	if item.Disabled {
		return
	}

	if item.Action != nil {
		item.Action()
	}

	if m.onSelect != nil {
		m.onSelect(m.selected, item)
	}
}

// Size returns the preferred size
func (m *Menu) Size() layout.Size {
	width := m.calculateWidth()
	if m.width > 0 {
		width = m.width
	}
	height := len(m.items)
	if m.showBorder {
		width += 2
		height += 2
	}
	return layout.NewSize(width, height)
}

// MinSize returns the minimum size
func (m *Menu) MinSize() layout.Size {
	if m.showBorder {
		return layout.NewSize(5, 4)
	}
	return layout.NewSize(3, 1)
}

// Separator returns a menu separator item
func Separator() *MenuItem {
	return &MenuItem{
		Label:    "────────────────",
		Disabled: true,
	}
}

