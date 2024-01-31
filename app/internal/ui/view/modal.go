package view

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/x/component"
	"github.com/partisiadev/partisiawallet/app/internal/ui/shared"
	"image/color"
	"time"
)

type Modal struct {
	btnBackdrop        widget.Clickable
	btnContent         widget.Clickable
	option             shared.ModalOption
	dismissWithoutAnim bool
}

func (m *Modal) Show(option shared.ModalOption) {
	if option.Widget == nil {
		return
	}
	m.option = option
	m.option.VisibilityAnimation.Appear(time.Now())
}

func (m *Modal) Layout(gtx layout.Context) layout.Dimensions {
	if m.option.Widget == nil {
		return layout.Dimensions{}
	}
	if m.dismissWithoutAnim {
		if m.option.AfterDismiss != nil {
			m.option.AfterDismiss()
		}
		m.dismissWithoutAnim = false
		m.option = shared.ModalOption{}
		return layout.Dimensions{}
	}
	if m.btnBackdrop.Clicked(gtx) && !m.btnContent.Clicked(gtx) {
		if m.option.OnBackdropClick != nil {
			m.option.OnBackdropClick()
		} else {
			m.option.VisibilityAnimation.Disappear(gtx.Now)
		}
	}
	//var finalPosY int
	state := m.option.State
	progress := m.option.Revealed(gtx)
	backdropColor := color.NRGBA{A: uint8(200.0 * progress)}
	d := layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return m.btnBackdrop.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					switch {
					case state == component.Invisible, progress == 0, m.option.Widget == nil:
						if !m.option.Animating() {
							if m.option.AfterDismiss != nil {
								m.option.AfterDismiss()
								m.option.AfterDismiss = nil
							}
							m.option.Widget = nil
						}
						return layout.Dimensions{}
					case state == component.Visible, state == component.Appearing, state == component.Disappearing:
						gtx.Constraints.Min = gtx.Constraints.Max
						paint.Fill(gtx.Ops, backdropColor)
						defer paint.PushOpacity(gtx.Ops, progress).Pop()
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return m.btnContent.Layout(gtx, m.option.Widget)
						})
					}
					return layout.Dimensions{}
				},
			)
		}),
	)
	return d
}

func (m *Modal) DismissWithAnim() {
	m.option.Disappear(time.Now())
}

func (m *Modal) DismissWithoutAnim() {
	m.dismissWithoutAnim = true
}
