package main

import (
	"log"
	"math"
	"image"
	"time"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	//"github.com/stefan-muehlebach/adatft/iliimg"
	"github.com/stefan-muehlebach/framebuffer"
	"gc9503cv/gc9503cv/geom"
	"gc9503cv/gc9503cv/iliimg"
)

const (
	defMOSIPinName = "GPIO23"
	defSCKPinName  = "GPIO24"
	defCSPinName   = "GPIO25"
	defRSTPinName  = "GPIO26"
	defFBDevName   = "/dev/fb1"

	xOffset   = 60
	yOffset   = 0
	longSide  = 960
	shortSide = 480
)

type (
	Point = geom.Point[float64]
	Rectangle = geom.Rectangle[float64]
)

//----------------------------------------------------------------------------

// Dies ist der Datentyp, welche für die Verbindung zum 6.2" TFT-Display
// steht. Voraussetzung ist das Device Tree Overlay 'vc4-kms-dpi-generic'
// und ein Framebuffer Overlay.
type GC9503CV struct {
	mosi, sck, cs, rst     gpio.PinIO
	fb                     *framebuffer.Device
	dispBounds, drawBounds geom.Rectangle[int]
	rot                    RotationType
}

// Geoeffnet wird die Verbindung zum neuen Monitor (zusammen mit weiteren
// Objekten) durch die Funktion Open(). Mit rot wird die gewuenschte
// Rotation des gesamten Monitors angegeben. Es gilt
func Open(rot RotationType) *GC9503CV {
	var offset geom.Point[int]
	var dispSize, drawSize geom.Point[int]
	var err error

	d := &GC9503CV{}
	d.rot = rot

	dispSize = geom.Point[int]{shortSide, longSide}
	switch d.rot {
	case Rot000, Rot180:
		offset = geom.Point[int]{xOffset, yOffset}
		drawSize = geom.Point[int]{shortSide, longSide}.Sub(offset.Mul(2))

	case Rot090, Rot270:
		offset = geom.Point[int]{yOffset, xOffset}
		drawSize = geom.Point[int]{longSide, shortSide}.Sub(offset.Mul(2))
	}
	d.dispBounds = geom.Rectangle[int]{Max: dispSize}
	d.drawBounds = geom.Rectangle[int]{Min: offset, Max: offset.Add(drawSize)}

	if d.mosi = gpioreg.ByName(defMOSIPinName); d.mosi == nil {
		log.Fatal("mosi: gpio pin not found")
	}
	if d.sck = gpioreg.ByName(defSCKPinName); d.mosi == nil {
		log.Fatal("sck: gpio pin not found")
	}
	if d.cs = gpioreg.ByName(defCSPinName); d.mosi == nil {
		log.Fatal("cs: gpio pin not found")
	}
	if d.rst = gpioreg.ByName(defRSTPinName); d.rst == nil {
		log.Fatal("rst: gpio pin not found")
	}
	d.mosi.Out(gpio.Low)
	d.sck.Out(gpio.Low)
	d.cs.Out(gpio.High)
	d.rst.Out(gpio.High)

	d.fb, err = framebuffer.Open(defFBDevName)
	if err != nil {
		log.Fatal("Couldn't open framebuffer")
	}

	return d
}

// Schliesst die Verbindung zum Display und gibt allozierte Ressourcen
// wieder frei.
func (d *GC9503CV) Close() {
	d.fb.Close()
}

// Führt die Initialisierung des Displays durch. Im Wesentlichen ist damit
// die Ausfuehrung einer sog. Initialisierungssequenz verbunden.
func (d *GC9503CV) Init(hwReset bool) (geom.Rectangle[int], geom.Rectangle[int]) {
	d.Reset(hwReset)
	for _, obj := range initCmds {
		d.Cmd(obj.cmd)
		if len(obj.data) > 0 {
			d.DataArray(obj.data)
		}
		if obj.delay > 0 {
			time.Sleep(time.Duration(obj.delay) * time.Millisecond)
		}
	}
	return d.dispBounds, d.drawBounds
}

func (d *GC9503CV) DispBounds() geom.Rectangle[int] {
	return d.dispBounds
}

func (d *GC9503CV) DrawBounds() geom.Rectangle[int] {
	return d.drawBounds
}

func (d *GC9503CV) DrawRect() geom.Rectangle[int] {
	return geom.Rectangle[int]{Max: d.drawBounds.Size()}
}

// Die folgenden Befehle beziehen sich auf den SPI-Bus, ueber welchen der
// Display konfiguriert werden kann. Ueber diesen SPI-Verbindung werden
// keine Bilddaten versendet, dies uebernimmt in der Folge der Framebuffer
// ueber die parallele DPI-Schnittstelle. Das wir durch die parallele
// RGB-Schnittstelle sehr viele Pins des GPIO belegen (insbesondere
// die SPI-Pins) muss SPI "zu Fuss" betrieben werden.

// Sendet 'cmd' als Befehl zum Display.
func (d *GC9503CV) Cmd(cmd uint8) {
	d.cs.Out(gpio.Low)
	d.mosi.Out(gpio.Low)
	d.sck.Out(gpio.Low)
	d.sck.Out(gpio.High)
	d.send(cmd)
	d.cs.Out(gpio.High)
}

// Sendet das Byte 'value' als Datenpaket zum Display
func (d *GC9503CV) Data8(value uint8) {
	d.cs.Out(gpio.Low)
	d.mosi.Out(gpio.High)
	d.sck.Out(gpio.Low)
	d.sck.Out(gpio.High)
	d.send(value)
	d.cs.Out(gpio.High)
}

// Sendet das 32 Bit-Wort 'value' als Datenpaket zum Display.
// ACHTUNG: nicht implementiert! Bitte 'Data8()' oder 'DataArray()' verwenden.
func (d *GC9503CV) Data32(value uint32) {
	log.Fatal("Data32() is not implemented yet!")
}

// Sendet die Daten aus dem Slice 'buf' als Daten zum ILI9341. Dies ist bloss
// eine Hilfsfunktion, damit das Senden von Daten aus einem Slice einfacher
// aufzurufen ist und die ganzen Konvertierungen nicht im Hauptprogramm
// sichtbar sind.
func (d *GC9503CV) DataArray(buf []byte) {
	for _, b := range buf {
		d.Data8(b)
	}
}

// Das eigentliche Versenden der Daten erfolgt mit dieser (nicht exportierten)
// Funktion - wie gesagt: SPI by hand...
func (d *GC9503CV) send(b byte) {
	for range 8 {
		d.sck.Out(gpio.Low)
		if b&0x80 != 0x00 {
			d.mosi.Out(gpio.High)
		} else {
			d.mosi.Out(gpio.Low)
		}
		d.sck.Out(gpio.High)
		b <<= 1
	}
}

// Reset ist eine Hilfsfunktion, mit welcher sich der Display in einen
// definierten Ausgangszustand versetzen laesst. Unser Display kennt
// Software- und auch Hardware-Reset. Mit dem Parameter 'hw' kann ein
// Hardware-Reset ausgefuehrt werden.
func (d *GC9503CV) Reset(hw bool) {
	if hw {
		d.rst.Out(gpio.Low)
		time.Sleep(120 * time.Millisecond)
		d.rst.Out(gpio.High)
		time.Sleep(120 * time.Millisecond)
	} else {
		d.Cmd(0x01)
		time.Sleep(120 * time.Millisecond)
	}
}

func (d *GC9503CV) PartialArea(rect image.Rectangle) {
	buffer   := []byte{0, 0, 0, 0}
	startRow := int16(d.dispBounds.Min.Y - rect.Min.Y)
	endRow   := int16(d.dispBounds.Max.Y - rect.Min.Y)
	buffer[0] = byte((startRow >> 8) & 0xff)
	buffer[1] = byte(startRow & 0x00ff)
	buffer[2] = byte((endRow >> 8) & 0xff)
	buffer[3] = byte(endRow & 0x00ff)
	d.Cmd(PARTON)
	d.Cmd(PARTAREA)
	d.DataArray(buffer)
}

func (d *GC9503CV) Canvas() (canv *Canvas) {
	canv = newCanvas(d.dispBounds.Size())
	switch d.rot {
	case Rot000:
	case Rot090:
		canv.GC.Translate(d.dispBounds.SW().ToFloat().AsCoord())
		canv.GC.Rotate(3.0 * math.Pi / 2.0)
	case Rot180:
		canv.GC.Translate(d.dispBounds.SE().ToFloat().AsCoord())
		canv.GC.Rotate(math.Pi)
	case Rot270:
		canv.GC.Translate(d.dispBounds.NE().ToFloat().AsCoord())
		canv.GC.Rotate(math.Pi / 2.0)
	default:
		log.Fatal("unknown rotation!")
	}
	offset := d.DrawBounds().Min
	canv.GC.Translate(offset.ToFloat().AsCoord())
	return canv
}

func (d *GC9503CV) Send(img *iliimg.ILIImage) {
	//log.Printf("Bounds of the image: %v", img.Bounds())
	if d.dispBounds.Min.Y != img.Rect.Min.Y ||
			d.dispBounds.Max.Y != img.Rect.Max.Y {
		d.PartialArea(img.Rect)
	}
	if d.fb.Bounds() != img.Bounds() {
		log.Fatal("By now, image and screen must have equal size")
	}
	//i := img.PixOffset(img.Rect.Min.X, img.Rect.Min.Y)
	copy(d.fb.Pix, img.Pix)
	d.Cmd(PARTOFF)
}

//----------------------------------------------------------------------------

// Die wenigen Befehle des Displays sind als Kontaten hier definiert koennen
// bspw. eingesetzt werden, um den Code etwas lesbarer zu gestalten.
const (
	SLEEPIN  = 0x10
	SLEEPOUT = 0x11
	PARTON   = 0x12
	PARTOFF  = 0x13
	ALLPOFF  = 0x22
	ALLPON   = 0x23
	DISPOFF  = 0x28
	DISPON   = 0x29
	PARTAREA = 0x30
	IDLEOFF  = 0x38
	IDLEON   = 0x39
	PIXFMT   = 0x3A
	EXTCMDEN = 0xF0
	RGBIFCTL = 0xB0
	DISPCTL  = 0xB1
)

// Ueber diese (nicht exportierten) Objekte wird die Initialisierungs-
// sequenz fuer den Display verwaltet.
type cmdSequence struct {
	cmd   byte
	data  []byte
	delay int
}

var (
	initCmds = []cmdSequence{
		{EXTCMDEN, []byte{0x55, 0xAA, 0x52, 0x08, 0x00}, 0},
		{0xF6, []byte{0x5A, 0x87}, 0},
		{0xC1, []byte{0x3F}, 0},
		{0xCD, []byte{0x25}, 0},
		{0xC9, []byte{0x10}, 0},
		{0xF8, []byte{0x8A}, 0},
		{0xAC, []byte{0x45}, 0},
		{0xA7, []byte{0x47}, 0},
		{0xA0, []byte{0xCC}, 0},
		{0x86, []byte{0x99, 0xA3, 0xA3, 0x31}, 0},
		{0xFA, []byte{0x08, 0x08, 0x00, 0x04}, 0},
		{0xA3, []byte{0x6E}, 0},
		{0xFD, []byte{0x28, 0x3C, 0x00}, 0},
		{0x9A, []byte{0x4a}, 0},
		{0x9B, []byte{0x22}, 0},
		{0x82, []byte{0x00, 0x00}, 0},
		{0x80, []byte{0x54}, 0},
		// DISPLAY_CTL
		{RGBIFCTL, []byte{0x00, 0x0a, 0x0a, 0x0a, 0x0a}, 0},
		// RGB Interface Signals Control
		// 0x33 - Bis auf die Raender: einwandfrei; Farben ok.
		{DISPCTL, []byte{0x33}, 0},
		// Interface Pixel Format
		// 0x60 - Use 18Bit RGB interface; corresponds to DPI[2:0]
		{PIXFMT, []byte{0x60}, 0},

		{0x7A, []byte{0x0F, 0x13}, 0},
		{0x7B, []byte{0x0F, 0x13}, 0},

		{0x6D, []byte{0x0c, 0x03, 0x1e, 0x02, 0x08, 0x1a, 0x19, 0x03, 0x0d,
			0x0e, 0x0f, 0x10, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E,
			0x11, 0x12, 0x13, 0x14, 0x03, 0x19, 0x1a, 0x07, 0x01, 0x1E, 0x03,
			0x0c}, 0},
		{0x64, []byte{0x38, 0x04, 0x03, 0xc4, 0x03, 0x03, 0x38, 0x02, 0x03,
			0xc6, 0x03, 0x03, 0x2C, 0x7A, 0x2C, 0x7A}, 0},
		{0x65, []byte{0x38, 0x08, 0x03, 0xc0, 0x03, 0x03, 0x38, 0x06, 0x03,
			0xc2, 0x03, 0x03, 0x2C, 0x7A, 0x2C, 0x7A}, 0},
		{0x66, []byte{0x83, 0xd0, 0x03, 0xc4, 0x03, 0x03, 0x83, 0xd0, 0x03,
			0xc4, 0x03, 0x03, 0x2C, 0x7A, 0x2C, 0x7A}, 0},
		{0x60, []byte{0x38, 0x0c, 0x3c, 0x3c, 0x38, 0x0b, 0x3c, 0x3c}, 0},
		{0x61, []byte{0xb3, 0xc4, 0x3c, 0x3c, 0xb3, 0xc4, 0x3c, 0x3c}, 0},
		{0x62, []byte{0xb3, 0xc4, 0x3c, 0x3c, 0xb3, 0xc4, 0x3c, 0x3c}, 0},
		{0x63, []byte{0x38, 0x0a, 0x3c, 0x3c, 0x38, 0x09, 0x3c, 0x3c}, 0},
		{0x68, []byte{0x77, 0x08, 0x0a, 0x08, 0x09, 0x00, 0x00, 0x18, 0x0a,
			0x08, 0x09, 0x00, 0x00}, 0},
		{0x69, []byte{0x14, 0x22, 0x14, 0x22, 0x44, 0x22, 0x08}, 0},
		{0x6B, []byte{0x07}, 0},

		{0xD1, []byte{0x00, 0x00, 0x00, 0x10, 0x00, 0x22, 0x00, 0x2c, 0x00,
			0x2e, 0x00, 0x56, 0x00, 0x58, 0x00, 0x7c, 0x00, 0x9a, 0x00, 0xce,
			0x00, 0xfa, 0x01, 0x4c, 0x01, 0x94, 0x01, 0x96, 0x01, 0xda, 0x02,
			0x32, 0x02, 0x76, 0x02, 0xcc, 0x03, 0x18, 0x03, 0x55, 0x03, 0x6b,
			0x03, 0x9b, 0x03, 0xac, 0x03, 0xb8, 0x03, 0xe0, 0x03, 0xFF}, 0},
		{0xD2, []byte{0x00, 0x00, 0x00, 0x10, 0x00, 0x22, 0x00, 0x2c, 0x00,
			0x2e, 0x00, 0x56, 0x00, 0x58, 0x00, 0x7c, 0x00, 0x9a, 0x00, 0xce,
			0x00, 0xfa, 0x01, 0x4c, 0x01, 0x94, 0x01, 0x96, 0x01, 0xda, 0x02,
			0x32, 0x02, 0x76, 0x02, 0xcc, 0x03, 0x18, 0x03, 0x55, 0x03, 0x6b,
			0x03, 0x9b, 0x03, 0xac, 0x03, 0xb8, 0x03, 0xe0, 0x03, 0xFF}, 0},
		{0xD3, []byte{0x00, 0x00, 0x00, 0x10, 0x00, 0x22, 0x00, 0x2c, 0x00,
			0x2e, 0x00, 0x56, 0x00, 0x58, 0x00, 0x7c, 0x00, 0x9a, 0x00, 0xce,
			0x00, 0xfa, 0x01, 0x4c, 0x01, 0x94, 0x01, 0x96, 0x01, 0xda, 0x02,
			0x32, 0x02, 0x76, 0x02, 0xcc, 0x03, 0x18, 0x03, 0x55, 0x03, 0x6b,
			0x03, 0x9b, 0x03, 0xac, 0x03, 0xb8, 0x03, 0xe0, 0x03, 0xFF}, 0},
		{0xD4, []byte{0x00, 0x00, 0x00, 0x10, 0x00, 0x22, 0x00, 0x2c, 0x00,
			0x2e, 0x00, 0x56, 0x00, 0x58, 0x00, 0x7c, 0x00, 0x9a, 0x00, 0xce,
			0x00, 0xfa, 0x01, 0x4c, 0x01, 0x94, 0x01, 0x96, 0x01, 0xda, 0x02,
			0x32, 0x02, 0x76, 0x02, 0xcc, 0x03, 0x18, 0x03, 0x55, 0x03, 0x6b,
			0x03, 0x9b, 0x03, 0xac, 0x03, 0xb8, 0x03, 0xe0, 0x03, 0xFF}, 0},
		{0xD5, []byte{0x00, 0x00, 0x00, 0x10, 0x00, 0x22, 0x00, 0x2c, 0x00,
			0x2e, 0x00, 0x56, 0x00, 0x58, 0x00, 0x7c, 0x00, 0x9a, 0x00, 0xce,
			0x00, 0xfa, 0x01, 0x4c, 0x01, 0x94, 0x01, 0x96, 0x01, 0xda, 0x02,
			0x32, 0x02, 0x76, 0x02, 0xcc, 0x03, 0x18, 0x03, 0x55, 0x03, 0x6b,
			0x03, 0x9b, 0x03, 0xac, 0x03, 0xb8, 0x03, 0xe0, 0x03, 0xFF}, 0},
		{0xD6, []byte{0x00, 0x00, 0x00, 0x10, 0x00, 0x22, 0x00, 0x2c, 0x00,
			0x2e, 0x00, 0x56, 0x00, 0x58, 0x00, 0x7c, 0x00, 0x9a, 0x00, 0xce,
			0x00, 0xfa, 0x01, 0x4c, 0x01, 0x94, 0x01, 0x96, 0x01, 0xda, 0x02,
			0x32, 0x02, 0x76, 0x02, 0xcc, 0x03, 0x18, 0x03, 0x55, 0x03, 0x6b,
			0x03, 0x9b, 0x03, 0xac, 0x03, 0xb8, 0x03, 0xe0, 0x03, 0xFF}, 0},

		{SLEEPOUT, []byte{0x00}, 120},
		{DISPON, []byte{0x00}, 20},
	}
)
