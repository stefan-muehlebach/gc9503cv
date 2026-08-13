package geom

import "errors"

//----------------------------------------------------------------------------

type CoordType interface {
	~int | ~int16 | ~int32 | ~float32 | ~float64
}

// Rotationsmöglichkeiten des Displays. Es gibt (logischerweise) 4
// Möglichkeiten das Display zu rotieren. Dies hat Auswirkungen auf die
// Initialisierung des Displays, auf die globalen Variablen Width und Height
// und auf die Konfigurationsdateien, in welchen die Daten für die
// Transformation von Touch-Koordinaten auf Display-Koordianten abgelegt
// sind, etc.
type RotationType byte

const (
	Rotate000 RotationType = iota
	Rotate090
	Rotate180
	Rotate270
)

func (rot RotationType) String() string {
	switch rot {
	case Rotate000:
		return "Rotate000"
	case Rotate090:
		return "Rotate090"
	case Rotate180:
		return "Rotate180"
	case Rotate270:
		return "Rotate270"
	default:
		return "(unknown rotation)"
	}
}

func (rot *RotationType) Set(s string) error {
	switch s {
	case "Rotate000", "000", "0":
		*rot = Rotate000
	case "Rotate090", "090", "90":
		*rot = Rotate090
	case "Rotate180", "180":
		*rot = Rotate180
	case "Rotate270", "270":
		*rot = Rotate270
	default:
		return errors.New("Unknown rotation: " + s)
	}
	return nil
}

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
	case Rotate000, Rotate180:
        s1.X, s1.Y = s0.X, s0.Y
	case Rotate090, Rotate270:
        s1.X, s1.Y = s0.Y, s0.X
    }
    return s1
}

func (t *Transform[T]) Point(p0 Point[T], dir TransfDirection) (p1 Point[T]) {
	switch t.Rot {
	case Rotate000:
		p1 = p0
	case Rotate090:
		if dir == Base2Sub {
			p1 = Pt(t.Rect.Max.Y-p0.Y, p0.X)
		} else {
			p1 = Pt(p0.Y, t.Rect.Max.Y-p0.X)
		}
	case Rotate180:
		p1 = Pt(t.Rect.Max.X-p0.X, t.Rect.Max.Y-p0.Y)
	case Rotate270:
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
	case Rotate000:
		r1 = r0
	case Rotate090:
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
	case Rotate180:
		r1 = Rectangle[T]{
			Min: t.Point(r0.Max, dir),
			Max: t.Point(r0.Min, dir),
		}
	case Rotate270:
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

