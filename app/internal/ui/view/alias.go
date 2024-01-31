package view

import (
	"gioui.org/layout"
	"gioui.org/x/component"
	"github.com/partisiadev/partisiawallet/app/internal/ui/shared"
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
	View      = shared.View
	Flex      = layout.Flex
)
