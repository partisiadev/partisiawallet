package fwk

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"net/url"
	"regexp"
)

type Handle func(path *url.URL) View

func (ch Handle) Handle(path *url.URL) View {
	return ch(path)
}

type Router interface {
	Handle(path *url.URL) View
}

type route struct {
	*url.URL
	path          string
	activePattern string
	View
}

type Nav struct {
	route    route
	handlers map[string]Router
}

func (n *Nav) getRoute() route {
	return n.route
}

func (n *Nav) setRoute(route route) {
	n.route = route
}

func NewNav() *Nav {
	return &Nav{handlers: make(map[string]Router)}
}

func (n *Nav) Register(pattern string, handler Router) bool {
	_, ok := n.handlers[pattern]
	if ok {
		return false
	}
	n.handlers[pattern] = handler
	return true
}

func (n *Nav) NavigateTo(path string) bool {
	uRL, err := url.ParseRequestURI(path)
	if err != nil {
		return false
	}
	for k, h := range n.handlers {
		ok := regexp.MustCompile(k).MatchString(path)
		if ok {
			v := h.Handle(uRL)
			if v != nil {
				n.setRoute(route{
					URL:           uRL,
					activePattern: n.getRoute().activePattern,
					View:          v,
					path:          path,
				})
				return true
			}
		}
	}
	return false
}

func (n *Nav) Layout(gtx layout.Context) layout.Dimensions {
	if n.route.View != nil {
		return n.route.View.Layout(gtx)
	}
	return material.Body1(material.NewTheme(), "View Not Found").Layout(gtx)
}
