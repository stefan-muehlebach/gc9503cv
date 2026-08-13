package main

import (
	"math"
	"math/rand/v2"
	"slices"
	"time"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	//"github.com/stefan-muehlebach/gg/geom"
)

const (
	defNumDots = 200
)

type PlatonicAnim struct {
	callbackEmbed
	m0  Matrix
	objList []Object3D
	alpha, dAlpha, beta, dBeta    float64
	minZoom, maxZoom, zoom, dZoom float64
}

func NewPlatonicAnimation(numDots int, rect Rectangle) *PlatonicAnim {
	a := &PlatonicAnim{}

	if numDots <= 0 {
		numDots = defNumDots
	}

	mp := rect.Center()
	zero := NewVector(0.0, 0.0, 0.0)
	//ex := NewVector(100.0, 0.0, 0.0)
	//ey := NewVector(0.0, 100.0, 0.0)
	//ez := NewVector(0.0, 0.0, 100.0)

	a.objList = make([]Object3D, 0)
	//a.objList = append(a.objList, NewSegment(zero, ex, colors.Red, 4.0))
	//a.objList = append(a.objList, NewSegment(zero, ey, colors.Green, 4.0))
	//a.objList = append(a.objList, NewSegment(zero, ez, colors.Blue, 4.0))
/*
    obj := NewPlatonicSolid(TetraederPoints, TetraederFaces, ey.Neg().Mul(4),
		100.0, colors.Plum, 3.0)
	a.objList = append(a.objList, obj)
	for _, f := range obj.FaceList() {
		a.objList = append(a.objList, f)
	}
	obj = NewPlatonicSolid(HexaederPoints, HexaederFaces, ey.Neg().Mul(2),
		100.0, colors.PaleTurquoise, 3.0)
	a.objList = append(a.objList, obj)
	for _, f := range obj.FaceList() {
		a.objList = append(a.objList, f)
	}
    obj = NewPlatonicSolid(OktaederPoints, OktaederFaces, zero,
		100.0, colors.LightSalmon, 3.0)
	a.objList = append(a.objList, obj)
	for _, f := range obj.FaceList() {
		a.objList = append(a.objList, f)
	}
	obj = NewPlatonicSolid(IkosaederPoints, IkosaederFaces, ey.Mul(2),
		100.0, colors.PaleGreen, 3.0)
	a.objList = append(a.objList, obj)
	for _, f := range obj.FaceList() {
		a.objList = append(a.objList, f)
	}
	obj = NewPlatonicSolid(DodekaederPoints, DodekaederFaces, ey.Mul(4),
		100.0, colors.Khaki, 3.0)
	a.objList = append(a.objList, obj)
	for _, f := range obj.FaceList() {
		a.objList = append(a.objList, f)
	}
*/
	obj := NewPlatonicSolid(SwissPoints, SwissFaces, zero,
		300.0, colors.LightPink, 3.0)
	obj.SetHoles(SwissHoles)
	a.objList = append(a.objList, obj)
	for _, f := range obj.FaceList() {
		a.objList = append(a.objList, f)
	}

	//a.objList = append(a.objList, NewCloud(0.0, 100.0, 0.0, 50.0,
	//	numDots, colors.YellowGreen, 4.0)...)

	a.m0 = Identity().
		Mul(Scale(NewVector(1.0, -1.0, 1.0))).
		Mul(Translate(NewVector(mp.X, -mp.Y, 0.0)))

	a.dAlpha = math.Pi / 162.0
	a.dBeta = math.Pi / 106.0

	a.minZoom = 0.4
	a.maxZoom = 1.5
	a.zoom = 1.0
	//a.dZoom = 0.005

	return a
}

func (a *PlatonicAnim) Init(gc *gg.Context) {
	a.alpha = 0.0 // math.Pi / 12.0
	a.beta = 0.0 // math.Pi / 18.0
}

func (a *PlatonicAnim) Update(dt time.Duration) {
	m := a.m0.
		//Mul(Scale(NewVector(a.zoom, a.zoom, 1.0))).
		Mul(RotateX(a.alpha)).
		Mul(RotateY(a.beta))
	for _, obj := range a.objList {
		obj.Transform(m)
	}
	a.alpha += a.dAlpha
	if a.alpha > 2*math.Pi {
		a.alpha -= 2 * math.Pi
	}
	a.beta += a.dBeta
	if a.beta > 2*math.Pi {
		a.beta -= 2*math.Pi
	}
	//a.zoom += a.dZoom
	//if (a.zoom < a.minZoom) || (a.zoom > a.maxZoom) {
	//	a.dZoom = -a.dZoom
	//	a.zoom += a.dZoom
	//}
}

func (a *PlatonicAnim) Draw(gc *gg.Context) {
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
		obj.Draw(gc)
	}
}

//----------------------------------------------------------------------------

type Object3D interface {
	Transform(m Matrix)
	Draw(gc *gg.Context)
	Z() float64
}

//----------------------------------------------------------------------------

var (
	// Groessen, Punkte und Flaechendefinitionen fuer....
	// Tetraeder
	f0 = math.Sqrt(3.0/8.0)
	a0 = (1.0/f0) * (1.0/2.0)
	b0 = (1.0/f0) * (math.Sqrt(2.0)/4.0)

	TetraederPoints = []Vector{
		Vector{a0, 0, -b0},
		Vector{-a0, 0, -b0},
		Vector{0, a0, b0},
		Vector{0, -a0, b0},
	}
	TetraederFaces = [][]int{
		{0, 1, 2},
		{0, 3, 1},
		{0, 2, 3},
		{1, 3, 2},
	}

	// Hexaeder
	f1 = math.Sqrt(3.0)/2.0
	a1 = (1.0/f1) * (1.0/2.0)

	HexaederPoints = []Vector{
		Vector{-a1, -a1, -a1},
		Vector{-a1, a1, -a1},
		Vector{a1, a1, -a1},
		Vector{a1, -a1, -a1},
		Vector{-a1, -a1, a1},
		Vector{-a1, a1, a1},
		Vector{a1, a1, a1},
		Vector{a1, -a1, a1},
	}
	HexaederFaces = [][]int{
		{0, 1, 2, 3},
		{0, 4, 5, 1},
		{1, 5, 6, 2},
		{2, 6, 7, 3},
		{3, 7, 4, 0},
		{4, 7, 6, 5},
	}

	// Oktaeder
	f2 = math.Sqrt(2.0)/2.0
	a2 = (1.0/f2) * (1.0/2.0)
	e2 = (1.0/f2) * (math.Sqrt(2.0)/2.0)

	OktaederPoints = []Vector{
		Vector{0, 0, -e2},
		Vector{a2, -a2, 0},
		Vector{a2, a2, 0},
		Vector{-a2, a2, 0},
		Vector{-a2, -a2, 0},
		Vector{0, 0, e2},
	}
	OktaederFaces = [][]int{
		{0, 2, 1},
		{0, 3, 2},
		{0, 4, 3},
		{0, 1, 4},
		{1, 2, 5},
		{2, 3, 5},
		{3, 4, 5},
		{4, 1, 5},
	}
		
	// Ikosaeder
	phi = (math.Sqrt(5)+1)/2
	d = 1.0 / phi
	c = 1.0

	f3 = 0.25 * math.Sqrt(10.0 + 2.0*math.Sqrt(5.0))
	a3 = (1.0/f3) * (1.0/2.0)
	c3 = (1.0/f3) * (1.0+math.Sqrt(5.0))/4.0

	IkosaederPoints = []Vector{
		Vector{0.0, -a3, -c3},
		Vector{0.0, a3, -c3},
		Vector{-c3, 0.0, -a3},
		Vector{c3, 0.0, -a3},

		Vector{-a3, -c3, 0.0},
		Vector{a3, -c3, 0.0},
		Vector{a3, c3, 0.0},
		Vector{-a3, c3, 0.0},

		Vector{-c3, 0.0, a3},
		Vector{c3, 0.0, a3},
		Vector{0.0, -a3, c3},
		Vector{0.0, a3, c3},
	}

	IkosaederFaces = [][]int{
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
	}

	// Dodekaeder
	f4 = 0.25 * (math.Sqrt(5.0) + 1.0) * math.Sqrt(3.0)
	a4 = (1.0/f4) * 1.0/2.0
	b4 = (1.0/f4) * (1.0+phi)/2.0
	c4 = (1.0/f4) * phi/2.0

	DodekaederPoints = []Vector{
		Vector{a4, 0, -b4},
		Vector{-a4, 0, -b4},
		Vector{c4, c4, -c4},
		Vector{-c4, c4, -c4},
		Vector{-c4, -c4, -c4},
		Vector{c4, -c4, -c4},
		Vector{0, b4, -a4},
		Vector{0, -b4, -a4},
		Vector{b4, -a4, 0},
		Vector{b4, a4, 0},
		Vector{-b4, a4, 0},
		Vector{-b4, -a4, 0},
		Vector{0, b4, a4},
		Vector{0, -b4, a4},
		Vector{c4, c4, c4},
		Vector{-c4, c4, c4},
		Vector{-c4, -c4, c4},
		Vector{c4, -c4, c4},
		Vector{a4, 0, b4},
		Vector{-a4, 0, b4},
	}
	DodekaederFaces = [][]int{
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
	}

	// Schweizerkreuz-Box (wird noch nicht korrekt gezeichnet, da
	f5 = 1.0/32.0
	a5 = f5 * 16.0
	b5 = f5 * 10.0
	c5 = f5 *  3.0
	d5 = f5 *  5.0

	SwissPoints = []Vector{
		// Front and Back
		Vector{a5, a5, -d5},
		Vector{-a5, a5, -d5},
		Vector{-a5, -a5, -d5},
		Vector{a5, -a5, -d5},

		// Sides
		Vector{a5, a5, d5},
		Vector{-a5, a5, d5},
		Vector{-a5, -a5, d5},
		Vector{a5, -a5, d5},

		Vector{-b5, -c5, -d5},
		Vector{-b5,  c5, -d5},
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
		Vector{-b5,  c5, d5},
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
	}
	SwissFaces = [][]int{
		{3, 2, 1, 0},
		{4, 5, 6, 7},

		{0, 1, 5, 4},
		{1, 2, 6, 5},
		{2, 3, 7, 6},
		{3, 0, 4, 7},

		// Coords of the 
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
	}
	SwissHoles = [][]int{
		{19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8},
		{20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
	}
)

//-----------------------------------------------------------------------

type PlatonicSolid struct {
	Pts, PtsT  []Vector
	Faces      []*Face
	ZMin, ZMax float64
	LineWidth  float64
	LineColor  colors.RGBA
	FillColor  colors.RGBA
}

func NewPlatonicSolid(points []Vector, faces [][]int, pos Vector,
		size float64, color colors.RGBA, lineWidth float64) *PlatonicSolid {
	o := &PlatonicSolid{}
	o.LineWidth = lineWidth
	o.LineColor = color
	o.FillColor = color.Dark(0.1)
	o.Pts = make([]Vector, len(points))
	o.PtsT = make([]Vector, len(points))
	o.Faces = make([]*Face, len(faces))

	m0 := Scale(Vector{size, size, size})
	m1 := Translate(pos)
	m := m1.Mul(m0)
	for i, pt := range points {
		o.Pts[i] = m.Transform(pt)
	}
	for i, face := range faces {
		o.Faces[i] = NewFace(o, face, color, lineWidth)
	}
	return o
}

func (o *PlatonicSolid) SetHoles(idx [][]int) {
    for i, holes := range idx {
		o.Faces[i].Hole = holes
	}
}

func (o *PlatonicSolid) FaceList() []*Face {
	return o.Faces
}

func (o *PlatonicSolid) Transform(m Matrix) {
	o.ZMin = +math.MaxFloat64
	o.ZMax = -math.MaxFloat64
	for i, pt := range o.Pts {
		ptT := m.Transform(pt)
		if ptT.Z < o.ZMin {
			o.ZMin = ptT.Z
		}
		if ptT.Z > o.ZMax {
			o.ZMax = ptT.Z
		}
		o.PtsT[i] = ptT
	}
}

func (o *PlatonicSolid) Z() float64 {
	return o.ZMax
}

func (o *PlatonicSolid) Draw(gc *gg.Context) {
	return
}

//-----------------------------------------------------------------------

type Face struct {
	o *PlatonicSolid
	Idx     []int
	Hole    []int
	Norm, NormT Vector
	ZMax, ZMin float64
	LineWidth  float64
	LineColor  colors.RGBA
	FillColor  colors.RGBA
}

func NewFace(o *PlatonicSolid, idx []int, color colors.RGBA,
		lineWidth float64) *Face {
	f := &Face{}

	f.o = o
	f.Idx = idx
	f.LineWidth = lineWidth
	f.LineColor = color
	f.FillColor = color.Dark(0.1)

	v1 := o.Pts[idx[2]].Sub(o.Pts[idx[1]])
	v2 := o.Pts[idx[0]].Sub(o.Pts[idx[1]])
	n := v1.Cross(v2)
	f.Norm = n.Div(n.Abs())

	return f
}

func (f *Face) Transform(m Matrix) {
	f.NormT = m.Transform(f.Norm)
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
	gc.SetFillColor(f.FillColor.Dark((1.0 - f.NormT.Z)/2.0))
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

type Dot struct {
	Pt, PtT    Vector
	ZMin, ZMax float64
	Color      colors.RGBA
	Size       float64
}

func NewDot(x, y, z float64, zMin, zMax float64, color colors.RGBA,
	size float64) *Dot {
	d := &Dot{}
	d.Pt = Vector{x, y, z}
	d.ZMin, d.ZMax = zMin, zMax
	d.Color = color
	d.Size = size
	return d
}

func (d *Dot) Transform(m Matrix) {
	d.PtT = m.Transform(d.Pt)
}

func (d *Dot) Z() float64 {
	return d.PtT.Z
}

func (d *Dot) Draw(gc *gg.Context) {
	dark := Map(d.PtT.Z, d.ZMin, d.ZMax, 0.7, 0.2)
	gc.SetFillColor(d.Color.Dark(dark))
	gc.DrawPoint(d.PtT.X, d.PtT.Y, d.Size)
	gc.Fill()
}

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
	s.PtsT[0] = m.Transform(s.Pts[0])
	s.PtsT[1] = m.Transform(s.Pts[1])
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

func NewCloud(x, y, z, w float64, numObjs int, color colors.RGBA,
	size float64) []Object3D {

	dots := make([]Object3D, 0)
	for range numObjs {
		px := rand.NormFloat64()*w + x
		py := rand.NormFloat64()*w + y
		pz := rand.NormFloat64()*w + z
		dots = append(dots, NewDot(px, py, pz, -w, w, color, size))
	}
	return dots
}

//-----------------------------------------------------------------------

type Mappable interface {
	~int | ~int16 | ~float64
}

func Map[I, O Mappable](valIn, lbIn, ubIn I, lbOut, ubOut O) (valOut O) {
	return lbOut + O(valIn-lbIn)*(ubOut-lbOut)/O(ubIn-lbIn)
}

