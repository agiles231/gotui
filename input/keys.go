package input

// Key represents a keyboard key
type Key int

// Special keys
const (
	KeyNone Key = iota

	// Control keys
	KeyEnter
	KeyTab
	KeyBackspace
	KeyEscape
	KeySpace
	KeyDelete

	// Arrow keys
	KeyUp
	KeyDown
	KeyLeft
	KeyRight

	// Navigation keys
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyInsert

	// Function keys
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12

	// Rune key (regular character)
	KeyRune
)

// Modifier keys
type Modifier int

const (
	ModNone  Modifier = 0
	ModShift Modifier = 1 << iota
	ModAlt
	ModCtrl
	ModMeta
)

// KeyName returns a human-readable name for a key
func KeyName(k Key) string {
	switch k {
	case KeyEnter:
		return "Enter"
	case KeyTab:
		return "Tab"
	case KeyBackspace:
		return "Backspace"
	case KeyEscape:
		return "Escape"
	case KeySpace:
		return "Space"
	case KeyDelete:
		return "Delete"
	case KeyUp:
		return "Up"
	case KeyDown:
		return "Down"
	case KeyLeft:
		return "Left"
	case KeyRight:
		return "Right"
	case KeyHome:
		return "Home"
	case KeyEnd:
		return "End"
	case KeyPageUp:
		return "PageUp"
	case KeyPageDown:
		return "PageDown"
	case KeyInsert:
		return "Insert"
	case KeyF1:
		return "F1"
	case KeyF2:
		return "F2"
	case KeyF3:
		return "F3"
	case KeyF4:
		return "F4"
	case KeyF5:
		return "F5"
	case KeyF6:
		return "F6"
	case KeyF7:
		return "F7"
	case KeyF8:
		return "F8"
	case KeyF9:
		return "F9"
	case KeyF10:
		return "F10"
	case KeyF11:
		return "F11"
	case KeyF12:
		return "F12"
	case KeyRune:
		return "Rune"
	default:
		return "Unknown"
	}
}

// CtrlKey returns the rune for Ctrl+letter (e.g., CtrlKey('c') returns 3)
func CtrlKey(c rune) rune {
	if c >= 'a' && c <= 'z' {
		return c - 'a' + 1
	}
	if c >= 'A' && c <= 'Z' {
		return c - 'A' + 1
	}
	return c
}

// IsCtrl checks if a rune is a control character
func IsCtrl(r rune) bool {
	return r >= 0 && r <= 31
}

// CtrlToLetter converts a control character back to its letter
func CtrlToLetter(r rune) rune {
	if r >= 1 && r <= 26 {
		return 'a' + r - 1
	}
	return r
}

