package editor

import "slices"

type Line struct {
	text   []rune
	cursor int
}

func (l *Line) CleanLine() {
	l.text = l.text[:0]
	l.cursor = 0
}

func (l *Line) AddCharacter(val rune) {
	l.text = slices.Insert(l.text, l.cursor, val)
	l.cursor++
}

func (l *Line) RemoveCharacter() {
	if l.cursor == 0 {
		return
	}
	l.text = slices.Delete(l.text, l.cursor-1, l.cursor)
	l.cursor--
}

func (l *Line) MoveLeft() {
	if l.cursor > 0 {
		l.cursor--
	}
}

func (l *Line) MoveRight() {
	if l.cursor < len(l.text) {
		l.cursor++
	}
}

func (l *Line) MoveUp() {

}

func (l *Line) MoveDown() {

}
