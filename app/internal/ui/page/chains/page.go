package chains

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/app/internal/state/wallet"
	"github.com/partisiadev/partisiawallet/app/internal/ui/shared"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"
	"github.com/partisiadev/partisiawallet/app/internal/ui/view"
	"image"
	"strings"
	"sync"
)

type (
	Gtx  = layout.Context
	Dim  = layout.Dimensions
	Page = shared.Page
)

type page struct {
	shared.Manager
	AppTheme                   theme.AppTheme
	title                      string
	pageInitialized            bool
	width                      int
	searchText                 view.Search
	AllChainsList              widget.List
	chainItems                 []*tabAllChainsConnItem
	chainItemsFilteredByFilter []*tabAllChainsConnItem
	rpcClients                 []wallet.RpcClient
	filter                     string
	sync.Once
}

func New(m shared.Manager) shared.View {
	themeAlt := theme.GlobalTheme
	p := page{
		Manager:  m,
		AppTheme: themeAlt,
		title:    "Chains",
	}
	return &p
}

func (p *page) Layout(gtx Gtx) Dim {
	if !p.pageInitialized {
		if p.AppTheme == nil {
			p.AppTheme = theme.GlobalTheme
		}
		p.AllChainsList.Axis = layout.Vertical
		p.pageInitialized = true
	}
	p.width = gtx.Constraints.Max.X
	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd, Alignment: layout.Start}
	d := flex.Layout(gtx,

		layout.Rigid(p.chainsLayout),
	)
	return d
}

func (p *page) chainsLayout(gtx Gtx) Dim {
	p.Once.Do(p.preventChainItemsNilValue)
	p.generateFilteredChainItemsIfRequired(gtx)

	inset := layout.UniformInset(16)
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		flex := layout.Flex{Axis: layout.Vertical}
		return flex.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				var isLoading bool
				//for _, chainItem := range p.chainItemsFilteredByFilter {
				//	for _, item := range chainItem.connStateItems {
				//		if item.ClientState().GetState() == evm.RPCClientConnStateConnecting ||
				//			item.ClientState().GetState() == evm.RPCClientConnStateDisconnecting {
				//			isLoading = true
				//			break
				//		}
				//	}
				//	if isLoading {
				//		break
				//	}
				//}
				if isLoading {
					gtx.Constraints.Max.Y = gtx.Dp(56)
					gtx.Constraints.Min.Y = gtx.Dp(56)
					loader := view.Loader{AppTheme: p.AppTheme, Size: image.Pt(gtx.Dp(28), gtx.Dp(28))}
					return loader.Layout(gtx)
				}
				return layout.Dimensions{}
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				ml := material.List(p.AppTheme.Theme(), &p.AllChainsList)
				return ml.Layout(gtx, len(p.chainItemsFilteredByFilter), func(gtx layout.Context, index int) layout.Dimensions {
					return layout.Dimensions{}
					//if len(p.chainItemsFilteredByFilter[index].ConnChain.Clients) == 0 {
					//	return layout.Dimensions{}
					//}
					//flex := layout.Flex{Axis: layout.Vertical}
					//return flex.Layout(gtx,
					//	layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					//		inset := layout.Inset{Top: 12, Bottom: 12}
					//		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					//			return p.chainItemsFilteredByFilter[index].Layout(gtx, index)
					//		})
					//	}),
					//	layout.Rigid(component.Divider(p.AppTheme.Theme()).Layout),
					//)
				})
			}),
		)
	})
}

func (p *page) preventChainItemsNilValue() {
	if p.chainItems == nil {
		p.chainItems = make([]*tabAllChainsConnItem, 0)
	}
	if p.chainItemsFilteredByFilter == nil {
		p.chainItemsFilteredByFilter = p.chainItems[:0]
	}
	if p.rpcClients == nil {
		//rpcClients := evm.NewChainRPCsFromChains(evm.LoadChainsFromAssets())
		//chainItems := p.chainItems[:0]
		//for _, cl := range rpcClients {
		//	chainItem := tabAllChainsConnItem{
		//		//ConnChain: cl,
		//		AppTheme:  p.AppTheme,
		//		parent:    p,
		//	}
		//	chainItems = append(chainItems, &chainItem)
		//}
		//p.chainItems = chainItems
		//p.rpcClients = rpcClients
	}
}

func (p *page) generateFilteredChainItemsIfRequired(gtx Gtx) {
	prevFilter := p.filter
	currentFilter := strings.TrimSpace(strings.ToLower(p.searchText.Text()))
	isCurrentTextEmpty := strings.TrimSpace(currentFilter) == ""
	updateRequired := prevFilter != currentFilter ||
		len(p.chainItems) != len(p.rpcClients) ||
		(isCurrentTextEmpty && len(p.chainItemsFilteredByFilter) != len(p.chainItems))
	if updateRequired {
		p.filter = currentFilter
		if isCurrentTextEmpty {
			p.chainItemsFilteredByFilter = append(p.chainItemsFilteredByFilter[:0], p.chainItems...)
		} else {
			filteredItems := p.chainItemsFilteredByFilter[:0]
			//for _, ch := range p.chainItems {
			//	if len(ch.ConnChain.Clients) == 0 {
			//		continue
			//	}
			//	text1 := strings.ToLower(ch.ConnChain.Chain.Name)
			//	text2 := strings.ToLower(ch.ConnChain.Chain.ShortName)
			//	text3 := strings.ToLower(ch.ConnChain.Chain.Chain)
			//	shouldContain := strings.Contains(text1, currentFilter) || strings.Contains(text2, currentFilter) ||
			//		strings.Contains(text3, currentFilter)
			//	if shouldContain {
			//		filteredItems = append(filteredItems, ch)
			//	}
			//}
			//sort.SliceStable(filteredItems, func(i, j int) bool {
			//	prevName := strings.ToLower(filteredItems[i].ConnChain.Chain.Name)
			//	nextName := strings.ToLower(filteredItems[j].ConnChain.Chain.Name)
			//	prevNameIndex := strings.Index(prevName, currentFilter)
			//	nextNameIndex := strings.Index(nextName, currentFilter)
			//	if prevNameIndex < 0 {
			//		return false
			//	}
			//	if prevNameIndex >= 0 && nextNameIndex < 0 {
			//		return true
			//	}
			//	if prevNameIndex == nextNameIndex {
			//		return prevName[prevNameIndex:] < nextName[nextNameIndex:]
			//	}
			//	return prevNameIndex < nextNameIndex
			//})
			p.chainItemsFilteredByFilter = filteredItems
		}
		if len(p.chainItemsFilteredByFilter) > 0 {
			p.AllChainsList.ScrollTo(0)
		}
	}
}
