package shared

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/x/component"
	"time"
)

type View interface {
	Layout(gtx layout.Context) layout.Dimensions
}

type Page struct {
	View
	URL string
}

type PagesStack struct {
	Page
	Stack []View
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

type Manager interface {
	CurrentPage() Page
	Push(v View)
	Snackbar() View
	Window() *app.Window
	GetWindowWidthInDp() unit.Dp
	GetWindowWidthInPx() int
	GetWindowHeightInDp() unit.Dp
	GetWindowHeightInPx() int
	PopUp()
	Modal() Modal
}

const (
	WalletPageURL = "/wallet"
	AboutPageURL  = "/about"
	ChainsPageURL = "/chains"
)
