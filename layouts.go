package main

import (
	"container/list"
	"math"

	. "github.com/stefan-muehlebach/gc9503cv/props"
	"github.com/stefan-muehlebach/gg"
)

//-----------------------------------------------------------------------------

var (
	LayoutProps = PropsMap["Layout"]
)

//-----------------------------------------------------------------------------

type LayoutManager interface {
	Layout(childList *list.List, size Point)
	MinSize(childList *list.List) Point
}

//-----------------------------------------------------------------------------

type NullLayout struct{}

func (l *NullLayout) Layout(childList *list.List, size Point) {}

func (l *NullLayout) MinSize(childList *list.List) Point {
	minSize := Point{}
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		if !n.IsVisible() {
			continue
		}
		minSize = minSize.Max(n.Bounds().Max)
	}
	return minSize
}

//-----------------------------------------------------------------------------

type PadLayout struct {
	pad [4]float64
}

func NewPadLayout(pads ...float64) *PadLayout {
	var l, t, r, b float64

	switch len(pads) {
	case 0:
		s := LayoutProps.Size(Padding)
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
	return &PadLayout{[4]float64{l, t, r, b}}
}

func (l *PadLayout) Layout(childList *list.List, size Point) {
	pos := Point{l.pad[0], l.pad[1]}
	siz := Point{size.X - l.pad[0] - l.pad[2],
		size.Y - l.pad[1] - l.pad[3]}
	for e := childList.Front(); e != nil; e = e.Next() {
		node := e.Value.(Node)
		node.SetSize(siz)
		node.SetPos(pos)
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
	return minSize.Add(Point{l.pad[0] + l.pad[2], l.pad[1] + l.pad[3]})
}

//-----------------------------------------------------------------------------

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

type BoxLayout struct {
	orient Orientation
	pad    float64
}

func NewHBoxLayout(pads ...float64) *BoxLayout {
	pad := LayoutProps.Size(InnerPadding)
	if len(pads) > 0 {
		pad = pads[0]
	}
	return &BoxLayout{Horizontal, pad}
}

func NewVBoxLayout(pads ...float64) *BoxLayout {
	pad := LayoutProps.Size(InnerPadding)
	if len(pads) > 0 {
		pad = pads[0]
	}
	return &BoxLayout{Vertical, pad}
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
		node := e.Value.(Node)
		if !node.IsVisible() {
			continue
		}
		if l.isSpacer(node) {
			spacers++
			continue
		}
		nodeSize := node.MinSize()
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
		extra = size.X - total - l.pad*float64(childList.Len()-
			spacers-1)
	case Vertical:
		extra = size.Y - total - l.pad*float64(childList.Len()-
			spacers-1)
	}
	if spacers > 0 {
		extraCell = extra / float64(spacers)
	}
	pos := Point{}
	for e := childList.Front(); e != nil; e = e.Next() {
		node := e.Value.(Node)
		if !node.IsVisible() {
			continue
		}
		if l.isSpacer(node) {
			switch l.orient {
			case Horizontal:
				pos.X += extraCell
			case Vertical:
				pos.Y += extraCell
			}
			continue
		}
		node.SetPos(pos)
		switch l.orient {
		case Horizontal:
			width := node.MinSize().X
			pos.X += width + l.pad
			node.SetSize(Point{width, size.Y})
		case Vertical:
			height := node.MinSize().Y
			pos.Y += height + l.pad
			node.SetSize(Point{size.X, height})
		}
	}
}

func (l *BoxLayout) MinSize(childList *list.List) Point {
	//log.Printf("BoxLayout.MinSize()")
	minSize := Point{}
	addPadding := false
	for e := childList.Front(); e != nil; e = e.Next() {
		node := e.Value.(Node)
		if !node.IsVisible() || l.isSpacer(node) {
			continue
		}
		childSize := node.MinSize()
		switch l.orient {
		case Horizontal:
			minSize.Y = max(minSize.Y, childSize.Y)
			minSize.X += childSize.X
			if addPadding {
				minSize.X += l.pad
			}
		case Vertical:
			minSize.X = max(minSize.X, childSize.X)
			minSize.Y += childSize.Y
			if addPadding {
				minSize.Y += l.pad
			}
		}
		addPadding = true
	}
	return minSize
}

//-----------------------------------------------------------------------------

// GridLayout ordnet die Kinder in einer bestimmten, fixen Anzahl Spalten
// (resp. Zeilen) an. Ueberschreitet die Anzahl der hinzugefuegten Kinder diese
// Groesse, dann wird eine weitere Zeile (resp. Spalte) erstellt und
// weitere Kinder analog zur ersten Zeile fortlaufend angeordnet.
type GridLayout struct {
	Cols   int
	orient Orientation
}

// Fixiert die Anzahl Spalten des GridLayouts.
func NewColumnGridLayout(cols int) LayoutManager {
	l := &GridLayout{Cols: cols, orient: Horizontal}
	return l
}

// Fixiert die Anzahl Zeilen des GridLayouts.
func NewRowGridLayout(rows int) LayoutManager {
	l := &GridLayout{Cols: rows, orient: Vertical}
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
		node := e.Value.(Node)
		if node.IsVisible() {
			count++
		}
	}
	return int(math.Ceil(float64(count) / float64(l.Cols)))
}

func getLeading(size float64, offset int) float64 {
	return (size + float64(LayoutProps.Size(Padding))) * float64(offset)
}

func getTrailing(size float64, offset int) float64 {
	return getLeading(size, offset+1) - LayoutProps.Size(Padding)
}

func (l *GridLayout) Layout(childList *list.List, size Point) {
	rows := l.countRows(childList)
	padding := LayoutProps.Size(Padding)
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
		node := e.Value.(Node)
		if !node.IsVisible() {
			continue
		}

		x1 := getLeading(cellWidth, col)
		y1 := getLeading(cellHeight, row)
		x2 := getTrailing(cellWidth, col)
		y2 := getTrailing(cellHeight, row)

		node.SetPos(Point{x1, y1})
		node.SetSize(Point{x2 - x1, y2 - y1})

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
		node := e.Value.(Node)
		if !node.IsVisible() {
			continue
		}
		minSize = minSize.Max(node.MinSize())
	}

	pad := LayoutProps.Size(Padding)
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
