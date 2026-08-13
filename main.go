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
	var prog Prog

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

	log.Printf("screen bounds : %v", dispBounds)
	log.Printf("drawing bounds: %v", drawBounds)
	
	switch progIdx {
	case 0:
		prog = NewStripeAnimation(colorList, drawBounds)
	case 1:
		prog = NewPolygonAnimation(numObjs, drawRect.ToFloat())
	case 2:
		prog = NewCircleAnimation(numObjs, drawRect.ToFloat())
	case 3:
		prog = NewPlatonicAnimation(numObjs, drawRect.ToFloat())
	case 4:
		prog = NewGoColorAnimation(drawRect)
	}
	app.AddProg(prog)

	log.Print("Starting animation")
	app.Run(timeout)

	app.PrintWatchStats()
}

