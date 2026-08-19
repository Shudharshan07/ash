package main

import (
	"ash/internal/shell"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	shell := shell.NewShell(ctx, cancel)

	shell.Run()
}
