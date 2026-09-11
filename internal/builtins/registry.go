package builtins

import (
	"ash/internal/terminal"
	"context"
	"fmt"
)

type Builtin func([]string) error

type Registry struct {
	ctx    context.Context
	cancel context.CancelFunc

	terminal *terminal.Terminal

	registry map[string]Builtin
}

func NewRegistry(term *terminal.Terminal, ctx context.Context, cancel context.CancelFunc) *Registry {
	registry := &Registry{
		ctx:      ctx,
		cancel:   cancel,
		terminal: term,
	}

	var reg = map[string]Builtin{
		"cd":   registry.builtinCd,
		"pwd":  registry.builtinPwd,
		"exit": registry.builtinExit,
		"cls":  registry.builtinClear,
	}

	registry.registry = reg

	return registry
}

func (r *Registry) Exists(name string) bool {
	_, ok := r.registry[name]

	return ok
}

func (r *Registry) Get(name string) (Builtin, bool) {
	fn, ok := r.registry[name]

	return fn, ok
}

func (r *Registry) Run(name string, args []string) error {
	fn, ok := r.registry[name]
	if !ok {
		return fmt.Errorf("command not found: %s", name) // we need to change this too
	}

	return fn(args)
}
