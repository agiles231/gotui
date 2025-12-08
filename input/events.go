package input

// EventType represents the type of an event
type EventType int

const (
	EventKey EventType = iota
	EventResize
	EventMouse
	EventError
	EventQuit
)

// Event is the interface for all events
type Event interface {
	Type() EventType
}

// KeyEvent represents a keyboard event
type KeyEvent struct {
	Key      Key
	Rune     rune
	Modifier Modifier
}

func (e KeyEvent) Type() EventType {
	return EventKey
}

// IsCtrl checks if the Ctrl modifier is pressed
func (e KeyEvent) IsCtrl() bool {
	return e.Modifier&ModCtrl != 0
}

// IsAlt checks if the Alt modifier is pressed
func (e KeyEvent) IsAlt() bool {
	return e.Modifier&ModAlt != 0
}

// IsShift checks if the Shift modifier is pressed
func (e KeyEvent) IsShift() bool {
	return e.Modifier&ModShift != 0
}

// Matches checks if the event matches a key and modifiers
func (e KeyEvent) Matches(key Key, mod Modifier) bool {
	return e.Key == key && e.Modifier == mod
}

// MatchesRune checks if the event matches a rune key
func (e KeyEvent) MatchesRune(r rune) bool {
	return e.Key == KeyRune && e.Rune == r
}

// ResizeEvent represents a terminal resize event
type ResizeEvent struct {
	Width  int
	Height int
}

func (e ResizeEvent) Type() EventType {
	return EventResize
}

// MouseButton represents mouse button states
type MouseButton int

const (
	MouseNone MouseButton = iota
	MouseLeft
	MouseMiddle
	MouseRight
	MouseWheelUp
	MouseWheelDown
	MouseRelease
)

// MouseEvent represents a mouse event
type MouseEvent struct {
	X      int
	Y      int
	Button MouseButton
	Mod    Modifier
}

func (e MouseEvent) Type() EventType {
	return EventMouse
}

// ErrorEvent represents an error that occurred during input handling
type ErrorEvent struct {
	Err error
}

func (e ErrorEvent) Type() EventType {
	return EventError
}

// QuitEvent represents a quit signal
type QuitEvent struct{}

func (e QuitEvent) Type() EventType {
	return EventQuit
}

