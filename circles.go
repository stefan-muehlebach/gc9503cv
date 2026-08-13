package main

import (
	"log"
	"math/rand/v2"
	"time"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	//"github.com/stefan-muehlebach/gg/geom"
)

const (
	defNumCircles = 30
)

func RandVel(low, high float64) float64 {
	v := 2.0*rand.Float64() - 1.0
	if v >= 0.0 {
		return low + v*(high-low)
	} else {
		return -low + v*(high-low)
	}
}

type CircleAnim struct {
	callbackEmbed
	circList []*Circle
}

func NewCircleAnimation(numObjs int, rect Rectangle) *CircleAnim {
	a := &CircleAnim{}

	if numObjs <= 0 {
		numObjs = defNumCircles
	}
	a.circList = make([]*Circle, numObjs)
	for i := range numObjs {
		a.circList[i] = NewCircle(rect)
	}
	a.SetOnClick(func(ev MouseEvent) {
		log.Printf("Click!")
		circ := NewCircle(rect)
		circ.pos = ev.Pos.ToFloat()
		circ.r = 10.0
		a.circList = append(a.circList , circ)
	})
	return a
}

func (a *CircleAnim) Init(gc *gg.Context) {
	gc.SetLineWidth(2.0)
	gc.SetLineCapRound()
	gc.SetLineJoinRound()
}

func (a *CircleAnim) Update(dt time.Duration) {
	for _, c := range a.circList {
		c.Update(dt)
	}
}

func (a *CircleAnim) Draw(gc *gg.Context) {
	for _, c := range a.circList {
		c.Draw(gc)
	}
}

//-----------------------------------------------------------------------

type Circle struct {
	rect                   Rectangle
	pos, vel               Point
	r, rMin, rMax, dr      float64
	strokeColor, fillColor colors.RGBA
}

func NewCircle(rect Rectangle) *Circle {
	c := &Circle{}
	c.rect = rect
	c.rMin = 10.0
	c.rMax = 50.0
	c.dr = 0.2 + 0.4*rand.Float64()
	c.r = c.rMin + (c.rMax-c.rMin)*rand.Float64()

	c.pos = rect.Inset(c.r, c.r).RelPos(rand.Float64(), rand.Float64())
	c.vel = Point{RandVel(1.0, 3.0), RandVel(1.0, 3.0)}
	c.strokeColor = colors.White
	c.fillColor = colors.RandColorByGroup(colors.Blues).Alpha(0.5)
	return c
}

func (c *Circle) Update(dt time.Duration) {
	c.pos.Move(c.vel)
	if c.pos.X-c.r < c.rect.Min.X || c.pos.X+c.r > c.rect.Max.X {
		c.vel.X *= -1.0
		c.pos.X += c.vel.X
	}
	if c.pos.Y-c.r < c.rect.Min.Y || c.pos.Y+c.r > c.rect.Max.Y {
		c.vel.Y *= -1.0
		c.pos.Y += c.vel.Y
	}
	c.r += c.dr
	if c.r > c.rMax || c.r < c.rMin {
		c.dr = -c.dr
		c.r += c.dr
	}
}

func (c *Circle) Draw(gc *gg.Context) {
	gc.DrawCircle(c.pos.X, c.pos.Y, c.r)
	gc.ClosePath()
	gc.SetLineColor(c.strokeColor)
	gc.SetFillColor(c.fillColor)
	gc.FillStroke()
}
