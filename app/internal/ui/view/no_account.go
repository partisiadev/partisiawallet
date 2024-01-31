package view

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/partisiadev/partisiawallet/app/internal/ui/shared"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"
	"github.com/partisiadev/partisiawallet/log"

	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"time"
)

type NoAccountView struct {
	shared.Manager
	buttonNewAccount *IconButton
	AppTheme         theme.AppTheme
	*widget.Icon
	inActiveTh      theme.AppTheme
	iconCreateNewID *widget.Icon
	AccountFormView View
	*ModalContent
}

func NewNoAccount(manager shared.Manager) *NoAccountView {
	acc := NoAccountView{AppTheme: theme.GlobalTheme, Manager: manager}
	acc.AccountFormView = NewAccountFormView(manager, acc.onSuccess)
	acc.ModalContent = NewModalContent(func() {
		acc.Modal().DismissWithAnim()
	})
	return &acc
}

func (na *NoAccountView) Layout(gtx Gtx) Dim {
	if na.AppTheme == nil {
		na.AppTheme = theme.GlobalTheme
	}
	if na.Icon == nil {
		na.Icon, _ = widget.NewIcon(icons.ActionAccountCircle)
	}
	if na.inActiveTh == nil {
		inActiveTh := theme.GlobalTheme.Clone()
		inActiveTh.Theme().ContrastBg = color.NRGBA(colornames.Grey500)
		na.inActiveTh = inActiveTh
	}
	if na.iconCreateNewID == nil {
		na.iconCreateNewID, _ = widget.NewIcon(icons.ContentCreate)
	}
	if na.buttonNewAccount == nil {
		na.buttonNewAccount = &IconButton{
			AppTheme: na.AppTheme,
			Icon:     na.Icon,
			Text:     "Add/Edit Account",
		}
	}

	flex := Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides, Alignment: layout.Middle}
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	if na.buttonNewAccount.Button.Clicked(gtx) {
		na.Manager.Modal().Show(shared.ModalOption{
			VisibilityAnimation: component.VisibilityAnimation{
				Duration: time.Millisecond * 250,
				State:    component.Invisible,
				Started:  time.Time{},
			},
			Widget: na.drawModalContent,
		})
	}
	d := flex.Layout(gtx,
		Rigid(func(gtx Gtx) Dim {
			return DrawAppImageCenter(gtx, na.AppTheme)
		}),
		Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		Rigid(func(gtx Gtx) Dim {
			return layout.Center.Layout(gtx, func(gtx Gtx) Dim {
				bdy := material.Body1(na.AppTheme.Theme(), "No Account(s) Created")
				bdy.Alignment = text.Middle
				bdy.Font.Weight = font.Black
				bdy.Color = color.NRGBA{R: 102, G: 117, B: 127, A: 255}
				return bdy.Layout(gtx)
			})
		}),
		Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		Rigid(func(gtx Gtx) Dim {
			return Flex{Spacing: layout.SpaceSides}.Layout(gtx, Rigid(func(gtx layout.Context) Dim {
				gtx.Constraints.Max.X = gtx.Dp(250)
				return na.buttonNewAccount.Layout(gtx)
			}))
		}),
		Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
	)
	return d
}

func (na *NoAccountView) drawModalContent(gtx Gtx) Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	return na.ModalContent.DrawContent(gtx, na.AppTheme, na.AccountFormView.Layout)
}

func (na *NoAccountView) onSuccess() {
	//na.Modal().DismissWithAnim()
	log.Logger().Println("After Success No Action Selected!")
}
