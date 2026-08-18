package main

import (
	"fmt"
	"math"
)

//----------------------------------------------------------------------------

type ProjMatrix struct {
	M11, M22, M33, M34, M43 float64
}

func NewProjMatrix(near, far float64) ProjMatrix {
	return ProjMatrix{1.0, 1.0, (near + far) / near, -far, 1.0 / near}
}

func (m ProjMatrix) Project(v Vector) Vector {
	w := Vector{m.M11 * v.X, m.M22 * v.Y, m.M33*v.Z - m.M34}
	f := m.M43 * v.Z
	return w.Div(f)
}

var (
	left, right = -500.0, 500.0
	top, bottom = 180.0, -180.0
	near, far   = -10.0, 10.0

	PM = NewProjMatrix(near, far)
)

//----------------------------------------------------------------------------

type Vector struct {
	X, Y, Z float64
}

func NewVector(x, y, z float64) Vector {
	return Vector{x, y, z}
}

func (v Vector) Add(w Vector) Vector {
	return Vector{v.X + w.X, v.Y + w.Y, v.Z + w.Z}
}

func (v Vector) Sub(w Vector) Vector {
	return Vector{v.X - w.X, v.Y - w.Y, v.Z - w.Z}
}

func (v Vector) Mul(k float64) Vector {
	return Vector{k * v.X, k * v.Y, k * v.Z}
}

func (v Vector) Div(k float64) Vector {
	return Vector{v.X / k, v.Y / k, v.Z / k}
}

func (v Vector) Neg() Vector {
	return Vector{-v.X, -v.Y, -v.Z}
}

func (v Vector) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector) Dot(w Vector) float64 {
	return v.X*w.X + v.Y*w.Y + v.Z*w.Z
}

func (v Vector) Cross(w Vector) Vector {
	return Vector{
		v.Y*w.Z - w.Y*v.Z,
		v.Z*w.X - w.Z*v.X,
		v.X*w.Y - w.X*v.Y,
	}
}

//----------------------------------------------------------------------------

type Matrix [12]float64

func Identity() Matrix {
	return Matrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
	}
}

func Translate(v Vector) Matrix {
	return Matrix{
		1, 0, 0, v.X,
		0, 1, 0, v.Y,
		0, 0, 1, v.Z,
	}
}

func Scale(sv Vector) Matrix {
	return Matrix{
		sv.X, 0, 0, 0,
		0, sv.Y, 0, 0,
		0, 0, sv.Z, 0,
	}
}

func RotateX(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix{
		1, 0, 0, 0,
		0, c, -s, 0,
		0, s, c, 0,
	}
}

func RotateY(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix{
		c, 0, -s, 0,
		0, 1, 0, 0,
		s, 0, c, 0,
	}
}

func RotateZ(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix{
		c, -s, 0, 0,
		s, c, 0, 0,
		0, 0, 1, 0,
	}
}

func (a Matrix) Mul(b Matrix) Matrix {
	return Matrix{
		a[0]*b[0] + a[1]*b[4] + a[2]*b[8],
		a[0]*b[1] + a[1]*b[5] + a[2]*b[9],
		a[0]*b[2] + a[1]*b[6] + a[2]*b[10],
		a[0]*b[3] + a[1]*b[7] + a[2]*b[11] + a[3],
		a[4]*b[0] + a[5]*b[4] + a[6]*b[8],
		a[4]*b[1] + a[5]*b[5] + a[6]*b[9],
		a[4]*b[2] + a[5]*b[6] + a[6]*b[10],
		a[4]*b[3] + a[5]*b[7] + a[6]*b[11] + a[7],
		a[8]*b[0] + a[9]*b[4] + a[10]*b[8],
		a[8]*b[1] + a[9]*b[5] + a[10]*b[9],
		a[8]*b[2] + a[9]*b[6] + a[10]*b[10],
		a[8]*b[3] + a[9]*b[7] + a[10]*b[11] + a[11],
	}
}

func (a Matrix) Transform(v Vector) Vector {
	return Vector{
		a[0]*v.X + a[1]*v.Y + a[2]*v.Z + a[3],
		a[4]*v.X + a[5]*v.Y + a[6]*v.Z + a[7],
		a[8]*v.X + a[9]*v.Y + a[10]*v.Z + a[11],
	}
}

func (a Matrix) Det() float64 {
	return a[0]*a[5]*a[10] + a[1]*a[6]*a[8] + a[2]*a[4]*a[9] -
		a[0]*a[6]*a[9] - a[1]*a[4]*a[10] - a[2]*a[5]*a[8]
}

func (a Matrix) Inv() Matrix {
	det := a.Det()
	adj := Matrix{}
	adj[0] = Matrix3{
		a[5], a[6], a[7],
		a[9], a[10], a[11],
	}.Det() / det
	adj[1] = -Matrix3{
		a[4], a[6], a[7],
		a[8], a[10], a[11],
	}.Det() / det
	adj[2] = Matrix3{
		a[4], a[5], a[7],
		a[8], a[9], a[11],
	}.Det() / det
	adj[3] = -Matrix3{
		a[4], a[5], a[6],
		a[8], a[9], a[10],
	}.Det() / det

	adj[4] = -Matrix3{
		a[1], a[2], a[3],
		a[9], a[10], a[11],
	}.Det() / det
	adj[5] = Matrix3{
		a[0], a[2], a[3],
		a[8], a[10], a[11],
	}.Det() / det
	adj[6] = -Matrix3{
		a[0], a[1], a[3],
		a[8], a[9], a[11],
	}.Det() / det
	adj[7] = Matrix3{
		a[0], a[1], a[2],
		a[8], a[9], a[10],
	}.Det() / det

	adj[8] = Matrix3{
		a[1], a[2], a[3],
		a[5], a[6], a[7],
	}.Det() / det
	adj[9] = -Matrix3{
		a[0], a[2], a[3],
		a[4], a[6], a[7],
	}.Det() / det
	adj[10] = Matrix3{
		a[0], a[1], a[3],
		a[4], a[5], a[7],
	}.Det() / det
	adj[11] = -Matrix3{
		a[0], a[1], a[2],
		a[4], a[5], a[6],
	}.Det() / det

	return Matrix{adj[0], adj[4], adj[8], adj[3],
		adj[1], adj[5], adj[9], adj[7],
		adj[2], adj[6], adj[10], adj[11],
	}
}

func (a Matrix) String() string {
	return fmt.Sprintf("[%.4v %.4v %.4v %.4v]\n[%.4v %.4v %.4v %.4v]\n[%.4v %.4v %.4v %.4v]",
		a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9],
		a[10], a[11])
}

//----------------------------------------------------------------------------

type Matrix3 [6]float64

func (a Matrix3) Det() float64 {
	return a[0]*a[4] - a[1]*a[3]
}

func (a Matrix3) String() string {
	return fmt.Sprintf("[%.4v %.4v %.4v]\n[%.4v %.4v %.4v]",
		a[0], a[1], a[2], a[3], a[4], a[5])
}
