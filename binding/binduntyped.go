package binding

import (
	"log"
	"reflect"
)

// Untyped supports binding a interface{} value.
type Untyped interface {
	DataItem
	Get() interface{}
	Set(interface{})
}

// ExternalUntyped supports binding a interface{} value to an external value.
type ExternalUntyped interface {
	Untyped
	Reload()
}

// NewUntyped returns a bindable interface{} value that is managed internally.
func NewUntyped() Untyped {
	var blank interface{} = nil
	v := &blank
	b := &boundUntyped{val: reflect.ValueOf(v).Elem()}
	b.Init(b)
	return b
}

// BindUntyped returns a new bindable value that controls the contents of the provided interface{} variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
func BindUntyped(v interface{}) ExternalUntyped {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		log.Println("Invalid type passed to BindUntyped, must be a pointer")
		return nil
	}
	if v == nil {
		var blank interface{}
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalUntyped{}
	b.val = reflect.ValueOf(v).Elem()
	b.old = b.val.Interface()
	b.Init(b)
	return b
}

type boundUntyped struct {
	base
	val reflect.Value
}

func (b *boundUntyped) Get() interface{} {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.val.Interface()
}

func (b *boundUntyped) Set(val interface{}) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.val.Interface() == val {
		return
	}
	b.val.Set(reflect.ValueOf(val))
	b.trigger()
}

type boundExternalUntyped struct {
	boundUntyped
	old interface{}
}

func (b *boundExternalUntyped) Set(val interface{}) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.old == val {
		return
	}
	b.val.Set(reflect.ValueOf(val))
	b.old = val
	b.trigger()
}

func (b *boundExternalUntyped) Reload() {
	b.Set(b.val.Interface())
}
