package main

import (
	"image"
	//"image/draw"
	"log"
	"time"

	//"fmt"
	"github.com/holoplot/go-evdev"
	"github.com/stefan-muehlebach/gc9503cv/geom"
)

// Mit diesem Datentyp werden die unterschiedlichen Event-Arten abgebildet,
// welche durch das Druecken auf den Bildschirm an die GUI-Elemente gesendet
// werden koennen. Es ist Aufgabe der Objekte 'Screen' und 'Window', aus den
// rohen Ereginisse vom Touchscreen (Press, Drag, Release) diese Events zu
// erzeugen.
type InputEventType uint32

const (
	// Mit Einfuehrung der Maus, wird dieses Event (fahren ohne eine Taste
	// gedrueckt zu halten) relevant.
	MoveEvent InputEventType = iota
	// PressEvent wird erzeugt, wenn der Finger oder Stift auf ein Objekt
	// gedrueckt wird. Dieses Event ist vergleichbar mit einem Press-Event
	// durch die Maus.
	PressEvent
	// ReleaseEvent wird erzeugt, wenn man den Finger oder Stift wieder vom
	// Bildschirm hebt - vergleichbar mit dem Loslassen der Maustaste.
	ReleaseEvent
	// Im Unterschied zur Maus, kann ein 'Wandern' des Druckpunktes (Finger
	// oder Maus) nur erkannt werden, wenn man den Bildschirm beruehrt.
	// D.h. es gibt nur Drag-Events - ein Move-Event ist unbekannt.
	DragEvent
	// Wird laengere Zeit auf ein bestimmtes Objekt gedrueckt, erzeugt dies
	// das LongPressEvent-Event (siehe auch die Konstanten LongPressThreshold
	// und NearThreshold).
	LongPressEvent
	// Wurde das Mausrad gedreht
	WheelEvent
	// Verlaesst oder Betritt der Druckpunkt beim Wandern ein Objekt, dann
	// werden die Events TypeLeave, resp. EnterEvent gesendet.
	EnterEvent
	LeaveEvent
	// Ein Tap, bzw. DoubleTap entspricht dem Klicken resp. Doppelklicken mit
	// der Maus. Es gibt sowohl zeitliche, als auch raeumliche Grenzen, wann
	// ein Tap, resp. DoubleTap erzeugt wird (siehe Konstanten weiter unten).
	ClickEvent
	DoubleClickEvent
	// Diesen Event-Typ verwendet die Applikation, um beim Stoppen mit diesem
	// Event allfaellige Handler korrekt zu terminieren.
	QuitEvent
	numEvents
)

type MouseButtonType byte

const (
	LeftButton MouseButtonType = (1 << iota)
	MiddleButton
	RightButton
)

func (b MouseButtonType) IsSet(val MouseButtonType) bool {
	if b&val != 0x00 {
		return true
	} else {
		return false
	}
}

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

// Damit beim Debuggen klar ist, um welchen Event-Typ es sich handelt,
// implementiert Type das Stringer-Interface.
func (t InputEventType) String() string {
	switch t {
	case MoveEvent:
		return "Move"
	case PressEvent:
		return "Press"
	case ReleaseEvent:
		return "Release"
	case DragEvent:
		return "Drag"
	case LongPressEvent:
		return "LongPress"
	case WheelEvent:
		return "Wheel"
	case EnterEvent:
		return "Enter"
	case LeaveEvent:
		return "Leave"
	case ClickEvent:
		return "Click"
	case DoubleClickEvent:
		return "DoubleClick"
	default:
		return "(Unknown Type)"
	}
}

const (
	eventQueueSize = 30
)

// Dies ist der Funktionstyp für den PenEvent-Handler - also jene Funktion,
// welche beim Eintreffen eines Interrupts vom STMPE610 aufgerufen werden
// soll.
type MouseEventChannelType chan InputEvent

// Dieser Typ steht fuer das SPI Interface zum STMPE - dem Touchscreen.
type Mouse struct {
	Pos                          geom.Point[int]
	Wheel                        int
	LeftBtn, MiddleBtn, RightBtn int32
	EventQ                       MouseEventChannelType
	dev                          *evdev.InputDevice
	posRange                     geom.Rectangle[int]
	minWheel, maxWheel           int
	cursor                       *Cursor
	isOpen                       bool
}

var (
	devFile = "/dev/input/event4"
)

// Funktionen
func OpenMouse() *Mouse {
	var err error

	m := &Mouse{}
	m.dev, err = evdev.Open(devFile)
	if err != nil {
		log.Fatalf("evdev.Open failed: %v", err)
	}
	m.LeftBtn, m.MiddleBtn, m.RightBtn = 0, 0, 0
	m.EventQ = make(chan InputEvent, eventQueueSize)
	m.isOpen = true
	return m
}

func (m *Mouse) Close() {
	m.isOpen = false
	close(m.EventQ)
	m.dev.Close()
}

func (m *Mouse) SetCursor(cursor *Cursor) (old *Cursor) {
	old = m.cursor
	m.cursor = cursor
	return
}

func (m *Mouse) Bounds() image.Rectangle {
	orig := m.Pos.Sub(m.cursor.hotspot).ToInt()
	return m.cursor.img.Bounds().Add(orig)
}

func (m *Mouse) Image() image.Image {
	return m.cursor.img
}

/*
func (m *Mouse) Draw(img *image.RGBA) {
	if m.cursor == nil {
		return
	}
	orig := m.Pos.Sub(m.cursor.hotspot).ToInt()
	dstRect := m.cursor.img.Bounds().Add(orig).Add(img.Rect.Min)
	draw.Draw(img, dstRect, m.cursor.img, image.Point{}, draw.Over)
}
*/

func (m *Mouse) SetPosRange(rect geom.Rectangle[int]) {
	m.posRange = rect
	m.Pos = rect.Center()
}

func (m *Mouse) SetWheelRange(ini, min, max int) {
	m.Wheel = ini
	m.minWheel = min
	m.maxWheel = max
}

func (m *Mouse) Start() {
	go m.processEvents()
}

func (m *Mouse) Stop() {
	m.enqueueEvent(InputEvent{Type: QuitEvent})
}

// Diese Funktion wird von 'aussen' aufgerufen und gibt das nächste Pen-Event
// zurück. Es ist eine Alternative zum Lesen aus der öffentlichen Event-Queue.
func (m *Mouse) WaitForEvent() InputEvent {
	return <-m.EventQ
}

// Mit dieser Funktion wird ein neues Pen-Event in die zentrale Event-Queue
// gestellt (welche dann von der Applikation ausgelesen werden muss).
// Diese Operation darf nicht blockierend ausgeführt werden, andernfalls
// würde der Event-Handler blockiert - was in meinen Augen gravierender ist.
//
// Mit dem auskommentierten Code kann für Testzwecke dafür gesorgt werden,
// dass bei einem Fehler ein Runtime-Panic ausgelöst wird.
func (m *Mouse) enqueueEvent(ev InputEvent) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("Runtime panic: %v\n", x)
		}
	}()
	select {
	case m.EventQ <- ev:
	default:
		log.Printf("Sending not possible: event queue full!\n")
	}
}

//-----------------------------------------------------------------------------

func (m *Mouse) processEvents() {
	var mev InputEvent

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
				if mev.Button > 0 {
					mev.Type = DragEvent
				} else {
					mev.Type = MoveEvent
				}
				mev.Pos = m.Pos.ToFloat()
				mev.Time = time.Now()
				m.enqueueEvent(mev)
			case evdev.REL_WHEEL:
				m.Wheel += int(e.Value)
				if m.Wheel < m.minWheel {
					m.Wheel = m.minWheel
				} else if m.Wheel >= m.maxWheel {
					m.Wheel = m.maxWheel - 1
				}
				mev.Type = WheelEvent
				mev.Wheel = m.Wheel
				mev.WheelRel = int(e.Value)
				mev.Time = time.Now()
				m.enqueueEvent(mev)
			default:
				continue
			}

		case evdev.EV_KEY:
			switch e.Value {
			case 1:
				mev.Type = PressEvent
				mev.InitTime = time.Now()
				mev.Time = mev.InitTime
				mev.InitPos = m.Pos.ToFloat()
				mev.Pos = mev.InitPos
				mev.LongPressed = false
				switch e.Code {
				case evdev.BTN_LEFT:
					mev.Button = LeftButton
				case evdev.BTN_RIGHT:
					mev.Button = RightButton
				case evdev.BTN_MIDDLE:
					mev.Button = MiddleButton
				}

				go func() {
					mevCopy := mev
					time.Sleep(LongPressThreshold)
					if (mev.Type == PressEvent || mev.Type == DragEvent) &&
						mev.Button == mevCopy.Button &&
						mev.InitPos.Distance(mev.Pos) < NearThreshold {
						mev.LongPressed = true
						mevCopy = mev
						mevCopy.Type = LongPressEvent
						mevCopy.Time = time.Now()
						m.enqueueEvent(mevCopy)
					}
				}()
				m.enqueueEvent(mev)

			case 0:
				mev.Type = ReleaseEvent
				mev.Time = time.Now()
				mev.Pos = m.Pos.ToFloat()
				m.enqueueEvent(mev)

				if mev.InitPos.Distance(mev.Pos) < NearThreshold &&
					mev.Time.Sub(mev.InitTime) < ClickDuration {
					mev.Type = ClickEvent
					m.enqueueEvent(mev)
				}
				mev.InitPos = Point{}
				mev.InitTime = time.Time{}
				mev.LongPressed = false
				mev.Button = 0x00

			case 2:
				// Fuer repetition eines buttons.
			}
		default:
			continue
		}
	}
}

//----------------------------------------------------------------------------

// Jedes Ereignis des Touchscreens wird durch eine Variable des Typs
// 'Event' repraesentiert.
type InputEvent struct {
	Type           InputEventType
	Time, InitTime time.Time
	InitPos, Pos   Point // geom.Point[int]
	Button         MouseButtonType
	LongPressed    bool
	Wheel          int
	WheelRel       int
}
