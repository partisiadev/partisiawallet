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

	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"time"
)

type NoContactView struct {
	shared.Manager
	buttonAddContact *IconButton
	AppTheme         theme.AppTheme
	*widget.Icon
	ContactFormView View
	*ModalContent
}

func NewNoContact(manager shared.Manager, contactForm View, onSuccess func(contactAddr string), btnText string) *NoContactView {
	btnIcon, _ := widget.NewIcon(icons.CommunicationContacts)
	if btnText == "" {
		btnText = "Add Contact"
	}
	nc := NoContactView{
		Manager:  manager,
		AppTheme: theme.GlobalTheme,
		buttonAddContact: &IconButton{
			AppTheme: theme.GlobalTheme,
			Icon:     btnIcon,
			Text:     btnText,
		},
		ContactFormView: contactForm,
	}
	nc.ModalContent = NewModalContent(func() { nc.Modal().DismissWithAnim() })
	return &nc
}

func (nc *NoContactView) Layout(gtx Gtx) Dim {
	flex := Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides, Alignment: layout.Middle}
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	if nc.buttonAddContact.Button.Clicked(gtx) {
		nc.Modal().Show(shared.ModalOption{
			Widget: nc.drawModalContent,
			VisibilityAnimation: Animation{
				Duration: time.Millisecond * 250,
				State:    component.Invisible,
				Started:  time.Time{},
			},
		})
	}
	d := flex.Layout(gtx,
		Rigid(func(gtx Gtx) Dim {
			return DrawAppImageCenter(gtx, nc.AppTheme)
		}),
		Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		Rigid(func(gtx Gtx) Dim {
			return layout.Center.Layout(gtx, func(gtx Gtx) Dim {
				bdy := material.Body1(nc.AppTheme.Theme(), "No Contact(s) Found")
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
				return nc.buttonAddContact.Layout(gtx)
			}))
		}),
		Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
	)
	return d
}

func (nc *NoContactView) drawModalContent(gtx Gtx) Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	return nc.ModalContent.DrawContent(gtx, nc.AppTheme, nc.ContactFormView.Layout)
}
