package main

import (
	"image"
	"image/draw"
	"log"
	"path/filepath"
	"time"

	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
)

var (
	ImageDir      = "images"
	ImageFileList = []string{
		"00_Fish.png",
		"01_Flowers.png",
		"02_Reef.png",
		"03_Space.png",
		"04_Hoehenweg.png",
		"05_SeeuferA.png",
		"06_SeeuferB.png",
		"12_TwoWord.png",
		"14_GoRegular.png",
	}
	ImageShowDuration = 3 * time.Second
)

//----------------------------------------------------------------------------

type ImageWin struct {
	windowEmbed
	img             *image.RGBA
	images          []image.Image
	curIdx, nextIdx int
	t0              time.Time
}

func NewImageWindow(bounds geom.Rectangle[int]) *ImageWin {
	a := &ImageWin{}
	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())
	a.img = a.gc.Image().(*image.RGBA)
	a.images = make([]image.Image, 0)
	for _, fileName := range ImageFileList {
		png, err := gg.LoadPNG(filepath.Join(ImageDir, fileName))
		if err != nil {
			log.Fatalf("Couldn't load %s: %v", fileName, err)
		}
		a.images = append(a.images, png)
	}

	a.Root = NewGroup()
	a.Root.SetSize(bounds.ToFloat().Size())

	return a
}

func (a *ImageWin) Init() {
	a.curIdx = -1
	a.nextIdx = 0
	a.t0 = time.Now()
}

func (a *ImageWin) Update(t time.Time) {
	if time.Since(a.t0) > ImageShowDuration {
		a.nextIdx = (a.nextIdx + 1) % len(a.images)
		a.t0 = time.Now()
	}
}

func (a *ImageWin) Refresh() {
	if a.curIdx != a.nextIdx {
		a.curIdx = a.nextIdx
		draw.Draw(a.img, a.images[a.curIdx].Bounds(), a.images[a.curIdx], image.Point{}, draw.Over)
	}
}
