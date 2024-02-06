package wallet

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/partisiadev/partisiawallet/db"
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

type Wallet struct {
	shared.Manager
	layout.List
}

func New(m shared.Manager) shared.View {
	return &Wallet{Manager: m}
}

func (p *Wallet) Layout(gtx Gtx) Dim {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.Body1(
			theme.GlobalTheme.Theme(),
			fmt.Sprintf("%T", db.Instance().State().DatabaseExists()),
		).Layout(gtx)
	})
}
