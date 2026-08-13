package geom

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

// Über diese Konstante kann der Vergleich von Fliesskommazahlen beeinflusst
// werden. Die Zahlen im Format IEEE 754 zeichnen sich insbesondere in
// Zusammenhang mit Geometrie durch eine nervige 'Unschärfe' aus, d.h. zwei
// Zahlen sollten gem. Konstruktion gleich sein (d.h. a == b ist true) - aber
// sie liegen oft um ein paar wenige 1.0e-10's auseinader.
const (
	eps           = 1.0e-12
	numObjs       = 10_001
	Width, Height = 480.0, 960.0
)

// type ct float64
// type Point[float64] Point[float64]
// type Rectangle[float64] Rectangle[float64]

var (
	p0, p1, p2, p3, p4, p5, p6, p7, p8, p9 Point[float64]
	r0, r1, r2, r3, r4, r5, r6, r7, r8, r9 Rectangle[float64]
	d, t, tx, ty                           float64
	x, y, z                                float64
	points                                 []Point[float64]
	rectangles                             []Rectangle[float64]
	isIn                                   bool

	rectList = []Rectangle[float64]{
		Rect(0.0, 0.0, 10.0, 10.0),
		Rect(10.0, 0, 20, 10),
		Rect(1.0, 2, 3, 4),
		Rect(4.0, 6, 10, 10),
		Rect(2.0, 3, 12, 5),
		Rect(-1.0, -2, 0, 0),
		Rect(-1.0, -2, 4, 6),
		Rect(-10.0, -20, 30, 40),
		Rect(8.0, 8, 8, 8),
		Rect(88.0, 88, 88, 88),
		Rect(6.0, 5, 4, 3),
		{Point[float64]{-1, -1}, Point[float64]{1, 1}},
		{Point[float64]{-1, 1}, Point[float64]{1, -1}},
		{Point[float64]{1, 1}, Point[float64]{-1, -1}},
		{Point[float64]{1, -1}, Point[float64]{-1, 1}},
		{Point[float64]{-1, 0}, Point[float64]{1, 0}},
		{Point[float64]{1, 0}, Point[float64]{-1, 0}},
		{Point[float64]{0, -1}, Point[float64]{0, 1}},
		{Point[float64]{0, 1}, Point[float64]{0, -1}},
	}
)

func init() {

}

func eq(a, b float64) bool {
	if math.Abs(a-b) < eps {
		return true
	} else {
		return false
	}
}

//----------------------------------------------------------------------------
//
// Tests for the Point datatype.
//

func TestPoint(t *testing.T) {
	var p Point[float64]
	var d float64
	var b bool

	t.Run("Add", func(t *testing.T) {
		testVect := [][]Point[float64]{
			{Pt(2.0, -2.0), Pt(-4.0, 4.0), Pt(-2.0, 2.0)},
		}
		for _, rec := range testVect {
			p = rec[0].Add(rec[1])
			if p != rec[2] {
				t.Errorf("%v.Add(%v); got %v, want %v", rec[0], rec[1], p, rec[2])
			}
		}
	})

	t.Run("Sub", func(t *testing.T) {
		testVect := [][]Point[float64]{
			{Pt(2.0, -2.0), Pt(-4.0, 4.0), Pt(6.0, -6.0)},
		}
		for _, rec := range testVect {
			p = rec[0].Sub(rec[1])
			if p != rec[2] {
				t.Errorf("%v.Sub(%v); got %v, want %v", rec[0], rec[1], p, rec[2])
			}
		}
	})

	t.Run("Mul", func(t *testing.T) {
		testVect := []struct {
			p Point[float64]
			t float64
			q Point[float64]
		}{
			{Pt(2.0, -2.0), 2.0, Pt(4.0, -4.0)},
			{Pt(0.5, 2.0), 0.5, Pt(0.25, 1.0)},
		}
		for _, rec := range testVect {
			p = rec.p.Mul(rec.t)
			if p != rec.q {
				t.Errorf("%v.Mul(%v); got %v, want %v", rec.p, rec.t, p, rec.q)
			}
		}
	})

	t.Run("Distance", func(t *testing.T) {
		testVect := []struct {
			p0 Point[float64]
			p1 Point[float64]
			d  float64
		}{
			{Pt(0.0, 0.0), Pt(1.0, 0.0), 1.0},
			{Pt(0.0, 0.0), Pt(-1.0, 0.0), 1.0},
			{Pt(0.0, 0.0), Pt(1.0, 1.0), math.Sqrt(2.0)},
			{Pt(0.0, 0.0), Pt(-1.0, -1.0), math.Sqrt(2.0)},
		}
		for _, rec := range testVect {
			d = rec.p0.Distance(rec.p1)
			if d != rec.d {
				t.Errorf("%v.Distance(%v); got %v, want %v", rec.p0, rec.p1, d, rec.d)
			}
		}
	})

	t.Run("In", func(t *testing.T) {
		testVect := []struct {
			p Point[float64]
			r Rectangle[float64]
			b bool
		}{
			{Pt(0.0, 0.0), Rect(-1.0, -1.0, 1.0, 1.0), true},
			{Pt(-1.0, -1.0), Rect(-1.0, -1.0, 1.0, 1.0), true},
			{Pt(1.0, 1.0), Rect(-1.0, -1.0, 1.0, 1.0), false},
			{Pt(0.0, 1.0), Rect(-1.0, -1.0, 1.0, 1.0), false},
			{Pt(1.0, 0.0), Rect(-1.0, -1.0, 1.0, 1.0), false},
		}
		for _, rec := range testVect {
			b = rec.p.In(rec.r)
			if b != rec.b {
				t.Errorf("%v.In(%v); got %v, want %v", rec.p, rec.r, b, rec.b)
			}
		}
	})
}

type testCase struct {
	str  string         // Die Textdarstellung des Punktes
	ok   bool           // Ist dies ueberhaupt eine korrekte Darstellung
	pt   Point[float64] // Die erwartete Darstellung als Point
	dist float64        // Der erwartete Betrag
}

var (
	testList = []testCase{
		{"(0;0)", true, Point[float64]{0, 0}, 0.0},
		{"(1;0)", true, Point[float64]{1, 0}, 1.0},
		{"(1;1)", true, Point[float64]{1, 1}, math.Sqrt(2)},
		{"(0;1)", true, Point[float64]{0, 1}, 1.0},
		{"(-1;1)", true, Point[float64]{-1, 1}, math.Sqrt(2)},
		{"(-1;0)", true, Point[float64]{-1, 0}, 1.0},
		{"(-1;-1)", true, Point[float64]{-1, -1}, math.Sqrt(2)},
		{"(0;-1)", true, Point[float64]{0, -1}, 1.0},
		{"(1;-1)", true, Point[float64]{1, -1}, math.Sqrt(2)},

		{"(-.5;.5)", true, Point[float64]{-0.5, 0.5}, math.Sqrt(0.5)},

		{"1;1", false, Point[float64]{}, 0.0},
		{"(1,1)", false, Point[float64]{}, 0.0},
		{"1,1)", false, Point[float64]{}, 0.0},
		{"(1,1", false, Point[float64]{}, 0.0},
	}
)

func TestPointSet(t *testing.T) {
	var err error

	for i, tst := range testList {
		err = p0.Set(tst.str)
		if tst.ok && err != nil {
			t.Errorf("%d: parsing %s failed: %v", i, tst.str, err)
			continue
		}
		if !tst.ok {
			if err == nil {
				t.Errorf("%d: parsing %s should fail!", i, tst.str)
			}
			continue
		}
		if !p0.Eq(tst.pt) {
			t.Errorf("%d: %v not equal to %v", i, p0, tst.pt)
			continue
		}
		// d := p0.Abs()
		// if d != tst.dist {
		// 	t.Errorf("%d: distance to (0;0): have %f, want %f", i, d, tst.dist)
		// 	continue
		// }

		if (p0.X == 0.0) && (p0.Y == 0.0) {
			continue
		}
		// pNorm := p0.Normalize()
		// lNorm := pNorm.Abs()
		// if !eq(lNorm, 1.0) {
		// 	t.Errorf("length of normalized vector %v should be 1.0, is %f", pNorm, lNorm)
		// }
	}
}

func BenchmarkPoint(b *testing.B) {
	var id int

	points = make([]Point[float64], numObjs)
	rectangles = make([]Rectangle[float64], numObjs)

	seed := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(seed)
	for i := 0; i < numObjs; i++ {
		points[i] = Point[float64]{Width * rnd.Float64(), Height * rnd.Float64()}
		x0, y0 := Width*rnd.Float64(), Height*rnd.Float64()
		x1, y1 := x0+(Width-x0)*rnd.Float64(), y0+(Height-y0)*rnd.Float64()
		rectangles[i] = Rect(x0, y0, x1, y1)
	}

	id = 0

	b.Run("Add", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			points[id+2] = points[id].Add(points[id+1])
		}
	})
	id += 3

	b.Run("Sub", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			points[id+2] = points[id].Sub(points[id+1])
		}
	})
	id += 3

	t = rnd.Float64()
	b.Run("Mul", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			points[id+1] = points[id].Mul(t)
		}
	})
	id += 2

	t = rnd.Float64()
	b.Run("Div", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			points[id+1] = points[id].Div(t)
		}
	})
	id += 2

	b.Run("Move", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			points[id].Move(points[id+1])
		}
	})
	id += 2

	b.Run("Distance", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			d = points[id].Distance(points[id+1])
		}
	})
	id += 2

	b.Run("Dist2", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			d = points[id].Dist2(points[id+1])
		}
	})
	id += 2

	b.Run("Min", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			points[id+2] = points[id].Min(points[id+1])
		}
	})
	id += 3

	b.Run("Max", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			points[id+2] = points[id].Max(points[id+1])
		}
	})
	id += 3

	b.Run("In", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			isIn = points[id].In(rectangles[id+1])
		}
	})
	id += 2

	// t = rnd.Float64()
	// b.Run("Interpolate", func(b *testing.B) {
	// 	for b.Loop() {
	// 		points[id+2] = points[id].Interpolate(points[id+1], t)
	// 	}
	// })
	// id += 3

}

//----------------------------------------------------------------------------
//
// Tests for the Rectangle datatype.
//

func TestRectangle(t *testing.T) {
	// in checks that every point in f is in g.
	in := func(f, g Rectangle[float64]) error {
		if !f.In(g) {
			return fmt.Errorf("f=%s, f.In(%s): got false, want true", f, g)
		}
		for y := f.Min.Y; y < f.Max.Y; y++ {
			for x := f.Min.X; x < f.Max.X; x++ {
				p := Point[float64]{x, y}
				if !p.In(g) {
					return fmt.Errorf("p=%s, p.In(%s): got false, want true", p, g)
				}
			}
		}
		return nil
	}

	// r.Eq(s) should be equivalent to every point in r being in s, and every
	// point in s being in r.
	for _, r := range rectList {
		for _, s := range rectList {
			got := r.Eq(s)
			want := in(r, s) == nil && in(s, r) == nil
			if got != want {
				t.Errorf("Eq: r=%s, s=%s: got %t, want %t", r, s, got, want)
			}
		}
	}

	// The intersection should be the largest rectangle a such that every point
	// in a is both in r and in s.
	// for _, r := range rectList {
	// 	for _, s := range rectList {
	// 		a := r.Intersect(s)
	// 		if err := in(a, r); err != nil {
	// 			t.Errorf("Intersect: r=%s, s=%s, a=%s, a not in r: %v", r, s, a, err)
	// 		}
	// 		if err := in(a, s); err != nil {
	// 			t.Errorf("Intersect: r=%s, s=%s, a=%s, a not in s: %v", r, s, a, err)
	// 		}
	// 		if isZero, overlaps := a == (Rectangle[float64]{}), r.Overlaps(s); isZero == overlaps {
	// 			t.Errorf("Intersect: r=%s, s=%s, a=%s: isZero=%t same as overlaps=%t",
	// 				r, s, a, isZero, overlaps)
	// 		}
	// 		largerThanA := [4]Rectangle[float64]{a, a, a, a}
	// 		largerThanA[0].Min.X--
	// 		largerThanA[1].Min.Y--
	// 		largerThanA[2].Max.X++
	// 		largerThanA[3].Max.Y++
	// 		for i, b := range largerThanA {
	// 			if b.Empty() {
	// 				// b isn't actually larger than a.
	// 				continue
	// 			}
	// 			if in(b, r) == nil && in(b, s) == nil {
	// 				t.Errorf("Intersect: r=%s, s=%s, a=%s, b=%s, i=%d: intersection could be larger",
	// 					r, s, a, b, i)
	// 			}
	// 		}
	// 	}
	// }

	// The union should be the smallest rectangle a such that every point in r
	// is in a and every point in s is in a.
	for _, r := range rectList {
		for _, s := range rectList {
			a := r.Union(s)
			if err := in(r, a); err != nil {
				t.Errorf("Union: r=%s, s=%s, a=%s, r not in a: %v", r, s, a, err)
			}
			if err := in(s, a); err != nil {
				t.Errorf("Union: r=%s, s=%s, a=%s, s not in a: %v", r, s, a, err)
			}
			if a.Empty() {
				// You can't get any smaller than a.
				continue
			}
			smallerThanA := [4]Rectangle[float64]{a, a, a, a}
			smallerThanA[0].Min.X++
			smallerThanA[1].Min.Y++
			smallerThanA[2].Max.X--
			smallerThanA[3].Max.Y--
			for i, b := range smallerThanA {
				if in(r, b) == nil && in(s, b) == nil {
					t.Errorf("Union: r=%s, s=%s, a=%s, b=%s, i=%d: union could be smaller",
						r, s, a, b, i)
				}
			}
		}
	}
}

func TestRectangleCanon(t *testing.T) {
	for _, rect := range rectList {
		r0 = rect.Canon()
		if (r0.Min.X > r0.Max.X) || (r0.Min.Y > r0.Max.Y) {
			t.Errorf("rectangle %v is not canonical", r0)
		}
	}
}

// func TestRectangleRelPos(t *testing.T) {
// 	r0 = Rect(-3, -2, 4, 3)
// 	relList := [][]float64{
// 		{0.0, 0.0},
// 		{1.0, 0.0},
// 		{0.0, 1.0},
// 		{1.0, 1.0},
// 		{-1.0, -1.0},
// 		{2.0, 2.0},
// 		{0.5, 0.5},
// 		{-0.5, -0.5},
// 	}
// 	posList := []Point[float64]{
// 		{-3.0, -2.0},
// 		{4.0, -2.0},
// 		{-3.0, 3.0},
// 		{4.0, 3.0},
// 		{-10.0, -7.0},
// 		{11.0, 8.0},
// 		{0.5, 0.5},
// 		{-6.5, -4.5},
// 	}

// 	for i, rel := range relList {
// 		p0 = r0.RelPos(rel[0], rel[1])
// 		if (p0.X != posList[i].X) || (p0.Y != posList[i].Y) {
// 			t.Errorf("%v.RelPos(%f, %f) = %v; want %v", r0, rel[0], rel[1], p0, posList[i])
// 		}
// 	}
// 	for i, pos := range posList {
// 		rx, ry := r0.PosRel(pos)
// 		if (rx != relList[i][0]) || (ry != relList[i][1]) {
// 			t.Errorf("%v.PosRel(%v) = %f, %f; want %f, %f", r0, pos, rx, ry, relList[i][0], relList[i][1])
// 		}
// 	}
// }

// func TestRectangleSetInside(t *testing.T) {
// 	r0 = Rect(-3, -2, 4, 3)
// 	posList := [][]Point{
// 		{Point{-4, -3}, Point{-3, -2}},
// 		{Point{-4, -2}, Point{-3, -2}},
// 		{Point{-4, 3}, Point{-3, 3}},
// 		{Point{-4, 4}, Point{-3, 3}},
// 		{Point{-3, -3}, Point{-3, -2}},
// 		{Point{-3, -2}, Point{-3, -2}},
// 		{Point{-3, 3}, Point{-3, 3}},
// 		{Point{-3, 4}, Point{-3, 3}},
// 		{Point{0, -3}, Point{0, -2}},
// 		{Point{0, -2}, Point{0, -2}},
// 		{Point{0, 3}, Point{0, 3}},
// 		{Point{0, 4}, Point{0, 3}},
// 		{Point{4, -3}, Point{4, -2}},
// 		{Point{4, -2}, Point{4, -2}},
// 		{Point{4, 3}, Point{4, 3}},
// 		{Point{4, 4}, Point{4, 3}},
// 	}
// 	for _, pos := range posList {
// 		p0 = r0.SetInside(pos[0])
// 		if !p0.Eq(pos[1]) {
// 			t.Errorf("SetInside of %v failed: got %v want %v", pos[0], p0, pos[1])
// 		}
// 	}
// }

func BenchmarkRectangle(b *testing.B) {
	var id int

	points = make([]Point[float64], numObjs)
	rectangles = make([]Rectangle[float64], numObjs)

	seed := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(seed)
	for i := 0; i < numObjs; i++ {
		points[i] = Pt(Width*rnd.Float64(), Height*rnd.Float64())
		x0, y0 := Width*rnd.Float64(), Height*rnd.Float64()
		x1, y1 := x0+(Width-x0)*rnd.Float64(), y0+(Height-y0)*rnd.Float64()
		rectangles[i] = Rect(x0, y0, x1, y1)
	}

	id = 0

	b.Run("Add", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			rectangles[id+2] = rectangles[id].Add(points[id+1])
		}
	})
	id += 3

	b.Run("Sub", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			rectangles[id+2] = rectangles[id].Sub(points[id+1])
		}
	})
	id += 3

	b.Run("Move", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			rectangles[id].Move(points[id])
		}
	})
	id += 1

	// b.Run("PosRel", func(b *testing.B) {
	// 	for b.Loop() {
	// 		tx, ty = rectangles[id].PosRel(points[id])
	// 	}
	// })
	// id += 1

	// tx, ty = rand.Float64(), rand.Float64()
	// b.Run("RelPos", func(b *testing.B) {
	// 	for b.Loop() {
	// 		points[id] = rectangles[id].RelPos(tx, ty)
	// 	}
	// })
	// id += 1

	// b.Run("Center", func(b *testing.B) {
	// 	for b.Loop() {
	// 		points[id] = rectangles[id].Center()
	// 	}
	// })
	// id += 1

	// b.Run("SetInside", func(b *testing.B) {
	// 	for b.Loop() {
	// 		points[id+1] = rectangles[id].SetInside(points[id])
	// 	}
	// })
	// id += 2

	b.Run("Union", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			rectangles[id+2] = rectangles[id].Union(rectangles[id+1])
		}
	})
	id += 3

	b.Run("Intersect", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			rectangles[id+2] = rectangles[id].Intersect(rectangles[id+1])
		}
	})
	id += 3

	b.Run("In", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			isIn = rectangles[id].In(rectangles[id+1])
		}
	})
	id += 2

}

//----------------------------------------------------------------------------
//
// Tests for the Transform datatype.
//

func TestTransform(t *testing.T) {
	var p1 Point[float64]
	var r1 Rectangle[float64]
	r := Rect(0.0, 0.0, 360.0, 960.0)
	tr := NewTransform(Rotate000, r)

	t.Run("Point", func(t *testing.T) {
		testVect := []struct {
			p0  Point[float64]
			rot RotationType
			dir TransfDirection
			p1  Point[float64]
		}{
			{Pt(30.0, 10.0), Rotate000, Base2Sub, Pt(30.0, 10.0)},
			{Pt(30.0, 10.0), Rotate090, Base2Sub, Pt(950.0, 30.0)},
			{Pt(30.0, 10.0), Rotate180, Base2Sub, Pt(330.0, 950.0)},
			{Pt(30.0, 10.0), Rotate270, Base2Sub, Pt(10.0, 330.0)},

			{Pt(30.0, 10.0), Rotate000, Sub2Base, Pt(30.0, 10.0)},
			{Pt(30.0, 10.0), Rotate090, Sub2Base, Pt(10.0, 930.0)},
			{Pt(30.0, 10.0), Rotate180, Sub2Base, Pt(330.0, 950.0)},
			{Pt(30.0, 10.0), Rotate270, Sub2Base, Pt(350.0, 30.0)},

			{Pt(0.0, 0.0), Rotate180, Sub2Base, Pt(360.0, 960.0)},
			{Pt(-10.0, -10.0), Rotate180, Sub2Base, Pt(370.0, 970.0)},
		}
		for _, rec := range testVect {
			tr.Rot = rec.rot
			p1 = tr.Point(rec.p0, rec.dir)
			if p1 != rec.p1 {
				t.Errorf("%v.Point(%v, %v); got %v, want %v", rec.rot, rec.p0, rec.dir, p1, rec.p1)
			}
		}
	})

	t.Run("Rectangle", func(t *testing.T) {
		testVect := []struct {
			r0  Rectangle[float64]
			rot RotationType
			dir TransfDirection
			r1  Rectangle[float64]
		}{
			{Rect(0.0, 0.0, 360.0, 960.0), Rotate000, Base2Sub, Rect(0.0, 0.0, 360.0, 960.0)},
			{Rect(0.0, 0.0, 360.0, 960.0), Rotate090, Base2Sub, Rect(0.0, 0.0, 960.0, 360.0)},
			{Rect(0.0, 0.0, 360.0, 960.0), Rotate180, Base2Sub, Rect(0.0, 0.0, 360.0, 960.0)},
			{Rect(0.0, 0.0, 360.0, 960.0), Rotate270, Base2Sub, Rect(0.0, 0.0, 960.0, 360.0)},

			{Rect(0.0, 0.0, 360.0, 960.0), Rotate000, Sub2Base, Rect(0.0, 0.0, 360.0, 960.0)},
			{Rect(0.0, 0.0, 960.0, 360.0), Rotate090, Sub2Base, Rect(0.0, 0.0, 360.0, 960.0)},
			{Rect(0.0, 0.0, 360.0, 960.0), Rotate180, Sub2Base, Rect(0.0, 0.0, 360.0, 960.0)},
			{Rect(0.0, 0.0, 960.0, 360.0), Rotate270, Sub2Base, Rect(0.0, 0.0, 360.0, 960.0)},

			{Rect(30.0, 10.0, 40.0, 20.0), Rotate000, Base2Sub, Rect(30.0, 10.0, 40.0, 20.0)},
			{Rect(30.0, 10.0, 40.0, 20.0), Rotate090, Base2Sub, Rect(940.0, 30.0, 950.0, 40.0)},
			{Rect(30.0, 10.0, 40.0, 20.0), Rotate180, Base2Sub, Rect(320.0, 940.0, 330.0, 950.0)},
			{Rect(30.0, 10.0, 40.0, 20.0), Rotate270, Base2Sub, Rect(10.0, 320.0, 20.0, 330.0)},

			{Rect(30.0, 10.0, 40.0, 20.0), Rotate000, Sub2Base, Rect(30.0, 10.0, 40.0, 20.0)},
			{Rect(30.0, 10.0, 40.0, 20.0), Rotate090, Sub2Base, Rect(10.0, 920.0, 20.0, 930.0)},
			{Rect(30.0, 10.0, 40.0, 20.0), Rotate180, Sub2Base, Rect(320.0, 940.0, 330.0, 950.0)},
			{Rect(30.0, 10.0, 40.0, 20.0), Rotate270, Sub2Base, Rect(340.0, 30.0, 350.0, 40.0)},
		}
		for _, rec := range testVect {
			tr.Rot = rec.rot
			r1 = tr.Rectangle(rec.r0, rec.dir)
			if r1 != rec.r1 {
				t.Errorf("%v.Rectangle(%v, %v); got %v, want %v", rec.rot, rec.r0, rec.dir, r1, rec.r1)
			}
		}
	})
}

func BenchmarkTransform(b *testing.B) {
	seed := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(seed)

	r := Rect(0.0, 0.0, Width, Height)
	tr := NewTransform(Rotate090, r)

	p0 = Pt(Width*rnd.Float64(), Height*rnd.Float64())
	x0, y0 := Width*rnd.Float64(), Height*rnd.Float64()
	x1, y1 := x0+(Width-x0)*rnd.Float64(), y0+(Height-y0)*rnd.Float64()
	r0 = Rect(x0, y0, x1, y1)

	b.Run("Point", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			p1 = tr.Point(p0, Base2Sub)
		}
	})

	b.Run("Rectangle", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			r1 = tr.Rectangle(r0, Base2Sub)
		}
	})
}
