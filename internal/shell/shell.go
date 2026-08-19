package shell

import (
	editor "ash/internal/editor"
	terminal "ash/internal/terminal"
	"context"
	"fmt"
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

	shell.editor = editor.NewEditor(term.Reader, ctx, cancel)

	return shell
}

func (s Shell) init() {
	// The init stuff that has to be run before the shell
	fmt.Printf("%s >", s.pwd)
}

func (s *Shell) Run() {
	s.term.EnableRawMode()

	defer s.term.DisableRawMode()

	s.init()
	go s.editor.Listen()
	<-s.ctx.Done()
}
