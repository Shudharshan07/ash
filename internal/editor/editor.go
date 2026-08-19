package editor

import (
	"ash/internal/parser"
	"ash/internal/terminal"
	"context"
)

type Editor struct {
	ctx    context.Context
	cancel context.CancelFunc

	prompt string // Will be replased
	line   Line
	reader terminal.Reader
	parser parser.Parser

	renderer *Renderer
}

func NewEditor(term *terminal.Terminal, ctx context.Context, cancel context.CancelFunc) *Editor {
	return &Editor{
		ctx:      ctx,
		cancel:   cancel,
		prompt:   "$",
		reader:   term.Reader,
		renderer: NewRenderer(term),
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
		e.Insert(key.Rune)
	case terminal.KeyBackspace:
		e.Delete()

	case terminal.KeyLeft:
		e.Left()

	case terminal.KeyRight:
		e.Right()

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

	if res == "exit" {
		e.cancel()
	}
	e.line.CleanLine()
	e.renderer.RenderPrompt(e.prompt)
}

func (e *Editor) Insert(r rune) {
	e.line.AddCharacter(r)
	e.renderer.RenderInsert(&e.line)
}

func (e *Editor) Delete() {
	if e.line.cursor > 0 {
		e.line.RemoveCharacter()
		e.renderer.RenderDelete(&e.line)
	}
}

func (e *Editor) Left() {
	if e.line.cursor > 0 {
		e.line.MoveLeft()
		e.renderer.RenderLeft()
	}
}

func (e *Editor) Right() {
	if e.line.cursor < len(e.line.text) {
		e.line.MoveRight()
		e.renderer.RenderRight()
	}
}
