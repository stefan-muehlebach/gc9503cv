package main

import (
	"image"
	"image/draw"
	"log"
	"time"
	//"fmt"
	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/holoplot/go-evdev"
)

// Mit diesem Datentyp werden die unterschiedlichen Event-Arten abgebildet,
// welche durch das Druecken auf den Bildschirm an die GUI-Elemente gesendet
// werden koennen. Es ist Aufgabe der Objekte 'Screen' und 'Window', aus den
// rohen Ereginisse vom Touchscreen (Press, Drag, Release) diese Events zu
// erzeugen.
type MouseEventType uint32

const (
	// Mit Einfuehrung der Maus, wird dieses Event (fahren ohne eine Taste
	// gedrueckt zu halten) relevant.
	TypeMove MouseEventType = iota
	// TypePress wird erzeugt, wenn der Finger oder Stift auf ein Objekt
	// gedrueckt wird. Dieses Event ist vergleichbar mit einem Press-Event
	// durch die Maus.
	TypePress
	// TypeRelease wird erzeugt, wenn man den Finger oder Stift wieder vom
	// Bildschirm hebt - vergleichbar mit dem Loslassen der Maustaste.
	TypeRelease
	// Im Unterschied zur Maus, kann ein 'Wandern' des Druckpunktes (Finger
	// oder Maus) nur erkannt werden, wenn man den Bildschirm beruehrt.
	// D.h. es gibt nur Drag-Events - ein Move-Event ist unbekannt.
	TypeDrag
	// Wird laengere Zeit auf ein bestimmtes Objekt gedrueckt, erzeugt dies
	// das TypeLongPress-Event (siehe auch die Konstanten LongPressThreshold
	// und NearThreshold).
	TypeLongPress
	// Wurde das Mausrad gedreht
	TypeWheel
	// Verlaesst oder Betritt der Druckpunkt beim Wandern ein Objekt, dann
	// werden die Events TypeLeave, resp. TypeEnter gesendet.
	TypeEnter
	TypeLeave
	// Ein Tap, bzw. DoubleTap entspricht dem Klicken resp. Doppelklicken mit
	// der Maus. Es gibt sowohl zeitliche, als auch raeumliche Grenzen, wann
	// ein Tap, resp. DoubleTap erzeugt wird (siehe Konstanten weiter unten).
	TypeClick
	TypeDoubleClick
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
func (t MouseEventType) String() string {
	switch t {
	case TypeMove:
		return "Move"
	case TypePress:
		return "Press"
	case TypeRelease:
		return "Release"
	case TypeDrag:
		return "Drag"
	case TypeLongPress:
		return "LongPress"
	case TypeWheel:
		return "Wheel"
	case TypeEnter:
		return "Enter"
	case TypeLeave:
		return "Leave"
	case TypeClick:
		return "Click"
	case TypeDoubleClick:
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
type MouseEventChannelType chan MouseEvent

// Dieser Typ steht fuer das SPI Interface zum STMPE - dem Touchscreen.
type Mouse struct {
	Pos                          geom.Point[int]
	Wheel                        int
	LeftBtn, MiddleBtn, RightBtn int32
	EventQ                       MouseEventChannelType
	dev                          *evdev.InputDevice
	posRange             geom.Rectangle[int]
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
	m.EventQ = make(chan MouseEvent, eventQueueSize)
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

func (m *Mouse) Draw(img *image.RGBA) {
	if m.cursor == nil {
		return
	}
	orig := m.Pos.Sub(m.cursor.hotspot).ToInt()
	dstRect := m.cursor.img.Bounds().Add(orig).Add(img.Rect.Min)
	draw.Draw(img, dstRect, m.cursor.img, image.Point{}, draw.Over)
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

func (m *Mouse) StartEvents() {
	go m.processEvents()
}

// Diese Funktion wird von 'aussen' aufgerufen und gibt das nächste Pen-Event
// zurück. Es ist eine Alternative zum Lesen aus der öffentlichen Event-Queue.
func (m *Mouse) WaitForEvent() MouseEvent {
	return <-m.EventQ
}

// Mit dieser Funktion wird ein neues Pen-Event in die zentrale Event-Queue
// gestellt (welche dann von der Applikation ausgelesen werden muss).
// Diese Operation darf nicht blockierend ausgeführt werden, andernfalls
// würde der Event-Handler blockiert - was in meinen Augen gravierender ist.
//
// Mit dem auskommentierten Code kann für Testzwecke dafür gesorgt werden,
// dass bei einem Fehler ein Runtime-Panic ausgelöst wird.
func (m *Mouse) enqueueEvent(ev MouseEvent) {
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
	var mev MouseEvent

	for {
		e, err := m.dev.ReadOne()
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
				if mev.Button > 0 {
					mev.Type = TypeDrag
				} else {
					mev.Type = TypeMove
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
				mev.Type = TypeWheel
				mev.Wheel = m.Wheel
				mev.Time = time.Now()
				m.enqueueEvent(mev)
			default:
				continue
			}

		case evdev.EV_KEY:
			if e.Value == 1 {
				mev.Type = TypePress
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
					if (mev.Type == TypePress || mev.Type == TypeDrag) &&
						mev.Button == mevCopy.Button &&
						mev.InitPos.Distance(mev.Pos) < NearThreshold {
						mev.LongPressed = true
						mevCopy = mev
						//mevCopy.Type = TypeLongPress
						mevCopy.Time = time.Now()
						m.enqueueEvent(mevCopy)
					}
				}()
				m.enqueueEvent(mev)

			} else {
				mev.Type = TypeRelease
				mev.Time = time.Now()
				mev.Pos = m.Pos.ToFloat()
				m.enqueueEvent(mev)

				if mev.InitPos.Distance(mev.Pos) < NearThreshold &&
					mev.Time.Sub(mev.InitTime) < ClickDuration {
					//log.Printf("Click will be generated")
					mev.Type = TypeClick
					m.enqueueEvent(mev)
				}
				mev.InitPos = Point{}
				mev.InitTime = time.Time{}
				mev.LongPressed = false
				mev.Button = 0x00
			}
		default:
			continue
		}
		//m.enqueueEvent(mev)
	}
}

//----------------------------------------------------------------------------

// In diesem Datentyp ist alles zusammengefasst, was an Touch-Events an die
// applikatorischen Elemente gesendet werden kann.
/*
type Event struct {
	// Der Typ des Events (moegliche Typen: die die Konstanten TypeXXX).
	Type MouseEventType
	// Alle Events werden sequenziell durchnumeriert, damit man bspw. das
	// Release-Event mit einem vorgaengigen Press-Event in Verbindung bringen
	// kann.
	//SeqNumber int
	// Dieses Feld wird auf true gesetzt, sobald ein LongPressed-Ereignis
	// erkannt wird.
	LongPressed bool
	// In InitTime und InitPos werden Zeitpunkt und Position des Press-Events
	// (des initialen Events) festgehalten.
	InitTime time.Time
	InitPos  geom.Point[int]
	// Wohingegen Time und Pos die Zeit und die Position des aktuellen Events
	// enthalten.
	Time time.Time
	Pos  geom.Point[int]
}
*/

// Jedes Ereignis des Touchscreens wird durch eine Variable des Typs
// 'Event' repraesentiert.
type MouseEvent struct {
	Type           MouseEventType
	Time, InitTime time.Time
	Pos, InitPos   Point // geom.Point[int]
	Button         MouseButtonType
	LongPressed    bool
	Wheel          int
}

