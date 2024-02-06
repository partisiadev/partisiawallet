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
	"github.com/partisiadev/partisiawallet/router"
	"github.com/partisiadev/partisiawallet/ui/page/wallet"
	"github.com/partisiadev/partisiawallet/ui/shared"
	"github.com/partisiadev/partisiawallet/ui/theme"
	"github.com/partisiadev/partisiawallet/ui/view"
	"image"
	"sync"
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
	router         *router.Router
	// registeredPaths map[Path]View
	registeredPaths sync.Map
	// view can be any view, it's main purpose
	// is to layout the current page
	view view.Slider
}

func newAppManager(window *app.Window) *manager {
	m := manager{}
	m.router = router.New()
	m.window = window
	var err error
	m.notifier, err = notify.NewNotifier()
	if err != nil {
		log.Logger().Errorln(err)
	}
	m.modalsStack = view.Modal{}
	walletPath := router.Path("/wallet")
	_, err = m.router.Register(router.Config{
		Path:    "/wallet",
		Pattern: router.DefaultPathParamPattern,
		OnActive: func(concretePath router.Path) {
			log.Logger().Println("on Active", concretePath)
		},
		Tag: "",
	})
	if err != nil {
		log.Logger().Fatal(err)
	}
	m.registeredPaths.Store(walletPath, wallet.New(&m))
	err = m.Router().SwitchPath(walletPath)
	if err != nil {
		log.Logger().Fatal(err)
	}
	//m.snackbar = view.NewSnackBar(theme.GlobalTheme)
	m.snackbar = &view.Modal{}
	log.Logger().Println(m.Router().StackSize())
	return &m
}

var swiper view.Swiper

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
				val, ok := m.registeredPaths.Load(m.Router().CurrentPath())
				if ok {
					switch vw := val.(type) {
					case shared.View:
						return swiper.Layout(gtx, 2, func(gtx layout.Context, index int) layout.Dimensions {
							return vw.Layout(gtx)
						})
					}
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

func (m *manager) Router() *router.Router {
	return m.router
}
