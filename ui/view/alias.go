package view

import (
	"gioui.org/layout"
	"gioui.org/x/component"
	"github.com/partisiadev/partisiawallet/ui/fwk"
)

var (
	Rigid   = layout.Rigid
	Flexed  = layout.Flexed
	Stacked = layout.Stacked
)

type (
	Gtx       = layout.Context
	Dim       = layout.Dimensions
	Animation = component.VisibilityAnimation
	View      = fwk.View
	Flex      = layout.Flex
)
