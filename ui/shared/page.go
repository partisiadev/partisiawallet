package shared

type NavStackChild struct {
	URL string
	View
}

type NavStack struct {
	parent      *NavStack
	url         string
	pattern     string
	children    []*NavStackChild
	activeIndex int
	listener    int
}

func (p *NavStack) Parent() *NavStack {
	return p.parent
}

func (p *NavStack) Url() string {
	return p.url
}

func (p *NavStack) Pattern() string {
	return p.pattern
}

func (p *NavStack) Children() []*NavStackChild {
	return p.children
}

func (p *NavStack) ActiveIndex() int {
	return p.activeIndex
}

func (p *NavStack) SetActiveIndex(activeIndex int) {
	if activeIndex < len(p.Children()) {
		p.activeIndex = activeIndex
	}
}

func (p *NavStack) ActiveChild() *NavStackChild {
	if p.ActiveIndex() < len(p.Children()) {
		return p.Children()[p.ActiveIndex()]
	}
	return nil
}

func (p *NavStack) PushView(url string, view View) {
	p.children = append(p.children, &NavStackChild{
		URL:  url,
		View: view,
	})
}
func (p *NavStack) SetChildren(children []*NavStackChild) {
	p.children = append(p.children, children...)
}
