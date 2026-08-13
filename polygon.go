package main

import (
	"math/rand/v2"
	"time"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
)

const (
	defNumPolygons = 20
)

type PolygonAnim struct {
	callbackEmbed
	polyList []*Polygon
}

func NewPolygonAnimation(numObjs int, rect Rectangle) *PolygonAnim {
	numEdges := 3
	a := &PolygonAnim{}

	if numObjs <= 0 {
		numObjs = defNumPolygons
	}
	a.polyList = make([]*Polygon, numObjs)
	for i := 0; i < numObjs; i++ {
		a.polyList[i] = NewPolygon(numEdges, rect)
	}
	return a
}

func (a *PolygonAnim) Init(gc *gg.Context) {
	gc.SetLineWidth(2.0)
	gc.SetLineCapRound()
	gc.SetLineJoinRound()
}

func (a *PolygonAnim) Update(dt time.Duration) {
	for _, p := range a.polyList {
		p.Update(dt)
	}
}

func (a *PolygonAnim) Draw(gc *gg.Context) {
	for _, p := range a.polyList {
		p.Draw(gc)
	}
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
