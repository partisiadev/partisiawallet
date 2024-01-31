package view

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"

	"golang.org/x/exp/shiny/materialdesign/icons"
	"time"
)

type Search struct {
	initialized bool
	widget.Clickable
	EditorAnimated
	Icon     *widget.Icon
	AppTheme theme.AppTheme
}

func (s *Search) Layout(gtx layout.Context) layout.Dimensions {
	if !s.initialized {
		if s.Icon == nil {
			icon, _ := widget.NewIcon(icons.ActionSearch)
			s.Icon = icon
		}
		if s.AppTheme == nil {
			s.AppTheme = theme.GlobalTheme
		}
		if s.Animation.Duration == time.Duration(0) {
			s.Animation.Duration = time.Millisecond * 150
			s.Animation.State = component.Invisible
		}
		s.EditorAnimated.SingleLine = true
		s.initialized = true
	}
	if s.Clicked(gtx) {
		if !s.EditorAnimated.Animation.Animating() {
			s.EditorAnimated.Focus()
		}
		s.Animation.ToggleVisibility(gtx.Now)
	}
	flex := Flex{Alignment: layout.Middle}
	return flex.Layout(gtx,
		Flexed(1, func(gtx layout.Context) Dim {
			return s.EditorAnimated.Layout(gtx)
		}),
		Rigid(func(gtx layout.Context) Dim {
			return material.IconButton(s.AppTheme.Theme(), &s.Clickable, s.Icon, "").Layout(gtx)
		}),
	)
}
