package slist

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/local"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
)

// SListIterator implements an iterator over the elements in the list.
type SListIterator[T any] struct {
	util.IteratorBase[T]
	list      *SList[T]
	current   *SListNode[T]
	startNode *SListNode[T]
	predicate functions.PredicateFunc[T]

	local.InternalImpl
}

func newForwardIterator[T any](list *SList[T], predicate functions.PredicateFunc[T]) collections.Iterator[T] {
	return &SListIterator[T]{
		list:      list,
		current:   list.First(),
		startNode: list.First(),
		predicate: predicate,
		IteratorBase: util.IteratorBase[T]{
			Version:    list.version,
			NilElement: nil,
		},
	}
}

// Iterator returns an iterator that walks the list from first to last element
//
//	iter := ll.Iterator()
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (l *SList[T]) Iterator() collections.Iterator[T] {

	return newForwardIterator(l, util.DefaultPredicate[T])
}

// TakeWhile returns a forward iterater that walks the collection returning only
// those elements for which predicate returns true.
//
//	ll := SList.New[int]()
//	// add values
//	iter := ll.TakeWhile(func (val int) bool { return val % 2 == 0 })
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (l *SList[T]) TakeWhile(predicate functions.PredicateFunc[T]) collections.Iterator[T] {

	return newForwardIterator(l, predicate)
}

// Start begins an iteration across the SList returning the fisrt element,
// which will be nil if the collection is empty.
//
// Panics if the collection has been modified since creation of the iterator.
func (i *SListIterator[T]) Start() collections.Element[T] {
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
func (i *SListIterator[T]) Next() collections.Element[T] {
	i.validateIterator()

	for {
		i.current = i.current.Next()

		if i.current == nil {
			return i.NilElement
		}

		if i.predicate(i.current.item) {
			return util.NewElementType[T](i.list, &i.current.item)
		}
	}
}

func (i *SListIterator[T]) validateIterator() {
	// util.ValidatePointerNotNil(unsafe.Pointer(i))
	if i.Version != i.list.version {
		panic(messages.COLLECTION_MODIFIED)
	}
}
