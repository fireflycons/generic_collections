/*
Package set defines the interface for set collections. Sets are collections of distinct values.
Sub-packages contain implementations.
*/
package sets

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/internal/local"
)

// Set is the abstract interface for collections of unique elements.
//
// Implemented by HashSet[T], OrderedSet[T].
type Set[T any] interface {
	// Set implements Collection
	collections.Collection[T]

	// Difference returns the difference between two sets.
	//
	// The new set consists of a shallow-copy of all elements that are in this set, but not other set.
	Difference(Set[T]) Set[T]

	// Intersection returns the intersection between two sets.
	//
	// The new set consists of a shallow-copy of all elements that are in both this set and the other.
	Intersection(Set[T]) Set[T]

	// Union returns the union of two sets.
	//
	// The new set consists of a shallow-copy of all elements that are in both this and the other set.
	Union(Set[T]) Set[T]

	// Tests whether the given value is contained within the set.
	//
	// Not indended to be used by client programs as it never takes a lock, even when thread safety
	// is enabled. Used to speed up the above set operations.
	UnlockedContains(T) bool

	// Prevent external implementations of this interface
	local.InternalInter
}
