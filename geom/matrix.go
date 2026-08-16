package geom

import (
	"fmt"
	"math"
)

// Contains the values of a 3x3 transformation matrix as a slice of float
// values. The slice is in row major order. The last row of the matrix is
// always [0, 0, 1] and will not be stored.
type Matrix [6]float64

// Creates the identity matrix.
func Identity() *Matrix {
	return &Matrix{1.0, 0.0, 0.0,
		0.0, 1.0, 0.0}
}

// Creates a translation matrix for a translation along the vector d.
func Translate(d Point[float64]) *Matrix {
	return &Matrix{1.0, 0.0, d.X,
		0.0, 1.0, d.Y}
}

// Creates a rotation matrix around the origin of an angle of a (in radians)
// counter clockwise.
func Rotate(a float64) *Matrix {
	s := math.Sin(a)
	c := math.Cos(a)
	return &Matrix{c, -s, 0.0,
		s, c, 0.0}
}

// Erzeugt eine Rotationsmatrix um den Winkel a (im Bogenmass) mit Drehpunkt
// bei rp.
func RotateAbout(rp Point[float64], a float64) *Matrix {
	m := Translate(rp)
	m = m.Rotate(a)
	m = m.Translate(rp.Neg())
	return m
}

// Erzeugt eine Skalierungsmatrix mit den Skalierungsfaktoren sx in X-Richtung
// und sy in Y-Richtung. Zentrum der Skalierung ist der Ursprung des
// Koordinatensystems.
func Scale(sx, sy float64) *Matrix {
	return &Matrix{sx, 0.0, 0.0,
		0.0, sy, 0.0}
}

// Erzeugt eine Skalierungsmatrix mit den Skalierungsfaktoren sx in X-Richtung
// und sy in Y-Richtung. Zentrum der Skalierung ist der Punkt sp.
func ScaleAbout(sp Point[float64], sx, sy float64) *Matrix {
	m := Translate(sp)
	m = m.Scale(sx, sy)
	m = m.Translate(sp.Neg())
	return m
}

// Returns the inverse of the matrix a.
func (a *Matrix) Inv() *Matrix {
	det := a[0]*a[4] - a[1]*a[3]
	return &Matrix{a[4] / det,
		-a[1] / det,
		(a[1]*a[5] - a[4]*a[2]) / det,

		-a[3] / det,
		a[0] / det,
		(a[3]*a[2] - a[0]*a[5]) / det}
}

// Returns the product of the matrices a and b.
func (a *Matrix) Multiply(b *Matrix) *Matrix {
	return &Matrix{a[0]*b[0] + a[1]*b[3],
		a[0]*b[1] + a[1]*b[4],
		a[0]*b[2] + a[1]*b[5] + a[2],

		a[3]*b[0] + a[4]*b[3],
		a[3]*b[1] + a[4]*b[4],
		a[3]*b[2] + a[4]*b[5] + a[5]}
}

// Translates the matrix a by d.
func (m *Matrix) Translate(d Point[float64]) *Matrix {
	return m.Multiply(Translate(d))
}

// Rotates the matrix m by angle counter clockwise around the origin.
func (m *Matrix) Rotate(angle float64) *Matrix {
	return m.Multiply(Rotate(angle))
}

// Siehe Translate
func (m *Matrix) RotateAbout(rp Point[float64], a float64) *Matrix {
	return m.Multiply(RotateAbout(rp, a))
}

// Siehe Translate
func (m *Matrix) Scale(sx, sy float64) *Matrix {
	return m.Multiply(Scale(sx, sy))
}

// Siehe Translate
func (m *Matrix) ScaleAbout(sp Point[float64], sx, sy float64) *Matrix {
	return m.Multiply(ScaleAbout(sp, sx, sy))
}

// Hier schliesslich spielt die Musik: eine Matrix wird für die Transformation
// eines Punktes verwendet. Es wird ein neuer Punkt erstellt, aktuell gibt
// noch keinen Bedarf nach 'in place' Transformation.
func (m *Matrix) Transform(p Point[float64]) Point[float64] {
	return Point[float64]{m[0]*p.X + m[1]*p.Y + m[2],
		m[3]*p.X + m[4]*p.Y + m[5]}
}

// Will man den zweiten geometrischen Typ (Rectangle) transformieren, dann
// ist dies die Methode der Wahl.
func (m *Matrix) TransformRect(r Rectangle[float64]) Rectangle[float64] {
	return Rectangle[float64]{m.Transform(r.Min), m.Transform(r.Max)}
}

// TransformPoint und TransformVector sind zwei Methoden, die aus dem
// gg-Package übernommen wurden: sie operieren auf Punkte, resp. Vektoren, die
// als Paar von Float-Werten angegeben werden.
func (m *Matrix) TransformPoint(x, y float64) (float64, float64) {
	return m[0]*x + m[1]*y + m[2], m[3]*x + m[4]*y + m[5]
}

// Während TransformPoint eine vollständige Transformation durchführt, wird
// bei TransformVector die Translation ignoriert.
func (m Matrix) TransformVector(x, y float64) (float64, float64) {
	return m[0]*x + m[1]*y, m[3]*x + m[4]*y
}

// Mit dieser Methode implementiert Matrix schliesslich das 'Stringer'
// Interface und kann bequem per Printf("%v", m) ausgegeben werden.
// Beachte, dass dieser String Zeilenumbrüche enthält!
func (m Matrix) String() string {
	return fmt.Sprintf("[ % .4f  % .4f  % .4f ]\n[ % .4f  % .4f  % .4f ]\n[ % .4f  % .4f  % .4f ]",
		m[0], m[1], m[2],
		m[3], m[4], m[5],
		0.0, 0.0, 1.0)
}

// Im Package golang.org/x/image/math/f64 sind affine Transformationen als
// Arrays von Float-Werten definiert. Für die Konvertierung in diesen
// Typ steht die folgende Methode zur Verfügung.
func (m Matrix) AsAff3() [6]float64 {
	return [...]float64{m[0], m[1], m[2], m[3], m[4], m[5]}
}

func (m Matrix) AsMat3() [9]float64 {
	return [...]float64{m[0], m[1], m[2], m[3], m[4], m[5], 0, 0, 1}
}
