package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
	"golang.org/x/image/font"
)

var (
	text = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."

	MinFontSize      = 10.0
	MaxFontSize      = 48.0
	InitFontSize     = 18.0
	FontColor        = colors.Black
	CaptionFont      = fonts.SeafordBold
	CaptionFontSize  = 40.0
	CaptionFontColor = colors.Black.Alpha(0.3)
	MinLineSpace     = 0.5
	MaxLineSpace     = 2.0
)

type FontsWin struct {
	// eventHandlerEmbed
	windowEmbed
	rect              Rectangle
	idx               int
	lineSpace         float64
	fontList          []string
	fontSize, margin  float64
	fontColor         colors.RGBA
	face, captionFace font.Face
	fontHasChanged    bool
	p0, p1            Point
}

func NewFontsWindow(bounds geom.Rectangle[int]) *FontsWin {
	a := &FontsWin{}
	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	a.lineSpace = 1.0
	a.fontList = fonts.Names
	a.fontSize = InitFontSize
	a.fontColor = FontColor
	a.margin = 5.0

	a.Root = NewGroup()
	a.Root.SetSize(bounds.ToFloat().Size())

	a.Root.SetOnClick(func(ev InputEvent) {
		if ev.Button.IsSet(LeftButton) {
			a.idx = (a.idx + 1) % len(a.fontList)
			a.fontHasChanged = true
		}
		if ev.Button.IsSet(RightButton) {
			a.idx = (a.idx - 1 + len(a.fontList)) % len(a.fontList)
			a.fontHasChanged = true
		}
	})

	a.Root.SetOnWheel(func(ev InputEvent) {
		t := float64(ev.Wheel) / 10.0
		if ev.Modifiers.IsSet(ModShift) {
			if (a.fontSize <= MinFontSize && ev.Wheel < 0) ||
				(a.fontSize >= MaxFontSize && ev.Wheel > 0) {
				return
			}
			a.fontSize += float64(ev.Wheel)
			a.fontHasChanged = true
		} else {
			if (a.lineSpace <= MinLineSpace && t < 0.0) ||
				(a.lineSpace >= MaxLineSpace && t > 0.0) {
				return
			}
			a.lineSpace += t
		}
	})

	return a
}

func (a *FontsWin) Init() {
	log.Printf("Interaction:")
	log.Printf("  LMB         - Next Font in List")
	log.Printf("  RMB         - Previous Font in List")
	log.Printf("  Wheel       - Change the line space (0.5-2.0)")
	log.Printf("  Wheel+SHIFT - Change the font size (12-22)")

	a.rect = geom.Rectangle[int]{Max: a.bounds.Size()}.ToFloat().Inset(a.margin, a.margin)
	a.idx = 0
	a.captionFace, _ = fonts.NewFace(CaptionFont, CaptionFontSize)
	a.face, _ = fonts.NewFace(fonts.Map[a.fontList[a.idx]], a.fontSize)
	a.p0 = a.rect.SE()
	a.p1 = a.rect.SW()
}

func (a *FontsWin) Update(t time.Time) {
	if a.fontHasChanged {
		a.face, _ = fonts.NewFace(fonts.Map[a.fontList[a.idx]], a.fontSize)
		a.fontHasChanged = false
	}
}

func (a *FontsWin) Refresh() {
	a.gc.Clear(colors.BurlyWood)
	a.gc.SetFontFace(a.face)
	a.gc.SetTextColor(a.fontColor)
	a.gc.DrawStringWrapped(text, a.rect.Min.X, a.rect.Min.Y, 0.0, 0.0,
		a.rect.Dx(), a.lineSpace, gg.AlignLeft)

	a.gc.SetFontFace(a.captionFace)
	a.gc.SetTextColor(CaptionFontColor)
	a.gc.Push()
	a.gc.RotateAbout(-math.Pi/2.0, a.p0.X, a.p0.Y)
	a.gc.DrawStringAnchored(a.fontList[a.idx], a.p0.X, a.p0.Y, 0.0, 0.0)
	a.gc.Pop()
	txt := fmt.Sprintf("%.0f | %.1f", a.fontSize, a.lineSpace)
	a.gc.DrawStringAnchored(txt, a.p1.X, a.p1.Y, 0.0, 0.0)
}
