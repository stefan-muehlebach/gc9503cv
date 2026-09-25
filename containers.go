package main

import (
	"container/list"
	"time"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
)

//-----------------------------------------------------------------------------

type Container interface {
	Node
	Add(nl ...Node)
	Del(n Node)
	FindTarget(pt Point) (Node, Point)
	Update(dt time.Duration)
	Draw(gc *gg.Context)
	//SetLayoutManager(layoutManager LayoutManager)
	//LayoutManager() LayoutManager
	layout()
}

//-----------------------------------------------------------------------------

type containerEmbed struct {
	nodeEmbed
	childList *list.List
	Layout    LayoutManager
}

func (c *containerEmbed) Init(n Node) {
	c.nodeEmbed.Init(n)
	c.childList = list.New()
	c.Layout = &NullLayout{}
}

func (c *containerEmbed) SetSize(size Point) {
	c.nodeEmbed.SetSize(size)
	c.layout()
}

func (c *containerEmbed) MinSize() Point {
	if c.minSize.Eq(Point{0, 0}) {
		return c.Layout.MinSize(c.childList)
	} else {
		return c.nodeEmbed.MinSize()
	}
}

func (c *containerEmbed) Add(nl ...Node) {
	for _, n := range nl {
		embed := n.Wrappee()
		embed.parent = c
		c.childList.PushBack(n)
		c.layout()
	}
}

func (c *containerEmbed) Del(n Node) {
	for e := c.childList.Front(); e != nil; e = e.Next() {
		if e.Value.(Node) == n {
			c.childList.Remove(e)
			break
		}
	}
	c.layout()
}

func (c *containerEmbed) FindTarget(pt Point) (Node, Point) {
	if n, _ := c.nodeEmbed.FindTarget(pt); n == nil {
		return nil, Point{}
	}
	pt = pt.Sub(c.pos)
	for e := c.childList.Back(); e != nil; e = e.Prev() {
		if n, ptRel := e.Value.(Node).FindTarget(pt); n != nil {
			return n, ptRel
		}
	}
	return c, pt
}

func (c *containerEmbed) Update(dt time.Duration) {
	for e := c.childList.Front(); e != nil; e = e.Next() {
		e.Value.(Node).Update(dt)
	}
}

func (c *containerEmbed) Draw(gc *gg.Context) {
	gc.Push()
	gc.Translate(c.pos.AsCoord())
	for e := c.childList.Front(); e != nil; e = e.Next() {
		e.Value.(Node).Draw(gc)
	}
	gc.Pop()
}

func (c *containerEmbed) layout() {
	if c.Layout == nil {
		return
	}
	c.Layout.Layout(c.childList, c.wrapper.Size())
}

//-----------------------------------------------------------------------------

type Group struct {
	containerEmbed
}

func NewGroup() *Group {
	g := &Group{}
	g.Init(g)
	return g
}

//-----------------------------------------------------------------------------

type Panel struct {
	containerEmbed
	fillColor, lineColor colors.RGBA
	IsClipping bool
}

func NewPanel(color colors.RGBA) *Panel {
	p := &Panel{}
	p.Init(p)
	p.InitProp("Panel")
	p.fillColor = color
	p.lineColor = color
	return p
}

func (p *Panel) Draw(gc *gg.Context) {
	gc.SetLineWidth(1)
	gc.SetLineColor(p.lineColor)
	gc.SetFillColor(p.fillColor)
	gc.DrawRectangle(p.Bounds().Inset(0.5,0.5).AsCoord())
	gc.FillStroke()

	if p.IsClipping {
		gc.DrawRectangle(p.Bounds().AsCoord())
		gc.Clip()
		p.containerEmbed.Draw(gc)
		gc.ResetClip()
	} else {
		p.containerEmbed.Draw(gc)
	}
}
