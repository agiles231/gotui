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
	helpItems   []string
	help        *Table
	helpVisible bool
	onChange    func(string)
	onSubmit    func(string)
	style       terminal.Style
}

func NewSearch() *Search {
	s := &Search{
		BaseWidget: NewBaseWidget(),
		style: terminal.DefaultStyle(),
		helpVisible: true,
	}
	return s
}

func (s *Search) SetHelpItems(helpItems []string) *Search {
	s.helpItems = helpItems
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

func (s *Search) getHelp(width int) *Table {
	if s.help == nil {
		s.help = NewTable()
		s.help.SetShowHeader(false)
		s.help.SetShowBorder(false)
		s.help.SetStyle(s.style)
		s.help.SetHeaderStyle(s.style)
		s.help.SetSelectedStyle(s.style)
		s.help.SetColumnBorders(false)
		s.help.SetRowBorders(false)
	}
	const numColumns = 3
	splitWidth := width / numColumns
	columns := make([]TableColumn, numColumns)
	for i := 0; i < numColumns; i++ {
		columns[i] = TableColumn{Title: "", Width: splitWidth - 1}
	}
	s.help.SetColumns(columns)
	numRows := getNumRows(len(s.helpItems), numColumns)
	rows := make([][]string, numRows)
	for i := 0; i < numRows; i++ {
		rows[i] = make([]string, numColumns)
		for j := 0; j < numColumns; j++ {
			if i*numColumns+j < len(s.helpItems) {
				rows[i][j] = s.helpItems[i*numColumns+j]
			} else {
				rows[i][j] = ""
			}
		}
	}
	s.help.SetRows(rows)
	return s.help
}

func getNumRows(numItems int, numColumns int) int {
	return (numItems + numColumns - 1) / numColumns
}

func (s *Search) SetHelpVisible(visible bool) *Search {
	s.helpVisible = visible
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

	// inset bounds for inner content
	bounds = bounds.InsetAll(1)
	// layout search and help
	flexLayout := layout.NewVFlex()
	flexSearch := layout.NewFixedChild(3)
	flexHelp := layout.NewFlexChild(10)
	rects := flexLayout.Layout(bounds, []layout.FlexChild{
		flexSearch,
		flexHelp,
	})
	searchBounds := rects[0]
	helpBounds := rects[1]

	// draw search box
	buf.DrawBox(searchBounds.X, searchBounds.Y, searchBounds.Width, searchBounds.Height, s.style)
	searchBounds = searchBounds.InsetAll(1)
	// placeholder or value
	content := s.value
	if content == "" {
		content = s.placeholder
	}
	buf.DrawString(searchBounds.X, searchBounds.Y, content, s.style)

	if s.helpVisible {
		s.help = s.getHelp(helpBounds.X)
		s.help.Render(buf, helpBounds)
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
