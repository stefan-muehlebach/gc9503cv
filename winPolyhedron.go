package main

import (
	"log"

	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"

	"math"
	"slices"
	"time"
)

var (
	Ex = NewVector(1, 0, 0)
	Ey = NewVector(0, 1, 0)
	Ez = NewVector(0, 0, 1)

	LightDir = NewVector(1, -1, 1).Normalize()
)

type PolyhedronWin struct {
	windowEmbed
	m0                         Matrix
	objList                    []Object3D
	alpha, dAlpha, beta, dBeta float64
	zoom                       float64
	isManual                   bool
}

func NewPolyhedronWindow(bounds geom.Rectangle[int]) *PolyhedronWin {
	var obj *Polyhedron

	a := &PolyhedronWin{}

	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	rect := a.bounds.Sub(a.bounds.Min).ToFloat()
	drawStyle := Faces

	a.Root = NewGroup()
	a.Root.SetSize(bounds.ToFloat().Size())

	mp := rect.Center()

	a.objList = make([]Object3D, 0)

	objectList := [][]struct {
		data  ObjectData
		color colors.RGBA
	}{
		{
			{HollowCube, colors.DarkGreen},
			{TestBody, colors.DarkKhaki},
			//{HiddenFaces2, colors.DarkOrange},
		},
		{
			{Tetraeder, colors.GoGopherBlue},
			{Hexaeder, colors.GoFuchsia},
			{Oktaeder, colors.GoTeal},
			{Dodekaeder, colors.GoYellow},
			{Ikosaeder, colors.GoAqua},
		},
		{
			{SwissFlag, colors.Red},
			{SwissCross, colors.White},
		},
	}

	objScale := 50.0
	objStep := 170.0
	colLen := float64(len(objectList)-1) * objStep
	for i, list := range objectList {
		listLen := float64(len(list)-1) * objStep
		posX := Ex.Mul(-colLen/2.0 + float64(i)*objStep)
		for j, dat := range list {
			posY := Ey.Mul(listLen/2.0 - float64(j)*objStep)
			obj = NewPolyhedron(dat.data, dat.color, 4.0)
			obj.Scale(Vector{objScale, objScale, objScale})
			obj.Translate(posX.Add(posY))
			obj.DrawStyle = drawStyle
			a.objList = append(a.objList, obj)
		}
	}

	a.m0 = Identity().
		Mul(Translate(NewVector(mp.X, mp.Y, 0.0))).
		Mul(Scale(NewVector(1.0, -1.0, 1.0)))

	a.dAlpha = math.Pi / 300.0
	a.dBeta = math.Pi / 225.0
	a.zoom = 1.5
	a.isManual = true

	var lastPos Point

	a.Root.SetOnPress(func(ev InputEvent) {
		switch ev.Button {
		case LeftButton:
			a.isManual = true
			lastPos = ev.Pos
		}
	})
	a.Root.SetOnClick(func(ev InputEvent) {
		switch ev.Button {
		case RightButton:
			a.isManual = false
		}
	})
	a.Root.SetOnDrag(func(ev InputEvent) {
		switch ev.Button {
		case LeftButton:
			dPos := ev.Pos.Sub(lastPos)
			a.beta -= math.Pi * (dPos.X / 300.0)
			a.alpha += math.Pi * (dPos.Y / 300.0)
			lastPos = ev.Pos
		}
	})
	a.Root.SetOnWheel(func(ev InputEvent) {
		if ev.WheelRel < 0 {
			a.zoom = a.zoom * 0.9
		} else {
			a.zoom = a.zoom / 0.9
		}
	})

	return a
}

func (a *PolyhedronWin) Init() {
	a.alpha = 0.0 // math.Pi / 12.0
	a.beta = 0.0  // math.Pi / 18.0
}

func (a *PolyhedronWin) Update(t time.Time) {
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

func (a *PolyhedronWin) Refresh() {
	a.gc.Clear(colors.Black)
	slices.SortFunc(a.objList, func(a, b Object3D) int {
		_, aMax := a.Z()
		_, bMax := b.Z()
		if aMax < bMax {
			return -1
		} else if aMax > bMax {
			return +1
		} else {
			return 0
		}
	})
	for _, obj := range a.objList {
		obj.Draw(a.gc)
	}
}

//-----------------------------------------------------------------------

type DrawStyleType int

const (
	Faces DrawStyleType = iota
	Edges
	Vertexes
)

//----------------------------------------------------------------------------

type Object3D interface {
	Transform(m Matrix)
	Draw(gc *gg.Context)
	Z() (float64, float64)
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

//-----------------------------------------------------------------------

type Polyhedron struct {
	node3DEmbed
	Pts, PtsT []Vector
	Faces     []*Face
	Shapes    []Object3D

	ZMin, ZMax float64
	LineWidth  float64
	LineColor  colors.RGBA
	FillColor  colors.RGBA
	DrawStyle  DrawStyleType
}

func NewPolyhedron(data ObjectData, color colors.RGBA, lineWidth float64) *Polyhedron {
	o := &Polyhedron{}
	o.LineWidth = lineWidth
	o.LineColor = color.Bright(0.3)
	o.FillColor = color
	o.Pts = make([]Vector, len(data.Points))
	o.PtsT = make([]Vector, len(data.Points))
	o.Faces = make([]*Face, len(data.TrigFaces)+len(data.SimpleFaces))
	o.Shapes = make([]Object3D, 0)

	o.m = Identity()
	//o.Scale(Vector{size, size, size})
	//o.Translate(pos)
	for i, pt := range data.Points {
		o.Pts[i] = pt
	}
	idx := 0
	if len(data.SimpleFaces) > 0 {
		for _, convFace := range data.SimpleFaces {
			o.Faces[idx] = NewSimpleFace(o, convFace)
			idx++
		}
	}
	if len(data.TrigFaces) > 0 {
		for _, faceData := range data.TrigFaces {
			o.Faces[idx] = NewTrigFace(o, faceData)
			idx++
		}
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

func (o *Polyhedron) Draw(gc *gg.Context) {
	switch o.DrawStyle {
	case Faces:
		slices.SortFunc(o.Shapes, func(a, b Object3D) int {
			_, aMax := a.Z()
			_, bMax := b.Z()
			if aMax < bMax {
				return -1
			} else if aMax > bMax {
				return +1
			} else {
				return 0
			}
		})
		for _, shape := range o.Shapes {
			switch obj := shape.(type) {
			case *Trig:
				obj.Draw(gc)
			}
		}
	case Edges:
		slices.SortFunc(o.Shapes, func(a, b Object3D) int {
			_, aMax := a.Z()
			_, bMax := b.Z()
			if aMax < bMax {
				return -1
			} else if aMax > bMax {
				return +1
			} else {
				return 0
			}
		})
		for _, shape := range o.Shapes {
			switch obj := shape.(type) {
			case *Segment:
				obj.Draw(gc)
			}
		}
	case Vertexes:
		gc.SetFillColor(o.LineColor)
		for _, pt := range o.PtsT {
			gc.DrawPoint(pt.X, pt.Y, 2.0)
			gc.Fill()
		}
	}
}

func (o *Polyhedron) Z() (float64, float64) {
	return o.ZMin, o.ZMax
}

//-----------------------------------------------------------------------

type Face struct {
	o           *Polyhedron
	Norm, NormT Vector
	ZMin, ZMax  float64
	LineWidth   float64
	LineColor   colors.RGBA
	FillColor   colors.RGBA
}

// Es wird aktuell nicht geprueft, ob die einzelnen Dreiecke alle in
// einer Ebene liegen und nach der gleichen Seite ausgerichtet sind.
func NewTrigFace(o *Polyhedron, face FaceData) *Face {
	f := &Face{}
	f.o = o

	// Berechne und hinterlege die Flaechennormale aufgrund der ersten
	// Dreiecksflaeche
	an := o.Pts[face.Trigs[0][1]].Sub(o.Pts[face.Trigs[0][0]])
	bn := o.Pts[face.Trigs[0][2]].Sub(o.Pts[face.Trigs[0][0]])
	f.Norm = an.Cross(bn).Normalize()

	// Erstelle Linien-Segmente (gem. 'Edges')
	segments := []*Segment{}
	for _, polyline := range face.Edges {
		for i := 1; i < len(polyline); i++ {
			idxA, idxB := polyline[i-1], polyline[i]
			if idxA > idxB {
				idxA, idxB = idxB, idxA
			}
			segments = append(segments, &Segment{
				Face: f, Idxs: [2]int{idxA, idxB},
			})
		}
	}

	for _, trigIdxs := range face.Trigs {
		trig := &Trig{Face: f, Idxs: [3]int(trigIdxs)}
		trig.findEdges(segments)
		o.Shapes = append(o.Shapes, trig)
	}
	f.LineWidth = o.LineWidth
	f.LineColor = o.LineColor
	f.FillColor = o.FillColor

	return f
}

// Da die Fliesskomma-Rechnerei mit Fehlern behaftet ist und in der
// Geometrie solche Fehler sofort zur Spassbremse werden, hält die
// Konstante 'Epsilon' fest, was wir in unserer Umgebung als 'die
// betragsmässig kleinste Zahl' betrachten. Alle Zahlen,
// deren Betraege kleiner sind als Epsilon, sind gleich Null.
const (
	Epsilon = 1e-6
)

func NewSimpleFace(o *Polyhedron, face []int) *Face {
	if GlobDebug {
		log.Printf("New simple face (face: %+v)", face)
	}
	f := &Face{}

	f.o = o
	//for i := 0; i < len(face); i++ {
	//	idxA := face[i]
	//	idxB := face[(i+1)%len(face)]
	//	o.Shapes = append(o.Shapes, &Segment{Face: f, Idxs: [2]int{idxA, idxB}})
	//}

	// Take the segment before and after point with index 0 on 'face'
	// to construct the normal vector on this face.
	an := o.Pts[face[1]].Sub(o.Pts[face[0]])
	bn := o.Pts[face[len(face)-1]].Sub(o.Pts[face[0]])
	f.Norm = an.Cross(bn).Normalize()

	// Check if all points of a face are in the same plane.
	for l := 1; l < len(face)-1; l++ {
		c := o.Pts[face[l]].Sub(o.Pts[face[0]])
		if prod := an.TripleProd(bn, c); prod > Epsilon {
			if GlobDebug {
				log.Printf("  point %d (%d) not in plane with the others (%f)",
					l, face[l], prod)
			}
			return nil
		}
	}

	i, j, k := 0, 1, 2
	for k != i && j != i {
		if GlobDebug {
			log.Printf("  i,j,k: %d, %d, %d, (%d, %d, %d)", i, j, k,
				face[i], face[j], face[k])
		}
		if !f.IsConvex(o.Pts[face[i]], o.Pts[face[j]], o.Pts[face[k]]) {
			if GlobDebug {
				log.Printf("  triple product is negativ at point %d (%d)",
					j, face[j])
			}
			i, j, k = j, k, nextPoint(face, k)
			continue
		}

		trig := &Trig{Face: f, Idxs: [3]int{face[i], face[j], face[k]}}
		if (i+1)%len(face) == j {
			trig.Edges = append(trig.Edges, &Segment{
				Face: f, Idxs: [2]int{face[i], face[j]},
			})
		}
		if (j+1)%len(face) == k {
			trig.Edges = append(trig.Edges, &Segment{
				Face: f, Idxs: [2]int{face[j], face[k]},
			})
		}
		if (k+1)%len(face) == i {
			trig.Edges = append(trig.Edges, &Segment{
				Face: f, Idxs: [2]int{face[k], face[i]},
			})
		}
		o.Shapes = append(o.Shapes, trig)

		face[j] = -1
		j, k = k, nextPoint(face, k)
		if GlobDebug {
			log.Printf("  face: %+v", face)
		}
	}

	f.LineWidth = o.LineWidth
	f.LineColor = o.LineColor
	f.FillColor = o.FillColor

	return f
}

func nextPoint(face []int, last int) int {
	for idx := (last + 1) % len(face); idx != last; idx = (idx + 1) % len(face) {
		if face[idx] >= 0 {
			return idx
		}
	}
	return -1
}

func (f *Face) IsVisible() bool {
	return f.NormT.Z > 0
}

func (f *Face) IsConvex(a, b, c Vector) bool {
	x := b.Sub(a)
	y := c.Sub(b)
	prod := x.TripleProd(y, f.Norm)
	return prod > 0.0
}

func (f *Face) Transform(m Matrix) {
	f.NormT = m.TransformLinear(f.Norm)
}

func (f *Face) Z() (float64, float64) {
	return f.ZMin, f.ZMax
}

func (f *Face) Draw(gc *gg.Context) {}

//-----------------------------------------------------------------------

// Damit wird ein Dreieck (Triangle) als Runtime-Datenstrukut beschrieben
type Trig struct {
	// Pointer auf die Flaeche, bei welcher dieses Dreieck angesiedelt ist.
	Face *Face
	// Ein Dreieck kann max. an allen drei Seiten eine Flaechenbegrenzung
	// haben (genau dann, wenn die Flaeche nur durch das Dreieck bestimmt
	// ist) und min. an keiner Seite (genau dann, wenn das Dreieck an allen
	// drei Seiten von weiteren Dreiecken dieser Flaeche umschlossen wird).
	// Regel: der Slice 'Edges' wird mit aufsteigendem Index durchlaufen;
	// beim ersten 'nil'-Eintrag, wird die Verarbeitung abgebrochen.
	Edges []*Segment
	// Dies schliesslich sind die Indizes der Punkte, welche dieses Dreieck
	// ausmachen. Die Zahlen beziehen sich auf den Slice 'Pts', resp. 'PtsT'
	// im darueberliegenden 'Polyhedron'-Objekt.
	Idxs TrigData
}

func (t *Trig) Transform(m Matrix) {
	for _, edge := range t.Edges {
		edge.finalDraw = false
	}
}

func (t *Trig) Draw(gc *gg.Context) {
	if !t.Face.IsVisible() {
		return
	}
	ptsT := t.Face.o.PtsT
	gc.SetLineWidth(1.0)
	//col := t.Face.FillColor
	col := t.Face.FillColor.Dark((1.0 - t.Face.NormT.Z/t.Face.NormT.Abs()) / 2.0)
	if GlobDebug {
		gc.SetLineColor(colors.Crimson)
	} else {
		gc.SetLineColor(col)
	}
	gc.SetFillColor(col)
	gc.MoveTo(ptsT[t.Idxs[0]].X, ptsT[t.Idxs[0]].Y)
	gc.LineTo(ptsT[t.Idxs[1]].X, ptsT[t.Idxs[1]].Y)
	gc.LineTo(ptsT[t.Idxs[2]].X, ptsT[t.Idxs[2]].Y)
	gc.ClosePath()
	gc.FillStroke()

	for _, edge := range t.Edges {
		edge.Draw(gc)
	}
}

func (t *Trig) Z() (float64, float64) {
	ptsT := t.Face.o.PtsT
	zMin := min(ptsT[t.Idxs[0]].Z, min(ptsT[t.Idxs[1]].Z,
		ptsT[t.Idxs[2]].Z))
	zMax := max(ptsT[t.Idxs[0]].Z, max(ptsT[t.Idxs[1]].Z,
		ptsT[t.Idxs[2]].Z))
	return zMin, zMax
}

func (t *Trig) findEdges(segments []*Segment) {
	pairs := [][2]int{
		{t.Idxs[0], t.Idxs[1]},
		{t.Idxs[1], t.Idxs[2]},
		{t.Idxs[0], t.Idxs[2]},
	}
	for _, pair := range pairs {
		idxA, idxB := pair[0], pair[1]
		if idxA > idxB {
			idxA, idxB = idxB, idxA
		}
		for _, segment := range segments {
			if idxA == segment.Idxs[0] && idxB == segment.Idxs[1] {
				t.Edges = append(t.Edges, segment)
			}
		}
	}
}

//-----------------------------------------------------------------------

type Segment struct {
	Face      *Face
	Idxs      [2]int
	finalDraw bool
}

func (s *Segment) Transform(m Matrix) {}

func (s *Segment) Draw(gc *gg.Context) {
	if !s.finalDraw {
		s.finalDraw = true
		return
	}
	ptsT := s.Face.o.PtsT
	gc.SetLineWidth(s.Face.LineWidth)
	col := s.Face.LineColor
	// col := s.Face.LineColor.Dark((1.0 -
	//     s.Face.NormT.Z/s.Face.NormT.Abs()) / 4.0)
	gc.SetLineColor(col)
	gc.DrawLine(ptsT[s.Idxs[0]].X, ptsT[s.Idxs[0]].Y,
		ptsT[s.Idxs[1]].X, ptsT[s.Idxs[1]].Y)
	gc.Stroke()
}

func (s *Segment) Z() (float64, float64) {
	ptsT := s.Face.o.PtsT
	zMin := min(ptsT[s.Idxs[0]].Z, ptsT[s.Idxs[1]].Z)
	zMax := max(ptsT[s.Idxs[0]].Z, ptsT[s.Idxs[1]].Z)
	return zMin, zMax
}

//-----------------------------------------------------------------------

type Mappable interface {
	~int | ~int16 | ~float64
}

func Map[I, O Mappable](valIn, lbIn, ubIn I, lbOut, ubOut O) (valOut O) {
	return lbOut + O(valIn-lbIn)*(ubOut-lbOut)/O(ubIn-lbIn)
}
