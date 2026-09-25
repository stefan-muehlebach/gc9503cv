package main

import (
	"image"
	"image/draw"
	"math"
	"sync"
	"time"

	"golang.org/x/image/math/fixed"

	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gc9503cv/iliimg"
	"github.com/stefan-muehlebach/gc9503cv/props"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
)

//-----------------------------------------------------------------------------

type Node interface {
	EventHandler
	Wrappee() *nodeEmbed
	Bounds() Rectangle
	Rect() Rectangle
	Pos() Point
	SetPos(pos Point)
	MinSize() Point
	SetMinSize(size Point)
	Size() Point
	SetSize(size Point)
	FindTarget(pt Point) (Node, Point)
	IsVisible() bool
	Update(dt time.Duration)
	Draw(gc *gg.Context)
	OnInputEvent(ev InputEvent)
}

//-----------------------------------------------------------------------------

type nodeEmbed struct {
	eventHandlerEmbed
	props.PropertyEmbed
	wrapper            Node
	parent             *containerEmbed
	pos, size, minSize Point
	isHidden           bool
	isDisabled         bool
}

func (n *nodeEmbed) Init(node Node) {
	n.wrapper = node
}

func (n *nodeEmbed) Wrappee() *nodeEmbed {
	return n
}

func (n *nodeEmbed) Bounds() Rectangle {
	return Rectangle{Min: n.pos, Max: n.pos.Add(n.Size())}
}

func (n *nodeEmbed) Rect() Rectangle {
	return Rectangle{Max: n.Size()}
}

func (n *nodeEmbed) Pos() Point {
	return n.pos
}

func (n *nodeEmbed) SetPos(pos Point) {
	n.pos = pos
}

func (n *nodeEmbed) MinSize() Point {
	return n.minSize
}

func (n *nodeEmbed) SetMinSize(size Point) {
	n.minSize = size
}

func (n *nodeEmbed) Size() Point {
	return n.size.Max(n.wrapper.MinSize())
}

func (n *nodeEmbed) SetSize(size Point) {
	n.size = size
}

func (n *nodeEmbed) FindTarget(pt Point) (Node, Point) {
	if pt.In(n.Bounds()) {
		return n.wrapper, pt
	} else {
		return nil, Point{}
	}
}

func (n *nodeEmbed) IsVisible() bool {
	return !n.isHidden
}
func (n *nodeEmbed) SetVisible(v bool) {
	n.isHidden = !v
}

func (n *nodeEmbed) IsEnabled() bool {
	return !n.isDisabled
}
func (n *nodeEmbed) SetEnabled(e bool) {
	n.isDisabled = !e
}

func (n *nodeEmbed) Update(dt time.Duration) {}

func (n *nodeEmbed) Draw(gc *gg.Context) {
	n.wrapper.Draw(gc)
}

//-----------------------------------------------------------------------------

type Window interface {
	Bounds() image.Rectangle
	Image() *image.RGBA
	SetMouse(mouse *Mouse)
	Run()
	Stop()
	Init()
	Update(dt time.Duration)
	Refresh()
}

//-----------------------------------------------------------------------------

type windowEmbed struct {
	gc         *gg.Context
	bounds     geom.Rectangle[int]
	isRunning  bool
	mouse      *Mouse
	Root       Container
	ActiveNode Node
}

func (a *windowEmbed) Bounds() image.Rectangle {
	return a.bounds.ToInt()
}

func (a *windowEmbed) Image() *image.RGBA {
	return a.gc.Image().(*image.RGBA)
}

func (w *windowEmbed) SetMouse(mouse *Mouse) {
	w.mouse = mouse
}

func (w *windowEmbed) Run() {
	var node Node

	w.isRunning = true
	for ev := range w.mouse.EventQ {
		if !w.isRunning {
			break
		}
		if w.Root == nil {
			continue
		}
		node, ev.Pos = w.Root.FindTarget(ev.Pos)
		switch ev.Type {
		case MoveEvent, DragEvent:
			if w.ActiveNode != node {
				if w.ActiveNode != nil {
					evNew := ev
					evNew.Type = LeaveEvent
					w.ActiveNode.OnInputEvent(evNew)
				}
				if node != nil {
					evNew := ev
					evNew.Type = EnterEvent
					node.OnInputEvent(evNew)
				}
				w.ActiveNode = node
			}
			if w.ActiveNode != nil {
				w.ActiveNode.OnInputEvent(ev)
			}

		default:
			if node != nil {
				node.OnInputEvent(ev)
			}
		}
	}
}

func (w *windowEmbed) Stop() {
	w.isRunning = false
}

func (w *windowEmbed) Init() {}

func (w *windowEmbed) Update(dt time.Duration) {}

func (w *windowEmbed) Refresh() {
	w.gc.Clear(colors.Black)
	w.Root.Draw(w.gc)
}

//-----------------------------------------------------------------------------

type Application struct {
	disp      *GC9503CV
	win       Window
	Timer     *Stopwatch
	mainImg   *image.RGBA
	pixBuf    *iliimg.ILIImage
	mouse     *Mouse
	isRunning bool
	wg        sync.WaitGroup
}

func NewApplication(disp *GC9503CV) *Application {
	a := &Application{}

	a.disp = disp
	a.Timer = NewStopwatch("Rendering Pipeline")
	a.Timer.SetLapNames(
		"Animate",
		"Refresh",
		"Compose",
		"Convert",
		"Send",
	)
	a.mainImg = image.NewRGBA(disp.DrawBounds().ToInt())
	a.pixBuf = iliimg.NewILIImage(disp.DispBounds().ToInt())
	a.pixBuf.SetLSBFirst()
	a.mouse = OpenMouse()
	a.mouse.SetPosRange(disp.DrawRect())
	a.mouse.SetWheelRange(0, 0, 100)
	a.mouse.SetCursor(LeftPtrCursor)

	return a
}

func (a *Application) SetWindow(win Window) {
	a.win = win
	a.win.SetMouse(a.mouse)
}

func (a *Application) Window() Window {
	return a.win
}

func (a *Application) Run() {
	a.isRunning = true
	a.wg.Go(a.drawThread)
	if a.win != nil {
		a.wg.Go(a.win.Run)
	}
	a.mouse.Start()
	a.wg.Wait()
}

func (a *Application) Stop() {
	a.mouse.Stop()
	a.win.Stop()
	a.isRunning = false
}

func (a *Application) drawThread() {
	dt := 30 * time.Millisecond
	ticker := time.NewTicker(dt)
	defer ticker.Stop()
	a.win.Init()
	for range ticker.C {
		if !a.isRunning {
			break
		}
		a.Timer.Start()
		a.win.Update(dt)
		a.Timer.Lap()

		a.win.Refresh()
		a.Timer.Lap()

		draw.Draw(a.mainImg, a.win.Bounds(), a.win.Image(),
			image.Point{}, draw.Src)
		draw.Draw(a.mainImg, a.mouse.Bounds().Add(a.mainImg.Rect.Min),
			a.mouse.Image(), image.Point{}, draw.Over)
		a.Timer.Lap()

		a.pixBuf.Convert(a.mainImg, a.disp.rot)
		a.Timer.Lap()

		a.disp.Send(a.pixBuf)
		a.Timer.Stop()
	}

}

/*
func (a *Application) eventThread() {
	for mev := range a.mouse.EventQ {
		if !a.isRunning {
			break
		}
		a.win.OnInputEvent(mev)
	}
	a.mouse.Stop()
}
*/

//-----------------------------------------------------------------------------

func flt2fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

func fix2flt(x fixed.Int26_6) float64 {
	const shift, mask = 6, 1<<6 - 1
	if x >= 0 {
		return float64(x>>shift) + float64(x&mask)/64
	}
	x = -x
	if x >= 0 {
		return -(float64(x>>shift) + float64(x&mask)/64)
	}
	return 0
}
