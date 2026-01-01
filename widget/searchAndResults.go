package widget

import (
	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
)

type SearchAndResults struct {
	BaseWidget
	search *Search
	results *Table
	searchFocused bool
	resultsFocused bool
}

func NewSearchAndResults() *SearchAndResults {
	s := &SearchAndResults{
		BaseWidget: NewBaseWidget(),
		search: NewSearch(),
		results: NewTable(),
	}
	s.interactive = true
	return s
}

func (s *SearchAndResults) SetSearch(search *Search) *SearchAndResults {
	s.search = search
	return s
}

func (s *SearchAndResults) SetTable(table *Table) *SearchAndResults {
	s.results = table
	return s
}


func (s *SearchAndResults) Render(buf *screen.Buffer, bounds layout.Rect) {
	search_height := 12
	vFlex := layout.NewVFlex().WithGap(2)
	searchFlex := layout.NewFixedChild(search_height)
	resultsFlex := layout.NewFlexChild(10)
	rects := vFlex.Layout(bounds, []layout.FlexChild{
		searchFlex,
		resultsFlex,
	})
	search_bounds := rects[0]
	results_bounds := rects[1]
	s.search.Render(buf, search_bounds.InsetAll(1))
	s.results.Render(buf, results_bounds.InsetAll(1))
}

func (s *SearchAndResults) HandleEvent(event input.Event) bool {
	return s.search.HandleEvent(event) || s.results.HandleEvent(event)
}

func (s *SearchAndResults) Size() layout.Size {
	search_size := s.search.Size()
	results_size := s.results.Size()
	return layout.NewSize(search_size.Width + results_size.Width, search_size.Height + results_size.Height)
}

func (s *SearchAndResults) MinSize() layout.Size {
	search_min_size := s.search.MinSize()
	results_min_size := s.results.MinSize()
	return layout.NewSize(search_min_size.Width + results_min_size.Width, search_min_size.Height + results_min_size.Height)
}

func (s *SearchAndResults) SetFocused(focused bool) {
	s.search.SetFocused(focused)
	s.results.SetFocused(focused)
}

func (s *SearchAndResults) IsFocused() bool {
	return s.search.IsFocused() || s.results.IsFocused()
}

func (s *SearchAndResults) IsInteractive() bool {
	return s.search.IsInteractive() || s.results.IsInteractive()
}
