package stack

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/internal/util"
)

// Sort performs an in-place sort of this collection.
//
// The item with the smallest value will be placed at the top of the stack.
func (s *Stack[T]) Sort() {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	// Stack is a reverse-ordered slice
	s.doSort(util.GosortDescending[T])
}

// Sorted returns a sorted copy of this stack as a new stack using the provided [functions.DeepCopyFunc] if any.
func (s *Stack[T]) Sorted() collections.Collection[T] {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	s1 := s.makeDeepCopy()

	if s1.size > 1 {
		s1.doSort(util.GosortDescending[T])
	}

	return s1
}

// SortDescending performs an in-place sort of this collection.
//
// The item with the largest value will be placed at the top of the stack.
func (s *Stack[T]) SortDescending() {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	s.doSort(util.Gosort[T])
}

// SortedDescending returns a descending order sorted copy of this stack as a new stack using the provided [functions.DeepCopyFunc] if any.
func (s *Stack[T]) SortedDescending() collections.Collection[T] {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	s1 := s.makeDeepCopy()

	if s1.size > 1 {
		s1.doSort(util.Gosort[T])
	}

	return s1
}

func (s *Stack[T]) doSort(f util.SortFunc[T]) {

	// bottom of stack (largest value ofter sorting) is at front of slice
	f(s.buffer, s.size, s.compare)
	s.version++
}
