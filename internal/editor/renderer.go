package editor

import (
	"ash/internal/prompt"
	"strconv"
	"unicode/utf8"
)

// maybe i should move the Render Cursor to the editor, so i get a clean unduplicated codde
// remove all fmt for printing

type termWriter interface {
	Write(p []byte) (int, error)
	GetSize() (width, height int, err error)
	GetWidth() (width int)
}

type Renderer struct {
	prompt *prompt.Prompt
	term   termWriter

	prev      []rune
	cursorAbs int
	buf       []byte
}

func NewRenderer(term termWriter) *Renderer {
	return &Renderer{
		prompt: prompt.NewPrompt(),
		term:   term,
		prev:   make([]rune, 0, 128),
		buf:    make([]byte, 0, 256),
	}
}

func (r *Renderer) RenderPrompt() {
	b := r.buf[:0]
	b = append(b, '\n')
	b = append(b, r.prompt.Text...)
	r.term.Write(b)
	r.buf = b
	r.prev = r.prev[:0]
	r.cursorAbs = r.prompt.Width()
}

func (r *Renderer) Draw(l *Line) {
	width := r.term.GetWidth()
	if width <= 0 {
		return
	}

	promptWidth := r.prompt.Width()
	lineEnd := promptWidth + len(l.text)

	isTextChanged := len(l.text) != len(r.prev)
	prefix := -1
	if !isTextChanged {
		prefix = prefixLen(l.text, r.prev)
		isTextChanged = prefix != len(l.text)
	}

	target := promptWidth + l.cursor

	if !isTextChanged && target == r.cursorAbs {
		return
	}

	b := r.buf[:0]

	if !isTextChanged {
		b = appendMove(b, r.cursorAbs, target, width)
		r.cursorAbs = target
		r.buf = b
		if len(b) > 0 {
			r.term.Write(b)
		}
		return
	}

	if prefix == -1 {
		prefix = prefixLen(l.text, r.prev)
	}

	cp := promptWidth + prefix
	wroteAny := prefix < len(l.text)

	b = appendMove(b, r.cursorAbs, cp, width)
	b = append(b, "\x1b[0J"...)

	for _, ch := range l.text[prefix:] {
		b = utf8.AppendRune(b, ch)
	}

	final := lineEnd
	if wroteAny && lineEnd > 0 && lineEnd%width == 0 {
		final = lineEnd - 1
	}

	b = appendMove(b, final, target, width)

	r.cursorAbs = target

	if cap(r.prev) < len(l.text) {
		r.prev = make([]rune, len(l.text))
	} else {
		r.prev = r.prev[:len(l.text)]
	}
	copy(r.prev, l.text)

	r.buf = b
	r.term.Write(b)
}

func prefixLen(a, b []rune) int {
	n := min(len(a), len(b))

	for i := range n {
		if a[i] != b[i] {
			return i
		}
	}

	return n
}

func appendCSI(b []byte, n int, final byte) []byte {
	if n <= 0 {
		return b
	}
	b = append(b, '\x1b', '[')
	b = strconv.AppendInt(b, int64(n), 10)
	return append(b, final)
}

func appendMove(b []byte, from, to, width int) []byte {
	if from == to {
		return b
	}

	fromRow := from / width
	toRow, toCol := to/width, to%width

	b = appendCSI(b, fromRow-toRow, 'A')
	b = appendCSI(b, toRow-fromRow, 'B')
	b = append(b, '\r')
	b = appendCSI(b, toCol, 'C')

	return b
}
