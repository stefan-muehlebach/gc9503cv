package event

//----------------------------------------------------------------------------

type Event struct {
	X, Y float64
	Sequence Sequence
	Type Type
}

type Sequence int64

type Type byte

const (
	TypeBegin Type = iota
	TypeMove
	TypeEnd
)


