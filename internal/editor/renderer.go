package editor

import (
	"ash/internal/terminal"
	"fmt"
)

type Renderer struct {
	term *terminal.Terminal
}

func NewRenderer(term *terminal.Terminal) *Renderer {
	return &Renderer{
		term: term,
	}
}

func (r *Renderer) RenderPrompt(prompt string) {
	r.term.WriteString(fmt.Sprintf("\n%s", prompt))
}

func (r *Renderer) RenderInsert(l *Line) {
	if l.cursor <= 0 {
		return
	}

	tail := string(l.text[l.cursor-1:])
	r.term.WriteString(tail)

	// Now move the cursor back to origin position
	pos := len(l.text) - l.cursor
	if pos > 0 {
		r.term.WriteString(fmt.Sprintf("\x1b[%dD", pos))
	}
}

func (r *Renderer) RenderDelete(l *Line) {
	if l.cursor < 0 {
		return
	}

	tail := string(l.text[l.cursor:]) + " "
	r.term.WriteString(fmt.Sprintf("\x1b[D%s\x1b[%dD", tail, len(tail)))
}

func (r *Renderer) RenderLeft() {
	r.term.WriteString("\x1b[D")
}

func (r *Renderer) RenderRight() {
	r.term.WriteString("\x1b[C")
}
