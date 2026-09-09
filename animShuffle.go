package main

import (
	"log"
	"math/rand/v2"
	"image"
	"image/draw"
	"time"

	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
)

type ShuffleAnim struct {
	eventHandlerEmbed
	windowEmbed
	size image.Rectangle
	img *image.RGBA
	fg image.Image
	originList []image.Point
	idxList []int
	sz image.Rectangle
	orig image.Point
	t0 time.Time
	dt time.Duration
}

func NewShuffleAnimation(bounds geom.Rectangle[int]) *ShuffleAnim {
	var err error

	a := &ShuffleAnim{}
	a.bounds = bounds
	a.size = image.Rectangle{Max: bounds.Size().ToInt()}
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	a.img = a.gc.Image().(*image.RGBA)
	a.fg, err = gg.LoadPNG("images/08_TwoWord.png")
	if err != nil {
		log.Fatal(err)
	}
	a.originList = make([]image.Point, 44)
	a.idxList = make([]int, len(a.originList))
	a.sz = image.Rect(0, 0, 80, 80)
	a.orig = image.Point{}

	a.SetOnClick(func(ev MouseEvent) {
		if ev.Button.IsSet(LeftButton) {
			a.dt = 30 * time.Millisecond
		}
		if ev.Button.IsSet(RightButton) {
			a.dt = time.Hour
			for i, _ := range a.idxList {
				a.idxList[i] = i
			}
		}
	})

	return a
}

func (a *ShuffleAnim) Init() {
	p0 := image.Point{8, 0}
	dx := image.Point{a.sz.Dx() + 8, 0}
	dy := image.Point{0, a.sz.Dy() + 8}

	for i, _ := range a.originList {
		col := i % 4
		row := i / 4
		a.originList[i] = p0.Add(dx.Mul(col)).Add(dy.Mul(row))
	}
	for i, _ := range a.idxList {
		a.idxList[i] = i
	}
	a.t0 = time.Now()
	a.dt = 30 * time.Millisecond
}

func (a *ShuffleAnim) Update(dt time.Duration) {
	if time.Since(a.t0) > a.dt {
		rand.Shuffle(len(a.idxList), func(i, j int) {
			a.idxList[i], a.idxList[j] = a.idxList[j], a.idxList[i]
		})
		a.dt = 6*a.dt/5
		a.t0 = time.Now()
	}
}

func (a *ShuffleAnim) Refresh() {
	//draw.Draw(a.img, a.size, a.bg, image.Point{}, draw.Over)
	a.gc.Clear(colors.DimGray.Dark(0.6))
	for i := range a.idxList {
		srcPt := a.originList[i]
		dstPt := a.originList[a.idxList[i]]
		draw.Draw(a.img, a.sz.Add(dstPt).Add(a.orig), a.fg,
			srcPt, draw.Over)
	}
}

