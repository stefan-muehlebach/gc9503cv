package main

import (
	"container/list"
	"image"
	"image/draw"
	"sync"
	"time"

    "golang.org/x/image/font"

	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gc9503cv/iliimg"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
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
	FindTarget(pt Point) Node
	IsVisible() bool
	Update(dt time.Duration)
	Draw(gc *gg.Context)
	OnInputEvent(ev InputEvent)
}

//-----------------------------------------------------------------------------

type nodeEmbed struct {
	eventHandlerEmbed
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
	return Rectangle{Min: n.pos, Max: n.pos.Add(n.size)}
}

func (n *nodeEmbed) Rect() Rectangle {
	return Rectangle{Max: n.size}
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

func (n *nodeEmbed) FindTarget(pt Point) Node {
	if pt.In(n.Bounds()) {
		return n.wrapper
	} else {
		return nil
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

const (
	cornerRadius         =  6.0
	minWidth             = 36.0
	minHeight            = 36.0

	borderWidth          =  0.0
	lineWidth            =  2.5
	pushedBorderWidth    =  0.0
	selectedBorderWidth  =  3.0

	padding              =  5.0
	innerPadding         =  5.0

	buttonFontSize       = 12.0
	radioboxFontSize     = 12.0

	radioboxWidth        = 20.0
	radioboxHeight       = 20.0
	radioboxInnerPadding =  5.0
	radioboxCornerRadius =  5.0
	radiobuttonLineWidth =  8.0
)

var (
	fillColor           = colors.GoTeal
	borderColor         = colors.GoTeal
	textColor           = colors.Gainsboro
	lineColor           = colors.Gainsboro

	pushedColor         = colors.GoTeal.Bright(0.6)
	pushedBorderColor   = colors.GoTeal.Bright(0.7)
	pushedTextColor     = colors.Gainsboro
	pushedLineColor     = colors.Gainsboro

	selectedColor       = colors.GoTeal
	selectedBorderColor = colors.GoTeal.Bright(0.6)
	selectedTextColor   = colors.Gainsboro
	selectedLineColor   = colors.Gainsboro

	buttonFontName      = fonts.GoBold
	buttonFontFace, _   = fonts.NewFace(buttonFontName, buttonFontSize)

	radioboxFontName    = fonts.GoRegular
	radioboxFontFace, _ = fonts.NewFace(radioboxFontName, radioboxFontSize)
)

type Button struct {
	nodeEmbed
	pushed  bool
	checked bool
}

func NewButton(w, h float64) *Button {
	b := &Button{}
	b.Init(b)
	b.SetMinSize(Point{w, h})
	return b
}

func (b *Button) Draw(gc *gg.Context) {
	gc.DrawRoundedRectangle(b.pos.X, b.pos.Y, b.size.X, b.size.Y, cornerRadius)
	if b.pushed {
		gc.SetFillColor(pushedColor)
		gc.SetLineColor(pushedBorderColor)
		gc.SetLineWidth(pushedBorderWidth)
	} else {
		if b.checked {
			gc.SetFillColor(selectedColor)
			gc.SetLineColor(selectedBorderColor)
			gc.SetLineWidth(selectedBorderWidth)
		} else {
			gc.SetFillColor(fillColor)
			gc.SetLineColor(borderColor)
			gc.SetLineWidth(borderWidth)
		}
	}
	gc.FillStroke()
}

func (b *Button) OnInputEvent(ev InputEvent) {
	if ev.Type == PressEvent {
		b.pushed = true
	}
	if ev.Type == ReleaseEvent || ev.Type == LeaveEvent {
		b.pushed = false
	}
	b.CallEventHandler(ev)
}

//-----------------------------------------------------------------------------

// Der Typ AlignType dient der Ausrichtung von Text.
type AlignType int

const (
    AlignLeft AlignType = 1 << iota
    AlignCenter
    AlignRight
    AlignTop
    AlignMiddle
    AlignBottom

    horizontalAlignMask = (AlignLeft | AlignCenter | AlignRight)
    verticalAlignMask   = (AlignTop | AlignMiddle | AlignBottom )
)

type TextButton struct {
	Button
	label string
	align AlignType
	ax, ay float64
}

func NewTextButton(label string) *TextButton {
	b := &TextButton{}
	b.Init(b)
	b.SetMinSize(Point{minWidth, minHeight})
	b.label = label
	b.align = AlignCenter | AlignMiddle
	b.ax, b.ay = 0.5, 0.5
	return b
}

func (b *TextButton) Draw(gc *gg.Context) {
	b.Button.Draw(gc)
	gc.SetFontFace(buttonFontFace)
	if b.pushed {
		gc.SetTextColor(pushedTextColor)
	} else {
		if b.checked {
			gc.SetTextColor(selectedTextColor)
		} else {
			gc.SetTextColor(textColor)
		}
	}
	mp := b.Bounds().Center()
	gc.DrawStringAnchored(b.label, mp.X, mp.Y, b.ax, b.ay)
}

//-----------------------------------------------------------------------------

type IconButton struct {
	Button
	img image.Image
}

func NewIconButton(imgFile string) *IconButton {
	b := &IconButton{}
	b.Init(b)
	b.img, _ = gg.LoadPNG(imgFile)
	size := b.img.Bounds().Inset(-innerPadding).Size()
	b.SetMinSize(Point{float64(size.X), float64(size.Y)})
	return b
}

func (b *IconButton) Draw(gc *gg.Context) {
	b.Button.Draw(gc)
	mp := b.Bounds().Center()
	gc.DrawImageAnchored(b.img, mp.X, mp.Y, 0.5, 0.5)
}

func (b *IconButton) OnInputEvent(ev InputEvent) {
	b.Button.OnInputEvent(ev)
	if ev.Type == ClickEvent {
		b.checked = !b.checked
	}
}

//-----------------------------------------------------------------------------

type RadioboxType int

const (
	Checkbox RadioboxType = iota
	Radiobutton
)

type Radiobox struct {
	Button
	label string
	typ RadioboxType
}

func NewCheckbox(label string) *Radiobox {
	b := newRadiobox(label)
	b.typ = Checkbox
	return b
}

func NewRadiobutton(label string) *Radiobox {
	b := newRadiobox(label)
	b.typ = Radiobutton
	return b
}

func newRadiobox(label string) *Radiobox {
	b := &Radiobox{}
	b.Init(b)
	b.label = label
	b.typ = Radiobutton
	w := float64(font.MeasureString(radioboxFontFace, b.label))/64.0
	b.SetMinSize(Point{radioboxWidth+radioboxInnerPadding+w, radioboxHeight})
	return b
}	

func (b *Radiobox) Draw(gc *gg.Context) {
	var mp Point

	if b.typ == Checkbox {
		gc.DrawRoundedRectangle(b.pos.X, b.pos.Y, radioboxWidth,
			radioboxHeight, radioboxCornerRadius)
	} else {
		mp = Point{b.pos.X+0.5*radioboxWidth, b.pos.Y+0.5*radioboxHeight}
		gc.DrawCircle(mp.X, mp.Y, 0.5*radioboxWidth)
	}
	if b.pushed {
		gc.SetFillColor(pushedColor)
		gc.SetLineColor(pushedBorderColor)
	} else {
		gc.SetFillColor(fillColor)
		gc.SetLineColor(borderColor)
	}
	gc.SetLineWidth(borderWidth)
	gc.FillStroke()
	if b.checked {
		if b.typ == Checkbox {
			gc.SetLineWidth(lineWidth)
			if b.pushed {
				gc.SetLineColor(pushedLineColor)
			} else {
				gc.SetLineColor(lineColor)
			}
			gc.MoveTo(b.pos.X+4, b.pos.Y+9)
			gc.LineTo(b.pos.X+8, b.pos.Y+14)
			gc.LineTo(b.pos.X+14, b.pos.Y+5)
			gc.Stroke()
		} else {
			if b.pushed {
				gc.SetFillColor(pushedLineColor)
			} else {
				gc.SetFillColor(lineColor)
			}
			gc.DrawCircle(mp.X, mp.Y, 0.5*radiobuttonLineWidth)
			gc.Fill()
		}
	}
	x := b.pos.X + radioboxWidth + radioboxInnerPadding
	y := b.pos.Y + 0.5 * radioboxHeight
	gc.SetTextColor(textColor)
	gc.SetFontFace(radioboxFontFace)
	gc.DrawStringAnchored(b.label, x, y, 0.0, 0.5)
}

func (b *Radiobox) OnInputEvent(ev InputEvent) {
	b.Button.OnInputEvent(ev)
	if ev.Type == ClickEvent {
		b.checked = !b.checked
	}
}

//-----------------------------------------------------------------------------

type Container interface {
	Node
	Add(nl ...Node)
	Del(n Node)
	FindTarget(pt Point) Node
	Update(dt time.Duration)
	Draw(gc *gg.Context)
	SetLayoutManager(layoutManager LayoutManager)
	LayoutManager() LayoutManager
	layout()
}

//-----------------------------------------------------------------------------

type containerEmbed struct {
	nodeEmbed
	childList *list.List
	Layout    LayoutManager
}

func (c *containerEmbed) Init(n Node) {
	c.nodeEmbed.Init(n)
	c.childList = list.New()
}

func (c *containerEmbed) SetSize(size Point) {
	c.nodeEmbed.SetSize(size)
	c.layout()
}

func (c *containerEmbed) MinSize() Point {
	if c.minSize.Eq(Point{0, 0}) {
		return c.Layout.MinSize(c.childList)
	} else {
		return c.nodeEmbed.MinSize()
	}
}

func (c *containerEmbed) Add(nl ...Node) {
	for _, n := range nl {
		embed := n.Wrappee()
		embed.parent = c
		c.childList.PushBack(n)
		c.layout()
	}
}

func (c *containerEmbed) Del(n Node) {
	for e := c.childList.Front(); e != nil; e = e.Next() {
		if e.Value.(Node) == n {
			c.childList.Remove(e)
			break
		}
	}
	c.layout()
}

func (c *containerEmbed) FindTarget(pt Point) Node {
	if n := c.nodeEmbed.FindTarget(pt); n == nil {
		return nil
	}
	pt = pt.Sub(c.pos)
	for e := c.childList.Back(); e != nil; e = e.Prev() {
		if n := e.Value.(Node).FindTarget(pt); n != nil {
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
	gc.Push()
	gc.Translate(c.pos.AsCoord())
	for e := c.childList.Front(); e != nil; e = e.Next() {
		e.Value.(Node).Draw(gc)
	}
	gc.Pop()
}

func (c *containerEmbed) SetLayoutManager(layoutManager LayoutManager) {
	c.Layout = layoutManager
}

func (c *containerEmbed) LayoutManager() LayoutManager {
	return c.Layout
}

func (c *containerEmbed) layout() {
	if c.Layout == nil {
		return
	}
	c.Layout.Layout(c.childList, c.wrapper.Size())
}

//-----------------------------------------------------------------------------

type LayoutManager interface {
	Layout(childList *list.List, size Point)
	MinSize(childList *list.List) Point
}

//-----------------------------------------------------------------------------

type NullLayout struct{}

func (l *NullLayout) Layout(childList *list.List, size Point) {}

func (l *NullLayout) MinSize(childList *list.List) Point {
	minSize := Point{}
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		if !n.IsVisible() {
			continue
		}
		minSize = minSize.Max(n.Bounds().Max)
	}
	return minSize
}

//-----------------------------------------------------------------------------

type PadLayout struct {
	pad [4]float64
}

func NewPadLayout(pads ...float64) *PadLayout {
	var l, t, r, b float64

	switch len(pads) {
	case 0:
		l, t, r, b = padding, padding, padding, padding
	case 1:
		l, t, r, b = pads[0], pads[0], pads[0], pads[0]
	case 2:
		l, t, r, b = pads[0], pads[1], pads[0], pads[1]
	case 3:
		l, t, r, b = pads[0], pads[1], pads[2], pads[1]
	case 4:
		l, t, r, b = pads[0], pads[1], pads[2], pads[3]
	}
	return &PadLayout{[4]float64{l, t, r, b}}
}

func (l *PadLayout) Layout(childList *list.List, size Point) {
	pos := Point{l.pad[0], l.pad[1]}
	siz := Point{size.X - l.pad[0] - l.pad[2],
		size.Y - l.pad[1] - l.pad[3]}
	for e := childList.Front(); e != nil; e = e.Next() {
		node := e.Value.(Node)
		node.SetSize(siz)
		node.SetPos(pos)
	}
}

func (l *PadLayout) MinSize(childList *list.List) Point {
	minSize := Point{}
	for e := childList.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		if !n.IsVisible() {
			continue
		}
		minSize = minSize.Max(n.MinSize())
	}
	return minSize.Add(Point{l.pad[0] + l.pad[2], l.pad[1] + l.pad[3]})
}

//-----------------------------------------------------------------------------

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

type BoxLayout struct {
	orient Orientation
	pad    float64
}

func NewHBoxLayout(pads ...float64) *BoxLayout {
	pad := innerPadding
	if len(pads) > 0 {
		pad = pads[0]
	}
	return &BoxLayout{Horizontal, pad}
}

func NewVBoxLayout(pads ...float64) *BoxLayout {
	pad := innerPadding
	if len(pads) > 0 {
		pad = pads[0]
	}
	return &BoxLayout{Vertical, pad}
}

func (l *BoxLayout) isSpacer(obj Node) bool {
	spc, ok := obj.(*Spacer)
	if !ok {
		return false
	}
	if l.orient == Horizontal {
		return spc.ExpandHorizontal()
	}
	return spc.ExpandVertical()
}

func (l *BoxLayout) Layout(childList *list.List, size Point) {
	spacers := 0
	total := 0.0
	for e := childList.Front(); e != nil; e = e.Next() {
		node := e.Value.(Node)
		if !node.IsVisible() {
			continue
		}
		if l.isSpacer(node) {
			spacers++
			continue
		}
		nodeSize := node.MinSize()
		switch l.orient {
		case Horizontal:
			total += nodeSize.X
		case Vertical:
			total += nodeSize.Y
		}
	}
	extra, extraCell := 0.0, 0.0
	switch l.orient {
	case Horizontal:
		extra = size.X - total - l.pad*float64(childList.Len()-
			spacers-1)
	case Vertical:
		extra = size.Y - total - l.pad*float64(childList.Len()-
			spacers-1)
	}
	if spacers > 0 {
		extraCell = extra / float64(spacers)
	}
	pos := Point{}
	for e := childList.Front(); e != nil; e = e.Next() {
		node := e.Value.(Node)
		if !node.IsVisible() {
			continue
		}
		if l.isSpacer(node) {
			switch l.orient {
			case Horizontal:
				pos.X += extraCell
			case Vertical:
				pos.Y += extraCell
			}
			continue
		}
		node.SetPos(pos)
		switch l.orient {
		case Horizontal:
			width := node.MinSize().X
			pos.X += width + l.pad
			node.SetSize(Point{width, size.Y})
		case Vertical:
			height := node.MinSize().Y
			pos.Y += height + l.pad
			node.SetSize(Point{size.X, height})
		}
	}
}

func (l *BoxLayout) MinSize(childList *list.List) Point {
	//log.Printf("BoxLayout.MinSize()")
	minSize := Point{}
	addPadding := false
	for e := childList.Front(); e != nil; e = e.Next() {
		node := e.Value.(Node)
		if !node.IsVisible() || l.isSpacer(node) {
			continue
		}
		childSize := node.MinSize()
		switch l.orient {
		case Horizontal:
			minSize.Y = max(minSize.Y, childSize.Y)
			minSize.X += childSize.X
			if addPadding {
				minSize.X += l.pad
			}
		case Vertical:
			minSize.X = max(minSize.X, childSize.X)
			minSize.Y += childSize.Y
			if addPadding {
				minSize.Y += l.pad
			}
		}
		addPadding = true
	}
	return minSize
}

//-----------------------------------------------------------------------------

// Nimmt den verfügbaren Platz (vertikal oder horizontal) in Box-Layouts
// ein. Ist zwar ein Widget, passt aber irgendwie besser zum Layout-Zeugs.
type Spacer struct {
	nodeEmbed
	FixHorizontal, FixVertical bool
}

func NewSpacer() *Spacer {
	s := &Spacer{}
	s.Init(s)
	return s
}

func (s *Spacer) Draw(gc *gg.Context) {}

func (s *Spacer) ExpandHorizontal() bool {
	return !s.FixHorizontal
}

func (s *Spacer) ExpandVertical() bool {
	return !s.FixVertical
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

type Panel struct {
	containerEmbed
	bgColor colors.RGBA
}

func NewPanel(color colors.RGBA) *Panel {
	p := &Panel{}
	p.Init(p)
	p.bgColor = color
	return p
}

func (p *Panel) Draw(gc *gg.Context) {
	gc.DrawRectangle(p.Bounds().AsCoord())
	gc.SetFillColor(p.bgColor)
	gc.Fill()
	p.containerEmbed.Draw(gc)
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
	w.isRunning = true
	for ev := range w.mouse.EventQ {
		if !w.isRunning {
			break
		}
		if w.Root == nil {
			continue
		}
		node := w.Root.FindTarget(ev.Pos)
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
	a.Timer = NewStopwatch("New Timer")
	a.mainImg = image.NewRGBA(disp.DrawBounds().ToInt())
	a.pixBuf = iliimg.NewILIImage(disp.DispBounds().ToInt())
	a.pixBuf.SetLSBFirst()
	a.mouse = OpenMouse()
	a.mouse.SetPosRange(disp.DrawRect())
	a.mouse.SetWheelRange(0, 0, 100)
	a.mouse.SetCursor(CrosshairCursor)

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
