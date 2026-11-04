package main

import (
	"math/rand/v2"
	"time"

	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
)

const (
	defNumPolygons = 20
)

type PolygonWin struct {
	windowEmbed
	canvas *Group
}

func NewPolygonWindow(bounds geom.Rectangle[int],
	numObjs int) *PolygonWin {
	numEdges := 3
	a := &PolygonWin{}

	a.bounds = bounds
	a.gc = gg.NewContext(bounds.Dx(), bounds.Dy())
	if numObjs <= 0 {
		numObjs = defNumPolygons
	}

	root := NewGroup()
	root.Layout = NewPadLayout(5)
	root.SetSize(bounds.ToFloat().Size())

	a.canvas = NewGroup()
	a.canvas.Layout = &NullLayout{}
	root.Add(a.canvas)

	for i := 0; i < numObjs; i++ {
		a.canvas.Add(NewPolygon(numEdges, a.canvas.Rect()))
	}
	a.Root = root

	return a
}

func (a *PolygonWin) Init() {
	a.gc.SetLineWidth(2.0)
	a.gc.SetLineCapRound()
	a.gc.SetLineJoinRound()
}

func (a *PolygonWin) Update(t time.Time) {
	a.Root.Update(t)
}

func (a *PolygonWin) Refresh() {
	a.gc.Clear(colors.DarkSlateGrey)
	a.Root.Draw(a.gc)
}

//-----------------------------------------------------------------------

type Polygon struct {
	nodeEmbed
	rect      Rectangle
	posList   []Point
	velList   []Point
	fillColor colors.RGBA
}

func NewPolygon(edges int, rect Rectangle) *Polygon {
	p := &Polygon{}
	p.Init(p)
	p.InitProp("Polygon")
	p.rect = rect
	p.posList = make([]Point, edges)
	p.velList = make([]Point, edges)
	for i := range edges {
		p.posList[i] = rect.RelPos(rand.Float64(), rand.Float64())
		p.velList[i] = Point{
			rand.Float64()*5.0 - 2.5,
			rand.Float64()*5.0 - 2.5,
		}
	}
	p.fillColor = colors.RandGroupColor(colors.Purples).Alpha(0.5)
	return p
}

func (p *Polygon) Update(t time.Time) {
	for i, pos := range p.posList {
		pos.Move(p.velList[i])
		if pos.X < p.rect.Min.X || pos.X >= p.rect.Max.X {
			p.velList[i].X *= -1
			pos.X += p.velList[i].X
		}
		if pos.Y < p.rect.Min.Y || pos.Y >= p.rect.Max.Y {
			p.velList[i].Y *= -1
			pos.Y += p.velList[i].Y
		}
		p.posList[i] = pos
	}
}

func (p *Polygon) Draw(gc *gg.Context) {
	gc.MoveTo(p.posList[0].X, p.posList[0].Y)
	for _, pos := range p.posList[1:] {
		gc.LineTo(pos.X, pos.Y)
	}
	gc.ClosePath()
	gc.SetLineColor(p.BorderColor())
	gc.SetLineWidth(p.BorderWidth())
	gc.SetFillColor(p.fillColor)
	gc.FillStroke()
}
