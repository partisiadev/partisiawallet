package wallet

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/db"
	"github.com/partisiadev/partisiawallet/ui/fwk"
	"github.com/partisiadev/partisiawallet/ui/page/newacc"
	"github.com/partisiadev/partisiawallet/ui/theme"
	"github.com/partisiadev/partisiawallet/ui/view"
)

type Wallet struct {
	fwk.Manager
	layout.List
	buttons []widget.Clickable
	*view.PasswordForm
	newAccountView view.View
}

func New(m fwk.Manager) fwk.View {
	return &Wallet{Manager: m, List: layout.List{Axis: layout.Vertical}}
}

func (p *Wallet) Layout(gtx layout.Context) layout.Dimensions {
	if !db.Instance().DBAccessor().DatabaseExists() ||
		!db.Instance().DBAccessor().IsDBOpen() {
		return p.NoWalletLayout(gtx)
	}

	accounts, err := db.Instance().Accounts()
	if err != nil {
		//material.Body1(
		//	theme.GlobalTheme.Theme(),
		//	fmt.Sprintf("%s", err.Error()),
		//).Layout(gtx)
	}
	if len(accounts) == 0 {
		if p.newAccountView == nil {
			p.newAccountView = newacc.New(p.Manager)
		}
		return p.newAccountView.Layout(gtx)
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

func (p *Wallet) NoWalletLayout(gtx layout.Context) layout.Dimensions {
	if p.PasswordForm == nil {
		p.PasswordForm = view.NewPasswordForm(nil)
	}
	p.List.Axis = layout.Vertical
	return p.List.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
		return p.PasswordForm.Layout(gtx)
	})
}
