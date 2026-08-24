package main

// Alle Callback-Handler fuer die Ereignisse vom Touchscreen, muessen folgendes
// Profil aufweisen.
type EventHandler func(ev MouseEvent)

// Alle GUI-Elemente, welche ueber den Touchscreen gesteuert werden sollen,
// muesssen diesen Datentyp einbetten. Damit werden auch alle unten
// aufgefuehrten Methoden geerbt und es koennen Handler fuer die diversen
// Touchscreen-Ereignisse hinterlegt werden. Im Array touchFuncList kann
// fuer jedes Ereginis max. eine Funktion hinterlegt werden.
type eventHandlerEmbed struct {
	eventHandlerList [numEvents]EventHandler
}

// Diese Methode wird durch AdaGui aufgerufen, um ein Touch-Ereignis an
// ein GUI-Element zu senden. Es gibt eine Default-Implementation, welche das
// Event via CallTouchFunc an registrierte Event-Handler sendet.
// Es ist jedoch ueblich, dass ein GUI-Element diese Methode ueberschreibt
// um bspw. visuelle Anpassungen zu machen und dann selber CallTouchFunc
// aufruft.
func (m *eventHandlerEmbed) OnInputEvent(ev MouseEvent) {
	m.CallEventHandler(ev)
}

func (m *eventHandlerEmbed) CallEventHandler(ev MouseEvent) {
	if fnc := m.eventHandlerList[ev.Type]; fnc != nil {
		fnc(ev)
	}
}

// Mit SetHandler wird die Funktion fnc als Handler fuer den Event typ
// registriert. Eine bereits registrierte Funktion wird damit ueberschrieben.
func (m *eventHandlerEmbed) SetHandler(fnc EventHandler, types ...MouseEventType) {
	for _, typ := range types {
		m.eventHandlerList[typ] = fnc
	}
}

// Registriert fnc als Handler fuer den Move-Event.
func (m *eventHandlerEmbed) SetOnMove(fnc EventHandler) {
	m.SetHandler(fnc, TypeMove)
}

// Registriert fnc als Handler fuer den Press-Event.
func (m *eventHandlerEmbed) SetOnPress(fnc EventHandler) {
	m.SetHandler(fnc, TypePress)
}

// Registriert fnc als Handler fuer den Release-Event.
func (m *eventHandlerEmbed) SetOnRelease(fnc EventHandler) {
	m.SetHandler(fnc, TypeRelease)
}

// Registriert fnc als Handler fuer den Drag-Event.
func (m *eventHandlerEmbed) SetOnDrag(fnc EventHandler) {
	m.SetHandler(fnc, TypeDrag)
}

// Registriert fnc als Handler fuer den LongPress-Event.
func (m *eventHandlerEmbed) SetOnLongPress(fnc EventHandler) {
	m.SetHandler(fnc, TypeLongPress)
}

// Registriert fnc als Handler fuer den Wheel-Event.
func (m *eventHandlerEmbed) SetOnWheel(fnc EventHandler) {
	m.SetHandler(fnc, TypeWheel)
}

// Registriert fnc als Handler fuer den Enter-Event.
func (m *eventHandlerEmbed) SetOnEnter(fnc EventHandler) {
	m.SetHandler(fnc, TypeEnter)
}

// Registriert fnc als Handler fuer den Leave-Event.
func (m *eventHandlerEmbed) SetOnLeave(fnc EventHandler) {
	m.SetHandler(fnc, TypeLeave)
}

// Registriert fnc als Handler fuer den Click-Event.
func (m *eventHandlerEmbed) SetOnClick(fnc EventHandler) {
	m.SetHandler(fnc, TypeClick)
}

// Registriert fnc als Handler fuer den DoubleClick-Event.
func (m *eventHandlerEmbed) SetOnDoubleClick(fnc EventHandler) {
	m.SetHandler(fnc, TypeDoubleClick)
}
