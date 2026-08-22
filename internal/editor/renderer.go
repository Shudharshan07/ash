package editor

import (
	"ash/internal/prompt"
	"ash/internal/terminal"
	"fmt"
)

// maybe i should move the Render Cursor to the editor, so i get a clean unduplicated codde
// remove all fmt for printing

type Renderer struct {
	prompt *prompt.Prompt
	term   *terminal.Terminal

	startRow int
	startCol int
}

func NewRenderer(term *terminal.Terminal) *Renderer {
	return &Renderer{
		prompt: prompt.NewPrompt(),
		term:   term,
	}
}

func (r *Renderer) RenderPrompt() {
	r.term.WriteString(fmt.Sprintf("\n%s", r.prompt.Text))
}

func (r *Renderer) moveCursor(dx, dy int) {
	if dx != 0 {
		if dx < 0 {
			r.term.WriteString(fmt.Sprintf("\x1b[%dD", -dx))
		} else {
			r.term.WriteString(fmt.Sprintf("\x1b[%dC", dx))
		}
	}
	if dy != 0 {
		if dy < 0 {
			r.term.WriteString(fmt.Sprintf("\x1b[%dA", -dy))
		} else {
			r.term.WriteString(fmt.Sprintf("\x1b[%dB", dy))
		}
	}
}

func (r *Renderer) moveCursorFromTo(from, to int) {
	width, _, _ := r.term.GetSize()
	if width <= 0 {
		return
	}

	fromAbs := r.prompt.Width() + from
	toAbs := r.prompt.Width() + to

	fromRow, fromCol := fromAbs/width, fromAbs%width
	toRow, toCol := toAbs/width, toAbs%width

	r.moveCursor(toCol-fromCol, toRow-fromRow)
}

func (r *Renderer) RenderInsert(l *Line) {
	if l.cursor <= 0 {
		return
	}

	inserted := l.text[l.cursor-1]
	r.term.WriteString(string(inserted))

	tail := string(l.text[l.cursor:])
	if tail == "" {
		return
	}

	r.term.WriteString("\x1b7")
	r.term.WriteString(tail)
	r.term.WriteString("\x1b8")
}

func (r *Renderer) RenderBackspace(l *Line) {
	r.term.WriteString("\x1b[D")
	r.term.WriteString("\x1b7")

	tail := string(l.text[l.cursor:])

	r.term.WriteString(tail)
	r.term.WriteString(" ")

	r.term.WriteString("\x1b8")
}

func (r *Renderer) RenderDelete(l *Line) {
	r.term.WriteString("\x1b7")
	tail := string(l.text[l.cursor:])

	r.term.WriteString(tail)
	r.term.WriteString(" ")
	r.term.WriteString("\x1b8")
}

func (r *Renderer) RenderCursor(l *Line, prev *Line) {
	r.moveCursorFromTo(prev.cursor, l.cursor)
}

func (r *Renderer) toCoords(l *Line) (row, col int) {
	width, _, _ := r.term.GetSize()

	if width <= 0 {
		return 0, 0
	}

	absolute := r.prompt.Width() + l.cursor

	return absolute / width, absolute % width
}
