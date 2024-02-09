package wallet

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/partisiadev/partisiawallet/db"
	"github.com/partisiadev/partisiawallet/ui/shared"
	"github.com/partisiadev/partisiawallet/ui/theme"
	"github.com/partisiadev/partisiawallet/ui/view"
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
	buttons []widget.Clickable
	*view.PasswordForm
}

func New(m shared.Manager) shared.View {
	return &Wallet{Manager: m, List: layout.List{Axis: layout.Vertical}}
}

func (p *Wallet) Layout(gtx Gtx) Dim {
	if !db.Instance().DBAccessor().DatabaseExists() ||
		!db.Instance().DBAccessor().IsDBOpen() {
		return p.NoWalletLayout(gtx)
	}

	accounts, err := db.Instance().Accounts()
	if err != nil {
		material.Body1(
			theme.GlobalTheme.Theme(),
			fmt.Sprintf("%s", err.Error()),
		).Layout(gtx)
	}
	if len(accounts) == 0 {
		material.Body1(
			theme.GlobalTheme.Theme(),
			"There are no accounts to view",
		).Layout(gtx)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return p.List.Layout(gtx, len(accounts), func(gtx layout.Context, index int) layout.Dimensions {
				if len(p.buttons) <= index {
					buttons := make([]widget.Clickable, index-len(p.buttons)+1)
					p.buttons = append(p.buttons, buttons...)
				}
				if p.buttons[index].Clicked(gtx) {
					_ = db.Instance().SetActiveAccount(accounts[index])
					op.InvalidateOp{}.Add(gtx.Ops)
				}
				return p.buttons[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.Body1(
						theme.GlobalTheme.Theme(),
						fmt.Sprintf("%s", accounts[index].PathID()),
					).Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var txt string
			acc, err := db.Instance().ActiveAccount()
			if err != nil {
				txt = err.Error()
			} else {
				txt = acc.PathID()
			}
			return material.Body1(
				theme.GlobalTheme.Theme(),
				fmt.Sprintf("%s", txt),
			).Layout(gtx)
		}),
	)
}

func (p *Wallet) NoWalletLayout(gtx Gtx) layout.Dimensions {
	if p.PasswordForm == nil {
		p.PasswordForm = view.NewPasswordForm(nil)
	}
	p.List.Axis = layout.Vertical
	return p.List.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
		return p.PasswordForm.Layout(gtx)
	})
}
