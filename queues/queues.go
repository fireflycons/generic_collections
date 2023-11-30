/*
Package lists defines the interface for FIFO queue collections. Sub-packages contain implementations.
*/
package queues

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/internal/local"
)

// Queue is the abstract interface for collections that operate as FIFO queues.
//
// Implemented by Queue[T], RingBuffer[T].
type Queue[T any] interface {
	// Queue implements Collection
	collections.Collection[T]

	// Queues can be sorted
	collections.Sortable[T]

	// Dequeue removes the value at the front of the queue and returns it.
	//
	// Panics if the queue is empty.
	Dequeue() T

	// TryDequeue removes and returns the value at the front of the queue and true if
	// the queue is not empty; else zero value of T and false.
	TryDequeue() (T, bool)

	// Enqueue adds a value to the back of the queue.
	Enqueue(value T)

	// Peek returns the value at the front of the queue without removing it.
	//
	// Panics if the queue is empty.
	Peek() T

	// TryPeek returns the value at the front of the queue and true if
	// the queue is not empty; else zero value of T and false.
	TryPeek() (T, bool)

	// Prevent external implementations of this interface
	local.InternalInter
}
