package main

import (
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
	"log"
	"path"
	"regexp"
	"strconv"

	"github.com/stefan-muehlebach/gc9503cv/geom"
)

const (
	cursorDir  = "cursors"
	cursorSize = 32
)

var (
	CrossCursor     = OpenCursor("cross", cursorSize)
	CrosshairCursor = OpenCursor("crosshair", cursorSize)
	//GrabCursor         = OpenCursor("grab", cursorSize)
	//GrabbingCursor     = OpenCursor("grabbing", cursorSize)
	//IbeamCursor        = OpenCursor("ibeam", cursorSize)
	LeftPtrCursor = OpenCursor("left_ptr", cursorSize)
	PlusCursor    = OpenCursor("plus", cursorSize)
	//PointingHandCursor = OpenCursor("pointing_hand", cursorSize)
	//RightPtrCursor     = OpenCursor("right_ptr", cursorSize)
	TcrossCursor = OpenCursor("tcross", cursorSize)

	CursorList = []*Cursor{
		CrossCursor,
		CrosshairCursor,
		//GrabCursor,
		//GrabbingCursor,
		//IbeamCursor,
		LeftPtrCursor,
		PlusCursor,
		//PointingHandCursor,
		//RightPtrCursor,
		TcrossCursor,
	}
)

//go:embed cursors/*/*.png
var cursorFS embed.FS

type Cursor struct {
	img     image.Image
	hotspot geom.Point[int]
}

func OpenCursor(cursorName string, size int) *Cursor {
	c := &Cursor{}
	sizeDir := fmt.Sprintf("%d", size)
	fileName := fmt.Sprintf("%s@*.png", cursorName)
	fileNamePattern := path.Join(cursorDir, path.Join(sizeDir, fileName))
	hotspotPattern := regexp.MustCompile("@([0-9]+),([0-9]+)\\.png")
	fileList, err := fs.Glob(cursorFS, fileNamePattern)
	if err != nil {
		log.Fatalf("%s: Couldn't glob filesystem: %v", cursorName, err)
	}
	if len(fileList) != 1 {
		log.Fatalf("%s: Expected one matching file; got %d", cursorName,
			len(fileList))
	}
	coordList := hotspotPattern.FindStringSubmatch(fileList[0])
	if coordList == nil {
		log.Fatalf("Filename does not meet format specs: %s", fileList[0])
	}
	x, _ := strconv.Atoi(coordList[1])
	y, _ := strconv.Atoi(coordList[2])
	c.hotspot = geom.Point[int]{x, y}
	fh, err := cursorFS.Open(fileList[0])
	if err != nil {
		log.Fatalf("Couldn't open file: %v", err)
	}
	defer fh.Close()
	c.img, _, err = image.Decode(fh)
	if err != nil {
		log.Fatalf("Couldn't decode image data: %v", err)
	}
	return c
}
