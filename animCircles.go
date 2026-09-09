package main

import (
	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	// "log"
	"math"
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
	windowEmbed
	selection Node
}

func NewCircleAnimation(b geom.Rectangle[int], numObjs int) *CircleAnim {
	a := &CircleAnim{}
	a.bounds = b
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	if numObjs <= 0 {
		numObjs = defNumCircles
	}
	bounds := a.bounds.ToFloat()
	rect := bounds.Sub(bounds.Min)

	a.Root = NewGroup()
	a.Root.Move(Point{})
	a.Root.Resize(bounds.Size())

	for range numObjs {
		a.Root.Add(NewCircle(rect, colors.Greens))
	}

	a.Root.SetOnClick(func(ev MouseEvent) {
		c := NewCircle(rect, colors.Yellows)
		c.Move(ev.Pos)
		a.Root.Add(c)
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

/*
	cont := a.Root.(*Group)
	cont.Update(dt)
	for el := cont.childList.Front(); el != nil; el = el.Next() {
		this := el.Value.(*Circle)
		for oel := cont.childList.Front(); oel != nil; oel = oel.Next() {
			other := oel.Value.(*Circle)
			if this.r <= other.r {
				continue
			}
			if this.Overlaps(other) {
				t := other.r / (this.r + other.r)
				this.fillColor = this.fillColor.Interpolate(other.fillColor, t)
				this.r += math.Cbrt(other.r)
				cont.childList.Remove(oel)
			}
		}
	}
*/
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

const (
	radMin = 15.0
	radMax = 60.0
	velMin =  0.7
	velMax =  3.5
)

type Circle struct {
	nodeEmbed
	rect                 Rectangle
	vel             Point
	r    float64
	lineColor, fillColor colors.RGBA
	isActive, isSelected, isDragging bool
}

func NewCircle(rect Rectangle, colGrp colors.ColorGroup) *Circle {
	var dp Point

	c := &Circle{}
	c.Init(c)
	c.rect = rect
	t := rand.Float64()
	c.r = (1.0-t)*radMin + t*radMax
	c.size = Point{2*c.r, 2*c.r}
	c.pos = rect.Inset(c.r, c.r).RelPos(rand.Float64(), rand.Float64())
	vel := (1.0-t)*velMax + t*velMin
	dir := 2.0*math.Pi*rand.Float64()
	c.vel = Point{vel*math.Sin(dir), vel*math.Cos(dir)}
	c.lineColor = colors.White
	c.fillColor = colors.RandColorByGroup(colGrp).Alpha(0.8)

	c.SetOnPress(func(ev MouseEvent) {
		if ev.Button.IsSet(LeftButton) {
			c.isDragging = true
			dp = ev.Pos.Sub(c.Pos())
		}
	})

	c.SetOnDrag(func(ev MouseEvent) {
		c.Move(ev.Pos.Sub(dp))
	})

	c.SetOnRelease(func(ev MouseEvent) {
		if ev.Button.IsSet(LeftButton) {
			c.isDragging = false
		}
	})

	c.SetOnClick(func(ev MouseEvent) {
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
	if c.pos.X-c.r < c.rect.Min.X || c.pos.X+c.r > c.rect.Max.X {
		c.vel.X *= -1.0
		c.pos.X += c.vel.X
	}
	if c.pos.Y-c.r < c.rect.Min.Y || c.pos.Y+c.r > c.rect.Max.Y {
		c.vel.Y *= -1.0
		c.pos.Y += c.vel.Y
	}
}

func (c *Circle) Bounds() Rectangle {
	return Rectangle{Min: c.pos.SubXY(c.r, c.r), Max: c.pos.AddXY(c.r, c.r)}
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
/*
	gc.ClosePath()
	if c.isActive || c.isSelected {
		gc.SetLineWidth(7.0)
	} else {
		gc.SetLineWidth(2.0)
	}
	if c.isSelected {
		gc.SetLineColor(colors.GoYellow)
	} else {
		gc.SetLineColor(c.lineColor)
	}
*/
	gc.SetLineWidth(2.0)
	gc.SetLineColor(c.lineColor)
	gc.SetFillColor(c.fillColor)
	gc.FillStroke()

	if c.isSelected {
		b := c.Bounds()
		l := c.r/3.0
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

		//gc.DrawLine(c.pos.X-l, c.pos.Y, c.pos.X+l, c.pos.Y)
		//gc.DrawLine(c.pos.X, c.pos.Y-l, c.pos.X, c.pos.Y+l)

		gc.Stroke()
	}
}

func (c *Circle) Overlaps(c2 *Circle) bool {
	d := c.pos.Distance(c2.pos)
	return c.r >= d + c2.r
}

