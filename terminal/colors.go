package terminal

import "fmt"

// Color represents a terminal color
type Color interface {
	// FG returns the ANSI escape sequence for foreground color
	FG() string
	// BG returns the ANSI escape sequence for background color
	BG() string
}

// Basic 16 colors (0-15)
type BasicColor uint8

const (
	ColorBlack BasicColor = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorBrightBlack
	ColorBrightRed
	ColorBrightGreen
	ColorBrightYellow
	ColorBrightBlue
	ColorBrightMagenta
	ColorBrightCyan
	ColorBrightWhite
)

// Default color (terminal default)
const ColorDefault BasicColor = 255

func (c BasicColor) FG() string {
	if c == ColorDefault {
		return CSI + "39m"
	}
	if c < 8 {
		return fmt.Sprintf("%s%dm", CSI, 30+c)
	}
	return fmt.Sprintf("%s%dm", CSI, 90+(c-8))
}

func (c BasicColor) BG() string {
	if c == ColorDefault {
		return CSI + "49m"
	}
	if c < 8 {
		return fmt.Sprintf("%s%dm", CSI, 40+c)
	}
	return fmt.Sprintf("%s%dm", CSI, 100+(c-8))
}

// Color256 represents a 256-color palette color (0-255)
type Color256 uint8

func (c Color256) FG() string {
	return fmt.Sprintf("%s38;5;%dm", CSI, c)
}

func (c Color256) BG() string {
	return fmt.Sprintf("%s48;5;%dm", CSI, c)
}

// RGB represents a true color (24-bit)
type RGB struct {
	R, G, B uint8
}

// NewRGB creates a new RGB color
func NewRGB(r, g, b uint8) RGB {
	return RGB{R: r, G: g, B: b}
}

// Hex creates an RGB color from a hex value (e.g., 0xFF5733)
func Hex(hex uint32) RGB {
	return RGB{
		R: uint8((hex >> 16) & 0xFF),
		G: uint8((hex >> 8) & 0xFF),
		B: uint8(hex & 0xFF),
	}
}

func (c RGB) FG() string {
	return fmt.Sprintf("%s38;2;%d;%d;%dm", CSI, c.R, c.G, c.B)
}

func (c RGB) BG() string {
	return fmt.Sprintf("%s48;2;%d;%d;%dm", CSI, c.R, c.G, c.B)
}

// Style represents text styling options
type Style struct {
	FG        Color
	BG        Color
	Bold      bool
	Dim       bool
	Italic    bool
	Underline bool
	Blink     bool
	Reverse   bool
	Strike    bool
}

// DefaultStyle returns a style with default colors and no attributes
func DefaultStyle() Style {
	return Style{
		FG: ColorDefault,
		BG: ColorDefault,
	}
}

// WithFG returns a copy of the style with the specified foreground color
func (s Style) WithFG(c Color) Style {
	s.FG = c
	return s
}

// WithBG returns a copy of the style with the specified background color
func (s Style) WithBG(c Color) Style {
	s.BG = c
	return s
}

// WithBold returns a copy of the style with bold enabled
func (s Style) WithBold() Style {
	s.Bold = true
	return s
}

// WithDim returns a copy of the style with dim enabled
func (s Style) WithDim() Style {
	s.Dim = true
	return s
}

// WithItalic returns a copy of the style with italic enabled
func (s Style) WithItalic() Style {
	s.Italic = true
	return s
}

// WithUnderline returns a copy of the style with underline enabled
func (s Style) WithUnderline() Style {
	s.Underline = true
	return s
}

// WithBlink returns a copy of the style with blink enabled
func (s Style) WithBlink() Style {
	s.Blink = true
	return s
}

// WithReverse returns a copy of the style with reverse enabled
func (s Style) WithReverse() Style {
	s.Reverse = true
	return s
}

// WithStrike returns a copy of the style with strikethrough enabled
func (s Style) WithStrike() Style {
	s.Strike = true
	return s
}

// Sequence returns the complete ANSI escape sequence for this style
func (s Style) Sequence() string {
	result := StyleReset

	if s.FG != nil {
		result += s.FG.FG()
	}
	if s.BG != nil {
		result += s.BG.BG()
	}
	if s.Bold {
		result += StyleBold
	}
	if s.Dim {
		result += StyleDim
	}
	if s.Italic {
		result += StyleItalic
	}
	if s.Underline {
		result += StyleUnderline
	}
	if s.Blink {
		result += StyleBlink
	}
	if s.Reverse {
		result += StyleReverse
	}
	if s.Strike {
		result += StyleStrike
	}

	return result
}

// Equals checks if two styles are identical
func (s Style) Equals(other Style) bool {
	return s.FG == other.FG &&
		s.BG == other.BG &&
		s.Bold == other.Bold &&
		s.Dim == other.Dim &&
		s.Italic == other.Italic &&
		s.Underline == other.Underline &&
		s.Blink == other.Blink &&
		s.Reverse == other.Reverse &&
		s.Strike == other.Strike
}

