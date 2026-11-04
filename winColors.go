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
	BackColor  = colors.DarkSlateGrey.Dark(0.6)
	SquareSize = 60.0
	HalveSize  = SquareSize / 2.0
)

type GoColorWin struct {
	windowEmbed
	t0        time.Time
	patIdx    int
	isManual  bool
	colorList []string
}

func NewGoColorWindow(bounds geom.Rectangle[int]) *GoColorWin {
	a := &GoColorWin{}
	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	a.patIdx = 0
	a.colorList = colors.Groups[colors.GoColors]

	a.Root = NewGroup()
	a.Root.SetSize(bounds.ToFloat().Size())

	a.Root.SetOnClick(func(ev InputEvent) {
		if ev.Button.IsSet(RightButton) {
			a.isManual = false
		}
	})
	a.Root.SetOnWheel(func(ev InputEvent) {
		a.isManual = true
		if ev.Wheel > 0 {
			a.patIdx = (a.patIdx + ev.Wheel) % NumPattern
		} else if ev.Wheel < 0 {
			a.patIdx = (a.patIdx + NumPattern + ev.Wheel) % NumPattern
		}
	})
	return a
}

func (a *GoColorWin) Init() {
	log.Printf("Interaction:")
	log.Printf("-----------------------------------------------------------")
	log.Printf("RMB   - Return to automatic switching")
	log.Printf("Wheel - Surf through the patterns")
	a.gc.SetLineWidth(2.0)
	a.gc.SetLineCapRound()
	a.gc.SetLineJoinRound()
	a.t0 = time.Now()
}

func (a *GoColorWin) Update(t time.Time) {
	if a.isManual {
		return
	}
	if time.Since(a.t0) > time.Second {
		a.patIdx = (a.patIdx + 1) % NumPattern
		a.t0 = time.Now()
	}
}

func (a *GoColorWin) Refresh() {
	a.gc.Clear(BackColor)
	idx := 0
	for row := range a.bounds.Dy() / int(SquareSize) {
		y := float64(row) * SquareSize
		for col := range a.bounds.Dx() / int(SquareSize) {
			x := float64(col) * SquareSize
			DrawerList[a.patIdx](a.gc, x, y, colors.Map[a.colorList[idx]])
			idx = (idx + 1) % len(a.colorList)
		}
	}
}

//-----------------------------------------------------------------------

type PatternDrawer func(gc *gg.Context, x, y float64, col colors.RGBA)

var (
	DrawerList = []PatternDrawer{
		// Draw a filled square
		func(gc *gg.Context, x, y float64, col colors.RGBA) {
			gc.SetFillColor(col)
			gc.DrawRectangle(x, y, SquareSize, SquareSize)
			gc.Fill()
		},
		// Draw a filled rounded square
		func(gc *gg.Context, x, y float64, col colors.RGBA) {
			gc.SetFillColor(col)
			gc.DrawRoundedRectangle(x, y, SquareSize, SquareSize, SquareSize/3)
			gc.Fill()
		},
		// Draw rhombuses instead of squares
		func(gc *gg.Context, x, y float64, col colors.RGBA) {
			overLength := SquareSize / 4.0
			gc.SetFillColor(col)
			gc.MoveTo(x+HalveSize, y-overLength)
			gc.LineTo(x+SquareSize+overLength, y+HalveSize)
			gc.LineTo(x+HalveSize, y+SquareSize+overLength)
			gc.LineTo(x-overLength, y+HalveSize)
			gc.ClosePath()
			gc.Fill()
		},
		// Draw a filled circle
		func(gc *gg.Context, x, y float64, col colors.RGBA) {
			gc.SetFillColor(col)
			gc.DrawCircle(x+HalveSize, y+HalveSize, HalveSize-1.0)
			gc.Fill()
		},
		// Draw a filled square with a black hole in it
		func(gc *gg.Context, x, y float64, col colors.RGBA) {
			gc.SetFillColor(col)
			gc.DrawRectangle(x, y, SquareSize, SquareSize)
			gc.Fill()
			gc.SetFillColor(BackColor)
			gc.DrawCircle(x+HalveSize, y+HalveSize, HalveSize-1.0)
			gc.Fill()
		},
		// Draw a quarter of a circle
		func(gc *gg.Context, x, y float64, col colors.RGBA) {
			radius := SquareSize
			gc.SetFillColor(col)
			gc.MoveTo(x, y)
			gc.LineTo(x+radius, y)
			gc.QuadraticTo(x+radius, y+radius, x, y+radius)
			gc.ClosePath()
			gc.Fill()
		},
		// Draw a circle with a bold pen
		func(gc *gg.Context, x, y float64, col colors.RGBA) {
			gc.SetFillColor(col)
			gc.DrawCircle(x+HalveSize, y+HalveSize, HalveSize-1.0)
			gc.Fill()
			gc.SetFillColor(BackColor)
			gc.DrawCircle(x+HalveSize, y+HalveSize, HalveSize-15.0)
			gc.Fill()
		},
		// Draw multiple circles within each other
		func(gc *gg.Context, x, y float64, col colors.RGBA) {
			step := 10.0
			lineWidth := step / 2.0
			gc.SetLineWidth(lineWidth)
			gc.SetLineColor(col)
			for radius := HalveSize - lineWidth/2.0 - 1.0; radius >= step/2.0; radius -= step {
				gc.DrawCircle(x+HalveSize, y+HalveSize, radius)
				gc.Stroke()
			}
		},
		// Draw a spiral
		func(gc *gg.Context, x, y float64, col colors.RGBA) {
			step := 10.0
			lineWidth := step / 2.0
			y1 := HalveSize
			y2 := HalveSize + step/2.0
			gc.SetLineWidth(lineWidth)
			gc.SetLineColor(col)
			for radius := HalveSize - lineWidth/2.0 - 1.0; radius >= step/2.0; radius -= step {
				gc.DrawArc(x+HalveSize, y+y1, radius, -math.Pi/2.0, math.Pi/2.0)
				gc.DrawArc(x+HalveSize, y+y2, radius-step/2.0, math.Pi/2.0, 3.0*math.Pi/2.0)
			}
			gc.Stroke()
		},
	}

	NumPattern = len(DrawerList)
)
