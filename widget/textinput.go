package widget

import (
	"strings"

	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
	"github.com/agiles231/gotui/terminal"
)

// TextInput is a single-line text input widget
type TextInput struct {
	BaseWidget
	value        []rune
	cursor       int
	offset       int // Scroll offset for long text
	placeholder  string
	style        terminal.Style
	focusedStyle terminal.Style
	cursorStyle  terminal.Style
	width        int
	mask         rune // For password fields
	onChange     func(string)
	onSubmit     func(string)
}

// NewTextInput creates a new text input widget
func NewTextInput() *TextInput {
	ti := &TextInput{
		BaseWidget:   NewBaseWidget(),
		value:        []rune{},
		style:        terminal.DefaultStyle(),
		focusedStyle: terminal.DefaultStyle().WithReverse(),
		cursorStyle:  terminal.DefaultStyle().WithReverse(),
		width:        20,
	}
	ti.SetInteractive(true)
	return ti
}

// SetValue sets the input value
func (ti *TextInput) SetValue(value string) *TextInput {
	ti.value = []rune(value)
	if ti.cursor > len(ti.value) {
		ti.cursor = len(ti.value)
	}
	ti.updateOffset()
	return ti
}

// Value returns the current input value
func (ti *TextInput) Value() string {
	return string(ti.value)
}

// SetPlaceholder sets placeholder text
func (ti *TextInput) SetPlaceholder(placeholder string) *TextInput {
	ti.placeholder = placeholder
	return ti
}

// SetWidth sets the input width
func (ti *TextInput) SetWidth(width int) *TextInput {
	ti.width = width
	return ti
}

// SetMask sets a mask character for password fields
func (ti *TextInput) SetMask(mask rune) *TextInput {
	ti.mask = mask
	return ti
}

// SetStyle sets the normal style
func (ti *TextInput) SetStyle(style terminal.Style) *TextInput {
	ti.style = style
	return ti
}

// SetFocusedStyle sets the focused style
func (ti *TextInput) SetFocusedStyle(style terminal.Style) *TextInput {
	ti.focusedStyle = style
	return ti
}

// OnChange sets the change callback
func (ti *TextInput) OnChange(fn func(string)) *TextInput {
	ti.onChange = fn
	return ti
}

// OnSubmit sets the submit callback (Enter key)
func (ti *TextInput) OnSubmit(fn func(string)) *TextInput {
	ti.onSubmit = fn
	return ti
}

// Render draws the text input
func (ti *TextInput) Render(buf *screen.Buffer, bounds layout.Rect) {
	if !ti.visible {
		return
	}

	style := ti.style
	if ti.focused {
		style = ti.focusedStyle
	}

	width := bounds.Width
	if ti.width > 0 && ti.width < width {
		width = ti.width
	}

	// Draw background
	for x := 0; x < width; x++ {
		buf.Set(bounds.X+x, bounds.Y, bounds.Z, screen.NewCell(' ', style))
	}

	// Get display text
	var displayText string
	if len(ti.value) == 0 && !ti.focused {
		// Show placeholder
		displayText = ti.placeholder
		style = style.WithDim()
	} else if ti.mask != 0 {
		// Masked text
		displayText = strings.Repeat(string(ti.mask), len(ti.value))
	} else {
		displayText = string(ti.value)
	}

	// Calculate visible portion
	visibleStart := ti.offset
	visibleEnd := ti.offset + width
	if visibleEnd > len(displayText) {
		visibleEnd = len(displayText)
	}

	if visibleStart < len(displayText) {
		visible := displayText[visibleStart:visibleEnd]
		buf.DrawStringClipped(bounds.X, bounds.Y, bounds.Z, visible, style, width)
	}

	// Draw cursor if focused
	if ti.focused {
		cursorX := bounds.X + ti.cursor - ti.offset
		if cursorX >= bounds.X && cursorX < bounds.X+width {
			var cursorChar rune = ' '
			if ti.cursor < len(ti.value) {
				if ti.mask != 0 {
					cursorChar = ti.mask
				} else {
					cursorChar = ti.value[ti.cursor]
				}
			}
			buf.Set(cursorX, bounds.Y, bounds.Z, screen.NewCell(cursorChar, ti.cursorStyle))
		}
	}
}

// HandleEvent handles input events
func (ti *TextInput) HandleEvent(event input.Event) bool {
	if !ti.focused {
		return false
	}

	keyEvent, ok := event.(input.KeyEvent)
	if !ok {
		return false
	}

	switch keyEvent.Key {
	case input.KeyRune:
		ti.insert(keyEvent.Rune)
		return true

	case input.KeyBackspace:
		ti.backspace()
		return true

	case input.KeyDelete:
		ti.delete()
		return true

	case input.KeyLeft:
		if keyEvent.IsCtrl() {
			ti.wordLeft()
		} else {
			ti.cursorLeft()
		}
		return true

	case input.KeyRight:
		if keyEvent.IsCtrl() {
			ti.wordRight()
		} else {
			ti.cursorRight()
		}
		return true

	case input.KeyHome:
		ti.cursor = 0
		ti.updateOffset()
		return true

	case input.KeyEnd:
		ti.cursor = len(ti.value)
		ti.updateOffset()
		return true

	case input.KeyEnter:
		if ti.onSubmit != nil {
			ti.onSubmit(string(ti.value))
		}
		return true
	}

	// Handle Ctrl+key
	if keyEvent.IsCtrl() && keyEvent.Key == input.KeyRune {
		switch keyEvent.Rune {
		case 'a': // Ctrl+A: home
			ti.cursor = 0
			ti.updateOffset()
			return true
		case 'e': // Ctrl+E: end
			ti.cursor = len(ti.value)
			ti.updateOffset()
			return true
		case 'k': // Ctrl+K: kill to end
			ti.value = ti.value[:ti.cursor]
			ti.notifyChange()
			return true
		case 'u': // Ctrl+U: kill to start
			ti.value = ti.value[ti.cursor:]
			ti.cursor = 0
			ti.updateOffset()
			ti.notifyChange()
			return true
		case 'w': // Ctrl+W: kill word
			ti.deleteWord()
			return true
		}
	}

	return false
}

// insert inserts a character at the cursor
func (ti *TextInput) insert(r rune) {
	ti.value = append(ti.value[:ti.cursor], append([]rune{r}, ti.value[ti.cursor:]...)...)
	ti.cursor++
	ti.updateOffset()
	ti.notifyChange()
}

// backspace deletes the character before the cursor
func (ti *TextInput) backspace() {
	if ti.cursor > 0 {
		ti.value = append(ti.value[:ti.cursor-1], ti.value[ti.cursor:]...)
		ti.cursor--
		ti.updateOffset()
		ti.notifyChange()
	}
}

// delete deletes the character at the cursor
func (ti *TextInput) delete() {
	if ti.cursor < len(ti.value) {
		ti.value = append(ti.value[:ti.cursor], ti.value[ti.cursor+1:]...)
		ti.notifyChange()
	}
}

// cursorLeft moves cursor left
func (ti *TextInput) cursorLeft() {
	if ti.cursor > 0 {
		ti.cursor--
		ti.updateOffset()
	}
}

// cursorRight moves cursor right
func (ti *TextInput) cursorRight() {
	if ti.cursor < len(ti.value) {
		ti.cursor++
		ti.updateOffset()
	}
}

// wordLeft moves cursor to the start of the previous word
func (ti *TextInput) wordLeft() {
	// Skip spaces
	for ti.cursor > 0 && ti.value[ti.cursor-1] == ' ' {
		ti.cursor--
	}
	// Skip word
	for ti.cursor > 0 && ti.value[ti.cursor-1] != ' ' {
		ti.cursor--
	}
	ti.updateOffset()
}

// wordRight moves cursor to the end of the next word
func (ti *TextInput) wordRight() {
	// Skip word
	for ti.cursor < len(ti.value) && ti.value[ti.cursor] != ' ' {
		ti.cursor++
	}
	// Skip spaces
	for ti.cursor < len(ti.value) && ti.value[ti.cursor] == ' ' {
		ti.cursor++
	}
	ti.updateOffset()
}

// deleteWord deletes the word before the cursor
func (ti *TextInput) deleteWord() {
	if ti.cursor == 0 {
		return
	}

	end := ti.cursor
	// Skip spaces
	for ti.cursor > 0 && ti.value[ti.cursor-1] == ' ' {
		ti.cursor--
	}
	// Skip word
	for ti.cursor > 0 && ti.value[ti.cursor-1] != ' ' {
		ti.cursor--
	}

	ti.value = append(ti.value[:ti.cursor], ti.value[end:]...)
	ti.updateOffset()
	ti.notifyChange()
}

// updateOffset updates the scroll offset
func (ti *TextInput) updateOffset() {
	width := ti.width
	if width <= 0 {
		width = 20
	}

	// Keep cursor visible
	if ti.cursor < ti.offset {
		ti.offset = ti.cursor
	}
	if ti.cursor >= ti.offset+width {
		ti.offset = ti.cursor - width + 1
	}
}

// notifyChange calls the change callback
func (ti *TextInput) notifyChange() {
	if ti.onChange != nil {
		ti.onChange(string(ti.value))
	}
}

// Size returns the preferred size
func (ti *TextInput) Size() layout.Size {
	return layout.NewSize(ti.width, 1)
}

// MinSize returns the minimum size
func (ti *TextInput) MinSize() layout.Size {
	return layout.NewSize(5, 1)
}

