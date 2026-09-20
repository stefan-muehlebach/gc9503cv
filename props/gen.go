//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"
	"path"
	"runtime"
	"text/template"

	"github.com/stefan-muehlebach/gc9503cv/props"
)

const colorPropTempl = `
func (pe *PropertyEmbed) {{ .Name }}() (colors.RGBA) {
    return pe.prop.Color({{ .Name }})
}
func (pe *PropertyEmbed) Set{{ .Name }}(c colors.RGBA) {
    pe.prop.SetColor({{ .Name }}, c)
}
`
const fontPropTempl = `
func (pe *PropertyEmbed) {{ .Name }}() (*fonts.Font) {
    return pe.prop.Font({{ .Name }})
}
func (pe *PropertyEmbed) Set{{ .Name }}(f *fonts.Font) {
    pe.prop.SetFont({{ .Name }}, f)
}
`
const sizePropTempl = `
func (pe *PropertyEmbed) {{ .Name }}() (float64) {
    return pe.prop.Size({{ .Name }})
}
func (pe *PropertyEmbed) Set{{ .Name }}(s float64) {
    pe.prop.SetSize({{ .Name }}, s)
}
`

type ColorPropInfo struct {
	Templ *template.Template
	List  []props.ColorPropertyName
}
type FontPropInfo struct {
	Templ *template.Template
	List  []props.FontPropertyName
}
type SizePropInfo struct {
	Templ *template.Template
	List  []props.SizePropertyName
}

type PropName struct {
	Name string
}

func main() {
	_, dirName, _, _ := runtime.Caller(0)
	filePath := path.Join(path.Dir(dirName), "propsEmbed.go")
	os.Remove(filePath)
	propFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
		return
	}
	defer propFile.Close()
	propFile.WriteString(`//
// THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDTI
//

package props

import (
    "github.com/stefan-muehlebach/gg/colors"
    "github.com/stefan-muehlebach/gg/fonts"
)

type PropertyEmbed struct {
    prop *Properties
}

func (pe *PropertyEmbed) Init(parent *Properties) {
    pe.prop = NewProperties(parent)
}
func (pe *PropertyEmbed) InitProp(name string) {
    pe.prop = PropsMap[name]
}

`)

	colorProps := ColorPropInfo{
		Templ: template.Must(template.New("color").Parse(colorPropTempl)),
		List:  props.EmbedColorProps,
	}
	fontProps := FontPropInfo{
		Templ: template.Must(template.New("font").Parse(fontPropTempl)),
		List:  props.EmbedFontProps,
	}
	sizeProps := SizePropInfo{
		Templ: template.Must(template.New("size").Parse(sizePropTempl)),
		List:  props.EmbedSizeProps,
	}

	for _, pr := range colorProps.List {
		pn := PropName{Name: pr.String()}
		if err := colorProps.Templ.Execute(propFile, pn); err != nil {
			log.Fatalf("Unable to write file: %v", err)
		}
	}
	for _, pr := range fontProps.List {
		pn := PropName{Name: pr.String()}
		if err := fontProps.Templ.Execute(propFile, pn); err != nil {
			log.Fatalf("Unable to write file: %v", err)
		}
	}
	for _, pr := range sizeProps.List {
		pn := PropName{Name: pr.String()}
		if err := sizeProps.Templ.Execute(propFile, pn); err != nil {
			log.Fatalf("Unable to write file: %v", err)
		}
	}
}
