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
)

type FontsAnim struct {
	// eventHandlerEmbed
	windowEmbed
	rect                              Rectangle
	idx, lastIdx                      int
	lineSpace                         float64
	fontList                          []string
	fontSize, captionFontSize, margin float64
	fontColor, captionFontColor       colors.RGBA
	captionFont                       *fonts.Font
	face, captionFontFace             font.Face
	p0, p1                            Point
}

func NewFontsAnimation(bounds geom.Rectangle[int]) *FontsAnim {
	a := &FontsAnim{}
	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	a.lineSpace = 1.0
	a.fontList = fonts.Names
	a.fontSize = 18.0
	a.fontColor = colors.Black
	a.captionFont = fonts.SeafordBold
	a.captionFontSize = 40.0
	a.captionFontColor = colors.Black.Alpha(0.3)
	a.margin = 5.0

	a.Root = NewGroup()
	a.Root.SetSize(bounds.ToFloat().Size())

	a.Root.SetOnClick(func(ev InputEvent) {
		if ev.Button.IsSet(LeftButton) {
			a.idx = (a.idx + 1) % len(a.fontList)
		}
		if ev.Button.IsSet(RightButton) {
			a.idx = (a.idx - 1 + len(a.fontList)) % len(a.fontList)
		}
	})

	a.Root.SetOnWheel(func(ev InputEvent) {
		t := float64(ev.Wheel) / 100.0
		a.lineSpace = 1.0 + t
	})

	return a
}

func (a *FontsAnim) Init() {
	log.Printf("Interaction:")
	log.Printf("  LMB   - Next Font in List")
	log.Printf("  RMB   - Previous Font in List")
	log.Printf("  Wheel - Change the line space (1.0..2.0)")

	a.rect = geom.Rectangle[int]{Max: a.bounds.Size()}.ToFloat().Inset(a.margin, a.margin)
	a.idx = 0
	a.lastIdx = -1
	a.captionFontFace, _ = fonts.NewFace(a.captionFont, a.captionFontSize)
	a.p0 = a.rect.SE()
	a.p1 = a.rect.SW()
}

func (a *FontsAnim) Update(dt time.Duration) {
	if a.idx != a.lastIdx {
		a.face, _ = fonts.NewFace(fonts.Map[a.fontList[a.idx]], a.fontSize)
		a.lastIdx = a.idx
	}
}

func (a *FontsAnim) Refresh() {
	a.gc.Clear(colors.BurlyWood)
	a.gc.SetFontFace(a.face)
	a.gc.SetTextColor(a.fontColor)
	a.gc.DrawStringWrapped(text, a.rect.Min.X, a.rect.Min.Y, 0.0, 0.0,
		a.rect.Dx(), a.lineSpace, gg.AlignLeft)

	a.gc.SetFontFace(a.captionFontFace)
	a.gc.SetTextColor(a.captionFontColor)
	a.gc.Push()
	a.gc.RotateAbout(-math.Pi/2.0, a.p0.X, a.p0.Y)
	a.gc.DrawStringAnchored(a.fontList[a.idx], a.p0.X, a.p0.Y, 0.0, 0.0)
	a.gc.Pop()
	txt := fmt.Sprintf("%.2f", a.lineSpace)
	a.gc.DrawStringAnchored(txt, a.p1.X, a.p1.Y, 0.0, 0.0)
}
