package utils

// StackItem is intended for switching behavior, unlike push and pop up where
// stack size changes. children field is read only for this very same reason.
type StackItem[T any] struct {
	activeIndex         int
	children            []T
	indexChangeListener func()
}

func NewStackItem[T any](activeIndex int, children []T, indexChangeListener func()) *StackItem[T] {
	return &StackItem[T]{
		activeIndex: activeIndex,
		children:    children,
	}
}

func (s *StackItem[T]) SetActiveIndex(newIndex int) {
	prevIndex := s.GetActiveIndex()
	if newIndex < len(s.children) && newIndex != prevIndex {
		s.activeIndex = newIndex
		if s.indexChangeListener != nil {
			s.indexChangeListener()
		}
	}
}
func (s *StackItem[T]) GetActiveIndex() int {
	return s.activeIndex
}

func (s *StackItem[T]) GetActiveItem() T {
	var t T
	i := s.GetActiveIndex()
	childSize := len(s.children)
	if i >= 0 && i < childSize {
		return s.children[i]
	}
	return t
}

func (s *StackItem[T]) GetChildren() []T {
	return s.children
}

type Stack[T any] struct {
	history []*StackItem[T]
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		history: make([]*StackItem[T], 0),
	}
}

func (s *Stack[T]) CurrentItem() *StackItem[T] {
	if len(s.history) > 0 {
		return s.history[len(s.history)-1]
	}
	return nil
}
func (s *Stack[T]) Push(t *StackItem[T]) {
	s.history = append(s.history, t)
}
func (s *Stack[T]) PopUp(force bool) {
	if len(s.history) > 0 {
		if len(s.history) == 1 {
			if force {
				s.history = s.history[0 : len(s.history)-1]
			}
		} else {
			s.history = s.history[0 : len(s.history)-1]
		}
	}
}

func (s *Stack[T]) CanPopUp() bool {
	return s.Size() > 0
}

func (s *Stack[T]) Size() int {
	return len(s.history)
}

func (s *Stack[T]) Clone() []*StackItem[T] {
	history := make([]*StackItem[T], len(s.history))
	copy(history, s.history)
	return history
}
