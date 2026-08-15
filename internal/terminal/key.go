package terminal

type KeyType uint8

const (
	KeyUnknown KeyType = iota

	KeyCharacter
	KeyEnter
	KeyBackspace
	KeyDelete

	KeyLeft
	KeyRight
	KeyUp
	KeyDown

	KeyHome
	KeyEnd

	KeyCtrlC
	KeyCtrlD
	KeyTab
	KeyEscape
)

type Key struct {
	Type KeyType
	Rune rune
}
