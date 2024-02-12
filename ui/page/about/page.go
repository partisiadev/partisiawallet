package about

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/ui/fwk"
	"github.com/partisiadev/partisiawallet/ui/theme"
)

type About struct {
	fwk.Manager
	layout.List
}

func New(m fwk.Manager) fwk.View {
	return &About{Manager: m}
}

func (p *About) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.Body1(theme.GlobalTheme.Theme(), "About Page").Layout(gtx)
	})
}
