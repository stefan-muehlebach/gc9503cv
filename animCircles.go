package main

import (
	"container/list"
	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	//"log"
	"math/rand/v2"
	"time"
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
	eventHandlerEmbed
	appletEmbed
	objList   *list.List
	activeObj *Circle
}

func NewCircleAnimation(bounds geom.Rectangle[int], numObjs int) *CircleAnim {
	//var circ *Circle

	a := &CircleAnim{}
	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	if numObjs <= 0 {
		numObjs = defNumCircles
	}
	rect := a.bounds.Sub(a.bounds.Min).ToFloat()
	a.Root = NewGroup()
	a.Root.SetPos(Point{})
	a.Root.SetSize(Point{rect.Dx(), rect.Dy()})
	grp := NewGroup()
	grp.SetPos(Point{})
	grp.SetSize(Point{rect.Dx(), rect.Dy()/2})
	a.Root.Add(grp)
	for range numObjs {
		grp.Add(NewCircle(rect))
	}
	return a
}

func (a *CircleAnim) Init() {
	a.gc.SetLineWidth(2.0)
	a.gc.SetLineCapRound()
	a.gc.SetLineJoinRound()
}

func (a *CircleAnim) Update(dt time.Duration) {
	a.Root.Update(dt)
}

func (a *CircleAnim) Refresh() {
	a.gc.Clear(colors.Black)
	a.Root.Draw(a.gc)
}

func (a *CircleAnim) OnInputEvent(ev MouseEvent) {
	if n := a.Root.Contains(ev.Pos); n != nil {
		n.OnInputEvent(ev)
	}
}

//-----------------------------------------------------------------------

type Circle struct {
	nodeEmbed
	rect                 Rectangle
	pos, vel             Point
	r, rMin, rMax, dr    float64
	lineColor, fillColor colors.RGBA
	isActive, isSelected, isBlocked bool
}

func NewCircle(rect Rectangle) *Circle {
	var dp Point

	c := &Circle{}
	c.Init(c)
	c.rect = rect
	c.rMin = 20.0
	c.rMax = 50.0
	c.r = c.rMin + (c.rMax-c.rMin)*rand.Float64()
	c.pos = rect.Inset(c.r, c.r).RelPos(rand.Float64(), rand.Float64())
	c.vel = Point{RandVel(2.0, 5.0), RandVel(2.0, 5.0)}
	c.lineColor = colors.White
	c.fillColor = colors.RandColorByGroup(colors.Blues).Alpha(0.8)

	c.SetOnPress(func(ev MouseEvent) {
		if ev.Button.IsSet(LeftButton) {
			c.isBlocked = true
			dp = ev.Pos.Sub(c.pos)
		}
	})

	c.SetOnDrag(func(ev MouseEvent) {
		c.SetPos(ev.Pos.Sub(dp))
	})

	c.SetOnRelease(func(ev MouseEvent) {
		if ev.Button.IsSet(LeftButton) {
			c.isBlocked = false
		}
	})

	c.SetOnClick(func(ev MouseEvent) {
		if ev.Button.IsSet(RightButton) {
			// Kreis loeschen
		}
	})
	
	return c
}

func (c *Circle) Update(dt time.Duration) {
	if c.isBlocked {
		return
	}
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

func (c *Circle) Contains(pt Point) Node {
	if c.pos.Distance(pt) <= c.r {
		return c
	} else {
		return nil
	}
}

func (c *Circle) Draw(gc *gg.Context) {
	gc.DrawCircle(c.pos.X, c.pos.Y, c.r)
	gc.ClosePath()
	if c.isActive {
		gc.SetLineWidth(7.0)
	} else {
		gc.SetLineWidth(2.0)
	}
	if c.isSelected {
		gc.SetLineColor(colors.GoFuchsia)
		gc.SetFillColor(colors.GoFuchsia.Alpha(0.8))
	} else {
		gc.SetLineColor(c.lineColor)
		gc.SetFillColor(c.fillColor)
	}
	gc.FillStroke()
}
