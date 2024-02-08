package ui

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/notify"
	"github.com/partisiadev/partisiawallet/log"
	"github.com/partisiadev/partisiawallet/ui/shared"
	"github.com/partisiadev/partisiawallet/ui/theme"
	"github.com/partisiadev/partisiawallet/ui/view"
	"image"
)

type manager struct {
	window         *app.Window
	constraints    layout.Constraints
	metric         unit.Metric
	notifier       notify.Notifier
	insets         system.Insets
	modalsStack    view.Modal
	snackbar       shared.View
	decoratedSize  layout.Dimensions
	isStageRunning bool
	//router         *router.Router
	nav shared.Nav
}

func newAppManager(window *app.Window) *manager {
	m := manager{}
	//m.router = router.New(nil)
	m.window = window
	var err error
	m.notifier, err = notify.NewNotifier()
	if err != nil {
		log.Logger().Errorln(err)
	}
	m.modalsStack = view.Modal{}
	homeTabsManager(&m)
	////m.snackbar = layoutView.NewSnackBar(theme.GlobalTheme)
	m.snackbar = &view.Modal{}
	//log.Logger().Println(m.Router().StackSize())
	return &m
}

func (m *manager) Layout(gtx layout.Context) layout.Dimensions {
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
				if m.Nav().CurrentPage() != nil &&
					m.Nav().CurrentPage().ActiveChild().Layout != nil {
					return m.Nav().CurrentPage().ActiveChild().Layout(gtx)
				}
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.H1(theme.GlobalTheme.Theme(), "Path not found").Layout(gtx)
				})
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

func (m *manager) Modal() shared.Modal {
	return &m.modalsStack
}

func (m *manager) Snackbar() shared.View {
	return m.snackbar
}
func (m *manager) Window() *app.Window {
	return m.window
}

func (m *manager) WindowDimensions() shared.WindowDimensions {
	return shared.WindowDimensions{
		WidthDp:  unit.Dp(float32(m.constraints.Max.X) / m.metric.PxPerDp),
		WidthPx:  m.constraints.Max.X,
		HeightDp: unit.Dp(float32(m.constraints.Max.Y) / m.metric.PxPerDp),
		HeightPx: m.constraints.Max.Y,
	}
}

func (m *manager) Nav() *shared.Nav {
	return &m.nav
}

//func (m *manager) Router() *router.Router {
//	return m.router
//}
