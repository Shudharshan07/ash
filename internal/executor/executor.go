package executor

import (
	"ash/internal/builtins"
	"ash/internal/parser"
	"ash/internal/terminal"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
)

type Executor struct {
	ctx    context.Context
	cancel context.CancelFunc

	mu sync.Mutex

	parser *parser.Parser
	term   *terminal.Terminal

	registry *builtins.Registry
}

func NewExecutor(term *terminal.Terminal, ctx context.Context, cancel context.CancelFunc) *Executor {
	e := &Executor{
		ctx:    ctx,
		cancel: cancel,

		parser:   parser.NewParser(),
		term:     term,
		registry: builtins.NewRegistry(term, ctx, cancel),
	}

	go e.HandleSignals()

	return e
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

	if _, err := exec.LookPath(cmd[0]); err != nil {
		fmt.Printf("%s: command not found\n", cmd[0])
		return
	}

	ctx, cancel := context.WithCancel(e.ctx)
	e.setCancel(cancel)

	exe := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	exe.Stdin = e.term.Stdin()
	exe.Stdout = e.term.Stdout()
	exe.Stderr = e.term.Stderr()

	e.term.DisableRawMode()
	err = exe.Run()
	if rerr := e.term.EnableRawMode(); rerr != nil {
		fmt.Fprintf(os.Stderr, "raw mode restore failed: %v\n", rerr)
	}

	cancel()
	e.setCancel(nil)
}

func (e *Executor) HandleError(err error) {
	fmt.Println(err.Error())
}

func (e *Executor) HandleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for range sigChan {
		if c := e.getCancel(); c != nil {
			c()
		}
	}
}

func (e *Executor) setCancel(c context.CancelFunc) {
	e.mu.Lock()
	e.cancel = c
	e.mu.Unlock()
}

func (e *Executor) getCancel() context.CancelFunc {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.cancel
}
