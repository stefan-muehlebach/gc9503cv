package main

import (
	"flag"
	"gc9503cv/gc9503cv/geom"
	"log"
	"time"

	"github.com/stefan-muehlebach/gg/colors"
	"periph.io/x/host/v3"
)

//-----------------------------------------------------------------------------

var (
	disp                             *GC9503CV
	app                              *Application
	dispBounds, drawBounds, drawRect geom.Rectangle[int]

	colorList = []colors.RGBA{
		colors.RGBA{0xff, 0x00, 0x00, 0xff},
		colors.RGBA{0x00, 0xff, 0x00, 0xff},
		colors.RGBA{0x00, 0x00, 0xff, 0xff},
		colors.RGBA{0xff, 0xff, 0x00, 0xff},
		colors.RGBA{0x00, 0xff, 0xff, 0xff},
		colors.RGBA{0xff, 0x00, 0xff, 0xff},
	}
)

func main() {
	var progIdx int
	var timeout time.Duration
	var rotate geom.RotationType
	var numObjs int
	var applet Applet
	//var pointList []Point

	flag.IntVar(&numObjs, "numObjs", 0, "Number of objects.")
	flag.IntVar(&progIdx, "prog", 0, "Index of program to play.")
	flag.Var(&rotate, "rotate", "Rotation of the screen")
	flag.DurationVar(&timeout, "timeout", 10*time.Second,
		"Duration (for animations)")
	flag.Parse()

	if _, err := host.Init(); err != nil {
		log.Fatalf("host.Init(): %v", err)
	}

	disp = Open(rotate)
	disp.Init(true)
	dispBounds = disp.DispBounds()
	drawBounds = disp.DrawBounds()
	drawRect = disp.DrawRect()
	//m := disp.Matrix()

	//log.Printf("display bounds: %v", dispBounds)
	//log.Printf("drawing bounds: %v", drawBounds)
	//log.Printf("drawing rect  : %v", drawRect)

	app = NewApplication(disp)
	//log.Printf("app.mainImg: %v", app.mainImg.Bounds())
	//log.Printf("app.pixBuf : %v", app.pixBuf.Bounds())

/*
	switch rotate {
	case geom.Rot000, geom.Rot180:
		pointList = []Point{
			{0, 0},
			{360, 0},
			{0, 960},
			{360, 960},
		}
	case geom.Rot090, geom.Rot270:
		pointList = []Point{
			{0, 0},
			{960, 0},
			{0, 360},
			{960, 360},
		}
	}
*/
   	//for _, pt := range pointList {
	//	ptNew := m.Transform(pt)
	//	log.Printf("%v -> %v", pt, ptNew)
	//}

	switch progIdx {
	case 0:
		applet = NewStripeAnimation(drawBounds, colorList)
	case 1:
		applet = NewPolygonAnimation(drawBounds, numObjs)
	case 2:
		applet = NewCircleAnimation(drawBounds, numObjs)
	case 3:
		applet = NewPlatonicAnimation(drawBounds, numObjs)
	case 4:
		applet = NewGoColorAnimation(drawBounds)
	case 5:
		applet = NewFontsAnimation(drawBounds)
	default:
		log.Fatalf("No Applet with index %d found", progIdx)
	}
	app.SetApplet(applet)

	log.Printf("Starting Applet Nr. %d", progIdx)
	app.Run(timeout)

	app.PrintWatchStats()
}
