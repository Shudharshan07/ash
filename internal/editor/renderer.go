package editor

import (
	"ash/internal/prompt"
	"unicode/utf8"
)

// remove all fmt for printing

type termWriter interface {
	Write(p []byte) (int, error)
	GetWidth() (width int)
}

type Renderer struct {
	prompt *prompt.Prompt
	term   termWriter

	state RenderState
	buf   []byte
}

func NewRenderer(term termWriter) *Renderer {
	return &Renderer{
		prompt: prompt.NewPrompt(),
		term:   term,
		state:  *NewRenderState(),
		buf:    make([]byte, 0, 256),
	}
}

func (r *Renderer) RenderPrompt() {
	b := r.buf[:0]
	b = append(b, '\n')
	b = append(b, r.prompt.Text()...)
	r.term.Write(b)
	r.buf = b
	r.state.invalidate()
}

func (r *Renderer) Draw(l *Line) {
	width := r.term.GetWidth()
	if width <= 0 {
		return
	}

	if width != r.state.width {
		r.state.invalidate()
	}

	promptWidth := r.prompt.Width()
	target := promptWidth + l.cursor
	cursorAbs := promptWidth + r.state.cursor

	prefix, changed := r.state.diff(l.text)

	if !changed {
		if target == cursorAbs {
			return
		}
		b := appendMove(r.buf[:0], cursorAbs, target, width)
		r.buf = b
		if len(b) > 0 {
			r.term.Write(b)
		}
		r.state.cursor = l.cursor
		r.state.width = width
		return
	}

	cp := promptWidth + prefix
	wroteAny := prefix < l.Len()

	b := appendMove(r.buf[:0], cursorAbs, cp, width)
	b = append(b, "\x1b[0J"...)

	for _, ch := range l.text[prefix:] {
		b = utf8.AppendRune(b, ch)
	}

	lineEnd := promptWidth + l.Len()

	if wroteAny && lineEnd > 0 && lineEnd%width == 0 {
		lineEnd--
	}

	b = appendMove(b, lineEnd, target, width)
	r.buf = b
	r.term.Write(b)

	r.state.commit(l.text, l.cursor, width)
}

func (r *Renderer) NewLine() {
	r.term.Write([]byte("\r\n"))
}
