// This package contains signatures for functions used by collections and their methods.
package functions

// ComparerFunc is the signature for a function that compares two values.
//
// This function allows the developer to supply a custom compare function to a collection for a type that is not one of the supported types.
// A compare function should return less than zero if a<b, zero if a==b or greater than zero if a>b and looks like this, e.g. for a collection typed on int
//
//	func myComparer(int a, int b) int {
//		return a-b
//	}
type ComparerFunc[T any] func(T, T) int

// PredicateFunc is the signature for the function used in filtering and searching collections.
//
// For many of the Enumerable methods, a predicate function must be given as an argument. A value is selected when the predicate function returns `true`. For instance, to filter all even numbers from a collection of int it might look like this
//
//	func filterEvens(int value) bool {
//		return value % 2 == 0
//	}
type PredicateFunc[T any] func(T) bool

// HashFunc is the signature for custom hash algorithms to use with HashSet.
//
// The hash algorithms for the supported types are exported as function variables by the hashset
// package so can be used to construct hashes for struct types.
type HashFunc[T any] func(T) uintptr

// Function signature for a function to deep copy a collection element.
//
// This function should return a new instance of type T copied from the original.
//
// The default action when making copies of collection elements is to
// value-copy the element. If the element type is a pointer, or a struct
// containing pointers this may not be what you want. Should you need to
// deep-copy elements, supply an implementation of this function to the
// collection's constructor.
type DeepCopyFunc[T any] func(T) T
