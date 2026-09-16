package main

import (
	"github.com/stefan-muehlebach/gc9503cv/geom"
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/colors"
	"log"
)

type GUI struct {
	windowEmbed
}

func NewGUI(bounds geom.Rectangle[int]) *GUI {
	a := &GUI{}
	a.bounds = bounds
	a.gc = gg.NewContext(a.bounds.Dx(), a.bounds.Dy())

	a.Root = NewPanel(colors.Black)
	a.Root.SetLayoutManager(NewPadLayout())
	a.Root.SetSize(bounds.ToFloat().Size())

	main := NewPanel(colors.Black)
	main.SetLayoutManager(NewVBoxLayout())
	a.Root.Add(main)

	log.Printf("main.Bounds()       : %v", main.Bounds())
	log.Printf("main.Bounds().Size(): %v", main.Bounds().Size())

/*
	btn1 := NewButton(24, 24)
	btn2 := NewButton(32, 32)
	btn3 := NewButton(48, 48)
	btn4 := NewButton(64, 64)
	btn5 := NewButton(96, 96)

	main.Add(btn1, btn2, btn3, btn4, btn5, NewSpacer())
*/

	// Toolbar with IconButtons
	//
	toolBar := NewGroup()
	toolBar.SetLayoutManager(NewHBoxLayout())
	icn1 := NewIconButton("icons/1.png")
	icn2 := NewIconButton("icons/2.png")
	icn3 := NewIconButton("icons/3.png")
	icn4 := NewIconButton("icons/4.png")
	icn5 := NewIconButton("icons/5.png")
	icn6 := NewIconButton("icons/6.png")
	icn7 := NewIconButton("icons/7.png")
	icn8 := NewIconButton("icons/8.png")
	icn9 := NewIconButton("icons/9.png")
	toolBar.Add(icn1, icn2, icn3, icn4, icn5, icn6, icn7, icn8, icn9)
	main.Add(toolBar)

	// Checkboxes and Radiobuttons
	//
	radboxGroup := NewGroup()
	radboxGroup.SetLayoutManager(NewHBoxLayout(30))

	chkGroup := NewGroup()
	chkGroup.SetLayoutManager(NewVBoxLayout())
	chk1 := NewCheckbox("Knoblibrot")
	chk2 := NewCheckbox("Poulet")
	chk3 := NewCheckbox("Fondue")
	chk4 := NewCheckbox("Döner")
	chk5 := NewCheckbox("Racelette")
	chkGroup.Add(chk1, chk2, chk3, chk4, chk5)

	radGroup := NewGroup()
	radGroup.SetLayoutManager(NewVBoxLayout())
	rad1 := NewRadiobutton("Zum selber abholen")
	rad2 := NewRadiobutton("Lieferung")
	rad3 := NewRadiobutton("Express-Post")
	rad4 := NewRadiobutton("Postlagernd")
	radGroup.Add(rad1, rad2, rad3, rad4)

	radboxGroup.Add(chkGroup, radGroup)
	main.Add(radboxGroup)

	// Text-Buttons
	//
	btnGroup := NewGroup()
	btnGroup.SetLayoutManager(NewHBoxLayout())
	btnA := NewTextButton("Hallo")
	btnB := NewTextButton("Benedict")
	btnC := NewTextButton("...hadigärn...")
	btnGroup.Add(btnA, btnB, NewSpacer(), btnC)
	main.Add(btnGroup)


	return a
}
