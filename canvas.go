package main

import (
	"log"
	"container/list"
	"github.com/stefan-muehlebach/gg"
)

//----------------------------------------------------------------------------

const (
	guiDebugging = false
)

//----------------------------------------------------------------------------

type Leaf interface {}

type LayoutManager interface {
	Layout(nodes *list.List, size Point)
	MinSize(nodes *list.List) (Point)
}

//----------------------------------------------------------------------------

type Node interface {
	Wrapper() Node
    Wrappee() *nodeEmbed
	ToBack()
	ToFront()
	Pos() Point
	SetPos(pos Point)
	MinSize() Point
	Size() Point
	SetSize(size Point)
	Hide()
	Show()
	IsVisible() bool
	Contains(pt Point) (Node)
	Draw(gc *gg.Context)
}

//----------------------------------------------------------------------------

type nodeEmbed struct {
	wrapper Node
	parent *containerEmbed
	pos, minSize, size Point
	rect Rectangle
	isHidden bool
	transl, rotate, scale, transf Matrix
}

func (c *nodeEmbed) Init(n Node) {
	c.wrapper = n
	c.transl = Identity()
	c.rotate = Identity()
	c.scale  = Identity()
	c.transf = Identity()
}

func (c *nodeEmbed) Wrapper() (Node) {
	return c.wrapper
}

func (c *nodeEmbed) Wrappee() (*nodeEmbed) {
	return c
}

func (c *nodeEmbed) ToBack() {
	var e *list.Element

	if c.parent == nil {
		log.Fatal("node: this child is not attached")
	}
	p := c.parent
	for e = p.childList.Front(); e != nil; e = e.Next() {
		if e.Value.(*nodeEmbed) == c {
			break
		}
	}
	if e == nil {
		return
	}
	p.childList.MoveToFront(e)
}

func (c *nodeEmbed) ToFront() {
	var e *list.Element

	if c.parent == nil {
		log.Fatal("node: this child is not attached")
	}
	p := c.parent
	for e = p.childList.Front(); e != nil; e = e.Next() {
		if e.Value.(*nodeEmbed) == c {
			break
		}
	}
	if e == nil {
		return
	}
	p.childList.MoveToBack(e)
}

func (c *nodeEmbed) Pos() Point {
	return c.pos
}

func (c *nodeEmbed) SetPos(pos Point) {
	c.pos = pos
	c.rect = Rectangle{Min: c.pos, Max: c.pos.Add(c.size)}
}

func (c *nodeEmbed) MinSize() Point {
	return c.minSize
}

func (c *nodeEmbed) Size() Point {
	return c.size.Max(c.minSize)
}

func (c *nodeEmbed) SetSize(size Point) {
	c.size = size
	c.rect = Rectangle{Min: c.pos, Max: c.pos.Add(c.size)}
}

func (c *nodeEmbed) Hide() {
	c.isHidden = true
}

func (c *nodeEmbed) Show() {
	c.isHidden = false
}

func (c *nodeEmbed) IsVisible() (bool) {
	return !c.isHidden
}

func (c *nodeEmbed) Contains(pt Point) (Node) {
	if pt.In(c.rect) {
		return c.wrapper
	} else {
		return nil
	}
}

//----------------------------------------------------------------------------

type Container interface {
	Add(nodes ...Node)
	Del(node Node)
	Purge()
	layout()
}

//----------------------------------------------------------------------------

type containerEmbed struct {
	nodeEmbed
	childList *list.List
	Layout LayoutManager
}

func (c *containerEmbed) Init(n Node) {
	c.nodeEmbed.Init(n)
	c.childList = list.New()
	c.Layout = &NullLayout{}
}

func (c *containerEmbed) Add(nodes ...Node) {
	for _, node := range nodes {
		embed := node.Wrappee()
		embed.parent = c
		c.childList.PushBack(node)
		c.layout()
	}
}

func (c *containerEmbed) Del(n Node) {
	for elem := c.childList.Front(); elem != nil; elem = elem.Next() {
		node := elem.Value.(Node)
		if n != node {
			continue
		}
		c.childList.Remove(elem)
		break
	}
	c.layout()
}

func (c *containerEmbed) Purge() {
	c.childList.Init()
	c.layout()
}

func (c *containerEmbed) layout() {
    if c.Layout == nil {
		return
	}
	c.Layout.Layout(c.childList, c.Wrapper().Size())
}

//----------------------------------------------------------------------------

// Der Typ Border wird fuer die Bezeichnung der vier Bildschirmseiten oder
// Richtungen verwendet.
type Border int

const (
    Left Border = iota
    Top
    Right
    Bottom
)

// Mit dem NullLayout werden die verwalteten Nodes per SetPos platziert und
// werden durch den Container nicht mehr weiter verwaltet. MinSize liefert
// die maximale Grösse aller verwalteten Nodes.
type NullLayout struct {}

func (l *NullLayout) Layout(childList *list.List, size Point) {

}

func (l *NullLayout) MinSize(childList *list.List) Point {
	minSize := Point{}
	for elem := childList.Front(); elem != nil; elem = elem.Next() {
		child := elem.Value.(*nodeEmbed).Wrapper()
		if !child.IsVisible() {
			continue
		}
		minSize = minSize.Max(child.Pos().Add(child.MinSize()))
	}
	return minSize
}

// Mit PaddedLayout kann sinnvollerweise nur ein Node verwaltet werden, der
// mit einem konfigurierbaren Abstand auf die ganze Grösse des Containers
// expandiert wird.
type PaddedLayout struct {
    pad [4]float64
}
// Mit den variablen Parametern pads können die Ränder definiert werden.
// Dabei gilt:
//   - : verwende das Property 'Padding'
//     a       : verwende a für alle Ränder
//     a,b     : verwende a für die horizontalen und b für die vertikalen
//     Ränder
//     a,b,c   : verwende a für links, b für oben und unten c für rechts
//     a,b,c,d : (dito) und d für den unteren Rand.
func NewPaddedLayout(pads ...float64) *PaddedLayout {
    var l, t, r, b float64
    switch len(pads) {
    case 0:
        l, t, r, b = 10, 10, 10, 10
    case 1:
        l, t, r, b = pads[0], pads[0], pads[0], pads[0]
    case 2:
        l, t, r, b = pads[0], pads[1], pads[0], pads[1]
    case 3:
        l, t, r, b = pads[0], pads[1], pads[2], pads[1]
    case 4:
        l, t, r, b = pads[0], pads[1], pads[2], pads[3]
    }
    return &PaddedLayout{[4]float64{l, t, r, b}}
}

func (l *PaddedLayout) Layout(childList *list.List, size Point) {
    pos := Point{l.pad[Left], l.pad[Top]}
    siz := Point{size.X - l.pad[Left] - l.pad[Right],
        size.Y - l.pad[Top] - l.pad[Bottom]}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*nodeEmbed).Wrapper()
        child.SetSize(siz)
        child.SetPos(pos)
    }
}

func (l *PaddedLayout) MinSize(childList *list.List) Point {
    minSize := Point{}
    for elem := childList.Front(); elem != nil; elem = elem.Next() {
        child := elem.Value.(*nodeEmbed).Wrapper()
        if !child.IsVisible() {
            continue
        }
        minSize = minSize.Max(child.MinSize())
    }
    return minSize.Add(Point{l.pad[Left] + l.pad[Right],
        l.pad[Top] + l.pad[Bottom]})
}

