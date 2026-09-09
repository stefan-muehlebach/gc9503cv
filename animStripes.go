package main

import (
	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"math"
	"time"
)

type StripeAnim struct {
	eventHandlerEmbed
	windowEmbed
	colorList   []colors.RGBA
	stripeWidth float64
	off, size   geom.Point[int]
	dt          float64
}

func NewStripeAnimation(bounds geom.Rectangle[int],
	colorList []colors.RGBA) *StripeAnim {
	a := &StripeAnim{}

	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	a.colorList = make([]colors.RGBA, len(colorList))
	copy(a.colorList, colorList)
	a.off = bounds.Min
	a.size = bounds.Size()
	a.stripeWidth = float64(a.size.X) / float64(len(a.colorList))
	a.dt = 0.0

	return a
}

func (a *StripeAnim) Init() {
	a.dt = 0.0
}

func (a *StripeAnim) Update(dt time.Duration) {
	a.dt += 20.0
}

func (a *StripeAnim) Refresh() {
	for row := range a.size.Y {
		t := (a.dt + float64(row)) / float64(a.size.Y-1)
		t = math.Mod(t, 2.0)
		if t >= 1.0 {
			t = 2.0 - t
		}
		val := byte(255.0 * t)
		for col := range a.size.X {
			idx := int(math.Floor(float64(col) / a.stripeWidth))
			color := a.colorList[idx]
			if color.R == 0xFF {
				color.R = val
			}
			if color.G == 0xFF {
				color.G = val
			}
			if color.B == 0xFF {
				color.B = val
			}
			a.gc.SetPixel(col, row, color)
		}
	}
}
