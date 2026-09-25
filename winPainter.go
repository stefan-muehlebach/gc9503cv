package main

import (
	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"

	//"log"
	"time"
)

var (
	bgColor = colors.Snow
)

type PainterWin struct {
	windowEmbed
}

func NewPainterWindow(bounds geom.Rectangle[int]) *PainterWin {
	a := &PainterWin{}

	a.bounds = bounds
	a.gc = gg.NewContext(bounds.Dx(), bounds.Dy())

	root := NewGroup()
	root.Layout = NewMaxLayout()
	root.SetSize(bounds.ToFloat().Size())
	a.Root = root

	//grp := NewGroup()
	//grp.Layout = NewColumnGridLayout(1)
	//a.Root.Add(grp)

	//log.Printf("grp.Size(): %v", grp.Size())
	//w := (grp.Size().X - 20.0)/2.0
	//h := (grp.Size().Y - 20.0)/2.0

	var p1 Point
	canvas := NewPanel(bgColor)
	canvas.SetOnClick(func(ev InputEvent) {
		switch {
		case ev.Button.IsSet(LeftButton):
			p2 := ev.Pos
			if p1.X != 0 || p1.Y != 0 {
				l := NewLine(p1, p2, linColor)
				canvas.Add(l)
			}
			d := NewDot(p2, dotColor)
			canvas.Add(d)
			p1 = p2
		}
	})
	root.Add(canvas)

	/*
		dotColors := []colors.RGBA{
			colors.SkyBlue,
			colors.WhiteSmoke,
			colors.PaleGreen,
			colors.DeepPink,
		}
		for i, dotColor := range dotColors {
			canvas[i] = NewPanel(dotColor.Dark(0.7))
			canvas[i].IsClipping = true
			canvas[i].SetOnClick(func(ev InputEvent) {
				switch {
				case ev.Button.IsSet(LeftButton):
					d := NewDot(ev.Pos, dotColor)
					d.isActive = true
					canvas[i].Add(d)
				}
			})
			//log.Printf("(before) canvas.Size(): %v", canvas.Size())
			//log.Printf("         canvas.Pos() : %v", canvas.Pos())
			grp.Add(canvas[i])
			//log.Printf("(after)  canvas.Size(): %v", canvas.Size())
			//log.Printf("         canvas.Pos() : %v", canvas.Pos())
			//canvas.Add(NewDot(Point{canvas.Size().X, 0}, colors.Red))
		}
	*/

	//for i, dotColor := range dotColors {
	/*
		for range 10 {
			p1 := canvas.Rect().Inset(dotRadius, dotRadius).RelPos(rand.Float64(), rand.Float64())
			p2 := canvas.Rect().Inset(dotRadius, dotRadius).RelPos(rand.Float64(), rand.Float64())
			d1 := NewDot(p1, dotColor)
			d2 := NewDot(p2, dotColor)
			l := NewLine(p1, p2, linColor)
			l.isActive = true
			canvas.Add(l, d1, d2)
		}
	*/
	//}

	return a
}

func (a *PainterWin) Init() {
	a.gc.SetLineWidth(2.0)
	a.gc.SetLineCapRound()
	a.gc.SetLineJoinRound()
}

func (a *PainterWin) Update(dt time.Duration) {
	a.Root.Update(dt)
}

func (a *PainterWin) Refresh() {
	a.gc.Clear(colors.Black)
	a.Root.Draw(a.gc)
}

//----------------------------------------------------------------------------

var (
	dotRadius  = 3.0
	haloRadius = 4.0 * dotRadius
	dotColor   = colors.DimGray
	linColor   = colors.Gray
	haloColor  = linColor.Alpha(0.15)
)

type Dot struct {
	nodeEmbed
	lineColor colors.RGBA
	isActive  bool
}

func NewDot(pos Point, color colors.RGBA) *Dot {
	d := &Dot{}
	d.Init(d)
	d.SetPos(pos.SubXY(haloRadius, haloRadius))
	d.SetMinSize(Point{haloRadius, haloRadius}.Mul(2))
	d.lineColor = color
	d.SetOnClick(func(ev InputEvent) {
		switch ev.Button {
		case RightButton:
			d.parent.Del(d)
		}
	})
	d.SetOnEnter(func(ev InputEvent) {
		d.isActive = true
	})
	d.SetOnLeave(func(ev InputEvent) {
		d.isActive = false
	})
	return d
}

func (d *Dot) Update(dt time.Duration) {}

func (d *Dot) Draw(gc *gg.Context) {
	mp := d.Bounds().Center()
	if d.isActive {
		gc.SetFillColor(haloColor)
		gc.DrawPoint(mp.X, mp.Y, haloRadius)
		gc.Fill()
	}
	gc.SetLineWidth(1.0)
	gc.SetLineColor(d.lineColor)
	gc.SetFillColor(bgColor)
	gc.DrawPoint(mp.X, mp.Y, dotRadius)
	gc.FillStroke()
}

//----------------------------------------------------------------------------

type Line struct {
	nodeEmbed
	lineColor colors.RGBA
	isActive  bool
}

func NewLine(p1, p2 Point, color colors.RGBA) *Line {
	l := &Line{}
	l.Init(l)
	l.SetPos(p1)
	l.SetSize(p2.Sub(p1))
	l.lineColor = color
	return l
}

func (l *Line) Update(dt time.Duration) {}

func (l *Line) Draw(gc *gg.Context) {
	if l.isActive {
		gc.SetLineWidth(12.0)
		gc.SetLineColor(haloColor)
		gc.DrawLine(l.pos.X, l.pos.Y, l.pos.X+l.size.X, l.pos.Y+l.size.Y)
		gc.Stroke()
	}
	gc.SetLineWidth(2.0)
	gc.SetLineColor(l.lineColor)
	gc.DrawLine(l.pos.X, l.pos.Y, l.pos.X+l.size.X, l.pos.Y+l.size.Y)
	gc.Stroke()
}
