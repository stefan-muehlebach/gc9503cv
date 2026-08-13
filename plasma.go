//go:build ignore

package main

import (
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/geom"
	"image"
	"math"
	"time"
)

// The famous plasma animation ------------------------------------------------
var (
	ColorFuncList = []ColorFuncType{
		ColorFunc01,
		ColorFunc02,
		ColorFunc03,
	}
)

const (
	numThreads = 3
	numShades  = 256
	dispWidth  = 1.6
	dispHeight = 4.2
	adt        = 0.05
)

type PlasmaAnim struct {
	t            float64
	pal          colors.Palette
	valFld       [][]float64
	valMin       geom.Point
	valDx, valDy float64
	rect         geom.Rectangle
	off          image.Point
}

func NewPlasmaAnimation(pal colors.Palette, rect geom.Rectangle, off geom.Point) *PlasmaAnim {
	a := &PlasmaAnim{}

	a.t = 0.0
	a.pal = pal
	a.valFld = make([][]float64, int(rect.Dx()))
	for i, _ := range a.valFld {
		a.valFld[i] = make([]float64, int(rect.Dy()))
	}
	a.valMin = geom.Point{-dispWidth / 2.0, -dispHeight / 2.0}
	a.valDx = dispWidth / rect.Dx()
	a.valDy = dispHeight / rect.Dy()
	a.rect = rect
	a.off = off.Int()
	return a
}

func (a *PlasmaAnim) Init(gc *gg.Context) {
	a.t = 0.0
}

func (a *PlasmaAnim) Update(dt time.Duration) {
	a.t += adt
	y := a.valMin.Y
	for row := range int(a.rect.Dy()) {
		x := a.valMin.X
		for col := range int(a.rect.Dx()) {
			v1 := ColorFunc01(x, y, a.t)
			v2 := ColorFunc02(x, y, a.t)
			v3 := ColorFunc03(x, y, a.t)
			v := (v1 + v2 + v3 + 3.0) / 6.0
			a.valFld[col][row] = v
			x += a.valDx
		}
		y += a.valDy
	}
}

func (a *PlasmaAnim) Draw(gc *gg.Context) {
	for row := range int(a.rect.Dy()) {
		for col := range int(a.rect.Dx()) {
			c := a.pal.Color(a.valFld[col][row])
			gc.SetPixel(col+a.off.X, row+a.off.Y, c)
		}
	}
}

type ColorFuncType func(x, y, t float64) float64

const (
	f1p1 = 10.0

	f2p1 = 10.0
	f2p2 = 2.0
	f2p3 = 3.0

	f3p1 = 5.0
	f3p2 = 3.0
)

func ColorFunc01(x, y, t float64) float64 {
	return math.Sin(x*f1p1 + t)
}

func ColorFunc02(x, y, t float64) float64 {
	return math.Sin(f2p1*(x*math.Sin(t/f2p2)+y*math.Cos(t/f2p3)) + t)
}

func ColorFunc03(x, y, t float64) float64 {
	cx := x + 0.5*math.Sin(t/f3p1)
	cy := y + 0.5*math.Cos(t/f3p2)
	return math.Sin(math.Sqrt(100.0*(cx*cx+cy*cy)+1.0) + t)
}

/*
type Palette struct {
	name      string
	colorList []colors.RGBA
	shadeList []colors.RGBA
}

func NewPalette(name string, colors ...colors.RGBA) *Palette {
	p := &Palette{}
	p.name = name
	for _, color := range colors {
		p.colorList = append(p.colorList, color)
	}
	p.CalcLinearShades()
	return p
}

func (p *Palette) GetColor(t float64) colors.RGBA {
	idx := int(t * float64(numShades*(len(p.colorList)-1)))
	return p.shadeList[idx]
}

// Berechnet die Abstufung zwischen den Stuetzfarben einfach, d.h.
// linear.
func (p *Palette) CalcLinearShades() {
	var i, j, k int
	var color1, color2 colors.RGBA
	var t float64

	p.shadeList = make([]colors.RGBA, numShades*(len(p.colorList)-1))
	for i = 0; i < len(p.colorList)-1; i++ {
		for j = 0; j < numShades; j++ {
			color1 = p.colorList[i]
			color2 = p.colorList[i+1]
			t = float64(j) / float64(numShades)
			k = i*numShades + j
			p.shadeList[k] = color1.Interpolate(color2, t)
		}
	}
}
*/
