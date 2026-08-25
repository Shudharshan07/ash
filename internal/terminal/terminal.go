package terminal

import (
	"os"

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
