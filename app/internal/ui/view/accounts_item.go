package view

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/partisiadev/partisiawallet/app/assets"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"

	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
)

type accountsItem struct {
	AppTheme theme.AppTheme
	widget.Clickable
	btnSetCurrentIdentity widget.Clickable
	*widget.Enum
}

func (i *accountsItem) Layout(gtx Gtx) Dim {
	if i.AppTheme == nil {
		i.AppTheme = theme.GlobalTheme
	}
	return i.layoutContent(gtx)
}

func (i *accountsItem) IsSelected() bool {
	return false
}

func (i *accountsItem) layoutContent(gtx Gtx) Dim {
	if i.btnSetCurrentIdentity.Clicked(gtx) {
		//_ = chains.GlobalWallet.AddUpdateAccount(&i.Account)
	}

	btnStyle := material.ButtonLayoutStyle{Background: i.AppTheme.Theme().ContrastBg, Button: &i.Clickable}

	if i.IsSelected() || i.Clickable.Hovered() {
		btnStyle.Background.A = 50
	} else {
		btnStyle.Background.A = 10
	}

	d := btnStyle.Layout(gtx, func(gtx Gtx) Dim {
		inset := layout.UniformInset(unit.Dp(16))
		d := inset.Layout(gtx, func(gtx Gtx) Dim {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			flex := Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}
			d := flex.Layout(gtx,
				Rigid(func(gtx Gtx) Dim {
					gtx.Constraints.Max.X = gtx.Constraints.Max.X - gtx.Dp(32)
					flex := Flex{Alignment: layout.Middle}
					d := flex.Layout(gtx,
						Rigid(func(gtx layout.Context) Dim {
							var img image.Image
							//var err error
							//img, _, err = image.Decode(bytes.NewReader(i.Account.PublicImage))
							//if err != nil {
							//	log.Logger().Errorln(err)
							//}
							if img == nil {
								img = assets.AppIconImage
							}
							radii := gtx.Dp(12)
							gtx.Constraints.Max.X, gtx.Constraints.Max.Y = radii*2, radii*2
							bounds := image.Rect(0, 0, radii*2, radii*2)
							clipOp := clip.UniformRRect(bounds, radii).Push(gtx.Ops)
							imgOps := paint.NewImageOp(img)
							imgWidget := widget.Image{Src: imgOps, Fit: widget.Contain, Position: layout.Center, Scale: 0}
							d := imgWidget.Layout(gtx)
							clipOp.Pop()
							return d
						}),
						Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
						Rigid(func(gtx Gtx) Dim {
							label := material.Label(i.AppTheme.Theme(), i.AppTheme.Theme().TextSize, "")
							label.Font.Weight = font.Bold
							return component.TruncatingLabelStyle(label).Layout(gtx)
						}),
					)
					return d
				}),
				Rigid(func(gtx layout.Context) Dim {
					if i.IsSelected() {
						icon, _ := widget.NewIcon(icons.ToggleRadioButtonChecked)
						return icon.Layout(gtx, i.AppTheme.Theme().ContrastBg)
					}
					icon, _ := widget.NewIcon(icons.ToggleRadioButtonUnchecked)
					return icon.Layout(gtx, i.AppTheme.Theme().ContrastBg)
				}),
			)
			return d
		})
		return d
	})

	gtx.Constraints.Max.Y = d.Size.Y
	return d
}
