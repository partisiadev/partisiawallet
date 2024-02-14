package view

import (
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/db"
	"github.com/partisiadev/partisiawallet/ui/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/image/colornames"
	"image/color"
)

type ImportAccountView struct {
	pastePvtKeyBtn  IconButton
	submitPvtKeyBtn IconButton
	clearPvtKeyBtn  IconButton
	initialized     bool
	pvtKeyStr       string
	pvtKeyImportErr error
}

func NewImportAccountView() *ImportAccountView {
	acc := ImportAccountView{}
	acc.initialize()
	return &acc
}

func (i *ImportAccountView) initialize() {
	iconAccount, _ := widget.NewIcon(icons.ContentContentPaste)
	iconSubmit, _ := widget.NewIcon(icons.ActionDone)
	iconClear, _ := widget.NewIcon(icons.ContentClear)
	i.pastePvtKeyBtn = IconButton{
		Theme: theme.GlobalTheme.Theme(),
		Icon:  iconAccount,
		Text:  `Paste From Clipboard`,
		Inset: layout.Inset{},
	}
	i.submitPvtKeyBtn = IconButton{
		Theme: theme.GlobalTheme.Theme(),
		Icon:  iconSubmit,
		Text:  `Submit Private Key`,
		Inset: layout.Inset{},
	}
	i.clearPvtKeyBtn = IconButton{
		Theme: theme.GlobalTheme.Theme(),
		Icon:  iconClear,
		Text:  `Clear Imported Key`,
		Inset: layout.Inset{},
	}
}

func (i *ImportAccountView) Layout(gtx Gtx) layout.Dimensions {
	if !i.initialized {
		i.initialize()
		i.initialized = true
	}
	flex := layout.Flex{Axis: layout.Vertical}
	spc := layout.Spacer{}
	if i.pvtKeyStr != "" {
		spc.Height = 16
	}
	i.handleEvents(gtx)
	return flex.Layout(gtx,
		layout.Rigid(i.pastePvtKeyBtn.Layout),
		layout.Rigid(spc.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if i.pvtKeyStr == "" {
				return layout.Dimensions{}
			}
			bdr := widget.Border{
				Color:        theme.GlobalTheme.Theme().ContrastBg,
				CornerRadius: 8,
				Width:        1,
			}
			return bdr.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				inset := layout.UniformInset(16)
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					flex := flex
					return flex.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if i.pvtKeyImportErr != nil {
								str := i.pvtKeyImportErr.Error()
								bdy := material.Body1(theme.GlobalTheme.Theme(), str)
								bdy.Color = color.NRGBA(colornames.Red)
								return bdy.Layout(gtx)
							}
							return material.Body1(theme.GlobalTheme.Theme(), i.pvtKeyStr).Layout(gtx)
						}),
						layout.Rigid(spc.Layout),
						layout.Rigid(i.submitPvtKeyBtn.Layout),
						layout.Rigid(spc.Layout),
						layout.Rigid(i.clearPvtKeyBtn.Layout),
					)
				})
			})
		}),
	)
}
func (i *ImportAccountView) handleEvents(gtx Gtx) {
	for _, e := range gtx.Events(&i.pastePvtKeyBtn) {
		switch e := e.(type) {
		case clipboard.Event:
			i.pvtKeyStr = e.Text
			clipboard.WriteOp{Text: ""}.Add(gtx.Ops)
			i.pvtKeyImportErr = nil
		}
	}
	if i.pastePvtKeyBtn.Button.Clicked(gtx) {
		clipboard.ReadOp{Tag: &i.pastePvtKeyBtn}.Add(gtx.Ops)
	}
	if i.clearPvtKeyBtn.Button.Clicked(gtx) {
		i.pvtKeyStr = ""
		i.pvtKeyImportErr = nil
	}
	if i.submitPvtKeyBtn.Button.Clicked(gtx) {
		if i.pvtKeyStr != "" {
			i.pvtKeyImportErr = db.Instance().ImportECDSAAccount(i.pvtKeyStr)
		}
	}
}
