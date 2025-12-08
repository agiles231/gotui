package widget

import (
	"strings"

	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
	"github.com/agiles231/gotui/terminal"
)

// Text is a static text label widget
type Text struct {
	BaseWidget
	text      string
	style     terminal.Style
	alignment layout.Alignment
	wrap      bool
}

// NewText creates a new text widget
func NewText(text string) *Text {
	return &Text{
		BaseWidget: NewBaseWidget(),
		text:       text,
		style:      terminal.DefaultStyle(),
		alignment:  layout.AlignStart,
	}
}

// SetText sets the text content
func (t *Text) SetText(text string) *Text {
	t.text = text
	return t
}

// Text returns the text content
func (t *Text) Text() string {
	return t.text
}

// SetStyle sets the text style
func (t *Text) SetStyle(style terminal.Style) *Text {
	t.style = style
	return t
}

// SetAlignment sets the text alignment
func (t *Text) SetAlignment(alignment layout.Alignment) *Text {
	t.alignment = alignment
	return t
}

// SetWrap enables or disables text wrapping
func (t *Text) SetWrap(wrap bool) *Text {
	t.wrap = wrap
	return t
}

// Render draws the text widget
func (t *Text) Render(buf *screen.Buffer, bounds layout.Rect) {
	if !t.visible {
		return
	}

	lines := t.getLines(bounds.Width)
	for i, line := range lines {
		if i >= bounds.Height {
			break
		}

		// Calculate x offset based on alignment
		x := bounds.X
		switch t.alignment {
		case layout.AlignCenter:
			x += (bounds.Width - len(line)) / 2
		case layout.AlignEnd:
			x += bounds.Width - len(line)
		}

		buf.DrawStringClipped(x, bounds.Y+i, line, t.style, bounds.Width)
	}
}

// getLines splits text into lines, optionally wrapping
func (t *Text) getLines(maxWidth int) []string {
	if !t.wrap || maxWidth <= 0 {
		return strings.Split(t.text, "\n")
	}

	var result []string
	for _, line := range strings.Split(t.text, "\n") {
		if len(line) <= maxWidth {
			result = append(result, line)
			continue
		}

		// Word wrap
		words := strings.Fields(line)
		current := ""
		for _, word := range words {
			if len(current) == 0 {
				current = word
			} else if len(current)+1+len(word) <= maxWidth {
				current += " " + word
			} else {
				result = append(result, current)
				current = word
			}
		}
		if len(current) > 0 {
			result = append(result, current)
		}
	}

	return result
}

// HandleEvent handles input events (text widgets don't handle input)
func (t *Text) HandleEvent(event input.Event) bool {
	return false
}

// Size returns the preferred size
func (t *Text) Size() layout.Size {
	lines := strings.Split(t.text, "\n")
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	return layout.NewSize(maxWidth, len(lines))
}

// MinSize returns the minimum size
func (t *Text) MinSize() layout.Size {
	return layout.NewSize(1, 1)
}

// Styled text helpers

// Bold creates bold text
func Bold(text string) *Text {
	return NewText(text).SetStyle(terminal.DefaultStyle().WithBold())
}

// Italic creates italic text
func Italic(text string) *Text {
	return NewText(text).SetStyle(terminal.DefaultStyle().WithItalic())
}

// Underline creates underlined text
func Underline(text string) *Text {
	return NewText(text).SetStyle(terminal.DefaultStyle().WithUnderline())
}

// Colored creates colored text
func Colored(text string, fg terminal.Color) *Text {
	return NewText(text).SetStyle(terminal.DefaultStyle().WithFG(fg))
}

