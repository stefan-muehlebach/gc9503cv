package main

import (
	"container/list"
	"image"
	"image/draw"
	"sync"
	"time"

	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gc9503cv/iliimg"
	"github.com/stefan-muehlebach/gg"
)

//-----------------------------------------------------------------------------

type Node interface {
	EventHandler
	Wrappee() *nodeEmbed
	Bounds() Rectangle
	Pos() Point
	Move(pos Point)
	MinSize() Point
	Size() Point
	Resize(size Point)
	Contains(pt Point) Node
	Update(dt time.Duration)
	Draw(gc *gg.Context)
	OnInputEvent(ev MouseEvent)
}

//-----------------------------------------------------------------------------

type nodeEmbed struct {
	eventHandlerEmbed
	wrapper            Node
	parent             Container
	pos, size, minSize Point
}

func (n *nodeEmbed) Init(node Node) {
	n.wrapper = node
}

func (n *nodeEmbed) Wrappee() *nodeEmbed {
	return n
}

func (n *nodeEmbed) Bounds() Rectangle {
	return Rectangle{Min: n.pos, Max: n.pos.Add(n.size)}
}

func (n *nodeEmbed) Pos() Point {
	return n.pos
}

func (n *nodeEmbed) Move(pos Point) {
	n.pos = pos
}

func (n *nodeEmbed) MinSize() Point {
	return n.minSize
}

func (n *nodeEmbed) Size() Point {
	return n.size
}

func (n *nodeEmbed) Resize(size Point) {
	n.size = size
}

func (n *nodeEmbed) Contains(pt Point) Node {
	if pt.In(n.Bounds()) {
		return n.wrapper
	} else {
		return nil
	}
}

func (n *nodeEmbed) Draw(gc *gg.Context) {
	n.wrapper.Draw(gc)
}

//-----------------------------------------------------------------------------

type Container interface {
	Node
	Add(nl ...Node)
	Del(n Node)
	Contains(pt Point) Node
	Update(dt time.Duration)
	Draw(gc *gg.Context)
}

//-----------------------------------------------------------------------------

type containerEmbed struct {
	nodeEmbed
	childList *list.List
}

func (c *containerEmbed) Init(n Node) {
	c.nodeEmbed.Init(n)
	c.childList = list.New()
}

func (c *containerEmbed) Add(nl ...Node) {
	for _, n := range nl {
		embed := n.Wrappee()
		embed.parent = c
		c.childList.PushBack(n)
	}
}

func (c *containerEmbed) Del(n Node) {
	for e := c.childList.Front(); e != nil; e = e.Next() {
		if e.Value.(Node) == n {
			c.childList.Remove(e)
			break
		}
	}
}

func (c *containerEmbed) Contains(pt Point) Node {
	if n := c.nodeEmbed.Contains(pt); n == nil {
		return nil
	}
	for e := c.childList.Back(); e != nil; e = e.Prev() {
		if n := e.Value.(Node).Contains(pt); n != nil {
			return n
		}
	}
	return c
}

func (c *containerEmbed) Update(dt time.Duration) {
	for e := c.childList.Front(); e != nil; e = e.Next() {
		e.Value.(Node).Update(dt)
	}
}

func (c *containerEmbed) Draw(gc *gg.Context) {
	for e := c.childList.Front(); e != nil; e = e.Next() {
		e.Value.(Node).Draw(gc)
	}
}

//-----------------------------------------------------------------------------

type Group struct {
	containerEmbed
}

func NewGroup() *Group {
	g := &Group{}
	g.Init(g)
	return g
}

/*
func (g *Group) Add(nl ...Node) {
	rect := g.Bounds()
	for _, n := range nl {
		rect = rect.Union(n.Bounds())
	}
	g.pos = rect.Min
	g.size = rect.Size()
	g.containerEmbed.Add(nl...)
}
*/

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
	a.Timer = NewStopwatch("New Timer")
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
}

func (a *Application) Window() Window {
	return a.win
}

func (a *Application) Run() {
	a.isRunning = true
	a.wg.Go(a.drawThread)
	if a.win != nil {
		a.wg.Go(a.eventThread)
	}
	a.mouse.StartEvents()
	a.wg.Wait()
}

func (a *Application) Stop() {
	a.isRunning = false
}

func (a *Application) drawThread() {
	defer a.wg.Done()
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

func (a *Application) eventThread() {
	defer a.wg.Done()
	for mev := range a.mouse.EventQ {
		if !a.isRunning {
			break
		}
		a.win.OnInputEvent(mev)
	}
	a.mouse.Close()
}

//-----------------------------------------------------------------------------

type Window interface {
	Bounds() image.Rectangle
	Image() *image.RGBA
	Init()
	Update(dt time.Duration)
	Refresh()
	OnInputEvent(ev MouseEvent)
}

//-----------------------------------------------------------------------------

type windowEmbed struct {
	gc     *gg.Context
	bounds geom.Rectangle[int]
	Root   Container
}

func (a *windowEmbed) Bounds() image.Rectangle {
	return a.bounds.ToInt()
}

func (a *windowEmbed) Image() *image.RGBA {
	return a.gc.Image().(*image.RGBA)
}

