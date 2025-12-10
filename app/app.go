package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agiles231/gotui/input"
	"github.com/agiles231/gotui/layout"
	"github.com/agiles231/gotui/screen"
	"github.com/agiles231/gotui/terminal"
	"github.com/agiles231/gotui/widget"
)

// App represents the main application
type App struct {
	terminal     *terminal.Terminal
	screen       *screen.Screen
	inputReader  *input.Reader
	root         widget.Widget
	focusManager *widget.FocusManager
	running      bool
	quitChan     chan struct{}
	renderChan   chan struct{}
	fps          int
	onInit       func(*App)
	onQuit       func(*App)
	onResize     func(*App, int, int)
	onTick       func(*App, time.Time) bool
	tickInterval time.Duration
}

// New creates a new application
func New() *App {
	return &App{
		terminal:     terminal.New(),
		focusManager: widget.NewFocusManager(),
		quitChan:     make(chan struct{}),
		renderChan:   make(chan struct{}, 1),
		fps:          60,
	}
}

// SetRoot sets the root widget
func (a *App) SetRoot(w widget.Widget) *App {
	a.root = w
	return a
}

// Root returns the root widget
func (a *App) Root() widget.Widget {
	return a.root
}

// SetFPS sets the target frames per second
func (a *App) SetFPS(fps int) *App {
	a.fps = fps
	return a
}

// OnInit sets the initialization callback
func (a *App) OnInit(fn func(*App)) *App {
	a.onInit = fn
	return a
}

// OnQuit sets the quit callback
func (a *App) OnQuit(fn func(*App)) *App {
	a.onQuit = fn
	return a
}

// OnResize sets the resize callback
func (a *App) OnResize(fn func(*App, int, int)) *App {
	a.onResize = fn
	return a
}

// OnTick sets a periodic tick callback
// Return true to request a render
func (a *App) OnTick(interval time.Duration, fn func(*App, time.Time) bool) *App {
	a.tickInterval = interval
	a.onTick = fn
	return a
}

// FocusManager returns the focus manager
func (a *App) FocusManager() *widget.FocusManager {
	return a.focusManager
}

// Screen returns the screen
func (a *App) Screen() *screen.Screen {
	return a.screen
}

// Terminal returns the terminal
func (a *App) Terminal() *terminal.Terminal {
	return a.terminal
}

// Width returns the terminal width
func (a *App) Width() int {
	if a.screen != nil {
		return a.screen.Width()
	}
	return 0
}

// Height returns the terminal height
func (a *App) Height() int {
	if a.screen != nil {
		return a.screen.Height()
	}
	return 0
}

// Quit signals the application to quit
func (a *App) Quit() {
	if a.running {
		close(a.quitChan)
	}
}

// RequestRender requests a screen render
func (a *App) RequestRender() {
	select {
	case a.renderChan <- struct{}{}:
	default:
	}
}

// Run starts the application event loop
func (a *App) Run() error {
	// Enter raw mode
	if err := a.terminal.EnterRawMode(); err != nil {
		return err
	}
	defer a.terminal.ExitRawMode()

	// Enter alternate screen
	a.terminal.EnterAltScreen()
	defer a.terminal.ExitAltScreen()

	// Hide cursor
	a.terminal.HideCursor()
	defer a.terminal.ShowCursor()

	// Create screen
	var err error
	a.screen, err = screen.NewScreen(a.terminal)
	if err != nil {
		return err
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGWINCH)
	defer signal.Stop(sigChan)

	// Start input reader
	a.inputReader = input.NewReader()
	a.inputReader.Start()
	defer a.inputReader.Stop()

	a.running = true
	defer func() { a.running = false }()

	// Call init callback
	if a.onInit != nil {
		a.onInit(a)
	}

	// Initial render
	a.render()

	// Calculate frame duration
	frameDuration := time.Second / time.Duration(a.fps)

	// Setup tick timer if needed
	var tickChan <-chan time.Time
	if a.onTick != nil && a.tickInterval > 0 {
		ticker := time.NewTicker(a.tickInterval)
		defer ticker.Stop()
		tickChan = ticker.C
	}

	// Main event loop
	for {
		select {
		case <-a.quitChan:
			if a.onQuit != nil {
				a.onQuit(a)
			}
			return nil

		case sig := <-sigChan:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				if a.onQuit != nil {
					a.onQuit(a)
				}
				return nil
			case syscall.SIGWINCH:
				a.handleResize()
			}

		case event := <-a.inputReader.Events():
			if a.handleEvent(event) {
				a.render()
			}

		case <-a.renderChan:
			a.render()

		case t := <-tickChan:
			if a.onTick != nil && a.onTick(a, t) {
				a.render()
			}

		case <-time.After(frameDuration):
			// Idle - no events
		}
	}
}

// handleResize handles terminal resize
func (a *App) handleResize() {
	width, height, err := a.terminal.Size()
	if err != nil {
		return
	}

	a.screen.Resize(width, height)

	if a.onResize != nil {
		a.onResize(a, width, height)
	}

	a.forceRender()
}

// handleEvent processes an input event
func (a *App) handleEvent(event input.Event) bool {
	// Handle quit keys (Ctrl+C, Ctrl+Q)
	if keyEvent, ok := event.(input.KeyEvent); ok {
		if keyEvent.IsCtrl() && keyEvent.Key == input.KeyRune {
			switch keyEvent.Rune {
			case 'c', 'q':
				a.Quit()
				return false
			}
		}
	}

	// Pass to root widget
	if a.root != nil {
		return a.root.HandleEvent(event)
	}

	return false
}

// render renders the application
func (a *App) render() {
	if a.screen == nil || a.root == nil {
		return
	}

	// Clear screen
	a.screen.Clear()

	// Render root widget
	bounds := layout.NewRect(0, 0, a.screen.Width(), a.screen.Height())
	a.root.Render(a.screen.Buffer(), bounds)

	// Render to terminal
	a.screen.Render()
	a.screen.Flush()
}

// forceRender renders the entire screen regardless of changes
// Used after resize to clear artifacts from the previous terminal size
func (a *App) forceRender() {
	if a.screen == nil || a.root == nil {
		return
	}

	// Clear screen
	a.screen.Clear()

	// Render root widget
	bounds := layout.NewRect(0, 0, a.screen.Width(), a.screen.Height())
	a.root.Render(a.screen.Buffer(), bounds)

	// Force render all cells to terminal
	a.screen.ForceRender()
	a.screen.Flush()
}

// SimpleApp provides a simpler API for basic applications
type SimpleApp struct {
	*App
	title      string
	statusBar  string
	content    widget.Widget
}

// NewSimple creates a simple application with title and status bar
func NewSimple(title string) *SimpleApp {
	return &SimpleApp{
		App:   New(),
		title: title,
	}
}

// SetContent sets the main content widget
func (s *SimpleApp) SetContent(w widget.Widget) *SimpleApp {
	s.content = w
	s.App.SetRoot(&simpleLayout{
		title:     s.title,
		statusBar: &s.statusBar,
		content:   w,
	})
	return s
}

// SetStatus sets the status bar text
func (s *SimpleApp) SetStatus(status string) {
	s.statusBar = status
	s.RequestRender()
}

// simpleLayout wraps content with title and status bar
type simpleLayout struct {
	widget.BaseWidget
	title     string
	statusBar *string
	content   widget.Widget
}

func (l *simpleLayout) Render(buf *screen.Buffer, bounds layout.Rect) {
	style := terminal.DefaultStyle()
	titleStyle := style.WithBold().WithReverse()
	statusStyle := style.WithReverse()

	// Draw title bar
	for x := 0; x < bounds.Width; x++ {
		buf.Set(bounds.X+x, bounds.Y, screen.NewCell(' ', titleStyle))
	}
	titleX := (bounds.Width - len(l.title)) / 2
	buf.DrawString(bounds.X+titleX, bounds.Y, l.title, titleStyle)

	// Draw status bar
	statusY := bounds.Y + bounds.Height - 1
	for x := 0; x < bounds.Width; x++ {
		buf.Set(bounds.X+x, statusY, screen.NewCell(' ', statusStyle))
	}
	if l.statusBar != nil {
		buf.DrawString(bounds.X+1, statusY, *l.statusBar, statusStyle)
	}

	// Draw content
	if l.content != nil {
		contentBounds := layout.NewRect(
			bounds.X,
			bounds.Y+1,
			bounds.Width,
			bounds.Height-2,
		)
		l.content.Render(buf, contentBounds)
	}
}

func (l *simpleLayout) HandleEvent(event input.Event) bool {
	if l.content != nil {
		return l.content.HandleEvent(event)
	}
	return false
}

func (l *simpleLayout) Size() layout.Size {
	return layout.NewSize(80, 24)
}

func (l *simpleLayout) MinSize() layout.Size {
	return layout.NewSize(40, 10)
}

