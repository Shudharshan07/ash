package builtins

import (
	"os"
)

func (r *Registry) builtinCd(args []string) error {
	if len(args) == 0 {
		return nil
	}
	target := args[0]

	return os.Chdir(target)
}

func (r *Registry) builtinPwd(args []string) error {
	val, err := os.Getwd()

	if err != nil {
		return err
	}

	r.terminal.WriteString(val)
	return nil
}

func (r *Registry) builtinClear(args []string) error {
	r.terminal.WriteString("\x1b[H\x1b[2J\x1b[3J")
	return nil
}

func (r *Registry) builtinExit(args []string) error {
	r.cancel()
	return nil
}

// ll, ls, dir, cat | type, where, touch, export, unset , echo
// pause
