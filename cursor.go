package main

import (
	"embed"
	"image"
	"io/fs"
	"log"
	"path"
	"regexp"
	"strconv"
	"gc9503cv/gc9503cv/geom"
)

const (
	cursorDir = "cursors"
)

var (
	CrossCursor = OpenCursor("cross")
	CrosshairCursor = OpenCursor("crosshair")
	GrabCursor = OpenCursor("grab")
	GrabbingCursor = OpenCursor("grabbing")
	IbeamCursor = OpenCursor("ibeam")
	LeftPtrCursor = OpenCursor("left_ptr")
	PointingHandCursor = OpenCursor("pointing_hand")
	RightPtrCursor = OpenCursor("right_ptr")

	CursorList = []*Cursor{
    	CrossCursor,
    	CrosshairCursor,
    	GrabCursor,
    	GrabbingCursor,
    	IbeamCursor,
    	LeftPtrCursor,
    	PointingHandCursor,
    	RightPtrCursor,
	}
)

//go:embed cursors/*.png
var cursorFS embed.FS

type Cursor struct {
	img     image.Image
	hotspot geom.Point[int]
}

func OpenCursor(cursorName string) *Cursor {
	c := &Cursor{}
	fileNamePattern := path.Join(cursorDir, cursorName+"@*.png")
	hotspotPattern := regexp.MustCompile("@([0-9]+),([0-9]+)\\.png")
	fileList, err := fs.Glob(cursorFS, fileNamePattern)
	if err != nil {
		log.Fatalf("Couldn't glob filesystem: %v", err)
	}
	if len(fileList) != 1 {
		log.Fatalf("Expected only one matching file: %v", err)
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

