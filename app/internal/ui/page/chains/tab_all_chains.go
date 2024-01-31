package chains

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/app/assets"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"
	_ "image/jpeg"
	_ "image/png"
)

type tabAllChainsConnItem struct {
	//ConnChain      evm.RpcClients
	connStateItems []*tabAllChainsConnStateItem
	layout.List
	theme.AppTheme
	initialized bool
	parent      *page
}

func (c *tabAllChainsConnItem) Layout(gtx Gtx, index int) Dim {
	if !c.initialized {
		if c.AppTheme == nil {
			c.AppTheme = theme.GlobalTheme
		}
		//if len(c.ConnChain.Clients) > 0 &&
		//	len(c.ConnChain.Clients) != len(c.connStateItems) {
		//	c.connStateItems = make([]*tabAllChainsConnStateItem, len(c.ConnChain.Clients))
		//	for i, connState := range c.ConnChain.Clients {
		//		c.connStateItems[i] = &tabAllChainsConnStateItem{RpcClient: connState, AppTheme: c.AppTheme}
		//	}
		//}
		c.initialized = true
	}
	flex := layout.Flex{Axis: layout.Vertical}
	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			flexAlt := layout.Flex{Alignment: layout.Middle}
			return flexAlt.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Dp(56)
					gtx.Constraints.Max.Y = gtx.Dp(56)
					imgOps := paint.NewImageOp(assets.AppIconImage)
					//first := c.parent.AllChainsList.List.Position.First
					//last := first + c.parent.AllChainsList.List.Position.Count
					//isViewVisible := index >= first && index <= last
					//if isViewVisible {
					//	img, ok := evm2.IconsEncSmallImageCache.Get(c.ConnChain.Chain.Icon)
					//	if img != nil && ok {
					//		imgOps = paint.NewImageOp(img)
					//	}
					//}
					imgWidget := widget.Image{Src: imgOps, Fit: widget.Contain, Position: layout.Center}
					return imgWidget.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: 8}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					//labelStyle := material.H5(c.AppTheme.Theme(), c.ConnChain.Chain.Name)
					//return labelStyle.Layout(gtx)
					return layout.Dimensions{}
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: 8}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			//txt := fmt.Sprintf("Currency: %s", c.ConnChain.Chain.NativeCurrency.Name)
			//w := material.Body1(c.AppTheme.Theme(), txt)
			//w.TextSize = unit.Sp(16)
			//return w.Layout(gtx)
			return layout.Dimensions{}
		}),
		layout.Rigid(layout.Spacer{Height: 8}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			c.List.Axis = layout.Vertical
			return c.List.Layout(gtx, len(c.connStateItems), func(gtx layout.Context, index int) layout.Dimensions {
				inset := layout.Inset{Bottom: 8}
				//if index == len(c.ConnChain.Clients)-1 {
				//	inset.Bottom = 0
				//}
				return inset.Layout(gtx, c.connStateItems[index].Layout)
			})
		}),
	)
}

type tabAllChainsConnStateItem struct {
	//evm.RpcClient
	btnConnect widget.Clickable
	theme.AppTheme
	layout.Inset
	initialized bool
}

func (c *tabAllChainsConnStateItem) Layout(gtx Gtx) Dim {
	if !c.initialized {
		if c.AppTheme.Theme() == nil {
			c.AppTheme = theme.GlobalTheme
		}
		c.initialized = true
	}
	//state := *c.ClientState()
	//isConnected := state.GetState() == evm.RPCClientConnStateConnected
	flex := layout.Flex{Axis: layout.Vertical}
	return flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			//w := material.Body1(c.AppTheme.Theme(), string(state.GetRPC()))
			//w.TextSize = unit.Sp(16)
			//return w.Layout(gtx)
			return layout.Dimensions{}
		}),
		layout.Rigid(layout.Spacer{Height: 2}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			flex := layout.Flex{}
			return flex.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Body1(c.AppTheme.Theme(), "Balance").Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: 8}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					//bal := c.bal
					//if !isConnected {
					//	bal = "Not Connected"
					//}
					//if c.err != nil {
					//	bal = c.err.Error()
					//}
					//if !evm.GlobalWallet.IsOpen() {
					//	bal = "No Active Account."
					//}
					//fetched := c.balFetched
					//if !fetched {
					//	go func() {
					//		acc, _ := evm.GlobalWallet.Account()
					//		c.bal, c.err = c.ShowBalance(acc)
					//		c.balFetched = true
					//	}()
					//}
					//return material.Body1(c.AppTheme, bal).Layout(gtx)
					return layout.Dimensions{}
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			flex := layout.Flex{Axis: layout.Vertical}
			return flex.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				flex := layout.Flex{Alignment: layout.Middle}
				return flex.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// txt := "Connect"
						//btnStyle := material.Button(c.AppTheme.Theme(), &c.btnConnect, txt)
						//if state.GetState() == evm.RPCClientConnStateConnecting ||
						//	state.GetState() == evm.RPCClientConnStateDisconnecting {
						//	rec := op.Record(gtx.Ops)
						//	dim := btnStyle.Layout(gtx)
						//	rec.Stop()
						//	gtx.Constraints.Max, gtx.Constraints.Min = dim.Size, dim.Size
						//	loader := view.Loader{AppTheme: c.AppTheme, Size: image.Pt(gtx.Dp(28), gtx.Dp(28))}
						//	return layout.Center.Layout(gtx, loader.Layout)
						//}
						//if c.btnConnect.Clicked(gtx) && (state.GetState() == evm.RPCClientConnStateIdle ||
						//	state.GetState() == evm.RPCClientConnStateConnected) {
						//	//c.balFetched = false
						//	if state.GetState() == evm.RPCClientConnStateIdle {
						//		go c.Connect()
						//	} else {
						//		go c.Disconnect()
						//	}
						//}
						//btnStyle.Background = color.NRGBA(colornames.Green)
						//if isConnected {
						//	txt = "Disconnect"
						//	btnStyle.Text = txt
						//	btnStyle.Background = color.NRGBA(colornames.Red)
						//}
						//return btnStyle.Layout(gtx)
						return layout.Dimensions{}
					}),
				)
			}))
		}),
	)
}
