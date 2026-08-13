//go:build ignore

package main

import (
	"log"
	"math"
	"strings"
	//"image"
	//"image/draw"
	//"github.com/stefan-muehlebach/adatft/gc9503cv"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
)

var (
	scrBarWidth     = 24.0
	scrBarBackColor = colors.GoIndigo
	scrBarColor     = colors.GoTeal.Alpha(0.25)
	scrSliderColor  = colors.GoTeal

	statusWidth     = 80.0
	statusBackColor = colors.GoTeal

	panelWidth     = (960.0 - 2*scrBarWidth - statusWidth) / 2.0
	panelBackColor = colors.GoIndigo
	panelPadding   = 10.0

	listFontSize   = 36.0
	textFontSize   = 20.0
	statusFontSize = 18.0

	listLineSpace   = 1.6 * listFontSize
	textLineSpace   = 1.4
	statusLineSpace = 1.2 * statusFontSize

	listFont, _   = fonts.NewFace(fonts.ComfortaaBold, listFontSize)
	textFont, _   = fonts.NewFace(fonts.Comfortaa, textFontSize)
	statusFont, _ = fonts.NewFace(fonts.ComfortaaBold, statusFontSize)

	listTextColor   = colors.GoLightGray
	listSelectColor = colors.GoYellow
	textColor       = colors.GoYellow
	statusTextColor = colors.GoLightBlue

	blindText = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet."

	stationList = []string{
		"SRF 4 News",
		"SRF 2 Kultur",
		"SRF Virus",
		"SRF 1 SG",
		"SRF 1 LU",
		"SRF 3",
		"Swiss Classic",
		"Swiss Jazz",
		"Swiss Pop",
	}
)

func DABGui(scr *GC9503CV, canv *Canvas, rect Rectangle) {

	gc := canv.GC

	p0 := geom.Point{}
	p1 := p0.AddXY(scrBarWidth, rect.Dy())
	listScrBarBounds := geom.Rectangle{
		Min: p0,
		Max: p1}
	p0.X = p1.X
	p1 = p0.AddXY(panelWidth, rect.Dy())
	listPanelBounds := geom.Rectangle{
		Min: p0,
		Max: p1}
	p0.X = p1.X
	p1 = p0.AddXY(scrBarWidth, rect.Dy())
	textScrBarBounds := geom.Rectangle{
		Min: p0,
		Max: p1}
	p0.X = p1.X
	p1 = p0.AddXY(panelWidth, rect.Dy())
	textPanelBounds := geom.Rectangle{
		Min: p0,
		Max: p1}
	p0.X = p1.X
	p1 = p0.AddXY(statusWidth, rect.Dy())
	statusPanelBounds := geom.Rectangle{
		Min: p0,
		Max: p1}

	DrawWatch.Start()
	gc.SetLineWidth(2.0)
	gc.SetLineCapRound()
	gc.SetLineJoinRound()

	// Draw Background
	gc.SetFillColor(panelBackColor)
	//backBounds := listScrBarBounds.Union(textPanelBounds)
	gc.DrawRectangle(listScrBarBounds.AsCoord())
	gc.DrawRectangle(listPanelBounds.AsCoord())
	gc.DrawRectangle(textScrBarBounds.AsCoord())
	gc.DrawRectangle(textPanelBounds.AsCoord())
	gc.Fill()
	gc.SetFillColor(statusBackColor)
	gc.DrawRectangle(statusPanelBounds.AsCoord())
	gc.Fill()
	//gc.SetFillColor(colors.GoAqua)
	//gc.DrawRectangle(drawBounds.Min.X+drawBounds.Dx()-64.0,
	//	drawBounds.Min.Y+drawBounds.Dy()/2, 64.0, drawBounds.Dy()/2)
	//gc.Fill()

	// Draw GUI elements (like Sliders)
	gc.SetLineWidth(scrBarWidth)

	gc.SetLineColor(scrBarColor)
	p0 = listScrBarBounds.N()
	p1 = listScrBarBounds.S()
	gc.DrawLine(p0.X, p0.Y, p1.X, p1.Y)
	gc.Stroke()
	gc.SetLineColor(scrSliderColor)
	p0 = p0.AddXY(0.0, 150.0)
	p1 = p1.SubXY(0.0, 100.0)
	gc.DrawLine(p0.X, p0.Y, p1.X, p1.Y)
	gc.Stroke()

	// Draw the name of some radio stations
	gc.SetFontFace(listFont)
	gc.SetTextColor(listSelectColor)
	panelBounds := listPanelBounds.Inset(panelPadding, panelPadding)
	p0 = panelBounds.W()
	gc.DrawStringAnchored("SRF 1 SG", p0.X, p0.Y, 0.0, 0.5)

	gc.SetTextColor(listTextColor)
	p1 = p0.SubXY(0.0, listLineSpace)
	gc.DrawStringAnchored("SRF Virus", p1.X, p1.Y, 0.0, 0.5)
	p1 = p1.SubXY(0.0, listLineSpace)
	gc.DrawStringAnchored("SRF 2 Kultur", p1.X, p1.Y, 0.0, 0.5)
	p1 = p1.SubXY(0.0, listLineSpace)
	gc.DrawStringAnchored("SRF 4 News", p1.X, p1.Y, 0.0, 0.5)

	p1 = p0.AddXY(0.0, listLineSpace)
	gc.DrawStringAnchored("SRF 1 LU", p1.X, p1.Y, 0.0, 0.5)
	p1 = p1.AddXY(0.0, listLineSpace)
	gc.DrawStringAnchored("SRF 3", p1.X, p1.Y, 0.0, 0.5)
	p1 = p1.AddXY(0.0, listLineSpace)
	gc.DrawStringAnchored("Swiss Classic", p1.X, p1.Y, 0.0, 0.5)

	// Draw some info text
	gc.SetFontFace(textFont)
	gc.SetTextColor(textColor)
	panelBounds = textPanelBounds.Inset(panelPadding, panelPadding)
	textSlc := gc.WordWrap(blindText, panelBounds.Dx())
	w, h := gc.MeasureMultilineString(strings.Join(textSlc, "\n"),
		textLineSpace)
	scrRatio := panelBounds.Dy() / h
	log.Printf("w, h, ratio: %f, %f, %f", w, h, scrRatio)
	p0 = panelBounds.NW()
	gc.DrawStringWrapped(blindText, p0.X, p0.Y, 0.0, 0.0, panelBounds.Dx(),
		textLineSpace, gg.AlignLeft)

	// And its scrollbar
	gc.SetLineWidth(scrBarWidth)
	scrSliderLen := scrRatio * (textScrBarBounds.Dy() - scrBarWidth)

	gc.SetLineColor(scrBarColor)
	p0 = textScrBarBounds.N()
	p1 = textScrBarBounds.S()
	gc.DrawLine(p0.X, p0.Y, p1.X, p1.Y)
	gc.Stroke()
	gc.SetLineColor(scrSliderColor)
	p0 = p0.AddXY(0.0, scrBarWidth/2.0)
	p1 = p1.SubXY(0.0, scrSliderLen)
	gc.DrawLine(p0.X, p0.Y, p1.X, p1.Y)
	gc.Stroke()

	// Draw the stuff in the status bar
	gc.SetFontFace(statusFont)
	gc.SetTextColor(statusTextColor)
	gc.SetLineWidth(7.0)
	gc.SetLineColor(colors.GoLightBlue)
	mp := rect.NE().AddXY(-32.0, 32.0)
	r := 27.0
	gc.DrawArc(mp.X, mp.Y, r, 0.0, -math.Pi)
	gc.Stroke()
	r -= 9.0
	gc.DrawArc(mp.X, mp.Y, r, 0.0, -math.Pi)
	gc.Stroke()
	r -= 9.0
	gc.SetLineColor(colors.GoYellow)
	gc.DrawArc(mp.X, mp.Y, r, 0.0, -math.Pi)
	gc.Stroke()
	gc.SetFillColor(colors.GoYellow)
	gc.DrawCircle(mp.X, mp.Y, 3.5)
	gc.Fill()

	mp = mp.AddXY(0.0, 20.0)
	gc.DrawStringAnchored("23:45", mp.X, mp.Y, 0.5, 0.5)
	DrawWatch.Stop()

	Update(scr, canv, pixBuf)

}
