package main

import (
	"fmt"
	"math"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
)

const (
	Rot000 = iota
	Rot090
	Rot180
	Rot270
)

func main() {
	var t, r *geom.Matrix

	rot := Rot180

	gc := gg.NewContext(480, 960)
	drawRect := gc.Bounds()

	switch rot {
	case Rot000:
		t = geom.Identity()
		r = geom.Identity()
	case Rot090:
		t = geom.Translate(gc.Bounds().NE())
		r = geom.Rotate(math.Pi / 2.0)
	case Rot180:
		t = geom.Translate(gc.Bounds().SE())
		r = geom.Rotate(math.Pi)
	case Rot270:
		t = geom.Translate(gc.Bounds().SW())
		r = geom.Rotate(3.0 * math.Pi / 2.0)
	}
	gc.Multiply(t)
	gc.Multiply(r)
	size := gc.Matrix().TransformRect(drawRect).Canon().Size()
	drawRect = geom.Rectangle{Max: size}

	gc.SetLineColor(colors.Red)
	gc.SetLineWidth(4.0)
	gc.DrawLine(0.0, 0.0, 50.0, 0.0)
	gc.Stroke()
	gc.SetLineColor(colors.Green)
	gc.DrawLine(0.0, 0.0, 0.0, 50.0)
	gc.Stroke()

	face, _ := fonts.NewFace(fonts.LucidaBright, 18.0)
	gc.SetFontFace(face)
	gc.DrawStringAnchored("Oben", drawRect.N().X, drawRect.N().Y, 0.5, 1.0)
	gc.DrawStringAnchored("Unten", drawRect.S().X, drawRect.S().Y, 0.5, 0.0)
	gc.DrawStringAnchored("Links", drawRect.W().X, drawRect.W().Y, 0.0, 0.5)
	gc.DrawStringAnchored("Rechts", drawRect.E().X, drawRect.E().Y, 1.0, 0.5)

	fmt.Printf("DrawRect: %v\n", drawRect)
	gc.SavePNG("foo01.png")
}
