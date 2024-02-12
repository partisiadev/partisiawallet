package fwk

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"regexp"
	"strings"
)

// RouteHandler if view is nil then the path is not changed
// If path is relative and handler returns nil View, then
// the handler may be called again with absolute path
type RouteHandler interface {
	Handle(concretePath string) View
}

type PathView struct {
	view         View
	ConcretePath string
}

type Navigator struct {
	list               layout.List
	activeView         PathView
	activePattern      string
	registeredPatterns map[string]RouteHandler
}

func (n *Navigator) ActiveView() PathView {
	return n.activeView
}

func (n *Navigator) setActiveView(activeView PathView) {
	n.activeView = activeView
}

func NewNavigator() *Navigator {
	navigator := &Navigator{
		registeredPatterns: make(map[string]RouteHandler),
	}
	navigator.list.Axis = layout.Vertical
	return navigator
}

func (n *Navigator) Register(pattern string, handler RouteHandler) bool {
	_, ok := n.registeredPatterns[pattern]
	if ok {
		return false
	}
	n.registeredPatterns[pattern] = handler
	return true
}

func (n *Navigator) ActivePattern() string {
	return n.activePattern
}

func (n *Navigator) SetActivePattern(activePattern string) {
	n.activePattern = activePattern
}

func (n *Navigator) Layout(gtx layout.Context) layout.Dimensions {
	if n.activeView.view == nil {
		return n.fallbackLayout(gtx)
	}
	return n.activeView.view.Layout(gtx)
}

// NavigateTo Refer to RouteHandler
func (n *Navigator) NavigateTo(pth string) bool {
	if pth == n.ActiveView().ConcretePath {
		return false
	}
	if !strings.HasPrefix(pth, "/") {
		hndlr, ok := n.registeredPatterns[n.ActivePattern()]
		if ok {
			vw := hndlr.Handle(pth)
			if vw != nil {
				n.setActiveView(PathView{
					view:         vw,
					ConcretePath: n.ActiveView().ConcretePath + "/" + pth,
				})
			}
			return true
		}
		pth = n.ActiveView().ConcretePath + "/" + pth
	}
	for ptn, hndlr := range n.registeredPatterns {
		if regexp.MustCompile(ptn).MatchString(pth) {
			vw := hndlr.Handle(pth)
			if vw != nil {
				n.setActiveView(PathView{
					view:         vw,
					ConcretePath: pth,
				})
			}
		}
	}
	return false
}

func (n *Navigator) fallbackLayout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max
	return n.list.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return material.Body1(material.NewTheme(), "This indicates error").Layout(gtx)
		})
	})
}
