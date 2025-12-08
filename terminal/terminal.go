package terminal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sys/unix"
)

// Terminal represents a terminal instance with raw mode support
type Terminal struct {
	fd          int
	originalTTY *unix.Termios
	width       int
	height      int
	inAltScreen bool
}

// New creates a new Terminal instance
func New() *Terminal {
	return &Terminal{
		fd: int(os.Stdin.Fd()),
	}
}

// EnterRawMode puts the terminal into raw mode
func (t *Terminal) EnterRawMode() error {
	termios, err := unix.IoctlGetTermios(t.fd, unix.TCGETS)
	if err != nil {
		return fmt.Errorf("failed to get termios: %w", err)
	}

	// Save original settings
	t.originalTTY = termios

	// Create a copy for raw mode
	raw := *termios

	// Input flags: disable break, CR to NL, parity check, strip, XON/XOFF
	raw.Iflag &^= unix.BRKINT | unix.ICRNL | unix.INPCK | unix.ISTRIP | unix.IXON

	// Output flags: disable post-processing
	raw.Oflag &^= unix.OPOST

	// Control flags: set 8-bit chars
	raw.Cflag |= unix.CS8

	// Local flags: disable echo, canonical mode, signals, extended processing
	raw.Lflag &^= unix.ECHO | unix.ICANON | unix.IEXTEN | unix.ISIG

	// Control chars: set minimum input to 1 byte, no timeout
	raw.Cc[unix.VMIN] = 1
	raw.Cc[unix.VTIME] = 0

	if err := unix.IoctlSetTermios(t.fd, unix.TCSETS, &raw); err != nil {
		return fmt.Errorf("failed to set raw mode: %w", err)
	}

	return nil
}

// ExitRawMode restores the terminal to its original state
func (t *Terminal) ExitRawMode() error {
	if t.originalTTY == nil {
		return nil
	}

	if err := unix.IoctlSetTermios(t.fd, unix.TCSETS, t.originalTTY); err != nil {
		return fmt.Errorf("failed to restore terminal: %w", err)
	}

	return nil
}

// EnterAltScreen switches to the alternate screen buffer
func (t *Terminal) EnterAltScreen() {
	fmt.Print(AltScreenEnter)
	t.inAltScreen = true
}

// ExitAltScreen returns to the main screen buffer
func (t *Terminal) ExitAltScreen() {
	fmt.Print(AltScreenExit)
	t.inAltScreen = false
}

// Size returns the current terminal size (width, height)
func (t *Terminal) Size() (int, int, error) {
	ws, err := unix.IoctlGetWinsize(t.fd, unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get window size: %w", err)
	}

	t.width = int(ws.Col)
	t.height = int(ws.Row)

	return t.width, t.height, nil
}

// Width returns the cached terminal width
func (t *Terminal) Width() int {
	return t.width
}

// Height returns the cached terminal height
func (t *Terminal) Height() int {
	return t.height
}

// Clear clears the entire screen
func (t *Terminal) Clear() {
	fmt.Print(ClearScreen)
}

// HideCursor hides the cursor
func (t *Terminal) HideCursor() {
	fmt.Print(CursorHide)
}

// ShowCursor shows the cursor
func (t *Terminal) ShowCursor() {
	fmt.Print(CursorShow)
}

// MoveCursor moves the cursor to the specified position (1-indexed)
func (t *Terminal) MoveCursor(x, y int) {
	fmt.Print(CursorMove(x, y))
}

// Flush ensures all output is written
func (t *Terminal) Flush() {
	os.Stdout.Sync()
}

// SetupResizeHandler sets up a handler for terminal resize signals
func (t *Terminal) SetupResizeHandler(callback func(width, height int)) chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)

	go func() {
		for range sigChan {
			if w, h, err := t.Size(); err == nil {
				callback(w, h)
			}
		}
	}()

	return sigChan
}

// StopResizeHandler stops the resize signal handler
func (t *Terminal) StopResizeHandler(sigChan chan os.Signal) {
	signal.Stop(sigChan)
	close(sigChan)
}

