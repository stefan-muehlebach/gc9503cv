package geom

//----------------------------------------------------------------------------

type CoordType interface {
	~int | ~int16 | ~int32 | ~float32 | ~float64
}

//----------------------------------------------------------------------------

type TransfDirection byte

const (
	Base2Sub TransfDirection = iota
	Sub2Base
)

func (dir TransfDirection) String() string {
	switch dir {
	case Base2Sub:
		return "Base2Sub"
	case Sub2Base:
		return "Sub2Base"
	default:
		return "(unknown transform direction)"
	}
}

type Transform[T CoordType] struct {
	Rect Rectangle[T]
	Rot  RotationType
}

func NewTransform[T CoordType](rot RotationType, rect Rectangle[T]) *Transform[T] {
	t := &Transform[T]{}
	t.Rot = rot
	t.Rect = rect
	return t
}

func (t *Transform[T]) Size(s0 Point[T], dir TransfDirection) (s1 Point[T]) {
	switch t.Rot {
	case Rot000, Rot180:
		s1.X, s1.Y = s0.X, s0.Y
	case Rot090, Rot270:
		s1.X, s1.Y = s0.Y, s0.X
	}
	return s1
}

func (t *Transform[T]) Point(p0 Point[T], dir TransfDirection) (p1 Point[T]) {
	switch t.Rot {
	case Rot000:
		p1 = p0
	case Rot090:
		if dir == Base2Sub {
			p1 = Pt(t.Rect.Max.Y-p0.Y, p0.X)
		} else {
			p1 = Pt(p0.Y, t.Rect.Max.Y-p0.X)
		}
	case Rot180:
		p1 = Pt(t.Rect.Max.X-p0.X, t.Rect.Max.Y-p0.Y)
	case Rot270:
		if dir == Base2Sub {
			p1 = Pt(p0.Y, t.Rect.Max.X-p0.X)
		} else {
			p1 = Pt(t.Rect.Max.X-p0.Y, p0.X)
		}
	}
	return p1
}

func (t *Transform[T]) Rectangle(r0 Rectangle[T], dir TransfDirection) (r1 Rectangle[T]) {
	switch t.Rot {
	case Rot000:
		r1 = r0
	case Rot090:
		if dir == Base2Sub {
			r1 = Rectangle[T]{
				Min: t.Point(Point[T]{r0.Min.X, r0.Max.Y}, dir),
				Max: t.Point(Point[T]{r0.Max.X, r0.Min.Y}, dir),
			}
		} else {
			r1 = Rectangle[T]{
				Min: t.Point(Point[T]{r0.Max.X, r0.Min.Y}, dir),
				Max: t.Point(Point[T]{r0.Min.X, r0.Max.Y}, dir),
			}
		}
	case Rot180:
		r1 = Rectangle[T]{
			Min: t.Point(r0.Max, dir),
			Max: t.Point(r0.Min, dir),
		}
	case Rot270:
		if dir == Base2Sub {
			r1 = Rectangle[T]{
				Min: t.Point(Point[T]{r0.Max.X, r0.Min.Y}, dir),
				Max: t.Point(Point[T]{r0.Min.X, r0.Max.Y}, dir),
			}
		} else {
			r1 = Rectangle[T]{
				Min: t.Point(Point[T]{r0.Min.X, r0.Max.Y}, dir),
				Max: t.Point(Point[T]{r0.Max.X, r0.Min.Y}, dir),
			}
		}
	}
	return r1
}
