package orderedset

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
)

// Assert interface implementation.
var _ collections.Enumerable[int] = (*OrderedSet[int])(nil)

// Any returns true for the first element found where the predicate function returns true.
// It returns false if no element matches the predicate.
func (s *OrderedSet[T]) Any(predicate functions.PredicateFunc[T]) bool {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	iter := newForwardIterator[T](s, predicate)

	return iter.Start() != nil
}

// All applies the predicate function to every element in the collection,
// and returns true if all elements match the predicate.
func (s *OrderedSet[T]) All(predicate functions.PredicateFunc[T]) bool {

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
func (s *OrderedSet[T]) ForEach(f func(collections.Element[T])) {

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
// and returns a new OrderedSet containing the result of f.
func (s *OrderedSet[T]) Map(f func(T) T) collections.Collection[T] {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	iter := newForwardIterator[T](s, util.DefaultPredicate[T])

	s1 := New[T](WithComparer[T](s.compare))

	for e := iter.Start(); e != nil; e = iter.Next() {
		s1.doInsert(f(e.Value()))
	}

	return s1
}

// Select returns a new OrderedSet containing only the items for which predicate is true.
func (s *OrderedSet[T]) Select(predicate functions.PredicateFunc[T]) collections.Collection[T] {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	return s.doSelect(predicate, false)
}

// SelectDeep returns a new OrderedSet containing only the items for which predicate is true
//
// Elements are deep copied to the new collection using the provided [functions.DeepCopyFunc] if any.
func (s *OrderedSet[T]) SelectDeep(predicate functions.PredicateFunc[T]) collections.Collection[T] {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	return s.doSelect(predicate, true)
}

// Find finds the first occurrence of an element matching the predicate.
//
// The function returns nil if no match.
func (s *OrderedSet[T]) Find(predicate functions.PredicateFunc[T]) collections.Element[T] {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	result := s.find(predicate, false)

	if len(result) == 0 {
		return nil
	}

	return result[0]
}

// FindAll finds all occurrences of an element matching the predicate.
//
// The function returns an empty slice if none match.
func (s *OrderedSet[T]) FindAll(predicate functions.PredicateFunc[T]) []collections.Element[T] {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	result := s.find(predicate, true)

	return result
}

// Max returns the maximum value in the collection according to the Comparer function.
func (s *OrderedSet[T]) Max() T {

	if s.root == nil {
		panic(messages.COLLECTION_EMPTY)
	}

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	current := s.root
	for current.right != nil {
		current = current.right
	}

	return current.item
}

// Min returns the minimum value in the collection according to the Comparer function.
func (s *OrderedSet[T]) Min() T {

	if s.root == nil {
		panic(messages.COLLECTION_EMPTY)
	}

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	current := s.root
	for current.left != nil {
		current = current.left
	}

	return current.item
}

func (s *OrderedSet[T]) find(predicate functions.PredicateFunc[T], all bool) []collections.Element[T] {

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

func (s *OrderedSet[T]) doSelect(predicate functions.PredicateFunc[T], deepCopy bool) collections.Collection[T] {
	s1 := New[T](WithComparer[T](s.compare))
	iter := newForwardIterator[T](s, predicate)

	for e := iter.Start(); e != nil; e = iter.Next() {
		if deepCopy {
			s1.doInsert(s.copy(e.Value()))
		} else {
			s1.doInsert(e.Value())
		}
	}

	return s1
}
