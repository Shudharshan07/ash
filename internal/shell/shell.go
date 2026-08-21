package shell

import (
	editor "ash/internal/editor"
	terminal "ash/internal/terminal"
	"context"
	"os"
)

type Shell struct {
	ctx    context.Context
	cancel context.CancelFunc

	pwd    string
	editor *editor.Editor
	term   *terminal.Terminal
}

func NewShell(ctx context.Context, cancel context.CancelFunc) *Shell {
	pwd, err := os.Getwd()
	if err != nil {
		panic("Error getting the working dir")
	}
	term := terminal.NewTerminal()
	shell := &Shell{
		ctx:    ctx,
		cancel: cancel,
		pwd:    pwd,
		term:   term,
	}

	shell.editor = editor.NewEditor(term, ctx, cancel)

	return shell
}

func (s *Shell) Run() {
	s.term.EnableRawMode()

	defer s.term.DisableRawMode()

	s.editor.Init()
	go s.editor.Listen()
	<-s.ctx.Done()
}
