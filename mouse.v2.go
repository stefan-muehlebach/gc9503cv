//go:build v2

package main

import (
	"image"
	//"image/draw"
	"log"
	"time"

	//"fmt"
	"github.com/holoplot/go-evdev"
	"github.com/stefan-muehlebach/gc9503cv/event"
	"github.com/stefan-muehlebach/gc9503cv/event/mouse"
	"github.com/stefan-muehlebach/gc9503cv/geom"
)

//----------------------------------------------------------------------------

const (
	// TapDuration ist die Zeit, welche max. zwischen Press und Release
	// vergehen darf, damit dieses Ereignis als Tap interpretiert wird.
	ClickDuration = 200 * time.Millisecond
	// Analog dazu ist DoubleTapDuration die max. Dauer welche zwischen zwei
	// Tap-Ereignissen vergehen darf, damit sie zusammen als DoubleTap inter-
	// pretiert werden.
	DoubleClickDuration = 200 * time.Millisecond
	// Drueckt der Benutzer fuer mehr als LongPressThreshold auf ein Objekt,
	// wird ein LongPress-Event erzeugt.
	LongPressThreshold = 400 * time.Millisecond
	// Fuer die Ereignisse Tap, DoubleTap und auch LongPress darf sich der
	// Finger auf dem TouchScreen nicht zu stark bewegen. Der maximale Abstand
	// zwischen dem Press-Event und der aktuellen Position darf nicht mehr
	// NearThreshold betragen.
	NearThreshold = 8.0
)

type Mouse struct {
	dev                *evdev.InputDevice
	Pos                geom.Point[int]
	posRange           geom.Rectangle[int]
	Wheel              int
	minWheel, maxWheel int
	cursor             *Cursor

	LeftBtn, MiddleBtn, RightBtn int32
	eventQ                       event.EventQueueType
	isOpen                       bool
}

var (
	devFile = "/dev/input/myMouse"
)

// Funktionen
func NewMouse() *Mouse {
	var err error

	m := &Mouse{}
	m.Open()
	return m
}

func (m *Mouse) Open() {
	var err error

	m.dev, err = evdev.Open(devFile)
	if err != nil {
		log.Fatalf("evdev.Open failed: %v", err)
	}
	m.isOpen = true
}

func (m *Mouse) Close() {
	m.isOpen = false
	m.dev.Close()
}

func (m *Mouse) SetCursor(cursor *Cursor) *Cursor {
	old := m.cursor
	m.cursor = cursor
	return old
}

func (m *Mouse) Bounds() image.Rectangle {
	orig := m.Pos.Sub(m.cursor.hotspot).ToInt()
	return m.cursor.img.Bounds().Add(orig)
}

func (m *Mouse) Image() image.Image {
	return m.cursor.img
}

func (m *Mouse) SetPosRange(rect geom.Rectangle[int]) {
	m.posRange = rect
	m.Pos = rect.Center()
}

func (m *Mouse) SetWheelRange(ini, min, max int) {
	m.Wheel = ini
	m.minWheel = min
	m.maxWheel = max
}

func (m *Mouse) Start(eventQ event.EventQueueType) {
	m.eventQ = eventQ
	go m.processEvents()
}

func (m *Mouse) Stop() {

}

//-----------------------------------------------------------------------------

func (m *Mouse) enqueueEvent(ev event.InputEvent) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("Runtime panic: %v\n", x)
		}
	}()
	select {
	case m.eventQ <- ev:
	default:
		log.Printf("Sending not possible: event queue full!\n")
	}
}

func (m *Mouse) processEvents() {
	var mev mouse.Event
	var pressPos geom.Point[int]
	var pressTime time.Time

	for {
		e, err := m.dev.ReadOne()
		if err != nil {
			log.Fatalf("ReadOne() failed: %v", err)
		}
		//log.Printf("system type: %d", e.Type)
		switch e.Type {
		case evdev.EV_REL:
			switch e.Code {
			case evdev.REL_X, evdev.REL_Y:
				if e.Code == evdev.REL_X {
					m.Pos.X += int(e.Value)
					if m.Pos.X < m.posRange.Min.X {
						m.Pos.X = m.posRange.Min.X
					} else if m.Pos.X >= m.posRange.Max.X {
						m.Pos.X = m.posRange.Max.X - 1
					}
				} else {
					m.Pos.Y += int(e.Value)
					if m.Pos.Y < m.posRange.Min.Y {
						m.Pos.Y = m.posRange.Min.Y
					} else if m.Pos.Y >= m.posRange.Max.Y {
						m.Pos.Y = m.posRange.Max.Y - 1
					}
				}
				mev.X, mev.Y = m.Pos.ToFloat().AsCoord()
				mev.Direction = mouse.DirNone
				m.enqueueEvent(mev)

			case evdev.REL_WHEEL:
				m.Wheel += int(e.Value)
				if m.Wheel < m.minWheel {
					m.Wheel = m.minWheel
				} else if m.Wheel >= m.maxWheel {
					m.Wheel = m.maxWheel - 1
				}
				if e.Value > 0 {
					mev.Button = mouse.ButtonWheelUp
				} else {
					mev.Button = mouse.ButtonWheelDown
				}
				mev.Direction = mouse.DirNone
				m.enqueueEvent(mev)

			default:
				continue
			}

		case evdev.EV_KEY:
			switch e.Value {
			case 1:
				mev.Direction = mouse.DirPress
				switch e.Code {
				case evdev.BTN_LEFT:
					mev.Button = mouse.ButtonLeft
				case evdev.BTN_RIGHT:
					mev.Button = mouse.ButtonRight
				case evdev.BTN_MIDDLE:
					mev.Button = mouse.ButtonMiddle
				}
				m.enqueueEvent(mev)

				pressPos = m.Pos
				pressTime = time.Now()

			case 0:
				mev.Direction = mouse.DirRelease
				switch e.Code {
				case evdev.BTN_LEFT:
					mev.Button = mouse.ButtonLeft
				case evdev.BTN_RIGHT:
					mev.Button = mouse.ButtonRight
				case evdev.BTN_MIDDLE:
					mev.Button = mouse.ButtonMiddle
				}
				m.enqueueEvent(mev)

				if pressPos.Distance(m.Pos) < NearThreshold &&
					time.Since(pressTime) < ClickDuration {
					mev.Dir = mouse.DirStep
					m.enqueueEvent(mev)
				}

			case 2:
				// Fuer repetition eines buttons.
			}
		default:
			continue
		}
	}
}
