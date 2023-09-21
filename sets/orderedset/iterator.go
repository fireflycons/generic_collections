package orderedset

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/local"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/fireflycons/generic_collections/stacks/stack"
)

type direction bool

// Assert interface implementation
var _ collections.Iterable[int] = (*OrderedSet[int])(nil)

const (
	forward, reverse direction = true, false
)

type OrderedSetIterator[T any] struct {
	util.IteratorBase[T]
	set       *OrderedSet[T]
	stack     *stack.Stack[*node[T]]
	direction direction
	predicate functions.PredicateFunc[T]
	local.InternalImpl
}

func newForwardIterator[T any](set *OrderedSet[T], predicate functions.PredicateFunc[T]) *OrderedSetIterator[T] {
	return &OrderedSetIterator[T]{
		set:       set,
		stack:     stack.New(stack.WithCapacity[*node[T]](2 * intlog2(set.size+1))),
		direction: forward,
		predicate: predicate,
		IteratorBase: util.IteratorBase[T]{
			Version:    set.version,
			NilElement: nil,
		},
	}
}

func newReverseIterator[T any](set *OrderedSet[T]) *OrderedSetIterator[T] {
	iter := newForwardIterator(set, util.DefaultPredicate[T])
	iter.direction = reverse
	return iter
}

// Iterator returns an iterator that walks the collection in ascending order of values.
func (s *OrderedSet[T]) Iterator() collections.Iterator[T] {

	return newForwardIterator(s, util.DefaultPredicate[T])
}

// ReverseIterator returns an iterator that walks the collection in descending order of values.
func (s *OrderedSet[T]) ReverseIterator() collections.Iterator[T] {

	return newReverseIterator(s)
}

// TakeWhile returns a forward iterater that walks the collection returning only
// those elements for which predicate returns true.
//
//	set := hashset.New[int]()
//	// add values
//	iter := set.TakeWhile(func (val int) bool { return val % 2 == 0 })
//
//	for e := iter.Start() ; e != nil; e = iter.Next() {
//		// do something with e.Value()
//	}
func (s *OrderedSet[T]) TakeWhile(predicate functions.PredicateFunc[T]) collections.Iterator[T] {

	return newForwardIterator(s, predicate)
}

// Start begins an iteration across the set returning the fisrt element,
// which will be nil if the collection is empty.
//
// Panics if the set has been modified since creation of the iterator.
func (i *OrderedSetIterator[T]) Start() collections.Element[T] {
	i.validateIterator()
	i.move(i.set.root)

	if i.stack.Count() == 0 {
		return i.NilElement
	}

	return i.Next()
}

// Next returns the next element in the set,
// which will be nil if the end has been reached.
//
// Panics if the set has been modified since creation of the iterator.
func (i *OrderedSetIterator[T]) Next() collections.Element[T] {
	i.validateIterator()

	for {
		if i.stack.Count() == 0 {
			return i.NilElement
		}

		current := i.stack.Pop()
		i.move(util.Iif(i.direction == reverse, current.left, current.right))

		if i.predicate(current.item) {
			return util.NewElementType[T](i.set, &current.item)
		}
	}
}

func (i *OrderedSetIterator[T]) move(n *node[T]) {
	var next *node[T]

	for n != nil {
		next = util.Iif(i.direction == reverse, n.right, n.left)
		i.stack.Push(n)
		n = next
	}
}

func (i *OrderedSetIterator[T]) validateIterator() {
	// util.ValidatePointerNotNil(unsafe.Pointer(i))
	if i.Version != i.set.version {
		panic(messages.COLLECTION_MODIFIED)
	}
}
