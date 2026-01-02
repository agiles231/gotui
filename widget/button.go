package widget

import (
	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
	"github.com/agiles231/gotui/terminal"
)

// Button is a clickable button widget
type Button struct {
	BaseWidget
	label        string
	onPress      func()
	style        terminal.Style
	focusedStyle terminal.Style
	width        int
}

// NewButton creates a new button with the given label
func NewButton(label string) *Button {
	b := &Button{
		BaseWidget:   NewBaseWidget(),
		label:        label,
		style:        terminal.DefaultStyle(),
		focusedStyle: terminal.DefaultStyle().WithReverse(),
	}
	b.SetInteractive(true)
	return b
}

// SetLabel sets the button label
func (b *Button) SetLabel(label string) *Button {
	b.label = label
	return b
}

// Label returns the button label
func (b *Button) Label() string {
	return b.label
}

// OnPress sets the callback for when the button is pressed
func (b *Button) OnPress(fn func()) *Button {
	b.onPress = fn
	return b
}

// SetStyle sets the normal style
func (b *Button) SetStyle(style terminal.Style) *Button {
	b.style = style
	return b
}

// SetFocusedStyle sets the style when focused
func (b *Button) SetFocusedStyle(style terminal.Style) *Button {
	b.focusedStyle = style
	return b
}

// SetWidth sets a fixed width for the button
func (b *Button) SetWidth(width int) *Button {
	b.width = width
	return b
}

// Render draws the button
func (b *Button) Render(buf *screen.Buffer, bounds layout.Rect) {
	if !b.visible {
		return
	}

	style := b.style
	if b.focused {
		style = b.focusedStyle
	}

	// Format: [ Label ]
	text := "[ " + b.label + " ]"
	
	// Pad to width if specified
	if b.width > 0 && len(text) < b.width {
		padding := b.width - len(text)
		leftPad := padding / 2
		rightPad := padding - leftPad
		for i := 0; i < leftPad; i++ {
			text = " " + text
		}
		for i := 0; i < rightPad; i++ {
			text = text + " "
		}
	}

	// Center in bounds if text is shorter than bounds width
	x := bounds.X
	if len(text) < bounds.Width {
		x = bounds.X + (bounds.Width-len(text))/2
	}

	buf.DrawString(x, bounds.Y, bounds.Z, text, style)
}

// HandleEvent handles input events
func (b *Button) HandleEvent(event input.Event) bool {
	if !b.focused {
		return false
	}

	keyEvent, ok := event.(input.KeyEvent)
	if !ok {
		return false
	}

	// Trigger on Enter or Space
	if keyEvent.Key == input.KeyEnter || (keyEvent.Key == input.KeyRune && keyEvent.Rune == ' ') {
		if b.onPress != nil {
			b.onPress()
		}
		return true
	}

	return false
}

// Size returns the preferred size
func (b *Button) Size() layout.Size {
	width := len(b.label) + 4 // "[ " + label + " ]"
	if b.width > 0 {
		width = b.width
	}
	return layout.NewSize(width, 1)
}

// MinSize returns the minimum size
func (b *Button) MinSize() layout.Size {
	return layout.NewSize(len(b.label)+4, 1)
}

