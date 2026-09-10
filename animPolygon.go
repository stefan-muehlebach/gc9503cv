package main

import (
	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"math/rand/v2"
	"time"
)

const (
	defNumPolygons = 20
)

type PolygonAnim struct {
	windowEmbed
	canvas *Panel
}

func NewPolygonAnimation(bounds geom.Rectangle[int],
	numObjs int) *PolygonAnim {
	numEdges := 3
	a := &PolygonAnim{}

	a.bounds = bounds
	a.gc = gg.NewContext(bounds.Dx(), bounds.Dy())
	if numObjs <= 0 {
		numObjs = defNumPolygons
	}

	a.Root = NewPanel(colors.SlateGray)
	a.Root.SetLayoutManager(NewPadLayout(20))
	a.Root.SetSize(bounds.ToFloat().Size())

	a.canvas = NewPanel(colors.DarkSlateGray)
	a.canvas.SetLayoutManager(&NullLayout{})
	a.Root.Add(a.canvas)

	for i := 0; i < numObjs; i++ {
		a.canvas.Add(NewPolygon(numEdges, a.canvas.Bounds()))
	}

	return a
}

func (a *PolygonAnim) Init() {
	a.gc.SetLineWidth(2.0)
	a.gc.SetLineCapRound()
	a.gc.SetLineJoinRound()
}

func (a *PolygonAnim) Update(dt time.Duration) {
	a.Root.Update(dt)
}

func (a *PolygonAnim) Refresh() {
	a.gc.Clear(colors.Black)
	a.Root.Draw(a.gc)
}

//-----------------------------------------------------------------------

type Polygon struct {
	nodeEmbed
	rect                   Rectangle
	posList                []Point
	velList                []Point
	strokeColor, fillColor colors.RGBA
}

func NewPolygon(edges int, rect Rectangle) *Polygon {
	p := &Polygon{}
	p.rect = rect
	p.posList = make([]Point, edges)
	p.velList = make([]Point, edges)
	for i := range edges {
		p.posList[i] = rect.RelPos(rand.Float64(), rand.Float64())
		p.velList[i] = Point{
			rand.Float64()*5.0 - 2.0,
			rand.Float64()*5.0 - 2.0,
		}
	}
	p.strokeColor = colors.White
	p.fillColor = colors.RandGroupColor(colors.Purples).Alpha(0.5)
	return p
}

func (p *Polygon) Update(dt time.Duration) {
	for i, pos := range p.posList {
		pos.Move(p.velList[i])
		if pos.X < p.rect.Min.X || pos.X > p.rect.Max.X {
			p.velList[i].X *= -1
			pos.X += p.velList[i].X
		}
		if pos.Y < p.rect.Min.Y || pos.Y > p.rect.Max.Y {
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
	gc.SetLineColor(p.strokeColor)
	gc.SetFillColor(p.fillColor)
	gc.FillStroke()
}
