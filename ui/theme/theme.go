package theme

import (
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/assets/fonts"
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
	return newTheme()
}

var globalTheme *appTheme

func newTheme() AppTheme {
	if globalTheme == nil {
		globalTheme = &appTheme{}
		globalTheme.th = material.NewTheme()
		globalTheme.th.Shaper = text.NewShaper(text.WithCollection(fonts.Collection))
		globalTheme.th.ContrastBg = color.NRGBA{R: 10, B: 40, A: 255}
	}
	return globalTheme
}
