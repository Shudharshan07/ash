package shell

import (
	editor "ash/internal/editor"
	terminal "ash/internal/terminal"
	"context"
	"fmt"
	"os"
)

type Shell struct {
	pwd    string
	editor *editor.Editor
	term   *terminal.Terminal
}

func NewShell() *Shell {
	pwd, err := os.Getwd()
	if err != nil {
		panic("Error getting the working dir")
	}
	term := terminal.NewTerminal()
	shell := &Shell{
		pwd:  pwd,
		term: term,
	}

	shell.editor = editor.NewEditor(term.Reader)

	return shell
}

func (s Shell) init() {
	// The init stuff that has to be run before the shell
	s.term.EnableRawMode()
	fmt.Printf("%s >", s.pwd)
}

func (s *Shell) Run(ctx context.Context, cancel context.CancelFunc) {
	s.init()
	go s.editor.Listen()
	<-ctx.Done()
}
