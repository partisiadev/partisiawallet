package view

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"
)

type ErrorView struct {
	AppTheme theme.AppTheme
	Error    string
}

func (i *ErrorView) Layout(gtx Gtx) (d Dim) {
	if i.AppTheme == nil {
		i.AppTheme = theme.GlobalTheme
	}
	if i.Error != "" {
		return Flex{Axis: layout.Vertical}.Layout(gtx,
			Rigid(func(gtx Gtx) Dim {
				return material.Body1(i.AppTheme.Theme(), i.Error).Layout(gtx)
			}),
		)
	}
	return d
}
