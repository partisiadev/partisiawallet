package newacc

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/partisiadev/partisiawallet/db"
	"github.com/partisiadev/partisiawallet/log"
	"github.com/partisiadev/partisiawallet/ui/fwk"
	"github.com/partisiadev/partisiawallet/ui/theme"
	"github.com/partisiadev/partisiawallet/ui/view"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type CreateAccount struct {
	fwk.Manager
	layout.List
	*view.PasswordForm
	btnNonCustodial view.IconButton
	*view.ImportAccountView
}

func New(m fwk.Manager) fwk.View {
	icon, _ := widget.NewIcon(icons.ActionAccountBox)
	return &CreateAccount{
		Manager: m,
		List:    layout.List{Axis: layout.Vertical},
		btnNonCustodial: view.IconButton{
			Theme: theme.GlobalTheme.Theme(),
			Text:  "Create New Account",
			Icon:  icon,
		},
		ImportAccountView: view.NewImportAccountView(),
	}
}

func (p *CreateAccount) Layout(gtx layout.Context) layout.Dimensions {
	if !db.Instance().DBAccessor().DatabaseExists() ||
		!db.Instance().DBAccessor().IsDBOpen() {
		return p.NoWalletLayout(gtx)
	}

	if p.btnNonCustodial.Button.Clicked(gtx) {
		err := db.Instance().AutoCreateEcdsaAccount()
		if err != nil {
			log.Logger().Println(err)
		}
	}
	inset := layout.UniformInset(16)
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		flex := layout.Flex{Axis: layout.Vertical}
		return flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return p.btnNonCustodial.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: 16}.Layout),
			layout.Rigid(p.drawImportKeyTextField),
		)
	})
}

func (p *CreateAccount) drawImportKeyTextField(gtx layout.Context) layout.Dimensions {
	return p.ImportAccountView.Layout(gtx)
}

func (p *CreateAccount) NoWalletLayout(gtx layout.Context) layout.Dimensions {
	if p.PasswordForm == nil {
		p.PasswordForm = view.NewPasswordForm(nil)
	}
	p.List.Axis = layout.Vertical
	return p.List.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
		return p.PasswordForm.Layout(gtx)
	})
}
