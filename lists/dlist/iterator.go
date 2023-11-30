package dlist

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/local"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
)

type direction bool

const (
	forward, reverse direction = true, false
)

// Assert interface implementation
var _ collections.Iterable[int] = (*DList[int])(nil)
var _ collections.ReverseIterable[int] = (*DList[int])(nil)

// DListIterator implements an iterator over the elements in the list.
type DListIterator[T any] struct {
	util.IteratorBase[T]
	list      *DList[T]
	current   *DListNode[T]
	startNode *DListNode[T]
	predicate functions.PredicateFunc[T]
	direction direction

	local.InternalImpl
}

func newForwardIterator[T any](list *DList[T], predicate functions.PredicateFunc[T]) collections.Iterator[T] {
	return &DListIterator[T]{
		list:      list,
		current:   list.First(),
		startNode: list.First(),
		direction: forward,
		predicate: predicate,
		IteratorBase: util.IteratorBase[T]{
			Version:    list.version,
			NilElement: nil,
		},
	}
}

func newReverseIterator[T any](list *DList[T]) collections.Iterator[T] {
	return &DListIterator[T]{
		list:      list,
		current:   list.Last(),
		startNode: list.Last(),
		direction: reverse,
		predicate: func(item T) bool { return true },
		IteratorBase: util.IteratorBase[T]{
			Version:    list.version,
			NilElement: nil,
		},
	}
}

// Iterator returns an iterator that walks the DList from first to last element
//
//	iter := ll.Iterator()
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (l *DList[T]) Iterator() collections.Iterator[T] {

	return newForwardIterator(l, util.DefaultPredicate[T])
}

// ReverseIterator returns an iterator that walks the DList from last to first element
//
//	iter := ll.ReverseIterator()
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (l *DList[T]) ReverseIterator() collections.Iterator[T] {

	return newReverseIterator(l)
}

// TakeWhile returns a forward iterater that walks the collection returning only
// those elements for which predicate returns true.
//
//	ll := linkedlist.New[int]()
//	// add values
//	iter := ll.TakeWhile(func (val int) bool { return val % 2 == 0 })
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (l *DList[T]) TakeWhile(predicate functions.PredicateFunc[T]) collections.Iterator[T] {

	return newForwardIterator(l, predicate)
}

// Start begins an iteration across the DList returning the fisrt element,
// which will be nil if the collection is empty.
//
// Panics if the collection has been modified since creation of the iterator.
func (i *DListIterator[T]) Start() collections.Element[T] {
	i.validateIterator()
	i.current = i.startNode
	if i.current == nil {
		return i.NilElement
	}

	if !i.predicate(i.current.item) {
		return i.Next()
	}

	return util.NewElementType[T](i.list, &i.current.item)
}

// Next returns the next element in the list,
// which will be nil if the end has been reached.
//
// Panics if the collection has been modified since creation of the iterator.
func (i *DListIterator[T]) Next() collections.Element[T] {
	i.validateIterator()

	for {
		if i.direction == forward {
			i.current = i.current.Next()
		} else {
			i.current = i.current.Previous()
		}

		if i.current == nil {
			return i.NilElement
		}

		if i.predicate(i.current.item) {
			return util.NewElementType[T](i.list, &i.current.item)
		}
	}
}

func (i *DListIterator[T]) validateIterator() {
	// util.ValidatePointerNotNil(unsafe.Pointer(i))
	if i.Version != i.list.version {
		panic(messages.COLLECTION_MODIFIED)
	}
}
