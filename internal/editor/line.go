package editor

type Line struct {
	text   []rune
	cursor int
}

func (l *Line) AddCharacter(val rune) {
	l.text = append(l.text, 0)

	copy(
		l.text[l.cursor+1:],
		l.text[l.cursor:],
	)

	l.text[l.cursor] = val
	l.cursor++
}

func (l *Line) RemoveCharacter() {
	if l.cursor == 0 {
		return
	}

	copy(
		l.text[l.cursor-1:],
		l.text[l.cursor:],
	)

	l.text = l.text[:len(l.text)-1]
	l.cursor--
}

func (l *Line) MoveLeft() {

}

func (l *Line) MoveRight() {

}

func (l *Line) MoveUp() {

}

func (l *Line) MoveDown() {

}
