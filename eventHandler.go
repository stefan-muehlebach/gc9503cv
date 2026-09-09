package main

// Konsequenterweise wurde auch fuer diese Funktionalitaet ein Interface
// gebaut.
type EventHandler interface {
	OnInputEvent(ev MouseEvent)
	SetOnMove(fnc EventHandlerFunc)
	SetOnPress(fnc EventHandlerFunc)
	SetOnRelease(fnc EventHandlerFunc)
	SetOnDrag(fnc EventHandlerFunc)
	SetOnLongPress(fnc EventHandlerFunc)
	SetOnWheel(fnc EventHandlerFunc)
	SetOnEnter(fnc EventHandlerFunc)
	SetOnLeave(fnc EventHandlerFunc)
	SetOnClick(fnc EventHandlerFunc)
	SetOnDoubleClick(fnc EventHandlerFunc)
}

// Alle Callback-Handler fuer die Ereignisse vom Touchscreen, muessen folgendes
// Profil aufweisen.
type EventHandlerFunc func(ev MouseEvent)

// Alle GUI-Elemente, welche ueber den Touchscreen gesteuert werden sollen,
// muesssen diesen Datentyp einbetten. Damit werden auch alle unten
// aufgefuehrten Methoden geerbt und es koennen Handler fuer die diversen
// Touchscreen-Ereignisse hinterlegt werden. Im Array touchFuncList kann
// fuer jedes Ereginis max. eine Funktion hinterlegt werden.
type eventHandlerEmbed struct {
	eventHandlerList [numEvents]EventHandlerFunc
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
func (m *eventHandlerEmbed) setHandler(fnc EventHandlerFunc,
		types ...MouseEventType) {
	for _, typ := range types {
		m.eventHandlerList[typ] = fnc
	}
}

// Registriert fnc als Handler fuer den Move-Event.
func (m *eventHandlerEmbed) SetOnMove(fnc EventHandlerFunc) {
	m.setHandler(fnc, TypeMove)
}

// Registriert fnc als Handler fuer den Press-Event.
func (m *eventHandlerEmbed) SetOnPress(fnc EventHandlerFunc) {
	m.setHandler(fnc, TypePress)
}

// Registriert fnc als Handler fuer den Release-Event.
func (m *eventHandlerEmbed) SetOnRelease(fnc EventHandlerFunc) {
	m.setHandler(fnc, TypeRelease)
}

// Registriert fnc als Handler fuer den Drag-Event.
func (m *eventHandlerEmbed) SetOnDrag(fnc EventHandlerFunc) {
	m.setHandler(fnc, TypeDrag)
}

// Registriert fnc als Handler fuer den LongPress-Event.
func (m *eventHandlerEmbed) SetOnLongPress(fnc EventHandlerFunc) {
	m.setHandler(fnc, TypeLongPress)
}

// Registriert fnc als Handler fuer den Wheel-Event.
func (m *eventHandlerEmbed) SetOnWheel(fnc EventHandlerFunc) {
	m.setHandler(fnc, TypeWheel)
}

// Registriert fnc als Handler fuer den Enter-Event.
func (m *eventHandlerEmbed) SetOnEnter(fnc EventHandlerFunc) {
	m.setHandler(fnc, TypeEnter)
}

// Registriert fnc als Handler fuer den Leave-Event.
func (m *eventHandlerEmbed) SetOnLeave(fnc EventHandlerFunc) {
	m.setHandler(fnc, TypeLeave)
}

// Registriert fnc als Handler fuer den Click-Event.
func (m *eventHandlerEmbed) SetOnClick(fnc EventHandlerFunc) {
	m.setHandler(fnc, TypeClick)
}

// Registriert fnc als Handler fuer den DoubleClick-Event.
func (m *eventHandlerEmbed) SetOnDoubleClick(fnc EventHandlerFunc) {
	m.setHandler(fnc, TypeDoubleClick)
}
