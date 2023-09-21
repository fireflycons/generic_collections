package queue

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/internal/util"
)

// Sort performs an in-place sort of this collection.
//
// The item with the smallest value will be placed at the head of the queue.
func (q *Queue[T]) Sort() {

	q.doSort(util.Gosort[T])

}

// Sorted returns a sorted copy of this queue as a new queue using the provided [functions.DeepCopyFunc] if any.
func (q *Queue[T]) Sorted() collections.Collection[T] {

	if q.lock != nil {
		q.lock.RLock()
		defer q.lock.RUnlock()
	}

	q1 := q.makeDeepCopy()

	if q1.size > 1 {
		q1.doSort(util.Gosort[T])
	}

	return q1
}

// Sort performs an in-place sort of this collection.
//
// The item with the largest value will be placed at the head of the queue.
func (q *Queue[T]) SortDescending() {

	q.doSort(util.GosortDescending[T])
}

// SortedDescending returns a descending order sorted copy of this queue as a new queue using the provided [functions.DeepCopyFunc] if any.
func (q *Queue[T]) SortedDescending() collections.Collection[T] {

	if q.lock != nil {
		q.lock.RLock()
		defer q.lock.RUnlock()
	}

	q1 := q.makeDeepCopy()

	if q1.size > 1 {
		q1.doSort(util.GosortDescending[T])
	}

	return q1
}

func (q *Queue[T]) doSort(f util.SortFunc[T]) {

	if q.size <= 1 {
		return
	}

	if q.lock != nil {
		q.lock.Lock()
		defer q.lock.Unlock()
	}

	length := len(q.buffer)
	slc := make([]T, length)
	q.copyTo(slc, true)
	f(slc, q.size, q.compare)
	q.head = 0
	q.tail = util.Iif(q.size == length, 0, q.size)
	q.buffer = slc
	q.version++
}
