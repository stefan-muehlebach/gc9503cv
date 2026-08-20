//go:build ignore

package main

/*
#include <sys/ioctl.h>
#include <linux/fb.h>

struct fb_fix_screeninfo getFixScreenInfo(int fd) {
	struct fb_fix_screeninfo info;
	ioctl(fd, FBIOGET_FSCREENINFO, &info);
	return info;
}

struct fb_var_screeninfo getVarScreenInfo(int fd) {
	struct fb_var_screeninfo info;
	ioctl(fd, FBIOGET_VSCREENINFO, &info);
	return info;
}
*/
import "C"

import (
	"image"
	"image/color"
	"log"
	"os"
)

const (
	fbFileName = "/dev/fb1"
)

// Device represents the frame buffer. It implements the draw.Image interface.
type Device struct {
	file       *os.File
	Pix     []byte
	Stride      int
	Rect     image.Rectangle
	colorModel color.Model
}

// Open expects a framebuffer device as its argument (such as "/dev/fb0"). The
// device will be memory-mapped to a local buffer. Writing to the device changes
// the screen output.
// The returned Device implements the draw.Image interface. This means that you
// can use it to copy to and from other images.
// The only supported color model for the specified frame buffer is RGB565.
// After you are done using the Device, call Close on it to unmap the memory and
// close the framebuffer file.
func Info(deviceFile string) (error) {
	file, err := os.OpenFile(deviceFile, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return err
	}
	defer file.Close()

	fixInfo := C.getFixScreenInfo(C.int(file.Fd()))
	varInfo := C.getVarScreenInfo(C.int(file.Fd()))

	log.Printf("fixed screen info:")
	log.Printf("  id")
	log.Printf("  smem_start")
	log.Printf("  smem_len")
	log.Printf("  type")
	log.Printf("  type_aux")
	log.Printf("  visual")
	log.Printf("  xpanstep")
	log.Printf("  ypanstep")
	log.Printf("  ywarpstep")
	log.Printf("  line_length")
	log.Printf("  mmio_start")
	log.Printf("  mmio_len")
	log.Printf("  accel")
	log.Printf("  capabilities")
	log.Printf("  reserved")

	log.Printf("variable screen info:")
	log.Printf("%+v", varInfo)

	return nil
}

func main() {
	Info(fbFileName)
}

