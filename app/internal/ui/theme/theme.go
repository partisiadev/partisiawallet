package theme

import (
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/app/assets/fonts"
	"image/color"
	"sync"
)

var GlobalTheme = newTheme()

type appTheme struct {
	sync.Once
	th *material.Theme
}
type AppTheme interface {
	Theme() *material.Theme
	Clone() AppTheme
}

func (t *appTheme) initialize() {
	if t.th == nil {
		t.th = newTheme().Theme()
	}
}

func (t *appTheme) Theme() *material.Theme {
	t.Do(t.initialize)
	return t.th
}

func (t *appTheme) Clone() AppTheme {
	t.Do(t.initialize)
	return newTheme()
}

func newTheme() AppTheme {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(fonts.Collection))
	th.ContrastBg = color.NRGBA{R: 10, B: 40, A: 255}
	return &appTheme{th: th}
}
