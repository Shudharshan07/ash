package editor

import "slices"

type Command []rune

type Line struct {
	text   Command
	cursor int
}

func NewLine() *Line {
	return &Line{
		text:   make(Command, 0),
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
	if l.cursor >= l.Len() {
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
	if l.cursor >= l.Len() {
		return false
	}
	l.cursor++
	return true
}

func (l *Line) Start() bool {
	if l.cursor == 0 {
		return false
	}

	l.cursor = 0
	return true
}

func (l *Line) End() bool {
	if l.cursor == l.Len() {
		return false
	}

	l.cursor = l.Len()
	return true
}

func (l *Line) Len() int {
	return len(l.text)
}

func (l *Line) IsEmpty() bool {
	return len(l.text) == 0
}
