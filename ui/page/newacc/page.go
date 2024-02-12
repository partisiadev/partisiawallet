package newacc

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/db"
	"github.com/partisiadev/partisiawallet/ui/fwk"
	"github.com/partisiadev/partisiawallet/ui/theme"
	"github.com/partisiadev/partisiawallet/ui/view"
)

type CreateAccountView struct {
	fwk.Manager
	layout.List
	*view.PasswordForm
}

func New(m fwk.Manager) fwk.View {
	return &CreateAccountView{Manager: m, List: layout.List{Axis: layout.Vertical}}
}

func (p *CreateAccountView) Layout(gtx layout.Context) layout.Dimensions {
	if !db.Instance().DBAccessor().DatabaseExists() ||
		!db.Instance().DBAccessor().IsDBOpen() {
		return p.NoWalletLayout(gtx)
	}

	return material.Body1(
		theme.GlobalTheme.Theme(),
		"Welcome to create accounts page",
	).Layout(gtx)

	//return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
	//	layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
	//		return p.List.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
	//			return material.Body1(
	//				theme.GlobalTheme.Theme(),
	//				"Welcome to create accounts page",
	//			).Layout(gtx)
	//		})
	//	}),
	//)
}

func (p *CreateAccountView) NoWalletLayout(gtx layout.Context) layout.Dimensions {
	if p.PasswordForm == nil {
		p.PasswordForm = view.NewPasswordForm(nil)
	}
	p.List.Axis = layout.Vertical
	return p.List.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
		return p.PasswordForm.Layout(gtx)
	})
}
