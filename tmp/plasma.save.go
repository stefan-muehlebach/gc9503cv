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
	t       float64
	pal     colors.Palette
	orderQ  [numThreads]chan float64
	doneQ   [numThreads]chan bool
	valFlds [numThreads]*ValFieldType
	rect    geom.Rectangle
	off     image.Point
}

func NewPlasmaAnimation(pal colors.Palette, rect geom.Rectangle, off geom.Point) *PlasmaAnim {
	a := &PlasmaAnim{}

	a.t = 0.0
	a.pal = pal
	for i := range numThreads {
		a.orderQ[i] = make(chan float64)
		a.doneQ[i] = make(chan bool)
		a.valFlds[i] = NewValField(int(rect.Dx()), int(rect.Dy()),
			-dispWidth/2.0, dispWidth/2.0,
			dispHeight/2.0, -dispHeight/2.0,
			ColorFuncList[i])
		go UpdateThread(a.valFlds[i], a.orderQ[i], a.doneQ[i])
	}
	a.rect = rect
	a.off = off.Int()
	return a
}

func (a *PlasmaAnim) Init(gc *gg.Context) {
	a.t = 0.0
}

func (a *PlasmaAnim) Animate(dt time.Duration) {
	a.t += adt
	for i := range numThreads {
		a.orderQ[i] <- a.t
	}
	for i := range numThreads {
		<-a.doneQ[i]
	}
}

func (a *PlasmaAnim) Draw(gc *gg.Context) {
	valIdx := 0
	for row := range int(a.rect.Dy()) {
		for col := range int(a.rect.Dx()) {
			v1 := a.valFlds[0].Vals[valIdx]
			v2 := a.valFlds[1].Vals[valIdx]
			v3 := a.valFlds[2].Vals[valIdx]
			v := (v1 + v2 + v3 + 3.0) / 6.0
			c := a.pal.Color(v)
			gc.SetPixel(col+a.off.X, row+a.off.Y, c)
			valIdx += 1
		}
	}
}

func UpdateThread(valFld *ValFieldType, orderQ chan float64, doneQ chan bool) {
	var t float64
	var ok bool

	for {
		if t, ok = <-orderQ; !ok {
			break
		}
		valFld.Update(t)
		doneQ <- true
	}
}

type ValFieldType struct {
	Vals               []float64
	Cols, Rows         int
	Xmin, Ymax, Dx, Dy float64
	Fnc                ColorFuncType
}

func NewValField(cols, rows int, xmin, xmax, ymin, ymax float64,
	fnc ColorFuncType) *ValFieldType {
	v := &ValFieldType{}
	v.Vals = make([]float64, cols*rows)
	v.Cols = cols
	v.Rows = rows
	v.Xmin = xmin
	v.Ymax = ymax
	v.Dx = (xmax - xmin) / float64(cols)
	v.Dy = (ymax - ymin) / float64(rows)
	v.Fnc = fnc
	return v
}

func (v *ValFieldType) Update(t float64) {
	var x, y float64
	var col, row, idx int

	y = v.Ymax
	idx = 0
	for row = 0; row < v.Rows; row++ {
		x = v.Xmin
		for col = 0; col < v.Cols; col++ {
			v.Vals[idx] = v.Fnc(x, y, t)
			x += v.Dx
			idx++
		}
		y -= v.Dy
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
