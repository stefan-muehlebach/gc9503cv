//go:build ignore

package main

import (
	//"flag"
	//"fmt"
	"image"
	"image/draw"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"
	xdraw "golang.org/x/image/draw"
	//"image/draw"
	//"github.com/stefan-muehlebach/adagui"
	//"github.com/stefan-muehlebach/adatft/gc9503cv"
	//"github.com/stefan-muehlebach/adatft/iliimg"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/fonts"
	//"github.com/stefan-muehlebach/gg/geom"
	//"periph.io/x/host/v3"
	"gc9503cv/gc9503cv/iliimg"
)

var (
	AnimWatch = NewStopwatch()
	DrawWatch = NewStopwatch()
	ConvWatch = NewStopwatch()
	SendWatch = NewStopwatch()

	// pointerEventFile = "/dev/input/event0"
	paletteFile      = "paletten.json"

	cursorList = []*Cursor{
		OpenCursor("cross"),
		OpenCursor("crosshair"),
		OpenCursor("grab"),
		OpenCursor("grabbing"),
		OpenCursor("ibeam"),
		OpenCursor("left_ptr"),
		OpenCursor("pointing_hand"),
		OpenCursor("right_ptr"),
	}
	cursorIdx int
)

func PrintWatchStats() {
	animAvg := AnimWatch.Avg()
	drawAvg := DrawWatch.Avg()
	convAvg := ConvWatch.Avg()
	sendAvg := SendWatch.Avg()

	log.Printf("animation : %v", animAvg)
	log.Printf("drawing   : %v", drawAvg)
	log.Printf("converting: %v", convAvg)
	log.Printf("sending   : %v", sendAvg)
	log.Printf("total     : %v", animAvg+drawAvg+convAvg+sendAvg)

	AnimWatch.Reset()
	DrawWatch.Reset()
	ConvWatch.Reset()
	SendWatch.Reset()
}

func Update(disp *GC9503CV, canv *Canvas, pixBuf *iliimg.ILIImage) {
	ConvWatch.Start()
	pixBuf.Convert(canv.GC.Image().(*image.RGBA))
	ConvWatch.Stop()
	SendWatch.Start()
	disp.Send(pixBuf)
	SendWatch.Stop()
}

func Drawing(disp *GC9503CV, canv *Canvas, idxList []int, doCycle bool,
		timeout time.Duration) {

	//draw.Draw(screen, screen.Bounds(), image.NewUniform(colors.Black),
	//	image.Point{}, draw.Src)

	//_, palMap, err := colors.ReadPaletteFile(paletteFile)
	//if err != nil {
	//	log.Fatal(err)
	//}

	gc := canv.GC
	// size := disp.DrawBounds().Size()
	drawBounds := disp.DrawBounds()
	drawRect := disp.DrawRect()
	drawSize := drawRect.Size()
	off := disp.DrawBounds().Min
	log.Printf("drawSize: %v", drawSize)

	mouse := OpenMouse()
	mouse.SetPosRange(drawRect)
	mouse.SetWheelRange(25, 0, 50)
	mouse.SetCursor(cursorList[cursorIdx])

	gc.SetFillColor(colors.Black)
	gc.SetLineColor(colors.WhiteSmoke)
	gc.SetLineWidth(1.0)
	gc.SetTextColor(colors.WhiteSmoke)
	gc.Clear()

	i := 0
	enterEventLoop := false
	for doCycle || i < len(idxList) {
		i = i % len(idxList)
		switch idxList[i] {
		case 0:
			leftMargin := 20
			topMargin := 0
			squareLen := 5
			numColorSteps := 64
			colorStep := 4

			c0 := colors.RGBA{A: 0xFF}
			t0 := time.Now()
			for time.Since(t0) < timeout {
				for field := range numColorSteps {
					idx := field % 3
					if idx == 0 {
						gc.SetFillColor(colors.Black)
						gc.Clear()
					}
					DrawWatch.Start()
					c0.B = uint8(colorStep * field)
					y0 := topMargin + idx*(numColorSteps*squareLen)
					for row := range numColorSteps {
						c0.G = uint8(colorStep * row)
						y1 := y0 + squareLen*row
						for col := range numColorSteps {
							c0.R = uint8(colorStep * col)
							x0 := leftMargin + squareLen*col
							for i := range squareLen {
								for j := range squareLen {
									gc.SetPixel(off.X+x0+i, off.Y+y1+j, c0)
								}
							}
						}
					}
					DrawWatch.Stop()
					if idx == 2 || field == numColorSteps-1 {
						Update(disp, canv, pixBuf)
						time.Sleep(200 * time.Millisecond)
					}
				}
			}

		case 1:
			colorList := []colors.RGBA{
				colors.RGBA{0xff, 0x00, 0x00, 0xff},
				colors.RGBA{0x00, 0xff, 0x00, 0xff},
				colors.RGBA{0x00, 0x00, 0xff, 0xff},
				colors.RGBA{0xff, 0xff, 0x00, 0xff},
				colors.RGBA{0x00, 0xff, 0xff, 0xff},
				colors.RGBA{0xff, 0x00, 0xff, 0xff},
			}
			anim := NewStripeAnimation(colorList, drawBounds)
			RunAnimation(anim, disp, canv, timeout)

		case 2:
			t0 := time.Now()
			for time.Since(t0) < timeout {
				c0 := colors.RandColor()
				c1 := colors.RandColor()
				gc.SetFillColor(colors.Black)
				gc.Clear()
				DrawWatch.Start()
				for row := range drawSize.Y {
					tRow := float64(row) / float64(drawSize.Y-1)
					c2 := colors.Black.Interpolate(c0, tRow)
					for col := range drawSize.X {
						tCol := float64(col) / float64(drawSize.X-1)
						c3 := c2.Interpolate(c1, tCol)
						gc.SetPixel(off.X+col, off.Y+row, c3)
					}
				}
				DrawWatch.Stop()
				Update(disp, canv, pixBuf)
				time.Sleep(1 * time.Second)
			}

		case 3:
			step := 7.0
			cornerRadius := 15.0
			lineWidth := 2.0

			//dX := geom.Point{step, 0}
			dY := Point{0, step}

			gc.SetLineWidth(lineWidth)
			r := drawRect.ToFloat().Inset(step, step)
			t0 := time.Now()
			for time.Since(t0) < timeout {
				// Moiree
				DrawWatch.Start()
				gc.SetFillColor(colors.Black)
				gc.Clear()
				gc.SetLineColor(colors.RandColor())
				p0 := r.Min
				p1 := p0.AddXY(r.Dx(), 0)
				p2 := p0.AddXY(0, r.Dy())
				p3 := p0.AddXY(r.Dx(), r.Dy())
				for p0.Y < r.Max.Y {
					gc.DrawLine(p0.X, p0.Y, p3.X, p3.Y)
					gc.DrawLine(p2.X, p2.Y, p1.X, p1.Y)
					gc.Stroke()
					p0 = p0.Add(dY)
					p1 = p1.Add(dY)
					p2 = p2.Sub(dY)
					p3 = p3.Sub(dY)
				}
				DrawWatch.Stop()
				Update(disp, canv, pixBuf)
				time.Sleep(1 * time.Second)

				// Rectangles
				DrawWatch.Start()
				gc.SetFillColor(colors.Black)
				gc.Clear()
				gc.SetLineColor(colors.RandColor())
				p0 = r.Min
				w, h := r.Dx(), r.Dy()
				for w > step && h > step {
					gc.DrawRectangle(p0.X, p0.Y, w, h)
					gc.Stroke()
					p0 = p0.AddXY(step, step)
					w = w - (2 * step)
					h = h - (2 * step)
				}
				DrawWatch.Stop()
				Update(disp, canv, pixBuf)
				time.Sleep(1 * time.Second)

				// Rounded rectangles
				DrawWatch.Start()
				gc.SetFillColor(colors.Black)
				gc.Clear()
				gc.SetLineColor(colors.RandColor())
				p0 = r.Min
				w, h = r.Dx(), r.Dy()
				for w > 2*cornerRadius && h > 2*cornerRadius {
					gc.DrawRoundedRectangle(p0.X, p0.Y, w, h, cornerRadius)
					gc.Stroke()
					p0 = p0.AddXY(step, step)
					w = w - (2 * step)
					h = h - (2 * step)
				}
				DrawWatch.Stop()
				Update(disp, canv, pixBuf)
				time.Sleep(1 * time.Second)

				// Ellipses
				DrawWatch.Start()
				gc.SetFillColor(colors.Black)
				gc.Clear()
				gc.SetLineColor(colors.RandColor())
				mp := r.Center()
				rx, ry := r.Dx()/2.0, r.Dy()/2.0
				for rx > step && ry > step {
					gc.DrawEllipse(mp.X, mp.Y, rx, ry)
					gc.Stroke()
					rx = rx - step
					ry = ry - step
				}
				DrawWatch.Stop()
				Update(disp, canv, pixBuf)
				time.Sleep(1 * time.Second)

				// Circles
				DrawWatch.Start()
				gc.SetFillColor(colors.Black)
				gc.Clear()
				gc.SetLineColor(colors.RandColor())
				mp = r.Center()
				rad := (r.Dx()/2.0) - 2.0
				for mp.Y < r.Dy()-rad {
					gc.DrawCircle(mp.X, mp.Y, rad)
					gc.Stroke()
					mp = mp.AddXY(0.0, step)
					rad = rad - 1.5
				}
				DrawWatch.Stop()
				Update(disp, canv, pixBuf)
				time.Sleep(1 * time.Second)
			}

		case 4:
			colorList := colors.Groups[colors.GoColors]
			idx := 0
			t0 := time.Now()
			for time.Since(t0) < timeout {
				for i := range 4 {
					gc.SetFillColor(color.Black)
					gc.Clear()
					DrawWatch.Start()
					for row := range drawSize.Y / 60 {
						y := float64(row * 60)
						for col := range drawSize.X / 60 {
							x := float64(col * 60)
							DrawColorSquare(gc, x, y, i, colors.Map[colorList[idx]])
							idx = (idx + 1) % len(colorList)
						}
					}
					DrawWatch.Stop()
					Update(disp, canv, pixBuf)
					time.Sleep(2 * time.Second)
				}
			}

		case 5:
			/*
				fontList := []string{
					"GoRegular",
					"GoMono",
					"LucidaBright",
					"LucidaSans",
					"Seaford",
					"WorkSans",
					"Garamond",
					"Comfortaa",
				}
			*/
			fontList := fonts.Names
			fontSize := 18.0
			fontColor := colors.Black
			captionFont := fonts.SeafordBold
			captionFontSize := 40.0
			captionFontColor := colors.Black.Alpha(0.3)
			margin := 5.0
			text := "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."

			captionFontFace, _ := fonts.NewFace(captionFont, captionFontSize)
			r := drawRect.ToFloat().Inset(margin, margin)
			p0 := r.SE()
			t0 := time.Now()
			for idx := 0; ; idx = (idx + 1) % len(fontList) {
				if time.Since(t0) > timeout {
					break
				}
				fontName := fontList[idx]
				font := fonts.Map[fontName]
				if font.Id >= 900 {
					continue
				}
				gc.SetFillColor(colors.BurlyWood)
				gc.Clear()
				DrawWatch.Start()
				face, _ := fonts.NewFace(font, fontSize)
				gc.SetFontFace(face)
				gc.SetTextColor(fontColor)
				gc.DrawStringWrapped(text, r.Min.X, r.Min.Y, 0.0, 0.0,
					r.Dx(), 1.3, gg.AlignLeft)
				gc.Push()
				gc.RotateAbout(-math.Pi/2.0, p0.X, p0.Y)
				gc.SetFontFace(captionFontFace)
				gc.SetTextColor(captionFontColor)
				gc.DrawStringAnchored(fontName, p0.X, p0.Y, 0.0, 0.0)
				gc.Pop()

				DrawWatch.Stop()
				Update(disp, canv, pixBuf)
				time.Sleep(2 * time.Second)
			}

		case 6:
			margin := 5.0
			font := fonts.GoMono
			fontSize := 25.0
			fontColor := colors.White
			backColor := colors.Black
			dotColor := colors.Crimson
			dotRadius := 5.0

			r := drawRect.ToFloat().Inset(margin, margin)
			DrawWatch.Start()
			gc.SetFillColor(backColor)
			gc.Clear()
			face, _ := fonts.NewFace(font, fontSize)
			gc.SetFontFace(face)
			gc.SetTextColor(fontColor)
			gc.DrawStringAnchored("NW()", r.NW().X, r.NW().Y, 0.0, 1.0)
			gc.DrawStringAnchored("N()", r.N().X, r.N().Y, 0.5, 1.0)
			gc.DrawStringAnchored("NE()", r.NE().X, r.NE().Y, 1.0, 1.0)
			gc.DrawStringAnchored("W()", r.W().X, r.W().Y, 0.0, 0.5)
			gc.DrawStringAnchored("E()", r.E().X, r.E().Y, 1.0, 0.5)
			gc.DrawStringAnchored("SW()", r.SW().X, r.SW().Y, 0.0, 0.0)
			gc.DrawStringAnchored("S()", r.S().X, r.S().Y, 0.5, 0.0)
			gc.DrawStringAnchored("SE()", r.SE().X, r.SE().Y, 1.0, 0.0)
			gc.SetFillColor(dotColor)
			gc.DrawPoint(r.NW().X, r.NW().Y, dotRadius)
			gc.DrawPoint(r.N().X, r.N().Y, dotRadius)
			gc.DrawPoint(r.NE().X, r.NE().Y, dotRadius)
			gc.DrawPoint(r.W().X, r.W().Y, dotRadius)
			gc.DrawPoint(r.E().X, r.E().Y, dotRadius)
			gc.DrawPoint(r.SW().X, r.SW().Y, dotRadius)
			gc.DrawPoint(r.S().X, r.S().Y, dotRadius)
			gc.DrawPoint(r.SE().X, r.SE().Y, dotRadius)
			gc.Fill()
			DrawWatch.Stop()
			Update(disp, canv, pixBuf)

		case 7:
			anim := NewPolygonAnimation(numObjs, drawRect.ToFloat())
			RunAnimation(anim, disp, canv, timeout)

		case 8:
			anim := NewCircleAnimation(numObjs, drawRect.ToFloat())
			RunAnimation(anim, disp, canv, timeout)

		case 9:
			anim := NewPlatonicAnimation(numObjs, drawRect.ToFloat())
			RunAnimation(anim, disp, canv, timeout)

		case 10:
			DABGui(disp, canv, drawRect.ToFloat())

		case 11:
			margin      := 10.0
			padding     := 10.0
			panelBounds := drawRect.ToFloat().Inset(margin, margin)
			buttonSize  := Point{
				(panelBounds.Dx()-padding)/2.0,
				(panelBounds.Dy()-4*padding)/5.0,
			}

			p0 := panelBounds.Min
			for range 5 {
				b := NewButton()
				b.SetPos(p0)
				b.SetSize(buttonSize)
				b.Text = "Stefan"
				canv.Add(b)
				p1 := p0.AddXY(buttonSize.X+padding, 0.0)
				b = NewButton()
				b.SetPos(p1)
				b.SetSize(buttonSize)
				b.Text = "Benedict"
				canv.Add(b)
				p0 = p0.AddXY(0.0, buttonSize.Y+padding)
			}
			enterEventLoop = true

		case 12:
			imgTwoWord, _ := gg.LoadPNG("images/14_TwoWord.png")
			imgReef, _    := gg.LoadPNG("images/02_Reef.png")
			imgA, imgB := imgTwoWord, imgReef

			log.Printf("gg.DrawImage()")
			DrawWatch.Start()
		    canv.Clear(canv.BackColor)
			gc.DrawImage(imgA, 0, 0)
			Update(disp, canv, pixBuf)
			DrawWatch.Stop()
			PrintWatchStats()
			imgA, imgB = imgB, imgA
			
			log.Printf("draw.Draw()")
			DrawWatch.Start()
		    canv.Clear(canv.BackColor)
			draw.Draw(canv.Img, drawBounds.ToInt(), imgA,
				image.Point{}, draw.Src)
			Update(disp, canv, pixBuf)
			DrawWatch.Stop()
			PrintWatchStats()
			imgA, imgB = imgB, imgA
			
			log.Printf("xdraw.Copy()")
			DrawWatch.Start()
		    canv.Clear(canv.BackColor)
			xdraw.Copy(canv.Img, drawBounds.Min.ToInt(), imgA,
				imgTwoWord.Bounds(), draw.Src, nil)
			Update(disp, canv, pixBuf)
			DrawWatch.Stop()
			PrintWatchStats()

		case 13:
			imgTwoWord, _    := gg.LoadPNG("images/14_TwoWord.png")
			imgBackground, _ := gg.LoadPNG("images/10_Background.png")
			fg, bg := imgTwoWord, imgBackground

        	p0 := image.Point{8, 0}
        	sz := image.Rect(0, 0, 80, 80)
        	dx := image.Point{sz.Dx()+8, 0}
        	dy := image.Point{0, sz.Dy()+8}
        	orig := drawBounds.Min.ToInt()
        
        	originList := make([]image.Point, 44)
        	for i := range 44 {
        		col := i % 4
        		row := i / 4
        		originList[i] = p0.Add(dx.Mul(col)).Add(dy.Mul(row))
        	}
        
        	//dev.SetDoubleBuffer(true)
        	draw.Draw(canv.Img, drawBounds.ToInt(), bg, image.Point{}, draw.Src)
        	for _, pt := range originList {
        		draw.Draw(canv.Img, sz.Add(pt).Add(orig), fg, pt, draw.Src)
        	}
        	//dev.SwitchBuffer()
			Update(disp, canv, pixBuf)
        	time.Sleep(1 * time.Second)
        
        	draw.Draw(canv.Img, drawBounds.ToInt(), bg, image.Point{}, draw.Src)
        	idxList := make([]int, len(originList))
        	for i := range idxList {
        		idxList[i] = i
        	}
        	maxMs := 1000
        	nSteps := 40
        	for i := range nSteps {
        		t := float64(i) / float64(nSteps-1)
        		t = (t * t * t) * float64(maxMs)
        		d := time.Duration(t) * time.Millisecond
        
				DrawWatch.Start()
        	    draw.Draw(canv.Img, drawBounds.ToInt(), bg, image.Point{},
					draw.Src)
        		rand.Shuffle(len(idxList), func(i, j int) {
        			idxList[i], idxList[j] = idxList[j], idxList[i]
        		})
        		for i := range idxList {
        			srcPt := originList[i]
        			dstPt := originList[idxList[i]]
        
        			draw.Draw(canv.Img, sz.Add(dstPt).Add(orig), fg,
						srcPt, draw.Src)
        		}
				DrawWatch.Stop()
        		//dev.SwitchBuffer()
				Update(disp, canv, pixBuf)
        		time.Sleep(d)
        		d += 20 * time.Millisecond
        	}
        	time.Sleep(3 * time.Second)
        	//dev.SetDoubleBuffer(false)
		default:
			log.Printf("No app found!")
		}

		if enterEventLoop {
			go func() {
				var activeNode Node

				for mev := range mouse.EventQ {
					mouse.Pos = mev.Pos
					target := canv.FindTarget(mev.Pos.ToFloat())
					if target == nil {
						continue
					}
					activeNode = target
					switch obj := activeNode.(type) {
					case *Button:
						obj.FillColor = colors.Gold
					}
					
					switch mev.Type {
					case TypeRelease:
						cursorIdx = (cursorIdx + 1) % len(cursorList)
						mouse.SetCursor(cursorList[cursorIdx])
					}
				}
			}()
			mouse.StartEvents()

			ticker := time.NewTicker(30 * time.Millisecond)
			defer ticker.Stop()
			t0 := time.Now()
			for range ticker.C {
				if time.Since(t0) > timeout {
					break
				}
				DrawWatch.Start()
				canv.Clear(canv.BackColor)
				canv.Refresh()
				mouse.Draw(gc)
				DrawWatch.Stop()
				Update(disp, canv, pixBuf)
			}
			mouse.Close()
			enterEventLoop = false
		}
		i += 1
		log.Printf("Time statistics:")
		PrintWatchStats()
	}
}

//-----------------------------------------------------------------------

func DrawColorSquare(gc *gg.Context, x, y float64, typ int, col colors.RGBA) {
	switch typ {
	case 0:
		gc.SetFillColor(col)
		gc.DrawRectangle(x, y, 60.0, 60.0)
		gc.Fill()
	case 1:
		gc.SetFillColor(col)
		gc.DrawCircle(x+30.0, y+30.0, 30.0-1.0)
		gc.Fill()
	case 2:
		gc.SetFillColor(col)
		gc.DrawRectangle(x, y, 60.0, 60.0)
		gc.Fill()
		gc.SetFillColor(colors.Black)
		gc.DrawCircle(x+30.0, y+30.0, 30.0-1.0)
		gc.Fill()
	case 3:
		gc.SetFillColor(col)
		gc.DrawCircle(x+30.0, y+30.0, 30.0)
		gc.Fill()
		gc.SetFillColor(colors.Black)
		gc.DrawCircle(x+30.0, y+30.0, 30.0-15.0)
		gc.Fill()
	}
}
