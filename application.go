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
	Root   Container
}

func (a *appletEmbed) Bounds() image.Rectangle {
	return a.bounds.ToInt()
}

func (a *appletEmbed) Image() *image.RGBA {
	return a.gc.Image().(*image.RGBA)
}

//-----------------------------------------------------------------------------

type Node interface {
	Wrappee() *nodeEmbed
	Bounds() Rectangle
	Pos() Point
	Size() Point
	SetPos(pos Point)
	SetSize(size Point)
	Contains(pt Point) Node
	OnInputEvent(ev MouseEvent)
	Update(dt time.Duration)
	Draw(gc *gg.Context)
}

//-----------------------------------------------------------------------------

type nodeEmbed struct {
	eventHandlerEmbed
	wrapper Node
	parent Container
	pos, size Point
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

func (n *nodeEmbed) Size() Point {
	return n.size
}

func (n *nodeEmbed) SetPos(pos Point) {
	n.pos = pos
}
	
func (n *nodeEmbed) SetSize(size Point) {
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

func (c *containerEmbed) Update(dt time.Duration) {
    for e := c.childList.Front(); e != nil; e = e.Next() {
        e.Value.(Node).Update(dt)
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

func (c *containerEmbed) Draw(gc *gg.Context) {
	//b := c.Bounds()
	//gc.DrawRectangle(b.Min.X, b.Min.Y, b.Dx(), b.Dy())
	//gc.Clip()
    for e := c.childList.Front(); e != nil; e = e.Next() {
        if n, ok := e.Value.(Node); ok {
        	n.Draw(gc)
		}
    }
	//gc.ResetClip()
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

//-----------------------------------------------------------------------------

type Application struct {
	disp                                       *GC9503CV
	applet                                     Applet
	Timer                                      *Stopwatch
	mainImg                                    *image.RGBA
	pixBuf                                     *iliimg.ILIImage
	mouse                                      *Mouse
	isRunning                                  bool
	wg                      				   sync.WaitGroup
}

func NewApplication(disp *GC9503CV) *Application {
	a := &Application{}

	a.disp = disp
	a.Timer = NewStopwatch("New Tiner")
	a.mainImg = image.NewRGBA(disp.DrawBounds().ToInt())
	a.pixBuf = iliimg.NewILIImage(disp.DispBounds().ToInt())
	a.pixBuf.SetLSBFirst()
	a.mouse = OpenMouse()
	a.mouse.SetPosRange(disp.DrawRect())
	a.mouse.SetWheelRange(0, 0, 100)
	a.mouse.SetCursor(LeftPtrCursor)

	return a
}

func (a *Application) SetApplet(applet Applet) {
	a.applet = applet
}

func (a *Application) Run() {
	a.isRunning = true
	a.wg.Go(a.drawThread)
	a.wg.Go(a.eventThread)
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
	a.applet.Init()
	for range ticker.C {
		if !a.isRunning {
			break
		}
		a.Timer.Start()
		a.applet.Update(dt)
		a.Timer.Lap()

		a.applet.Refresh()
		a.Timer.Lap()
		
		draw.Draw(a.mainImg, a.applet.Bounds(), a.applet.Image(),
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
		a.applet.OnInputEvent(mev)
	}
	a.mouse.Close()
}


