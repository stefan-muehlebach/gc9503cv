package main

import (
	"gc9503cv/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
	"golang.org/x/image/font"
	"math"
	"time"
)

var (
	text = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."
)

type FontsAnim struct {
	callbackEmbed
	appletEmbed
	rect                              Rectangle
	t0                                time.Time
	idx, lastIdx                      int
	fontList                          []string
	fontSize, captionFontSize, margin float64
	fontColor, captionFontColor       colors.RGBA
	captionFont                       *fonts.Font
	face, captionFontFace             font.Face
	p0                                Point
}

func NewFontsAnimation(bounds geom.Rectangle[int]) *FontsAnim {
	a := &FontsAnim{}
	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	a.fontList = fonts.Names
	a.fontSize = 18.0
	a.fontColor = colors.Black
	a.captionFont = fonts.SeafordBold
	a.captionFontSize = 40.0
	a.captionFontColor = colors.Black.Alpha(0.3)
	a.margin = 5.0
	return a
}

func (a *FontsAnim) Init() {
	a.t0 = time.Now()
	a.rect = geom.Rectangle[int]{Max: a.bounds.Size()}.ToFloat().Inset(a.margin, a.margin)
	a.lastIdx = -1
	a.captionFontFace, _ = fonts.NewFace(a.captionFont, a.captionFontSize)
	a.p0 = a.rect.SE()
}

func (a *FontsAnim) Update(dt time.Duration) {
	a.idx = int(0.5*time.Since(a.t0).Seconds()) % len(a.fontList)
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
		a.rect.Dx(), 1.3, gg.AlignLeft)
	a.gc.Push()
	a.gc.RotateAbout(-math.Pi/2.0, a.p0.X, a.p0.Y)
	a.gc.SetFontFace(a.captionFontFace)
	a.gc.SetTextColor(a.captionFontColor)
	a.gc.DrawStringAnchored(a.fontList[a.idx], a.p0.X, a.p0.Y, 0.0, 0.0)
	a.gc.Pop()

	//draw.Draw(img, a.bounds.ToInt(), a.gc.Image().(*image.RGBA),
	//	image.Point{}, draw.Over)
}
