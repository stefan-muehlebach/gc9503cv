package main

import (
	"image"
	"math"

	"github.com/stefan-muehlebach/gc9503cv/binding"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/fonts"

	//	"github.com/stefan-muehlebach/gc9503cv/geom"
	"golang.org/x/image/font"
)

//-----------------------------------------------------------------------------

type Label struct {
	nodeEmbed
	label    binding.String
	fontFace font.Face
	align    AlignType
	ax, ay   float64
}

func NewLabel(label string) *Label {
	data := binding.NewString()
	data.Set(label)
	return NewLabelWithData(data)
}

func NewLabelWithData(data binding.String) *Label {
	l := &Label{}
	l.Init(l)
	l.InitProp("Label")
	l.label = data
	l.fontFace, _ = fonts.NewFace(l.RegularFont(), l.FontSize())
	l.align = AlignCenter | AlignMiddle
	l.ax, l.ay = 0.5, 0.5
	w := fix2flt(font.MeasureString(l.fontFace, l.label.Get()))
	h := fix2flt(l.fontFace.Metrics().CapHeight + l.fontFace.Metrics().Descent)
	l.SetMinSize(Point{w + 2*l.InnerPadding(), h + 2*l.InnerPadding()})
	return l
}

func (l *Label) Draw(gc *gg.Context) {
	gc.DrawRectangle(l.Bounds().AsCoord())
	gc.SetLineColor(l.BorderColor())
	gc.SetLineWidth(l.BorderWidth())
	gc.SetFillColor(l.FillColor())
	gc.FillStroke()
	gc.SetFontFace(l.fontFace)
	gc.SetTextColor(l.TextColor())
	mp := l.Bounds().Center()
	gc.DrawStringAnchored(l.label.Get(), mp.X, mp.Y, l.ax, l.ay)
}

//-----------------------------------------------------------------------------

type Button struct {
	nodeEmbed
	pushed  bool
	checked bool
}

func NewButton() *Button {
	b := &Button{}
	b.Init(b)
	b.InitProp("Button")
	b.SetMinSize(Point{b.Width(), b.Height()})
	return b
}

func (b *Button) Draw(gc *gg.Context) {
	gc.DrawRoundedRectangle(b.pos.X, b.pos.Y, b.size.X, b.size.Y,
		b.CornerRadius())
	if b.pushed {
		gc.SetFillColor(b.PushedColor())
		gc.SetLineColor(b.PushedBorderColor())
		gc.SetLineWidth(b.PushedBorderWidth())
	} else {
		if b.checked {
			gc.SetFillColor(b.SelectedColor())
			gc.SetLineColor(b.SelectedBorderColor())
			gc.SetLineWidth(b.SelectedBorderWidth())
		} else {
			gc.SetFillColor(b.FillColor())
			gc.SetLineColor(b.BorderColor())
			gc.SetLineWidth(b.BorderWidth())
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
	verticalAlignMask   = (AlignTop | AlignMiddle | AlignBottom)
)

type TextButton struct {
	Button
	fontFace font.Face
	label    string
	align    AlignType
	ax, ay   float64
}

func NewTextButton(label string) *TextButton {
	b := &TextButton{}
	b.Init(b)
	b.InitProp("TextButton")
	b.fontFace, _ = fonts.NewFace(b.BoldFont(), b.FontSize())
	b.label = label
	b.align = AlignCenter | AlignMiddle
	b.ax, b.ay = 0.5, 0.5
	w := fix2flt(font.MeasureString(b.fontFace, b.label))
	h := b.Height()
	b.SetMinSize(Point{w + 2*b.InnerPadding(), h})
	return b
}

func (b *TextButton) Draw(gc *gg.Context) {
	b.Button.Draw(gc)
	gc.SetFontFace(b.fontFace)
	if b.pushed {
		gc.SetTextColor(b.PushedTextColor())
	} else {
		if b.checked {
			gc.SetTextColor(b.SelectedTextColor())
		} else {
			gc.SetTextColor(b.TextColor())
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
	b.InitProp("IconButton")
	b.img, _ = gg.LoadPNG(imgFile)
	size := b.img.Bounds().Inset(-int(b.InnerPadding())).Size()
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
	fontFace font.Face
	label    string
	typ      RadioboxType
}

func NewCheckbox(label string) *Radiobox {
	b := newRadiobox(label)
	b.InitProp("Checkbox")
	b.typ = Checkbox
	return b
}

func NewRadiobutton(label string) *Radiobox {
	b := newRadiobox(label)
	b.InitProp("RadioButton")
	b.typ = Radiobutton
	return b
}

func newRadiobox(label string) *Radiobox {
	b := &Radiobox{}
	b.Init(b)
	b.InitProp("Radiobox")
	b.fontFace, _ = fonts.NewFace(b.RegularFont(), b.FontSize())
	b.label = label
	b.typ = Radiobutton
	w := fix2flt(font.MeasureString(b.fontFace, b.label))
	b.SetMinSize(Point{b.Width() + b.InnerPadding() + w, b.Height()})
	return b
}

func (b *Radiobox) Draw(gc *gg.Context) {
	var mp Point

	if b.typ == Checkbox {
		gc.DrawRoundedRectangle(b.pos.X, b.pos.Y, b.Width(),
			b.Height(), b.CornerRadius())
	} else {
		mp = Point{b.pos.X + 0.5*b.Width(), b.pos.Y + 0.5*b.Height()}
		gc.DrawCircle(mp.X, mp.Y, 0.5*b.Width())
	}
	if b.pushed {
		gc.SetFillColor(b.PushedColor())
		gc.SetLineColor(b.PushedBorderColor())
	} else {
		gc.SetFillColor(b.FillColor())
		gc.SetLineColor(b.BorderColor())
	}
	gc.SetLineWidth(b.BorderWidth())
	gc.FillStroke()
	if b.checked {
		if b.typ == Checkbox {
			gc.SetLineWidth(b.LineWidth())
			if b.pushed {
				gc.SetLineColor(b.PushedLineColor())
			} else {
				gc.SetLineColor(b.LineColor())
			}
			gc.MoveTo(b.pos.X+4, b.pos.Y+9)
			gc.LineTo(b.pos.X+8, b.pos.Y+14)
			gc.LineTo(b.pos.X+14, b.pos.Y+5)
			gc.Stroke()
		} else {
			if b.pushed {
				gc.SetFillColor(b.PushedLineColor())
			} else {
				gc.SetFillColor(b.LineColor())
			}
			gc.DrawCircle(mp.X, mp.Y, 0.5*b.LineWidth())
			gc.Fill()
		}
	}
	pt := b.pos.AddXY(b.Width()+b.InnerPadding(), 0.5*b.Height())
	//x := b.pos.X + b.Width() + b.InnerPadding()
	//y := b.pos.Y + 0.5*b.Height()
	gc.SetTextColor(b.TextColor())
	gc.SetFontFace(b.fontFace)
	gc.DrawStringAnchored(b.label, pt.X, pt.Y, 0.0, 0.5)
}

func (b *Radiobox) OnInputEvent(ev InputEvent) {
	b.Button.OnInputEvent(ev)
	if ev.Type == ClickEvent {
		b.checked = !b.checked
	}
}

//----------------------------------------------------------------------------

// Mit Slider kann man einen Schieberegler beliebiger Laenge horizontal oder
// vertikal im GUI positionieren. Als Werte sind aktuell nur Fliesskommazahlen
// vorgesehen.
type Slider struct {
	nodeEmbed
	pushed                                  bool
	orient                                  Orientation
	initValue, minValue, maxValue, stepSize float64
	value                                   binding.Float
	barStart, barEnd                        Point
	barLen                                  float64
	ctrlPos                                 Point
}

func NewSlider(len float64, orient Orientation) *Slider {
	s := &Slider{}
	s.Init(s)
	s.InitProp("Slider")
	s.orient = orient
	d1 := max(0.5*s.BarSize(), 0.5*s.CtrlSize())
	d2 := min(0.5*s.BarSize(), 0.5*s.CtrlSize())
	if s.orient == Horizontal {
		s.SetMinSize(Point{len, s.Height()})
		s.barStart = Point{d2, d1}
	} else {
		s.SetMinSize(Point{s.Width(), len})
		s.barStart = Point{d1, d2}
	}
	s.initValue = 0.0
	s.minValue = 0.0
	s.maxValue = 1.0
	s.stepSize = 0.1
	s.value = binding.NewFloat()
	s.SetValue(s.initValue)
	return s
}

func NewSliderWithData(len float64, orient Orientation, dat binding.Float) *Slider {
	s := NewSlider(len, orient)
	s.value = dat
	return s
}

func NewSliderWithCallback(len float64, orient Orientation,
	callback func(float64)) *Slider {
	s := NewSlider(len, orient)
	s.value.AddCallback(func(data binding.DataItem) {
		callback(data.(binding.Float).Get())
	})
	return s
}

func (s *Slider) SetSize(size Point) {
	s.nodeEmbed.SetSize(size)
	s.updateCtrl()
}

func (s *Slider) updateCtrl() {
	s.barEnd = s.Size().Sub(s.barStart)
	d := s.CtrlSize() - s.BarSize()
	if s.orient == Horizontal {
		s.barLen = s.barEnd.X - s.barStart.X
		p0 := s.barStart.AddXY(0.5*d, 0)
		p1 := s.barEnd.AddXY(-0.5*d, 0)
		s.ctrlPos = p0.Interpolate(p1, s.Factor())
	} else {
		s.barLen = s.barEnd.Y - s.barStart.Y
		p0 := s.barStart.AddXY(0, 0.5*d)
		p1 := s.barEnd.AddXY(0, -0.5*d)
		s.ctrlPos = p0.Interpolate(p1, 1.0-s.Factor())
	}
}

func (s *Slider) Draw(gc *gg.Context) {
	//log.Printf("Slider.Paint()")
	if s.pushed {
		gc.SetLineColor(s.PushedBarColor())
	} else {
		gc.SetLineColor(s.BarColor())
	}
	gc.SetLineWidth(s.BarSize())
	b := s.Bounds()
	if s.orient == Horizontal {
		gc.DrawLine(b.W().X, b.W().Y, b.E().X, b.E().Y)
	} else {
		gc.DrawLine(b.N().X, b.N().Y, b.S().X, b.S().Y)
	}
	gc.Stroke()

	if s.pushed {
		gc.SetLineColor(s.PushedColor())
	} else {
		gc.SetLineColor(s.FillColor())
	}
	gc.SetLineWidth(s.CtrlSize())
	ctrlPos := s.ctrlPos.Add(b.Min)
	if s.orient == Horizontal {
		gc.DrawLine(ctrlPos.X-0.5, ctrlPos.Y, ctrlPos.X+0.5, ctrlPos.Y)
	} else {
		gc.DrawLine(ctrlPos.X, ctrlPos.Y-0.5, ctrlPos.X, ctrlPos.Y+0.5)
	}
	gc.Stroke()
}

func (s *Slider) SetRange(min, max, step float64) {
	s.minValue = min
	s.maxValue = max
	s.stepSize = step
	if s.Value() < s.minValue {
		s.SetValue(min)
	}
	if s.Value() > s.maxValue {
		s.SetValue(max)
	}
}

func (s *Slider) Range() (float64, float64, float64) {
	return s.minValue, s.maxValue, s.stepSize
}

func (s *Slider) SetValue(v float64) {
	v = math.Round(v/s.stepSize) * s.stepSize
	if v > s.maxValue {
		v = s.maxValue
	}
	if v < s.minValue {
		v = s.minValue
	}
	s.value.Set(v)
	s.updateCtrl()
}

func (s *Slider) Value() float64 {
	return s.value.Get()
}

func (s *Slider) SetInitValue(v float64) {
	s.initValue = v
	s.SetValue(v)
}

func (s *Slider) InitValue() float64 {
	return s.initValue
}

func (s *Slider) SetFactor(f float64) {
	if f > 1.0 {
		f = 1.0
	}
	if f < 0.0 {
		f = 0.0
	}
	s.SetValue((1.0-f)*s.minValue + f*s.maxValue)

}

func (s *Slider) Factor() float64 {
	return (s.Value() - s.minValue) / (s.maxValue - s.minValue)
}

func (s *Slider) OnInputEvent(ev InputEvent) {
	switch ev.Type {
	case PressEvent:
		s.pushed = true
	case ReleaseEvent:
		s.pushed = false
	case WheelEvent:
		if ev.WheelRel > 0 {
			s.SetValue(s.Value() + s.stepSize)
		}
		if ev.WheelRel < 0 {
			s.SetValue(s.Value() - s.stepSize)
		}
	case DragEvent, ClickEvent:
		v := 0.0
		r := s.Bounds().Inset(0.5*s.CtrlSize(), 0.5*s.CtrlSize())
		if s.orient == Horizontal {
			v, _ = r.PosRel(ev.Pos)
		} else {
			_, v = r.PosRel(ev.Pos)
		}
		s.SetFactor(v)
		//s.Mark(MarkNeedsPaint)
	case DoubleClickEvent:
		s.SetValue(s.initValue)
		//s.Mark(MarkNeedsPaint)
	}
	s.CallEventHandler(ev)
}
