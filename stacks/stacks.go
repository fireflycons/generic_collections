/*
Package stacks defines the interface for LIFO stack collections. Sub-packages contain implementations.
*/
package stacks

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/internal/local"
)

// Stack is the abstract interface for collections that operate as LIFO queues.
//
// Implemented by Stack[T]
type Stack[T any] interface {
	// Stack implements Collection
	collections.Collection[T]

	// Stacks can be sorted
	collections.Sortable[T]

	// Push adds a value to the top of the stack.
	Push(value T)

	// Pop removes and returns the value at the top of the stack.
	//
	// Panics if the stack is empty.
	Pop() T

	// TryDequeue removes and returns the value at the front of the queue and true if
	// the queue is not empty; else zero value of T and false.
	TryPop() (T, bool)

	// Peek returns the value at the top of the stack without adjusting the stack.
	//
	// Panics if the queue is empty.
	Peek() T

	// TryPeek returns the value at the top of the stack and true if
	// the stack is not empty; else zero value of T and false.
	TryPeek() (T, bool)

	// Prevent external implementations of this interface
	local.InternalInter
}
