package editor

import "slices"

type History struct {
	commands [][]rune
	index    int
}

func NewHistory() *History {
	return &History{
		commands: make([][]rune, 0, 10),
		index:    0,
	}
}

func (h *History) SaveHistory(l []rune) {
	if len(l) == 0 {
		return
	}
	if !h.isEmpty() && slices.Equal(h.commands[len(h.commands)-1], l) {
		return
	}

	h.commands = append(h.commands, slices.Clone(l))
	h.index = len(h.commands)
}

func (h *History) MoveUp() []rune {
	if h.isEmpty() {
		return nil
	}
	if h.index > 0 {
		h.index--
	}
	return slices.Clone(h.commands[h.index])
}

func (h *History) MoveDown() []rune {
	if h.isEmpty() {
		return nil
	}
	h.index++
	if h.index >= len(h.commands) {
		h.index = len(h.commands)
		return nil
	}
	return slices.Clone(h.commands[h.index])
}

func (h *History) GetHistory() []rune {
	if h.isEmpty() || h.index < 0 || h.index >= len(h.commands) {
		return nil
	}
	return slices.Clone(h.commands[h.index])
}

func (h *History) isEmpty() bool {
	return len(h.commands) == 0
}

func (h *History) Len() int {
	return len(h.commands)
}
