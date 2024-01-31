package view

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"

	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image/color"
)

type PromptContent struct {
	theme.AppTheme
	btnYes      *widget.Clickable
	btnNo       *widget.Clickable
	HeaderTxt   string
	ContentText string
}

func NewPromptContent(theme theme.AppTheme, headerText string, contentText string, btnYes *widget.Clickable, btnNo *widget.Clickable) View {
	return &PromptContent{
		AppTheme:    theme,
		btnYes:      btnYes,
		btnNo:       btnNo,
		HeaderTxt:   headerText,
		ContentText: contentText,
	}
}

func (p *PromptContent) Layout(gtx Gtx) Dim {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	inset := layout.UniformInset(unit.Dp(16))
	d := inset.Layout(gtx, func(gtx Gtx) Dim {
		return Flex{Axis: layout.Vertical}.Layout(gtx,
			Rigid(func(gtx Gtx) Dim {
				if p.HeaderTxt == "" {
					return Dim{}
				}
				bd := material.Body1(p.AppTheme.Theme(), p.HeaderTxt)
				bd.Font.Weight = font.Bold
				bd.Alignment = text.Middle
				return bd.Layout(gtx)
			}),
			Rigid(func(gtx Gtx) Dim {
				return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
			}),
			Rigid(func(gtx Gtx) Dim {
				if p.ContentText == "" {
					return Dim{}
				}
				bd := material.Body1(p.AppTheme.Theme(), p.ContentText)
				bd.Alignment = text.Middle
				return bd.Layout(gtx)
			}),
			Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			Rigid(func(gtx Gtx) Dim {
				return Flex{Spacing: layout.SpaceSides, Alignment: layout.Middle}.Layout(gtx,
					Rigid(func(gtx Gtx) Dim {
						btn := material.Button(p.AppTheme.Theme(), p.btnYes, "Yes")
						btn.Background = color.NRGBA(colornames.Red500)
						return btn.Layout(gtx)
					}),
					Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
					Rigid(func(gtx Gtx) Dim {
						btn := material.Button(p.AppTheme.Theme(), p.btnNo, "No")
						btn.Background = color.NRGBA(colornames.Green500)
						return btn.Layout(gtx)
					}),
				)
			}),
		)
	})
	return d
}
