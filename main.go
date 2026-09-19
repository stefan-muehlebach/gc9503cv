package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/signal"
	"time"

	"periph.io/x/host/v3"

	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gc9503cv/props"
)

//----------------------------------------------------------------------------

const (
	debugFlag = false
)

//----------------------------------------------------------------------------

var (
	windowList = []string{
		"Fading color stripes in RGB",
		"Moving polygons",
		"Moving circles",
		"3D Animation of Polyhedrons",
		"Geometric patterns using GoColors",
		"Showing all available fonts",
		"Shuffle parts of an image randomly",
		"AdaGUI",
	}

	colorList = []color.RGBA{
		color.RGBA{0xff, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0x00, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff},
		color.RGBA{0x00, 0xff, 0xff, 0xff},
		color.RGBA{0x00, 0x00, 0xff, 0xff},
		color.RGBA{0xff, 0x00, 0xff, 0xff},
	}
)

//----------------------------------------------------------------------------

func main() {
	var disp *GC9503CV
	var app *Application
	var dispBounds, drawBounds, drawRect geom.Rectangle[int]
	var progIdx int
	var timeout time.Duration
	var rotate geom.RotationType
	var numObjs int
	var win Window
	var winInfo string
	var sigChan chan os.Signal

	for i, txt := range windowList {
		winInfo += fmt.Sprintf("\n%d - %s", i, txt)
	}

	flag.IntVar(&numObjs, "numObjs", 0, "Number of objects.")
	flag.IntVar(&progIdx, "prog", 0, "Index of program to play."+winInfo)
	flag.Var(&rotate, "rotate", "Rotation of the screen")
	flag.DurationVar(&timeout, "timeout", 10*time.Second,
		"Duration (for animations)")
	flag.Parse()

	log.Printf("PropsMap: %v\n", props.PropsMap)

	if _, err := host.Init(); err != nil {
		log.Fatalf("host.Init(): %v", err)
	}

	sigChan = make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		app.Stop()
	}()

	disp = Open(rotate)
	disp.Init(true)
	dispBounds = disp.DispBounds()
	drawBounds = disp.DrawBounds()
	drawRect = disp.DrawRect()

	if debugFlag {
		log.Printf("display bounds: %v", dispBounds)
		log.Printf("drawing bounds: %v", drawBounds)
		log.Printf("drawing rect  : %v", drawRect)
	}

	app = NewApplication(disp)

	switch progIdx {
	case 0:
		win = NewStripeAnimation(drawBounds, colorList)
	case 1:
		win = NewPolygonAnimation(drawBounds, numObjs)
	case 2:
		win = NewCircleAnimation(drawBounds, 0)
	case 3:
		win = NewPolyhedronAnimation(drawBounds, numObjs)
	case 4:
		win = NewGoColorAnimation(drawBounds)
	case 5:
		win = NewFontsAnimation(drawBounds)
	case 6:
		win = NewShuffleAnimation(drawBounds)
	case 7:
		win = NewGUI(drawBounds)
	default:
		log.Fatalf("No Window with index %d found", progIdx)
	}
	app.SetWindow(win)

	time.AfterFunc(timeout, app.Stop)
	log.Printf("Showing Window Nr. %d", progIdx)
	app.Run()

	log.Printf("----------------------------------------------------")
	log.Printf("Timing statistics:")
	for i := range app.Timer.NumLaps {
		aver := app.Timer.Avg(i + 1)
		mini := app.Timer.Min(i + 1)
		maxi := app.Timer.Max(i + 1)
		log.Printf("  %5d: %v  (%v .. %v)", i+1, aver, mini, maxi)
	}
	aver := app.Timer.Avg(0)
	mini := app.Timer.Min(0)
	maxi := app.Timer.Max(0)
	log.Printf("  total: %v  (%v .. %v)", aver, mini, maxi)
	log.Printf("----------------------------------------------------")
}
