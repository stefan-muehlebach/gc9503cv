package geom

import (
	"fmt"
	"image"
	"math"
)

//----------------------------------------------------------------------------

type Size[T CoordType] = Point[T]

//----------------------------------------------------------------------------

type Point[T CoordType] struct {
	X, Y T
}

func Pt[T CoordType](x, y T) Point[T] {
	return Point[T]{x, y}
}

func (p Point[T]) Add(q Point[T]) Point[T] {
	return Point[T]{p.X + q.X, p.Y + q.Y}
}

func (p Point[T]) AddXY(x, y T) Point[T] {
	return Point[T]{p.X + x, p.Y + y}
}

func (p Point[T]) Sub(q Point[T]) Point[T] {
	return Point[T]{p.X - q.X, p.Y - q.Y}
}

func (p Point[T]) SubXY(x, y T) Point[T] {
	return Point[T]{p.X - x, p.Y - y}
}

func (p Point[T]) Mul(t T) Point[T] {
	return Point[T]{t * p.X, t * p.Y}
}

func (p Point[T]) Scale(s Point[T]) Point[T] {
	return Point[T]{p.X * s.X, p.Y * s.Y}
}

func (p Point[T]) Div(t T) Point[T] {
	return Point[T]{p.X / t, p.Y / t}
}

func (p *Point[T]) Move(dp Point[T]) {
	p.X += dp.X
	p.Y += dp.Y
}

func (p *Point[T]) Neg() Point[T] {
	return Point[T]{-p.X, -p.Y}
}

func (p Point[T]) Eq(q Point[T]) bool {
	return p.X == q.X && p.Y == q.Y
}

func (p Point[T]) In(r Rectangle[T]) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

func (p Point[T]) Distance(q Point[T]) float64 {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return math.Hypot(float64(dx), float64(dy))
}

func (p Point[T]) Dist2(q Point[T]) T {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return dx*dx + dy*dy
}

func (p Point[T]) Abs() float64 {
	return math.Hypot(float64(p.X), float64(p.Y))
}

func (p Point[T]) Interpolate(q Point[T], t float64) Point[T] {
	u := 1.0 - t

	x := float64(p.X)*u + float64(q.X)*t
	y := float64(p.Y)*u + float64(q.Y)*t
	return Point[T]{T(x), T(y)}
}

// Vergleicht die X- sowie die Y-Werte der Punkte p und q und retourniert einen
// neuen Punkt mit den jeweils kleineren Werten.
func (p Point[T]) Min(q Point[T]) Point[T] {
	return Point[T]{min(p.X, q.X), min(p.Y, q.Y)}
}

// Vergleicht die X- sowie die Y-Werte der Punkte p und q und retourniert einen
// neuen Punkt mit den jeweils grösseren Werten.
func (p Point[T]) Max(q Point[T]) Point[T] {
	return Point[T]{max(p.X, q.X), max(p.Y, q.Y)}
}

func (p Point[T]) AsCoord() (x, y T) {
	return p.X, p.Y
}

func (p Point[T]) ToFloat() Point[float64] {
	return Point[float64]{float64(p.X), float64(p.Y)}
}

func (p Point[T]) ToInt() image.Point {
	return image.Point{int(p.X), int(p.Y)}
}

// Gibt die Koordinaten des Punktes in der Form '(x; y)' zurück. Implementiert
// das Stringer-Interface.
func (p Point[T]) String() string {
	return fmt.Sprintf("(%v; %v)", p.X, p.Y)
}

// Damit können Punkte (resp. die Koordinaten dazu) auch über ein Textfile
// oder die Kommandozeile eingelesen werden. Mit String zusammen implementiert
// Point somit das Getter-Interface aus flag.
func (p *Point[T]) Set(s string) error {
	var x, y T

	_, err := fmt.Sscanf(s, "(%f;%f)", &x, &y)
	p.X, p.Y = x, y
	return err
}
