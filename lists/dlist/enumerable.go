package dlist

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
)

// Assert interface implementation
var _ collections.Enumerable[int] = (*DList[int])(nil)

// Any returns true for the first element found where the predicate function returns true.
// It returns false if no element matches the predicate.
func (l *DList[T]) Any(predicate functions.PredicateFunc[T]) bool {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	iter := newForwardIterator[T](l, predicate)

	return iter.Start() != nil
}

// All applies the predicate function to every element in the collection,
// and returns true if all elements match the predicate.
func (l *DList[T]) All(predicate functions.PredicateFunc[T]) bool {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	iter := newForwardIterator[T](l, util.DefaultPredicate[T])

	for e := iter.Start(); e != nil; e = iter.Next() {
		if !predicate(e.Value()) {
			return false
		}
	}

	return true
}

// ForEach applies function f to all elements in the collection
func (l *DList[T]) ForEach(f func(collections.Element[T])) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	iter := newForwardIterator[T](l, util.DefaultPredicate[T])

	for e := iter.Start(); e != nil; e = iter.Next() {
		f(e)
	}
}

// Map applies function f to all elements in the collection
// and returns a new DList containing the result of f
func (l *DList[T]) Map(f func(T) T) collections.Collection[T] {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	iter := newForwardIterator[T](l, util.DefaultPredicate[T])

	ll1 := New[T](WithComparer[T](l.compare))

	for e := iter.Start(); e != nil; e = iter.Next() {
		ll1.addItemLast(f(e.Value()))
	}

	return ll1
}

// Select returns a new DList containing only the items for which predicate is true
func (l *DList[T]) Select(predicate functions.PredicateFunc[T]) collections.Collection[T] {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	return l.doSelect(predicate, false)
}

// Select returns a new DList containing only the items for which predicate is true
//
// Elements are deep copied to the new collection using the provided [functions.DeepCopyFunc] if any.
func (l *DList[T]) SelectDeep(predicate functions.PredicateFunc[T]) collections.Collection[T] {

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	return l.doSelect(predicate, true)
}

// Find finds the first occurrence of an element matching the predicate.
//
// The function returns nil if no match.
func (l *DList[T]) Find(predicate functions.PredicateFunc[T]) collections.Element[T] {

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	result := l.doFind(predicate, forward, false)

	if len(result) == 0 {
		return nil
	}

	return result[0]
}

// FindLast searches for the given value in the list searchinbg backwards from the tail. Up to O(n).
//
// Returns the last node that contains the value; else nil.
func (l *DList[T]) FindLast(predicate functions.PredicateFunc[T]) collections.Element[T] {

	if l.head == nil {
		return nil
	}

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	result := l.doFind(predicate, reverse, false)

	if len(result) == 0 {
		return nil
	}

	return result[0]
}

// FindAll searches the entire collection for elements matching the predicate.
//
// The function returns an empty slice if none match.
func (l *DList[T]) FindAll(predicate functions.PredicateFunc[T]) []collections.Element[T] {

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	result := l.doFind(predicate, forward, true)

	return result
}

// Min returns the minimum value in the collection according to the Comparer function
func (l *DList[T]) Min() T {

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	if l.head == nil {
		panic(messages.COLLECTION_EMPTY)
	}

	m := l.head.item

	for current := l.head; current != nil; current = current.next {
		if l.compare(m, current.item) > 0 {
			m = current.item
		}
	}

	return m
}

// Max returns the maximum value in the collection according to the Comparer function
func (l *DList[T]) Max() T {

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	if l.head == nil {
		panic(messages.COLLECTION_EMPTY)
	}

	m := l.head.item

	for current := l.head; current != nil; current = current.next {
		if l.compare(m, current.item) < 0 {
			m = current.item
		}
	}

	return m
}

func (l *DList[T]) findNode(value T, direction direction) *DListNode[T] {

	node := util.Iif(direction == forward, l.head, l.tail)

	for node != nil {
		if l.compare(node.item, value) == 0 {
			return node
		}
		node = util.Iif(direction == forward, node.next, node.prev)
	}

	return nil
}

func (l *DList[T]) doFind(predicate functions.PredicateFunc[T], direction direction, all bool) []collections.Element[T] {

	result := make([]collections.Element[T], 0, util.DefaultCapacity)
	node := util.Iif(direction == forward, l.head, l.tail)

	for node != nil {
		if predicate(node.item) {
			result = append(result, util.NewElementType[T](l, &node.item))
		}

		if !all {
			break
		}

		node = util.Iif(direction == forward, node.next, node.prev)
	}

	return result
}

func (l *DList[T]) doSelect(predicate functions.PredicateFunc[T], deepCopy bool) collections.Collection[T] {
	ll1 := New[T](WithComparer[T](l.compare))
	iter := newForwardIterator[T](l, predicate)

	for e := iter.Start(); e != nil; e = iter.Next() {
		if deepCopy {
			ll1.addItemLast(l.copy(e.Value()))
		} else {
			ll1.addItemLast(e.Value())
		}
	}

	return ll1
}
