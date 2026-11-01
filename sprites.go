package main

import (
	"image"
	"time"

	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"golang.org/x/image/draw"
)

//----------------------------------------------------------------------------

type ImageSeq struct {
	imgList   []image.Image
	Size      image.Point
	PosOffset image.Point
	//Divider     int
	//AutoReverse bool
}

func NewImageSeq(src *image.NRGBA, start, size, num, gap,
	posOffset image.Point) *ImageSeq {
	s := &ImageSeq{}
	s.imgList = make([]image.Image, 0)

	i := 0
	pt := image.Point{}
	for row := range num.Y {
		pt.Y = start.Y + row*(size.Y+gap.Y)
		for col := range num.X {
			pt.X = start.X + col*(size.X+gap.X)
			rect := image.Rectangle{Min: pt, Max: pt.Add(size)}
			s.imgList = append(s.imgList, src.SubImage(rect))
			i++
		}
	}
	s.Size = size
	s.PosOffset = posOffset
	//s.Divider = div
	//s.AutoReverse = false

	return s
}

func (s *ImageSeq) Len() int {
	return len(s.imgList)
}

func (s *ImageSeq) Image(idx int) image.Image {
	return s.imgList[idx%len(s.imgList)]
}

//----------------------------------------------------------------------------

type Sprite struct {
	Src             *ImageSeq
	Idx             int
	Divider, divPos int
	RepeatCount     int
	repeatsLeft     int
	AutoReverse     bool
	reverse         bool
	isRunning       bool
}

func NewSprite(imgs *ImageSeq) *Sprite {
	s := &Sprite{}
	s.Src = imgs
	return s
}

func NewFixedSprite(imgs *ImageSeq, idx int) *Sprite {
	s := &Sprite{}
	s.Src = imgs
	s.Idx = idx
	s.Divider = -1
	return s
}

func (s *Sprite) Start() {
	if s.isRunning {
		return
	}
	s.repeatsLeft = s.RepeatCount
	s.reverse = false
	s.isRunning = true
}

func (s *Sprite) Stop() {
	if !s.isRunning {
		return
	}
	s.isRunning = false
}

func (s *Sprite) Update(t time.Time) {
	if !s.isRunning || s.Divider < 0 {
		return
	}
	if s.divPos == 0 {
		s.divPos = s.Divider
		if s.reverse {
			s.Idx -= 1
			if s.Idx == 0 {
				if s.repeatsLeft > 0 {
					s.repeatsLeft--
				}
				if s.repeatsLeft == 0 {
					s.isRunning = false
					return
				}
				s.reverse = false
			}
		} else {
			s.Idx += 1
			if s.Idx == s.Src.Len() {
				if s.AutoReverse {
					s.reverse = true
					s.Idx--
				} else {
					if s.repeatsLeft > 0 {
						s.repeatsLeft--
					}
					s.Idx = 0
					if s.repeatsLeft == 0 {
						s.isRunning = false
						return
					}
				}
			}
		}
	} else {
		s.divPos--
	}
}

func (s *Sprite) Image() image.Image {
	return s.Src.Image(s.Idx)
}

//----------------------------------------------------------------------------

type Object struct {
	nodeEmbed
	Spr *Sprite
}

func NewObject() *Object {
	o := &Object{}
	o.Init(o)
	return o
}

func NewObject2(spr *Sprite, pos Point) *Object {
	o := &Object{}
	o.Init(o)
	o.SetSprite(spr)
	o.SetPos(pos)
	return o
}

func (o *Object) SetSprite(spr *Sprite) {
	o.Spr = spr
	o.SetSize(geom.Pt[float64](o.Spr.Src.Size))
	o.SetPosOffset(geom.Pt[float64](o.Spr.Src.PosOffset))
	o.Spr.RepeatCount = -1
	o.Spr.Start()
}

func (o *Object) Update(t time.Time) {
	o.Spr.Update(t)
}

func (o *Object) Draw(gc *gg.Context) {
	img := gc.Image().(*image.RGBA)
	p0 := o.pos.Sub(o.posOffset)
	draw.Copy(img, p0.ToInt(), o.Spr.Image(), o.Spr.Image().Bounds(),
		draw.Over, nil)

	if GlobDebug {
		gc.SetLineColor(colors.Red)
		gc.SetLineWidth(1.0)
		gc.DrawLine(o.pos.X-15, o.pos.Y, o.pos.X+15, o.pos.Y)
		gc.DrawLine(o.pos.X, o.pos.Y-15, o.pos.X, o.pos.Y+15)
		gc.Stroke()
		gc.DrawCircle(o.pos.X, o.pos.Y, 10)
		gc.Stroke()
	}
}
