package view

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

type TabItem struct {
	*widget.Icon
	Title string
}

type Tabs struct {
	layout.List
	TabItems []TabItem
}
