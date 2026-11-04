//go:build !v2

package main

import (
	"image"
	//"image/draw"
	"log"
	"time"

	//"fmt"
	evdev "github.com/mikegio27/go-evdev"
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

type MouseButton byte

const (
	LeftButton MouseButton = (1 << iota)
	MiddleButton
	RightButton
)

func (b MouseButton) IsSet(val MouseButton) bool {
	return b&val != 0x00
}

type Modifiers uint32

const (
	ModShift Modifiers = (1 << iota)
	ModCtrl
	ModAlt
)

func (m Modifiers) IsSet(val Modifiers) bool {
	return m&val != 0x00
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
type InputEventChan chan InputEvent

// Dieser Typ steht fuer das SPI Interface zum STMPE - dem Touchscreen.
type Mouse struct {
	devs                         []*evdev.Device
	EventQ                       InputEventChan
	LeftBtn, MiddleBtn, RightBtn int32
	Pos                          geom.Point[int]
	posRange                     geom.Rectangle[int]
	Wheel                        int
	//minWheel, maxWheel           int
	cursor *Cursor
	isOpen bool
}

var (
	devFiles = []string{
		"/dev/input/myMouse",
		"/dev/input/myKeyboard",
	}
)

// Funktionen
func OpenMouse() *Mouse {
	m := &Mouse{}
	for _, devFile := range devFiles {
		dev, err := evdev.Open(devFile)
		if err != nil {
			log.Fatalf("evdev.Open failed: %v", err)
		}
		m.devs = append(m.devs, dev)
	}
	m.LeftBtn, m.MiddleBtn, m.RightBtn = 0, 0, 0
	m.EventQ = make(chan InputEvent, eventQueueSize)
	m.isOpen = true
	return m
}

func (m *Mouse) Close() {
	m.isOpen = false
	close(m.EventQ)
	for _, dev := range m.devs {
		dev.Close()
	}
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

//func (m *Mouse) SetWheelRange(ini, min, max int) {
//	m.Wheel = ini
//	m.minWheel = min
//	m.maxWheel = max
//}

func (m *Mouse) Start() {
	for _, dev := range m.devs {
		go m.processEvents(dev)
	}
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

var (
	ev InputEvent
)

func (m *Mouse) processEvents(dev *evdev.Device) {
	for {
		e, err := dev.ReadOne()
		if err != nil {
			log.Fatalf("ReadOne() failed: %v", err)
		}
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
				if ev.Button > 0 {
					ev.Type = DragEvent
				} else {
					ev.Type = MoveEvent
				}
				ev.Pos = m.Pos.ToFloat()
				ev.Time = time.Now()
				m.enqueueEvent(ev)
			case evdev.REL_WHEEL:
				m.Wheel += int(e.Value)
				ev.Type = WheelEvent
				ev.Wheel = int(e.Value)
				ev.Time = time.Now()
				m.enqueueEvent(ev)
				ev.Wheel = 0
			default:
				continue
			}

		case evdev.EV_KEY:
			switch e.Code {
			case evdev.BTN_LEFT, evdev.BTN_RIGHT, evdev.BTN_MIDDLE:
				switch e.Value {
				case 1:
					ev.Type = PressEvent
					ev.InitTime = time.Now()
					ev.Time = ev.InitTime
					ev.InitPos = m.Pos.ToFloat()
					ev.Pos = ev.InitPos
					ev.LongPressed = false
					switch e.Code {
					case evdev.BTN_LEFT:
						ev.Button = LeftButton
					case evdev.BTN_RIGHT:
						ev.Button = RightButton
					case evdev.BTN_MIDDLE:
						ev.Button = MiddleButton
					}

					go func() {
						evCopy := ev
						time.Sleep(LongPressThreshold)
						if (ev.Type == PressEvent || ev.Type == DragEvent) &&
							ev.Button == evCopy.Button &&
							ev.InitPos.Distance(ev.Pos) < NearThreshold {
							ev.LongPressed = true
							evCopy = ev
							evCopy.Type = LongPressEvent
							evCopy.Time = time.Now()
							m.enqueueEvent(evCopy)
						}
					}()
					m.enqueueEvent(ev)

				case 0:
					ev.Type = ReleaseEvent
					ev.Time = time.Now()
					ev.Pos = m.Pos.ToFloat()
					m.enqueueEvent(ev)

					if ev.InitPos.Distance(ev.Pos) < NearThreshold &&
						ev.Time.Sub(ev.InitTime) < ClickDuration {
						ev.Type = ClickEvent
						m.enqueueEvent(ev)
					}
					ev.InitPos = Point{}
					ev.InitTime = time.Time{}
					ev.LongPressed = false
					ev.Button = 0x00

				case 2:
					// Fuer repetition eines buttons.
				}
			case evdev.KEY_LEFTSHIFT, evdev.KEY_RIGHTSHIFT:
				switch e.Value {
				case 0:
					ev.Modifiers = ev.Modifiers &^ ModShift
				case 1:
					ev.Modifiers = ev.Modifiers | ModShift
				}
				//log.Printf("%+v", ev)
			case evdev.KEY_LEFTCTRL, evdev.KEY_RIGHTCTRL:
				switch e.Value {
				case 0:
					ev.Modifiers = ev.Modifiers &^ ModCtrl
				case 1:
					ev.Modifiers = ev.Modifiers | ModCtrl
				}
			case evdev.KEY_LEFTALT, evdev.KEY_RIGHTALT:
				switch e.Value {
				case 0:
					ev.Modifiers = ev.Modifiers &^ ModAlt
				case 1:
					ev.Modifiers = ev.Modifiers | ModAlt
				}
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
	InitTime, Time time.Time
	InitPos, Pos   Point // geom.Point[int]
	Button         MouseButton
	Wheel          int
	LongPressed    bool
	Modifiers      Modifiers
}
