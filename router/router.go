package router

import (
	"github.com/google/uuid"
	"github.com/partisiadev/partisiawallet/utils"
	"sync"
)

type Router struct {
	stack *utils.Stack[Path]
	// routes map[Path] config (Path can be wildcard and abstract)
	routes sync.Map
}

func New() *Router {
	return &Router{
		stack:  utils.NewStack[Path](),
		routes: sync.Map{},
	}
}

// Register returns updated config with Tag as unique string
func (r *Router) Register(config Config) (Config, error) {
	ok := r.IsRegistered(config.Path)
	if ok {
		return config, ErrPathAlreadyRegistered
	} else {
		config.Tag = uuid.New().String()
		r.routes.Store(config.Path, config)
	}
	return config, nil
}

func (r *Router) IsRegistered(path Path) (ok bool) {
	_, ok = r.routes.Load(path)
	return ok
}

func (r *Router) PushPath(path Path) error {
	err := path.Validate()
	if err != nil {
		return err
	}
	r.stack.Push(path)
	go func(path Path) {
		r.routes.Range(func(key, value any) bool {
			shouldContinue := true
			switch c := value.(type) {
			case Config:
				var isMatch bool
				isMatch, err = c.MatchesPath(path)
				if isMatch && err != nil {
					shouldContinue = false
					if c.OnActive != nil {
						c.OnActive(path)
					}
				}
			}
			return shouldContinue
		})
	}(path)
	return nil
}

func (r *Router) CurrentPath() Path {
	return r.stack.CurrentItem()
}

func (r *Router) PopUp() (didPopUP bool) {
	return r.stack.PopUp()
}

func (r *Router) StackSize() int {
	return r.stack.Size()
}

func (r *Router) SwitchPath(path Path) error {
	err := path.Validate()
	if err != nil {
		return err
	}
	r.stack.Replace(path)
	return nil
}
