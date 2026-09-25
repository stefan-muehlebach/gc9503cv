package main

import (
	"container/list"
	//"log"
	"math"

	"github.com/stefan-muehlebach/gc9503cv/props"
	"github.com/stefan-muehlebach/gg"
)

//-----------------------------------------------------------------------------

type LayoutManager interface {
	Layout(childList *list.List, size Point)
	MinSize(childList *list.List) Point
}

//-----------------------------------------------------------------------------

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

//-----------------------------------------------------------------------------

type NullLayout struct {
	sz Point
}

func (l *NullLayout) Layout(childList *list.List, size Point) {
	//log.Printf("NullLayout.Layout(size: %v)", size)
	l.sz = size
}

func (l *NullLayout) MinSize(childList *list.List) Point {
	minSize := Point{}
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		if !n.IsVisible() {
			continue
		}
		minSize = minSize.Max(n.Bounds().Max)
	}
	return minSize.Min(l.sz)
}

//-----------------------------------------------------------------------------

type PadLayout struct {
	props.PropertyEmbed
	pads [4]float64
}

func NewPadLayout(pads ...float64) *PadLayout {
	var l, t, r, b float64

	lay := &PadLayout{}
	lay.InitProp("Layout")

	switch len(pads) {
	case 0:
		s := lay.Padding()
		l, t, r, b = s, s, s, s
	case 1:
		l, t, r, b = pads[0], pads[0], pads[0], pads[0]
	case 2:
		l, t, r, b = pads[0], pads[1], pads[0], pads[1]
	case 3:
		l, t, r, b = pads[0], pads[1], pads[2], pads[1]
	case 4:
		l, t, r, b = pads[0], pads[1], pads[2], pads[3]
	}
	lay.pads = [4]float64{l, t, r, b}

	return lay
}

func (l *PadLayout) Layout(childList *list.List, size Point) {
	pos := Point{l.pads[0], l.pads[1]}
	siz := Point{size.X - l.pads[0] - l.pads[2],
		size.Y - l.pads[1] - l.pads[3]}
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		n.SetSize(siz)
		n.SetPos(pos)
	}
}

func (l *PadLayout) MinSize(childList *list.List) Point {
	minSize := Point{}
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		if !n.IsVisible() {
			continue
		}
		minSize = minSize.Max(n.MinSize())
	}
	return minSize.Add(Point{l.pads[0] + l.pads[2], l.pads[1] + l.pads[3]})
}

//-----------------------------------------------------------------------------

type BoxLayout struct {
	props.PropertyEmbed
	orient Orientation
}

func NewHBoxLayout() *BoxLayout {
	l := &BoxLayout{orient: Horizontal}
	l.InitProp("Layout")
	return l
}

func NewVBoxLayout() *BoxLayout {
	l := &BoxLayout{orient: Vertical}
	l.InitProp("Layout")
	return l
}

func (l *BoxLayout) isSpacer(obj Node) bool {
	spc, ok := obj.(*Spacer)
	if !ok {
		return false
	}
	if l.orient == Horizontal {
		return spc.ExpandHorizontal()
	}
	return spc.ExpandVertical()
}

func (l *BoxLayout) Layout(childList *list.List, size Point) {
	spacers := 0
	total := 0.0
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		if !n.IsVisible() {
			continue
		}
		if l.isSpacer(n) {
			spacers++
			continue
		}
		nodeSize := n.MinSize()
		switch l.orient {
		case Horizontal:
			total += nodeSize.X
		case Vertical:
			total += nodeSize.Y
		}
	}
	extra, extraCell := 0.0, 0.0
	switch l.orient {
	case Horizontal:
		extra = size.X - total - l.InnerPadding()*float64(childList.Len()-
			spacers-1)
	case Vertical:
		extra = size.Y - total - l.InnerPadding()*float64(childList.Len()-
			spacers-1)
	}
	if spacers > 0 {
		extraCell = extra / float64(spacers)
	}
	pos := Point{}
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		if !n.IsVisible() {
			continue
		}
		if l.isSpacer(n) {
			switch l.orient {
			case Horizontal:
				pos.X += extraCell
			case Vertical:
				pos.Y += extraCell
			}
			continue
		}
		n.SetPos(pos)
		switch l.orient {
		case Horizontal:
			width := n.MinSize().X
			pos.X += width + l.InnerPadding()
			n.SetSize(Point{width, size.Y})
		case Vertical:
			height := n.MinSize().Y
			pos.Y += height + l.InnerPadding()
			n.SetSize(Point{size.X, height})
		}
	}
}

func (l *BoxLayout) MinSize(childList *list.List) Point {
	//log.Printf("BoxLayout.MinSize()")
	minSize := Point{}
	addPadding := false
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		if !n.IsVisible() || l.isSpacer(n) {
			continue
		}
		childSize := n.MinSize()
		switch l.orient {
		case Horizontal:
			minSize.Y = max(minSize.Y, childSize.Y)
			minSize.X += childSize.X
			if addPadding {
				minSize.X += l.InnerPadding()
			}
		case Vertical:
			minSize.X = max(minSize.X, childSize.X)
			minSize.Y += childSize.Y
			if addPadding {
				minSize.Y += l.InnerPadding()
			}
		}
		addPadding = true
	}
	return minSize
}

//-----------------------------------------------------------------------------

type MaxLayout struct{}

func NewMaxLayout() LayoutManager {
	return &MaxLayout{}
}

func (l *MaxLayout) Layout(childList *list.List, size Point) {
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		n.SetSize(size)
		n.SetPos(Point{0, 0})
	}
}

func (l *MaxLayout) MinSize(childList *list.List) Point {
	minSize := Point{0, 0}
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		if !n.IsVisible() {
			continue
		}
		minSize = minSize.Max(n.MinSize())
	}
	return minSize
}

//-----------------------------------------------------------------------------

// GridLayout ordnet die Kinder in einer bestimmten, fixen Anzahl Spalten
// (resp. Zeilen) an. Ueberschreitet die Anzahl der hinzugefuegten Kinder diese
// Groesse, dann wird eine weitere Zeile (resp. Spalte) erstellt und
// weitere Kinder analog zur ersten Zeile fortlaufend angeordnet.
type GridLayout struct {
	props.PropertyEmbed
	Cols   int
	orient Orientation
}

// Fixiert die Anzahl Spalten des GridLayouts.
func NewColumnGridLayout(cols int) LayoutManager {
	l := &GridLayout{Cols: cols, orient: Horizontal}
	l.InitProp("Layout")
	return l
}

// Fixiert die Anzahl Zeilen des GridLayouts.
func NewRowGridLayout(rows int) LayoutManager {
	l := &GridLayout{Cols: rows, orient: Vertical}
	l.InitProp("Layout")
	return l
}

func (l *GridLayout) horizontal() bool {
	return l.orient == Horizontal
}

func (l *GridLayout) countRows(childList *list.List) int {
	if l.Cols < 1 {
		l.Cols = 1
	}
	count := 0
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		if n.IsVisible() {
			count++
		}
	}
	return int(math.Ceil(float64(count) / float64(l.Cols)))
}

func (l *GridLayout) getLeading(size float64, offset int) float64 {
	return (size + l.Padding()) * float64(offset)
}

func (l *GridLayout) getTrailing(size float64, offset int) float64 {
	return l.getLeading(size, offset+1) - l.Padding()
}

func (l *GridLayout) Layout(childList *list.List, size Point) {
	rows := l.countRows(childList)
	padding := l.InnerPadding()
	padWidth := float64(l.Cols-1) * padding
	padHeight := float64(rows-1) * padding
	cellWidth := float64(size.X-padWidth) / float64(l.Cols)
	cellHeight := float64(size.Y-padHeight) / float64(rows)

	if !l.horizontal() {
		padWidth, padHeight = padHeight, padWidth
		cellWidth = float64(size.X-padWidth) / float64(rows)
		cellHeight = float64(size.Y-padHeight) / float64(l.Cols)
	}
	row, col := 0, 0
	i := 0
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		if !n.IsVisible() {
			continue
		}

		x1 := l.getLeading(cellWidth, col)
		y1 := l.getLeading(cellHeight, row)
		x2 := l.getTrailing(cellWidth, col)
		y2 := l.getTrailing(cellHeight, row)

		n.SetPos(Point{x1, y1})
		n.SetSize(Point{x2 - x1, y2 - y1})

		if l.horizontal() {
			if (i+1)%l.Cols == 0 {
				row++
				col = 0
			} else {
				col++
			}
		} else {
			if (i+1)%l.Cols == 0 {
				col++
				row = 0
			} else {
				row++
			}
		}
		i++
	}
}

func (l *GridLayout) MinSize(childList *list.List) Point {
	rows := l.countRows(childList)
	minSize := Point{0, 0}
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		if !n.IsVisible() {
			continue
		}
		minSize = minSize.Max(n.MinSize())
	}

	pad := l.Padding()
	if l.horizontal() {
		minContentSize := Point{minSize.X * float64(l.Cols),
			minSize.Y * float64(rows)}
		return minContentSize.Add(Point{pad * math.Max(float64(l.Cols-1), 0),
			pad * math.Max(float64(rows-1), 0.0)})
	}

	minContentSize := Point{minSize.X * float64(rows),
		minSize.Y * float64(l.Cols)}
	return minContentSize.Add(Point{pad * math.Max(float64(rows-1), 0.0),
		pad * math.Max(float64(l.Cols-1), 0.0)})
}

//-----------------------------------------------------------------------------

// Nimmt den verfügbaren Platz (vertikal oder horizontal) in Box-Layouts
// ein. Ist zwar ein Widget, passt aber irgendwie besser zum Layout-Zeugs.
type Spacer struct {
	nodeEmbed
	FixHorizontal, FixVertical bool
}

func NewSpacer() *Spacer {
	s := &Spacer{}
	s.Init(s)
	return s
}

func (s *Spacer) Draw(gc *gg.Context) {}

func (s *Spacer) ExpandHorizontal() bool {
	return !s.FixHorizontal
}

func (s *Spacer) ExpandVertical() bool {
	return !s.FixVertical
}
