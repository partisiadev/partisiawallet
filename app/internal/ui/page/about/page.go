package about

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/partisiadev/partisiawallet/app/internal/ui/shared"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"
)

type (
	Gtx       = layout.Context
	Dim       = layout.Dimensions
	Animation = component.VisibilityAnimation
	Page      = shared.Page
)

type page struct {
	shared.Manager
	Theme theme.AppTheme
	title string
}

func New(m shared.Manager) shared.View {
	return &page{
		Manager: m,
		Theme:   theme.GlobalTheme,
		title:   "About",
	}
}

func (p *page) Layout(gtx Gtx) Dim {
	if p.Theme == nil {
		p.Theme = theme.GlobalTheme
	}
	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd, Alignment: layout.Start}
	d := flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.H5(p.Theme.Theme(), "About Page").Layout(gtx)
		}),
	)
	return d
}
