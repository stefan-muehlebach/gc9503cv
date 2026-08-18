package main

import (
	"container/list"
	"gc9503cv/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"log"
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
	callbackEmbed
	appletEmbed
	objList   *list.List
	activeObj *Circle
}

func NewCircleAnimation(bounds geom.Rectangle[int], numObjs int) *CircleAnim {
	a := &CircleAnim{}

	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	if numObjs <= 0 {
		numObjs = defNumCircles
	}
	rect := a.bounds.Sub(a.bounds.Min).ToFloat()
	a.objList = list.New()
	for range numObjs {
		a.objList.PushBack(NewCircle(rect))
	}
	a.SetOnMove(func(ev MouseEvent) {
		for e := a.objList.Back(); e != nil; e = e.Prev() {
			c := e.Value.(*Circle)
			if c.Contains(ev.Pos.ToFloat()) {
				if a.activeObj != nil {
					if a.activeObj == c {
						return
					}
					a.activeObj.isActive = false
					a.activeObj = nil
				}
				a.activeObj = c
				a.activeObj.isActive = true
				return
			}
		}
		if a.activeObj != nil {
			a.activeObj.isActive = false
			a.activeObj = nil
		}
	})
	a.SetOnDrag(func(ev MouseEvent) {
		log.Printf("'Drag' received")

	})
	a.SetOnPress(func(ev MouseEvent) {
		log.Printf("'Press' received")
		
	})
	a.SetOnLongPress(func(ev MouseEvent) {
		log.Printf("'LongPress' received")
		if a.activeObj == nil {
			return
		}
		for e := a.objList.Front(); e != nil; e = e.Next() {
			if e.Value.(*Circle) == a.activeObj {
				a.objList.MoveToFront(e)
				break
			}
		}
	})
	a.SetOnClick(func(ev MouseEvent) {
		log.Printf("'Click' received")
		//log.Printf("%v", ev)
		switch ev.Button {
		case LeftButton:
			if a.activeObj == nil {
				circ := NewCircle(rect)
				circ.pos = ev.Pos.ToFloat()
				circ.r = 40.0
				a.objList.PushBack(circ)
				a.activeObj = circ
				a.activeObj.isActive = true
			} else {
				if !a.activeObj.isSelected {
					a.activeObj.isSelected = true
				} else {
					a.activeObj.isSelected = false
				}
			}

		case RightButton:
			if a.activeObj == nil {
				return
			}
			for e := a.objList.Front(); e != nil; e = e.Next() {
				if e.Value.(*Circle) == a.activeObj {
					a.objList.Remove(e)
					a.activeObj = nil
					break
				}
			}
		}
	})
	return a
}

func (a *CircleAnim) Init() {
	a.gc.SetLineWidth(2.0)
	a.gc.SetLineCapRound()
	a.gc.SetLineJoinRound()
}

func (a *CircleAnim) Update(dt time.Duration) {
	for e := a.objList.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Circle)
		c.Update(dt)
	}
}

func (a *CircleAnim) Refresh() {
	a.gc.Clear(colors.Black)
	for e := a.objList.Front(); e != nil; e = e.Next() {
		c := e.Value.(*Circle)
		c.Draw(a.gc)
	}
	
	//draw.Draw(img, a.bounds.ToInt(), a.gc.Image().(*image.RGBA),
	//	image.Point{}, draw.Over)
}

//-----------------------------------------------------------------------

type Circle struct {
	rect                 Rectangle
	pos, vel             Point
	r, rMin, rMax, dr    float64
	lineColor, fillColor colors.RGBA
	isActive, isSelected bool
}

func NewCircle(rect Rectangle) *Circle {
	c := &Circle{}
	c.rect = rect
	c.rMin = 30.0
	c.rMax = 50.0
	c.r = c.rMin + (c.rMax-c.rMin)*rand.Float64()
	//c.dr = 0.2 + 0.4*rand.Float64()

	c.pos = rect.Inset(c.r, c.r).RelPos(rand.Float64(), rand.Float64())
	//c.vel = Point{RandVel(1.0, 3.0), RandVel(1.0, 3.0)}
	c.lineColor = colors.White
	c.fillColor = colors.RandColorByGroup(colors.Blues).Alpha(0.8)
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

func (c *Circle) Contains(pt Point) bool {
	return c.pos.Distance(pt) <= c.r
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
