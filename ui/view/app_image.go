package view

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/partisiadev/partisiawallet/assets"
	"github.com/partisiadev/partisiawallet/ui/theme"
)

func DrawAppImageCenter(gtx Gtx, theme theme.AppTheme) Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.20)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.20)
	imgOps := paint.NewImageOp(assets.AppIconImage)
	imgWidget := widget.Image{Src: imgOps, Fit: widget.Contain, Position: layout.Center, Scale: 0}
	return imgWidget.Layout(gtx)
}
