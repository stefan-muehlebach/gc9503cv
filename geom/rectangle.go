package geom

import (
	"image"
)

//----------------------------------------------------------------------------

type Rectangle[T CoordType] struct {
	Min, Max Point[T]
}

func Rect[T CoordType](x0, y0, x1, y1 T) Rectangle[T] {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle[T]{
		Pt(x0, y0),
		Pt(x1, y1),
	}
}

// Kanonisiert das Rechteck r. Das heisst, dass im resultierenden Rechteck
// die Koordinaten des Punktes Min auf jeden Fall kleiner sind als die Koord.
// des Punktes Max.
func (r Rectangle[T]) Canon() Rectangle[T] {
	return Rectangle[T]{
		Min: r.Min.Min(r.Max),
		Max: r.Min.Max(r.Max),
	}
}

func (r Rectangle[T]) Add(p Point[T]) Rectangle[T] {
	return Rectangle[T]{
		r.Min.Add(p),
		r.Max.Add(p),
	}
}

func (r Rectangle[T]) Sub(p Point[T]) Rectangle[T] {
	return Rectangle[T]{
		r.Min.Sub(p),
		r.Max.Sub(p),
	}
}

func (r Rectangle[T]) AddXY(x, y T) Rectangle[T] {
	return Rectangle[T]{
		r.Min.AddXY(x, y),
		r.Max.AddXY(x, y),
	}
}

func (r Rectangle[T]) SubXY(x, y T) Rectangle[T] {
	return Rectangle[T]{
		r.Min.SubXY(x, y),
		r.Max.SubXY(x, y),
	}
}

func (r *Rectangle[T]) Move(dp Point[T]) {
	r.Min.Move(dp)
	r.Max.Move(dp)
}

func (r Rectangle[T]) Dx() T {
	return r.Max.X - r.Min.X
}

func (r Rectangle[T]) Dy() T {
	return r.Max.Y - r.Min.Y
}

func (r Rectangle[T]) Size() Point[T] {
	return Point[T]{r.Max.X - r.Min.X, r.Max.Y - r.Min.Y}
}

func (r Rectangle[T]) NW() Point[T] {
    return r.Min
}
func (r Rectangle[T]) N() Point[T] {
    return Point[T]{(r.Min.X+r.Max.X)/2, r.Min.Y}
}
func (r Rectangle[T]) NE() Point[T] {
    return Point[T]{r.Max.X, r.Min.Y}
}
func (r Rectangle[T]) W() Point[T] {
    return Point[T]{r.Min.X, (r.Min.Y+r.Max.Y)/2}
}
func (r Rectangle[T]) E() Point[T] {
    return Point[T]{r.Max.X, (r.Min.Y+r.Max.Y)/2}
}
func (r Rectangle[T]) SW() Point[T] {
    return Point[T]{r.Min.X, r.Max.Y}
}
func (r Rectangle[T]) S() Point[T] {
    return Point[T]{(r.Min.X+r.Max.X)/2, r.Max.Y}
}
func (r Rectangle[T]) SE() Point[T] {
    return r.Max
}
func (r Rectangle[T]) Center() Point[T] {
    return Point[T]{(r.Min.X+r.Max.X)/2, (r.Min.Y+r.Max.Y)/2}
}
func (r Rectangle[T]) RelPos(rx, ry float64) Point[T] {
	return r.Min.AddXY(T(rx*float64(r.Dx())), T(ry*float64(r.Dy())))
}

func (r Rectangle[T]) Intersect(s Rectangle[T]) Rectangle[T] {
	if r.Empty() {
		return Rectangle[T]{}
	} else {
		return Rect(max(r.Min.X, s.Min.X), max(r.Min.Y, s.Min.Y),
			min(r.Max.X, s.Max.X), min(r.Max.Y, s.Max.Y))
	}
}

func (r Rectangle[T]) Union(s Rectangle[T]) Rectangle[T] {
	if r.Empty() {
		return s
	}
	if s.Empty() {
		return r
	}
	return Rect(min(r.Min.X, s.Min.X), min(r.Min.Y, s.Min.Y),
		max(r.Max.X, s.Max.X), max(r.Max.Y, s.Max.Y))
}

func (r Rectangle[T]) Inset(dx, dy T) Rectangle[T] {
	return Rectangle[T]{Min: r.Min.AddXY(dx, dy), Max: r.Max.SubXY(dx, dy)}
}

// Prüft, ob zwei Rechteck gleich sind, dh. die exakt gleichen Koordinaten
// haben.
func (r Rectangle[T]) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

func (r Rectangle[T]) Eq(s Rectangle[T]) bool {
	return r == s || r.Empty() && s.Empty()
}

func (r Rectangle[T]) Overlaps(s Rectangle[T]) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

func (r Rectangle[T]) In(s Rectangle[T]) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	return r.Min.X >= s.Min.X && r.Max.X <= s.Max.X &&
		r.Min.Y >= s.Min.Y && r.Max.Y <= s.Max.Y
}

func (r Rectangle[T]) AsCoord() (x, y, w, h T) {
	return r.Min.X, r.Min.Y, r.Dx(), r.Dy()
}

func (r Rectangle[T]) ToFloat() (Rectangle[float64]) {
	return Rectangle[float64]{r.Min.ToFloat(), r.Max.ToFloat()}
}

func (r Rectangle[T]) ToInt() (image.Rectangle) {
	return image.Rectangle{r.Min.ToInt(), r.Max.ToInt()}
}

func (r Rectangle[T]) AsMinMax() (xMin, yMin, xMax, yMax T) {
	return r.Min.X, r.Min.Y, r.Max.X, r.Max.Y
}

