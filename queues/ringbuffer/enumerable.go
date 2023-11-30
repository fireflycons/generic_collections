package ringbuffer

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
)

// Assert interface implementation.
var _ collections.Enumerable[int] = (*RingBuffer[int])(nil)

// Any returns true for the first element found where the predicate function returns true.
// It returns false if no element matches the predicate.
func (buf *RingBuffer[T]) Any(predicate functions.PredicateFunc[T]) bool {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.Lock()
		defer buf.lock.Unlock()
	}

	iter := newForwardIterator[T](buf, predicate)

	return iter.Start() != nil
}

// All applies the predicate function to every element in the collection,
// and returns true if all elements match the predicate.
func (buf *RingBuffer[T]) All(predicate functions.PredicateFunc[T]) bool {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.Lock()
		defer buf.lock.Unlock()
	}

	iter := newForwardIterator[T](buf, util.DefaultPredicate[T])

	for e := iter.Start(); e != nil; e = iter.Next() {
		if !predicate(e.Value()) {
			return false
		}
	}

	return true
}

// ForEach applies function f to all elements in the collection.
func (buf *RingBuffer[T]) ForEach(f func(collections.Element[T])) {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.Lock()
		defer buf.lock.Unlock()
	}

	iter := newForwardIterator[T](buf, util.DefaultPredicate[T])

	for e := iter.Start(); e != nil; e = iter.Next() {
		f(e)
	}
}

// Map applies function f to all elements in the collection
// and returns a new RingBuffer containing the result of f.
func (buf *RingBuffer[T]) Map(f func(T) T) collections.Collection[T] {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.Lock()
		defer buf.lock.Unlock()
	}

	iter := newForwardIterator[T](buf, util.DefaultPredicate[T])

	buf1 := New[T](buf.maxSize, WithComparer[T](buf.compare))

	for e := iter.Start(); e != nil; e = iter.Next() {
		buf1.enqueue(f(e.Value()))
	}

	return buf1
}

// Select returns a new RingBuffer containing only the items for which predicate is true.
// Buffer capacity is the same as the source buffer.
func (buf *RingBuffer[T]) Select(predicate functions.PredicateFunc[T]) collections.Collection[T] {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.RLock()
		defer buf.lock.RUnlock()
	}

	return buf.doSelect(predicate, false)
}

// Select returns a new RingBuffer containing only the items for which predicate is true.
// Buffer capacity is the same as the source buffer.
//
// Elements are deep copied to the new collection using the provided [functions.DeepCopyFunc] if any.
func (buf *RingBuffer[T]) SelectDeep(predicate functions.PredicateFunc[T]) collections.Collection[T] {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.RLock()
		defer buf.lock.RUnlock()
	}

	return buf.doSelect(predicate, true)
}

// Find finds the first occurrence of an element matching the predicate.
//
// The function returns nil if no match.
func (buf *RingBuffer[T]) Find(predicate functions.PredicateFunc[T]) collections.Element[T] {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.RLock()
		defer buf.lock.RUnlock()
	}

	result := buf.doFind(predicate, false)

	if len(result) == 0 {
		return nil
	}

	return result[0]
}

// FindAll finds all occurrences of an element matching the predicate.
//
// The function returns an empty slice if none match.
func (buf *RingBuffer[T]) FindAll(predicate functions.PredicateFunc[T]) []collections.Element[T] {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.RLock()
		defer buf.lock.RUnlock()
	}

	result := buf.doFind(predicate, true)

	return result
}

// Min returns the minimum value in the collection according to the Comparer function.
func (buf *RingBuffer[T]) Min() T {

	if buf.lock != nil {
		buf.lock.RLock()
		defer buf.lock.RUnlock()
	}

	if buf.size == 0 {
		panic(messages.COLLECTION_EMPTY)
	}

	m := buf.buffer[buf.head]
	l := len(buf.buffer)

	for i := 0; i < buf.size; i++ {
		ind := (buf.head + i) % l
		v := buf.buffer[ind]
		if buf.compare(m, v) > 0 {
			m = v
		}

	}

	return m
}

// Max returns the maximum value in the collection according to the Comparer function.
func (buf *RingBuffer[T]) Max() T {

	if buf.lock != nil {
		buf.lock.RLock()
		defer buf.lock.RUnlock()
	}

	if buf.size == 0 {
		panic(messages.COLLECTION_EMPTY)
	}

	m := buf.buffer[buf.head]
	l := len(buf.buffer)

	for i := 0; i < buf.size; i++ {
		ind := (buf.head + i) % l
		v := buf.buffer[ind]
		if buf.compare(m, v) < 0 {
			m = v
		}

	}

	return m
}

func (q *RingBuffer[T]) doFind(predicate functions.PredicateFunc[T], all bool) []collections.Element[T] {

	iter := newForwardIterator[T](q, predicate)
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

// Select returns a new RingBuffer containing only the items for which predicate is true.
// Buffer capacity is the same as the source buffer.
func (buf *RingBuffer[T]) doSelect(predicate functions.PredicateFunc[T], deepCopy bool) collections.Collection[T] {

	buf1 := New[T](buf.maxSize, WithComparer[T](buf.compare))
	iter := newForwardIterator[T](buf, predicate)

	for e := iter.Start(); e != nil; e = iter.Next() {
		if deepCopy {
			buf1.enqueue(buf.copy(e.Value()))
		} else {
			buf1.enqueue(e.Value())
		}
	}

	return buf1
}
