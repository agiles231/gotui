package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/agiles231/gotui/app"
	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
	"github.com/agiles231/gotui/terminal"
	"github.com/agiles231/gotui/widget"
)

func main() {
	// Create demo app
	demo := NewDemoApp()

	if err := demo.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// DemoApp demonstrates all the TUI components
type DemoApp struct {
	app         *app.App
	tabs        []string
	currentTab  int
	list        *widget.List
	table       *widget.Table
	searchAndResults *widget.SearchAndResults
	progress    *widget.Progress
	spinner     *widget.Spinner
	form        *widget.Form
	textInput   *widget.TextInput
	statusText  string
	// Z-Order demo state
	zOrderRedInFront bool
}

// NewDemoApp creates the demo application
func NewDemoApp() *DemoApp {
	d := &DemoApp{
		app:        app.New(),
		tabs:       []string{"List", "Table", "Search and Results", "Progress", "Form", "Z-Order", "About"},
		currentTab: 0,
		statusText: "Press ←/→ to switch panels, Ctrl+Q to quit",
		zOrderRedInFront: true,
	}

	// Create list widget
	d.list = widget.NewList().
		SetStrings([]string{
			"Item 1 - First item",
			"Item 2 - Second item",
			"Item 3 - Third item",
			"Item 4 - Fourth item",
			"Item 5 - Fifth item",
			"Item 6 - Sixth item",
			"Item 7 - Seventh item",
			"Item 8 - Eighth item",
		}).
		SetHeight(8).
		OnSelect(func(i int, item widget.ListItem) {
			d.statusText = fmt.Sprintf("Selected: %s", item.Text)
		})

	// Create table widget
	d.table = widget.NewTable().
		SetColumns([]widget.TableColumn{
			{Title: "ID", Width: 5},
			{Title: "Name", Flex: 2},
			{Title: "Status", Flex: 1},
		}).
		SetRows([][]string{
			{"1", "Alice", "Active"},
			{"2", "Bob", "Inactive"},
			{"3", "Charlie", "Active"},
			{"4", "Diana", "Pending"},
			{"5", "Eve", "Active"},
		}).
		SetHeight(8).
		SetShowHeader(true)

	// Create search and results widget
	searchResultsTable := widget.NewTable().
		SetColumns([]widget.TableColumn{
			{Title: "ID", Width: 5},
			{Title: "Name", Flex: 2},
			{Title: "Status", Flex: 1},
		}).
		SetRows([][]string{
			{"1", "Alice", "Active"},
			{"2", "Bob", "Inactive"},
			{"3", "Charlie", "Active"},
			{"4", "Diana", "Pending"},
			{"5", "Eve", "Active"},
		})
	search := widget.NewSearch().SetPlaceholder("Search")
	search.SetOnSubmit(func(value string) {
		rows := d.table.Rows()
		resultRows := [][]string{}
		for _, row := range rows {
			if strings.Contains(strings.ToLower(row[1]), strings.ToLower(value)) {
				resultRows = append(resultRows, row)
			}
		}
		searchResultsTable.SetRows(resultRows)
	}).
	SetHelpItems([]string{
		"/items : Search items",
		"/item-specs : Search item specifications",
		"/storage : Search storage",
		"/storage-specs : Search storage specs",
		"/tags: Search tags",
	})
	d.searchAndResults = widget.NewSearchAndResults()
	d.searchAndResults.SetSearch(search)
	d.searchAndResults.SetTable(searchResultsTable)

	// Create progress widget
	d.progress = widget.NewProgress().
		SetWidth(30).
		SetValue(0.0)

	// Create spinner
	d.spinner = widget.NewSpinner().
		SetLabel("Loading...")

	// Create form widget
	d.form = widget.NewForm().
		SetTitle("User Registration").
		SetShowBorder(true).
		SetLabelWidth(12)

	d.textInput = d.form.AddTextInput("Username", "Enter username")
	d.form.AddPasswordInput("Password", "Enter password")
	d.form.AddTextInput("Email", "user@example.com")

	// Setup app
	d.app.
		SetRoot(d).
		SetFPS(30).
		OnTick(100*time.Millisecond, func(a *app.App, t time.Time) bool {
			// Animate progress and spinner
			val := d.progress.Value() + 0.01
			if val > 1 {
				val = 0
			}
			d.progress.SetValue(val)
			d.spinner.Advance()
			return d.currentTab == 3 // Only render if on progress tab
		})

	// Set initial focus
	d.list.SetFocused(true)

	return d
}

// Run starts the demo app
func (d *DemoApp) Run() error {
	return d.app.Run()
}

// Render implements widget.Widget
func (d *DemoApp) Render(buf *screen.Buffer, bounds layout.Rect) {
	style := terminal.DefaultStyle()
	titleStyle := style.WithBold().WithFG(terminal.ColorCyan)
	tabStyle := style.WithDim()
	activeTabStyle := style.WithBold().WithReverse()
	statusStyle := style.WithReverse()
	borderStyle := style.WithFG(terminal.ColorBlue)

	// Draw title
	title := "╔═══ GoTUI Demo ═══╗"
	titleX := (bounds.Width - len(title)) / 2
	buf.DrawString(titleX, bounds.Y, bounds.Z, title, titleStyle)

	// Draw tabs
	tabY := bounds.Y + 2
	tabX := 2
	for i, tab := range d.tabs {
		tabText := fmt.Sprintf(" %s ", tab)
		style := tabStyle
		if i == d.currentTab {
			style = activeTabStyle
		}
		buf.DrawString(tabX, tabY, bounds.Z, tabText, style)
		tabX += len(tabText) + 1
	}

	// Draw content area border
	contentY := tabY + 2
	contentHeight := bounds.Height - contentY - 2
	buf.DrawBox(1, contentY, bounds.Z, bounds.Width-2, contentHeight, borderStyle)

	// Draw content
	contentBounds := layout.NewRect(2, contentY+1, bounds.Z, bounds.Width-4, contentHeight-2)
	switch d.currentTab {
	case 0:
		d.renderListTab(buf, contentBounds)
	case 1:
		d.renderTableTab(buf, contentBounds)
	case 2:
		d.renderSearchAndResultsTab(buf, contentBounds)
	case 3:
		d.renderProgressTab(buf, contentBounds)
	case 4:
		d.renderFormTab(buf, contentBounds)
	case 5:
		d.renderZOrderTab(buf, contentBounds)
	case 6:
		d.renderAboutTab(buf, contentBounds)
	}

	// Draw status bar
	statusY := bounds.Y + bounds.Height - 1
	for x := 0; x < bounds.Width; x++ {
		buf.Set(x, statusY, bounds.Z, screen.NewCell(' ', statusStyle))
	}
	buf.DrawString(1, statusY, bounds.Z, d.statusText, statusStyle)

	// Draw help on right side of status bar
	help := "←/→: Switch tabs | Ctrl+Q: Quit"
	if d.currentTab == 5 {
		help = "SPACE: Swap z-order | ←/→: Switch tabs | Ctrl+Q: Quit"
	}
	buf.DrawString(bounds.Width-len(help)-1, statusY, bounds.Z, help, statusStyle)
}

func (d *DemoApp) renderListTab(buf *screen.Buffer, bounds layout.Rect) {
	style := terminal.DefaultStyle()

	// Title
	buf.DrawString(bounds.X, bounds.Y, bounds.Z, "Selectable List:", style.WithBold())

	// List
	listBounds := layout.NewRect(bounds.X, bounds.Y+2, bounds.Z, bounds.Width/2, bounds.Height-2)
	d.list.Render(buf, listBounds)

	// Instructions
	instX := bounds.X + bounds.Width/2 + 2
	buf.DrawString(instX, bounds.Y+2, bounds.Z, "Controls:", style.WithBold())
	buf.DrawString(instX, bounds.Y+4, bounds.Z, "↑/↓  Navigate items", style)
	buf.DrawString(instX, bounds.Y+5, bounds.Z, "Enter  Select item", style)
	buf.DrawString(instX, bounds.Y+6, bounds.Z, "PgUp/PgDn  Page navigation", style)
}

func (d *DemoApp) renderTableTab(buf *screen.Buffer, bounds layout.Rect) {
	style := terminal.DefaultStyle()

	// Title
	buf.DrawString(bounds.X, bounds.Y, bounds.Z, "Data Table:", style.WithBold())

	// Table
	tableBounds := layout.NewRect(bounds.X, bounds.Y+2, bounds.Z, bounds.Width, bounds.Height-2)
	d.table.Render(buf, tableBounds)
}

func (d *DemoApp) renderSearchAndResultsTab(buf *screen.Buffer, bounds layout.Rect) {
	style := terminal.DefaultStyle()

	// Title
	buf.DrawString(bounds.X, bounds.Y, bounds.Z, "Search and Results:", style.WithBold())
	searchAndResultsBounds := layout.NewRect(bounds.X, bounds.Y+2, bounds.Z, bounds.Width, bounds.Height-2)
	d.searchAndResults.Render(buf, searchAndResultsBounds)
}

func (d *DemoApp) renderProgressTab(buf *screen.Buffer, bounds layout.Rect) {
	style := terminal.DefaultStyle()

	// Title
	buf.DrawString(bounds.X, bounds.Y, bounds.Z, "Progress Indicators:", style.WithBold())

	// Progress bar
	buf.DrawString(bounds.X, bounds.Y+3, bounds.Z, "Progress Bar:", style)
	progressBounds := layout.NewRect(bounds.X, bounds.Y+4, bounds.Z, bounds.Width, 1)
	d.progress.Render(buf, progressBounds)

	// Spinner
	buf.DrawString(bounds.X, bounds.Y+7, bounds.Z, "Spinner:", style)
	spinnerBounds := layout.NewRect(bounds.X, bounds.Y+8, bounds.Z, bounds.Width, 1)
	d.spinner.Render(buf, spinnerBounds)

	// Static progress bars at different values
	buf.DrawString(bounds.X, bounds.Y+11, bounds.Z, "Static Progress:", style)
	
	staticProgress := widget.NewProgress().SetWidth(25)
	for i, val := range []float64{0.25, 0.50, 0.75, 1.0} {
		staticProgress.SetValue(val)
		staticProgress.Render(buf, layout.NewRect(bounds.X, bounds.Y+12+i, bounds.Z, bounds.Width, 1))
	}
}

func (d *DemoApp) renderFormTab(buf *screen.Buffer, bounds layout.Rect) {
	formBounds := layout.NewRect(bounds.X, bounds.Y, bounds.Z, 50, 8)
	d.form.Render(buf, formBounds)

	// Help text
	style := terminal.DefaultStyle()
	buf.DrawString(bounds.X, bounds.Y+10, bounds.Z, "Press Tab to move between fields", style.WithDim())
	buf.DrawString(bounds.X, bounds.Y+11, bounds.Z, "Press Ctrl+Enter to submit", style.WithDim())
}

func (d *DemoApp) renderZOrderTab(buf *screen.Buffer, bounds layout.Rect) {
	style := terminal.DefaultStyle()
	
	// Title
	buf.DrawString(bounds.X, bounds.Y, bounds.Z, "Z-Order Demo:", style.WithBold())
	buf.DrawString(bounds.X, bounds.Y+1, bounds.Z, "Press SPACE to swap which box is in front", style.WithDim())

	// Define two overlapping boxes
	// Red box (left side, partially overlapped)
	redStyle := terminal.DefaultStyle().WithFG(terminal.ColorRed).WithBold()
	redBgStyle := terminal.DefaultStyle().WithBG(terminal.ColorRed).WithFG(terminal.ColorWhite)
	
	// Blue box (right side, partially overlapped)
	blueStyle := terminal.DefaultStyle().WithFG(terminal.ColorBlue).WithBold()
	blueBgStyle := terminal.DefaultStyle().WithBG(terminal.ColorBlue).WithFG(terminal.ColorWhite)

	// Box dimensions
	boxWidth := 25
	boxHeight := 10
	
	// Positions - overlapping by about 10 characters
	redX := bounds.X + 2
	redY := bounds.Y + 3
	blueX := bounds.X + 17  // Overlaps red box
	blueY := bounds.Y + 5   // Slightly lower

	// Determine z-order based on state
	var redZ, blueZ int
	if d.zOrderRedInFront {
		redZ = 1
		blueZ = 0
	} else {
		redZ = 0
		blueZ = 1
	}

	// Draw RED box
	buf.DrawBox(redX, redY, redZ, boxWidth, boxHeight, redStyle)
	// Fill interior
	for y := 1; y < boxHeight-1; y++ {
		for x := 1; x < boxWidth-1; x++ {
			buf.Set(redX+x, redY+y, redZ, screen.NewCell(' ', redBgStyle))
		}
	}
	// Draw label
	redLabel := "RED BOX"
	buf.DrawString(redX+(boxWidth-len(redLabel))/2, redY+boxHeight/2, redZ, redLabel, redBgStyle.WithBold())
	if d.zOrderRedInFront {
		buf.DrawString(redX+(boxWidth-len("(FRONT)"))/2, redY+boxHeight/2+1, redZ, "(FRONT)", redBgStyle)
	} else {
		buf.DrawString(redX+(boxWidth-len("(BACK)"))/2, redY+boxHeight/2+1, redZ, "(BACK)", redBgStyle)
	}

	// Draw BLUE box
	buf.DrawBox(blueX, blueY, blueZ, boxWidth, boxHeight, blueStyle)
	// Fill interior
	for y := 1; y < boxHeight-1; y++ {
		for x := 1; x < boxWidth-1; x++ {
			buf.Set(blueX+x, blueY+y, blueZ, screen.NewCell(' ', blueBgStyle))
		}
	}
	// Draw label
	blueLabel := "BLUE BOX"
	buf.DrawString(blueX+(boxWidth-len(blueLabel))/2, blueY+boxHeight/2, blueZ, blueLabel, blueBgStyle.WithBold())
	if !d.zOrderRedInFront {
		buf.DrawString(blueX+(boxWidth-len("(FRONT)"))/2, blueY+boxHeight/2+1, blueZ, "(FRONT)", blueBgStyle)
	} else {
		buf.DrawString(blueX+(boxWidth-len("(BACK)"))/2, blueY+boxHeight/2+1, blueZ, "(BACK)", blueBgStyle)
	}

	// Status indicator
	statusY := bounds.Y + bounds.Height - 2
	if d.zOrderRedInFront {
		buf.DrawString(bounds.X, statusY, bounds.Z, "Current: RED is in front (z=1), BLUE is behind (z=0)", style)
	} else {
		buf.DrawString(bounds.X, statusY, bounds.Z, "Current: BLUE is in front (z=1), RED is behind (z=0)", style)
	}
}

func (d *DemoApp) renderAboutTab(buf *screen.Buffer, bounds layout.Rect) {
	style := terminal.DefaultStyle()
	highlightStyle := style.WithFG(terminal.ColorGreen)

	lines := []string{
		"GoTUI - A Terminal User Interface Framework",
		"",
		"Built from scratch using ANSI escape codes",
		"",
		"Features:",
		"  • Raw terminal mode handling",
		"  • Keyboard input with escape sequence parsing",
		"  • Double-buffered rendering with diff updates",
		"  • Flexbox-inspired layout system",
		"  • Widget library:",
		"    - Text labels",
		"    - Text input fields",
		"    - Lists with scrolling",
		"    - Tables with columns",
		"    - Progress bars & spinners",
		"    - Menus",
		"    - Forms",
		"",
		"Created with ❤️ in Go",
	}

	for i, line := range lines {
		if i >= bounds.Height {
			break
		}
		s := style
		if i == 0 {
			s = highlightStyle.WithBold()
		}
		buf.DrawString(bounds.X, bounds.Y+i, bounds.Z, line, s)
	}
}

// HandleEvent implements widget.Widget
func (d *DemoApp) HandleEvent(event input.Event) bool {
	keyEvent, ok := event.(input.KeyEvent)
	if !ok {
		return false
	}

	// Handle tab switching with left/right arrows
	if keyEvent.Key == input.KeyLeft {
		d.currentTab--
		if d.currentTab < 0 {
			d.currentTab = len(d.tabs) - 1
		}
		d.updateFocus()
		return true
	}
	if keyEvent.Key == input.KeyRight {
		d.currentTab++
		if d.currentTab >= len(d.tabs) {
			d.currentTab = 0
		}
		d.updateFocus()
		return true
	}

	// Pass event to current tab's widget
	switch d.currentTab {
	case 0:
		return d.list.HandleEvent(event)
	case 1:
		return d.table.HandleEvent(event)
	case 2:
		return d.searchAndResults.HandleEvent(event)
	case 4:
		return d.form.HandleEvent(event)
	case 5:
		// Z-Order tab - handle spacebar to swap z-order
		if keyEvent.Key == input.KeyRune && keyEvent.Rune == ' ' {
			d.zOrderRedInFront = !d.zOrderRedInFront
			return true
		}
	}

	return false
}

func (d *DemoApp) updateFocus() {
	// Unfocus all
	d.list.SetFocused(false)
	d.table.SetFocused(false)
	d.searchAndResults.SetFocused(false)
	d.form.SetFocused(false)

	// Focus current tab's widget
	switch d.currentTab {
	case 0:
		d.list.SetFocused(true)
	case 1:
		d.table.SetFocused(true)
	case 2:
		d.searchAndResults.SetFocused(true)
	case 4:
		d.form.SetFocused(true)
	}
}

// Size implements widget.Widget
func (d *DemoApp) Size() layout.Size {
	return layout.NewSize(80, 24)
}

// MinSize implements widget.Widget
func (d *DemoApp) MinSize() layout.Size {
	return layout.NewSize(60, 20)
}

// SetFocused implements widget.Widget
func (d *DemoApp) SetFocused(focused bool) {}

// IsFocused implements widget.Widget
func (d *DemoApp) IsFocused() bool { return true }

// IsInteractive implements widget.Widget
func (d *DemoApp) IsInteractive() bool { return true }
