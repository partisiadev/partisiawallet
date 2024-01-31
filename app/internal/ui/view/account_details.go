package view

import (
	"errors"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/partisiadev/partisiawallet/app/internal/state/wallet"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"strings"
)

type AccountDetails struct {
	AppTheme                theme.AppTheme
	buttonCopyPvtKey        IconButton
	buttonCopyPubKey        IconButton
	buttonPrivateKeyVisible IconButton
	buttonPrivateKeyHidden  IconButton
	inputPassword           *component.TextField
	Account                 wallet.Account
	inputPasswordStr        string
	pvtKeyStr               string
	pvtKeyListLayout        layout.List
	pubKeyListLayout        layout.List
}

func NewAccountDetails(account wallet.Account) *AccountDetails {
	iconCopy, _ := widget.NewIcon(icons.ContentContentCopy)
	iconVisible, _ := widget.NewIcon(icons.ActionVisibility)
	iconHidden, _ := widget.NewIcon(icons.ActionVisibilityOff)
	accountDetails := AccountDetails{
		AppTheme:      theme.GlobalTheme,
		Account:       account,
		inputPassword: &component.TextField{Editor: widget.Editor{SingleLine: false}},
		buttonCopyPvtKey: IconButton{
			AppTheme: theme.GlobalTheme,
			Icon:     iconCopy,
			Text:     "Copy Private Key",
		},
		buttonCopyPubKey: IconButton{
			AppTheme: theme.GlobalTheme,
			Icon:     iconCopy,
			Text:     "Copy Public Key",
		},
		buttonPrivateKeyVisible: IconButton{
			AppTheme: theme.GlobalTheme,
			Icon:     iconVisible,
			Text:     "Hide Private Key",
		},
		buttonPrivateKeyHidden: IconButton{
			AppTheme: theme.GlobalTheme,
			Icon:     iconHidden,
			Text:     "Show Private Key",
		},
	}

	return &accountDetails
}

func (ad *AccountDetails) Layout(gtx Gtx) Dim {
	if ad.AppTheme == nil {
		ad.AppTheme =
			theme.GlobalTheme
	}
	if ad.inputPassword.Text() != ad.inputPasswordStr {
		ad.inputPassword.ClearError()
	}
	ad.inputPasswordStr = ad.inputPassword.Text()

	inset := layout.UniformInset(unit.Dp(16))
	flex := Flex{Axis: layout.Vertical, Alignment: layout.Start}
	d := flex.Layout(gtx,
		Rigid(func(gtx Gtx) Dim {
			inset := inset
			return inset.Layout(gtx, ad.drawPasswordField)
		}),
		Rigid(func(gtx Gtx) Dim {
			inset := inset
			return inset.Layout(gtx, ad.drawPvtKeyField)
		}),
		Rigid(func(gtx Gtx) Dim {
			inset := inset
			return inset.Layout(gtx, ad.drawPubKeyField)
		}),
	)
	return d
}

func (ad *AccountDetails) drawPasswordField(gtx Gtx) Dim {
	if ad.buttonPrivateKeyHidden.Button.Clicked(gtx) {
		var err error
		if strings.TrimSpace(ad.inputPassword.Text()) == "" {
			err = errors.New("password is empty")
			ad.inputPassword.SetError(err.Error())
			ad.pvtKeyStr = ""
		} else {
			err = wallet.GlobalWallet.VerifyPassword(ad.inputPasswordStr)
			if err != nil {
				ad.pvtKeyStr = ""
				ad.inputPassword.SetError(err.Error())
			} else {
				ad.pvtKeyStr = ad.Account.PrivateKey
			}
		}
	}
	if ad.buttonPrivateKeyVisible.Button.Clicked(gtx) {
		ad.pvtKeyStr = ""
	}
	labelPasswordText := "Enter Password"
	flex := Flex{Axis: layout.Vertical}
	return flex.Layout(gtx,
		Rigid(func(gtx Gtx) Dim {
			th := *ad.AppTheme.Theme()
			origSize := th.TextSize
			if strings.TrimSpace(ad.inputPassword.Text()) == "" && !ad.inputPassword.Focused() {
				th.TextSize = unit.Sp(12)
			} else {
				th.TextSize = origSize
			}
			return ad.inputPassword.Layout(gtx, &th, labelPasswordText)
		}),
		Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		Rigid(func(gtx Gtx) Dim {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			btn := &ad.buttonPrivateKeyHidden
			if ad.pvtKeyStr != "" {
				btn = &ad.buttonPrivateKeyVisible
			}
			return btn.Layout(gtx)
		}),
	)
}

func (ad *AccountDetails) drawPvtKeyField(gtx Gtx) Dim {
	if ad.buttonCopyPvtKey.Button.Clicked(gtx) {
		//manager.Window().WriteClipboard(ad.pvtKeyStr)
	}
	flex := Flex{Axis: layout.Vertical}
	return flex.Layout(gtx,
		Rigid(func(gtx Gtx) Dim {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			var txt string
			txt = strings.TrimSpace(ad.pvtKeyStr)
			txtColor := ad.AppTheme.Theme().Fg
			if txt == "" {
				txt = "Your Private Key"
				txtColor = color.NRGBA(colornames.Grey500)
			}
			inset := layout.UniformInset(unit.Dp(16))
			mac := op.Record(gtx.Ops)
			d := inset.Layout(gtx,
				func(gtx Gtx) Dim {
					lbl := material.Label(ad.AppTheme.Theme(), ad.AppTheme.Theme().TextSize, txt)
					lbl.Color = txtColor
					return ad.pvtKeyListLayout.Layout(gtx, 1, func(gtx layout.Context, index int) Dim {
						return lbl.Layout(gtx)
					})
				})
			stop := mac.Stop()
			bounds := image.Rect(0, 0, d.Size.X, d.Size.Y)
			rect := clip.UniformRRect(bounds, gtx.Dp(4))
			paint.FillShape(gtx.Ops,
				ad.AppTheme.Theme().Fg,
				clip.Stroke{Path: rect.Path(gtx.Ops), Width: float32(gtx.Dp(1))}.Op(),
			)
			stop.Add(gtx.Ops)
			return d
		}),
		Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		Rigid(func(gtx Gtx) Dim {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return ad.buttonCopyPvtKey.Layout(gtx)
		}),
	)
}

func (ad *AccountDetails) drawPubKeyField(gtx Gtx) Dim {
	publicKey := ad.Account.PublicKey
	if ad.buttonCopyPubKey.Button.Clicked(gtx) {
		//manager.Window().WriteClipboard(publicKey)
	}
	flex := Flex{Axis: layout.Vertical}
	return flex.Layout(gtx,
		Rigid(func(gtx Gtx) Dim {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			var txt string
			txt = publicKey
			txtColor := ad.AppTheme.Theme().Fg
			if txt == "" {
				txt = "Your Public Key"
				txtColor = color.NRGBA(colornames.Grey500)
			}
			inset := layout.UniformInset(unit.Dp(16))
			mac := op.Record(gtx.Ops)
			d := inset.Layout(gtx,
				func(gtx Gtx) Dim {
					lbl := material.Label(ad.AppTheme.Theme(), ad.AppTheme.Theme().TextSize, txt)
					lbl.Color = txtColor
					return ad.pubKeyListLayout.Layout(gtx, 1, func(gtx layout.Context, index int) Dim {
						return lbl.Layout(gtx)
					})
				})
			stop := mac.Stop()
			bounds := image.Rect(0, 0, d.Size.X, d.Size.Y)
			rect := clip.UniformRRect(bounds, gtx.Dp(4))
			paint.FillShape(gtx.Ops,
				ad.AppTheme.Theme().Fg,
				clip.Stroke{Path: rect.Path(gtx.Ops), Width: float32(gtx.Dp(1))}.Op(),
			)
			stop.Add(gtx.Ops)
			return d
		}),
		Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		Rigid(func(gtx Gtx) Dim {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return ad.buttonCopyPubKey.Layout(gtx)
		}),
	)
}
