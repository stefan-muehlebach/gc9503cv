package main

import (
	"log"
	"math"
	"time"

	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
)

var (
	BackgroundColor = colors.DimGray.Dark(0.6)
	NumPattern = 8
)

type GoColorAnim struct {
	eventHandlerEmbed
	appletEmbed
	t0        time.Time
	patIdx, numPat int
	colorList []string
}

func NewGoColorAnimation(bounds geom.Rectangle[int]) *GoColorAnim {
	a := &GoColorAnim{}
	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	a.patIdx = 0
	a.colorList = colors.Groups[colors.GoColors]

	a.SetOnClick(func(ev MouseEvent) {
		if ev.Button.IsSet(LeftButton) {
			a.patIdx = (a.patIdx + 1) % NumPattern
		}
		if ev.Button.IsSet(RightButton) {
			a.patIdx = (a.patIdx - 1 + NumPattern) % NumPattern
		}
	})
	return a
}

func (a *GoColorAnim) Init() {
	log.Printf("Interaction")
	log.Printf("-----------------------------------------------------------")
	log.Printf("LMB - Jump to the next pattern")
	log.Printf("RMB - Jump to the previous pattern")
	a.gc.SetLineWidth(2.0)
	a.gc.SetLineCapRound()
	a.gc.SetLineJoinRound()
	a.t0 = time.Now()
}

func (a *GoColorAnim) Update(dt time.Duration) {
	//if time.Since(a.t0) > 2*time.Second {
	//	a.patIdx = (a.patIdx + 1) % a.numPat
	//	a.t0 = time.Now()
	//}
}

func (a *GoColorAnim) Refresh() {
	a.gc.Clear(BackgroundColor)
	idx := 0
	for row := range a.bounds.Dy() / 60 {
		y := float64(row * 60)
		for col := range a.bounds.Dx() / 60 {
			x := float64(col * 60)
			DrawPattern(a.gc, x, y, a.patIdx, colors.Map[a.colorList[idx]])
			idx = (idx + 1) % len(a.colorList)
		}
	}
	//draw.Draw(img, a.bounds.ToInt(), a.gc.Image().(*image.RGBA),
	//	image.Point{}, draw.Over)
}

//-----------------------------------------------------------------------

func DrawPattern(gc *gg.Context, x, y float64, typ int, col colors.RGBA) {
	switch typ {
	case 0:
		gc.SetFillColor(col)
		gc.DrawRectangle(x, y, 60.0, 60.0)
		gc.Fill()
	case 1:
		gc.SetFillColor(col)
		gc.DrawCircle(x+30.0, y+30.0, 30.0-1.0)
		gc.Fill()
	case 2:
		gc.SetFillColor(col)
		gc.DrawRectangle(x, y, 60.0, 60.0)
		gc.Fill()
		gc.SetFillColor(BackgroundColor)
		gc.DrawCircle(x+30.0, y+30.0, 30.0-1.0)
		gc.Fill()
	case 3:
		gc.SetFillColor(col)
		gc.DrawCircle(x+30.0, y+30.0, 30.0)
		gc.Fill()
		gc.SetFillColor(BackgroundColor)
		gc.DrawCircle(x+30.0, y+30.0, 30.0-15.0)
		gc.Fill()
	case 4:
		for rad:=30.0; rad>5.0; rad-=10.0 {
			gc.SetFillColor(col)
			gc.DrawCircle(x+30.0, y+30.0, rad)
			gc.Fill()
			gc.SetFillColor(BackgroundColor)
			gc.DrawCircle(x+30.0, y+30.0, rad-5.0)
			gc.Fill()
		}
	case 5:
		overLength := 15.0
		gc.SetFillColor(col)
		gc.MoveTo(x+30.0, y-overLength)
		gc.LineTo(x+60.0+overLength, y+30.0)
		gc.LineTo(x+30.0, y+60.0+overLength)
		gc.LineTo(x-overLength, y+30.0)
		gc.ClosePath()
		gc.Fill()
	case 6:
		radius := 60.0
		gc.SetFillColor(col)
		gc.MoveTo(x, y)
		gc.LineTo(x+radius, y)
		gc.QuadraticTo(x+radius, y+radius, x, y+radius)
		gc.ClosePath()
		gc.Fill()
	case 7:
		step := 10.0
		x1 := 30.0
		x2 := 30.0 + step/2.0
		gc.SetLineWidth(step/2.0)
		gc.SetLineColor(col)
		for radius:=27.0; radius>=step/2.0; radius-=step {
			gc.DrawArc(x+x1, y+30.0, radius, -math.Pi, 0)
			gc.DrawArc(x+x2, y+30.0, radius-step/2.0, 0, math.Pi)
		}
		gc.Stroke()
	}
}
