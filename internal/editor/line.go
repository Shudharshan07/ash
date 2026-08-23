package editor

import "slices"

type Line struct {
	text   []rune
	cursor int
}

func NewLine() *Line {
	return &Line{
		text:   make([]rune, 0),
		cursor: 0,
	}
}

func (l *Line) CleanLine() {
	l.text = l.text[:0]
	l.cursor = 0
}

func (l *Line) Insert(r rune) {
	l.text = slices.Insert(l.text, l.cursor, r)
	l.cursor++
}

func (l *Line) Backspace() bool {
	if l.cursor == 0 {
		return false
	}

	l.text = slices.Delete(l.text, l.cursor-1, l.cursor)
	l.cursor--
	return true
}

func (l *Line) Delete() bool {
	if l.cursor >= len(l.text) {
		return false
	}

	l.text = slices.Delete(l.text, l.cursor, l.cursor+1)
	return true
}

func (l *Line) MoveLeft() bool {
	if l.cursor <= 0 {
		return false
	}
	l.cursor--
	return true
}

func (l *Line) MoveRight() bool {
	if l.cursor >= len(l.text) {
		return false
	}
	l.cursor++
	return true
}

func (l *Line) Start() {
	l.cursor = 0
}

func (l *Line) End() {
	l.cursor = len(l.text)
}
