package main

import (
    "github.com/stefan-muehlebach/gg"
    "github.com/stefan-muehlebach/gg/colors"
    "github.com/stefan-muehlebach/gg/fonts"
    //"github.com/stefan-muehlebach/gg/geom"
)

//----------------------------------------------------------------------------

var (
	boldFace, _ = fonts.NewFace(fonts.SeafordBold, 18.0)
)

type Button struct {
	nodeEmbed
	BorderColor, FillColor, TextColor colors.RGBA
	BorderWidth, CornerRadius float64
	Text string
}

func NewButton() *Button {
	b := &Button{}
	b.Init(b)
	b.minSize = Point{36, 36}
	b.BorderColor = colors.Teal.Bright(0.3)
	b.FillColor = colors.Teal.Dark(0.3)
	b.TextColor = colors.White
	b.BorderWidth = 0.0
	b.CornerRadius = 15.0
	b.wrapper = b
	return b
}

func (b *Button) Draw(gc *gg.Context) {
	iRect := b.rect.Inset(b.BorderWidth/2.0, b.BorderWidth/2.0)
    gc.SetLineWidth(b.BorderWidth)
    gc.SetLineColor(b.BorderColor)
    gc.SetFillColor(b.FillColor)
    gc.DrawRoundedRectangle(iRect.Min.X, iRect.Min.Y, iRect.Dx(), iRect.Dy(),
		b.CornerRadius)
    gc.FillStroke()

	if b.Text != "" {
		mp := iRect.Center()
		gc.SetFontFace(boldFace)
		gc.SetTextColor(b.TextColor)
		gc.DrawStringAnchored(b.Text, mp.X, mp.Y, 0.5, 0.5)
	}

    if guiDebugging {
        gc.SetLineWidth(2.0)
        gc.SetLineColor(colors.GoFuchsia)
        gc.DrawRectangle(b.pos.X, b.pos.Y, b.size.X, b.size.Y)
        gc.Stroke()
    }
}

