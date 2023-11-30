package stack

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/util"
)

// Assert interface implementation.
var _ collections.Enumerable[int] = (*Stack[int])(nil)

// Any returns true for the first element found where the predicate function returns true.
// It returns false if no element matches the predicate.
func (s *Stack[T]) Any(predicate functions.PredicateFunc[T]) bool {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	iter := newForwardIterator[T](s, predicate)

	return iter.Start() != nil
}

// All applies the predicate function to every element in the collection,
// and returns true if all elements match the predicate.
func (s *Stack[T]) All(predicate functions.PredicateFunc[T]) bool {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	iter := newForwardIterator[T](s, util.DefaultPredicate[T])

	for e := iter.Start(); e != nil; e = iter.Next() {
		if !predicate(e.Value()) {
			return false
		}
	}

	return true
}

// ForEach applies function f to all elements in the collection.
func (s *Stack[T]) ForEach(f func(collections.Element[T])) {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	iter := newForwardIterator[T](s, util.DefaultPredicate[T])

	for e := iter.Start(); e != nil; e = iter.Next() {
		f(e)
	}
}

// Map applies function f to all elements in the collection
// and returns a new Stack containing the result of f.
func (s *Stack[T]) Map(f func(T) T) collections.Collection[T] {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	iter := newForwardIterator[T](s, util.DefaultPredicate[T])

	buf1 := New[T](WithCapacity[T](len(s.buffer)), WithComparer[T](s.compare))

	for e := iter.Start(); e != nil; e = iter.Next() {
		buf1.push(f(e.Value()))
	}

	return buf1
}

// Select returns a new Stack containing only the items for which predicate is true.
func (s *Stack[T]) Select(predicate functions.PredicateFunc[T]) collections.Collection[T] {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	return s.doSelect(predicate, false)
}

// SelectDeep returns a new Stack containing only the items for which predicate is true
//
// Elements are deep copied to the new collection using the provided [functions.DeepCopyFunc] if any.
func (s *Stack[T]) SelectDeep(predicate functions.PredicateFunc[T]) collections.Collection[T] {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	return s.doSelect(predicate, true)
}

// Find finds the first occurrence of an element matching the predicate.
//
// The function returns nil if no match.
func (s *Stack[T]) Find(predicate functions.PredicateFunc[T]) collections.Element[T] {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	result := s.doFind(predicate, false)

	if len(result) == 0 {
		return nil
	}

	return result[0]
}

// FindAll finds all occurrences of an element matching the predicate.
//
// The function returns an empty slice if none match.
func (s *Stack[T]) FindAll(predicate functions.PredicateFunc[T]) []collections.Element[T] {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	result := s.doFind(predicate, true)

	return result
}

func (s *Stack[T]) doFind(predicate functions.PredicateFunc[T], all bool) []collections.Element[T] {

	iter := newForwardIterator[T](s, predicate)
	result := make([]collections.Element[T], 0, util.DefaultCapacity)
	for e := iter.Start(); e != nil; e = iter.Next() {

		if predicate(e.Value()) {
			result = append(result, e)
		}

		if !all {
			break
		}
	}

	return result
}

func (s *Stack[T]) doSelect(predicate functions.PredicateFunc[T], deepCopy bool) collections.Collection[T] {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	s1 := New[T](WithCapacity[T](len(s.buffer)), WithComparer[T](s.compare))
	iter := newForwardIterator[T](s, predicate)

	for e := iter.Start(); e != nil; e = iter.Next() {
		if deepCopy {
			s1.push(s.copy(e.Value()))
		} else {
			s1.push(e.Value())
		}
	}

	return s1
}

// Min returns the minimum value in the collection according to the Comparer function.
func (s *Stack[T]) Min() T {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	return util.Min(s.buffer[0:s.size], s.compare, s.concurrent)
}

// Max returns the maximum value in the collection according to the Comparer function.
func (s *Stack[T]) Max() T {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	return util.Max(s.buffer[0:s.size], s.compare, s.concurrent)
}
