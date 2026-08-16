// Das Package touch enthaelt alles, um Events vom Touchscreen bequemer zu
// verarbeiten.
//
// Eine Interaktion auf dem Bildschirm laeuft grundsaetzlich wie folgt ab.
// In Klammern sind Events gesetzt, die optional erzeugt werden:
//
//   TypePress
//    |
//   TypeDrag
//   ...
//    |
//   (TypeLongPress)
//    |
//   TypeDrag
//   ...
//    |
//   TypeRelease
//    |
//   (TypeClick)
//
package point

import (
    "fmt"
    "time"
    "github.com/stefan-muehlebach/gg/geom"
)

// Mit diesem Datentyp werden die unterschiedlichen Event-Arten abgebildet,
// welche durch das Druecken auf den Bildschirm an die GUI-Elemente gesendet
// werden koennen. Es ist Aufgabe der Objekte 'Screen' und 'Window', aus den
// rohen Ereginisse vom Touchscreen (Press, Drag, Release) diese Events zu
// erzeugen.
type Type uint32

const (
	// Mit Einfuehrung der Maus, wird dieses Event (fahren ohne eine Taste
	// gedrueckt zu halten) relevant.
	TypeMove Type = iota
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
    ClickDuration          = 200 * time.Millisecond
    // Analog dazu ist DoubleTapDuration die max. Dauer welche zwischen zwei
    // Tap-Ereignissen vergehen darf, damit sie zusammen als DoubleTap inter-
    // pretiert werden.
    DoubleClickDuration    = 200 * time.Millisecond
    // Drueckt der Benutzer fuer mehr als LongPressThreshold auf ein Objekt,
    // wird ein LongPress-Event erzeugt.
    LongPressThreshold   = 400 * time.Millisecond
    // Fuer die Ereignisse Tap, DoubleTap und auch LongPress darf sich der
    // Finger auf dem TouchScreen nicht zu stark bewegen. Der maximale Abstand
    // zwischen dem Press-Event und der aktuellen Position darf nicht mehr
    // NearThreshold betragen.
    NearThreshold        = 8.0
)

// Damit beim Debuggen klar ist, um welchen Event-Typ es sich handelt,
// implementiert Type das Stringer-Interface.
func (t Type) String() (string) {
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

// In diesem Datentyp ist alles zusammengefasst, was an Touch-Events an die
// applikatorischen Elemente gesendet werden kann.
type Event struct {
    // Der Typ des Events (moegliche Typen: die die Konstanten TypeXXX).
    Type Type
    // Alle Events werden sequenziell durchnumeriert, damit man bspw. das
    // Release-Event mit einem vorgaengigen Press-Event in Verbindung bringen
    // kann.
    SeqNumber int
    // Dieses Feld wird auf true gesetzt, sobald ein LongPressed-Ereignis
    // erkannt wird.
    LongPressed bool
    // In InitTime und InitPos werden Zeitpunkt und Position des Press-Events
    // (des initialen Events) festgehalten.
    InitTime time.Time
    InitPos  geom.Point
    // Wohingegen Time und Pos die Zeit und die Position des aktuellen Events
    // enthalten.
    Time time.Time
    Pos  geom.Point
}

// Fuer das Debugging implementiert Event das Stringer-Interface.
func (evt Event) String() (string) {
    return fmt.Sprintf("%v %d %v %s %v %s %v", evt.Type, evt.SeqNumber,
            evt.LongPressed, evt.InitTime.Format("15:04:05.000000"),
            evt.InitPos, evt.Time.Format("15:04:05.000000"), evt.Pos)
}

// Alle Callback-Handler fuer die Ereignisse vom Touchscreen, muessen folgendes
// Profil aufweisen.
type PointFunction func(evt Event)

// Alle GUI-Elemente, welche ueber den Touchscreen gesteuert werden sollen,
// muesssen diesen Datentyp einbetten. Damit werden auch alle unten
// aufgefuehrten Methoden geerbt und es koennen Handler fuer die diversen
// Touchscreen-Ereignisse hinterlegt werden. Im Array touchFuncList kann
// fuer jedes Ereginis max. eine Funktion hinterlegt werden.
type PointEmbed struct {
    pointFuncList [numEvents]PointFunction
}

// Diese Methode wird durch AdaGui aufgerufen, um ein Touch-Ereignis an
// ein GUI-Element zu senden. Es gibt eine Default-Implementation, welche das
// Event via CallTouchFunc an registrierte Event-Handler sendet.
// Es ist jedoch ueblich, dass ein GUI-Element diese Methode ueberschreibt
// um bspw. visuelle Anpassungen zu machen und dann selber CallTouchFunc
// aufruft.
func (m *PointEmbed) OnInputEvent(evt Event) {
    m.CallPointFunc(evt)
}

// Diese Methode schliesslich ruft (sofern vorhanden) den entsprechenden
// Event-Handler auf. Hier bestuende theoretisch die Moeglichkeit, die Aufrufe
// resp. die Verarbeitung der Events mittels Go-Routinen zu paralellisieren.
// Die notwendige Synchronisation in den GUI-Elementen stelle ich mir jedoch
// ziemlich anspruchsvoll vor...
func (m *PointEmbed) CallPointFunc(evt Event) {
    if fnc := m.pointFuncList[evt.Type]; fnc != nil {
        fnc(evt)
    }
}

// Mit SetPointFunc wird die Funktion fnc als Handler fuer den Event typ
// registriert. Eine bereits registrierte Funktion wird damit ueberschrieben.
func (m *PointEmbed) SetPointFunc(fnc PointFunction, types ...Type) {
    for _, typ := range types {
        m.pointFuncList[typ] = fnc
    }
}

func (m *PointEmbed) SetOnMove(fnc PointFunction) {
    m.SetPointFunc(fnc, TypeMove)
}

// Registriert fnc als Handler fuer den Press-Event.
func (m *PointEmbed) SetOnPress(fnc PointFunction) {
    m.SetPointFunc(fnc, TypePress)
}

// Registriert fnc als Handler fuer den Release-Event.
func (m *PointEmbed) SetOnRelease(fnc PointFunction) {
    m.SetPointFunc(fnc, TypeRelease)
}

// Registriert fnc als Handler fuer den Drag-Event.
func (m *PointEmbed) SetOnDrag(fnc PointFunction) {
    m.SetPointFunc(fnc, TypeDrag)
}

// Registriert fnc als Handler fuer den LongPress-Event.
func (m *PointEmbed) SetOnLongPress(fnc PointFunction) {
    m.SetPointFunc(fnc, TypeLongPress)
}

// Registriert fnc als Handler fuer den Enter-Event.
func (m *PointEmbed) SetOnEnter(fnc PointFunction) {
    m.SetPointFunc(fnc, TypeEnter)
}

// Registriert fnc als Handler fuer den Leave-Event.
func (m *PointEmbed) SetOnLeave(fnc PointFunction) {
    m.SetPointFunc(fnc, TypeLeave)
}

// Registriert fnc als Handler fuer den Tap-Event.
func (m *PointEmbed) SetOnClick(fnc PointFunction) {
    m.SetPointFunc(fnc, TypeClick)
}

// Registriert fnc als Handler fuer den DoubleTap-Event.
func (m *PointEmbed) SetOnDoubleClick(fnc PointFunction) {
    m.SetPointFunc(fnc, TypeDoubleClick)
}

