package widget

import (
	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
)

type WidgetAndLayout struct {
	bounds layout.Rect
	widget Widget
}

type Tab struct {
	BaseWidget
	name string
	widgetAndLayouts []WidgetAndLayout
	focusedWidget    int
}

func (t *Tab) Render(buf *screen.Buffer, bounds layout.Rect) {
	if !t.visible {
		return
	}
	for _, widgetAndLayout := range t.widgetAndLayouts {
		bounds := widgetAndLayout.bounds
		widget := widgetAndLayout.widget
		widget.Render(buf, bounds)
	}
}

// HandleEvent handles input events
func (t *Tab) HandleEvent(event input.Event) bool {
	if t.focusedWidget >= 0 && t.focusedWidget < len(t.widgetAndLayouts) {
		return t.widgetAndLayouts[t.focusedWidget].widget.HandleEvent(event)
	}
	return false
}

// Size returns the preferred size
func (t *Tab) Size() layout.Size {
	return layout.NewSize(80, 24)
}

// MinSize returns the minimum size
func (t *Tab) MinSize() layout.Size {
	return layout.NewSize(1, 1)
}

