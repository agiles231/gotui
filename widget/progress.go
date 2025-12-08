package widget

import (
	"fmt"

	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
	"github.com/agiles231/gotui/terminal"
)

// Progress is a progress bar widget
type Progress struct {
	BaseWidget
	value       float64 // 0.0 to 1.0
	width       int
	showPercent bool
	showValue   bool
	style       terminal.Style
	fillStyle   terminal.Style
	fillChar    rune
	emptyChar   rune
	label       string
}

// NewProgress creates a new progress bar
func NewProgress() *Progress {
	return &Progress{
		BaseWidget: NewBaseWidget(),
		width:      20,
		showPercent: true,
		style:      terminal.DefaultStyle(),
		fillStyle:  terminal.DefaultStyle().WithFG(terminal.ColorGreen),
		fillChar:   '█',
		emptyChar:  '░',
	}
}

// SetValue sets the progress value (0.0 to 1.0)
func (p *Progress) SetValue(value float64) *Progress {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	p.value = value
	return p
}

// Value returns the current progress value
func (p *Progress) Value() float64 {
	return p.value
}

// SetPercent sets the progress as a percentage (0 to 100)
func (p *Progress) SetPercent(percent int) *Progress {
	return p.SetValue(float64(percent) / 100)
}

// Percent returns the progress as a percentage
func (p *Progress) Percent() int {
	return int(p.value * 100)
}

// SetWidth sets the progress bar width
func (p *Progress) SetWidth(width int) *Progress {
	p.width = width
	return p
}

// SetShowPercent enables or disables percentage display
func (p *Progress) SetShowPercent(show bool) *Progress {
	p.showPercent = show
	return p
}

// SetShowValue enables or disables value display
func (p *Progress) SetShowValue(show bool) *Progress {
	p.showValue = show
	return p
}

// SetLabel sets a label for the progress bar
func (p *Progress) SetLabel(label string) *Progress {
	p.label = label
	return p
}

// SetStyle sets the background style
func (p *Progress) SetStyle(style terminal.Style) *Progress {
	p.style = style
	return p
}

// SetFillStyle sets the fill style
func (p *Progress) SetFillStyle(style terminal.Style) *Progress {
	p.fillStyle = style
	return p
}

// SetChars sets the fill and empty characters
func (p *Progress) SetChars(fill, empty rune) *Progress {
	p.fillChar = fill
	p.emptyChar = empty
	return p
}

// Render draws the progress bar
func (p *Progress) Render(buf *screen.Buffer, bounds layout.Rect) {
	if !p.visible {
		return
	}

	x := bounds.X
	y := bounds.Y
	width := p.width
	if width > bounds.Width {
		width = bounds.Width
	}

	// Draw label if present
	if p.label != "" {
		buf.DrawString(x, y, p.label+": ", p.style)
		x += len(p.label) + 2
		width -= len(p.label) + 2
	}

	// Reserve space for percentage
	percentText := ""
	if p.showPercent {
		percentText = fmt.Sprintf(" %3d%%", p.Percent())
		width -= len(percentText)
	}

	// Calculate filled portion
	filled := int(float64(width) * p.value)

	// Draw progress bar
	for i := 0; i < width; i++ {
		var cell screen.Cell
		if i < filled {
			cell = screen.NewCell(p.fillChar, p.fillStyle)
		} else {
			cell = screen.NewCell(p.emptyChar, p.style)
		}
		buf.Set(x+i, y, cell)
	}

	// Draw percentage
	if p.showPercent {
		buf.DrawString(x+width, y, percentText, p.style)
	}
}

// HandleEvent handles input events (progress bars don't handle input)
func (p *Progress) HandleEvent(event input.Event) bool {
	return false
}

// Size returns the preferred size
func (p *Progress) Size() layout.Size {
	width := p.width
	if p.label != "" {
		width += len(p.label) + 2
	}
	if p.showPercent {
		width += 5
	}
	return layout.NewSize(width, 1)
}

// MinSize returns the minimum size
func (p *Progress) MinSize() layout.Size {
	return layout.NewSize(10, 1)
}

// Spinner is an indeterminate progress indicator
type Spinner struct {
	BaseWidget
	frames  []rune
	current int
	style   terminal.Style
	label   string
}

// NewSpinner creates a new spinner widget
func NewSpinner() *Spinner {
	return &Spinner{
		BaseWidget: NewBaseWidget(),
		frames:     []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'},
		style:      terminal.DefaultStyle(),
	}
}

// SetFrames sets the animation frames
func (s *Spinner) SetFrames(frames []rune) *Spinner {
	s.frames = frames
	return s
}

// SetStyle sets the spinner style
func (s *Spinner) SetStyle(style terminal.Style) *Spinner {
	s.style = style
	return s
}

// SetLabel sets the spinner label
func (s *Spinner) SetLabel(label string) *Spinner {
	s.label = label
	return s
}

// Advance moves to the next frame
func (s *Spinner) Advance() {
	s.current = (s.current + 1) % len(s.frames)
}

// Render draws the spinner
func (s *Spinner) Render(buf *screen.Buffer, bounds layout.Rect) {
	if !s.visible || len(s.frames) == 0 {
		return
	}

	buf.Set(bounds.X, bounds.Y, screen.NewCell(s.frames[s.current], s.style))

	if s.label != "" {
		buf.DrawString(bounds.X+2, bounds.Y, s.label, s.style)
	}
}

// HandleEvent handles input events (spinners don't handle input)
func (s *Spinner) HandleEvent(event input.Event) bool {
	return false
}

// Size returns the preferred size
func (s *Spinner) Size() layout.Size {
	width := 1
	if s.label != "" {
		width += 1 + len(s.label)
	}
	return layout.NewSize(width, 1)
}

// MinSize returns the minimum size
func (s *Spinner) MinSize() layout.Size {
	return layout.NewSize(1, 1)
}

