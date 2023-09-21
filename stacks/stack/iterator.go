package stack

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/local"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
)

type direction int

const (
	forward, reverse direction = -1, 1
)

// StackIterator implements an iterator over the elements in the stack.
type StackIterator[T any] struct {
	util.IteratorBase[T]
	stack     *Stack[T]
	index     int
	predicate functions.PredicateFunc[T]
	direction direction

	local.InternalImpl
}

func newForwardIterator[T any](stack *Stack[T], predicate functions.PredicateFunc[T]) collections.Iterator[T] {
	return &StackIterator[T]{
		stack:     stack,
		index:     stack.size - 1,
		direction: forward,
		predicate: predicate,
		IteratorBase: util.IteratorBase[T]{
			Version:    stack.version,
			NilElement: nil,
		},
	}
}

func newReverseIterator[T any](stack *Stack[T]) collections.Iterator[T] {
	return &StackIterator[T]{
		stack:     stack,
		index:     0,
		direction: reverse,
		predicate: func(item T) bool { return true },
		IteratorBase: util.IteratorBase[T]{
			Version:    stack.version,
			NilElement: nil,
		},
	}
}

// Iterator retuns an iterator that walks the stack from first to last element
//
//	iter := stack.Iterator()
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (s *Stack[T]) Iterator() collections.Iterator[T] {

	return newForwardIterator(s, util.DefaultPredicate[T])
}

// ReverseIterator retuns an iterator that walks the stack from last to first element
//
//	iter := stack.ReverseIterator()
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (s *Stack[T]) ReverseIterator() collections.Iterator[T] {

	return newReverseIterator(s)
}

// TakeWhile returns a forward iterater that walks the collection returning only
// those elements for which predicate returns true.
//
//	stack := hashset.New[int]()
//	// add values
//	iter := stack.TakeWhile(func (val int) bool { return val % 2 == 0 })
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (s *Stack[T]) TakeWhile(predicate functions.PredicateFunc[T]) collections.Iterator[T] {

	return newForwardIterator(s, predicate)
}

// Start begins iteration across the stack returning the fisrt element,
// which will be nil if the stack is empty.
//
// Panics if the underlying stack is modified between iteration creation and call to Start()
func (i *StackIterator[T]) Start() collections.Element[T] {
	i.validateIterator()

	if i.stack.size == 0 {
		return i.NilElement
	}

	i.index = util.Iif(i.direction == reverse, 0, i.stack.size-1)

	valPtr := &i.stack.buffer[i.index]

	if !i.predicate(*valPtr) {
		return i.Next()
	}

	return util.NewElementType[T](i.stack, valPtr)
}

// Next returns the next element from the iterator,
// which will be nil if the end has been reached.
//
// Panics if the underlying stack is modified between calls to Next.
func (i *StackIterator[T]) Next() collections.Element[T] {
	i.validateIterator()

	for {
		i.index += int(i.direction)

		if i.stack.size == 0 || i.index >= i.stack.size || i.index < 0 {
			return i.NilElement
		}

		valPtr := &i.stack.buffer[i.index]

		if i.predicate(*valPtr) {
			return util.NewElementType[T](i.stack, valPtr)
		}
	}
}

func (i *StackIterator[T]) validateIterator() {
	// util.ValidatePointerNotNil(unsafe.Pointer(i))
	if i.Version != i.stack.version {
		panic(messages.COLLECTION_MODIFIED)
	}
}
