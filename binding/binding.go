//go:generate go run gen.go

// Enthält eine Sammlung von gebundenen Datentypen.
package binding

import (
	"log"
	"sync"
)

// Die init-Funktion ist vorallem hier, damit der Import des Log-Packages
// nicht bei jedem Aktivieren von Debug-Meldungen (aus)kommentiert werden muss.
func init() {
	log.SetFlags(log.Lmicroseconds | log.Lmsgprefix)
	log.SetPrefix(": ")
}

// Mit DataItem ist das minimale Interface definiert welches ein Bind-Objekt
// implementieren muss. Einem Bind-Objekt können einerseits Objekte
// hinzugefügt werden, die das DataListener Interface implementieren,
// andererseits aber auch direkt Funktionen vom Typ CallbackFunc. Diese beiden
// Möglichkeiten können auch zusammen, d.h. synchron verwendet werden.
type DataItem interface {
	AddListener(l DataListener)
	RemoveListener(l DataListener)
	AddCallback(f CallbackFunc)
	RemoveCallback(f CallbackFunc)
}

// Alle Typen, welche über die Aenderungen von Bind-Objekten informiert werden
// wollen, müssen das DataListener-Interface implementieren. Im Wesentlichen
// eine einzige Methode (DataChanged), welche als Argument das veränderte
// Bind-Objekt hat.
type DataListener interface {
	DataChanged(data DataItem)
}

// Mit NewDataListener kann einfach und bequem ein neues Listener-Objekt
// erzeugt werden.
func NewDataListener(fn func(data DataItem)) DataListener {
	return &listener{fn}
}

// Der (private) Typ listener ist nun eine der möglichen Implementationen des
// DataListener-Interfaces...
type listener struct {
	callback func(data DataItem)
}

// ... und implementiert logischerweise das DataListener Interface.
func (l *listener) DataChanged(data DataItem) {
	l.callback(data)
}

// Neben den DataListeners gibt es noch eine schlanke Variante seinen Code bei
// Veränderungen eines Bind-Objektes aufrufen zu lassen: man kann einfach eine
// Funktion/Methode, welche dem Typ CallbackFunc entspricht, mit AddCallback
// hinzufügen.
type CallbackFunc func(data DataItem)

// base ist der Basistyp, welche die Methoden des DataItem-Interfaces
// implementiert. Er ist nicht öffentlich, sondern wird von den weiter
// unten gezeigten konkreten Bind-Type verwendet.
type base struct {
	super DataItem
	// listeners und callbacks sind synchronisierte Maps, in welchen die
	// DataListener-Objekte oder Callback-Funktionen hinterlegt werden.
	listeners sync.Map
	callbacks sync.Map
	// Der Zugriff auf weitere Strukturen dieses Typs wird über das Mutex
	// lock gesteuert.
	lock sync.RWMutex
}

// Mit Init werden wichtige Initialisierungen vorgenommen. Init ist bei jeder
// Erstellung eines Bind-Objektes aufzurufen!
func (b *base) Init(super DataItem) {
	b.super = super
}

// AddListener fügt ein DataListener hinzu und ruft die DataChanged-Methode
// auch gleich auf!
func (b *base) AddListener(l DataListener) {
	b.listeners.Store(l, true)
	l.DataChanged(b.super)
}

// Mit RemoveListener können DataListener wieder entfernt werden.
func (b *base) RemoveListener(l DataListener) {
	b.listeners.Delete(l)
}

// AddCallback fügt eine Callback-Funktion hinzu und ruft diese auch gleich
// das erste mal auf.
func (b *base) AddCallback(f CallbackFunc) {
	b.callbacks.Store(&f, true)
	f(b.super)
}

// Mit RemoveCallback schliesslich können Callback-Funktionen wieder
// entfernt werden.
func (b *base) RemoveCallback(f CallbackFunc) {
	b.callbacks.Delete(&f)
}

// Die (private) Methode trigger wird immer dann aufgerufen, wenn sich der
// Wert des Bind-Objektes ändert. Diese Methode ruft die registrierten
// Listener-, resp. Callback-Methoden auf.
func (b *base) trigger() {
	b.listeners.Range(func(key, _ any) bool {
		go key.(DataListener).DataChanged(b.super)
		return true
	})
	b.callbacks.Range(func(f, _ any) bool {
		go (*f.(*CallbackFunc))(b.super)
		return true
	})
}
