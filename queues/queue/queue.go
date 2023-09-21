/*
Package queue provides a slice-backed FIFO queue.
*/
package queue

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

// Assert Queue implements required interfaces
var _ queues.Queue[int] = (*Queue[int])(nil)
var _ collections.ReverseIterable[int] = (*Queue[int])(nil)

const (
	growFactor  = 200
	minimumGrow = 4
)

// QueueOptionFunc is the signature of a function
// for providing options to the Queue constructor.
type QueueOptionFunc[T any] func(*Queue[T])

// Queue implements a first-in, first-out collection.
type Queue[T any] struct {
	version         int
	lock            *sync.RWMutex
	head            int
	tail            int
	size            int
	initialCapacity int
	compare         functions.ComparerFunc[T]
	copy            functions.DeepCopyFunc[T]
	buffer          []T
	concurrent      bool

	local.InternalImpl
}

// New creates a new queue.
//
// The queue is backed by a circular slice of T.
func New[T any](options ...QueueOptionFunc[T]) *Queue[T] {
	queue := &Queue[T]{
		initialCapacity: util.DefaultCapacity,
	}

	for _, o := range options {
		o(queue)
	}

	if queue.copy == nil {
		queue.copy = util.DefaultDeepCopy[T]
	}

	queue.buffer = make([]T, queue.initialCapacity)

	if queue.compare == nil {
		queue.compare = util.GetDefaultComparer[T]()
	}

	return queue
}

// Option function for New to make the collection thread-safe. Adds overhead.
func WithThreadSafe[T any]() QueueOptionFunc[T] {
	return func(s *Queue[T]) {
		s.lock = &sync.RWMutex{}
	}
}

// Option function to enable concurrency feature
func WithConcurrent[T any]() QueueOptionFunc[T] {
	return func(q *Queue[T]) {
		q.concurrent = true
	}
}

// Option function to set initial capacity to
// something other than the default 16 elements.
func WithCapacity[T any](capacity int) QueueOptionFunc[T] {
	if capacity < 0 {
		panic(messages.NEGATIVE_CAPACITY)
	}
	return func(q *Queue[T]) {
		q.initialCapacity = capacity
	}
}

// Option function to provide a comparer function for values of type T.
// Required if the element type is not a supported type.
func WithComparer[T any](comparer functions.ComparerFunc[T]) QueueOptionFunc[T] {
	if comparer == nil {
		panic(messages.COMP_FN_NIL)
	}
	return func(q *Queue[T]) {
		q.compare = comparer
	}
}

// Option func to provide a deep copy implementation for collection elements.
func WithDeepCopy[T any](copier functions.DeepCopyFunc[T]) QueueOptionFunc[T] {
	// Can be nil
	return func(q *Queue[T]) {
		q.copy = copier
	}
}

// Add enqueues a value in the queue.
//
// Always returns true.
func (q *Queue[T]) Add(value T) bool {

	q.Enqueue(value)
	return true
}

// AddCollection adds the values of the given collection to the end of this queue.
//
// Values are added in the order defined by the other collection.
func (q *Queue[T]) AddCollection(collection collections.Collection[T]) {

	q.AddRange(collection.ToSliceDeep())
}

// AddRange enqueues the values in the given slice
func (q *Queue[T]) AddRange(values []T) {

	if len(values) == 0 {
		return
	}

	if q.lock != nil {
		q.lock.Lock()
		defer q.lock.Unlock()
	}

	var newBufferSize int

	lv := len(values)
	newSize := q.size + lv

	if len(q.buffer) < newSize {
		newBufferSize = util.Iif(newSize > util.DefaultCapacity, newSize*growFactor/100, util.DefaultCapacity)
	} else {
		newBufferSize = len(q.buffer)
	}

	if q.size == 0 && q.head == 0 {
		q.buffer = make([]T, newBufferSize)
		copy(q.buffer, values)
		q.size = lv
		q.head = 0
		q.tail = util.Iif(q.size == len(q.buffer), 0, q.size)
		return
	}

	// Realloc buffer and shift items to front of slice
	q.setLength(newBufferSize)
	util.PartialCopy(values, 0, q.buffer, q.size, lv)
	q.size += lv
	q.tail = util.Iif(q.size == len(q.buffer), 0, q.size)
}

// Clear removes all values from the queue.
func (q *Queue[T]) Clear() {

	if q.lock != nil {
		q.lock.Lock()
		defer q.lock.Unlock()
	}

	q.buffer = make([]T, cap(q.buffer))
	q.head = 0
	q.tail = 0
	q.size = 0
}

// Contains returns true if the given value is in the queue; else false
func (q *Queue[T]) Contains(value T) bool {

	if q.lock != nil {
		q.lock.RLock()
		defer q.lock.RUnlock()
	}

	return q.find(value) != -1
}

// Count returns the number of elements in the queue
func (q *Queue[T]) Count() int {
	return q.size
}

// IsEmpty returns true if the collection has no elements
func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}

// Dequeue removes the value at the front of the queue and returns it.
//
// Panics if the queue is empty.
func (q *Queue[T]) Dequeue() T {

	if q.lock != nil {
		q.lock.Lock()
		defer q.lock.Unlock()
	}

	if q.size == 0 {
		panic(messages.COLLECTION_EMPTY)
	}

	return q.removeItem()
}

// TryDequeue removes and returns the value at the front of the queue and true if
// the queue is not empty; else zero value of T and false.
func (q *Queue[T]) TryDequeue() (T, bool) {

	if q.lock != nil {
		q.lock.Lock()
		defer q.lock.Unlock()
	}

	if q.size == 0 {
		var empty T
		return empty, false
	}

	return q.removeItem(), true
}

// Enqueue adds a value to the back of the queue.
func (q *Queue[T]) Enqueue(value T) {

	if q.lock != nil {
		q.lock.Lock()
		defer q.lock.Unlock()
	}

	q.enqueue(value)
}

// Peek returns the value at the front of the queue without removing it.
//
// Panics if the queue is empty.
func (q *Queue[T]) Peek() T {

	if q.lock != nil {
		q.lock.RLock()
		defer q.lock.RUnlock()
	}

	if q.size == 0 {
		panic(messages.COLLECTION_EMPTY)
	}

	return q.buffer[q.head]
}

// TryPeek returns the value at the front of the queue and true if
// the queue is not empty; else zero value of T and false.
func (q *Queue[T]) TryPeek() (T, bool) {

	if q.lock != nil {
		q.lock.RLock()
		defer q.lock.RUnlock()
	}

	if q.size == 0 {
		var empty T
		return empty, false
	}

	return q.buffer[q.head], true
}

// Remove removes the first occurrence of the given value from the queue, searching from front.
//
// Returns true if the value was present and was removed; else false.
func (q *Queue[T]) Remove(value T) bool {

	if q.size == 0 {
		return false
	}

	if q.lock != nil {
		q.lock.Lock()
		defer q.lock.Unlock()
	}

	index := q.find(value)

	if index == -1 {
		return false
	}

	var empty T
	q.buffer[index] = empty
	q.size--
	buf := make([]T, len(q.buffer))

	switch {
	case q.head < q.tail || (q.head == 0 && q.tail == 0):
		util.PartialCopy(q.buffer, 0, buf, 0, index)
		util.PartialCopy(q.buffer, index+1, buf, index, len(q.buffer)-index-1)
		if q.head < q.tail {
			q.tail--
		}
	case index == q.head:
		q.head++
		headToEnd := q.size - q.head + (q.head - q.tail)
		util.PartialCopy(q.buffer, q.head, buf, 0, headToEnd)
		util.PartialCopy(q.buffer, 0, buf, headToEnd, q.tail)
		q.head = 0
		q.tail = q.size
	case index > q.head:
		removedToEnd := q.size - index + util.Iif(q.tail < q.head, q.head-q.tail, 0)
		util.PartialCopy(q.buffer, q.head, buf, 0, index-q.head)
		util.PartialCopy(q.buffer, index+1, buf, index-q.head, removedToEnd)
		util.PartialCopy(q.buffer, 0, buf, index-q.head+removedToEnd, q.tail)
		q.head = 0
		q.tail = q.size
	case q.head > q.tail && q.tail > index:
		headToEnd := q.size - q.head + (q.head - q.tail + 1)
		util.PartialCopy(q.buffer, q.head, buf, 0, headToEnd)
		util.PartialCopy(q.buffer, 0, buf, headToEnd, index)
		util.PartialCopy(q.buffer, index+1, buf, headToEnd+index, q.tail-index-1)
		q.head = 0
		q.tail = q.size
	default:
		panic(fmt.Sprintf("BUG: Queue.Remove - Head: %d, tail %d, index: %d", q.head, q.tail, index))
	}

	q.buffer = buf
	return true
}

// ToSlice returns a copy of the queue content as a slice
func (q *Queue[T]) ToSlice() []T {

	if q.lock != nil {
		q.lock.Lock()
		defer q.lock.Unlock()
	}
	return q.toSlice(false)
}

// ToSliceDeep returns a copy of the queue content as a slice using the provided [functions.DeepCopyFunc] if any.
func (q *Queue[T]) ToSliceDeep() []T {

	if q.lock != nil {
		q.lock.Lock()
		defer q.lock.Unlock()
	}

	return q.toSlice(false)
}

func (q *Queue[T]) toSlice(deepCopy bool) []T {
	slc := make([]T, q.size)
	q.copyTo(slc, deepCopy)
	return slc
}

// Type returns the type of this collection
func (*Queue[T]) Type() collections.CollectionType {
	return collections.COLLECTION_QUEUE
}

// String returns a string representation of container
func (q *Queue[T]) String() string {

	if q.lock != nil {
		q.lock.RLock()
		defer q.lock.RUnlock()
	}

	var values []string
	for _, value := range q.toSlice(false) {
		values = append(values, fmt.Sprintf("%v", value))
	}

	return "Queue\n" + strings.Join(values, ", ")
}

func (q *Queue[T]) find(value T) int {
	index := q.head

	for count := q.size; count > 0; count-- {
		if q.compare(q.buffer[index], value) == 0 {
			return index
		}

		index = (index + 1) % len(q.buffer)
	}

	return -1
}

// Reallocate the queue buffer, moving the head to to beginning of the slice.
func (q *Queue[T]) setLength(capacity int) {
	newBuffer := make([]T, capacity)

	q.copyTo(newBuffer, false)
	q.buffer = newBuffer
	q.head = 0
	q.tail = util.Iif(q.size == capacity, 0, q.size)
	q.version++
}

func (q *Queue[T]) copyTo(slc []T, deepCopy bool) {
	if q.size > 0 {
		if q.head < q.tail || (q.head == 0 && q.tail == 0) {
			if deepCopy {
				util.DeepCopySlice(slc, q.buffer, q.copy)
			} else {
				copy(slc, q.buffer)
			}
		} else {
			headToEnd := q.size - q.head + (q.head - q.tail)
			if deepCopy {
				util.PartialCopyDeep(q.buffer, q.head, slc, 0, headToEnd, q.copy)
				util.PartialCopyDeep(q.buffer, 0, slc, headToEnd, q.tail, q.copy)
			} else {
				util.PartialCopy(q.buffer, q.head, slc, 0, headToEnd)
				util.PartialCopy(q.buffer, 0, slc, headToEnd, q.tail)
			}
		}
	}
}

func (q *Queue[T]) enqueue(value T) {
	if q.size == len(q.buffer) {
		newCapacity := len(q.buffer) * growFactor / 100
		if newCapacity < len(q.buffer)+minimumGrow {
			newCapacity = len(q.buffer) + minimumGrow
		}
		q.setLength(newCapacity)
	}

	q.buffer[q.tail] = value
	q.tail = (q.tail + 1) % len(q.buffer)
	q.size++
	q.version++
}

func (q *Queue[T]) removeItem() T {
	var empty T
	removed := q.buffer[q.head]
	q.buffer[q.head] = empty
	q.head = (q.head + 1) % len(q.buffer)
	q.size--
	q.version++
	return removed
}

func (q *Queue[T]) makeDeepCopy() *Queue[T] {
	other := &Queue[T]{
		head:            q.head,
		tail:            q.tail,
		size:            q.size,
		version:         0,
		initialCapacity: q.initialCapacity,
		compare:         q.compare,
		copy:            q.copy,
	}

	if q.lock != nil {
		other.lock = &sync.RWMutex{}
	}

	other.buffer = make([]T, len(q.buffer), cap(q.buffer))
	util.DeepCopySlice(other.buffer, q.buffer, q.copy)
	return other
}
