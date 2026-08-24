package editor

import (
	"slices"
)

type HistoryEntry []rune

func NewEntry(l *Line) HistoryEntry {
	return slices.Clone(l.text)
}

func (h HistoryEntry) equal(l *Line) bool {
	return slices.Equal(l.text, h)
}

type History struct {
	commands []HistoryEntry
	index    int
}

func NewHistory() *History {
	return &History{
		commands: make([]HistoryEntry, 0, 10),
		index:    0,
	}
}

func (h *History) SaveHistory(l *Line) {
	if l.IsEmpty() {
		return
	}

	if !h.IsEmpty() && h.commands[h.Len()-1].equal(l) {
		return
	}

	h.commands = append(h.commands, NewEntry(l))
	h.index = h.Len()
}

func (h *History) MoveUp() []rune {
	if h.IsEmpty() {
		return nil
	}
	if h.index > 0 {
		h.index--
	}
	return slices.Clone(h.commands[h.index])
}

func (h *History) MoveDown() []rune {
	if h.IsEmpty() {
		return nil
	}
	h.index++
	if h.index >= h.Len() {
		h.index = h.Len()
		return nil
	}
	return slices.Clone(h.commands[h.index])
}

func (h *History) GetHistory() []rune {
	if h.IsEmpty() || h.index < 0 || h.index >= h.Len() {
		return nil
	}
	return slices.Clone(h.commands[h.index])
}

func (h *History) IsEmpty() bool {
	return len(h.commands) == 0
}

func (h *History) Len() int {
	return len(h.commands)
}
