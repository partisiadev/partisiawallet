package fwk

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/x/component"
	"time"
)

// Page ---> Has a unique Path name among the siblings
// (ex like a folder/file which cannot be same inside the same parent folder)
type Page interface {
	PathName() string
	View
}

// View is anything that can be displayed on the screen
type View interface {
	Layout(gtx layout.Context) layout.Dimensions
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

const (
	WalletPageName        = `wallet`
	ChainsPageName        = `chains`
	AboutPageName         = `about`
	CreateAccountPageName = `createAccount`
	HomeTabsPageName      = `homeTabs`
)

type Manager interface {
	Snackbar() View
	Window() *app.Window
	WindowDimensions() WindowDimensions
	Modal() Modal
	Navigator() *Navigator
}
