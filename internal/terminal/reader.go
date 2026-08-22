package terminal

type Reader struct {
	term *Terminal
}

func (r *Reader) ReadKey() (Key, error) {
	b, err := r.term.ReadByte()
	if err != nil {
		return Key{}, err
	}

	switch b {
	case 13:
		return Key{Type: KeyEnter}, nil

	case 127:
		return Key{Type: KeyBackspace}, nil

	case 3:
		return Key{Type: KeyCtrlC}, nil

	case 4:
		return Key{Type: KeyCtrlD}, nil

	case 9:
		return Key{Type: KeyTab}, nil

	case 27:
		return r.readEscapeSequence()

	default:
		return r.readCharacter(b)
	}
}

func (r *Reader) readEscapeSequence() (Key, error) {
	b, err := r.term.ReadByte()
	if err != nil {
		return Key{}, err
	}

	if b != '[' {
		return Key{Type: KeyEscape}, nil
	}

	b, err = r.term.ReadByte()
	if err != nil {
		return Key{}, err
	}

	switch b {
	case 'A':
		return Key{Type: KeyUp}, nil

	case 'B':
		return Key{Type: KeyDown}, nil

	case 'C':
		return Key{Type: KeyRight}, nil

	case 'D':
		return Key{Type: KeyLeft}, nil

	case 'H':
		return Key{Type: KeyHome}, nil

	case 'F':
		return Key{Type: KeyEnd}, nil

	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return r.readTildeSequence(b)

	default:
		return Key{Type: KeyUnknown}, nil
	}
}

func (r *Reader) readCharacter(first byte) (Key, error) {
	return Key{
		Type: KeyCharacter,
		Rune: rune(first),
	}, nil
}

func (r *Reader) readTildeSequence(first byte) (Key, error) {
	digits := []byte{first}

	for {
		b, err := r.term.ReadByte()
		if err != nil {
			return Key{}, err
		}

		if b == '~' {
			break
		}
		// Some terminals send multi-digit codes (e.g. modifiers like "3;5~").
		// Keep consuming until the terminating '~' either way.
		digits = append(digits, b)
	}

	switch digits[0] {
	case '3':
		return Key{Type: KeyDelete}, nil
	case '1', '7':
		return Key{Type: KeyHome}, nil
	case '4', '8':
		return Key{Type: KeyEnd}, nil
	default:
		return Key{Type: KeyUnknown}, nil
	}
}
