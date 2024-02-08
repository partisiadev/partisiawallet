package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/log"
	"github.com/partisiadev/partisiawallet/ui/page/about"
	"github.com/partisiadev/partisiawallet/ui/page/chains"
	"github.com/partisiadev/partisiawallet/ui/page/wallet"
	"github.com/partisiadev/partisiawallet/ui/shared"
	"github.com/partisiadev/partisiawallet/ui/theme"
	"github.com/partisiadev/partisiawallet/ui/view"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type homeTabsLayout struct {
	view.Slider
	shared.Manager
	tabItems    [3]*tabItem
	pagesFound  map[string]struct{}
	activeIndex int
}

type tabItem struct {
	m   *manager
	URL string
	widget.Clickable
	shared.View
	title  string
	parent *homeTabsLayout
	icon   *widget.Icon
}

func (t *tabItem) Layout(gtx layout.Context, index int) layout.Dimensions {
	if t.Clickable.Clicked(gtx) {
		if index < t.parent.activeIndex {
			t.parent.Slider.PushRight()
		} else if index > t.parent.activeIndex {
			t.parent.Slider.PushLeft()
		}
		t.parent.activeIndex = index
	}

	btnStyle := material.ButtonLayoutStyle{Button: &t.Clickable}
	iconColor := theme.GlobalTheme.Theme().ContrastBg
	btnStyle.Background = theme.GlobalTheme.Theme().ContrastFg
	txtColor := theme.GlobalTheme.Theme().ContrastBg
	if t.parent.activeIndex == index {
		iconColor, btnStyle.Background, txtColor = btnStyle.Background, iconColor, btnStyle.Background
	}
	return btnStyle.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		inset := layout.Inset{Top: 4, Bottom: 4}
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			flex := layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceSides}
			return flex.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					d := t.icon.Layout(gtx, iconColor)
					log.Logger().Println(d.Size)
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

func homeTabsManager(m *manager) *homeTabsLayout {
	walletIcon, _ := widget.NewIcon(icons.ActionAccountBalanceWallet)
	chainsIcon, _ := widget.NewIcon(icons.ContentLink)
	aboutIcon, _ := widget.NewIcon(icons.ActionInfo)
	hmTabs := &homeTabsLayout{Manager: m, tabItems: [3]*tabItem{
		{URL: shared.WalletPagePattern, m: m, title: "Wallet", View: wallet.New(m), icon: walletIcon},
		{URL: shared.ChainsPagePattern, m: m, title: "Chains", View: chains.New(m), icon: chainsIcon},
		{URL: shared.AboutPagePattern, m: m, title: "About", View: about.New(m), icon: aboutIcon},
	}}
	hmTabs.tabItems[0].parent = hmTabs
	hmTabs.tabItems[1].parent = hmTabs
	hmTabs.tabItems[2].parent = hmTabs
	m.Nav().Register(shared.WalletPagePattern, func(concretePath string) shared.View {
		hmTabs.tabItems[0].View = wallet.New(m)
		return hmTabs
	})
	m.Nav().Register(shared.ChainsPagePattern, func(concretePath string) shared.View {
		hmTabs.tabItems[1].View = chains.New(m)
		return hmTabs
	})
	m.Nav().Register(shared.AboutPagePattern, func(concretePath string) shared.View {
		hmTabs.tabItems[2].View = about.New(m)
		return hmTabs
	})
	m.Nav().NavigateToPath(shared.WalletPagePattern)
	return hmTabs
}

func (p *homeTabsLayout) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = gtx.Constraints.Max
			return p.Slider.Layout(gtx, p.tabItems[p.activeIndex].View.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			inset := layout.Inset{Top: 0, Bottom: 0}
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Spacing: layout.SpaceEvenly, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return p.tabItems[0].Layout(gtx, 0)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return p.tabItems[1].Layout(gtx, 1)
					}), layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return p.tabItems[2].Layout(gtx, 2)
					}),
				)
			})
		}),
	)
}
