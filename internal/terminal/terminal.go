package terminal

import (
	"os"
	"unicode/utf8"

	"golang.org/x/term"
)

type Terminal struct {
	file   *os.File
	out    *os.File
	state  *term.State
	Reader Reader
}

func NewTerminal() *Terminal {
	term := &Terminal{
		file: os.Stdin,
		out:  os.Stdout,
	}

	term.Reader.term = term

	return term
}

func (t *Terminal) EnableRawMode() error {
	state, err := term.MakeRaw(int(t.file.Fd()))
	t.state = state

	return err
}

func (t *Terminal) DisableRawMode() error {
	return term.Restore(int(t.file.Fd()), t.state)
}

func (t *Terminal) ReadByte() (byte, error) {
	var buf [1]byte

	_, err := t.file.Read(buf[:])
	if err != nil {
		return 0, err
	}

	return buf[0], nil
}

func (t *Terminal) ReadRune() (rune, int, error) {
	var buf [4]byte

	// Read the first byte
	_, err := t.file.Read(buf[:1])
	if err != nil {
		return 0, 0, err
	}

	first := buf[0]

	// If it's standard ASCII (0-127), it's only 1 byte long.
	if first < utf8.RuneSelf {
		return rune(first), 1, nil
	}

	// Otherwise, determine how many bytes this character uses
	var size int
	if first >= 0xC0 && first <= 0xDF {
		size = 2
	} else if first >= 0xE0 && first <= 0xEF {
		size = 3
	} else if first >= 0xF0 && first <= 0xF7 {
		size = 4
	} else {
		// Invalid UTF-8 start byte
		return utf8.RuneError, 1, nil
	}

	// Read the remaining bytes for this specific character
	for i := 1; i < size; i++ {
		_, err := t.file.Read(buf[i : i+1])
		if err != nil {
			return 0, i, err
		}
	}

	// Decode the full byte slice into a single rune
	r, _ := utf8.DecodeRune(buf[:size])
	return r, size, nil
}

func (t *Terminal) WriteString(s string) (int, error) {
	return t.out.WriteString(s)
}

func (t *Terminal) GetSize() (width, height int, err error) {
	return term.GetSize(int(t.out.Fd()))
}

func (t *Terminal) GetWidth() (width int) {
	w, _, _ := term.GetSize(int(t.out.Fd()))
	return w
}

func (t *Terminal) Write(p []byte) (int, error) {
	return t.out.Write(p)
}

func (t *Terminal) Stdin() *os.File {
	return t.file
}

func (t *Terminal) Stdout() *os.File {
	return t.out
}

func (t *Terminal) Stderr() *os.File {
	return t.out
}
