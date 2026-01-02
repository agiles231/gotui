package widget

import (
	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
	"github.com/agiles231/gotui/terminal"
)

// FormField represents a field in a form
type FormField struct {
	Label  string
	Widget Widget
}

// Form is a container for form fields
type Form struct {
	BaseWidget
	fields        []FormField
	buttons       []*Button
	focusedField  int
	focusedButton int // -1 means no button focused
	labelWidth    int
	style         terminal.Style
	labelStyle    terminal.Style
	showBorder    bool
	title         string
	onSubmit      func(values map[string]string)
}

// NewForm creates a new form widget
func NewForm() *Form {
	f := &Form{
		BaseWidget:    NewBaseWidget(),
		labelWidth:    15,
		style:         terminal.DefaultStyle(),
		labelStyle:    terminal.DefaultStyle().WithBold(),
		focusedButton: -1, // No button focused initially
	}
	f.SetInteractive(true)
	return f
}

// AddField adds a field to the form
func (f *Form) AddField(label string, widget Widget) *Form {
	f.fields = append(f.fields, FormField{
		Label:  label,
		Widget: widget,
	})
	
	// Auto-focus first interactive field
	if len(f.fields) == 1 && widget.IsInteractive() {
		f.focusedField = 0
		widget.SetFocused(true)
	}
	
	return f
}

// AddTextInput adds a text input field
func (f *Form) AddTextInput(label, placeholder string) *TextInput {
	ti := NewTextInput().SetPlaceholder(placeholder)
	f.AddField(label, ti)
	return ti
}

// AddPasswordInput adds a password input field
func (f *Form) AddPasswordInput(label, placeholder string) *TextInput {
	ti := NewTextInput().SetPlaceholder(placeholder).SetMask('*')
	f.AddField(label, ti)
	return ti
}

// AddButton adds a button to the form's button row
func (f *Form) AddButton(label string, onPress func()) *Button {
	btn := NewButton(label).OnPress(onPress)
	f.buttons = append(f.buttons, btn)
	return btn
}

// Fields returns the form fields
func (f *Form) Fields() []FormField {
	return f.fields
}

// SetLabelWidth sets the label column width
func (f *Form) SetLabelWidth(width int) *Form {
	f.labelWidth = width
	return f
}

// SetStyle sets the form style
func (f *Form) SetStyle(style terminal.Style) *Form {
	f.style = style
	return f
}

// SetLabelStyle sets the label style
func (f *Form) SetLabelStyle(style terminal.Style) *Form {
	f.labelStyle = style
	return f
}

// SetShowBorder enables or disables the border
func (f *Form) SetShowBorder(show bool) *Form {
	f.showBorder = show
	return f
}

// SetTitle sets the form title
func (f *Form) SetTitle(title string) *Form {
	f.title = title
	return f
}

// OnSubmit sets the submit callback
func (f *Form) OnSubmit(fn func(values map[string]string)) *Form {
	f.onSubmit = fn
	return f
}

// Values returns a map of field labels to values
func (f *Form) Values() map[string]string {
	values := make(map[string]string)
	for _, field := range f.fields {
		if ti, ok := field.Widget.(*TextInput); ok {
			values[field.Label] = ti.Value()
		}
	}
	return values
}

// Render draws the form
func (f *Form) Render(buf *screen.Buffer, bounds layout.Rect) {
	if !f.visible {
		return
	}

	innerBounds := bounds
	y := bounds.Y

	if f.showBorder {
		height := len(f.fields) + 2
		if f.title != "" {
			height++
		}
		if len(f.buttons) > 0 {
			height += 2 // Extra row for buttons + spacing
		}
		buf.DrawBox(bounds.X, bounds.Y, bounds.Z, bounds.Width, height, f.style)
		innerBounds = bounds.Inset(1, 1, 1, 1)
		y = innerBounds.Y
	}

	// Draw title
	if f.title != "" {
		titleX := innerBounds.X + (innerBounds.Width-len(f.title))/2
		buf.DrawString(titleX, y, innerBounds.Z, f.title, f.labelStyle)
		y++
	}

	// Draw fields
	for i, field := range f.fields {
		// Draw label
		label := field.Label
		if len(label) > f.labelWidth-2 {
			label = label[:f.labelWidth-2]
		}
		label += ": "
		buf.DrawString(innerBounds.X, y+i, innerBounds.Z, label, f.labelStyle)

		// Calculate widget bounds
		widgetBounds := layout.NewRect(
			innerBounds.X+f.labelWidth,
			y+i,
			innerBounds.Z,
			innerBounds.Width-f.labelWidth,
			1,
		)

		// Render widget
		field.Widget.Render(buf, widgetBounds)
	}

	// Draw buttons row
	if len(f.buttons) > 0 {
		buttonY := y + len(f.fields) + 1 // Leave a blank line
		buttonX := innerBounds.X
		
		for i, btn := range f.buttons {
			btnWidth := btn.Size().Width
			btnBounds := layout.NewRect(buttonX, buttonY, innerBounds.Z, btnWidth, 1)
			
			// Set button focus state
			btn.SetFocused(f.focusedButton == i)
			btn.Render(buf, btnBounds)
			
			buttonX += btnWidth + 2 // 2 spaces between buttons
		}
	}
}

// HandleEvent handles input events
func (f *Form) HandleEvent(event input.Event) bool {
	if !f.focused {
		return false
	}
	
	// Need at least fields or buttons
	if len(f.fields) == 0 && len(f.buttons) == 0 {
		return false
	}

	keyEvent, ok := event.(input.KeyEvent)
	if !ok {
		return false
	}

	// Handle tab/shift-tab for navigation
	if keyEvent.Key == input.KeyTab {
		if keyEvent.IsShift() {
			f.focusPrev()
		} else {
			f.focusNext()
		}
		return true
	}

	// If on button row, use Left/Right for button navigation
	if f.focusedButton >= 0 {
		if keyEvent.Key == input.KeyLeft {
			f.focusPrevButton()
			return true
		}
		if keyEvent.Key == input.KeyRight {
			f.focusNextButton()
			return true
		}
		// Up moves from buttons back to last field
		if keyEvent.Key == input.KeyUp {
			f.focusLastField()
			return true
		}
		// Enter/Space triggers button
		if keyEvent.Key == input.KeyEnter || (keyEvent.Key == input.KeyRune && keyEvent.Rune == ' ') {
			btn := f.buttons[f.focusedButton]
			return btn.HandleEvent(event)
		}
		// Consume Down arrow on buttons (nowhere to go)
		if keyEvent.Key == input.KeyDown {
			return true
		}
	}

	// If on a field, handle Up/Down for field navigation
	if f.focusedButton < 0 && f.focusedField >= 0 {
		if keyEvent.Key == input.KeyUp {
			if f.focusedField > 0 {
				f.focusPrevField()
			}
			return true // Consume even if at top
		}
		if keyEvent.Key == input.KeyDown {
			if f.focusedField < len(f.fields)-1 {
				f.focusNextField()
			} else if len(f.buttons) > 0 {
				// Move to first button
				f.focusFirstButton()
			}
			return true
		}
		// Consume Left/Right in fields (let TextInput handle cursor, but don't propagate)
		if keyEvent.Key == input.KeyLeft || keyEvent.Key == input.KeyRight {
			// Pass to field widget for cursor movement
			if f.focusedField >= 0 && f.focusedField < len(f.fields) {
				f.fields[f.focusedField].Widget.HandleEvent(event)
			}
			return true // Always consume to prevent tab switching
		}
	}

	// Pass other events to focused field
	if f.focusedButton < 0 && f.focusedField >= 0 && f.focusedField < len(f.fields) {
		return f.fields[f.focusedField].Widget.HandleEvent(event)
	}

	return false
}

// focusNextButton moves focus to next button
func (f *Form) focusNextButton() {
	if len(f.buttons) == 0 {
		return
	}
	if f.focusedButton >= 0 {
		f.buttons[f.focusedButton].SetFocused(false)
	}
	f.focusedButton++
	if f.focusedButton >= len(f.buttons) {
		f.focusedButton = 0
	}
	f.buttons[f.focusedButton].SetFocused(true)
}

// focusPrevButton moves focus to previous button
func (f *Form) focusPrevButton() {
	if len(f.buttons) == 0 {
		return
	}
	if f.focusedButton >= 0 {
		f.buttons[f.focusedButton].SetFocused(false)
	}
	f.focusedButton--
	if f.focusedButton < 0 {
		f.focusedButton = len(f.buttons) - 1
	}
	f.buttons[f.focusedButton].SetFocused(true)
}

// focusFirstButton moves focus to first button
func (f *Form) focusFirstButton() {
	if len(f.buttons) == 0 {
		return
	}
	// Unfocus current field
	if f.focusedField >= 0 && f.focusedField < len(f.fields) {
		f.fields[f.focusedField].Widget.SetFocused(false)
	}
	f.focusedField = -1
	f.focusedButton = 0
	f.buttons[0].SetFocused(true)
}

// focusLastField moves focus from buttons back to last field
func (f *Form) focusLastField() {
	if len(f.fields) == 0 {
		return
	}
	// Unfocus current button
	if f.focusedButton >= 0 && f.focusedButton < len(f.buttons) {
		f.buttons[f.focusedButton].SetFocused(false)
	}
	f.focusedButton = -1
	f.focusedField = len(f.fields) - 1
	f.fields[f.focusedField].Widget.SetFocused(true)
}

// focusNextField moves focus to next field
func (f *Form) focusNextField() {
	if len(f.fields) == 0 {
		return
	}
	if f.focusedField >= 0 && f.focusedField < len(f.fields) {
		f.fields[f.focusedField].Widget.SetFocused(false)
	}
	f.focusedField++
	if f.focusedField >= len(f.fields) {
		f.focusedField = len(f.fields) - 1
	}
	f.fields[f.focusedField].Widget.SetFocused(true)
}

// focusPrevField moves focus to previous field
func (f *Form) focusPrevField() {
	if len(f.fields) == 0 {
		return
	}
	if f.focusedField >= 0 && f.focusedField < len(f.fields) {
		f.fields[f.focusedField].Widget.SetFocused(false)
	}
	f.focusedField--
	if f.focusedField < 0 {
		f.focusedField = 0
	}
	f.fields[f.focusedField].Widget.SetFocused(true)
}

func (f *Form) focusNext() {
	totalFields := len(f.fields)
	totalButtons := len(f.buttons)
	
	if totalFields == 0 && totalButtons == 0 {
		return
	}

	// Unfocus current field or button
	if f.focusedButton >= 0 && f.focusedButton < totalButtons {
		f.buttons[f.focusedButton].SetFocused(false)
	} else if f.focusedField >= 0 && f.focusedField < totalFields {
		f.fields[f.focusedField].Widget.SetFocused(false)
	}

	// Currently on a button?
	if f.focusedButton >= 0 {
		// Move to next button or wrap to first field
		f.focusedButton++
		if f.focusedButton >= totalButtons {
			// Wrap to first field
			f.focusedButton = -1
			f.focusedField = 0
			if totalFields > 0 {
				f.fields[0].Widget.SetFocused(true)
			}
			return
		}
		f.buttons[f.focusedButton].SetFocused(true)
		return
	}

	// Currently on a field - move to next field or first button
	f.focusedField++
	if f.focusedField >= totalFields {
		// Move to buttons if any, otherwise wrap to first field
		if totalButtons > 0 {
			f.focusedField = -1
			f.focusedButton = 0
			f.buttons[0].SetFocused(true)
			return
		}
		f.focusedField = 0
	}
	
	if f.focusedField >= 0 && f.focusedField < totalFields {
		f.fields[f.focusedField].Widget.SetFocused(true)
	}
}

func (f *Form) focusPrev() {
	totalFields := len(f.fields)
	totalButtons := len(f.buttons)
	
	if totalFields == 0 && totalButtons == 0 {
		return
	}

	// Unfocus current field or button
	if f.focusedButton >= 0 && f.focusedButton < totalButtons {
		f.buttons[f.focusedButton].SetFocused(false)
	} else if f.focusedField >= 0 && f.focusedField < totalFields {
		f.fields[f.focusedField].Widget.SetFocused(false)
	}

	// Currently on a button?
	if f.focusedButton >= 0 {
		// Move to previous button or wrap to last field
		f.focusedButton--
		if f.focusedButton < 0 {
			// Wrap to last field
			if totalFields > 0 {
				f.focusedField = totalFields - 1
				f.fields[f.focusedField].Widget.SetFocused(true)
			} else {
				// No fields, wrap to last button
				f.focusedButton = totalButtons - 1
				f.buttons[f.focusedButton].SetFocused(true)
			}
			return
		}
		f.buttons[f.focusedButton].SetFocused(true)
		return
	}

	// Currently on a field - move to previous field or last button
	f.focusedField--
	if f.focusedField < 0 {
		// Move to last button if any, otherwise wrap to last field
		if totalButtons > 0 {
			f.focusedButton = totalButtons - 1
			f.buttons[f.focusedButton].SetFocused(true)
			return
		}
		f.focusedField = totalFields - 1
	}
	
	if f.focusedField >= 0 && f.focusedField < totalFields {
		f.fields[f.focusedField].Widget.SetFocused(true)
	}
}

// SetFocused sets focus state
func (f *Form) SetFocused(focused bool) {
	f.focused = focused
	
	// Focus/unfocus the current field or button
	if f.focusedButton >= 0 && f.focusedButton < len(f.buttons) {
		f.buttons[f.focusedButton].SetFocused(focused)
	} else if f.focusedField >= 0 && f.focusedField < len(f.fields) {
		f.fields[f.focusedField].Widget.SetFocused(focused)
	}
}

// Size returns the preferred size
func (f *Form) Size() layout.Size {
	height := len(f.fields)
	width := f.labelWidth + 20 // Default input width

	if f.title != "" {
		height++
	}
	if len(f.buttons) > 0 {
		height += 2 // Spacing + button row
	}
	if f.showBorder {
		width += 2
		height += 2
	}

	return layout.NewSize(width, height)
}

// MinSize returns the minimum size
func (f *Form) MinSize() layout.Size {
	return layout.NewSize(f.labelWidth+10, len(f.fields)+2)
}

