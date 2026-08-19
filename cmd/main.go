package main

import (
	"ash/internal/shell"
	"context"
)

func main() {
	shell := shell.NewShell()
	shell.Run(context.WithCancel(context.Background()))
}
