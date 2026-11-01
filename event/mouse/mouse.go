package mouse

import (
	"fmt"

	"github.com/stefan-muehlebach/gc9503cv/event/key"
)

//----------------------------------------------------------------------------

type Event struct {
	X, Y      float64
	Button    Button
	Modifiers key.Modifiers
	Direction Direction
}

type Button int32

const (
	ButtonNone Button = iota
	ButtonLeft
	ButtonMiddle
	ButtonRight

	ButtonWheelUp    Button = -1
	ButtonWheelDown  Button = -2
	ButtonWheelLeft  Button = -3
	ButtonWheelRight Button = -4
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

func (d Direction) String() string {
	switch d {
	case DirNone:
		return "None"
	case DirPress:
		return "Press"
	case DirRelease:
		return "Release"
	case DirStep:
		return "Step"
	default:
		return fmt.Sprintf("mouse.Direction(%d)", d)
	}
}
