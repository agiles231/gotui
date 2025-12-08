package widget

import (
	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
)

// Widget is the interface all UI components implement
type Widget interface {
	// Render draws the widget to the buffer within the given bounds
	Render(buf *screen.Buffer, bounds layout.Rect)

	// HandleEvent processes an input event
	// Returns true if the event was handled
	HandleEvent(event input.Event) bool

	// Size returns the preferred size of the widget
	Size() layout.Size

	// MinSize returns the minimum size the widget needs
	MinSize() layout.Size

	// SetFocused sets whether the widget has focus
	SetFocused(focused bool)

	// IsFocused returns whether the widget is focused
	IsFocused() bool

	// IsInteractive returns whether the widget can receive input
	IsInteractive() bool
}

// BaseWidget provides common functionality for widgets
type BaseWidget struct {
	focused     bool
	interactive bool
	visible     bool
}

// NewBaseWidget creates a new base widget
func NewBaseWidget() BaseWidget {
	return BaseWidget{
		visible: true,
	}
}

// SetFocused sets the focus state
func (w *BaseWidget) SetFocused(focused bool) {
	w.focused = focused
}

// IsFocused returns the focus state
func (w *BaseWidget) IsFocused() bool {
	return w.focused
}

// IsInteractive returns whether the widget can be interacted with
func (w *BaseWidget) IsInteractive() bool {
	return w.interactive
}

// SetInteractive sets whether the widget can be interacted with
func (w *BaseWidget) SetInteractive(interactive bool) {
	w.interactive = interactive
}

// IsVisible returns whether the widget is visible
func (w *BaseWidget) IsVisible() bool {
	return w.visible
}

// SetVisible sets the visibility
func (w *BaseWidget) SetVisible(visible bool) {
	w.visible = visible
}

// Container is a widget that contains other widgets
type Container interface {
	Widget
	// Children returns the child widgets
	Children() []Widget
	// AddChild adds a child widget
	AddChild(w Widget)
	// RemoveChild removes a child widget
	RemoveChild(w Widget)
}

// FocusManager manages focus between widgets
type FocusManager struct {
	widgets       []Widget
	focusedIndex  int
	cycleFocus    bool
}

// NewFocusManager creates a new focus manager
func NewFocusManager() *FocusManager {
	return &FocusManager{
		focusedIndex: -1,
		cycleFocus:   true,
	}
}

// Add adds a widget to the focus manager
func (fm *FocusManager) Add(w Widget) {
	if w.IsInteractive() {
		fm.widgets = append(fm.widgets, w)
	}
}

// Remove removes a widget from the focus manager
func (fm *FocusManager) Remove(w Widget) {
	for i, widget := range fm.widgets {
		if widget == w {
			fm.widgets = append(fm.widgets[:i], fm.widgets[i+1:]...)
			if fm.focusedIndex >= len(fm.widgets) {
				fm.focusedIndex = len(fm.widgets) - 1
			}
			break
		}
	}
}

// Focus sets focus to a specific widget
func (fm *FocusManager) Focus(w Widget) {
	for i, widget := range fm.widgets {
		if widget == w {
			if fm.focusedIndex >= 0 && fm.focusedIndex < len(fm.widgets) {
				fm.widgets[fm.focusedIndex].SetFocused(false)
			}
			fm.focusedIndex = i
			widget.SetFocused(true)
			return
		}
	}
}

// FocusNext moves focus to the next widget
func (fm *FocusManager) FocusNext() {
	if len(fm.widgets) == 0 {
		return
	}

	// Unfocus current
	if fm.focusedIndex >= 0 && fm.focusedIndex < len(fm.widgets) {
		fm.widgets[fm.focusedIndex].SetFocused(false)
	}

	// Find next
	fm.focusedIndex++
	if fm.focusedIndex >= len(fm.widgets) {
		if fm.cycleFocus {
			fm.focusedIndex = 0
		} else {
			fm.focusedIndex = len(fm.widgets) - 1
		}
	}

	fm.widgets[fm.focusedIndex].SetFocused(true)
}

// FocusPrev moves focus to the previous widget
func (fm *FocusManager) FocusPrev() {
	if len(fm.widgets) == 0 {
		return
	}

	// Unfocus current
	if fm.focusedIndex >= 0 && fm.focusedIndex < len(fm.widgets) {
		fm.widgets[fm.focusedIndex].SetFocused(false)
	}

	// Find previous
	fm.focusedIndex--
	if fm.focusedIndex < 0 {
		if fm.cycleFocus {
			fm.focusedIndex = len(fm.widgets) - 1
		} else {
			fm.focusedIndex = 0
		}
	}

	fm.widgets[fm.focusedIndex].SetFocused(true)
}

// Focused returns the currently focused widget
func (fm *FocusManager) Focused() Widget {
	if fm.focusedIndex >= 0 && fm.focusedIndex < len(fm.widgets) {
		return fm.widgets[fm.focusedIndex]
	}
	return nil
}

// HandleEvent passes the event to the focused widget
func (fm *FocusManager) HandleEvent(event input.Event) bool {
	// Handle tab for focus cycling
	if keyEvent, ok := event.(input.KeyEvent); ok {
		if keyEvent.Key == input.KeyTab {
			if keyEvent.IsShift() {
				fm.FocusPrev()
			} else {
				fm.FocusNext()
			}
			return true
		}
	}

	// Pass to focused widget
	if focused := fm.Focused(); focused != nil {
		return focused.HandleEvent(event)
	}

	return false
}

