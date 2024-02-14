package fwk

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/x/component"
	"time"
)

type Widget layout.Widget

func (w Widget) Layout(gtx layout.Context) layout.Dimensions {
	return w(gtx)
}

// View is anything that can be displayed on the screen
type View interface {
	Layout(gtx layout.Context) layout.Dimensions
}

type PathView interface {
	Name() string
	View
}

type SnackbarOption struct {
	View
	Duration time.Duration
}
type Snackbar interface {
	Show(option SnackbarOption)
}

type ModalOption struct {
	OnBackdropClick func()
	AfterDismiss    func()
	component.VisibilityAnimation
	Widget layout.Widget
}

type Modal interface {
	Show(option ModalOption)
	DismissWithAnim()
	DismissWithoutAnim()
}

type WindowDimensions struct {
	WidthDp  unit.Dp
	WidthPx  int
	HeightDp unit.Dp
	HeightPx int
}

type Manager interface {
	Snackbar() View
	Window() *app.Window
	WindowDimensions() WindowDimensions
	Modal() Modal
	Nav() *Nav
}
