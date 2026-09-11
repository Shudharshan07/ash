package executor

import (
	"ash/internal/parser"
	"ash/internal/terminal"
	"fmt"
	"os/exec"
)

type Executor struct {
	parser *parser.Parser
	term   *terminal.Terminal
}

func NewExecutor(term *terminal.Terminal) *Executor {
	return &Executor{
		parser: parser.NewParser(),
		term:   term,
	}
}

func (e *Executor) Run(command []rune) {
	cmd, err := e.parser.Parse(command)

	fmt.Println(cmd)
	if err != nil {
		e.HandleError(err)
		return
	}

	if len(cmd) == 0 {
		// send some error
		return
	}

	exe := exec.Command(cmd[0], cmd[1:]...)

	exe.Stdin = e.term.Stdin()
	exe.Stdout = e.term.Stdout()
	exe.Stderr = e.term.Stderr()

	// the error and the output should be handled properly
	e.term.DisableRawMode()
	err = exe.Run()
	e.term.EnableRawMode()
	if err != nil {
		e.HandleError(err)
	}
}

func (e *Executor) HandleError(err error) {
	fmt.Print("\n" + err.Error())
}
