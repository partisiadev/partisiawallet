package wallet

import (
	"fmt"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/partisiadev/partisiawallet/app/internal/state/wallet"
	"github.com/partisiadev/partisiawallet/app/internal/ui/shared"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"
	"github.com/partisiadev/partisiawallet/app/internal/ui/view"
	"github.com/partisiadev/partisiawallet/log"

	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"time"
)

type (
	Gtx         = layout.Context
	Dim         = layout.Dimensions
	Animation   = component.VisibilityAnimation
	View        = shared.View
	ModalOption = shared.ModalOption
)

type page struct {
	shared.Manager
	layout.List
	AppTheme           theme.AppTheme
	title              string
	iconNewChat        *widget.Icon
	btnBackdrop        widget.Clickable
	btnMenuIcon        widget.Clickable
	btnMenuContent     widget.Clickable
	btnAddAccount      widget.Clickable
	btnDeleteAccounts  widget.Clickable
	btnCloseSelection  widget.Clickable
	btnYes             widget.Clickable
	btnNo              widget.Clickable
	btnSelectAll       widget.Clickable
	btnDeleteAll       widget.Clickable
	btnSelectionMode   widget.Clickable
	menuIcon           *widget.Icon
	closeIcon          *widget.Icon
	menuVisibilityAnim component.VisibilityAnimation
	PasswordForm       View
	navigationIcon     *widget.Icon
	accountsView       []*pageItem
	NoAccount          View
	AccountForm        View
	ModalContent       *view.ModalContent
	SelectionMode      bool
	initialized        bool
}

func New(m shared.Manager) shared.View {
	closeIcon, _ := widget.NewIcon(icons.ContentClear)
	iconNewChat, _ := widget.NewIcon(icons.ContentCreate)
	iconMenu, _ := widget.NewIcon(icons.NavigationMoreVert)
	errorTh := *material.NewTheme()
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	themeAlt := theme.GlobalTheme.Clone()
	p := page{
		Manager:      m,
		AppTheme:     themeAlt,
		title:        "Wallet",
		iconNewChat:  iconNewChat,
		List:         layout.List{Axis: layout.Vertical},
		accountsView: []*pageItem{},
		menuIcon:     iconMenu,
		closeIcon:    closeIcon,
		menuVisibilityAnim: component.VisibilityAnimation{
			Duration: time.Millisecond * 250,
			State:    component.Invisible,
			Started:  time.Time{},
		},
	}
	p.AccountForm = view.NewAccountFormView(p.Manager, p.onSuccess)
	p.ModalContent = view.NewModalContent(func() {
		p.Modal().DismissWithAnim()
		p.AccountForm = view.NewAccountFormView(p.Manager, p.onSuccess)
	})
	p.NoAccount = view.NewNoAccount(p.Manager)
	p.PasswordForm = view.NewPasswordForm(func() {})
	return &p
}

func (p *page) Layout(gtx Gtx) Dim {
	if !p.initialized {
		if p.AppTheme == nil {
			p.AppTheme = theme.GlobalTheme
		}
		p.loadAccountsView()
		p.initialized = true
	}
	p.handleSelectionMode()
	p.handleAddAccountClick(gtx)
	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd, Alignment: layout.Start}

	d := flex.Layout(gtx,
		layout.Rigid(p.drawIdentitiesItems),
	)
	p.drawMenuLayout(gtx)
	p.handleEvents(gtx)
	return d
}

func (p *page) drawIdentitiesItems(gtx Gtx) Dim {
	isPasswordSet := wallet.GlobalWallet.IsOpen()
	if !isPasswordSet {
		return p.PasswordForm.Layout(gtx)
	}
	if len(p.accountsView) == 0 {
		return p.NoAccount.Layout(gtx)
	}
	return p.List.Layout(gtx, len(p.accountsView), func(gtx Gtx, index int) (d Dim) {
		inset := layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0)}
		return inset.Layout(gtx, func(gtx Gtx) Dim {
			return p.accountsView[index].Layout(gtx)
		})
	})
}

func (p *page) drawMenuLayout(gtx Gtx) Dim {
	if p.btnBackdrop.Clicked(gtx) {
		if !p.btnMenuContent.Pressed() {
			p.menuVisibilityAnim.Disappear(gtx.Now)
		}
		for _, idView := range p.accountsView {
			if !idView.btnMenuContent.Pressed() && !idView.Hovered() {
				idView.menuVisibilityAnim.Disappear(gtx.Now)
			}
		}
	}
	layout.Stack{Alignment: layout.NE}.Layout(gtx,
		layout.Stacked(func(gtx Gtx) Dim {
			return p.btnBackdrop.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					progress := p.menuVisibilityAnim.Revealed(gtx)
					gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * progress)
					gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * progress)
					return component.Rect{Size: gtx.Constraints.Max, Color: color.NRGBA{A: 200}}.Layout(gtx)
				},
			)
		}),
		layout.Stacked(func(gtx Gtx) Dim {
			progress := p.menuVisibilityAnim.Revealed(gtx)
			macro := op.Record(gtx.Ops)
			d := p.btnMenuContent.Layout(gtx, p.drawMenuItems)
			call := macro.Stop()
			d.Size.X = int(float32(d.Size.X) * progress)
			d.Size.Y = int(float32(d.Size.Y) * progress)
			component.Rect{Size: d.Size, Color: color.NRGBA(colornames.White)}.Layout(gtx)
			clipOp := clip.Rect{Max: d.Size}.Push(gtx.Ops)
			call.Add(gtx.Ops)
			clipOp.Pop()
			return d
		}),
	)
	return Dim{}
}

func (p *page) drawMenuItems(gtx Gtx) Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) / 1.5)
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	if p.SelectionMode {
		return p.drawSelectionMenuItems(gtx)
	}
	return p.drawNormalMenuItems(gtx)
}

func (p *page) drawNormalMenuItems(gtx Gtx) Dim {
	if p.btnSelectAll.Clicked(gtx) {
		p.selectAll()
		p.menuVisibilityAnim.Disappear(gtx.Now)
	}
	if p.btnSelectionMode.Clicked(gtx) {
		p.SelectionMode = true
		p.menuVisibilityAnim.Disappear(gtx.Now)
	}
	if p.btnDeleteAll.Clicked(gtx) {
		p.selectAll()
		p.Manager.Modal().Show(shared.ModalOption{Widget: p.drawDeleteAccountsModal,
			VisibilityAnimation: component.VisibilityAnimation{
				Duration: time.Millisecond * 250,
				State:    component.Invisible,
				Started:  time.Time{},
			}},
		)
		p.menuVisibilityAnim.Disappear(gtx.Now)
	}

	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
		p.drawMenuItem("Add Account", &p.btnAddAccount),
		p.drawMenuItem("Selection Mode", &p.btnSelectionMode),
		p.drawMenuItem("Select All Accounts", &p.btnSelectAll),
		p.drawMenuItem("Delete All Accounts", &p.btnDeleteAll),
	)
}
func (p *page) drawSelectionMenuItems(gtx Gtx) Dim {
	if p.btnDeleteAccounts.Clicked(gtx) {
		p.Manager.Modal().Show(shared.ModalOption{
			Widget: p.drawDeleteAccountsModal,
			VisibilityAnimation: Animation{
				Duration: time.Millisecond * 250,
				State:    component.Invisible,
				Started:  time.Time{},
			}},
		)
		p.menuVisibilityAnim.Disappear(gtx.Now)
	}

	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
		p.drawMenuItem("Delete Selected Accounts", &p.btnDeleteAccounts),
		p.drawMenuItem("Clear Selection", &p.btnCloseSelection),
	)
}
func (p *page) drawMenuItem(txt string, btn *widget.Clickable) layout.FlexChild {
	inset := layout.UniformInset(unit.Dp(12))
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		btnStyle := material.ButtonLayoutStyle{Button: btn}
		btnStyle.Background = color.NRGBA(colornames.White)
		return btnStyle.Layout(gtx,
			func(gtx Gtx) Dim {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				inset := inset
				return inset.Layout(gtx, func(gtx Gtx) Dim {
					return layout.Flex{Spacing: layout.SpaceEnd}.Layout(gtx,
						layout.Rigid(func(gtx Gtx) Dim {
							bd := material.Body1(p.AppTheme.Theme(), txt)
							bd.Color = color.NRGBA(colornames.Black)
							bd.Alignment = text.Start
							return bd.Layout(gtx)
						}),
					)
				})
			},
		)
	})
}

func (p *page) drawAddAccountModal(gtx Gtx) Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	return p.ModalContent.DrawContent(gtx, p.AppTheme, p.AccountForm.Layout)
}

func (p *page) drawDeleteAccountsModal(gtx Gtx) Dim {
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.85)
	gtx.Constraints.Max.Y = int(float32(gtx.Constraints.Max.Y) * 0.85)
	if p.btnYes.Clicked(gtx) {
		accounts := make([]wallet.Account, 0)
		accountsViewSize := len(p.accountsView)
		for _, eachView := range p.accountsView {
			if eachView.Selected {
				accounts = append(accounts, eachView.Account)
			}
		}
		err := wallet.GlobalWallet.DeleteAccounts(accounts)
		if err != nil {
			log.Logger().Errorln(err)
		}
		p.Manager.Modal().DismissWithAnim()
		p.clearAllSelection()
		//	var txtTmp string
		if len(accounts) == accountsViewSize {
			//	txtTmp = "all wallet."
		} else {
			//txtTmp = fmt.Sprintf("%d wallet.", len(accounts))
		}
		if len(accounts) == 1 {
			//txtTmp = "1 account."
		}
		//		txt := fmt.Sprintf("Successfully deleted %s", txtTmp)
		//p.Manager.Snackbar().Show(txt, nil, color.NRGBA{}, "")
		p.loadAccountsView()
	}
	if p.btnNo.Clicked(gtx) {
		p.Manager.Modal().DismissWithAnim()
		p.clearAllSelection()
	}
	count := p.getSelectionCount()
	accountsSize := len(p.accountsView)
	var txt string
	if count == accountsSize {
		txt = "all wallet"
	} else {
		txt = fmt.Sprintf("%d selected wallet", count)
	}
	if count == 1 {
		txt = "the selected account"
	}
	promptContent := view.NewPromptContent(p.AppTheme,
		"Account Deletion!",
		fmt.Sprintf("Are you sure you want to delete %s?", txt),
		&p.btnYes, &p.btnNo)
	return p.ModalContent.DrawContent(gtx, p.AppTheme, promptContent.Layout)
}

func (p *page) onSuccess() {
	p.Manager.Modal().DismissWithAnim()
	p.AccountForm = view.NewAccountFormView(p.Manager, p.onSuccess)
	//a, _ := wallet.GlobalWallet.Account()
	//txt := fmt.Sprintf("Successfully created %s", a.PublicKey)
	p.Manager.Window().Invalidate()
	//p.Manager.Snackbar().Show(txt, nil, color.NRGBA{}, "")
}
func (p *page) getSelectionCount() (count int) {
	for _, item := range p.accountsView {
		if item.Selected {
			count++
		}
	}
	return count
}
func (p *page) clearAllSelection() {
	p.SelectionMode = false
	for _, item := range p.accountsView {
		item.Selected = false
		item.SelectionMode = false
	}
}
func (p *page) selectAll() {
	p.SelectionMode = true
	for _, item := range p.accountsView {
		item.Selected = true
		item.SelectionMode = true
	}
}

func (p *page) handleSelectionMode() {
	for _, item := range p.accountsView {
		if p.SelectionMode {
			item.SelectionMode = p.SelectionMode
		} else if item.SelectionMode {
			p.SelectionMode = item.SelectionMode
			break
		}
	}
	if p.SelectionMode {
		p.AppTheme.Theme().ContrastBg = color.NRGBA{A: 255}
	} else {
		p.AppTheme.Theme().ContrastBg = theme.GlobalTheme.Theme().ContrastBg
	}
}

func (p *page) handleAddAccountClick(gtx Gtx) {
	if p.btnAddAccount.Clicked(gtx) {
		p.Manager.Modal().Show(shared.ModalOption{
			Widget: p.drawAddAccountModal,
			AfterDismiss: func() {
				p.AccountForm = view.NewAccountFormView(p.Manager, p.onSuccess)
				p.Manager.Modal().DismissWithAnim()
			},
			VisibilityAnimation: Animation{
				Duration: time.Millisecond * 250,
				State:    component.Invisible,
				Started:  time.Time{},
			}},
		)
		p.menuVisibilityAnim.Disappear(gtx.Now)
	}
}

func (p *page) handleEvents(gtx Gtx) {
	for _, e := range gtx.Queue.Events(p) {
		if e, ok := e.(pointer.Event); ok {
			if e.Kind == pointer.Press {
				if !p.btnMenuContent.Pressed() {
					p.menuVisibilityAnim.Disappear(gtx.Now)
				}
				for _, idView := range p.accountsView {
					if !idView.btnMenuContent.Pressed() && !idView.Hovered() {
						idView.menuVisibilityAnim.Disappear(gtx.Now)
					}
				}
			}
		}
	}
}

func (p *page) loadAccountsView() {
	accs, _ := wallet.GlobalWallet.Accounts()
	accsView := make([]*pageItem, len(accs))
	for i := range accsView {
		accsView[i] = &pageItem{
			Manager:      p.Manager,
			AppTheme:     p.AppTheme,
			Account:      accs[i],
			ModalContent: p.ModalContent,
		}
	}
	p.accountsView = accsView
}
