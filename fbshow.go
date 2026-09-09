//go:build ignore

package main

import (
	"flag"
	"os"
	"fmt"
	"log"
	"image"
	"image/draw"
	_ "image/png"
	"github.com/stefan-muehlebach/gc9503cv/framebuffer"
	"github.com/stefan-muehlebach/gc9503cv/iliimg"
	"github.com/stefan-muehlebach/gc9503cv/geom"
)

func main() {
	var count int
	var method int
	var rotate geom.RotationType

	flag.IntVar(&count, "count", 1, "number of iterations")
	flag.IntVar(&method, "method", 0, "technology for copying")
	flag.Var(&rotate, "rotate", "rotation of the screen")
	flag.Parse()

	fb, err := framebuffer.Open("/dev/fb0")
	if err != nil {
		log.Fatal(err)
	}
	fh, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(fh)
	if err != nil {
		log.Fatal(err)
	}
	fh.Close()

	//img := image.NewRGBA(fb.Bounds())
	ili := iliimg.NewILIImage(fb.Bounds())
	//off := image.Point{60, 0}

	watch := NewStopwatch("Framebuffer")
	
	for range count {
		switch method {
		case 0:
			watch.Start()
			draw.Draw(fb, image.Rectangle{Min: image.Point{60, 0},
				Max: image.Point{420, 960}}, img, image.Point{0, 0},
				draw.Over)
			watch.Stop()
		case 1:
			watch.Start()
			draw.Draw(fb, image.Rectangle{Min: image.Point{60, 0},
				Max: image.Point{420, 960}}, img, image.Point{0, 0},
				draw.Src)
			watch.Stop()
		case 2:
			watch.Start()
			ili.Convert(img.(*image.RGBA), rotate)
			copy(fb.Pix, ili.Pix)
			watch.Stop()
		default:
			log.Fatalf("no method with index %d found", method)
		}
	}

	fmt.Printf("watch: %v  (%v .. %v)\n", watch.Avg(0), watch.Min(0), watch.Max(0))
	fb.Close()
}

