//go:build v2

package main

import (
	"log"

	. "github.com/stefan-muehlebach/gc9503cv/event"
	"github.com/stefan-muehlebach/gc9503cv/event/key"
	"github.com/stefan-muehlebach/gc9503cv/event/mouse"
	"github.com/stefan-muehlebach/gc9503cv/event/touch"
)

// Konsequenterweise wurde auch fuer diese Funktionalitaet ein Interface
// gebaut.
type EventHandler interface {
	OnInputEvent(ev InputEvent)
	SetOnKeyEvent(fnc KeyEventHandlerFunc)
	SetOnMouseEvent(fnc MouseEventHandlerFunc)
	SetOnTouchEvent(fnc TouchEventHandlerFunc)
}

// Alle Callback-Handler fuer die Ereignisse vom Touchscreen, muessen folgendes
// Profil aufweisen.
type KeyEventHandlerFunc func(ev key.Event)
type MouseEventHandlerFunc func(ev mouse.Event)
type TouchEventHandlerFunc func(ev touch.Event)

// Alle GUI-Elemente, welche ueber den Touchscreen gesteuert werden sollen,
// muesssen diesen Datentyp einbetten. Damit werden auch alle unten
// aufgefuehrten Methoden geerbt und es koennen Handler fuer die diversen
// Touchscreen-Ereignisse hinterlegt werden. Im Array touchFuncList kann
// fuer jedes Ereginis max. eine Funktion hinterlegt werden.
type eventHandlerEmbed struct {
	keyEventHandler KeyEventHandlerFunc
	mouseEventHandler MouseEventHandlerFunc
	touchEventHandler TouchEventHandlerFunc
}

// Diese Methode wird durch AdaGui aufgerufen, um ein Touch-Ereignis an
// ein GUI-Element zu senden. Es gibt eine Default-Implementation, welche das
// Event via CallTouchFunc an registrierte Event-Handler sendet.
// Es ist jedoch ueblich, dass ein GUI-Element diese Methode ueberschreibt
// um bspw. visuelle Anpassungen zu machen und dann selber CallTouchFunc
// aufruft.
func (m *eventHandlerEmbed) OnInputEvent(ev InputEvent) {
	m.CallEventHandler(ev)
}

func (m *eventHandlerEmbed) CallEventHandler(ev InputEvent) {
	switch e := ev.(type) {
	case key.Event:
		if m.keyEventHandler != nil {
			m.keyEventHandler(e)
		}
	case mouse.Event:
		if m.mouseEventHandler != nil {
		m.mouseEventHandler(e)
		}
	case touch.Event:
		if m.touchEventHandler != nil {
		m.touchEventHandler(e)
		}
	default:
		log.Printf("No handler for event type: %T", ev)
	}
}

func (m *eventHandlerEmbed) SetOnKeyEvent(fnc KeyEventHandlerFunc) {
	m.keyEventHandler = fnc
}

func (m *eventHandlerEmbed) SetOnMouseEvent(fnc MouseEventHandlerFunc) {
	m.mouseEventHandler = fnc
}

func (m *eventHandlerEmbed) SetOnTouchEvent(fnc TouchEventHandlerFunc) {
	m.touchEventHandler = fnc
}

