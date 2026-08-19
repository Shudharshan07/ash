package editor

import (
	"ash/internal/parser"
	"ash/internal/terminal"
	"context"
	"fmt"
)

type Editor struct {
	ctx    context.Context
	cancel context.CancelFunc

	line   Line
	reader terminal.Reader
	parser parser.Parser
}

func NewEditor(reader terminal.Reader, ctx context.Context, cancel context.CancelFunc) *Editor {
	return &Editor{
		ctx:    ctx,
		cancel: cancel,
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
	case terminal.KeyEnter:
		e.ExecuteCommand(e.line.text)
	}
}

func (e *Editor) ExecuteCommand(cmd []rune) { // temp input
	res := e.parser.Parse(cmd)

	fmt.Print(res)

	if res == "exit" {
		e.cancel()
	}
}
