package key

//----------------------------------------------------------------------------

type Event struct {
	Rune rune
	Code Code
	Modifiers Modifiers
	Direction Direction
}

type Code uint32

type Modifiers uint32

const (
	ModShift Modifiers = 1 << iota
	ModControl
	ModAlt
	ModMeta
)

type Direction uint8

const (
	DirNone Direction = iota
	DirPress
	DirRelease
)

