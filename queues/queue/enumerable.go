package queue

import (
	"sync"

	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
)

// Assert interface implementation
var _ collections.Enumerable[int] = (*Queue[int])(nil)

// Any returns true for the first element found where the predicate function returns true.
// It returns false if no element matches the predicate.
func (q *Queue[T]) Any(predicate functions.PredicateFunc[T]) bool {

	if q.lock != nil {
		q.lock.Lock()
		defer q.lock.Unlock()
	}

	iter := newForwardIterator[T](q, predicate)

	return iter.Start() != nil
}

// All applies the predicate function to every element in the collection,
// and returns true if all elements match the predicate.
func (q *Queue[T]) All(predicate functions.PredicateFunc[T]) bool {

	if q.lock != nil {
		q.lock.Lock()
		defer q.lock.Unlock()
	}

	iter := newForwardIterator[T](q, util.DefaultPredicate[T])

	for e := iter.Start(); e != nil; e = iter.Next() {
		if !predicate(e.Value()) {
			return false
		}
	}

	return true
}

// ForEach applies function f to all elements in the collection
func (q *Queue[T]) ForEach(f func(collections.Element[T])) {

	if q.lock != nil {
		q.lock.Lock()
		defer q.lock.Unlock()
	}

	iter := newForwardIterator[T](q, util.DefaultPredicate[T])

	for e := iter.Start(); e != nil; e = iter.Next() {
		f(e)
	}
}

// Map applies function f to all elements in the collection
// and returns a new Queue containing the result of f
func (q *Queue[T]) Map(f func(T) T) collections.Collection[T] {

	if q.lock != nil {
		q.lock.Lock()
		defer q.lock.Unlock()
	}

	iter := newForwardIterator[T](q, util.DefaultPredicate[T])

	q1 := New[T](WithComparer[T](q.compare))

	for e := iter.Start(); e != nil; e = iter.Next() {
		q1.enqueue(f(e.Value()))
	}

	return q1
}

// Select returns a new Queue containing only the items for which predicate is true
func (q *Queue[T]) Select(predicate functions.PredicateFunc[T]) collections.Collection[T] {

	if q.lock != nil {
		q.lock.RLock()
		defer q.lock.RUnlock()
	}

	return q.doSelect(predicate, false)
}

// SelectDeep returns a new Queue containing only the items for which predicate is true
//
// Elements are deep copied to the new collection using the provided [functions.DeepCopyFunc] if any.
func (q *Queue[T]) SelectDeep(predicate functions.PredicateFunc[T]) collections.Collection[T] {

	if q.lock != nil {
		q.lock.RLock()
		defer q.lock.RUnlock()
	}

	return q.doSelect(predicate, true)
}

// Find finds the first occurrence of an element matching the predicate.
//
// The function returns nil if no match.
func (q *Queue[T]) Find(predicate functions.PredicateFunc[T]) collections.Element[T] {

	if q.lock != nil {
		q.lock.RLock()
		defer q.lock.RUnlock()
	}

	result := q.doFind(predicate, false)

	if len(result) == 0 {
		return nil
	}

	return result[0]
}

// FindAll finds all occurrences of an element matching the predicate.
//
// The function returns an empty slice if none match.
func (q *Queue[T]) FindAll(predicate functions.PredicateFunc[T]) []collections.Element[T] {

	if q.lock != nil {
		q.lock.RLock()
		defer q.lock.RUnlock()
	}

	result := q.doFind(predicate, true)

	return result
}

// Min returns the minimum value in the collection according to the Comparer function
func (q *Queue[T]) Min() T {

	if q.lock != nil {
		q.lock.RLock()
		defer q.lock.RUnlock()
	}

	if q.size == 0 {
		panic(messages.COLLECTION_EMPTY)
	}

	l := len(q.buffer)

	if l <= 100 {
		// Do it directly
		m := q.buffer[q.head]

		for i := 0; i < q.size; i++ {
			ind := (q.head + i) % l
			v := q.buffer[ind]
			if q.compare(m, v) > 0 {
				m = v
			}
		}

		return m
	}

	// If buffer has not wrapped
	if q.head+q.size <= l {
		return util.Min(q.buffer[q.head:q.head+q.size], q.compare, q.concurrent)
	}

	// Else buffer has wrapped and tail is before head.
	// Aggregate each segment of the buffer concurrently.
	var m1, m2 T
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		m1 = util.Min(q.buffer[q.head:l], q.compare, q.concurrent)
		wg.Done()
	}()

	go func() {
		m2 = util.Min(q.buffer[0:(q.head+q.size)%l], q.compare, q.concurrent)
		wg.Done()
	}()

	wg.Wait()

	if q.compare(m1, m2) < 0 {
		return m1
	}

	return m2
}

// Max returns the maximum value in the collection according to the Comparer function
func (q *Queue[T]) Max() T {

	if q.lock != nil {
		q.lock.RLock()
		defer q.lock.RUnlock()
	}

	if q.size == 0 {
		panic(messages.COLLECTION_EMPTY)
	}

	l := len(q.buffer)

	if l <= 100 {
		// Do it directly
		m := q.buffer[q.head]

		for i := 0; i < q.size; i++ {
			ind := (q.head + i) % l
			v := q.buffer[ind]
			if q.compare(m, v) < 0 {
				m = v
			}
		}

		return m
	}

	// If buffer has not wrapped
	if q.head+q.size <= l {
		return util.Max(q.buffer[q.head:q.head+q.size], q.compare, q.concurrent)
	}

	// Else buffer has wrapped and tail is before head.
	// Aggregate each segment of the buffer concurrently.
	var m1, m2 T
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		m1 = util.Max(q.buffer[q.head:l], q.compare, q.concurrent)
		wg.Done()
	}()

	go func() {
		m2 = util.Max(q.buffer[0:(q.head+q.size)%l], q.compare, q.concurrent)
		wg.Done()
	}()

	wg.Wait()

	if q.compare(m1, m2) > 0 {
		return m1
	}

	return m2
}

func (q *Queue[T]) doFind(predicate functions.PredicateFunc[T], all bool) []collections.Element[T] {

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

func (q *Queue[T]) doSelect(predicate functions.PredicateFunc[T], deepCopy bool) collections.Collection[T] {

	q1 := New[T](WithComparer[T](q.compare), WithCapacity[T](util.Iif[int](q.initialCapacity > q.size, q.initialCapacity, q.size)))
	iter := newForwardIterator[T](q, predicate)

	for e := iter.Start(); e != nil; e = iter.Next() {
		if deepCopy {
			q1.enqueue(q.copy(e.Value()))
		} else {
			q1.enqueue(e.Value())
		}
	}

	return q1
}
