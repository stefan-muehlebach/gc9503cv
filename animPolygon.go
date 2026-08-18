package main

import (
	"gc9503cv/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"math/rand/v2"
	"time"
)

const (
	defNumPolygons = 20
)

type PolygonAnim struct {
	callbackEmbed
	appletEmbed
	polyList []*Polygon
}

func NewPolygonAnimation(bounds geom.Rectangle[int],
	numObjs int) *PolygonAnim {
	numEdges := 3
	a := &PolygonAnim{}

	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	if numObjs <= 0 {
		numObjs = defNumPolygons
	}
	rect := a.bounds.Sub(a.bounds.Min).ToFloat()
	a.polyList = make([]*Polygon, numObjs)
	for i := 0; i < numObjs; i++ {
		a.polyList[i] = NewPolygon(numEdges, rect)
	}
	return a
}

func (a *PolygonAnim) Init() {
	a.gc.SetLineWidth(2.0)
	a.gc.SetLineCapRound()
	a.gc.SetLineJoinRound()
}

func (a *PolygonAnim) Update(dt time.Duration) {
	for _, p := range a.polyList {
		p.Update(dt)
	}
}

func (a *PolygonAnim) Refresh() {
	a.gc.Clear(colors.Black)
	for _, p := range a.polyList {
		p.Draw(a.gc)
	}
	//draw.Draw(img, a.bounds.ToInt(), a.gc.Image().(*image.RGBA),
	//	image.Point{}, draw.Over)
}

//-----------------------------------------------------------------------

type Polygon struct {
	rect                   Rectangle
	pos, vel               []Point
	strokeColor, fillColor colors.RGBA
}

func NewPolygon(edges int, rect Rectangle) *Polygon {
	p := &Polygon{}
	p.rect = rect
	p.pos = make([]Point, edges)
	p.vel = make([]Point, edges)
	for i := range edges {
		p.pos[i] = rect.RelPos(rand.Float64(), rand.Float64())
		p.vel[i] = Point{
			rand.Float64()*5.0 - 2.0,
			rand.Float64()*5.0 - 2.0,
		}
	}
	p.strokeColor = colors.White
	p.fillColor = colors.RandGroupColor(colors.Purples).Alpha(0.5)
	return p
}

func (p *Polygon) Update(dt time.Duration) {
	for i, pos := range p.pos {
		pos.Move(p.vel[i])
		if pos.X < p.rect.Min.X || pos.X > p.rect.Max.X {
			p.vel[i].X *= -1
			pos.X += p.vel[i].X
		}
		if pos.Y < p.rect.Min.Y || pos.Y > p.rect.Max.Y {
			p.vel[i].Y *= -1
			pos.Y += p.vel[i].Y
		}
		p.pos[i] = pos
	}
}

func (p *Polygon) Draw(gc *gg.Context) {
	gc.MoveTo(p.pos[0].X, p.pos[0].Y)
	for _, pos := range p.pos[1:] {
		gc.LineTo(pos.X, pos.Y)
	}
	gc.ClosePath()
	gc.SetLineColor(p.strokeColor)
	gc.SetFillColor(p.fillColor)
	gc.FillStroke()
}
