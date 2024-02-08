package shared

import (
	"errors"
	"github.com/partisiadev/partisiawallet/utils"
	"regexp"
)

var (
	ErrPatternAlreadyRegistered = errors.New("pattern already registered")
)

type PageGenerator func(concretePath string) View

type Nav struct {
	history    []*NavStack
	generators utils.Map[string, PageGenerator]
}

// Register only registers if the pattern wasn't registered
func (n *Nav) Register(pattern string, generator PageGenerator) {
	n.generators.LoadOrStore(pattern, generator)
}

func (n *Nav) CurrentPage() *NavStack {
	if len(n.history) > 0 {
		return n.history[len(n.history)-1]
	}
	return nil
}

func (n *Nav) pushPage(page *NavStack) {
	n.history = append(n.history, page)
}

func (n *Nav) PopUp() {
	if len(n.history) > 1 {
		n.history = n.history[0 : len(n.history)-1]
	}
}

func (n *Nav) ReplacePage(p *NavStack) {
	if len(n.history) > 0 {
		n.history[len(n.history)-1] = p
	} else {
		n.history = append(n.history, p)
	}
}

func (n *Nav) NavigateToPath(p string) {
	if n.CurrentPage() != nil {
		if n.CurrentPage().Url() == p {
			return
		}
	}
	n.generators.Range(func(key string, value PageGenerator) bool {
		if ok := regexp.MustCompile(key).MatchString(p); ok {
			vw := value(p)
			page := NavStack{
				url:         p,
				pattern:     key,
				children:    make([]*NavStackChild, 0),
				activeIndex: 0,
				listener:    0,
			}
			page.PushView(p, vw)
			n.pushPage(&page)
			return false
		}
		return true
	})
}
