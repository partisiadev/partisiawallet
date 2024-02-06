package utils

import "sync"

type Stack[T any] struct {
	history      []T
	historyMutex sync.RWMutex
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		history: make([]T, 0),
	}
}

//func (s *Stack[T]) SetHistory(history []T) {
//	s.historyMutex.Lock()
//	defer s.historyMutex.Unlock()
//	s.history = append(s.history[:0], history...)
//}

//func (s *Stack[T]) GetHistory() []T {
//	s.historyMutex.RLock()
//	his := make([]T, len(s.history))
//	copy(his, s.history)
//	s.historyMutex.RUnlock()
//	return his
//}

func (s *Stack[T]) CurrentItem() T {
	s.historyMutex.RLock()
	defer s.historyMutex.RUnlock()
	if len(s.history) > 0 {
		return s.history[len(s.history)-1]
	}
	var t T
	return t
}
func (s *Stack[T]) Push(t T) {
	s.historyMutex.Lock()
	defer s.historyMutex.Unlock()
	s.history = append(s.history, t)
}
func (s *Stack[T]) PopUp() (didPopUP bool) {
	s.historyMutex.Lock()
	defer s.historyMutex.Unlock()
	if len(s.history) > 0 {
		s.history = s.history[0 : len(s.history)-1]
		didPopUP = true
	}
	return didPopUP
}

// Replace replaces the last item, if the stack is empty, it will add t
func (s *Stack[T]) Replace(t T) {
	s.historyMutex.Lock()
	defer s.historyMutex.Unlock()
	if len(s.history) == 0 {
		s.history = make([]T, 1)
	}
	s.history[len(s.history)-1] = t
}

func (s *Stack[T]) Size() int {
	s.historyMutex.RLock()
	defer s.historyMutex.RUnlock()
	return len(s.history)
}

func (s *Stack[T]) Clone() []T {
	s.historyMutex.Lock()
	defer s.historyMutex.Unlock()
	history := make([]T, len(s.history))
	copy(history, s.history)
	return history
}
