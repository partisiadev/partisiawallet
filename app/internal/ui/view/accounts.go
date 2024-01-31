package view

import (
	"bytes"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/app/assets"
	"github.com/partisiadev/partisiawallet/app/internal/state/wallet"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"
	"image"
)

type accountsView struct {
	layout.List
	AppTheme              theme.AppTheme
	title                 string
	accountsItems         []*accountsItem
	currentAccountLayout  layout.List
	enum                  widget.Enum
	accountChangeCallback func()
}

func NewAccountsView(accountChangeCallback func()) View {
	p := accountsView{
		AppTheme:              theme.GlobalTheme,
		title:                 "Accounts",
		List:                  layout.List{Axis: layout.Vertical},
		accountsItems:         []*accountsItem{},
		accountChangeCallback: accountChangeCallback,
	}
	return &p
}

func (p *accountsView) Layout(gtx Gtx) Dim {
	a, _ := wallet.GlobalWallet.Account()
	p.enum.Value = a.PublicKey
	flex := Flex{Axis: layout.Vertical,
		Spacing:   layout.SpaceEnd,
		Alignment: layout.Start,
	}

	d := flex.Layout(gtx,
		Rigid(p.drawIdentitiesItems),
	)

	return d
}

func (p *accountsView) drawIdentitiesItems(gtx Gtx) Dim {
	if p.isProcessingRequired() {
		accs, _ := wallet.GlobalWallet.Accounts()
		p.accountsItems = make([]*accountsItem, 0, len(accs))
		for _, userID := range accs {
			_ = userID
			p.accountsItems = append(p.accountsItems, &accountsItem{
				AppTheme: p.AppTheme,
				//Account: userID,
				Enum: &p.enum,
			})
		}
	}
	return Flex{Axis: layout.Vertical}.Layout(gtx,
		Rigid(func(gtx Gtx) Dim {
			inset := layout.UniformInset(unit.Dp(16))
			return inset.Layout(gtx, func(gtx layout.Context) Dim {
				flex := Flex{Alignment: layout.Middle}
				a, _ := wallet.GlobalWallet.Account()
				d := flex.Layout(gtx,
					Rigid(func(gtx layout.Context) Dim {
						var img image.Image
						img, _, _ = image.Decode(bytes.NewReader(a.PublicImage))
						if img == nil {
							img = assets.AppIconImage
						}
						radii := gtx.Dp(24)
						gtx.Constraints.Max.X, gtx.Constraints.Max.Y = radii*2, radii*2
						bounds := image.Rect(0, 0, radii*2, radii*2)
						clipOp := clip.UniformRRect(bounds, radii).Push(gtx.Ops)
						imgOps := paint.NewImageOp(img)
						imgWidget := widget.Image{Src: imgOps, Fit: widget.Contain, Position: layout.Center, Scale: 0}
						d := imgWidget.Layout(gtx)
						clipOp.Pop()
						return d
					}),
					Rigid(func(gtx Gtx) Dim {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return p.currentAccountLayout.Layout(gtx, 1, func(gtx layout.Context, index int) Dim {
							flex := Flex{Spacing: layout.SpaceSides, Alignment: layout.Start, Axis: layout.Vertical}
							inset := layout.Inset{Right: unit.Dp(8), Left: unit.Dp(8)}
							d := inset.Layout(gtx, func(gtx Gtx) Dim {
								d := flex.Layout(gtx,
									Rigid(func(gtx Gtx) Dim {
										b := material.Body1(p.AppTheme.Theme(), a.PublicKey)
										b.Font.Weight = font.Bold
										return b.Layout(gtx)
									}),
									//Rigid(func(gtx Gtx) Dim {
									//	b := material.Body1(p.AppTheme(), strings.Trim(string(p.currentAccount.Contents), "\n"))
									//	b.Color = color.NRGBA(colornames.Grey600)
									//	return b.Layout(gtx)
									//}),
								)
								return d
							})
							return d
						})
					}),
				)
				return d
			})
		}),
		Rigid(func(gtx layout.Context) Dim {
			return p.List.Layout(gtx, len(p.accountsItems), func(gtx Gtx, index int) (d Dim) {
				accountItem := p.accountsItems[index]
				if accountItem.Clickable.Pressed() {
					//_ = chains.GlobalWallet.AddUpdateAccount(&accountItem.Account)
					//if p.accountChangeCallback != nil {
					//	p.accountChangeCallback()
					//}
				}
				return p.accountsItems[index].Layout(gtx)
			})
		}),
	)
}

// isProcessingRequired
func (p *accountsView) isProcessingRequired() bool {
	accs, _ := wallet.GlobalWallet.Accounts()
	isRequired := len(accs) != len(p.accountsItems)
	return isRequired
}
