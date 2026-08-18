package main

import (
	"image"
	"image/draw"
	"log"
	"time"
	"gc9503cv/gc9503cv/geom"
	"gc9503cv/gc9503cv/iliimg"
	"github.com/stefan-muehlebach/gg"
)

//----------------------------------------------------------------------------

func IsSet(value, pattern MouseButtonType) bool {
	if value&pattern != 0x00 {
		return true
	} else {
		return false
	}
}

//-----------------------------------------------------------------------------

type Applet interface {
	Bounds() image.Rectangle
	Image() *image.RGBA
	Init()
	Update(dt time.Duration)
	Refresh()
	OnInputEvent(ev MouseEvent)
}

//-----------------------------------------------------------------------------

type appletEmbed struct {
	gc     *gg.Context
	bounds geom.Rectangle[int]
}

func (a *appletEmbed) Bounds() image.Rectangle {
	return a.bounds.ToInt()
}

func (a *appletEmbed) Image() *image.RGBA {
	return a.gc.Image().(*image.RGBA)
}

//-----------------------------------------------------------------------------

type Application struct {
	disp                                       *GC9503CV
	applet                                     Applet
	AnimWatch, DrawWatch, ConvWatch, SendWatch *Stopwatch
	mainImg                                    *image.RGBA
	pixBuf                                     *iliimg.ILIImage
	mouse                                      *Mouse
	isRunning                                  bool
}

func NewApplication(disp *GC9503CV) *Application {
	a := &Application{}

	a.disp = disp
	a.AnimWatch = NewStopwatch()
	a.DrawWatch = NewStopwatch()
	a.ConvWatch = NewStopwatch()
	a.SendWatch = NewStopwatch()
	a.mainImg = image.NewRGBA(disp.DrawBounds().ToInt())
	a.pixBuf = iliimg.NewILIImage(disp.DispBounds().ToInt())
	a.pixBuf.SetLSBFirst()
	a.mouse = OpenMouse()
	a.mouse.SetPosRange(disp.DrawRect())
	a.mouse.SetWheelRange(25, 0, 50)
	a.mouse.SetCursor(LeftPtrCursor)

	return a
}

func (a *Application) PrintWatchStats() {
	animAvg := a.AnimWatch.Avg()
	drawAvg := a.DrawWatch.Avg()
	convAvg := a.ConvWatch.Avg()
	sendAvg := a.SendWatch.Avg()

	log.Printf("animation : %v", animAvg)
	log.Printf("drawing   : %v", drawAvg)
	log.Printf("converting: %v", convAvg)
	log.Printf("sending   : %v", sendAvg)
	log.Printf("total     : %v", animAvg+drawAvg+convAvg+sendAvg)

	a.AnimWatch.Reset()
	a.DrawWatch.Reset()
	a.ConvWatch.Reset()
	a.SendWatch.Reset()
}

func (a *Application) SetApplet(applet Applet) {
	a.applet = applet
}

func (a *Application) drawThread() {
	dt := 30 * time.Millisecond
	ticker := time.NewTicker(dt)
	defer ticker.Stop()
	a.applet.Init()
	for range ticker.C {
		if !a.isRunning {
			break
		}
		a.AnimWatch.Start()
		a.applet.Update(dt)
		a.AnimWatch.Stop()

		a.DrawWatch.Start()
		a.applet.Refresh()
		draw.Draw(a.mainImg, a.applet.Bounds(), a.applet.Image(),
			image.Point{}, draw.Over)
		draw.Draw(a.mainImg, a.mouse.Bounds().Add(a.mainImg.Rect.Min), a.mouse.Image(),
			image.Point{}, draw.Over)
		a.DrawWatch.Stop()

		a.ConvWatch.Start()
		a.pixBuf.Convert(a.mainImg, a.disp.rot)
		a.ConvWatch.Stop()

		a.SendWatch.Start()
		a.disp.Send(a.pixBuf)
		a.SendWatch.Stop()
	}
}

func (a *Application) eventThread() {
	for mev := range a.mouse.EventQ {
		if !a.isRunning {
			break
		}
		a.applet.OnInputEvent(mev)
	}
	a.mouse.Close()
}

func (a *Application) Run(timeout time.Duration) {
	a.isRunning = true
	go a.drawThread()
	go a.eventThread()
	a.mouse.StartEvents()
	time.Sleep(timeout)
	a.isRunning = false
}
