package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"time"

	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
)

type TimeTable struct {
	t0      time.Time
	dur     time.Duration
	actions []TimeTableEntry
}

type TimeTableEntry struct {
	ts   time.Duration
	fnc  func()
	done bool
}

type SpriteWin struct {
	windowEmbed
	canvas *Group
	tt     TimeTable
}

const (
	F = 2
)

func Pt(x, y int) image.Point {
	return image.Pt(F*x, F*y)
}

var (
	Imgs = []struct {
		Img  image.Image
		File string
	}{
		{File: "images/F%1d_SuperMario.png"},
		{File: "images/F%1d_Ground.png"},
		{File: "images/F%1d_Summer.png"},
		{File: "images/F%1d_Fire.png"},
		{File: "images/F%1d_Glint.png"},
		{File: "images/F%1d_Sky.png"},
	}
)

func InitBackground(scr *Group) {
	var err error

	for i, data := range Imgs {
		fileName := fmt.Sprintf(data.File, F)
		Imgs[i].Img, err = gg.LoadImage(fileName)
		if err != nil {
			log.Fatalf("Could't load image: %v", err)
		}
	}

	// Several background sprites
	stoneSeq := NewImageSeq(Imgs[1].Img.(*image.NRGBA), Pt(2, 21),
		Pt(16, 16), image.Pt(8, 6), Pt(1, 1), Pt(0, 0))
	stoneL := NewFixedSprite(stoneSeq, 24)
	stoneM := NewFixedSprite(stoneSeq, 25)
	stoneR := NewFixedSprite(stoneSeq, 26)

	// Sprites for the plants and so
	bushSeq := NewImageSeq(Imgs[2].Img.(*image.NRGBA), Pt(229, 66),
		Pt(16, 16), image.Pt(3, 1), Pt(1, 0), Pt(0, 0))
	bushL := NewFixedSprite(bushSeq, 0)
	bushM := NewFixedSprite(bushSeq, 1)
	bushR := NewFixedSprite(bushSeq, 2)

	grassSeq := NewImageSeq(Imgs[2].Img.(*image.NRGBA), Pt(229, 83),
		Pt(16, 16), image.Pt(3, 1), Pt(1, 0), Pt(0, 0))
	grassL := NewFixedSprite(grassSeq, 0)
	grassM := NewFixedSprite(grassSeq, 1)
	grassR := NewFixedSprite(grassSeq, 2)

	pipeSeq := NewImageSeq(Imgs[2].Img.(*image.NRGBA), Pt(2, 21),
		Pt(16, 16), image.Pt(5, 5), Pt(1, 1), Pt(0, 0))
	pipeTL := NewFixedSprite(pipeSeq, 2)
	pipeTR := NewFixedSprite(pipeSeq, 3)
	pipeL := NewFixedSprite(pipeSeq, 7)
	pipeR := NewFixedSprite(pipeSeq, 8)

	Up := geom.Pt[float64](Pt(0, -16))
	Right := geom.Pt[float64](Pt(16, 0))

	p := Point{50, 200}
	for i := range 20 {
		var obj *Object

		pt := p.Add(Right.Mul(float64(i)))
		switch i {
		case 0:
			obj = NewObject2(stoneL, pt)
		case 19:
			obj = NewObject2(stoneR, pt)
		default:
			obj = NewObject2(stoneM, pt)
		}
		scr.Add(obj)

		switch i {
		case 2:
			pt1 := pt.Add(Up)
			green := NewObject2(grassL, pt1)
			scr.Add(green)
			pt1 = pt1.Add(Right)
			green = NewObject2(grassM, pt1)
			scr.Add(green)
			pt1 = pt1.Add(Right)
			green = NewObject2(grassR, pt1)
			scr.Add(green)
		case 6:
			pt1 := pt.Add(Up)
			green := NewObject2(bushL, pt1)
			scr.Add(green)
			pt1 = pt1.Add(Right)
			green = NewObject2(bushM, pt1)
			scr.Add(green)
			pt1 = pt1.Add(Right)
			green = NewObject2(bushR, pt1)
			scr.Add(green)
		case 10:
			pt1 := pt.Add(Up.Mul(2))
			pipe := NewObject2(pipeTL, pt1)
			scr.Add(pipe)
			pt1 = pt1.Add(Right)
			pipe = NewObject2(pipeTR, pt1)
			scr.Add(pipe)
			pt1 = pt.Add(Up)
			pipe = NewObject2(pipeL, pt1)
			scr.Add(pipe)
			pt1 = pt1.Add(Right)
			pipe = NewObject2(pipeR, pt1)
			scr.Add(pipe)
		}
	}
}

func NewSpriteWindow(bounds geom.Rectangle[int]) *SpriteWin {
	a := &SpriteWin{}
	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())

	root := NewGroup()
	root.Layout = NewMaxLayout()
	root.SetSize(bounds.ToFloat().Size())
	a.Root = root

	a.canvas = NewGroup()
	a.canvas.Layout = &NullLayout{}
	a.Root.Add(a.canvas)

	InitBackground(a.canvas)

	// Different shapes of Mario
	waitRightImg := NewImageSeq(Imgs[0].Img.(*image.NRGBA), Pt(2, 11),
		Pt(35, 35), image.Pt(4, 1), Pt(1, 0), Pt(14, 34))
	waitRight := NewSprite(waitRightImg)
	waitRight.Divider = 9
	waitRight.AutoReverse = true

	walkRightImg := NewImageSeq(Imgs[0].Img.(*image.NRGBA), Pt(2, 105),
		Pt(35, 35), image.Pt(8, 1), Pt(1, 0), Pt(14, 34))
	walkRight := NewSprite(walkRightImg)
	walkRight.Divider = 4

	runRightImg := NewImageSeq(Imgs[0].Img.(*image.NRGBA), Pt(2, 58),
		Pt(35, 35), image.Pt(8, 1), Pt(1, 0), Pt(14, 34))
	runRight := NewSprite(runRightImg)
	runRight.Divider = 2

	waitLeftImg := NewImageSeq(Imgs[0].Img.(*image.NRGBA), Pt(2, 199),
		Pt(35, 35), image.Pt(4, 1), Pt(1, 0), Pt(14, 34))
	waitLeft := NewSprite(waitLeftImg)
	waitLeft.Divider = 9
	waitLeft.AutoReverse = true

	walkLeftImg := NewImageSeq(Imgs[0].Img.(*image.NRGBA), Pt(2, 293),
		Pt(35, 35), image.Pt(8, 1), Pt(1, 0), Pt(14, 34))
	walkLeft := NewSprite(walkLeftImg)
	walkLeft.Divider = 4

	runLeftImg := NewImageSeq(Imgs[0].Img.(*image.NRGBA), Pt(2, 246),
		Pt(35, 35), image.Pt(8, 1), Pt(1, 0), Pt(14, 34))
	runLeft := NewSprite(runLeftImg)
	runLeft.Divider = 2

	// Fire objects
	fireBlastImg := NewImageSeq(Imgs[3].Img.(*image.NRGBA), Pt(1, 37),
		Pt(16, 48), image.Pt(13, 1), Pt(1, 0), Pt(8, 48))
	fireBlast := NewSprite(fireBlastImg)
	fireBlast.Divider = 12

	fireSwirlImg := NewImageSeq(Imgs[3].Img.(*image.NRGBA), Pt(1, 86),
		Pt(8, 8), image.Pt(8, 1), Pt(1, 0), Pt(4, 4))
	fireSwirl := NewSprite(fireSwirlImg)
	fireSwirl.Divider = 2

	// Now to the objects themselve
	mario := NewObject2(waitRight, Point{100, 100})

	a.tt.actions = []TimeTableEntry{
		{ts: 0 * time.Second, fnc: func() {
			mario.SetSprite(waitRight)
		}},
		{ts: 5 * time.Second, fnc: func() {
			mario.SetSprite(walkRight)
		}},
		{ts: 10 * time.Second, fnc: func() {
			mario.SetSprite(runRight)
		}},
		{ts: 15 * time.Second, fnc: func() {
			mario.SetSprite(waitLeft)
		}},
		{ts: 20 * time.Second, fnc: func() {
			mario.SetSprite(walkLeft)
		}},
		{ts: 25 * time.Second, fnc: func() {
			mario.SetSprite(runLeft)
		}},
	}
	a.tt.dur = 30 * time.Second

	for i := range 5 {
		fireObj := NewObject2(fireBlast, Point{float64(150 + i*72), 100.0})
		a.canvas.Add(fireObj)
	}

	fire2 := NewObject2(fireSwirl, Point{50, 100})
	a.canvas.Add(mario, fire2)

	return a
}

func (a *SpriteWin) Init() {
	for i, _ := range a.tt.actions {
		a.tt.actions[i].done = false
	}
	a.tt.t0 = time.Now()
}

func (a *SpriteWin) Update(t time.Time) {
	if time.Since(a.tt.t0) > a.tt.dur {
		for i, _ := range a.tt.actions {
			a.tt.actions[i].done = false
		}
		a.tt.t0 = time.Now()
	} else {
		for i, _ := range a.tt.actions {
			if a.tt.actions[i].done {
				continue
			}
			if time.Since(a.tt.t0) > a.tt.actions[i].ts {
				a.tt.actions[i].done = true
				a.tt.actions[i].fnc()
			}
		}
	}
	a.Root.Update(t)
}

func (a *SpriteWin) Refresh() {
	a.gc.Clear(colors.RGBA{0, 188, 188, 0xFF})
	a.Root.Draw(a.gc)
}
