package view

import (
	"gioui.org/font"
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/app/internal/state/wallet"
	"github.com/partisiadev/partisiawallet/app/internal/ui/shared"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"

	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"strings"
)

type accountForm struct {
	shared.Manager
	AppTheme              theme.AppTheme
	InActiveTheme         theme.AppTheme
	iconCreateNewID       *widget.Icon
	iconImportFile        *widget.Icon
	pvtKeyStr             string
	title                 string
	importLabelText       string
	btnClear              IconButton
	btnNewID              IconButton
	btnSubmitImportKey    IconButton
	btnPasteKey           IconButton
	navigationIcon        *widget.Icon
	iDDetailsView         AccountDetails
	errorCreateNewID      error
	errorImportKey        error
	creatingNewID         bool
	submittingImportedKey bool
	OnSuccess             func()
	*ModalContent
}

func NewAccountFormView(manager shared.Manager, onSuccess func()) View {
	clearIcon, _ := widget.NewIcon(icons.ContentClear)
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	iconCreateNewID, _ := widget.NewIcon(icons.ActionDone)
	iconImportFile, _ := widget.NewIcon(icons.FileFileUpload)
	pasteIcon, _ := widget.NewIcon(icons.ContentContentPaste)
	errorTh := *theme.GlobalTheme.Theme()
	th := theme.GlobalTheme
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	inActiveTh := theme.GlobalTheme.Clone()
	inActiveTh.Theme().ContrastBg = color.NRGBA(colornames.Grey500)
	s := accountForm{
		Manager:         manager,
		AppTheme:        theme.GlobalTheme,
		InActiveTheme:   inActiveTh,
		title:           "Account",
		navigationIcon:  navIcon,
		iconCreateNewID: iconCreateNewID,
		iconImportFile:  iconImportFile,
		importLabelText: "Import Key",
		OnSuccess:       onSuccess,
		btnSubmitImportKey: IconButton{
			AppTheme: theme.GlobalTheme,
			Icon:     iconCreateNewID,
			Text:     "Submit",
		},
		btnPasteKey: IconButton{
			AppTheme: th,
			Icon:     pasteIcon,
			Text:     "Paste",
		},
		btnNewID: IconButton{
			AppTheme: th,
			Icon:     iconCreateNewID,
			Text:     "Auto Create New Account",
		},
		btnClear: IconButton{
			AppTheme: th,
			Icon:     clearIcon,
			Text:     "Clear",
		},
		iDDetailsView: AccountDetails{
			AppTheme: th,
		},
	}
	s.ModalContent = NewModalContent(func() {
		s.Modal().DismissWithAnim()
		s.creatingNewID = false
		s.submittingImportedKey = false
		if s.OnSuccess != nil {
			s.OnSuccess()
		}
	})
	return &s
}

func (p *accountForm) Layout(gtx Gtx) Dim {
	if p.AppTheme == nil {
		p.AppTheme = theme.GlobalTheme
	}

	inset := layout.UniformInset(unit.Dp(16))
	flex := Flex{Axis: layout.Vertical, Alignment: layout.Start}
	d := flex.Layout(gtx,
		Rigid(func(gtx Gtx) Dim {
			inset := inset
			return inset.Layout(gtx, p.drawImportKeyTextField)
		}),
		Rigid(func(gtx Gtx) Dim {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			bd := material.Body1(p.AppTheme.Theme(), "Or")
			bd.Font.Weight = font.Bold
			bd.Alignment = text.Middle
			bd.TextSize = unit.Sp(20)
			return bd.Layout(gtx)
		}),
		Rigid(func(gtx Gtx) Dim {
			inset := inset
			return inset.Layout(gtx, p.drawAutoCreateField)
		}),
	)
	if p.creatingNewID || p.submittingImportedKey {
		layout.Stack{}.Layout(gtx,
			Stacked(func(gtx layout.Context) Dim {
				loader := Loader{}
				gtx.Constraints.Max, gtx.Constraints.Min = d.Size, d.Size
				return Flex{Alignment: layout.Middle, Spacing: layout.SpaceSides}.Layout(gtx,
					Rigid(func(gtx Gtx) Dim {
						return loader.Layout(gtx)
					}))
			}),
		)
		return d
	}
	return d
}

func (p *accountForm) drawImportKeyTextField(gtx Gtx) Dim {

	for _, e := range gtx.Events(&p.btnPasteKey) {
		switch e := e.(type) {
		case clipboard.Event:
			_ = e
			p.pvtKeyStr = e.Text
			// Clear the clipboard
			clipboard.WriteOp{Text: ""}.Add(gtx.Ops)
			p.errorImportKey = nil
		}
	}
	if p.btnPasteKey.Button.Clicked(gtx) {
		p.btnPasteKey.Button.Focus()
		clipboard.ReadOp{Tag: &p.btnPasteKey}.Add(gtx.Ops)
	}

	if p.btnClear.Button.Clicked(gtx) {
		p.pvtKeyStr = ""
		p.errorImportKey = nil
	}

	if p.btnSubmitImportKey.Button.Clicked(gtx) && !p.submittingImportedKey {
		p.submittingImportedKey = true
		p.createAccountFromPvtKeyHexStr()
	}
	flex := Flex{Axis: layout.Vertical}
	return flex.Layout(gtx,
		Rigid(func(gtx Gtx) Dim {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			var txt string
			txt = strings.TrimSpace(p.pvtKeyStr)
			txtColor := p.AppTheme.Theme().Fg
			if txt == "" {
				txt = "Paste key file contents here"
				txtColor = color.NRGBA(colornames.Grey500)
			}
			if p.errorImportKey != nil {
				txt = p.errorImportKey.Error()
				txtColor = color.NRGBA(colornames.Red500)
			}
			inset := layout.UniformInset(unit.Dp(16))
			mac := op.Record(gtx.Ops)
			d := inset.Layout(gtx,
				func(gtx layout.Context) Dim {
					lbl := material.Label(p.AppTheme.Theme(), p.AppTheme.Theme().TextSize, txt)
					lbl.MaxLines = 10
					lbl.Color = txtColor
					return lbl.Layout(gtx)
				})
			stop := mac.Stop()
			bounds := image.Rect(0, 0, d.Size.X, d.Size.Y)
			rect := clip.UniformRRect(bounds, gtx.Dp(4))
			paint.FillShape(gtx.Ops,
				p.AppTheme.Theme().Fg,
				clip.Stroke{Path: rect.Path(gtx.Ops), Width: float32(gtx.Dp(1))}.Op(),
			)
			stop.Add(gtx.Ops)
			return d
		}),
		Rigid(func(gtx layout.Context) Dim {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			mobileWidth := gtx.Dp(350)
			flex := Flex{Spacing: layout.SpaceBetween}
			spacerLayout := layout.Spacer{Width: unit.Dp(16)}
			submitLayout := Flexed(1, func(gtx layout.Context) Dim {
				return p.btnSubmitImportKey.Layout(gtx)
			})
			pasteLayout := Flexed(1, func(gtx Gtx) Dim {
				return p.btnPasteKey.Layout(gtx)
			})
			clearLayout := Flexed(1, func(gtx Gtx) Dim {
				return p.btnClear.Layout(gtx)
			})
			if gtx.Constraints.Max.X <= mobileWidth {
				flex.Axis = layout.Vertical
				spacerLayout.Width = 0
				spacerLayout.Height = 8
				submitLayout = Rigid(func(gtx Gtx) Dim {
					return p.btnSubmitImportKey.Layout(gtx)
				})
				pasteLayout = Rigid(func(gtx Gtx) Dim {
					return p.btnPasteKey.Layout(gtx)
				})
				clearLayout = Rigid(func(gtx Gtx) Dim {
					return p.btnClear.Layout(gtx)
				})
			}
			inset := layout.Inset{Top: unit.Dp(16)}
			return inset.Layout(gtx, func(gtx layout.Context) Dim {
				return flex.Layout(gtx,
					submitLayout,
					Rigid(spacerLayout.Layout),
					pasteLayout,
					Rigid(spacerLayout.Layout),
					clearLayout,
				)
			})
		}),
	)

}

func (p *accountForm) drawAutoCreateField(gtx Gtx) Dim {
	var button *IconButton
	if p.errorCreateNewID != nil {
		button = &IconButton{
			AppTheme: p.InActiveTheme,
			Icon:     p.iconCreateNewID,
			Text:     "Auto Create New Account",
		}
	} else {
		button = &p.btnNewID
	}
	if p.btnNewID.Button.Clicked(gtx) && !p.creatingNewID {
		p.creatingNewID = true
		p.autoCreateAccount()
	}
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	flex := Flex{Axis: layout.Vertical}
	return flex.Layout(gtx,
		Rigid(func(gtx Gtx) Dim {
			flex := Flex{Spacing: layout.SpaceEnd}
			inset := layout.Inset{Top: unit.Dp(16)}
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return inset.Layout(gtx, func(gtx Gtx) Dim {
				return flex.Layout(gtx,
					Rigid(func(gtx Gtx) Dim {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return button.Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (p *accountForm) createAccountFromPvtKeyHexStr() {
	p.submittingImportedKey = true
	go func() {
		p.errorImportKey = wallet.GlobalWallet.CreateAccount(p.pvtKeyStr)
		p.submittingImportedKey = false
		if p.errorImportKey == nil {
			//p.account.PrivateKey = p.Service().Account().PrivateKey
			if p.OnSuccess != nil {
				p.OnSuccess()
			}
		}
		p.Window().Invalidate()
	}()
}

func (p *accountForm) autoCreateAccount() {
	p.creatingNewID = true
	go func() {
		p.errorCreateNewID = wallet.GlobalWallet.AutoCreateAccount()
		p.creatingNewID = false
		if p.errorCreateNewID == nil {
			if p.OnSuccess != nil {
				p.OnSuccess()
			}
		} else {
			//p.Snackbar().Show(p.errorCreateNewID.Error(), &widget.Clickable{}, color.NRGBA{}, "")
			p.errorCreateNewID = nil
		}
		p.Window().Invalidate()
	}()
}
