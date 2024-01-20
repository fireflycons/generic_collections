/*
Package orderedset provides a red-black tree backed ordered collection of unique items.
*/
package orderedset

import (
	"fmt"
	"sync"

	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/local"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/fireflycons/generic_collections/sets"
	"github.com/fireflycons/generic_collections/stacks/stack"
)

// Assert OrderedSet implements required interfaces.
var _ sets.Set[int] = (*OrderedSet[int])(nil)
var _ collections.ReverseIterable[int] = (*OrderedSet[int])(nil)

type color bool

const (
	black, red color = true, false
)

type treeWalkPredicate[T any] func(*node[T]) bool

// OrderedSetOptionFunc is the signature of a function
// for providing options to the OrderedSet constructor.
type OrderedSetOptionFunc[T any] func(*OrderedSet[T])

// OrderedSet stores an ordered collection of unique elements.
type OrderedSet[T any] struct {
	version    int
	lock       *sync.RWMutex
	root       *node[T]
	size       int
	compare    functions.ComparerFunc[T]
	copy       functions.DeepCopyFunc[T]
	concurrent bool
	local.InternalImpl
}

// node is a single element within the tree.
type node[T any] struct {
	item   T
	color  color
	left   *node[T]
	right  *node[T]
	Parent *node[T]
}

func newNode[T any](value T) *node[T] {
	return &node[T]{
		item:  value,
		color: red,
	}
}

// Constructs a new OrderedSet[T].
func New[T any](options ...OrderedSetOptionFunc[T]) *OrderedSet[T] {
	set := &OrderedSet[T]{}

	for _, o := range options {
		o(set)
	}

	if set.copy == nil {
		set.copy = util.DefaultDeepCopy[T]
	}

	if set.compare == nil {
		set.compare = util.GetDefaultComparer[T]()
	}

	return set
}

// Option function for New to make the collection thread-safe. Adds overhead.
func WithThreadSafe[T any]() OrderedSetOptionFunc[T] {
	return func(s *OrderedSet[T]) {
		s.lock = &sync.RWMutex{}
	}
}

// Option function to enable concurrency feature.
func WithConcurrent[T any]() OrderedSetOptionFunc[T] {
	return func(s *OrderedSet[T]) {
		s.concurrent = true
	}
}

// Option function for NewOrderedSet to provide a comparer function for values of type T.
// Required if the element type is not numeric, bool, pointer or string.
func WithComparer[T any](comparer functions.ComparerFunc[T]) OrderedSetOptionFunc[T] {
	if comparer == nil {
		panic(messages.COMP_FN_NIL)
	}
	return func(s *OrderedSet[T]) {
		s.compare = comparer
	}
}

// Option func to provide a deep copy implementation for collection elements.
func WithDeepCopy[T any](copier functions.DeepCopyFunc[T]) OrderedSetOptionFunc[T] {
	// Can be nil
	return func(s *OrderedSet[T]) {
		s.copy = copier
	}
}

// AddRange adds a slice of values to the set.
func (s *OrderedSet[T]) AddRange(values []T) {

	if len(values) == 0 {
		return
	}

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	s.version++

	for _, v := range values {
		s.doInsert(v)
	}
}

// AddCollection inserts the values of the given collection into this set.
func (s *OrderedSet[T]) AddCollection(collection collections.Collection[T]) {

	s.AddRange(collection.ToSliceDeep())
}

// Add adds a value into the collection.
// Returns false if the value already exists; else true if it was added.
func (s *OrderedSet[T]) Add(value T) bool {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	inserted := s.doInsert(value)
	s.version++
	return inserted
}

// Contains returns true if the given value exists in the set.
func (s *OrderedSet[T]) Contains(value T) bool {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	return s.lookup(value) != nil
}

func (s *OrderedSet[T]) UnlockedContains(value T) bool {
	return s.lookup(value) != nil
}

// Get returns the collection element that matches the given value, or nil if it is not found.
// Useful if the set contains struct elements you want to modify in-place.
func (s *OrderedSet[T]) Get(value T) collections.Element[T] {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	n := s.lookup(value)

	if n == nil {
		return nil
	}

	return util.NewElementType[T](s, &n.item)
}

// Remove removes a value from the set.
//
// Returns true if the value was present and was removed;
// else false.
func (s *OrderedSet[T]) Remove(key T) bool {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	s.version++
	var child *node[T]
	n := s.lookup(key)
	if n == nil {
		return false
	}
	if n.left != nil && n.right != nil {
		pred := n.left.maximumNode()
		n.item = pred.item
		n = pred
	}
	if n.left == nil || n.right == nil {
		if n.right == nil {
			child = n.left
		} else {
			child = n.right
		}
		if n.color == black {
			n.color = nodeColor(child)
			// Delete as per https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
			s.deleteCase1(n)
		}
		s.replaceNode(n, child)
		if n.Parent == nil && child != nil {
			child.color = black
		}
	}
	s.size--
	return true
}

// Empty returns true if tree does not contain any nodes.
func (s *OrderedSet[T]) Empty() bool {

	return s.size == 0
}

// Count returns the number of values stored in the collection.
func (s *OrderedSet[T]) Count() int {

	return s.size
}

// IsEmpty returns true if the collection has no elements.
func (s *OrderedSet[T]) IsEmpty() bool {
	return s.size == 0
}

// ToSlice returns the collection content as a slice.
// The values will be in ascending order.
func (s *OrderedSet[T]) ToSlice() []T {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	slc := make([]T, s.size)
	s.copyTo(slc, 0, s.size, false)
	return slc
}

// ToSliceDeep returns the collection content as a slice.
// The values will be in ascending order.
// Elements are deep copied using the provided [functions.DeepCopyFunc] if any.
func (s *OrderedSet[T]) ToSliceDeep() []T {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	slc := make([]T, s.size)
	s.copyTo(slc, 0, s.size, true)
	return slc
}

// Clear removes all nodes from the tree.
func (s *OrderedSet[T]) Clear() {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	s.root = nil
	s.size = 0
	s.version++
}

// String returns a string representation of container.
func (s *OrderedSet[T]) String() string {
	str := "OrderedSet\n"
	if !s.Empty() {
		output(s.root, "", true, &str)
	}
	return str
}

func (n *node[T]) String() string {
	return fmt.Sprintf("%v", n.item)
}

// Type returns the type of the collection (to avoid reflecting).
func (s *OrderedSet[T]) Type() collections.CollectionType {
	return collections.COLLECTION_ORDEREDSET
}

type containsFnT[T any] func(T) bool

// Difference returns the difference between two sets.
// The new set consists of all elements that are in this set, but not other set.
//
// The argument can be any implementation of Set[T]. The result is a new OrderedSet with the same properties as this one.
// Items are shallow-copied.
func (s *OrderedSet[T]) Difference(other sets.Set[T]) sets.Set[T] {

	ol := util.GetLock[T](other)

	if ol != nil {
		ol.RLock()
		defer ol.RUnlock()
	}

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	result := s.makeEmptyCopy()

	osOther, otherIsOrderedSet := other.(*OrderedSet[T])

	var otherContains containsFnT[T]

	if otherIsOrderedSet {
		otherContains = func(val T) bool {
			return osOther.lookup(val) != nil
		}
	} else {
		otherContains = func(val T) bool {
			return other.UnlockedContains(val)
		}
	}

	s.inOrderTreeWalk(func(n *node[T]) bool {
		if !otherContains(n.item) {
			result.doInsert(n.item)
		}
		return true
	})

	return result
}

// Intersection returns the intersection between two sets.
// The new set consists of all elements that are in both this set and the other.
//
// The argument can be any implementation of Set[T]. The result is a new OrderedSet with the same properties as this one.
// Items are shallow-copied.
func (s *OrderedSet[T]) Intersection(other sets.Set[T]) sets.Set[T] {

	ol := util.GetLock[T](other)

	if ol != nil {
		ol.RLock()
		defer ol.RUnlock()
	}

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	// It's much quicker to scan the smaller collection
	// and look up values in the larger one as lookup
	// is very fast in sets.
	var smaller, larger sets.Set[T]

	if s.Count() < other.Count() {
		smaller = s
		larger = other
	} else {
		smaller = other
		larger = s
	}

	result := s.makeEmptyCopy()

	// Where we know the set to walk is an OrderedSet
	// a treewalk is faster than a conversion to slice first.
	osSml, smallerIsOrderedSet := smaller.(*OrderedSet[T])
	osLrg, largerIsOrderedSet := larger.(*OrderedSet[T])

	if smallerIsOrderedSet {
		if largerIsOrderedSet {
			osSml.inOrderTreeWalk(func(n *node[T]) bool {
				if osLrg.lookup(n.item) != nil {
					result.doInsert(n.item)
				}
				return true
			})

			return result
		}

		osSml.inOrderTreeWalk(func(n *node[T]) bool {
			if larger.UnlockedContains(n.item) {
				result.doInsert(n.item)
			}
			return true
		})

		return result
	}

	// Smaller is not OrderedSet, therefore larger must be OrderedSet
	for _, value := range smaller.ToSlice() {
		if osLrg.lookup(value) != nil {
			result.doInsert(value)
		}
	}

	return result
}

// Union returns the union of two sets.
// The new set consists of all elements that are in buth this and the other set.
//
// The argument can be any implementation of Set[T]. The result is a new OrderedSet with the same properties as this one.
// Items are shallow-copied.
func (s *OrderedSet[T]) Union(other sets.Set[T]) sets.Set[T] {

	ol := util.GetLock[T](other)

	if ol != nil {
		ol.RLock()
		defer ol.RUnlock()
	}

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	result := s.makeEmptyCopy()
	result.AddCollection(s)
	result.AddCollection(other)
	return result
}

func (s *OrderedSet[T]) copyTo(slc []T, index, count int, deepCopy bool) {
	if slc == nil {
		panic(fmt.Sprintf(messages.ARG_NIL_FMT, "array"))
	}

	if index < 0 {
		panic(fmt.Sprintf(messages.ARG_OUT_OF_RANGE_FMT, "index"))
	}

	if count < 0 {
		panic(fmt.Sprintf(messages.ARG_OUT_OF_RANGE_FMT, "count"))
	}

	if index > len(slc) || count > len(slc)-index {
		panic(messages.SLICE_TOO_SMALL)
	}

	count += index

	s.inOrderTreeWalk(func(n *node[T]) bool {
		if index > count {
			return false
		} else {
			if deepCopy {
				slc[index] = util.DeepCopy(n.item, s.copy)
			} else {
				slc[index] = n.item
			}
			index++
			return true
		}
	})
}

// TreeWalk walks the underlying tree implemnetation of this set
// from smallest to largest value calling the delegate for each value.
// If the action delegate returns false, stop the walk.
//
// Returns true if the entire tree has been walked.
// Otherwise returns false.
func (s *OrderedSet[T]) TreeWalk(action func(T) bool) bool {
	return s.inOrderTreeWalkWithDirection(func(n *node[T]) bool {
		return action(n.item)
	}, false)
}

// Do an in order walk on tree and calls the delegate for each node.
// If the action delegate returns false, stop the walk.
//
// Return true if the entire tree has been walked.
// Otherwise returns false.
func (s *OrderedSet[T]) inOrderTreeWalk(action treeWalkPredicate[T]) bool {
	return s.inOrderTreeWalkWithDirection(action, false)
}

func (s *OrderedSet[T]) inOrderTreeWalkWithDirection(action treeWalkPredicate[T], reverse bool) bool {
	if s.root == nil {
		return true
	}

	// The maximum height of a red-black tree is 2*lg(n+1).
	// See page 264 of "Introduction to algorithms" by Thomas H. Cormen
	// note: this should be logbase2, but since the stack grows itself, we
	// don't want the extra cost
	stack := stack.New(stack.WithCapacity[*node[T]](2 * intlog2(s.size+1)))
	current := s.root

	for current != nil {
		stack.Push(current)
		current = util.Iif(reverse, current.right, current.left)
	}

	for stack.Count() != 0 {
		current = stack.Pop()
		if !action(current) {
			return false
		}

		n := util.Iif(reverse, current.left, current.right)

		for n != nil {
			stack.Push(n)
			n = util.Iif(reverse, n.right, n.left)
		}
	}

	return true
}

func intlog2(value int) int {
	c := 0
	for value > 0 {
		c++
		value >>= 1
	}
	return c
}

func output[T any](n *node[T], prefix string, isTail bool, str *string) {
	if n.right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(n.right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += n.String() + "\n"
	if n.left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(n.left, newPrefix, true, str)
	}
}

func (s *OrderedSet[T]) lookup(key T) *node[T] {
	node := s.root
	for node != nil {
		compare := s.compare(key, node.item)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.left
		case compare > 0:
			node = node.right
		}
	}
	return nil
}

func (n *node[T]) grandparent() *node[T] {
	if n != nil && n.Parent != nil {
		return n.Parent.Parent
	}
	return nil
}

func (n *node[T]) uncle() *node[T] {
	if n == nil || n.Parent == nil || n.Parent.Parent == nil {
		return nil
	}
	return n.Parent.sibling()
}

func (n *node[T]) sibling() *node[T] {
	if n == nil || n.Parent == nil {
		return nil
	}
	if n == n.Parent.left {
		return n.Parent.right
	}
	return n.Parent.left
}

func (s *OrderedSet[T]) rotateLeft(n *node[T]) {
	rightNode := n.right
	s.replaceNode(n, rightNode)
	n.right = rightNode.left
	if rightNode.left != nil {
		rightNode.left.Parent = n
	}
	rightNode.left = n
	n.Parent = rightNode
}

func (s *OrderedSet[T]) rotateRight(n *node[T]) {
	leftNode := n.left
	s.replaceNode(n, leftNode)
	n.left = leftNode.right
	if leftNode.right != nil {
		leftNode.right.Parent = n
	}
	leftNode.right = n
	n.Parent = leftNode
}

func (s *OrderedSet[T]) replaceNode(oldNode *node[T], newNode *node[T]) {
	if oldNode.Parent == nil {
		s.root = newNode
	} else {
		if oldNode == oldNode.Parent.left {
			oldNode.Parent.left = newNode
		} else {
			oldNode.Parent.right = newNode
		}
	}
	if newNode != nil {
		newNode.Parent = oldNode.Parent
	}
}

func (s *OrderedSet[T]) doInsert(value T) bool {
	var insertedNode *node[T]
	if s.root == nil {
		s.root = newNode(value)
		insertedNode = s.root
	} else {
		n := s.root
		loop := true
		for loop {
			order := s.compare(value, n.item)
			switch {
			case order == 0:
				return false
			case order < 0:
				if n.left == nil {
					n.left = newNode(value)
					insertedNode = n.left
					loop = false
				} else {
					n = n.left
				}
			case order > 0:
				if n.right == nil {
					n.right = newNode(value)
					insertedNode = n.right
					loop = false
				} else {
					n = n.right
				}
			}
		}
		insertedNode.Parent = n
	}
	// Insertion as per https://en.wikipedia.org/wiki/Red%E2%80%93black_tree
	s.insertCase1(insertedNode)
	s.size++
	return true
}

func (s *OrderedSet[T]) insertCase1(n *node[T]) {
	if n.Parent == nil {
		n.color = black
	} else {
		s.insertCase2(n)
	}
}

func (s *OrderedSet[T]) insertCase2(n *node[T]) {
	if nodeColor(n.Parent) == black {
		return
	}
	s.insertCase3(n)
}

func (s *OrderedSet[T]) insertCase3(n *node[T]) {
	uncle := n.uncle()
	if nodeColor(uncle) == red {
		n.Parent.color = black
		uncle.color = black
		n.grandparent().color = red
		s.insertCase1(n.grandparent())
	} else {
		s.insertCase4(n)
	}
}

func (s *OrderedSet[T]) insertCase4(n *node[T]) {
	grandparent := n.grandparent()
	if n == n.Parent.right && n.Parent == grandparent.left {
		s.rotateLeft(n.Parent)
		n = n.left
	} else if n == n.Parent.left && n.Parent == grandparent.right {
		s.rotateRight(n.Parent)
		n = n.right
	}
	s.insertCase5(n)
}

func (s *OrderedSet[T]) insertCase5(n *node[T]) {
	n.Parent.color = black
	grandparent := n.grandparent()
	grandparent.color = red
	if n == n.Parent.left && n.Parent == grandparent.left {
		s.rotateRight(grandparent)
	} else if n == n.Parent.right && n.Parent == grandparent.right {
		s.rotateLeft(grandparent)
	}
}

func (n *node[T]) maximumNode() *node[T] {
	if n == nil {
		return nil
	}
	for n.right != nil {
		n = n.right
	}
	return n
}

func (s *OrderedSet[T]) deleteCase1(n *node[T]) {
	if n.Parent == nil {
		return
	}
	s.deleteCase2(n)
}

func (s *OrderedSet[T]) deleteCase2(n *node[T]) {
	sibling := n.sibling()
	if nodeColor(sibling) == red {
		n.Parent.color = red
		sibling.color = black
		if n == n.Parent.left {
			s.rotateLeft(n.Parent)
		} else {
			s.rotateRight(n.Parent)
		}
	}
	s.deleteCase3(n)
}

func (s *OrderedSet[T]) deleteCase3(n *node[T]) {
	sibling := n.sibling()
	if nodeColor(n.Parent) == black &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.left) == black &&
		nodeColor(sibling.right) == black {
		sibling.color = red
		s.deleteCase1(n.Parent)
	} else {
		s.deleteCase4(n)
	}
}

func (s *OrderedSet[T]) deleteCase4(n *node[T]) {
	sibling := n.sibling()
	if nodeColor(n.Parent) == red &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.left) == black &&
		nodeColor(sibling.right) == black {
		sibling.color = red
		n.Parent.color = black
	} else {
		s.deleteCase5(n)
	}
}

func (s *OrderedSet[T]) deleteCase5(n *node[T]) {
	sibling := n.sibling()
	if n == n.Parent.left &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.left) == red &&
		nodeColor(sibling.right) == black {
		sibling.color = red
		sibling.left.color = black
		s.rotateRight(sibling)
	} else if n == n.Parent.right &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.right) == red &&
		nodeColor(sibling.left) == black {
		sibling.color = red
		sibling.right.color = black
		s.rotateLeft(sibling)
	}
	s.deleteCase6(n)
}

func (s *OrderedSet[T]) deleteCase6(n *node[T]) {
	sibling := n.sibling()
	sibling.color = nodeColor(n.Parent)
	n.Parent.color = black
	if n == n.Parent.left && nodeColor(sibling.right) == red {
		sibling.right.color = black
		s.rotateLeft(n.Parent)
	} else if nodeColor(sibling.left) == red {
		sibling.left.color = black
		s.rotateRight(n.Parent)
	}
}

func (s *OrderedSet[T]) makeEmptyCopy() *OrderedSet[T] {
	other := &OrderedSet[T]{
		compare:    s.compare,
		copy:       s.copy,
		concurrent: s.concurrent,
	}

	if s.lock != nil {
		other.lock = &sync.RWMutex{}
	}

	return other
}

func nodeColor[T any](n *node[T]) color {
	if n == nil {
		return black
	}
	return n.color
}
