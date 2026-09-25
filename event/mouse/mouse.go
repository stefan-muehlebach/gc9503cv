package mouse

import (
	"github.com/stefan-muehlebach/gc9503cv/event/key"
)

//----------------------------------------------------------------------------

type Event struct {
	X, Y float64
	Button Button
	Modifiers key.Modifiers
	Direction Direction
}

type Button int32

const (
	ButtonNone Button = +0
	ButtonLeft Button = +1
	ButtonMiddle Button = +2
	ButtonRight Button = +3

	ButtonWheelUp Button = -1
	ButtonWheelDown Button = -2
)

func (b Button) IsWheel() bool {
	return b < 0
}

type Direction uint8

const (
	DirNone Direction = iota
	DirPress
	DirRelease
	DirStep
)

