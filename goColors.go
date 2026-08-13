package main

import (
	"time"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"gc9503cv/gc9503cv/geom"
)

type GoColorAnim struct {
	callbackEmbed
	rect geom.Rectangle[int]
	t0 time.Time
	patIdx, prevPatIdx int
	colorList []string
}

func NewGoColorAnimation(rect geom.Rectangle[int]) *GoColorAnim {
	a := &GoColorAnim{}
	a.rect = rect
	a.patIdx = -1
	a.prevPatIdx = -1
	a.colorList = colors.Groups[colors.GoColors]
	return a
}

func (a *GoColorAnim) Init(gc *gg.Context) {
	gc.SetLineWidth(2.0)
	gc.SetLineCapRound()
	gc.SetLineJoinRound()
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

func (a *GoColorAnim) Draw(gc *gg.Context) {
	if a.patIdx == a.prevPatIdx {
		return
	}
	a.prevPatIdx = a.patIdx
	gc.SetFillColor(colors.Black)
	gc.Clear()
	idx := 0
	for row := range a.rect.Dy() / 60 {
		y := float64(row * 60)
		for col := range a.rect.Dx() / 60 {
			x := float64(col * 60)
			DrawColorSquare(gc, x, y, a.patIdx, colors.Map[a.colorList[idx]])
			idx = (idx + 1) % len(a.colorList)
		}
	}
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

