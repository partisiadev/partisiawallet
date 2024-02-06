package view

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/ui/theme"
	"image"
)

type Loader struct {
	AppTheme theme.AppTheme
	loader   material.LoaderStyle
	Size     image.Point
}

func (l *Loader) Layout(gtx layout.Context) layout.Dimensions {
	var th theme.AppTheme
	if l.AppTheme == nil {
		l.AppTheme = theme.GlobalTheme
	}
	th = l.AppTheme
	return layout.Flex{Alignment: layout.Middle,
		Axis:    layout.Vertical,
		Spacing: layout.SpaceSides}.Layout(gtx,
		layout.Flexed(1.0, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					if l.Size == (image.Point{}) {
						l.Size = image.Point{X: gtx.Dp(56), Y: gtx.Dp(56)}
					}
					gtx.Constraints.Min = l.Size
					l.loader.Color = th.Theme().ContrastBg
					return l.loader.Layout(gtx)
				},
			)
		}),
	)
}
