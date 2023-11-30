/*
Package lists defines the interface for list collections. Sub-packages contain implementations.
*/
package lists

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/internal/local"
)

// Type List defines the interface to list collections.
type List[T any] interface {
	// List implements Collection
	collections.Collection[T]

	// Lists can be sorted
	collections.Sortable[T]

	// Prevent external implementations of this interface
	local.InternalInter

	// AddItemFirst adds the given value at the head of the list.
	AddItemFirst(value T)

	// AddItemLast adds the given value at the end of the list.
	AddItemLast(value T)

	// RemoveItem searches the list for the first occurrence of value
	// and removes the node containing that value. Up to O(n).
	//
	// Returns true if a node was removed; else false.
	RemoveItem(value T) bool

	// RemoveFirst removes the node at the head of the list and returns the value that was stored
	//
	// Panics if list is empty.
	RemoveFirst() T

	// RemoveLast removes the node at the end of the list and returns the value that was stored.
	//
	// Panics if list is empty.
	RemoveLast() T

	// TryRemoveFirst removes the node at the head of the list and returns the value that was stored and true,
	// or the zero value of T and false if the list is empty.
	TryRemoveFirst() (T, bool)

	// TryRemoveLast removes the node at the end of the list and returns the value that was stored and true,
	// or the zero value of T and false if the list is empty.
	TryRemoveLast() (T, bool)
}
