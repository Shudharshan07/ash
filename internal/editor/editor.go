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

	history *History
}

func NewEditor(term *terminal.Terminal, ctx context.Context, cancel context.CancelFunc) *Editor {
	return &Editor{
		ctx:      ctx,
		cancel:   cancel,
		reader:   term.Reader,
		line:     NewLine(),
		renderer: NewRenderer(term),
		history:  NewHistory(),
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
		e.Up()

	case terminal.KeyDown:
		e.Down()

	case terminal.KeyEnter:
		e.ExecuteCommand()
	}
}

func (e *Editor) ExecuteCommand() { // temp input
	// if no text no need to save to history
	cmd := e.line.text
	e.End()
	e.history.SaveHistory(cmd)

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
	e.renderer.Draw(e.line)
}

func (e *Editor) Backspace() {
	if e.line.Backspace() {
		e.renderer.Draw(e.line)
	}
}

func (e *Editor) Delete() {
	if e.line.Delete() {
		e.renderer.Draw(e.line)
	}
}

func (e *Editor) Left() {
	if e.line.MoveLeft() {
		e.renderer.Draw(e.line)
	}
}

func (e *Editor) Right() {
	if e.line.MoveRight() {
		e.renderer.Draw(e.line)
	}
}

func (e *Editor) Start() {
	e.line.Start()
	e.renderer.Draw(e.line)
}

func (e *Editor) End() {
	e.line.End()
	e.renderer.Draw(e.line)
}

func (e *Editor) Up() {
	line := e.history.MoveUp()
	if line == nil {
		return
	}
	e.line.text = line
	e.End()
	e.renderer.Draw(e.line)
}

func (e *Editor) Down() {
	line := e.history.MoveDown()
	if line == nil {
		return
	}
	e.line.text = line
	e.End()
	e.renderer.Draw(e.line)
}
