package main

import (
	"image"
	"image/draw"
	"time"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"gc9503cv/gc9503cv/geom"
)

type GoColorAnim struct {
	callbackEmbed
    gc *gg.Context
    bounds geom.Rectangle[int]
	//rect geom.Rectangle[int]
	t0 time.Time
	patIdx, prevPatIdx int
	colorList []string
}

func NewGoColorAnimation(bounds geom.Rectangle[int]) *GoColorAnim {
	a := &GoColorAnim{}
    a.bounds = bounds
    a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	//a.rect = rect
	a.patIdx = -1
	a.prevPatIdx = -1
	a.colorList = colors.Groups[colors.GoColors]
	return a
}

func (a *GoColorAnim) Init() {
	a.gc.SetLineWidth(2.0)
	a.gc.SetLineCapRound()
	a.gc.SetLineJoinRound()
	a.t0 = time.Now()
}

func (a *GoColorAnim) Update(dt time.Duration) {
	l := time.Since(a.t0)
	switch {
	case l < 2*time.Second:
		a.patIdx = 0
	case l < 4*time.Second:
		a.patIdx = 1
	case l < 6*time.Second:
		a.patIdx = 2
	case l < 8*time.Second:
		a.patIdx = 3
	default:
		a.t0 = time.Now()
	}
}

func (a *GoColorAnim) Draw(img *image.RGBA) {
	if a.patIdx == a.prevPatIdx {
		return
	}
	a.prevPatIdx = a.patIdx
	a.gc.Clear(colors.Black)
	idx := 0
	for row := range a.bounds.Dy() / 60 {
		y := float64(row * 60)
		for col := range a.bounds.Dx() / 60 {
			x := float64(col * 60)
			DrawColorSquare(a.gc, x, y, a.patIdx, colors.Map[a.colorList[idx]])
			idx = (idx + 1) % len(a.colorList)
		}
	}
    draw.Draw(img, a.bounds.ToInt(), a.gc.Image().(*image.RGBA),
        image.Point{}, draw.Over)
}

//-----------------------------------------------------------------------

func DrawColorSquare(gc *gg.Context, x, y float64, typ int, col colors.RGBA) {
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
        gc.SetFillColor(colors.Black)
        gc.DrawCircle(x+30.0, y+30.0, 30.0-1.0)
        gc.Fill()
    case 3:
        gc.SetFillColor(col)
        gc.DrawCircle(x+30.0, y+30.0, 30.0)
        gc.Fill()
        gc.SetFillColor(colors.Black)
        gc.DrawCircle(x+30.0, y+30.0, 30.0-15.0)
        gc.Fill()
    }
}

