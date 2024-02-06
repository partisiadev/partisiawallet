package chains

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/partisiadev/partisiawallet/ui/shared"
	"github.com/partisiadev/partisiawallet/ui/theme"
)

type (
	Gtx         = layout.Context
	Dim         = layout.Dimensions
	Animation   = component.VisibilityAnimation
	View        = shared.View
	ModalOption = shared.ModalOption
)

type Chains struct {
	shared.Manager
	layout.List
}

func New(m shared.Manager) shared.View {
	return &Chains{Manager: m}
}

func (p *Chains) Layout(gtx Gtx) Dim {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.Body1(theme.GlobalTheme.Theme(), "Chains Page").Layout(gtx)
	})
}
