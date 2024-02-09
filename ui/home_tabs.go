package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
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
	tabItems [3]*tabItem
	*shared.NavStack
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

func homeTabsManager(m *manager) *homeTabsLayout {
	walletIcon, _ := widget.NewIcon(icons.ActionAccountBalanceWallet)
	chainsIcon, _ := widget.NewIcon(icons.ContentLink)
	aboutIcon, _ := widget.NewIcon(icons.ActionInfo)
	tabItems := [3]*tabItem{
		{URL: shared.WalletPagePattern, m: m, title: "Wallet", View: wallet.New(m), icon: walletIcon, index: 0},
		{URL: shared.ChainsPagePattern, m: m, title: "Chains", View: chains.New(m), icon: chainsIcon, index: 1},
		{URL: shared.AboutPagePattern, m: m, title: "About", View: about.New(m), icon: aboutIcon, index: 2},
	}
	hmTabs := &homeTabsLayout{Manager: m, tabItems: tabItems}
	hmTabs.tabItems[0].parent = hmTabs
	hmTabs.tabItems[1].parent = hmTabs
	hmTabs.tabItems[2].parent = hmTabs
	stackChildren := []*shared.NavStackChild{
		{URL: shared.WalletPagePattern, View: hmTabs},
		{URL: shared.ChainsPagePattern, View: hmTabs},
		{URL: shared.AboutPagePattern, View: hmTabs},
	}
	m.Nav().Register(shared.WalletPagePattern,
		func(concretePath string) func(stack *shared.NavStack) shared.View {
			return func(stack *shared.NavStack) shared.View {
				hmTabs.NavStack = stack
				vw := wallet.New(m)
				hmTabs.tabItems[0].View = vw
				stackChildren[0].URL = concretePath
				hmTabs.SetChildren(stackChildren)
				hmTabs.SetActiveIndex(hmTabs.activeIndex)
				return vw
			}
		})
	m.Nav().Register(shared.ChainsPagePattern, func(concretePath string) func(stack *shared.NavStack) shared.View {
		return func(stack *shared.NavStack) shared.View {
			hmTabs.NavStack = stack
			vw := chains.New(m)
			hmTabs.tabItems[1].View = vw
			stackChildren[1].URL = concretePath
			hmTabs.SetChildren(stackChildren)
			hmTabs.SetActiveIndex(hmTabs.activeIndex)
			return vw
		}
	})
	m.Nav().Register(shared.AboutPagePattern, func(concretePath string) func(stack *shared.NavStack) shared.View {
		return func(stack *shared.NavStack) shared.View {
			hmTabs.NavStack = stack
			vw := about.New(m)
			hmTabs.tabItems[2].View = vw
			stackChildren[2].URL = concretePath
			hmTabs.SetChildren(stackChildren)
			hmTabs.SetActiveIndex(hmTabs.activeIndex)
			return vw
		}
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
						return p.tabItems[0].Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return p.tabItems[1].Layout(gtx)
					}), layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return p.tabItems[2].Layout(gtx)
					}),
				)
			})
		}),
	)
}
