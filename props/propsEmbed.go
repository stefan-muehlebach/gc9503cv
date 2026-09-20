//
// THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDTI
//

package props

import (
    "github.com/stefan-muehlebach/gg/colors"
    "github.com/stefan-muehlebach/gg/fonts"
)

type PropertyEmbed struct {
    prop *Properties
}

func (pe *PropertyEmbed) Init(parent *Properties) {
    pe.prop = NewProperties(parent)
}
func (pe *PropertyEmbed) InitProp(name string) {
    pe.prop = PropsMap[name]
}


func (pe *PropertyEmbed) FillColor() (colors.RGBA) {
    return pe.prop.Color(FillColor)
}
func (pe *PropertyEmbed) SetFillColor(c colors.RGBA) {
    pe.prop.SetColor(FillColor, c)
}

func (pe *PropertyEmbed) PushedColor() (colors.RGBA) {
    return pe.prop.Color(PushedColor)
}
func (pe *PropertyEmbed) SetPushedColor(c colors.RGBA) {
    pe.prop.SetColor(PushedColor, c)
}

func (pe *PropertyEmbed) SelectedColor() (colors.RGBA) {
    return pe.prop.Color(SelectedColor)
}
func (pe *PropertyEmbed) SetSelectedColor(c colors.RGBA) {
    pe.prop.SetColor(SelectedColor, c)
}

func (pe *PropertyEmbed) BorderColor() (colors.RGBA) {
    return pe.prop.Color(BorderColor)
}
func (pe *PropertyEmbed) SetBorderColor(c colors.RGBA) {
    pe.prop.SetColor(BorderColor, c)
}

func (pe *PropertyEmbed) PushedBorderColor() (colors.RGBA) {
    return pe.prop.Color(PushedBorderColor)
}
func (pe *PropertyEmbed) SetPushedBorderColor(c colors.RGBA) {
    pe.prop.SetColor(PushedBorderColor, c)
}

func (pe *PropertyEmbed) SelectedBorderColor() (colors.RGBA) {
    return pe.prop.Color(SelectedBorderColor)
}
func (pe *PropertyEmbed) SetSelectedBorderColor(c colors.RGBA) {
    pe.prop.SetColor(SelectedBorderColor, c)
}

func (pe *PropertyEmbed) TextColor() (colors.RGBA) {
    return pe.prop.Color(TextColor)
}
func (pe *PropertyEmbed) SetTextColor(c colors.RGBA) {
    pe.prop.SetColor(TextColor, c)
}

func (pe *PropertyEmbed) PushedTextColor() (colors.RGBA) {
    return pe.prop.Color(PushedTextColor)
}
func (pe *PropertyEmbed) SetPushedTextColor(c colors.RGBA) {
    pe.prop.SetColor(PushedTextColor, c)
}

func (pe *PropertyEmbed) SelectedTextColor() (colors.RGBA) {
    return pe.prop.Color(SelectedTextColor)
}
func (pe *PropertyEmbed) SetSelectedTextColor(c colors.RGBA) {
    pe.prop.SetColor(SelectedTextColor, c)
}

func (pe *PropertyEmbed) LineColor() (colors.RGBA) {
    return pe.prop.Color(LineColor)
}
func (pe *PropertyEmbed) SetLineColor(c colors.RGBA) {
    pe.prop.SetColor(LineColor, c)
}

func (pe *PropertyEmbed) PushedLineColor() (colors.RGBA) {
    return pe.prop.Color(PushedLineColor)
}
func (pe *PropertyEmbed) SetPushedLineColor(c colors.RGBA) {
    pe.prop.SetColor(PushedLineColor, c)
}

func (pe *PropertyEmbed) SelectedLineColor() (colors.RGBA) {
    return pe.prop.Color(SelectedLineColor)
}
func (pe *PropertyEmbed) SetSelectedLineColor(c colors.RGBA) {
    pe.prop.SetColor(SelectedLineColor, c)
}

func (pe *PropertyEmbed) BarColor() (colors.RGBA) {
    return pe.prop.Color(BarColor)
}
func (pe *PropertyEmbed) SetBarColor(c colors.RGBA) {
    pe.prop.SetColor(BarColor, c)
}

func (pe *PropertyEmbed) PushedBarColor() (colors.RGBA) {
    return pe.prop.Color(PushedBarColor)
}
func (pe *PropertyEmbed) SetPushedBarColor(c colors.RGBA) {
    pe.prop.SetColor(PushedBarColor, c)
}

func (pe *PropertyEmbed) BackgroundColor() (colors.RGBA) {
    return pe.prop.Color(BackgroundColor)
}
func (pe *PropertyEmbed) SetBackgroundColor(c colors.RGBA) {
    pe.prop.SetColor(BackgroundColor, c)
}

func (pe *PropertyEmbed) MenuBackgroundColor() (colors.RGBA) {
    return pe.prop.Color(MenuBackgroundColor)
}
func (pe *PropertyEmbed) SetMenuBackgroundColor(c colors.RGBA) {
    pe.prop.SetColor(MenuBackgroundColor, c)
}

func (pe *PropertyEmbed) RegularFont() (*fonts.Font) {
    return pe.prop.Font(RegularFont)
}
func (pe *PropertyEmbed) SetRegularFont(f *fonts.Font) {
    pe.prop.SetFont(RegularFont, f)
}

func (pe *PropertyEmbed) BoldFont() (*fonts.Font) {
    return pe.prop.Font(BoldFont)
}
func (pe *PropertyEmbed) SetBoldFont(f *fonts.Font) {
    pe.prop.SetFont(BoldFont, f)
}

func (pe *PropertyEmbed) Width() (float64) {
    return pe.prop.Size(Width)
}
func (pe *PropertyEmbed) SetWidth(s float64) {
    pe.prop.SetSize(Width, s)
}

func (pe *PropertyEmbed) Height() (float64) {
    return pe.prop.Size(Height)
}
func (pe *PropertyEmbed) SetHeight(s float64) {
    pe.prop.SetSize(Height, s)
}

func (pe *PropertyEmbed) BorderWidth() (float64) {
    return pe.prop.Size(BorderWidth)
}
func (pe *PropertyEmbed) SetBorderWidth(s float64) {
    pe.prop.SetSize(BorderWidth, s)
}

func (pe *PropertyEmbed) PushedBorderWidth() (float64) {
    return pe.prop.Size(PushedBorderWidth)
}
func (pe *PropertyEmbed) SetPushedBorderWidth(s float64) {
    pe.prop.SetSize(PushedBorderWidth, s)
}

func (pe *PropertyEmbed) SelectedBorderWidth() (float64) {
    return pe.prop.Size(SelectedBorderWidth)
}
func (pe *PropertyEmbed) SetSelectedBorderWidth(s float64) {
    pe.prop.SetSize(SelectedBorderWidth, s)
}

func (pe *PropertyEmbed) LineWidth() (float64) {
    return pe.prop.Size(LineWidth)
}
func (pe *PropertyEmbed) SetLineWidth(s float64) {
    pe.prop.SetSize(LineWidth, s)
}

func (pe *PropertyEmbed) InnerPadding() (float64) {
    return pe.prop.Size(InnerPadding)
}
func (pe *PropertyEmbed) SetInnerPadding(s float64) {
    pe.prop.SetSize(InnerPadding, s)
}

func (pe *PropertyEmbed) Padding() (float64) {
    return pe.prop.Size(Padding)
}
func (pe *PropertyEmbed) SetPadding(s float64) {
    pe.prop.SetSize(Padding, s)
}

func (pe *PropertyEmbed) CornerRadius() (float64) {
    return pe.prop.Size(CornerRadius)
}
func (pe *PropertyEmbed) SetCornerRadius(s float64) {
    pe.prop.SetSize(CornerRadius, s)
}

func (pe *PropertyEmbed) FontSize() (float64) {
    return pe.prop.Size(FontSize)
}
func (pe *PropertyEmbed) SetFontSize(s float64) {
    pe.prop.SetSize(FontSize, s)
}

func (pe *PropertyEmbed) BarSize() (float64) {
    return pe.prop.Size(BarSize)
}
func (pe *PropertyEmbed) SetBarSize(s float64) {
    pe.prop.SetSize(BarSize, s)
}

func (pe *PropertyEmbed) CtrlSize() (float64) {
    return pe.prop.Size(CtrlSize)
}
func (pe *PropertyEmbed) SetCtrlSize(s float64) {
    pe.prop.SetSize(CtrlSize, s)
}

func (pe *PropertyEmbed) FieldSize() (float64) {
    return pe.prop.Size(FieldSize)
}
func (pe *PropertyEmbed) SetFieldSize(s float64) {
    pe.prop.SetSize(FieldSize, s)
}
