//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"
	"path"
	"runtime"
	"text/template"
)

const itemBindTemplate = `
// {{ .Name }} supports binding a {{ .Type }} value.
type {{ .Name }} interface {
    DataItem
    Get() ({{ .Type }})
    Set({{ .Type }})
}

// External{{ .Name }} supports binding a {{ .Type }} value to an external value.
type External{{ .Name }} interface {
    {{ .Name }}
    Reload()
}

// New{{ .Name }} returns a bindable {{ .Type }} value that is managed internally.
func New{{ .Name }}() {{ .Name }} {
    var blank {{ .Type }} = {{ .Default }}
    b := &bound{{ .Name }}{val: &blank}
    b.Init(b)
    return b
}

// Bind{{ .Name }} returns a new bindable value that controls the contents of the provided {{ .Type }} variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
func Bind{{ .Name }}(v *{{ .Type }}) External{{ .Name }} {
    if v == nil {
        var blank {{ .Type }} = {{ .Default }}
        v = &blank // never allow a nil value pointer
    }
    b := &boundExternal{{ .Name }}{}
    b.val = v
    b.old = *v
    b.Init(b)
    return b
}

type bound{{ .Name }} struct {
    base
    val *{{ .Type }}
}

func (b *bound{{ .Name }}) Get() ({{ .Type }}) {
    b.lock.RLock()
    defer b.lock.RUnlock()
    if b.val == nil {
        return {{ .Default }}
    }
    return *b.val
}

func (b *bound{{ .Name }}) Set(val {{ .Type }}) {
    b.lock.Lock()
    defer b.lock.Unlock()
    {{- if eq .Comparator "" }}
    if *b.val == val {
        return
    }
    {{- else }}
    if {{ .Comparator }}(*b.val, val) {
        return
    }
    {{- end }}
    *b.val = val
    b.trigger()
}

type boundExternal{{ .Name }} struct {
    bound{{ .Name }}
    old {{ .Type }}
}

func (b *boundExternal{{ .Name }}) Set(val {{ .Type }}) {
    b.lock.Lock()
    defer b.lock.Unlock()
    {{- if eq .Comparator "" }}
    if b.old == val {
        return
    }
    {{- else }}
    if {{ .Comparator }}(b.old, val) {
        return
    }
    {{- end }}
    *b.val = val
    b.old = val
    b.trigger()
}

func (b *boundExternal{{ .Name }}) Reload() {
    b.Set(*b.val)
}
`

const toStringTemplate = `
type stringFrom{{ .Name }} struct {
    base
{{ if .Format }}
    format string
{{ end }}
    from {{ .Name }}
}

// {{ .Name }}ToString creates a binding that connects a {{ .Name }} data item to a String.
// Changes to the {{ .Name }} will be pushed to the String and setting the string will parse and set the
// {{ .Name }} if the parse was successful.
//
func {{ .Name }}ToString(v {{ .Name }}) String {
    str := &stringFrom{{ .Name }}{from: v}
    v.AddListener(str)
    return str
}
{{ if .Format }}
// {{ .Name }}ToStringWithFormat creates a binding that connects a {{ .Name }} data item to a String and is
// presented using the specified format. Changes to the {{ .Name }} will be pushed to the String and setting
// the string will parse and set the {{ .Name }} if the string matches the format and its parse was successful.
//
func {{ .Name }}ToStringWithFormat(v {{ .Name }}, format string) String {
    if format == "{{ .Format }}" { // Same as not using custom formatting.
        return {{ .Name }}ToString(v)
    }

    str := &stringFrom{{ .Name }}{from: v, format: format}
    v.AddListener(str)
    return str
}
{{ end }}
func (s *stringFrom{{ .Name }}) Get() (string) {
    val := s.from.Get()
{{ if .ToString }}
    return {{ .ToString }}(val)
{{- else }}
    if s.format != "" {
        return fmt.Sprintf(s.format, val)
    }
    return format{{ .Name }}(val)
{{- end }}
}

func (s *stringFrom{{ .Name }}) Set(str string) {
{{- if .FromString }}
    val := {{ .FromString }}(str)
{{ else }}
    var val {{ .Type }}
    if s.format != "" {
        safe := stripFormatPrecision(s.format)
        n, _ := fmt.Sscanf(str, safe+" ", &val) // " " denotes match to end of string
        if n != 1 {
            return
        }
    } else {
        new, _ := parse{{ .Name }}(str)
        val = new
    }
{{ end }}
    old := s.from.Get()
    if val == old {
        return
    }
    s.from.Set(val)
    s.DataChanged(s.super)
}

func (s *stringFrom{{ .Name }}) DataChanged(data DataItem) {
    s.lock.RLock()
    defer s.lock.RUnlock()
    s.trigger()
}
`

const fromStringTemplate = `
type stringTo{{ .Name }} struct {
    base
{{ if .Format }}
    format string
{{ end }}
    from String
}

// StringTo{{ .Name }} creates a binding that connects a String data item to a {{ .Name }}.
// Changes to the String will be parsed and pushed to the {{ .Name }} if the parse was successful, and setting
// the {{ .Name }} update the String binding.
//
func StringTo{{ .Name }}(str String) {{ .Name }} {
    v := &stringTo{{ .Name }}{from: str}
    str.AddListener(v)
    return v
}
{{ if .Format }}
// StringTo{{ .Name }}WithFormat creates a binding that connects a String data item to a {{ .Name }} and is
// presented using the specified format. Changes to the {{ .Name }} will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the {{ .Name }} will push a formatted value
// into the String.
//
func StringTo{{ .Name }}WithFormat(str String, format string) {{ .Name }} {
    if format == "{{ .Format }}" { // Same as not using custom format.
        return StringTo{{ .Name }}(str)
    }

    v := &stringTo{{ .Name }}{from: str, format: format}
    str.AddListener(v)
    return v
}
{{ end }}
func (s *stringTo{{ .Name }}) Get() ({{ .Type }}) {
    str := s.from.Get()
    if str == "" {
        return {{ .Default }}
    }
{{ if .FromString }}
    return {{ .FromString }}(str)
{{- else }}
    var val {{ .Type }}
    if s.format != "" {
        n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
        if err != nil {
            return {{ .Default }}
        }
        if n != 1 {
            return {{ .Default }}
        }
    } else {
        new, err := parse{{ .Name }}(str)
        if err != nil {
            return {{ .Default }}
        }
        val = new
    }
    return val
{{- end }}
}

func (s *stringTo{{ .Name }}) Set(val {{ .Type }}) {
{{- if .ToString }}
    str := {{ .ToString }}(val)
{{- else }}
    var str string
    if s.format != "" {
        str = fmt.Sprintf(s.format, val)
    } else {
        str = format{{ .Name }}(val)
    }
{{ end }}
    old := s.from.Get()
    if str == old {
        return
    }
    s.from.Set(str)
    s.DataChanged(s.super)
}

func (s *stringTo{{ .Name }}) DataChanged(data DataItem) {
    s.lock.RLock()
    defer s.lock.RUnlock()
    s.trigger()
}
`

type bindValues struct {
	Name, Type, Default  string
	Format               string
	FromString, ToString string // function names...
	Comparator           string // comparator function name
}

func newFile(name string) (*os.File, error) {
	_, dirname, _, _ := runtime.Caller(0)
	filepath := path.Join(path.Dir(dirname), name+".go")
	os.Remove(filepath)
	f, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("Unable to open file %s: %v", f.Name(), err)
		return nil, err
	}

	f.WriteString(`// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding
`)
	return f, nil
}

func writeFile(f *os.File, t *template.Template, d interface{}) {
	if err := t.Execute(f, d); err != nil {
		log.Fatalf("Unable to write file %s: %v", f.Name(), err)
	}
}

func main() {
	itemFile, err := newFile("binditems")
	if err != nil {
		return
	}
	defer itemFile.Close()
	itemFile.WriteString(`
import (
    "bytes"
)
`)

	convertFile, err := newFile("convert")
	if err != nil {
		return
	}
	defer convertFile.Close()
	convertFile.WriteString(`
import (
    "fmt"
)
`)

	item := template.Must(template.New("item").Parse(itemBindTemplate))
	fromString := template.Must(template.New("fromString").Parse(fromStringTemplate))
	toString := template.Must(template.New("toString").Parse(toStringTemplate))

	binds := []bindValues{
		bindValues{Name: "Bool", Type: "bool", Default: "false", Format: "%t"},
		bindValues{Name: "Bytes", Type: "[]byte", Default: "nil", Comparator: "bytes.Equal"},
		bindValues{Name: "Float", Type: "float64", Default: "0.0", Format: "%f"},
		bindValues{Name: "Int", Type: "int", Default: "0", Format: "%d"},
		bindValues{Name: "Rune", Type: "rune", Default: "rune(0)"},
		bindValues{Name: "String", Type: "string", Default: "\"\""},
	}
	for _, b := range binds {
		writeFile(itemFile, item, b)

		if b.Format != "" || b.ToString != "" {
			writeFile(convertFile, toString, b)
		}
	}
	for _, b := range binds {
		if b.Format != "" || b.ToString != "" {
			writeFile(convertFile, fromString, b)
		}
	}
}
