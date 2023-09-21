/*
Package collections provides the interfaces that define the collections in this module.
*/
package collections

import (
	"fmt"

	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/local"
)

type CollectionType int

// CollectionType is a enueration used to identify the type of a collection
// implementated by an instance of the abstract interface [collections.Collection].
const (
	_ CollectionType = iota
	COLLECTION_STACK
	COLLECTION_DLIST
	COLLECTION_SLIST
	COLLECTION_QUEUE
	COLLECTION_RINGBUFFER
	COLLECTION_HASHSET
	COLLECTION_ORDEREDSET
)

// Collection is the abstract interface to all collection types defined in this package.
// Defines methods implemented by all collections.
type Collection[T any] interface {

	// All collections are enumerable.
	Enumerable[T]

	// All collections are iterable.
	Iterable[T]

	// All collections have a string representation.
	fmt.Stringer

	// Add adds a value to the collection.
	//
	// How the value is added is specific to the type of this collection.
	// For linear collections the value is generally added at the end (stacks are pushed).
	// For non-linear collections like sets, the value will be inserted
	// according to the rules of the backing data structure.
	//
	// Returns true if a new element was added to the collection.
	Add(T) bool

	// AddRange adds a slice of elements to the collection.
	//
	// Insertion rules are as per calling Add for each value in the slice.
	AddRange([]T)

	// AddCollection inserts the values of the given collection into this collection.
	//
	// How the new values are inserted is specific to the type of this collection,
	// and how the values for the other collection are obtained is specific to the retrieval order
	// for that collection.
	AddCollection(Collection[T])

	// Clear removes all elements from the collection
	Clear()

	// Contains returns true if the given value is present in the collection; else false.
	Contains(T) bool

	// Count returns the number of values stored in the collection.
	Count() int

	// IsEmpty returns true if the collection has no elements
	IsEmpty() bool

	// Remove removes the first occurrence of the given value from the collection.
	//
	// Returns true if the value was removed.
	Remove(T) bool

	// ToSlice returns the content of the collection as a slice.
	ToSlice() []T

	// ToSlice returns the content of the collection as a slice.
	//
	// If a DeepCopyFunc[T] was provided to the collection constructor it will be used,
	// else a by-value copy is made, i.e. works the same as ToSlice.
	ToSliceDeep() []T

	// Type returns the type of the collection (to avoid unnecessary reflecting).
	Type() CollectionType

	// Prevent external implementations of this interface
	local.InternalInter
}

// Enumerable describes methods that enumerate across a collection.
type Enumerable[T any] interface {

	// Any returns true for the first element found where the predicate function returns true.
	// It returns false if no element matches the predicate.
	Any(predicate functions.PredicateFunc[T]) bool

	// All applies the predicate function to every element in the collection,
	// and returns true if all elements match the predicate.
	All(predicate functions.PredicateFunc[T]) bool

	// Find finds the first occurrence of an element matching the predicate.
	//
	// The function returns nil if no match.
	Find(predicate functions.PredicateFunc[T]) Element[T]

	// FindAll finds all occurrences of an element matching the predicate.
	//
	// The function returns an empty slice if none match.
	FindAll(predicate functions.PredicateFunc[T]) []Element[T]

	// ForEach applies function f to all elements in the collection
	//
	// You can modify the value in the collection via the Element interface,
	// except for implementations of Set which will panic as modification of
	// values in a set breaks the set structure.
	ForEach(func(Element[T]))

	// Min returns the minimum value in the collection according to the Comparer function.
	Min() T

	// Max returns the maximum value in the collection according to the Comparer function.
	Max() T

	// Map applies function f to all elements in the collection
	// and returns a new collection containing the results of f
	// applied to each value in the source collection.
	Map(func(T) T) Collection[T]

	// Select returns a new collection of the same type
	// containing only the items for which predicate is true.
	Select(functions.PredicateFunc[T]) Collection[T]

	// SelectDeep returns a new collection of the same type
	// containing only the items for which predicate is true.
	//
	// If a DeepCopyFunc[T] was provided to the collection constructor it will be used,
	// else a by-value copy is made, i.e. works the same as Select.
	SelectDeep(functions.PredicateFunc[T]) Collection[T]

	// Prevent external implementations of this interface
	local.InternalInter
}

// Iterable defines collections that can be iterated from start to end.
type Iterable[T any] interface {
	// Iterator returns an iterator that walks the collection from start to end.
	Iterator() Iterator[T]

	// TakeWhile returns a forward iterater that walks the collection returning only
	// those elements for which predicate returns true.
	TakeWhile(functions.PredicateFunc[T]) Iterator[T]

	// Prevent external implementations of this interface
	local.InternalInter
}

// ReverseIterable defines collections that can be iterated from end to start.
type ReverseIterable[T any] interface {
	// ReverseIterator returns an iterator that walks the collection from end to start.
	ReverseIterator() Iterator[T]
}

// Sortable defines collections that can have their values sorted.
//
// Built-in implementation from sort package is used, which is a variation of introspective sort
// using pattern-defeating quicksort instead of regular quicksort. This was found to be the
// fastest algorithm during testing. Worst case should be O(n log n).
//
// Sortable is implemented by all collections except implementations of [sets.Set], as sets by their
// nature are either always ordered or always unordered, depending on the strategy used.
//
// Note that List types use a merge sort moving node pointers around, as this is faster
// than converting to a slice and using pdqsort, then back to a list. Still worst O(n log n).
type Sortable[T any] interface {

	// All Sortables are collections.
	Collection[T]

	// Sort performs an in-place, non-stable sort of the collection in ascending order.
	Sort()

	// Sorted returns a sorted copy of this collection as a new collection of the same type.
	Sorted() Collection[T]

	// Sort performs an in-place, non-stable sort of the collection in descending order.
	SortDescending()

	// Sorted returns a descending order sorted copy of this collection as a new collection of the same type.
	SortedDescending() Collection[T]

	// Prevent external implementations of this interface
	local.InternalInter
}

// Iterator describes the mechanism of iterating across a collection.
type Iterator[T any] interface {
	// Begins iteration at the first item in the collection and returns that elemment.
	// First item is determined by the iterator direction.
	Start() Element[T]

	// Moves the iteration to the next element in the collection and returns it.
	Next() Element[T]

	// Prevent external implementations of this interface
	local.InternalInter
}

// Element is an interface describing an element of a collection at a given position.
type Element[T any] interface {
	// Gets the value of this element.
	Value() T

	// Gets a pointer to the actual value in the collection,
	// allowing modification of the stored value.
	//
	// Note that this method will panic if the element represents a value in any implementation of Set[T],
	// since modifying the value will break the set implementation.
	ValuePtr() *T

	// Prevent external implementations of this interface
	local.InternalInter
}
