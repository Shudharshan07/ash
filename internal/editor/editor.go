package editor

import (
	"ash/internal/executor"
	"ash/internal/terminal"
	"context"
)

type Editor struct {
	ctx    context.Context
	cancel context.CancelFunc

	line     *Line
	reader   terminal.Reader
	executor *executor.Executor

	renderer *Renderer

	history *History
}

func NewEditor(term *terminal.Terminal, ctx context.Context, cancel context.CancelFunc) *Editor {
	return &Editor{
		ctx:      ctx,
		cancel:   cancel,
		reader:   term.Reader,
		line:     NewLine(),
		executor: executor.NewExecutor(term),
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
	isRun := true

	switch key.Type {
	case terminal.KeyCharacter:
		e.line.Insert(key.Rune)
	case terminal.KeyBackspace:
		isRun = e.line.Backspace()

	case terminal.KeyDelete:
		isRun = e.line.Delete()

	case terminal.KeyLeft:
		isRun = e.line.MoveLeft()

	case terminal.KeyRight:
		isRun = e.line.MoveRight()

	case terminal.KeyHome:
		isRun = e.line.Start()

	case terminal.KeyEnd:
		isRun = e.line.End()

	case terminal.KeyUp:
		e.Up()

	case terminal.KeyDown:
		e.Down()

	case terminal.KeyEnter:
		e.ExecuteCommand()
	}

	if isRun {
		e.renderer.Draw(e.line)
	}
}

func (e *Editor) ExecuteCommand() error {
	e.End() // will render the cursor to the end of line

	e.history.SaveHistory(e.line)

	cmd := e.line.text
	if string(cmd) == "exit" {
		e.cancel()
	}

	e.executor.Run(cmd)
	// show the output using some write, we need to manage the context so we can kiil this sheel (we need to implement the inbuilts)

	e.line.CleanLine()
	e.renderer.RenderPrompt()

	return nil
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

func (e *Editor) End() {
	e.line.End()
	e.renderer.Draw(e.line)
}
