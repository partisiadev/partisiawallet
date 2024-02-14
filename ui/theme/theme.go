package theme

import (
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/assets/fonts"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image/color"
)

var GlobalTheme = newTheme()

type appTheme struct {
	th *material.Theme
}
type AppTheme interface {
	Theme() *material.Theme
	Clone() AppTheme
}

func (t *appTheme) Theme() *material.Theme {
	return t.th
}

func (t *appTheme) Clone() AppTheme {
	var appTh appTheme
	if t == nil {
		appTh.th = material.NewTheme()
		appTh.th.Shaper = text.NewShaper(text.WithCollection(fonts.Collection))
		appTh.th.ContrastBg = color.NRGBA{R: 10, B: 40, A: 255}
		appTh.th.ContrastFg = color.NRGBA(colornames.White)
		return &appTh
	}
	appTh = *t
	th := *t.Theme()
	appTh.th = &th
	return &appTh
}

var globalTheme *appTheme

func newTheme() AppTheme {
	if globalTheme == nil {
		globalTheme = &appTheme{}
		globalTheme.th = material.NewTheme()
		globalTheme.th.Shaper = text.NewShaper(text.WithCollection(fonts.Collection))
		globalTheme.th.ContrastBg = color.NRGBA{R: 10, B: 40, A: 255}
		globalTheme.th.ContrastFg = color.NRGBA(colornames.White)
	}
	return globalTheme
}
