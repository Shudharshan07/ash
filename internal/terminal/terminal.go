package terminal

import (
	"os"

	"golang.org/x/term"
)

type Terminal struct {
	file   *os.File
	state  *term.State
	Reader Reader
}

func NewTerminal() *Terminal {
	term := &Terminal{
		file: os.Stdin,
	}

	term.Reader.term = term

	return term
}

func (t *Terminal) EnableRawMode() error {
	state, err := term.MakeRaw(int(t.file.Fd()))
	if err != nil {
		return err
	}

	t.state = state
	return nil
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
