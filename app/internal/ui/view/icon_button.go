package view

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"

	"image"
)

type IconButton struct {
	AppTheme theme.AppTheme
	Button   widget.Clickable
	Icon     *widget.Icon
	Text     string
	layout.Inset
}

func (b *IconButton) Layout(gtx Gtx) Dim {
	btnLayoutStyle := material.ButtonLayout(b.AppTheme.Theme(), &b.Button)
	btnLayoutStyle.CornerRadius = unit.Dp(8)
	return btnLayoutStyle.Layout(gtx, func(gtx Gtx) Dim {
		inset := b.Inset
		if b.Inset == (layout.Inset{}) {
			inset = layout.UniformInset(unit.Dp(12))
		}
		return inset.Layout(gtx, func(gtx Gtx) Dim {
			iconAndLabel := Flex{Alignment: layout.Middle, Spacing: layout.SpaceSides}
			textIconSpacer := unit.Dp(5)

			layIcon := Rigid(func(gtx Gtx) Dim {
				return layout.Inset{Right: textIconSpacer}.Layout(gtx, func(gtx Gtx) Dim {
					var d Dim
					if b.Icon != nil {
						size := gtx.Dp(24)
						d = Dim{Size: image.Pt(size, size)}
						gtx.Constraints = layout.Exact(d.Size)
						d = b.Icon.Layout(gtx, b.AppTheme.Theme().ContrastFg)
					}
					return d
				})
			})

			layLabel := Rigid(func(gtx Gtx) Dim {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx Gtx) Dim {
					l := material.Label(b.AppTheme.Theme(), b.AppTheme.Theme().TextSize, b.Text)
					l.Alignment = text.Middle
					l.Color = b.AppTheme.Theme().Palette.ContrastFg
					return l.Layout(gtx)
				})
			})

			return iconAndLabel.Layout(gtx, layIcon, layLabel)
		})
	})
}
