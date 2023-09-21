package ringbuffer

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/internal/util"
)

// Sort performs an in-place sort of this collection.
//
// The item with the smallest value will be placed at the head of the buffer.
func (buf *RingBuffer[T]) Sort() {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))
	buf.doSort(util.Gosort[T])
}

// Sorted returns a sorted copy of this queue as a new queue using the provided [functions.DeepCopyFunc] if any.
func (buf *RingBuffer[T]) Sorted() collections.Collection[T] {

	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.RLock()
		defer buf.lock.RUnlock()
	}

	buf1 := buf.makeDeepCopy()

	if buf1.size > 1 {
		buf1.doSort(util.Gosort[T])
	}

	return buf1
}

// SortDescending performs an in-place decending sort of this collection.
//
// The item with the largest value will be placed at the head of the buffer.
func (buf *RingBuffer[T]) SortDescending() {
	// util.ValidatePointerNotNil(unsafe.Pointer(buf))
	buf.doSort(util.GosortDescending[T])
}

// SortedDescending returns a sorted copy of this ringbuffer as a new ringbuffer using the provided [functions.DeepCopyFunc] if any.
func (buf *RingBuffer[T]) SortedDescending() collections.Collection[T] {

	// util.ValidatePointerNotNil(unsafe.Pointer(buf))

	if buf.lock != nil {
		buf.lock.RLock()
		defer buf.lock.RUnlock()
	}

	buf1 := buf.makeDeepCopy()

	if buf1.size > 1 {
		buf1.doSort(util.GosortDescending[T])
	}

	return buf1
}

func (buf *RingBuffer[T]) doSort(f util.SortFunc[T]) {

	if buf.size <= 1 {
		return
	}

	if buf.lock != nil {
		buf.lock.Lock()
		defer buf.lock.Unlock()
	}

	slc := buf.toSlice(true, false)
	f(slc, buf.size, buf.compare)
	buf.head = 0
	buf.tail = buf.size % buf.maxSize
	buf.buffer = slc
	buf.version++
}
