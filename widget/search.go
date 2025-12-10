package widget

import (
	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
	"github.com/agiles231/gotui/terminal"
)

type Search struct {
	BaseWidget
	placeholder string
	value       string
	cursor      int
	onChange    func(string)
	onSubmit    func(string)
	style       terminal.Style
}

func NewSearch() *Search {
	s := &Search{
		BaseWidget: NewBaseWidget(),
		style: terminal.DefaultStyle(),
	}
	return s
}

func (s *Search) SetPlaceholder(placeholder string) *Search {
	s.placeholder = placeholder
	return s
}

func (s *Search) SetValue(value string) *Search {
	s.value = value
	return s
}

func (s *Search) SetOnChange(onChange func(string)) *Search {
	s.onChange = onChange
	return s
}

func (s *Search) SetOnSubmit(onSubmit func(string)) *Search {
	s.onSubmit = onSubmit
	return s
}

func (s *Search) SetStyle(style terminal.Style) *Search {
	s.style = style
	return s
}

// Meet interface for Widget
func (s *Search) Render(buf *screen.Buffer, bounds layout.Rect) {
	if !s.IsVisible() {
		return
	}

	// border
	buf.DrawBox(bounds.X, bounds.Y, bounds.Width, bounds.Height, s.style)
	bounds = bounds.InsetAll(1)

	// placeholder or value
	if s.value == "" {
		buf.DrawString(bounds.X, bounds.Y, s.placeholder, s.style)
	} else {
		buf.DrawString(bounds.X, bounds.Y, s.value, s.style)
	}

}

func (s *Search) HandleEvent(event input.Event) bool {
	event_key, ok := event.(input.KeyEvent)
	if !ok {
		return false
	}
	if event_key.Key == input.KeyEnter {
		if s.onSubmit != nil {
			s.onSubmit(s.value)
		}
		return true
	}
	if event_key.Key == input.KeyBackspace {
		if s.cursor == 0 {
			return true
		}
		if s.cursor == len(s.value) {
			s.value = s.value[:s.cursor-1]
			s.cursor--
			return true
		}
		s.value = s.value[:s.cursor-1] + s.value[s.cursor:]
		return true
	}
	if event_key.Key == input.KeyDelete {
		s.value = s.value[:s.cursor] + s.value[s.cursor+1:]
		return true
	}
	if event_key.Key == input.KeyLeft {
		s.cursor--
		return true
	}
	if event_key.Key == input.KeyRight {
		s.cursor++
		if s.cursor > len(s.value) {
			s.cursor = len(s.value)
		}
		return true
	}
	if event_key.Key == input.KeyHome {
		s.cursor = 0
		return true
	}
	if event_key.Key == input.KeyEnd {
		s.cursor = len(s.value)
		return true
	}
	if event_key.Key == input.KeyRune {
		if len(s.value) == 0 {
			s.value = string(event_key.Rune)
			s.cursor = 1
			return true
		}
		if len(s.value) == s.cursor {
			s.value = s.value + string(event_key.Rune)
			s.cursor++
			return true
		}
		if s.cursor == 0 {
			s.value = string(event_key.Rune) + s.value
			s.cursor = 1
			return true
		}
		s.value = s.value[:s.cursor] + string(event_key.Rune) + s.value[s.cursor:]
		s.cursor++
		return true
	}
	return false
}

func (s *Search) Size() layout.Size {
	return layout.NewSize(100, 20)
}

func (s *Search) MinSize() layout.Size {
	return layout.NewSize(100, 20)
}

func (s *Search) SetFocused(focused bool) {
	s.BaseWidget.SetFocused(focused)
}

func (s *Search) IsFocused() bool {
	return s.BaseWidget.IsFocused()
}

func (s *Search) IsInteractive() bool {
	return s.BaseWidget.IsInteractive()
}