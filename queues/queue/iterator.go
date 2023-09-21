package queue

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/local"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
)

type direction int

// Assert interface implementation
var _ collections.Iterable[int] = (*Queue[int])(nil)

const (
	forward, reverse direction = 1, -1
)

// DListIterator implements an iterator over the elements in the queue.
type QueueIterator[T any] struct {
	util.IteratorBase[T]
	queue     *Queue[T]
	index     int
	version   int
	predicate functions.PredicateFunc[T]
	direction direction

	local.InternalImpl
}

func newForwardIterator[T any](q *Queue[T], predicate functions.PredicateFunc[T]) collections.Iterator[T] {
	return &QueueIterator[T]{
		queue:     q,
		index:     0,
		direction: forward,
		version:   q.version,
		predicate: predicate,
		IteratorBase: util.IteratorBase[T]{
			Version:    q.version,
			NilElement: nil,
		},
	}
}

func newReverseIterator[T any](q *Queue[T]) collections.Iterator[T] {
	return &QueueIterator[T]{
		queue:     q,
		direction: reverse,
		version:   q.version,
		predicate: func(item T) bool { return true },
		IteratorBase: util.IteratorBase[T]{
			Version:    q.version,
			NilElement: nil,
		},
	}
}

// Iterator retuns an iterator that walks the Queue from head to tail element
//
//	iter := q.Iterator()
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (q *Queue[T]) Iterator() collections.Iterator[T] {

	return newForwardIterator(q, util.DefaultPredicate[T])
}

// ReverseIterator retuns an iterator that walks the Queue from tail to head element
//
//	iter := q.ReverseIterator()
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (q *Queue[T]) ReverseIterator() collections.Iterator[T] {

	return newReverseIterator(q)
}

// TakeWhile returns a forward iterater that walks the collection returning only
// those elements for which predicate returns true.
//
//	queue := queue.New[int]()
//	// add values
//	iter := queue.TakeWhile(func (val int) bool { return val % 2 == 0 })
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (q *Queue[T]) TakeWhile(predicate functions.PredicateFunc[T]) collections.Iterator[T] {

	return newForwardIterator(q, predicate)
}

// Start begins an iteration across the queue returning the fisrt element,
// which will be nil if the collection is empty.
//
// Panics if the collection has been modified since creation of the iterator.
func (i *QueueIterator[T]) Start() collections.Element[T] {
	i.validateIterator()

	if i.queue.size == 0 {
		return i.NilElement
	}

	i.index = util.Iif(i.direction == forward, 0, i.queue.size-1)
	valPtr := &(i.queue.buffer[i.toBufferPosition()])

	if !i.predicate(*valPtr) {
		return i.Next()
	}

	return util.NewElementType[T](i.queue, valPtr)
}

// Next returns the next element in the collection,
// which will be nil if the end has been reached.
//
// Panics if the collection has been modified since creation of the iterator.
func (i *QueueIterator[T]) Next() collections.Element[T] {
	i.validateIterator()

	for {
		i.index += int(i.direction)

		if i.queue.size == 0 || i.index >= i.queue.size || i.index < 0 {
			return i.NilElement
		}

		valPtr := &(i.queue.buffer[i.toBufferPosition()])

		if i.predicate(*valPtr) {
			return util.NewElementType[T](i.queue, valPtr)
		}
	}
}

func (i *QueueIterator[T]) validateIterator() {
	// util.ValidatePointerNotNil(unsafe.Pointer(i))
	if i.Version != i.queue.version {
		panic(messages.COLLECTION_MODIFIED)
	}
}

func (i *QueueIterator[T]) toBufferPosition() int {
	return (i.queue.head + i.index) % len(i.queue.buffer)
}
