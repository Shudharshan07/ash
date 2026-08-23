package editor

type RenderState struct {
	text   []rune
	cursor int
	width  int
}

func NewRenderState() *RenderState {
	return &RenderState{
		text: make([]rune, 0, 128),
	}
}

func (s *RenderState) diff(text []rune) (prefix int, changed bool) {
	prefix = prefixLen(text, s.text)
	changed = prefix != len(text) || len(text) != len(s.text)
	return prefix, changed
}

func (s *RenderState) invalidate() {
	s.text = s.text[:0]
	s.cursor = 0
	s.width = 0
}

func (s *RenderState) commit(text []rune, cursor, width int) {
	if cap(s.text) < len(text) {
		s.text = make([]rune, len(text))
	} else {
		s.text = s.text[:len(text)]
	}
	copy(s.text, text)
	s.cursor = cursor
	s.width = width
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
