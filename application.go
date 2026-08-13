package main

import (
	"image"
	"log"
	"time"
	"github.com/stefan-muehlebach/gg"
	"gc9503cv/gc9503cv/iliimg"
)

//-----------------------------------------------------------------------------

type Prog interface {
    Init(gc *gg.Context)
    Update(dt time.Duration)
	Draw(gc *gg.Context)
	OnInputEvent(ev MouseEvent)
}

//-----------------------------------------------------------------------------

type Application struct {
	disp *GC9503CV
	canv *Canvas
	prog Prog
	AnimWatch, DrawWatch, ConvWatch, SendWatch *Stopwatch
	pixBuf *iliimg.ILIImage
	mouse *Mouse
	isRunning bool
}

func NewApplication(disp *GC9503CV) (*Application) {
	a := &Application{}

	a.disp = disp
	a.canv = disp.Canvas()
	a.AnimWatch = NewStopwatch()
	a.DrawWatch = NewStopwatch()
	a.ConvWatch = NewStopwatch()
	a.SendWatch = NewStopwatch()
	a.pixBuf = iliimg.NewILIImage(disp.DispBounds().ToInt())
	a.pixBuf.SetLSBFirst()
	a.mouse = OpenMouse()
	a.mouse.SetPosRange(disp.DrawRect())
	a.mouse.SetWheelRange(25, 0, 50)
	a.mouse.SetCursor(LeftPtrCursor)

	return a
}

func (a *Application) Canvas() *Canvas {
	return a.canv
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

func (a *Application) AddProg(prog Prog) {
	a.prog = prog
}

func (a *Application) drawThread() {
	dt := 30 * time.Millisecond
	ticker := time.NewTicker(dt)
	//subImgRect := a.disp.DispBounds().Inset(150, 300).ToInt()
	defer ticker.Stop()
	if a.prog != nil {
		a.prog.Init(a.canv.GC)
	}
	for range ticker.C {
		if !a.isRunning {
			break
		}
		a.AnimWatch.Start()
		if a.prog != nil {
			a.prog.Update(dt)
		} else {
			// a.canv.Update(dt)
		}
		a.AnimWatch.Stop()
		a.DrawWatch.Start()
		//a.canv.Clear(a.canv.BackColor)
		if a.prog != nil {
			a.prog.Draw(a.canv.GC)
		} else {
			a.canv.Refresh()
		}
		a.mouse.Draw(a.canv.GC)
		a.DrawWatch.Stop()
		a.ConvWatch.Start()
		rgba := a.canv.Img.(*image.RGBA)
		a.pixBuf.Convert(rgba)
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
		a.mouse.Pos = mev.Pos
		a.prog.OnInputEvent(mev)
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

