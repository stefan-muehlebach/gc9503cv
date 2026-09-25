package main

import (
	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	//"log"
	"math"
	//"math/rand/v2"
	"slices"
	"time"
)

const (
	defNumDots = 200
)

type PolyhedronAnim struct {
	// eventHandlerEmbed
	windowEmbed
	m0                         Matrix
	objList                    []Object3D
	alpha, dAlpha, beta, dBeta float64
	zoom                       float64
	isManual                   bool
}

func NewPolyhedronAnimation(bounds geom.Rectangle[int],
	numDots int) *PolyhedronAnim {
	var obj *Polyhedron

	a := &PolyhedronAnim{}

	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	rect := a.bounds.Sub(a.bounds.Min).ToFloat()
	drawStyle := Faces

	if numDots <= 0 {
		numDots = defNumDots
	}

	a.Root = NewGroup()
	a.Root.SetSize(bounds.ToFloat().Size())

	mp := rect.Center()
	zero := NewVector(0.0, 0.0, 0.0)
	ex := NewVector(1.0, 0.0, 0.0)
	ey := NewVector(0.0, 1.0, 0.0)
	//ez := NewVector(0.0, 0.0, 1.0)

	a.objList = make([]Object3D, 0)
	//a.objList = append(a.objList, NewSegment(zero, ex, colors.Red, 4.0))
	//a.objList = append(a.objList, NewSegment(zero, ey, colors.Green, 4.0))
	//a.objList = append(a.objList, NewSegment(zero, ez, colors.Blue, 4.0))

	obj = NewPolyhedron(Cross, zero, 50.0, colors.GoWhite, 3.0)
	obj.DrawStyle = drawStyle
	a.objList = append(a.objList, obj)

	objectList := []ObjectData{Tetraeder, Hexaeder, Oktaeder, Ikosaeder, Dodekaeder}
	colorList := []colors.RGBA{colors.GoGopherBlue, colors.GoYellow,
		colors.GoAqua, colors.GoFuchsia, colors.GoLightBlue}

	for i, solid := range objectList {
		angle := float64(i) * (2.0 * math.Pi / 5.0)
		dx := ex.Mul(150.0 * math.Sin(angle))
		dy := ey.Mul(150.0 * math.Cos(angle))
		v := dx.Add(dy)

		obj = NewPolyhedron(solid, v, 50.0, colorList[i], 3.0)
		obj.DrawStyle = drawStyle
		a.objList = append(a.objList, obj)
	}

	/*
		obj = NewPolyhedron(Tetraeder, ey.Neg().Mul(450.0), 50.0,
			colors.GoGopherBlue, 3.0)
		obj.DrawStyle = drawStyle
		a.objList = append(a.objList, obj)

		obj = NewPolyhedron(Hexaeder, ey.Neg().Mul(300.0), 30.0,
			colors.GoYellow, 3.0)
		obj.DrawStyle = drawStyle
		a.objList = append(a.objList, obj)

		obj = NewPolyhedron(Oktaeder, ey.Neg().Mul(150.0), 50.0,
			colors.GoAqua, 3.0)
		obj.DrawStyle = drawStyle
		a.objList = append(a.objList, obj)

		obj = NewPolyhedron(Ikosaeder, ey.Mul(150.0), 50.0,
			colors.GoFuchsia, 3.0)
		obj.DrawStyle = drawStyle
		a.objList = append(a.objList, obj)

		obj = NewPolyhedron(Dodekaeder, ey.Mul(300.0), 50.0,
			colors.GoLightBlue, 3.0)
		obj.DrawStyle = drawStyle
		a.objList = append(a.objList, obj)
	*/

	//obj = NewPolyhedron(Staircube, ey.Mul(100.0), 50.0,
	//	colors.GoWhite, 3.0)
	//obj.DrawStyle = drawStyle
	//a.objList = append(a.objList, obj)

	//a.objList = append(a.objList, NewDot(0.0, 0.0, 0.0, colors.Red, 5.0))
	//a.objList = append(a.objList, NewDot(0.0, 100.0, 0.0, colors.Red, 5.0))
	//a.objList = append(a.objList, NewDot(0.0, -100.0, 0.0, colors.Red, 5.0))

	a.m0 = Identity().
		Mul(Scale(NewVector(1.0, -1.0, 1.0))).
		Mul(Translate(NewVector(mp.X, -mp.Y, 0.0)))

	a.dAlpha = math.Pi / 250.0
	a.dBeta = math.Pi / 180.0
	a.zoom = 1.0

	var lastPos Point

	a.Root.SetOnPress(func(ev InputEvent) {
		a.isManual = true
		lastPos = ev.Pos
	})
	a.Root.SetOnDrag(func(ev InputEvent) {
		dPos := ev.Pos.Sub(lastPos)
		a.beta += math.Pi * (dPos.X / 300.0)
		a.alpha += math.Pi * (dPos.Y / 300.0)
		lastPos = ev.Pos
	})
	a.Root.SetOnWheel(func(ev InputEvent) {
		a.zoom = float64(ev.Wheel) / 10.0
	})

	return a
}

func (a *PolyhedronAnim) Init() {
	a.alpha = 0.0 // math.Pi / 12.0
	a.beta = 0.0  // math.Pi / 18.0
}

func (a *PolyhedronAnim) Update(dt time.Duration) {
	m := a.m0.
		Mul(Scale(Vector{a.zoom, a.zoom, a.zoom})).
		Mul(RotateX(a.alpha)).
		Mul(RotateY(a.beta))
	for _, obj := range a.objList {
		obj.Transform(m)
	}
	if a.isManual {
		return
	}
	a.alpha += a.dAlpha
	if a.alpha >= 2*math.Pi {
		a.alpha -= 2 * math.Pi
	}
	a.beta += a.dBeta
	if a.beta >= 2*math.Pi {
		a.beta -= 2 * math.Pi
	}
}

func (a *PolyhedronAnim) Refresh() {
	a.gc.Clear(colors.Black)
	slices.SortFunc(a.objList, func(a, b Object3D) int {
		if a.Z() < b.Z() {
			return -1
		} else if a.Z() > b.Z() {
			return +1
		} else {
			return 0
		}
	})
	for _, obj := range a.objList {
		obj.Draw(a.gc)
	}
}

//----------------------------------------------------------------------------

type Object3D interface {
	Transform(m Matrix)
	Draw(gc *gg.Context)
	Z() float64
}

type Node3D interface {
	Translate(t Vector)
	Scale(s Vector)
	Rotate(axis AxisType, angle float64)
}

type AxisType int

const (
	XAxis AxisType = iota
	YAxis
	ZAxis
)

//----------------------------------------------------------------------------

type node3DEmbed struct {
	m Matrix
}

func (n *node3DEmbed) Translate(t Vector) {
	n.m = Translate(t).Mul(n.m)
}

func (n *node3DEmbed) Scale(s Vector) {
	n.m = Scale(s).Mul(n.m)
}

func (n *node3DEmbed) Rotate(axis AxisType, angle float64) {
	switch axis {
	case XAxis:
		n.m = RotateX(angle).Mul(n.m)
	case YAxis:
		n.m = RotateY(angle).Mul(n.m)
	case ZAxis:
		n.m = RotateZ(angle).Mul(n.m)
	}
}

//----------------------------------------------------------------------------

type ObjectData struct {
	Points []Vector
	Edges  [][]int
	Faces  [][]int
	Holes  [][]int
}

var (
	// The golden ratio
	phi = (math.Sqrt(5) + 1) / 2

	// Bei den folgenden Polyedern (den platonischen Koerpern) wird davon
	// ausgegangen, dass die Seitenlaenge der begrenzenden Flaechen immer
	// Laenge = 1 hat. Anschl. wird das Polyeder so positioniert und skaliert
	// dass der Mittelpunkt des Koerpers (0,0,0) ist und der gesamte Koerper
	// in einen Wuerfel von (-1,-1,-1) bis (1,1,1) passt.
	//
	// Groessen, Punkte und Flaechendefinitionen fuer....
	// Tetraeder
	A0 = 1.0 / 2.0
	B0 = math.Sqrt(2.0) / 4.0

	f0 = 1.0 / max(A0, B0)
	a0 = f0 * A0
	b0 = f0 * B0

	Tetraeder = ObjectData{
		Points: []Vector{
			Vector{a0, 0, -b0},
			Vector{-a0, 0, -b0},
			Vector{0, a0, b0},
			Vector{0, -a0, b0},
		},
		Edges: [][]int{
			{0, 1, 2, 3},
			{2, 0, 3, 1},
		},
		Faces: [][]int{
			{0, 1, 2},
			{0, 3, 1},
			{0, 2, 3},
			{1, 3, 2},
		},
	}

	// Hexaeder (oder einfach Wuerfel, gell...)
	A1 = 1.0 / 2.0

	f1 = 1.0 / A1
	a1 = f1 * A1

	Hexaeder = ObjectData{
		Points: []Vector{
			Vector{-a1, -a1, -a1},
			Vector{-a1, a1, -a1},
			Vector{a1, a1, -a1},
			Vector{a1, -a1, -a1},
			Vector{-a1, -a1, a1},
			Vector{-a1, a1, a1},
			Vector{a1, a1, a1},
			Vector{a1, -a1, a1},
		},
		Edges: [][]int{
			{0, 1, 2, 3, 0},
			{4, 5, 6, 7, 4},
			{0, 4}, {1, 5}, {2, 6}, {3, 7},
		},
		Faces: [][]int{
			{0, 1, 2, 3},
			{0, 4, 5, 1},
			{1, 5, 6, 2},
			{2, 6, 7, 3},
			{3, 7, 4, 0},
			{4, 7, 6, 5},
		},
	}

	// Oktaeder
	A2 = 1.0 / 2.0
	E2 = math.Sqrt(2.0) / 2.0

	f2 = 1.0 / max(A2, E2)
	a2 = f2 * A2
	e2 = f2 * E2

	Oktaeder = ObjectData{
		Points: []Vector{
			Vector{0, 0, -e2},
			Vector{a2, -a2, 0},
			Vector{a2, a2, 0},
			Vector{-a2, a2, 0},
			Vector{-a2, -a2, 0},
			Vector{0, 0, e2},
		},
		Edges: [][]int{
			{0, 1}, {0, 2}, {0, 3}, {0, 4},
			{1, 2}, {2, 3}, {3, 4}, {4, 1},
			{5, 1}, {5, 2}, {5, 3}, {5, 4},
		},
		Faces: [][]int{
			{0, 2, 1},
			{0, 3, 2},
			{0, 4, 3},
			{0, 1, 4},
			{1, 2, 5},
			{2, 3, 5},
			{3, 4, 5},
			{4, 1, 5},
		},
	}

	// Ikosaeder
	A3 = 1.0 / 2.0
	C3 = (1.0 + math.Sqrt(5.0)) / 4.0

	f3 = 1.0 / max(A3, C3)
	a3 = f3 * A3
	c3 = f3 * C3

	Ikosaeder = ObjectData{
		Points: []Vector{
			Vector{0.0, -a3, -c3}, // 0
			Vector{0.0, a3, -c3},  // 1
			Vector{-c3, 0.0, -a3}, // 2
			Vector{c3, 0.0, -a3},  // 3

			Vector{-a3, -c3, 0.0}, // 4
			Vector{a3, -c3, 0.0},  // 5
			Vector{a3, c3, 0.0},   // 6
			Vector{-a3, c3, 0.0},  // 7

			Vector{-c3, 0.0, a3}, // 8
			Vector{c3, 0.0, a3},  // 9
			Vector{0.0, -a3, c3}, // 10
			Vector{0.0, a3, c3},  // 11
		},
		Edges: [][]int{
			{0, 1, 2, 0, 3, 1},
			{10, 11, 9, 10, 8, 11},
			{4, 5, 0, 4, 10, 5},
			{6, 7, 1, 6, 11, 7},
			{9, 3, 6, 9, 5, 3},
			{8, 2, 7, 8, 4, 2},
		},
		Faces: [][]int{
			{0, 1, 3},
			{0, 3, 5},
			{0, 5, 4},
			{0, 4, 2},
			{0, 2, 1},
			{1, 2, 7},
			{2, 8, 7},
			{2, 4, 8},
			{4, 10, 8},
			{4, 5, 10},
			{5, 9, 10},
			{5, 3, 9},
			{3, 6, 9},
			{3, 1, 6},
			{1, 7, 6},
			{8, 10, 11},
			{10, 9, 11},
			{9, 6, 11},
			{6, 7, 11},
			{7, 8, 11},
		},
	}

	// Dodekaeder
	A4 = 1.0 / 2.0
	B4 = (1.0 + phi) / 2.0
	C4 = phi / 2.0

	f4 = 1.0 / max(A4, max(B4, C4))
	a4 = f4 * A4
	b4 = f4 * B4
	c4 = f4 * C4

	Dodekaeder = ObjectData{
		Points: []Vector{
			Vector{a4, 0, -b4},    // 0
			Vector{-a4, 0, -b4},   // 1
			Vector{c4, c4, -c4},   // 2
			Vector{-c4, c4, -c4},  // 3
			Vector{-c4, -c4, -c4}, // 4
			Vector{c4, -c4, -c4},  // 5
			Vector{0, b4, -a4},    // 6
			Vector{0, -b4, -a4},   // 7
			Vector{b4, -a4, 0},    // 8
			Vector{b4, a4, 0},     // 9
			Vector{-b4, a4, 0},    // 10
			Vector{-b4, -a4, 0},   // 11
			Vector{0, b4, a4},     // 12
			Vector{0, -b4, a4},    // 13
			Vector{c4, c4, c4},    // 14
			Vector{-c4, c4, c4},   // 15
			Vector{-c4, -c4, c4},  // 16
			Vector{c4, -c4, c4},   // 17
			Vector{a4, 0, b4},     // 18
			Vector{-a4, 0, b4},    // 19
		},
		Edges: [][]int{
			{0, 1, 4, 7, 5, 0, 2, 6, 3, 1},
			{18, 19, 16, 13, 17, 18, 14, 12, 15, 19},
			{5, 8, 17}, {2, 9, 14}, {3, 10, 15}, {4, 11, 16},
			{7, 13}, {6, 12}, {8, 9}, {10, 11},
		},
		Faces: [][]int{
			{0, 2, 9, 8, 5},
			{1, 4, 11, 10, 3},
			{0, 1, 3, 6, 2},
			{0, 5, 7, 4, 1},
			{2, 6, 12, 14, 9},
			{3, 10, 15, 12, 6},
			{5, 8, 17, 13, 7},
			{4, 7, 13, 16, 11},
			{8, 9, 14, 18, 17},
			{10, 11, 16, 19, 15},
			{12, 15, 19, 18, 14},
			{13, 17, 18, 19, 16},
		},
	}

	// Schweizerkreuz-Box

	A5 = 16.0
	B5 = 10.0
	C5 = 3.0
	D5 = 3.0

	f5 = 1.0 / max(B5, max(C5, D5))
	b5 = f5 * B5
	c5 = f5 * C5
	d5 = f5 * D5

	Cross = ObjectData{
		Points: []Vector{
			Vector{-b5, -c5, -d5}, // 0
			Vector{-b5, c5, -d5},  // 1
			Vector{-c5, c5, -d5},  // 2
			Vector{-c5, b5, -d5},  // 3
			Vector{c5, b5, -d5},   // 4
			Vector{c5, c5, -d5},   // 5
			Vector{b5, c5, -d5},   // 6
			Vector{b5, -c5, -d5},  // 7
			Vector{c5, -c5, -d5},  // 8
			Vector{c5, -b5, -d5},  // 9
			Vector{-c5, -b5, -d5}, // 10
			Vector{-c5, -c5, -d5}, // 11

			Vector{-b5, -c5, d5}, // 12
			Vector{-b5, c5, d5},  // 13
			Vector{-c5, c5, d5},  // 14
			Vector{-c5, b5, d5},  // 15
			Vector{c5, b5, d5},   // 16
			Vector{c5, c5, d5},   // 17
			Vector{b5, c5, d5},   // 18
			Vector{b5, -c5, d5},  // 19
			Vector{c5, -c5, d5},  // 20
			Vector{c5, -b5, d5},  // 21
			Vector{-c5, -b5, d5}, // 22
			Vector{-c5, -c5, d5}, // 23

			Vector{-c5, c5, -b5},  // 24
			Vector{c5, c5, -b5},   // 25
			Vector{c5, -c5, -b5},  // 26
			Vector{-c5, -c5, -b5}, // 27

			Vector{-c5, c5, b5},  // 28
			Vector{c5, c5, b5},   // 29
			Vector{c5, -c5, b5},  // 30
			Vector{-c5, -c5, b5}, // 31
		},
		Edges: [][]int{
			// Entlang X-Achse
			{0, 1, 6, 7, 0}, {12, 13, 18, 19, 12}, {0, 12}, {1, 13}, {6, 18},
			{7, 19},
			// Entlang Y-Achse
			{3, 4, 9, 10, 3}, {15, 16, 21, 22, 15}, {3, 15}, {4, 16}, {9, 21},
			{10, 22},
			// Entlang Z-Achse
			{24, 25, 29, 28, 24}, {26, 27, 31, 30, 26}, {25, 26}, {24, 27},
			{29, 30}, {28, 31},
		},
		Faces: [][]int{
			// Untere Kreuzflaeche
			{0, 1, 2, 11}, {2, 3, 4, 5}, {5, 6, 7, 8}, {8, 9, 10, 11},
			// Obere Kreuzflaeche
			{12, 23, 14, 13}, {14, 17, 16, 15}, {17, 20, 19, 18},
			{20, 23, 22, 21},
			// Seitliche Begrenzungsflaechen
			{1, 0, 12, 13}, {2, 1, 13, 14}, {3, 2, 14, 15}, {4, 3, 15, 16},
			{5, 4, 16, 17}, {6, 5, 17, 18}, {7, 6, 18, 19}, {8, 7, 19, 20},
			{9, 8, 20, 21}, {10, 9, 21, 22}, {11, 10, 22, 23},
			{0, 11, 23, 12},
			// Unterer Turm
			{2, 24, 27, 11}, {2, 5, 25, 24}, {5, 8, 26, 25}, {8, 11, 27, 26},
			{24, 25, 26, 27},
			// Oberer Turn
			{14, 23, 31, 28}, {17, 14, 28, 29}, {20, 17, 29, 30},
			{23, 20, 30, 31}, {28, 31, 30, 29},
		},
	}

	// Staircube
	Staircube = ObjectData{
		Points: []Vector{
			Vector{0, 0, 0}, // 0
			Vector{1, 0, 0}, // 1
			Vector{1, 1, 0}, // 2
			Vector{0, 1, 0}, // 3

			Vector{1, 0, 0.2}, // 4
			Vector{1, 1, 0.2}, // 5
			Vector{0, 1, 0.2}, // 6

			Vector{0.8, 0, 0.2},   // 7
			Vector{0.8, 0.8, 0.2}, // 8
			Vector{0, 0.8, 0.2},   // 9

			Vector{0.8, 0, 0.4},   // 10
			Vector{0.8, 0.8, 0.4}, // 11
			Vector{0, 0.8, 0.4},   // 12

			Vector{0.6, 0, 0.4},   // 13
			Vector{0.6, 0.6, 0.4}, // 14
			Vector{0, 0.6, 0.4},   // 15

			Vector{0.6, 0, 0.6}, // 16
			Vector{0.6, 0.6, 0.6},
			Vector{0, 0.6, 0.6},

			Vector{0.4, 0, 0.6}, // 19
			Vector{0.4, 0.4, 0.6},
			Vector{0, 0.4, 0.6},

			Vector{0.4, 0, 0.8}, // 22
			Vector{0.4, 0.4, 0.8},
			Vector{0, 0.4, 0.8},

			Vector{0.2, 0, 0.8}, // 25
			Vector{0.2, 0.2, 0.8},
			Vector{0, 0.2, 0.8},

			Vector{0.2, 0, 1.0}, // 28
			Vector{0.2, 0.2, 1.0},
			Vector{0, 0.2, 1.0},

			Vector{0, 0, 1.0}, // 31
		},
		Edges: [][]int{
			{0, 1, 2, 3, 0},
			{4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}, {16, 17, 18}, {19, 20, 21},
			{22, 23, 24}, {25, 26, 27}, {28, 29, 30, 31, 28},
		},
		Faces: [][]int{
			{0, 3, 2, 1},
			{1, 2, 5, 4}, {2, 3, 6, 5},
			{4, 5, 6, 9, 8, 7},
			{7, 8, 11, 10}, {8, 9, 12, 11},
			{10, 11, 12, 15, 14, 13},
			{13, 14, 17, 16}, {14, 15, 18, 17},
			{16, 17, 18, 21, 20, 19},
			{19, 20, 23, 22}, {20, 21, 24, 23},
			{22, 23, 24, 27, 26, 25},
			{25, 26, 29, 28}, {26, 27, 30, 29},
			{28, 29, 30, 31, 28},
			{0, 1, 4, 7, 10, 13, 16, 19, 22, 25, 28, 31},
			{0, 31, 30, 27, 24, 21, 18, 15, 12, 9, 6, 3},
		},
	}

/*
	Flag = ObjectData{
		Points: []Vector{
			Vector{a5, a5, -d5},
			Vector{-a5, a5, -d5},
			Vector{-a5, -a5, -d5},
			Vector{a5, -a5, -d5},

			Vector{a5, a5, d5},
			Vector{-a5, a5, d5},
			Vector{-a5, -a5, d5},
			Vector{a5, -a5, d5},

			Vector{-b5, -c5, -d5},
			Vector{-b5, c5, -d5},
			Vector{-c5, c5, -d5},
			Vector{-c5, b5, -d5},
			Vector{c5, b5, -d5},
			Vector{c5, c5, -d5},
			Vector{b5, c5, -d5},
			Vector{b5, -c5, -d5},
			Vector{c5, -c5, -d5},
			Vector{c5, -b5, -d5},
			Vector{-c5, -b5, -d5},
			Vector{-c5, -c5, -d5},

			Vector{-b5, -c5, d5},
			Vector{-b5, c5, d5},
			Vector{-c5, c5, d5},
			Vector{-c5, b5, d5},
			Vector{c5, b5, d5},
			Vector{c5, c5, d5},
			Vector{b5, c5, d5},
			Vector{b5, -c5, d5},
			Vector{c5, -c5, d5},
			Vector{c5, -b5, d5},
			Vector{-c5, -b5, d5},
			Vector{-c5, -c5, d5},
		},
		Edges: [][]int{
			{0, 1, 2, 3, 0},
			{4, 5, 6, 7, 4},
			{0, 4}, {1, 5}, {2, 6}, {3, 7},
		},
		Faces: [][]int{
			// Vorder- und Rueckseite
			{3, 2, 1, 0},
			{4, 5, 6, 7},
			// Umlaufende Seitenflaechen
			{0, 1, 5, 4},
			{1, 2, 6, 5},
			{2, 3, 7, 6},
			{3, 0, 4, 7},
			// Begrenzungsflaechen zum kreuzfoermigen Loch
			{8, 9, 21, 20},
			{9, 10, 22, 21},
			{10, 11, 23, 22},
			{11, 12, 24, 23},
			{12, 13, 25, 24},
			{13, 14, 26, 25},
			{14, 15, 27, 26},
			{15, 16, 28, 27},
			{16, 17, 29, 28},
			{17, 18, 30, 29},
			{18, 19, 31, 30},
			{19, 8, 20, 31},
		},
		Holes: [][]int{
			// Kreuzfoermige Ausschnitte in Vorder- und Rueckseite
			{19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8},
			{20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
		},
	}
*/
)

//-----------------------------------------------------------------------

type DrawStyleType int

const (
	Faces DrawStyleType = iota
	Edges
	Corners
)

//-----------------------------------------------------------------------

type Polyhedron struct {
	node3DEmbed
	Pts, PtsT  []Vector
	Faces      []*Face
	Edges      [][]int
	ZMin, ZMax float64
	LineWidth  float64
	LineColor  colors.RGBA
	FillColor  colors.RGBA
	DrawStyle  DrawStyleType
}

func NewPolyhedron(data ObjectData, pos Vector,
	size float64, color colors.RGBA, lineWidth float64) *Polyhedron {
	o := &Polyhedron{}
	o.LineWidth = lineWidth
	o.LineColor = color
	o.FillColor = color.Dark(0.1)
	o.Pts = make([]Vector, len(data.Points))
	o.PtsT = make([]Vector, len(data.Points))
	o.Edges = make([][]int, len(data.Edges))
	o.Faces = make([]*Face, len(data.Faces))

	o.m = Identity()
	o.Scale(Vector{size, size, size})
	o.Translate(pos)
	for i, pt := range data.Points {
		o.Pts[i] = pt
	}
	for i, edge := range data.Edges {
		o.Edges[i] = edge
	}
	for i, face := range data.Faces {
		o.Faces[i] = NewFace(o, face, color, lineWidth)
	}
	for i, holes := range data.Holes {
		o.Faces[i].Hole = holes
	}
	return o
}

func (o *Polyhedron) Transform(m Matrix) {
	M := m.Mul(o.m)
	o.ZMin = +math.MaxFloat64
	o.ZMax = -math.MaxFloat64
	for i, pt := range o.Pts {
		ptT := M.TransformAffine(pt)
		if ptT.Z < o.ZMin {
			o.ZMin = ptT.Z
		}
		if ptT.Z > o.ZMax {
			o.ZMax = ptT.Z
		}
		o.PtsT[i] = ptT
	}
	for _, face := range o.Faces {
		face.Transform(M)
	}
}

func (o *Polyhedron) Z() float64 {
	return o.ZMax
}

func (o *Polyhedron) Draw(gc *gg.Context) {
	switch o.DrawStyle {
	case Faces:
		slices.SortFunc(o.Faces, func(a, b *Face) int {
			if a.Z() < b.Z() {
				return -1
			} else if a.Z() > b.Z() {
				return +1
			} else {
				return 0
			}
		})
		for _, face := range o.Faces {
			face.Draw(gc)
		}
	case Edges:
		gc.SetLineWidth(o.LineWidth)
		gc.SetLineColor(o.LineColor)
		for _, edge := range o.Edges {
			pt := o.PtsT[edge[0]]
			gc.MoveTo(pt.X, pt.Y)
			for _, idx := range edge[1:] {
				pt := o.PtsT[idx]
				gc.LineTo(pt.X, pt.Y)
			}
			gc.Stroke()
		}
	case Corners:
		gc.SetFillColor(o.LineColor)
		for _, pt := range o.PtsT {
			gc.DrawPoint(pt.X, pt.Y, 2.0)
			gc.Fill()
		}
	}
}

//-----------------------------------------------------------------------

type Face struct {
	o           *Polyhedron
	Idx         []int
	Hole        []int
	Norm, NormT Vector
	ZMax, ZMin  float64
	LineWidth   float64
	LineColor   colors.RGBA
	FillColor   colors.RGBA
}

func NewFace(o *Polyhedron, idx []int, color colors.RGBA,
	lineWidth float64) *Face {
	f := &Face{}

	f.o = o
	f.Idx = idx
	f.LineWidth = lineWidth
	f.LineColor = color
	f.FillColor = color.Dark(0.1)

	v1 := o.Pts[idx[2]].Sub(o.Pts[idx[1]])
	v2 := o.Pts[idx[0]].Sub(o.Pts[idx[1]])
	f.Norm = v1.Cross(v2)

	return f
}

func (f *Face) Transform(m Matrix) {
	f.NormT = m.TransformLinear(f.Norm)
}

func (f *Face) Z() float64 {
	f.ZMin = +math.MaxFloat64
	f.ZMax = -math.MaxFloat64
	for _, idx := range f.Idx {
		z := f.o.PtsT[idx].Z
		if z < f.ZMin {
			f.ZMin = z
		}
		if z > f.ZMax {
			f.ZMax = z
		}
	}
	return f.ZMax
}

func (f *Face) Draw(gc *gg.Context) {
	if f.NormT.Z < 0 {
		return
	}
	gc.SetLineWidth(f.LineWidth)
	gc.SetLineColor(f.LineColor)
	gc.SetFillColor(f.FillColor.Dark((1.0 - f.NormT.Z/f.NormT.Abs()) / 2.0))
	gc.MoveTo(f.o.PtsT[f.Idx[0]].X, f.o.PtsT[f.Idx[0]].Y)
	for _, idx := range f.Idx[1:] {
		gc.LineTo(f.o.PtsT[idx].X, f.o.PtsT[idx].Y)
	}
	gc.ClosePath()
	if len(f.Hole) > 0 {
		gc.SetFillRule(gg.FillRuleEvenOdd)
		gc.MoveTo(f.o.PtsT[f.Hole[0]].X, f.o.PtsT[f.Hole[0]].Y)
		for _, idx := range f.Hole[1:] {
			gc.LineTo(f.o.PtsT[idx].X, f.o.PtsT[idx].Y)
		}
		gc.ClosePath()
	}
	gc.FillStroke()
}

//-----------------------------------------------------------------------

/*
type Dot struct {
	Pt, PtT Vector
	Color   colors.RGBA
	Size    float64
}

func NewDot(x, y, z float64, color colors.RGBA, size float64) *Dot {
	d := &Dot{}
	d.Pt = Vector{x, y, z}
	d.Color = color
	d.Size = size
	return d
}

func (d *Dot) Transform(m Matrix) {
	d.PtT = m.TransformAffine(d.Pt)
}

func (d *Dot) Z() float64 {
	return d.PtT.Z
}

func (d *Dot) Draw(gc *gg.Context) {
	dark := Map(d.PtT.Z, -200.0, 200.0, 0.7, 0.2)
	gc.SetFillColor(d.Color.Dark(dark))
	gc.DrawPoint(d.PtT.X, d.PtT.Y, d.Size)
	gc.Fill()
}
*/

//-----------------------------------------------------------------------

type Segment struct {
	Pts, PtsT [2]Vector
	Color     colors.RGBA
	Width     float64
}

func NewSegment(p0, p1 Vector, color colors.RGBA, width float64) *Segment {
	s := &Segment{}
	s.Pts[0] = p0
	s.Pts[1] = p1
	s.Color = color
	s.Width = width
	return s
}

func (s *Segment) Transform(m Matrix) {
	s.PtsT[0] = m.TransformAffine(s.Pts[0])
	s.PtsT[1] = m.TransformAffine(s.Pts[1])
}

func (s *Segment) Z() float64 {
	return max(s.PtsT[0].Z, s.PtsT[1].Z)
}

func (s *Segment) Draw(gc *gg.Context) {
	gc.SetLineColor(s.Color)
	gc.SetLineWidth(s.Width)
	gc.DrawLine(s.PtsT[0].X, s.PtsT[0].Y, s.PtsT[1].X, s.PtsT[1].Y)
	gc.Stroke()
}

//-----------------------------------------------------------------------

/*
func NewCloud(x, y, z, w float64, numObjs int, color colors.RGBA,
	size float64) []Object3D {

	dots := make([]Object3D, 0)
	for range numObjs {
		px := rand.NormFloat64()*w + x
		py := rand.NormFloat64()*w + y
		pz := rand.NormFloat64()*w + z
		dots = append(dots, NewDot(px, py, pz, color, size))
	}
	return dots
}
*/

//-----------------------------------------------------------------------

type Mappable interface {
	~int | ~int16 | ~float64
}

func Map[I, O Mappable](valIn, lbIn, ubIn I, lbOut, ubOut O) (valOut O) {
	return lbOut + O(valIn-lbIn)*(ubOut-lbOut)/O(ubIn-lbIn)
}
