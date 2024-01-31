package ui

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/x/component"
	"gioui.org/x/notify"
	"github.com/partisiadev/partisiawallet/app/internal/ui/page/wallet"
	"github.com/partisiadev/partisiawallet/app/internal/ui/shared"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"
	"github.com/partisiadev/partisiawallet/app/internal/ui/view"
	"github.com/partisiadev/partisiawallet/log"
	"image"
)

type Manager struct {
	window         *app.Window
	constraints    layout.Constraints
	metric         unit.Metric
	notifier       notify.Notifier
	insets         system.Insets
	modalsStack    view.Modal
	pagesStack     shared.PagesStack
	snackbar       shared.View
	decoratedSize  layout.Dimensions
	isStageRunning bool
}

func newAppManager(window *app.Window) *Manager {
	m := Manager{}
	m.window = window
	m.pagesStack = shared.PagesStack{
		Page: shared.Page{
			View: wallet.New(&m),
			URL:  shared.WalletPageURL,
		},
		Stack: nil,
	}
	var err error
	m.notifier, err = notify.NewNotifier()
	if err != nil {
		log.Logger().Errorln(err)
	}
	m.modalsStack = view.Modal{}
	//m.snackbar = view.NewSnackBar(theme.GlobalTheme)
	m.snackbar = &view.Modal{}
	return &m
}

func (m *Manager) Layout(gtx layout.Context) layout.Dimensions {
	stackLayout := layout.Stack{}
	maxDim := gtx.Constraints.Max
	topInsets := gtx.Dp(m.insets.Top)
	topDecoration := m.decoratedSize.Size.Y
	bottomInsets := gtx.Dp(m.insets.Bottom)
	d := stackLayout.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			d := layout.Flex{Axis: layout.Vertical}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.Rect{
					Color: theme.GlobalTheme.Theme().ContrastBg,
					Size:  image.Point{X: maxDim.X, Y: topInsets + topDecoration},
					Radii: 0,
				}.Layout(gtx)
			}), layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = gtx.Constraints.Max
				return m.CurrentPage().Layout(gtx)
			}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.Rect{
					Color: theme.GlobalTheme.Theme().ContrastBg,
					Size:  image.Point{X: maxDim.X, Y: bottomInsets},
					Radii: 0,
				}.Layout(gtx)
			}))
			return d
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = maxDim
			return m.snackbar.Layout(gtx)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = maxDim
			return m.modalsStack.Layout(gtx)
		}))

	return d
}

func (m *Manager) CurrentPage() shared.Page {
	return m.pagesStack.Page
}

func (m *Manager) Modal() shared.Modal {
	return &m.modalsStack
}

func (m *Manager) Push(pg shared.View) {
	m.pagesStack.Stack = append(m.pagesStack.Stack, pg)
}

func (m *Manager) Snackbar() view.View {
	return m.snackbar
}
func (m *Manager) Window() *app.Window {
	return m.window
}

func (m *Manager) GetWindowWidthInDp() unit.Dp {
	width := unit.Dp(float32(m.constraints.Max.X) / m.metric.PxPerDp)
	return width
}

func (m *Manager) GetWindowWidthInPx() int {
	return m.constraints.Max.X
}

func (m *Manager) GetWindowHeightInDp() unit.Dp {
	width := unit.Dp(float32(m.constraints.Max.Y) / m.metric.PxPerDp)
	return width
}

func (m *Manager) GetWindowHeightInPx() int {
	return m.constraints.Max.Y
}

func (m *Manager) PopUp() {
	if len(m.pagesStack.Stack) > 0 {
		m.pagesStack.Stack = m.pagesStack.Stack[:len(m.pagesStack.Stack)-1]
	}
}
