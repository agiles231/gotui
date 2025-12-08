package terminal

import "fmt"

// ANSI escape code constants
const (
	// Escape sequence start
	ESC = "\x1b"
	CSI = ESC + "["

	// Screen control
	ClearScreen      = CSI + "2J"
	ClearLine        = CSI + "2K"
	ClearToEndLine   = CSI + "K"
	ClearToStartLine = CSI + "1K"

	// Alternate screen buffer
	AltScreenEnter = CSI + "?1049h"
	AltScreenExit  = CSI + "?1049l"

	// Cursor visibility
	CursorHide = CSI + "?25l"
	CursorShow = CSI + "?25h"

	// Cursor position save/restore
	CursorSave    = CSI + "s"
	CursorRestore = CSI + "u"

	// Cursor movement
	CursorHome = CSI + "H"

	// Text styles
	StyleReset     = CSI + "0m"
	StyleBold      = CSI + "1m"
	StyleDim       = CSI + "2m"
	StyleItalic    = CSI + "3m"
	StyleUnderline = CSI + "4m"
	StyleBlink     = CSI + "5m"
	StyleReverse   = CSI + "7m"
	StyleHidden    = CSI + "8m"
	StyleStrike    = CSI + "9m"

	// Style reset individual
	StyleBoldOff      = CSI + "22m"
	StyleDimOff       = CSI + "22m"
	StyleItalicOff    = CSI + "23m"
	StyleUnderlineOff = CSI + "24m"
	StyleBlinkOff     = CSI + "25m"
	StyleReverseOff   = CSI + "27m"
	StyleHiddenOff    = CSI + "28m"
	StyleStrikeOff    = CSI + "29m"

	// Mouse tracking
	MouseEnable       = CSI + "?1000h"
	MouseDisable      = CSI + "?1000l"
	MouseExtended     = CSI + "?1006h"
	MouseExtendedOff  = CSI + "?1006l"
	MouseAllMotion    = CSI + "?1003h"
	MouseAllMotionOff = CSI + "?1003l"
)

// CursorMove returns the escape sequence to move cursor to (x, y)
// Coordinates are 1-indexed (top-left is 1,1)
func CursorMove(x, y int) string {
	return fmt.Sprintf("%s%d;%dH", CSI, y, x)
}

// CursorUp returns the escape sequence to move cursor up n rows
func CursorUp(n int) string {
	return fmt.Sprintf("%s%dA", CSI, n)
}

// CursorDown returns the escape sequence to move cursor down n rows
func CursorDown(n int) string {
	return fmt.Sprintf("%s%dB", CSI, n)
}

// CursorForward returns the escape sequence to move cursor right n columns
func CursorForward(n int) string {
	return fmt.Sprintf("%s%dC", CSI, n)
}

// CursorBack returns the escape sequence to move cursor left n columns
func CursorBack(n int) string {
	return fmt.Sprintf("%s%dD", CSI, n)
}

// CursorNextLine moves cursor to beginning of line n lines down
func CursorNextLine(n int) string {
	return fmt.Sprintf("%s%dE", CSI, n)
}

// CursorPrevLine moves cursor to beginning of line n lines up
func CursorPrevLine(n int) string {
	return fmt.Sprintf("%s%dF", CSI, n)
}

// CursorColumn moves cursor to column n
func CursorColumn(n int) string {
	return fmt.Sprintf("%s%dG", CSI, n)
}

// ScrollUp scrolls the screen up n lines
func ScrollUp(n int) string {
	return fmt.Sprintf("%s%dS", CSI, n)
}

// ScrollDown scrolls the screen down n lines
func ScrollDown(n int) string {
	return fmt.Sprintf("%s%dT", CSI, n)
}

// SetScrollRegion sets the scroll region from top to bottom
func SetScrollRegion(top, bottom int) string {
	return fmt.Sprintf("%s%d;%dr", CSI, top, bottom)
}

// ResetScrollRegion resets the scroll region to the full screen
func ResetScrollRegion() string {
	return CSI + "r"
}

