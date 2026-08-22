package editor

import (
	"ash/internal/parser"
	"ash/internal/terminal"
	"context"
)

type Editor struct {
	ctx    context.Context
	cancel context.CancelFunc

	line   *Line
	reader terminal.Reader
	parser parser.Parser

	renderer *Renderer
}

func NewEditor(term *terminal.Terminal, ctx context.Context, cancel context.CancelFunc) *Editor {
	return &Editor{
		ctx:      ctx,
		cancel:   cancel,
		reader:   term.Reader,
		line:     NewLine(),
		renderer: NewRenderer(term),
	}
}

func (e *Editor) Init() {
	// Add some cool stuff
	e.renderer.RenderPrompt()
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
		e.Backspace()

	case terminal.KeyDelete:
		e.Delete()

	case terminal.KeyLeft:
		e.Left()

	case terminal.KeyRight:
		e.Right()

	case terminal.KeyHome:
		e.Start()

	case terminal.KeyEnd:
		e.End()

	case terminal.KeyUp:
		e.line.MoveUp()

	case terminal.KeyDown:
		e.line.MoveDown()
	case terminal.KeyEnter:
		e.ExecuteCommand(e.line.text)
	}
}

func (e *Editor) ExecuteCommand(cmd []rune) { // temp input
	// e.End()
	res := e.parser.Parse(cmd)

	if res == "exit" {
		e.cancel()
	}
	e.line.CleanLine()
	e.renderer.RenderPrompt()
}

// editor.go (relevant methods)
func (e *Editor) Insert(r rune) {
	e.line.Insert(r)
	e.renderer.RenderInsert(e.line)
}

func (e *Editor) Backspace() {
	if e.line.Backspace() {
		e.renderer.RenderBackspace(e.line)
	}
}

func (e *Editor) Delete() {
	if e.line.Delete() {
		e.renderer.RenderDelete(e.line)
	}
}

func (e *Editor) Left() {
	if e.line.cursor == 0 {
		return
	}
	prev := *e.line
	e.line.MoveLeft()
	e.renderer.RenderCursor(e.line, &prev)
}

func (e *Editor) Right() {
	if e.line.cursor == len(e.line.text) {
		return
	}
	prev := *e.line
	e.line.MoveRight()
	e.renderer.RenderCursor(e.line, &prev)
}

func (e *Editor) Start() {
	prev := *e.line
	e.line.Start()
	e.renderer.RenderCursor(e.line, &prev)
}

func (e *Editor) End() {
	prev := *e.line
	e.line.End()
	e.renderer.RenderCursor(e.line, &prev)
}
