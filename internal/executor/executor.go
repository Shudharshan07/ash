package executor

import (
	"ash/internal/builtins"
	"ash/internal/parser"
	"ash/internal/terminal"
	"context"
	"fmt"
	"os/exec"
)

type Executor struct {
	// ctx    context.Context
	// cancel context.CancelFunc

	parser *parser.Parser
	term   *terminal.Terminal

	registry *builtins.Registry
}

func NewExecutor(term *terminal.Terminal, ctx context.Context, cancel context.CancelFunc) *Executor {
	return &Executor{
		parser:   parser.NewParser(),
		term:     term,
		registry: builtins.NewRegistry(term, ctx, cancel),
	}
}

func (e *Executor) Run(command []rune) {
	cmd, err := e.parser.Parse(command)

	if err != nil {
		e.HandleError(err)
		return
	}

	if len(cmd) == 0 {
		return
	}

	if e.registry.Exists(cmd[0]) {
		e.registry.Run(cmd[0], cmd[1:])
		return
	}

	exe := exec.Command(cmd[0], cmd[1:]...)

	exe.Stdin = e.term.Stdin()
	exe.Stdout = e.term.Stdout()
	exe.Stderr = e.term.Stderr()

	e.term.DisableRawMode()
	err = exe.Run()
	e.term.EnableRawMode()

	// the error and the output should be handled properly
}

func (e *Executor) HandleError(err error) {
	fmt.Println(err.Error())
}
