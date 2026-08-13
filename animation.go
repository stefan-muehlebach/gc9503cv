//go:build ignore

package main

import (
	//"github.com/stefan-muehlebach/adatft/gc9503cv"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"time"
)

//-----------------------------------------------------------------------------

type Animation interface {
	Init(gc *gg.Context)
	Update(dt time.Duration)
	Draw(gc *gg.Context)
}

func RunAnimation(anim Animation, disp *GC9503CV, canv *Canvas,
	timeout time.Duration) {
	dt := 25 * time.Millisecond

	anim.Init(canv.GC)
	ticker := time.NewTicker(dt)
	defer ticker.Stop()
	t0 := time.Now()
	for range ticker.C {
		if time.Since(t0) > timeout {
			break
		}
		AnimWatch.Start()
		anim.Update(dt)
		AnimWatch.Stop()

		DrawWatch.Start()
		canv.Clear(colors.Black)
		anim.Draw(canv.GC)
		DrawWatch.Stop()

		Update(disp, canv, pixBuf)
	}
}
