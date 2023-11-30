/*
Package ringbuffer provides a fixed length FIFO circular buffer.
*/
package ringbuffer

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/local"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/fireflycons/generic_collections/queues"
)

var _ queues.Queue[int] = (*RingBuffer[int])(nil)
var _ collections.ReverseIterable[int] = (*RingBuffer[int])(nil)

// RingBufferOptionFunc is the signature of a function
// for providing options to the RingBuffer constructor.
type RingBufferOptionFunc[T any] func(*RingBuffer[T])

// RingBuffer implements a first-in, first-out collection,
// of fixed size. When the buffer is full, items added to
// the end displace items at the front.
type RingBuffer[T any] struct {
	version int
	lock    *sync.RWMutex
	head    int
	tail    int
	full    bool
	maxSize int
	size    int
	compare functions.ComparerFunc[T]
	copy    functions.DeepCopyFunc[T]
	buffer  []T

	local.InternalImpl
}

// New instantiates a new empty buffer with the specified size of maximum number of elements that it can hold.
// This max size of the buffer cannot be changed.
func New[T any](maxSize int, options ...RingBufferOptionFunc[T]) *RingBuffer[T] {
	if maxSize < 1 {
		panic("Invalid maxSize, should be at least 1")
	}

	buf := &RingBuffer[T]{
		maxSize: maxSize,
		buffer:  make([]T, maxSize),
	}

	for _, o := range options {
		o(buf)
	}

	if buf.copy == nil {
		buf.copy = util.DefaultDeepCopy[T]
	}

	if buf.compare == nil {
		buf.compare = util.GetDefaultComparer[T]()
	}

	return buf
}

// Option function for New to provide a comparer function for values of type T.
// Required if the element type is not numeric, bool, pointer or string.
func WithComparer[T any](comparer functions.ComparerFunc[T]) RingBufferOptionFunc[T] {
	if comparer == nil {
		panic(messages.COMP_FN_NIL)
	}
	return func(s *RingBuffer[T]) {
		s.compare = comparer
	}
}

// Option function for New to make the collection thread-safe. Adds overhead.
func WithThreadSafe[T any]() RingBufferOptionFunc[T] {
	return func(s *RingBuffer[T]) {
		s.lock = &sync.RWMutex{}
	}
}

// Option func to provide a deep copy implementation for collection elements.
func WithDeepCopy[T any](copier functions.DeepCopyFunc[T]) RingBufferOptionFunc[T] {
	// Can be nil
	return func(buf *RingBuffer[T]) {
		buf.copy = copier
	}
}

// Add enqueues a value in the buffer. It is an alias for Enqueue.
//
// Always returns true.
func (buf *RingBuffer[T]) Add(value T) bool {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))
	buf.Enqueue(value)
	return true
}

// AddCollection adds the values of the given collection to the end of this buffer.
//
// If the collection is larger or equal to the size of the buffer,
// then all existing elements will be displaced and the end-most
// portion of the collection (governed by any ordering of the collection)
// that will fit into the buffer will be the values enqueued.
// If the collection is smaller than the buffer, then
// all collection elements will be enqueued, displacing elements from the
// head of the buffer as necessary.
func (buf *RingBuffer[T]) AddCollection(collection collections.Collection[T]) {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))
	buf.AddRange(collection.ToSliceDeep())
}

// AddRange enqueues the values in the given slice.
//
// If the slice is larger or equal to the size of the buffer,
// then all existing elements will be displaced and the end-most
// portion of the slice that will fit into the buffer will be the
// values enqueued. If the slice is smaller than the buffer, then
// all slice elements will be enqueued, displacing elements from the
// head of the buffer as necessary.
func (buf *RingBuffer[T]) AddRange(values []T) {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))
	if len(values) == 0 {
		return
	}

	if buf.lock != nil {
		buf.lock.Lock()
		defer buf.lock.Unlock()
	}

	if len(values) >= buf.maxSize {
		// Buffer will be filled from incoming slice and any
		// existing values completely displaced
		startIndex := len(values) - buf.maxSize
		util.PartialCopy(values, startIndex, buf.buffer, 0, buf.maxSize)
		buf.full = true
		buf.size = buf.maxSize
		buf.head = 0
		buf.tail = 0
	} else {
		for _, v := range values {
			if buf.full {
				buf.head = (buf.head + 1) % buf.maxSize
			}
			buf.buffer[buf.tail] = v
			buf.tail = (buf.tail + 1) % buf.maxSize
			if buf.tail == buf.head {
				buf.full = true
			} else {
				buf.size++
			}
		}
	}

	buf.version++
}

// Contains returns true if the given value is in the queue; else false
//
// Up to O(n).
func (buf *RingBuffer[T]) Contains(value T) bool {

	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.RLock()
		defer buf.lock.RUnlock()
	}

	return buf.find(value) != -1
}

func (buf *RingBuffer[T]) CurrentVersion() int {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))
	return buf.version
}

// Enqueue adds a value to the end of the buffer
//
// If the buffer is full, the item at the head
// is discarded.
func (buf *RingBuffer[T]) Enqueue(value T) {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.Lock()
		defer buf.lock.Unlock()
	}

	buf.enqueue(value)
}

func (buf *RingBuffer[T]) enqueue(value T) {

	if buf.full {
		// increments version
		buf.removeHead()
	} else {
		buf.version++
	}

	buf.append(value)
}

// Offer offers a value to the buffer.
//
// If the buffer is full, then false is returned;
// else the value is enqueued and true is returned.
func (buf *RingBuffer[T]) Offer(value T) bool {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.Lock()
		defer buf.lock.Unlock()
	}

	if buf.full {
		return false
	}

	buf.append(value)
	buf.version++
	return true
}

// Dequeue removes first element of the buffer and returns it, or nil if buffer is empty.
// Second return parameter is true, unless the buffer was empty and there was nothing to dequeue.
func (buf *RingBuffer[T]) Dequeue() T {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.Lock()
		defer buf.lock.Unlock()
	}

	return buf.removeHead()
}

// TryDequeue removes and returns the value at the front of the buffer and true if
// the buffer is not empty; else zero value of T and false.
func (buf *RingBuffer[T]) TryDequeue() (T, bool) {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.Lock()
		defer buf.lock.Unlock()
	}

	if buf.size == 0 {
		var empty T
		return empty, false
	}

	return buf.removeHead(), true
}

// Peek returns the value at the front of the buffer without removing it.
//
// Panics if the buffer is empty.
func (buf *RingBuffer[T]) Peek() T {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))
	if buf.size == 0 {
		panic(messages.COLLECTION_EMPTY)
	}
	return buf.buffer[buf.head]
}

// TryPeek returns the value at the front of the buffer and true if
// the buffer is not empty; else zero value of T and false.
func (buf *RingBuffer[T]) TryPeek() (T, bool) {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.RLock()
		defer buf.lock.RUnlock()
	}

	if buf.size == 0 {
		var empty T
		return empty, false
	}

	return buf.buffer[buf.head], true
}

// Remove removes the first occurrence of the given value from the buffer, searching from front.
//
// Returns true if the value was present and was removed; else false.
func (buf *RingBuffer[T]) Remove(value T) bool {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.Lock()
		defer buf.lock.Unlock()
	}

	index := buf.find(value)

	if index == -1 {
		return false
	}

	var empty T
	buf.buffer[index] = empty
	buf.size--
	newBuffer := make([]T, buf.maxSize)

	switch {
	case buf.head < buf.tail || (buf.head == 0 && buf.tail == 0):
		util.PartialCopy(buf.buffer, 0, newBuffer, 0, index)
		util.PartialCopy(buf.buffer, index+1, newBuffer, index, len(buf.buffer)-index-1)
		if buf.head < buf.tail {
			buf.tail--
		}
	case index == buf.head:
		buf.head++
		headToEnd := buf.size - buf.head + (buf.head - buf.tail)
		util.PartialCopy(buf.buffer, buf.head, newBuffer, 0, headToEnd)
		util.PartialCopy(buf.buffer, 0, newBuffer, headToEnd, buf.tail)
		buf.head = 0
		buf.tail = buf.size
	case index > buf.head:
		removedToEnd := buf.size - index + util.Iif(buf.tail < buf.head, buf.head-buf.tail, 0)
		util.PartialCopy(buf.buffer, buf.head, newBuffer, 0, index-buf.head)
		util.PartialCopy(buf.buffer, index+1, newBuffer, index-buf.head, removedToEnd)
		util.PartialCopy(buf.buffer, 0, newBuffer, index-buf.head+removedToEnd, buf.tail)
		buf.head = 0
		buf.tail = buf.size
	case buf.head > buf.tail && buf.tail > index:
		headToEnd := buf.size - buf.head + (buf.head - buf.tail + 1)
		util.PartialCopy(buf.buffer, buf.head, newBuffer, 0, headToEnd)
		util.PartialCopy(buf.buffer, 0, newBuffer, headToEnd, index)
		util.PartialCopy(buf.buffer, index+1, newBuffer, headToEnd+index, buf.tail-index-1)
		buf.head = 0
		buf.tail = buf.size
	default:
		panic(fmt.Sprintf("BUG: RingBuffer.Remove - Head: %d, tail %d, index: %d", buf.head, buf.tail, index))
	}

	buf.buffer = newBuffer
	buf.full = false
	return true

}

// Empty returns true if buffer does not contain any elements.
func (buf *RingBuffer[T]) Empty() bool {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))
	return buf.size == 0
}

// Full returns true if the buffer is full, i.e. has reached the maximum number of elements that it can hold.
func (buf *RingBuffer[T]) Full() bool {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))
	return buf.size == buf.maxSize
}

// Count returns number of elements within the buffer.
func (buf *RingBuffer[T]) Count() int {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))
	return buf.size
}

// IsEmpty returns true if the collection has no elements.
func (buf *RingBuffer[T]) IsEmpty() bool {
	return buf.size == 0
}

// Clear removes all elements from the buffer.
func (buf *RingBuffer[T]) Clear() {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.Lock()
		defer buf.lock.Unlock()
	}

	buf.buffer = make([]T, buf.maxSize)
	buf.head = 0
	buf.tail = 0
	buf.full = false
	buf.size = 0
}

// ToSlice returns a copy of the buffer content as a slice
//
// O(n).
func (buf *RingBuffer[T]) ToSlice() []T {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.RLock()
		defer buf.lock.RUnlock()
	}

	return buf.toSlice(false, false)
}

// ToSlice returns a copy of the buffer content as a slice
//
// O(n).
func (buf *RingBuffer[T]) ToSliceDeep() []T {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.RLock()
		defer buf.lock.RUnlock()
	}

	return buf.toSlice(false, true)
}

// String returns a string representation of container.
func (buf *RingBuffer[T]) String() string {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.RLock()
		defer buf.lock.RUnlock()
	}

	str := "RingBuffer\n"
	var values []string
	for _, value := range buf.toSlice(false, false) {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

// Type returns the type of this collection.
func (*RingBuffer[T]) Type() collections.CollectionType {
	return collections.COLLECTION_RINGBUFFER
}

func (buf *RingBuffer[T]) append(value T) {
	buf.buffer[buf.tail] = value
	buf.tail = (buf.tail + 1) % buf.maxSize

	if buf.tail == buf.head {
		buf.full = true
	}

	buf.size = buf.calculateSize()
}

func (buf *RingBuffer[T]) find(value T) int {
	index := buf.head

	for count := buf.size; count > 0; count-- {
		if buf.compare(buf.buffer[index], value) == 0 {
			return index
		}

		index = (index + 1) % buf.maxSize
	}

	return -1
}

func (buf *RingBuffer[T]) removeHead() T {

	if buf.size == 0 {
		panic(messages.COLLECTION_EMPTY)
	}

	var empty T

	value := buf.buffer[buf.head]
	buf.buffer[buf.head] = empty

	buf.head = buf.head + 1
	if buf.head >= buf.maxSize {
		buf.head = 0
	}

	buf.full = false
	buf.size = buf.size - 1

	buf.version++
	return value
}

func (buf *RingBuffer[T]) toSlice(keepCapacity, deepCopy bool) []T {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	values := make([]T, util.Iif(keepCapacity, buf.maxSize, buf.size))
	for i := 0; i < buf.size; i++ {
		if deepCopy {
			values[i] = util.DeepCopy(buf.buffer[(buf.head+i)%buf.maxSize], buf.copy)
		} else {
			values[i] = buf.buffer[(buf.head+i)%buf.maxSize]
		}
	}
	return values
}

func (buf *RingBuffer[T]) calculateSize() int {
	if buf.tail < buf.head {
		return buf.maxSize - buf.head + buf.tail
	} else if buf.tail == buf.head {
		if buf.full {
			return buf.maxSize
		}
		return 0
	}
	return buf.tail - buf.head
}

func (buf *RingBuffer[T]) makeDeepCopy() *RingBuffer[T] {
	other := &RingBuffer[T]{
		head:    buf.head,
		tail:    buf.tail,
		size:    buf.size,
		version: 0,
		maxSize: buf.maxSize,
		full:    buf.full,
		compare: buf.compare,
	}

	other.buffer = make([]T, len(buf.buffer), cap(buf.buffer))

	if buf.lock != nil {
		other.lock = &sync.RWMutex{}
	}

	util.DeepCopySlice(other.buffer, buf.buffer, buf.copy)
	return other
}
