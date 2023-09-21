package slist

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
)

// Assert interface implementation
var _ collections.Enumerable[int] = (*SList[int])(nil)

// Any returns true for the first element found where the predicate function returns true.
// It returns false if no element matches the predicate.
func (l *SList[T]) Any(predicate functions.PredicateFunc[T]) bool {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	iter := newForwardIterator[T](l, predicate)

	return iter.Start() != nil
}

// All applies the predicate function to every element in the collection,
// and returns true if all elements match the predicate.
func (l *SList[T]) All(predicate functions.PredicateFunc[T]) bool {

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
func (l *SList[T]) ForEach(f func(collections.Element[T])) {

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
// and returns a new SList containing the result of f
func (l *SList[T]) Map(f func(T) T) collections.Collection[T] {

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

// Select returns a new SList containing only the items for which predicate is true
func (l *SList[T]) Select(predicate functions.PredicateFunc[T]) collections.Collection[T] {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	return l.doSelect(predicate, false)
}

// Select returns a new SList containing only the items for which predicate is true
//
// Elements are deep copied to the new collection using the provided [functions.DeepCopyFunc] if any.
func (l *SList[T]) SelectDeep(predicate functions.PredicateFunc[T]) collections.Collection[T] {

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	return l.doSelect(predicate, true)
}

// Find finds the first occurrence of an element matching the predicate.
//
// The function returns nil if no match.
func (l *SList[T]) Find(predicate functions.PredicateFunc[T]) collections.Element[T] {

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	result := l.doFind(predicate, false)

	if len(result) == 0 {
		return nil
	}

	return result[0]
}

// FindLast searches for the given value in the list searchinbg backwards from the tail. O(n).
//
// Returns the last node that contains the value; else nil.
func (l *SList[T]) FindLast(predicate functions.PredicateFunc[T]) collections.Element[T] {

	if l.head == nil {
		return nil
	}

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	result := l.doFind(predicate, true)

	length := len(result)

	if length == 0 {
		return nil
	}

	return result[length-1]
}

// FindAll searches the entire collection for elements matching the predicate.
//
// The function returns an empty slice if none match.
func (l *SList[T]) FindAll(predicate functions.PredicateFunc[T]) []collections.Element[T] {

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	result := l.doFind(predicate, true)

	return result
}

// Min returns the minimum value in the collection according to the Comparer function
func (l *SList[T]) Min() T {

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
func (l *SList[T]) Max() T {

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

func (l *SList[T]) findNode(value T) *SListNode[T] {

	for node := l.head; node != nil; node = node.next {
		if l.compare(node.item, value) == 0 {
			return node
		}
	}

	return nil
}

func (l *SList[T]) findNodeLast(value T) *SListNode[T] {

	var result *SListNode[T]

	for node := l.head; node != nil; node = node.next {
		if l.compare(node.item, value) == 0 {
			result = node
		}
	}

	return result
}

func (l *SList[T]) doFind(predicate functions.PredicateFunc[T], all bool) []collections.Element[T] {

	result := make([]collections.Element[T], 0, util.DefaultCapacity)

	for node := l.head; node != nil; node = node.next {
		if predicate(node.item) {
			result = append(result, util.NewElementType[T](l, &node.item))
		}

		if !all {
			break
		}

	}

	return result
}

func (l *SList[T]) doSelect(predicate functions.PredicateFunc[T], deepCopy bool) collections.Collection[T] {
	sl1 := New[T](WithComparer[T](l.compare))
	iter := newForwardIterator[T](l, predicate)

	for e := iter.Start(); e != nil; e = iter.Next() {
		if deepCopy {
			sl1.addItemLast(l.copy(e.Value()))
		} else {
			sl1.addItemLast(e.Value())
		}
	}

	return sl1
}
