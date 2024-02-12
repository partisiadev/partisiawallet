package ui

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/log"
	"github.com/partisiadev/partisiawallet/ui/fwk"
	"github.com/partisiadev/partisiawallet/ui/page/about"
	"github.com/partisiadev/partisiawallet/ui/page/chains"
	"github.com/partisiadev/partisiawallet/ui/page/wallet"
	"github.com/partisiadev/partisiawallet/ui/theme"
	"github.com/partisiadev/partisiawallet/ui/view"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type homeTabsLayout struct {
	view.Slider
	fwk.Manager
	tabItems    [3]*tabItem
	activeIndex int
}

type tabItem struct {
	URL string
	widget.Clickable
	fwk.View
	title  string
	parent *homeTabsLayout
	icon   *widget.Icon
	index  int
}

func (t *tabItem) Layout(gtx layout.Context) layout.Dimensions {
	if t.Clickable.Clicked(gtx) {
		if t.index < t.parent.activeIndex {
			t.parent.Slider.PushRight()
		} else if t.index > t.parent.activeIndex {
			t.parent.Slider.PushLeft()
		}
		t.parent.activeIndex = t.index
		t.parent.Navigator().NavigateTo(t.URL)
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

func newHomeTabsManager(m *manager) *homeTabsLayout {
	walletIcon, _ := widget.NewIcon(icons.ActionAccountBalanceWallet)
	chainsIcon, _ := widget.NewIcon(icons.ContentLink)
	aboutIcon, _ := widget.NewIcon(icons.ActionInfo)
	tabItems := [3]*tabItem{
		{URL: "/" + fwk.HomeTabsPageName + "/" + fwk.WalletPageName, title: "Wallet", View: wallet.New(m), icon: walletIcon, index: 0},
		{URL: "/" + fwk.HomeTabsPageName + "/" + fwk.ChainsPageName, title: "Chains", View: chains.New(m), icon: chainsIcon, index: 1},
		{URL: "/" + fwk.HomeTabsPageName + "/" + fwk.AboutPageName, title: "About", View: about.New(m), icon: aboutIcon, index: 2},
	}
	hmTabs := &homeTabsLayout{Manager: m, tabItems: tabItems}
	hmTabs.tabItems[0].parent = hmTabs
	hmTabs.tabItems[1].parent = hmTabs
	hmTabs.tabItems[2].parent = hmTabs
	hmTabs.Navigator().Register(fmt.Sprintf(`^/%s`, fwk.HomeTabsPageName), hmTabs)
	hmTabs.Navigator().NavigateTo(tabItems[0].URL)
	return hmTabs
}

func (hm *homeTabsLayout) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = gtx.Constraints.Max
			return hm.Slider.Layout(gtx, hm.tabItems[hm.activeIndex].View.Layout)
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

func (hm *homeTabsLayout) Handle(concretePath string) fwk.View {
	log.Logger().Println(concretePath)
	return hm
}
