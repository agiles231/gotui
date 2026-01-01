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
	focusedField  int
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
		BaseWidget: NewBaseWidget(),
		labelWidth: 15,
		style:      terminal.DefaultStyle(),
		labelStyle: terminal.DefaultStyle().WithBold(),
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
}

// HandleEvent handles input events
func (f *Form) HandleEvent(event input.Event) bool {
	if !f.focused || len(f.fields) == 0 {
		return false
	}

	keyEvent, ok := event.(input.KeyEvent)
	if !ok {
		return false
	}

	// Handle tab/shift-tab for field navigation
	if keyEvent.Key == input.KeyTab {
		if keyEvent.IsShift() {
			f.focusPrev()
		} else {
			f.focusNext()
		}
		return true
	}

	// Handle up/down for field navigation
	if keyEvent.Key == input.KeyUp {
		f.focusPrev()
		return true
	}
	if keyEvent.Key == input.KeyDown {
		f.focusNext()
		return true
	}

	// Handle submit (Ctrl+Enter or just Enter if on non-input field)
	if keyEvent.Key == input.KeyEnter && keyEvent.IsCtrl() {
		if f.onSubmit != nil {
			f.onSubmit(f.Values())
		}
		return true
	}

	// Pass event to focused field
	if f.focusedField >= 0 && f.focusedField < len(f.fields) {
		return f.fields[f.focusedField].Widget.HandleEvent(event)
	}

	return false
}

func (f *Form) focusNext() {
	if len(f.fields) == 0 {
		return
	}

	// Unfocus current
	if f.focusedField >= 0 && f.focusedField < len(f.fields) {
		f.fields[f.focusedField].Widget.SetFocused(false)
	}

	// Find next interactive field
	start := f.focusedField
	for {
		f.focusedField++
		if f.focusedField >= len(f.fields) {
			f.focusedField = 0
		}
		
		if f.fields[f.focusedField].Widget.IsInteractive() {
			break
		}
		
		if f.focusedField == start {
			break // No interactive fields
		}
	}

	f.fields[f.focusedField].Widget.SetFocused(true)
}

func (f *Form) focusPrev() {
	if len(f.fields) == 0 {
		return
	}

	// Unfocus current
	if f.focusedField >= 0 && f.focusedField < len(f.fields) {
		f.fields[f.focusedField].Widget.SetFocused(false)
	}

	// Find previous interactive field
	start := f.focusedField
	for {
		f.focusedField--
		if f.focusedField < 0 {
			f.focusedField = len(f.fields) - 1
		}
		
		if f.fields[f.focusedField].Widget.IsInteractive() {
			break
		}
		
		if f.focusedField == start {
			break // No interactive fields
		}
	}

	f.fields[f.focusedField].Widget.SetFocused(true)
}

// SetFocused sets focus state
func (f *Form) SetFocused(focused bool) {
	f.focused = focused
	
	// Focus/unfocus the current field
	if f.focusedField >= 0 && f.focusedField < len(f.fields) {
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

