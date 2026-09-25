package main

import (
	"log"
	"math"
	"math/rand/v2"
	"time"

	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
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
	windowEmbed
	selection Node
	canvas    *Panel
}

func NewCircleAnimation(bounds geom.Rectangle[int], numObjs int) *CircleAnim {
	a := &CircleAnim{}
	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	if numObjs <= 0 {
		numObjs = defNumCircles
	}

	root := NewPanel(colors.SlateGray)
	root.Layout = NewPadLayout(20)
	root.SetSize(bounds.ToFloat().Size())
	a.Root = root

	a.canvas = NewPanel(colors.DarkSlateGray)
	a.canvas.Layout = &NullLayout{}
	a.Root.Add(a.canvas)

	log.Printf("a.canvas.Bounds(): %v", a.canvas.Bounds())
	log.Printf("a.canvas.Rect()  : %v", a.canvas.Rect())

	for range numObjs {
		a.canvas.Add(NewCircle(a.canvas.Rect(), colors.Greens))
	}

	a.Root.SetOnClick(func(ev InputEvent) {
		log.Printf("Click in the root object")
	})

	a.canvas.SetOnClick(func(ev InputEvent) {
		switch {
		case ev.Button.IsSet(LeftButton):
			c := NewCircle(a.canvas.Rect(), colors.Reds)
			pt := ev.Pos.Sub(a.canvas.Pos())
			c.SetPos(pt.SubXY(c.r, c.r))
			a.canvas.Add(c)
		default:
			log.Printf("No click with left button")
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
	a.Root.Update(dt)

	for el := a.canvas.childList.Front(); el != nil; el = el.Next() {
		this := el.Value.(*Circle)
		for oel := a.canvas.childList.Front(); oel != nil; oel = oel.Next() {
			other := oel.Value.(*Circle)
			if this.r <= other.r {
				continue
			}
			if this.Overlaps(other) {
				t := other.r / (this.r + other.r)
				this.fillColor = this.fillColor.Interpolate(other.fillColor, t)
				this.r += math.Cbrt(other.r)
				this.SetSize(Point{2 * this.r, 2 * this.r})
				a.canvas.childList.Remove(oel)
			}
		}
	}
}

func (a *CircleAnim) Refresh() {
	a.gc.Clear(colors.Black)
	a.Root.Draw(a.gc)
}

//-----------------------------------------------------------------------

const (
	radMin = 15.0
	radMax = 60.0
	velMin = 0.7
	velMax = 3.5
)

type Circle struct {
	nodeEmbed
	rect                             Rectangle
	vel                              Point
	r                                float64
	fillColor                        colors.RGBA
	isActive, isSelected, isDragging bool
}

func NewCircle(rect Rectangle, colGrp colors.ColorGroup) *Circle {
	var dp Point

	c := &Circle{}
	c.Init(c)
	c.InitProp("Circle")
	c.rect = rect
	t := rand.Float64()
	c.r = (1.0-t)*radMin + t*radMax
	c.size = Point{2 * c.r, 2 * c.r}
	r := rect.Inset(c.r, c.r)
	mp := r.RelPos(rand.Float64(), rand.Float64())
	c.pos = mp.SubXY(c.r, c.r)
	vel := (1.0-t)*velMax + t*velMin
	dir := 2.0 * math.Pi * rand.Float64()
	c.vel = Point{vel * math.Sin(dir), vel * math.Cos(dir)}
	c.fillColor = colors.RandColorByGroup(colGrp).Alpha(0.8)

	c.SetOnPress(func(ev InputEvent) {
		if ev.Button.IsSet(LeftButton) {
			c.isDragging = true
			dp = ev.Pos.Sub(c.Pos())
		}
	})

	c.SetOnDrag(func(ev InputEvent) {
		if c.isDragging {
			c.SetPos(ev.Pos.Sub(dp))
		}
	})

	c.SetOnRelease(func(ev InputEvent) {
		if ev.Button.IsSet(LeftButton) {
			c.isDragging = false
		}
	})

	c.SetOnEnter(func(ev InputEvent) {
		c.isActive = true
	})
	c.SetOnLeave(func(ev InputEvent) {
		c.isActive = false
	})

	c.SetOnClick(func(ev InputEvent) {
		if ev.Button.IsSet(LeftButton) {
			c.isSelected = !c.isSelected
		}
		if ev.Button.IsSet(RightButton) {
			c.parent.Del(c)
		}
	})

	return c
}

func (c *Circle) Update(dt time.Duration) {
	if c.isDragging {
		return
	}
	c.pos.Move(c.vel)
	b := c.Bounds()
	if b.Min.X < c.rect.Min.X || b.Max.X > c.rect.Max.X {
		c.vel.X *= -1.0
		c.pos.X += c.vel.X
	}
	if b.Min.Y < c.rect.Min.Y || b.Max.Y > c.rect.Max.Y {
		c.vel.Y *= -1.0
		c.pos.Y += c.vel.Y
	}
}

func (c *Circle) SetSize(size Point) {
	sz := min(size.X, size.Y)
	c.r = sz / 2.0
	ds := c.size.X - sz
	c.pos = c.pos.AddXY(ds/2.0, ds/2.0)
	c.size = Point{2 * c.r, 2 * c.r}
}

func (c *Circle) Bounds() Rectangle {
	return Rectangle{Min: c.pos, Max: c.pos.Add(c.size)}
}

func (c *Circle) FindTarget(pt Point) (Node, Point) {
	mp := c.pos.AddXY(c.r, c.r)
	if mp.Distance(pt) <= c.r {
		return c, pt
	} else {
		return nil, Point{}
	}
}

func (c *Circle) Draw(gc *gg.Context) {
	mp := c.pos.AddXY(c.r, c.r)
	gc.DrawCircle(mp.X, mp.Y, c.r)

	if c.isActive {
		gc.SetLineWidth(c.PushedBorderWidth())
	} else {
		gc.SetLineWidth(c.BorderWidth())
	}
	gc.SetLineColor(c.BorderColor())
	gc.SetFillColor(c.fillColor)
	gc.FillStroke()

	if c.isSelected {
		b := c.Bounds()
		l := c.r / 3.0
		//mp := c.pos.AddXY(c.r, c.r)

		gc.SetLineWidth(3.0)
		gc.SetLineColor(colors.Red)
		gc.MoveTo(b.Min.X, b.Min.Y+l)
		gc.LineTo(b.Min.X, b.Min.Y)
		gc.LineTo(b.Min.X+l, b.Min.Y)

		gc.MoveTo(b.Max.X-l, b.Min.Y)
		gc.LineTo(b.Max.X, b.Min.Y)
		gc.LineTo(b.Max.X, b.Min.Y+l)

		gc.MoveTo(b.Max.X, b.Max.Y-l)
		gc.LineTo(b.Max.X, b.Max.Y)
		gc.LineTo(b.Max.X-l, b.Max.Y)

		gc.MoveTo(b.Min.X+l, b.Max.Y)
		gc.LineTo(b.Min.X, b.Max.Y)
		gc.LineTo(b.Min.X, b.Max.Y-l)

		//gc.DrawLine(mp.X-l, mp.Y, mp.X+l, mp.Y)
		//gc.DrawLine(mp.X, mp.Y-l, mp.X, mp.Y+l)

		gc.Stroke()
	}
}

func (c *Circle) Overlaps(c2 *Circle) bool {
	mp1 := c.pos.AddXY(c.r, c.r)
	mp2 := c2.pos.AddXY(c2.r, c2.r)
	d := mp1.Distance(mp2)
	return c.r >= d+c2.r
}
