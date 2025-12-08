package input

import (
	"io"
	"os"
	"time"
)

// Reader reads input from the terminal and produces events
type Reader struct {
	reader     io.Reader
	eventChan  chan Event
	stopChan   chan struct{}
	buf        []byte
	escapeTime time.Duration
}

// NewReader creates a new input reader
func NewReader() *Reader {
	return &Reader{
		reader:     os.Stdin,
		eventChan:  make(chan Event, 100),
		stopChan:   make(chan struct{}),
		buf:        make([]byte, 256),
		escapeTime: 50 * time.Millisecond,
	}
}

// Start begins reading input in a goroutine
func (r *Reader) Start() {
	go r.readLoop()
}

// Stop stops the input reader
func (r *Reader) Stop() {
	close(r.stopChan)
}

// Events returns the channel for receiving events
func (r *Reader) Events() <-chan Event {
	return r.eventChan
}

// readLoop continuously reads from stdin and parses input
func (r *Reader) readLoop() {
	for {
		select {
		case <-r.stopChan:
			return
		default:
			n, err := r.reader.Read(r.buf)
			if err != nil {
				if err != io.EOF {
					r.eventChan <- ErrorEvent{Err: err}
				}
				continue
			}

			r.parseInput(r.buf[:n])
		}
	}
}

// parseInput parses raw input bytes into events
func (r *Reader) parseInput(data []byte) {
	for len(data) > 0 {
		event, consumed := r.parseSequence(data)
		if event != nil {
			r.eventChan <- event
		}
		if consumed == 0 {
			consumed = 1
		}
		data = data[consumed:]
	}
}

// parseSequence attempts to parse an escape sequence or single key
func (r *Reader) parseSequence(data []byte) (Event, int) {
	if len(data) == 0 {
		return nil, 0
	}

	b := data[0]

	// Escape sequence
	if b == 0x1b {
		return r.parseEscape(data)
	}

	// Control characters
	if b < 32 {
		return r.parseControl(b), 1
	}

	// Delete key
	if b == 127 {
		return KeyEvent{Key: KeyBackspace}, 1
	}

	// Regular ASCII/UTF-8 character
	return r.parseRune(data)
}

// parseEscape parses escape sequences
func (r *Reader) parseEscape(data []byte) (Event, int) {
	if len(data) == 1 {
		return KeyEvent{Key: KeyEscape}, 1
	}

	// CSI sequences (ESC [)
	if data[1] == '[' {
		return r.parseCSI(data)
	}

	// SS3 sequences (ESC O) - for some function keys
	if data[1] == 'O' {
		return r.parseSS3(data)
	}

	// Alt + key
	if len(data) >= 2 {
		event, consumed := r.parseSequence(data[1:])
		if keyEvent, ok := event.(KeyEvent); ok {
			keyEvent.Modifier |= ModAlt
			return keyEvent, consumed + 1
		}
		return event, consumed + 1
	}

	return KeyEvent{Key: KeyEscape}, 1
}

// parseCSI parses CSI (Control Sequence Introducer) sequences
func (r *Reader) parseCSI(data []byte) (Event, int) {
	if len(data) < 3 {
		return KeyEvent{Key: KeyEscape}, 1
	}

	// Parse the sequence to find the final byte
	i := 2
	for i < len(data) && data[i] >= 0x30 && data[i] <= 0x3f {
		i++
	}
	for i < len(data) && data[i] >= 0x20 && data[i] <= 0x2f {
		i++
	}
	if i >= len(data) {
		return KeyEvent{Key: KeyEscape}, 1
	}

	finalByte := data[i]
	params := data[2:i]

	switch finalByte {
	case 'A':
		return KeyEvent{Key: KeyUp, Modifier: parseCSIModifier(params)}, i + 1
	case 'B':
		return KeyEvent{Key: KeyDown, Modifier: parseCSIModifier(params)}, i + 1
	case 'C':
		return KeyEvent{Key: KeyRight, Modifier: parseCSIModifier(params)}, i + 1
	case 'D':
		return KeyEvent{Key: KeyLeft, Modifier: parseCSIModifier(params)}, i + 1
	case 'H':
		return KeyEvent{Key: KeyHome, Modifier: parseCSIModifier(params)}, i + 1
	case 'F':
		return KeyEvent{Key: KeyEnd, Modifier: parseCSIModifier(params)}, i + 1
	case '~':
		return r.parseTildeSequence(params), i + 1
	case 'Z':
		return KeyEvent{Key: KeyTab, Modifier: ModShift}, i + 1
	}

	// Function keys (some terminals)
	if finalByte >= 'P' && finalByte <= 'S' {
		return KeyEvent{Key: Key(int(KeyF1) + int(finalByte-'P'))}, i + 1
	}

	return KeyEvent{Key: KeyEscape}, 1
}

// parseTildeSequence parses CSI n ~ sequences
func (r *Reader) parseTildeSequence(params []byte) Event {
	// Parse the number before the tilde
	n := 0
	mod := ModNone
	semicolonIdx := -1

	for i, b := range params {
		if b == ';' {
			semicolonIdx = i
			break
		}
		if b >= '0' && b <= '9' {
			n = n*10 + int(b-'0')
		}
	}

	// Parse modifier if present (after semicolon)
	if semicolonIdx >= 0 {
		modNum := 0
		for _, b := range params[semicolonIdx+1:] {
			if b >= '0' && b <= '9' {
				modNum = modNum*10 + int(b-'0')
			}
		}
		mod = decodeModifier(modNum)
	}

	switch n {
	case 1:
		return KeyEvent{Key: KeyHome, Modifier: mod}
	case 2:
		return KeyEvent{Key: KeyInsert, Modifier: mod}
	case 3:
		return KeyEvent{Key: KeyDelete, Modifier: mod}
	case 4:
		return KeyEvent{Key: KeyEnd, Modifier: mod}
	case 5:
		return KeyEvent{Key: KeyPageUp, Modifier: mod}
	case 6:
		return KeyEvent{Key: KeyPageDown, Modifier: mod}
	case 7:
		return KeyEvent{Key: KeyHome, Modifier: mod}
	case 8:
		return KeyEvent{Key: KeyEnd, Modifier: mod}
	case 11:
		return KeyEvent{Key: KeyF1, Modifier: mod}
	case 12:
		return KeyEvent{Key: KeyF2, Modifier: mod}
	case 13:
		return KeyEvent{Key: KeyF3, Modifier: mod}
	case 14:
		return KeyEvent{Key: KeyF4, Modifier: mod}
	case 15:
		return KeyEvent{Key: KeyF5, Modifier: mod}
	case 17:
		return KeyEvent{Key: KeyF6, Modifier: mod}
	case 18:
		return KeyEvent{Key: KeyF7, Modifier: mod}
	case 19:
		return KeyEvent{Key: KeyF8, Modifier: mod}
	case 20:
		return KeyEvent{Key: KeyF9, Modifier: mod}
	case 21:
		return KeyEvent{Key: KeyF10, Modifier: mod}
	case 23:
		return KeyEvent{Key: KeyF11, Modifier: mod}
	case 24:
		return KeyEvent{Key: KeyF12, Modifier: mod}
	}

	return KeyEvent{Key: KeyEscape}
}

// parseSS3 parses SS3 (Single Shift 3) sequences
func (r *Reader) parseSS3(data []byte) (Event, int) {
	if len(data) < 3 {
		return KeyEvent{Key: KeyEscape}, 1
	}

	switch data[2] {
	case 'A':
		return KeyEvent{Key: KeyUp}, 3
	case 'B':
		return KeyEvent{Key: KeyDown}, 3
	case 'C':
		return KeyEvent{Key: KeyRight}, 3
	case 'D':
		return KeyEvent{Key: KeyLeft}, 3
	case 'H':
		return KeyEvent{Key: KeyHome}, 3
	case 'F':
		return KeyEvent{Key: KeyEnd}, 3
	case 'P':
		return KeyEvent{Key: KeyF1}, 3
	case 'Q':
		return KeyEvent{Key: KeyF2}, 3
	case 'R':
		return KeyEvent{Key: KeyF3}, 3
	case 'S':
		return KeyEvent{Key: KeyF4}, 3
	}

	return KeyEvent{Key: KeyEscape}, 1
}

// parseControl parses control characters
func (r *Reader) parseControl(b byte) Event {
	switch b {
	case 0x00: // Ctrl+Space or Ctrl+@
		return KeyEvent{Key: KeySpace, Modifier: ModCtrl}
	case 0x09: // Tab
		return KeyEvent{Key: KeyTab}
	case 0x0a, 0x0d: // Enter (LF or CR)
		return KeyEvent{Key: KeyEnter}
	case 0x1b: // Escape
		return KeyEvent{Key: KeyEscape}
	case 0x7f: // Backspace (DEL)
		return KeyEvent{Key: KeyBackspace}
	default:
		// Ctrl+letter (Ctrl+A = 1, Ctrl+B = 2, etc.)
		if b >= 1 && b <= 26 {
			return KeyEvent{
				Key:      KeyRune,
				Rune:     rune('a' + b - 1),
				Modifier: ModCtrl,
			}
		}
	}

	return KeyEvent{Key: KeyNone}
}

// parseRune parses a UTF-8 rune
func (r *Reader) parseRune(data []byte) (Event, int) {
	// Determine UTF-8 byte length
	b := data[0]
	var runeLen int

	switch {
	case b < 0x80:
		runeLen = 1
	case b < 0xe0:
		runeLen = 2
	case b < 0xf0:
		runeLen = 3
	default:
		runeLen = 4
	}

	if len(data) < runeLen {
		return KeyEvent{Key: KeyRune, Rune: rune(b)}, 1
	}

	// Decode UTF-8
	var r2 rune
	switch runeLen {
	case 1:
		r2 = rune(b)
	case 2:
		r2 = rune(b&0x1f)<<6 | rune(data[1]&0x3f)
	case 3:
		r2 = rune(b&0x0f)<<12 | rune(data[1]&0x3f)<<6 | rune(data[2]&0x3f)
	case 4:
		r2 = rune(b&0x07)<<18 | rune(data[1]&0x3f)<<12 | rune(data[2]&0x3f)<<6 | rune(data[3]&0x3f)
	}

	return KeyEvent{Key: KeyRune, Rune: r2}, runeLen
}

// parseCSIModifier extracts modifier from CSI params
func parseCSIModifier(params []byte) Modifier {
	// Look for ;n pattern where n is modifier
	for i, b := range params {
		if b == ';' && i+1 < len(params) {
			modNum := 0
			for _, c := range params[i+1:] {
				if c >= '0' && c <= '9' {
					modNum = modNum*10 + int(c-'0')
				}
			}
			return decodeModifier(modNum)
		}
	}
	return ModNone
}

// decodeModifier decodes xterm modifier encoding
// 2=Shift, 3=Alt, 4=Shift+Alt, 5=Ctrl, 6=Shift+Ctrl, 7=Alt+Ctrl, 8=Shift+Alt+Ctrl
func decodeModifier(n int) Modifier {
	if n <= 1 {
		return ModNone
	}
	n--
	mod := ModNone
	if n&1 != 0 {
		mod |= ModShift
	}
	if n&2 != 0 {
		mod |= ModAlt
	}
	if n&4 != 0 {
		mod |= ModCtrl
	}
	return mod
}

