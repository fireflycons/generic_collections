package ringbuffer

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/local"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
)

type direction int

// Assert interface implementation
var _ collections.Iterable[int] = (*RingBuffer[int])(nil)

const (
	forward, reverse direction = 1, -1
)

// RingBufferIterator implements an iterator over the elements in the queue.
type RingBufferIterator[T any] struct {
	util.IteratorBase[T]
	buffer    *RingBuffer[T]
	index     int
	version   int
	predicate functions.PredicateFunc[T]
	direction direction

	local.InternalImpl
}

func newForwardIterator[T any](buf *RingBuffer[T], predicate functions.PredicateFunc[T]) collections.Iterator[T] {
	return &RingBufferIterator[T]{
		buffer:    buf,
		index:     0,
		direction: forward,
		version:   buf.version,
		predicate: predicate,
		IteratorBase: util.IteratorBase[T]{
			Version:    buf.version,
			NilElement: nil,
		},
	}
}

func newReverseIterator[T any](buf *RingBuffer[T]) collections.Iterator[T] {
	return &RingBufferIterator[T]{
		buffer:    buf,
		direction: reverse,
		version:   buf.version,
		predicate: func(item T) bool { return true },
		IteratorBase: util.IteratorBase[T]{
			Version:    buf.version,
			NilElement: nil,
		},
	}
}

func (buf *RingBuffer[T]) Iterator() collections.Iterator[T] {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))
	return newForwardIterator[T](buf, util.DefaultPredicate[T])
}

func (buf *RingBuffer[T]) ReverseIterator() collections.Iterator[T] {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))
	return newReverseIterator[T](buf)
}

// TakeWhile returns a forward iterater that walks the collection returning only
// those elements for which predicate returns true.
//
//	buf := ringBuffer.New[int]()
//	// add values
//	iter := buf.TakeWhile(func (val int) bool { return val % 2 == 0 })
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (buf *RingBuffer[T]) TakeWhile(predicate functions.PredicateFunc[T]) collections.Iterator[T] {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))
	return newForwardIterator(buf, predicate)
}

// Start begins an iteration across the queue returning the fisrt element,
// which will be nil if the collection is empty.
//
// Panics if the collection has been modified since creation of the iterator.
func (i *RingBufferIterator[T]) Start() collections.Element[T] {
	i.validateIterator()

	if i.buffer.size == 0 {
		return i.NilElement
	}

	i.index = util.Iif(i.direction == forward, 0, i.buffer.size-1)

	valPtr := &(i.buffer.buffer[i.toBufferPosition()])

	if !i.predicate(*valPtr) {
		return i.Next()
	}

	return util.NewElementType[T](i.buffer, valPtr)
}

// Next returns the next element in the collection,
// which will be nil if the end has been reached.
//
// Panics if the collection has been modified since creation of the iterator.
func (i *RingBufferIterator[T]) Next() collections.Element[T] {
	i.validateIterator()

	for {
		i.index += int(i.direction)

		if i.buffer.size == 0 || i.index >= i.buffer.size || i.index < 0 {
			return i.NilElement
		}

		valPtr := &(i.buffer.buffer[i.toBufferPosition()])

		if i.predicate(*valPtr) {
			return util.NewElementType[T](i.buffer, valPtr)
		}
	}
}

func (i *RingBufferIterator[T]) validateIterator() {
	// util.ValidatePointerNotNil(unsafe.Pointer(i))
	if i.Version != i.buffer.version {
		panic(messages.COLLECTION_MODIFIED)
	}
}

func (i *RingBufferIterator[T]) toBufferPosition() int {
	return (i.buffer.head + i.index) % len(i.buffer.buffer)
}
