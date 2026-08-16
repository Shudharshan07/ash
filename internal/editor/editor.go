package editor

import (
	"ash/internal/terminal"
)

type Editor struct {
	line   Line
	reader terminal.Reader
}

func NewEditor(reader terminal.Reader) *Editor {
	return &Editor{
		reader: reader,
	}
}

func (e *Editor) Listen() error {
	for {
		key, err := e.reader.ReadKey()
		if err != nil {
			return err
		}

		e.handleKey(key)
	}
}

func (e *Editor) handleKey(key terminal.Key) {
	switch key.Type {
	case terminal.KeyCharacter:
		e.line.AddCharacter(key.Rune)
	case terminal.KeyBackspace:
		e.line.RemoveCharacter()

	case terminal.KeyLeft:
		e.line.MoveLeft()

	case terminal.KeyRight:
		e.line.MoveRight()

	case terminal.KeyUp:
		e.line.MoveUp()

	case terminal.KeyDown:
		e.line.MoveDown()
	}
}
