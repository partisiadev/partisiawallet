package view

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"

	"image"
	"time"
)

type EditorAnimated struct {
	initialized bool
	Animation   component.VisibilityAnimation
	widget.Editor
	AppTheme theme.AppTheme
	layout.Inset
}

func (e *EditorAnimated) Layout(gtx Gtx) Dim {
	if !e.initialized {
		if e.AppTheme == nil {
			e.AppTheme = theme.GlobalTheme
		}
		if e.Animation.Duration == time.Duration(0) {
			e.Animation.Duration = time.Millisecond * 150
			e.Animation.State = component.Invisible
		}
		if e.Inset == (layout.Inset{}) {
			e.Inset = layout.Inset{Left: 16, Right: 16}
		}
		e.Editor.SingleLine = true
		e.initialized = true
	}
	return e.Inset.Layout(gtx, func(gtx layout.Context) Dim {
		progress := e.Animation.Revealed(gtx)
		if progress == 0 {
			return Dim{}
		}
		inset := layout.Inset{Top: 8, Bottom: 8, Left: 12, Right: 12}
		rec := op.Record(gtx.Ops)
		dims := inset.Layout(gtx, func(gtx layout.Context) Dim {
			return material.Editor(e.AppTheme.Theme(), &e.Editor, "").Layout(gtx)
		})
		call := rec.Stop()
		radii := 8
		dims.Size.X = int(float32(dims.Size.X) * progress)
		//dims.Size.Y = int(float32(dims.Size.Y) * progress)
		radii = int(float32(radii) * progress)
		rect := component.Rect{Color: e.AppTheme.Theme().Bg, Size: dims.Size, Radii: radii}
		rect.Layout(gtx)
		rRect := clip.RRect{Rect: image.Rect(0, 0,
			dims.Size.X,
			dims.Size.Y),
			SE: radii, SW: radii, NW: radii, NE: radii}.Push(gtx.Ops)
		call.Add(gtx.Ops)
		rRect.Pop()
		return dims
	})
}
