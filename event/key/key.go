package key

//----------------------------------------------------------------------------

type Event struct {
	Rune      rune
	Code      Code
	Modifiers Modifiers
	Direction Direction
}

/*
func (e Event) String() string {
	if e.Rune >= 0 {
		return fmt.Sprintf("key.Event{%q (%v), %v, %v}", e.Rune, e.Code,
			e.Modifiers, e.Direction)
	}
	return fmt.Sprintf("key.Event{(%v), %v, %v}", e.Code, e.Modifiers,
		e.Direction)
}
*/

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
