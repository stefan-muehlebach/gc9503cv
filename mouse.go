package main

import (
	//"image"
	"log"
	"time"
	"fmt"
	"github.com/holoplot/go-evdev"
	"github.com/stefan-muehlebach/gg"
	// "github.com/stefan-muehlebach/gg/geom"
	"gc9503cv/gc9503cv/geom"
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
	posRange                     geom.Rectangle[int]
	Wheel                        int
	LeftBtn, MiddleBtn, RightBtn int32
	EventQ                       MouseEventChannelType
	dev                          *evdev.InputDevice
	isOpen                       bool
	minWheel, maxWheel           int
	cursor *Cursor
}

var (
	devFile = "/dev/input/event0"
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

func (m *Mouse) Draw(gc *gg.Context) {
	if m.cursor == nil {
		return
	}
    orig := m.Pos.Sub(m.cursor.hotspot)
    gc.DrawImage(m.cursor.img, float64(orig.X), float64(orig.Y))
}

func (m *Mouse) SetPosRange(rect geom.Rectangle[int]) {
	m.posRange = rect
	m.Pos = rect.Size().Div(2)
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
	ev.Time = time.Now()
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

/*
type MouseEventIface interface {
	Type() MouseEventType
	Time() time.Time
}

type mouseEventEmbed struct {
	typ MouseEventType
	tim time.Time
}

func (ev *mouseEventEmbed) Type() MouseEventType {
	return ev.typ
}

func (ev *mouseEventEmbed) Time() time.Time {
	return ev.tim
}

type MouseMoveEvent struct {
	mouseEventEmbed
	Step geom.Point[int]
}

func NewMouseMoveEvent(step geom.Point[int]) MouseEventIface {
	ev := &MouseMoveEvent{}
	ev.typ = TypeMove
	ev.tim = time.Now()
	ev.Step = step
	return ev
}
*/

func (m *Mouse) processEvents() {
	var mev MouseEvent
	var pressPos geom.Point[int]
	var pressTime time.Time

	for {
		e, err := m.dev.ReadOne()
		if err != nil {
			log.Fatalf("ReadOne() failed: %v", err)
		}
		switch e.Type {
		case evdev.EV_REL:
			switch e.Code {
			case evdev.REL_X:
				m.Pos.X += int(e.Value)
				if m.Pos.X < 0 {
					m.Pos.X = 0
				} else if m.Pos.X >= m.posRange.Max.X {
					m.Pos.X = m.posRange.Max.X - 1
				}
			case evdev.REL_Y:
				m.Pos.Y += int(e.Value)
				if m.Pos.Y < 0 {
					m.Pos.Y = 0
				} else if m.Pos.Y >= m.posRange.Max.Y {
					m.Pos.Y = m.posRange.Max.Y - 1
				}
			case evdev.REL_WHEEL:
				m.Wheel += int(e.Value)
				if m.Wheel < m.minWheel {
					m.Wheel = m.minWheel
				} else if m.Wheel >= m.maxWheel {
					m.Wheel = m.maxWheel - 1
				}
			default:
				continue
			}
			if m.LeftBtn != 0 || m.MiddleBtn != 0 || m.RightBtn != 0 {
				mev.Type = TypeDrag
			} else {
				mev.Type = TypeMove
			}
			mev.Pos = m.Pos
			mev.Wheel = m.Wheel

		case evdev.EV_KEY:
			switch e.Code {
			case evdev.BTN_LEFT, evdev.BTN_MIDDLE, evdev.BTN_RIGHT:
				if e.Value == 1 {
					mev.Type = TypePress
					pressPos = mev.Pos
					pressTime = time.Now()
				} else {
					if pressPos.Distance(mev.Pos) < NearThreshold &&
							time.Since(pressTime) < ClickDuration {
						mev.Type = TypeClick
					} else {
						mev.Type = TypeRelease
					}
				}
				switch e.Code {
				case evdev.BTN_LEFT:
					m.LeftBtn = e.Value
				case evdev.BTN_MIDDLE:
					m.MiddleBtn = e.Value
				case evdev.BTN_RIGHT:
					m.RightBtn = e.Value
				}
				mev.LeftBtn = m.LeftBtn
				mev.MiddleBtn = m.MiddleBtn
				mev.RightBtn = m.RightBtn
			default:
				continue
			}

		default:
			continue
		}
		m.enqueueEvent(mev)
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
	Type                         MouseEventType
	Time                         time.Time
	Pos                          geom.Point[int]
	Wheel                        int
	LeftBtn, MiddleBtn, RightBtn int32
}

// Fuer das Debugging implementiert Event das Stringer-Interface.
func (evt MouseEvent) String() string {
	return fmt.Sprintf("%v %s %v", evt.Type,
		evt.Time.Format("15:04:05.000000"), evt.Pos)
}

// Alle Callback-Handler fuer die Ereignisse vom Touchscreen, muessen folgendes
// Profil aufweisen.
type CallbackType func(evt MouseEvent)

// Alle GUI-Elemente, welche ueber den Touchscreen gesteuert werden sollen,
// muesssen diesen Datentyp einbetten. Damit werden auch alle unten
// aufgefuehrten Methoden geerbt und es koennen Handler fuer die diversen
// Touchscreen-Ereignisse hinterlegt werden. Im Array touchFuncList kann
// fuer jedes Ereginis max. eine Funktion hinterlegt werden.
type callbackEmbed struct {
	callbackList [numEvents]CallbackType
}

// Diese Methode wird durch AdaGui aufgerufen, um ein Touch-Ereignis an
// ein GUI-Element zu senden. Es gibt eine Default-Implementation, welche das
// Event via CallTouchFunc an registrierte Event-Handler sendet.
// Es ist jedoch ueblich, dass ein GUI-Element diese Methode ueberschreibt
// um bspw. visuelle Anpassungen zu machen und dann selber CallTouchFunc
// aufruft.
func (m *callbackEmbed) OnInputEvent(evt MouseEvent) {
	if fnc := m.callbackList[evt.Type]; fnc != nil {
		fnc(evt)
	}
}

// Mit SetCallback wird die Funktion fnc als Handler fuer den Event typ
// registriert. Eine bereits registrierte Funktion wird damit ueberschrieben.
func (m *callbackEmbed) SetCallback(fnc CallbackType, types ...MouseEventType) {
	for _, typ := range types {
		m.callbackList[typ] = fnc
	}
}

func (m *callbackEmbed) SetOnMove(fnc CallbackType) {
	m.SetCallback(fnc, TypeMove)
}

// Registriert fnc als Handler fuer den Press-Event.
func (m *callbackEmbed) SetOnPress(fnc CallbackType) {
	m.SetCallback(fnc, TypePress)
}

// Registriert fnc als Handler fuer den Release-Event.
func (m *callbackEmbed) SetOnRelease(fnc CallbackType) {
	m.SetCallback(fnc, TypeRelease)
}

// Registriert fnc als Handler fuer den Drag-Event.
func (m *callbackEmbed) SetOnDrag(fnc CallbackType) {
	m.SetCallback(fnc, TypeDrag)
}

// Registriert fnc als Handler fuer den LongPress-Event.
func (m *callbackEmbed) SetOnLongPress(fnc CallbackType) {
	m.SetCallback(fnc, TypeLongPress)
}

// Registriert fnc als Handler fuer den Enter-Event.
func (m *callbackEmbed) SetOnEnter(fnc CallbackType) {
	m.SetCallback(fnc, TypeEnter)
}

// Registriert fnc als Handler fuer den Leave-Event.
func (m *callbackEmbed) SetOnLeave(fnc CallbackType) {
	m.SetCallback(fnc, TypeLeave)
}

// Registriert fnc als Handler fuer den Click-Event.
func (m *callbackEmbed) SetOnClick(fnc CallbackType) {
	m.SetCallback(fnc, TypeClick)
}

// Registriert fnc als Handler fuer den DoubleClick-Event.
func (m *callbackEmbed) SetOnDoubleClick(fnc CallbackType) {
	m.SetCallback(fnc, TypeDoubleClick)
}
