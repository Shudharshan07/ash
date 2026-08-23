package editor

import "strconv"

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
