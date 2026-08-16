package main

import (
	"flag"
	"log"
	"time"
	"github.com/stefan-muehlebach/gg/colors"
	"periph.io/x/host/v3"
	"gc9503cv/gc9503cv/geom"
)

//-----------------------------------------------------------------------------

var (
	disp                   *GC9503CV
	app *Application
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
	var rotate RotationType
	var numObjs int
	var applet Applet

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
	dispBounds, drawBounds = disp.Init(true)
	drawRect = disp.DrawRect()

	app = NewApplication(disp)

	//log.Printf("screen bounds : %v", dispBounds)
	//log.Printf("drawing bounds: %v", drawBounds)
	
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
	}
	app.SetApplet(applet)

	log.Printf("Starting Applet Nr. %d", progIdx)
	app.Run(timeout)

	app.PrintWatchStats()
}

