package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/ui/fwk"
	"github.com/partisiadev/partisiawallet/ui/page/about"
	"github.com/partisiadev/partisiawallet/ui/page/chains"
	"github.com/partisiadev/partisiawallet/ui/page/newacc"
	"github.com/partisiadev/partisiawallet/ui/page/wallet"
	"github.com/partisiadev/partisiawallet/ui/theme"
	"github.com/partisiadev/partisiawallet/ui/view"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"net/url"
	"strings"
)

type AppLayout struct {
	view.Slider
	fwk.Manager
	tabItems    [3]*tabItem
	view        view.View
	activeIndex int
}

func NewAppLayout(m *manager) *AppLayout {
	appLayout := &AppLayout{Manager: m}
	walletIcon, _ := widget.NewIcon(icons.ActionAccountBalanceWallet)
	chainsIcon, _ := widget.NewIcon(icons.ContentLink)
	aboutIcon, _ := widget.NewIcon(icons.ActionInfo)
	tabItems := [3]*tabItem{
		{title: "Wallet", icon: walletIcon, parent: appLayout, index: 0, path: `/homeTabs/wallet`, view: wallet.New(m)},
		{title: "Chains", icon: chainsIcon, parent: appLayout, index: 1, path: `/homeTabs/chains`, view: chains.New(m)},
		{title: "About", icon: aboutIcon, parent: appLayout, index: 2, path: `/homeTabs/about`, view: about.New(m)},
	}
	appLayout.tabItems = tabItems
	appLayout.Nav().Register(`^/wallet`, appLayout)
	appLayout.Nav().Register(`^/chains`, appLayout)
	appLayout.Nav().Register(`^/about`, appLayout)
	appLayout.Nav().Register(`^/createAccount`, appLayout)
	appLayout.Nav().Register(`^/homeTabs/wallet`, appLayout)
	appLayout.Nav().Register(`^/homeTabs/chains`, appLayout)
	appLayout.Nav().Register(`^/homeTabs/about`, appLayout)
	appLayout.Nav().Register(`^/homeTabs/createAccount`, appLayout)
	return appLayout
}

func (hm *AppLayout) Handle(path *url.URL) fwk.View {
	var v fwk.View
	if path == nil {
		return v
	}
	pth := path.Path
	switch {
	case strings.HasPrefix(pth, `/homeTabs/wallet`):
		hm.activeIndex = 0
		hm.view = hm.tabItems[0].view
		hm.tabItems[0].path = pth
		return hm
	case strings.HasPrefix(pth, `/homeTabs/chains`):
		hm.activeIndex = 1
		hm.view = hm.tabItems[1].view
		hm.tabItems[1].path = pth
		return hm
	case strings.HasPrefix(pth, `/homeTabs/about`):
		hm.activeIndex = 2
		hm.view = hm.tabItems[2].view
		hm.tabItems[2].path = pth
		return hm
	case strings.HasPrefix(pth, `/homeTabs/createAccount`):
		hm.view = newacc.New(hm)
		return hm
	}
	return v
}

func (hm *AppLayout) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = gtx.Constraints.Max
			if hm.view == nil {
				return layout.Dimensions{}
			}
			return hm.Slider.Layout(gtx, hm.view.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			inset := layout.Inset{Top: 0, Bottom: 0}
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Spacing: layout.SpaceEvenly, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return hm.tabItems[0].Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return hm.tabItems[1].Layout(gtx)
					}), layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return hm.tabItems[2].Layout(gtx)
					}),
				)
			})
		}),
	)
}

type tabItem struct {
	widget.Clickable
	title  string
	icon   *widget.Icon
	index  int
	parent *AppLayout
	view   view.View
	path   string
}

func (t *tabItem) Layout(gtx layout.Context) layout.Dimensions {
	if t.Clickable.Clicked(gtx) {
		if t.index < t.parent.activeIndex {
			t.parent.Slider.PushRight()
			t.parent.Nav().NavigateTo(t.path)
		} else if t.index > t.parent.activeIndex {
			t.parent.Slider.PushLeft()
			t.parent.Nav().NavigateTo(t.path)
		}
	}

	btnStyle := material.ButtonLayoutStyle{Button: &t.Clickable}
	iconColor := theme.GlobalTheme.Theme().ContrastBg
	btnStyle.Background = theme.GlobalTheme.Theme().ContrastFg
	txtColor := theme.GlobalTheme.Theme().ContrastBg
	if t.parent.activeIndex == t.index {
		iconColor, btnStyle.Background, txtColor = btnStyle.Background, iconColor, btnStyle.Background
	}
	return btnStyle.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		inset := layout.Inset{Top: 4, Bottom: 4}
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			flex := layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceSides}
			return flex.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					d := t.icon.Layout(gtx, iconColor)
					return d
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(theme.GlobalTheme.Theme(), 12, t.title)
					lbl.Color = txtColor
					return lbl.Layout(gtx)
				}),
			)
		})
	})
}
