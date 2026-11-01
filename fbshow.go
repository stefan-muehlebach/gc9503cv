//go:build ignore

package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"

	"github.com/stefan-muehlebach/gc9503cv/framebuffer"
	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gc9503cv/iliimg"
)

func main() {
	var count int
	var method int
	var rotate geom.RotationType

	flag.IntVar(&count, "count", 1, "number of iterations")
	flag.IntVar(&method, "method", 1, "technology for copying, (1: Draw-Over, 2: Draw-Src, 3: Convert-Copy")
	flag.Var(&rotate, "rotate", "rotation of the screen")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatalf("usage: %s <png-file>", os.Args[0])
	}

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
	offset := image.Point{60, 0}
	size := image.Point{360, 960}
	imgBounds := image.Rectangle{Min: offset, Max: size.Add(offset)}

	watch := NewStopwatch("Framebuffer")

	for range count {
		switch method {
		case 1:
			watch.Start()
			draw.Draw(fb, imgBounds, img, image.Point{0, 0}, draw.Over)
			watch.Stop()
		case 2:
			watch.Start()
			draw.Draw(fb, imgBounds, img, image.Point{0, 0}, draw.Src)
			watch.Stop()
		case 3:
			watch.Start()
			ili.Convert(img.(*image.NRGBA), rotate)
			copy(fb.Pix, ili.Pix)
			watch.Stop()
		default:
			log.Fatalf("no method with index %d found", method)
		}
	}

	fmt.Printf("watch: %v  (%v .. %v)\n", watch.Avg(0), watch.Min(0), watch.Max(0))
	fb.Close()
}
