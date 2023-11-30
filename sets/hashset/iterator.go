package hashset

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"

	"github.com/fireflycons/generic_collections/internal/local"
	"golang.org/x/exp/maps"
)

// Assert interface implementation
var _ collections.Iterable[int] = (*HashSet[int])(nil)

// HashSetIterator implements an iterator over the elements in the set.
type HashSetIterator[T any] struct {
	util.IteratorBase[T]
	set            *HashSet[T]
	position       int
	bucketPosition int
	endPosition    int
	predicate      functions.PredicateFunc[T]
	keys           []uintptr

	local.InternalImpl
}

func newForwardIterator[T any](s *HashSet[T], predicate functions.PredicateFunc[T]) collections.Iterator[T] {
	keys := maps.Keys(s.buffer)

	return &HashSetIterator[T]{
		set:            s,
		position:       0,
		bucketPosition: 0,
		endPosition:    len(keys) - 1,
		keys:           keys,
		predicate:      predicate,
		IteratorBase: util.IteratorBase[T]{
			Version:    s.version,
			NilElement: nil,
		},
	}
}

// Iterator returns a forward iterator that walks the set from first to last element
//
//	iter := set.Iterator()
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (s *HashSet[T]) Iterator() collections.Iterator[T] {

	return newForwardIterator(s, util.DefaultPredicate[T])
}

// TakeWhile returns a forward iterater that walks the collection returning only
// those elements for which predicate returns true.
//
//	set := hashset.New[int]()
//	// add values
//	iter := set.TakeWhile(func (val int) bool { return val % 2 == 0 })
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (s *HashSet[T]) TakeWhile(predicate functions.PredicateFunc[T]) collections.Iterator[T] {

	return newForwardIterator(s, predicate)
}

// Start begins iteration across the set returning the fisrt element,
// which will be nil if the set is empty.
//
// Panics if the underlying set is modified between iteration creation and call to Start()
func (i *HashSetIterator[T]) Start() collections.Element[T] {
	i.validateIterator()
	i.position = 0
	if i.set.size == 0 || !moveToNextPopulatedBucket(i) {
		return i.NilElement
	}

	valPtr := &i.set.buffer[i.keys[i.position]][i.bucketPosition]

	if !i.predicate(*valPtr) {
		return i.Next()
	}

	elem := util.NewElementType[T](i.set, valPtr)
	i.bucketPosition++
	return elem
}

// Next returns the next element from the iterator,
// which will be nil if the end has been reached.
//
// Panics if the underlying set is modified between calls to Next.
func (i *HashSetIterator[T]) Next() collections.Element[T] {
	i.validateIterator()

	for {
		e := moveForward(i)

		if e == i.NilElement || i.predicate(e.Value()) {
			return e
		}
	}
}

func (i *HashSetIterator[T]) validateIterator() {
	// util.ValidatePointerNotNil(unsafe.Pointer(i))
	if i.Version != i.set.version {
		panic(messages.COLLECTION_MODIFIED)
	}
}

func moveForward[T any](i *HashSetIterator[T]) collections.Element[T] {

	if i.bucketPosition < len(i.set.buffer[i.keys[i.position]]) {
		val := util.NewElementType[T](i.set, &i.set.buffer[i.keys[i.position]][i.bucketPosition])
		i.bucketPosition++
		return val
	}

	i.position++

	if !moveToNextPopulatedBucket(i) {
		return i.NilElement
	}

	if i.position >= len(i.keys) {
		return i.NilElement
	}

	val := util.NewElementType[T](i.set, &i.set.buffer[i.keys[i.position]][0])
	i.bucketPosition = 1
	return val
}

func moveToNextPopulatedBucket[T any](i *HashSetIterator[T]) bool {
	for j := i.position; j <= i.endPosition; j++ {
		if len(i.set.buffer[i.keys[j]]) > 0 {
			i.position = j
			i.bucketPosition = 0
			return true
		}
	}

	return false
}

func max(x, y int) int {
	if x > y {
		return x
	}

	return y
}
