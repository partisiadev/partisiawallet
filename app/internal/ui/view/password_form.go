package view

import (
	"errors"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/partisiadev/partisiawallet/app/internal/state/wallet"
	"github.com/partisiadev/partisiawallet/app/internal/ui/theme"

	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"strings"
	"sync"
)

const setNewPassword = "Set new password"
const enterPassword = "Enter current password"
const reEnterPassword = "Re-enter current password"
const reEnterNewPassword = "Re-enter new password"

type PasswordForm struct {
	AppTheme               theme.AppTheme
	inActiveTheme          theme.AppTheme
	inputAuthErrStr        string
	buttonSubmit           IconButton
	inputPasswdState       inputPasswdFieldState
	inputRepeatPasswdState inputPasswdFieldState
	errorAuth              error
	authenticating         bool
	OnSuccess              func()
	initialized            bool
}

func NewPasswordForm(OnSuccess func()) *PasswordForm {
	iconSubmit, _ := widget.NewIcon(icons.ActionDone)
	inActiveTheme := theme.GlobalTheme.Clone()
	inActiveTheme.Theme().ContrastBg = color.NRGBA(colornames.Grey500)
	passForm := PasswordForm{
		AppTheme:      theme.GlobalTheme,
		OnSuccess:     OnSuccess,
		inActiveTheme: inActiveTheme,
		buttonSubmit: IconButton{
			AppTheme: theme.GlobalTheme,
			Icon:     iconSubmit,
			Text:     "Submit",
		},
	}
	return &passForm
}

func (p *PasswordForm) Layout(gtx Gtx) Dim {
	if !p.initialized {
		if p.AppTheme == nil {
			p.AppTheme = theme.GlobalTheme
		}
		for i, inputState := range []*inputPasswdFieldState{
			&p.inputPasswdState,
			&p.inputRepeatPasswdState,
		} {
			inputState.editor.SingleLine = true
			inputState.appTheme = p.AppTheme
			inputState.border = widget.Border{
				Color:        p.AppTheme.Theme().ContrastBg,
				CornerRadius: 8,
				Width:        1,
			}
			inputState.border.Color.A = 100
			inputState.labelStyle = material.Label(p.AppTheme.Theme(), 16, setNewPassword)
			inputState.hintText = setNewPassword
			if i == 1 {
				inputState.labelStyle.Text = reEnterNewPassword
				inputState.hintText = reEnterNewPassword
			}
		}
		p.initialized = true
	}
	p.handleFieldEvents(gtx)
	inset := layout.UniformInset(unit.Dp(16))
	flex := Flex{Axis: layout.Vertical, Alignment: layout.Start}
	d := flex.Layout(gtx,
		Rigid(func(gtx Gtx) Dim {
			return inset.Layout(gtx, p.drawPasswordTextFields)
		}),
	)
	if p.authenticating {
		layout.Stack{}.Layout(gtx,
			Stacked(func(gtx layout.Context) Dim {
				loader := Loader{}
				gtx.Constraints.Max, gtx.Constraints.Min = d.Size, d.Size
				return Flex{Alignment: layout.Middle, Spacing: layout.SpaceSides}.Layout(gtx,
					Rigid(func(gtx Gtx) Dim {
						return loader.Layout(gtx)
					}))
			}),
		)
		return d
	}
	return d
}

func (p *PasswordForm) handleFieldEvents(gtx Gtx) {
	if p.inputPasswdState.editor.Text() != p.inputPasswdState.inputStr ||
		p.inputRepeatPasswdState.editor.Text() != p.inputRepeatPasswdState.inputStr {
		p.errorAuth = nil
		p.inputAuthErrStr = ""
	}

	p.inputPasswdState.inputStr = p.inputPasswdState.editor.Text()
	p.inputRepeatPasswdState.inputStr = p.inputRepeatPasswdState.editor.Text()
	p.inputPasswdState.labelStyle.Text = setNewPassword
	p.inputRepeatPasswdState.labelStyle.Text = reEnterPassword
	dbExists := wallet.GlobalWallet.DatabaseExists()
	if dbExists {
		p.inputPasswdState.labelStyle.Text = enterPassword
		p.inputPasswdState.hintText = enterPassword
		p.inputRepeatPasswdState.labelStyle.Text = reEnterPassword
		p.inputRepeatPasswdState.hintText = reEnterPassword
	}

	for _, in := range []*inputPasswdFieldState{
		&p.inputPasswdState, &p.inputRepeatPasswdState} {
		in.labelStyle.Font.Weight = font.Normal
		in.border.Color.A = 100
		if in.btnClear.Clicked(gtx) {
			in.editor.SetText("")
			p.inputAuthErrStr = ""
			p.errorAuth = nil
			in.editor.Focus()
			in.labelStyle.Font.Weight = font.Bold
			in.border.Color.A = 255
		}
		if in.btnShowHide.Clicked(gtx) {
			in.labelStyle.Font.Weight = font.Bold
			in.border.Color.A = 255
			in.editor.Focus()
			if in.editor.Mask == '*' {
				in.editor.Mask = '\x00'
			} else {
				in.editor.Mask = '*'
			}
		}
		if in.editor.Focused() ||
			in.btnInput.Hovered() ||
			in.btnClear.Focused() ||
			in.btnClear.Hovered() ||
			in.btnShowHide.Focused() ||
			in.btnShowHide.Hovered() {
			in.labelStyle.Font.Weight = font.Bold
			in.border.Color.A = 255
		}
	}
	if p.inputPasswdState.btnInput.Clicked(gtx) {
		p.inputPasswdState.labelStyle.Font.Weight = font.Bold
		p.inputPasswdState.editor.Focus()
	}
	if p.inputRepeatPasswdState.btnInput.Clicked(gtx) {
		p.inputRepeatPasswdState.editor.Focus()
	}
}

func (p *PasswordForm) drawPasswordTextFields(gtx Gtx) Dim {
	if p.buttonSubmit.Button.Clicked(gtx) {
		p.authenticating = true
		if strings.TrimSpace(p.inputPasswdState.editor.Text()) != strings.TrimSpace(p.inputRepeatPasswdState.editor.Text()) {
			p.errorAuth = errors.New("Password mismatch!\n Please make sure password matches in both the inputs")
			p.authenticating = false
			p.inputAuthErrStr = p.errorAuth.Error()
		} else {
			p.errorAuth = wallet.GlobalWallet.OpenFromPassword(strings.TrimSpace(p.inputPasswdState.editor.Text()))
			p.authenticating = false
			if p.errorAuth != nil {
				p.inputAuthErrStr = p.errorAuth.Error()
			}
			if p.errorAuth == nil {
				p.inputAuthErrStr = ""
				if p.OnSuccess != nil {
					p.OnSuccess()
				}
			}
		}
	}
	gtx.Constraints.Min = gtx.Constraints.Max
	return Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceEnd}.Layout(gtx,
		Rigid(func(gtx Gtx) Dim {
			return DrawAppImageCenter(gtx, p.AppTheme)
		}),
		Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		Rigid(func(gtx layout.Context) Dim {
			return p.inputPasswdState.Layout(gtx)
		}),
		Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		Rigid(func(gtx layout.Context) Dim {
			return p.inputRepeatPasswdState.Layout(gtx)
		}),
		Rigid(func(gtx Gtx) Dim {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			inset := layout.Inset{Top: 32}
			return inset.Layout(gtx, func(gtx Gtx) Dim {
				return p.buttonSubmit.Layout(gtx)
			})
		}),
		Rigid(func(gtx layout.Context) Dim {
			if p.inputAuthErrStr != "" {
				lbl := material.Label(p.AppTheme.Theme(), 16, p.inputAuthErrStr)
				lbl.Color = color.NRGBA(colornames.Red500)
				inset := layout.UniformInset(8)
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return lbl.Layout(gtx)
				})
			}
			return Dim{}
		}),
	)
}

type inputPasswdFieldState struct {
	sync.Once
	editor      widget.Editor
	inputStr    string
	btnClear    widget.Clickable
	btnShowHide widget.Clickable
	btnInput    widget.Clickable
	border      widget.Border
	labelStyle  material.LabelStyle
	hintText    string
	appTheme    theme.AppTheme
}

func (i *inputPasswdFieldState) Layout(gtx Gtx) Dim {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceEnd}.Layout(gtx,
		Rigid(func(gtx layout.Context) Dim {
			return i.labelStyle.Layout(gtx)
		}),
		Rigid(layout.Spacer{Height: 2}.Layout),
		Rigid(func(gtx layout.Context) Dim {
			return i.btnInput.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return i.drawFormField(gtx)
			})
		}),
	)
}

func (i *inputPasswdFieldState) drawFormField(gtx Gtx) Dim {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return Flex{Alignment: layout.Middle}.Layout(gtx,
		Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return i.border.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					inset := layout.Inset{Top: 8, Bottom: 8, Left: 12, Right: 12}
					return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return Flex{Alignment: layout.Middle}.Layout(gtx,
							Flexed(1, func(gtx layout.Context) Dim {
								gtx.Constraints.Min.X = gtx.Constraints.Max.X
								return material.Editor(i.appTheme.Theme(),
									&i.editor,
									i.hintText,
								).Layout(gtx)
							}),
							Rigid(layout.Spacer{Width: 8}.Layout),
							Rigid(func(gtx layout.Context) Dim {
								icon, _ := widget.NewIcon(icons.ActionVisibility)
								if i.editor.Mask == '*' {
									icon, _ = widget.NewIcon(icons.ActionVisibilityOff)
								}
								btn := material.IconButton(i.appTheme.Theme(),
									&i.btnShowHide, icon, "Show/Hide Password")
								btn.Size = unit.Dp(25)
								btn.Inset = layout.Inset{}
								return btn.Layout(gtx)
							}),
							Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
							Rigid(func(gtx layout.Context) Dim {
								clearIcon, _ := widget.NewIcon(icons.ContentClear)
								btn := material.IconButton(i.appTheme.Theme(),
									&i.btnClear, clearIcon, "Clear Password")
								btn.Size = unit.Dp(25)
								btn.Inset = layout.Inset{}
								return btn.Layout(gtx)
							}),
						)
					})
				},
			)
		}),
	)
}
