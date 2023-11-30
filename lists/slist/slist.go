/*
Package slist implements a singly linked list.
*/
package slist

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/local"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/fireflycons/generic_collections/lists"
)

// Assert SList implements required interfaces.
var _ lists.List[int] = (*SList[int])(nil)

// SListOptionFunc is the signature of a function
// for providing options to the SList constructor.
type SListOptionFunc[T any] func(*SList[T])

type SList[T any] struct {
	version int
	lock    *sync.RWMutex
	head    *SListNode[T]
	tail    *SListNode[T]
	count   int
	compare functions.ComparerFunc[T]
	copy    functions.DeepCopyFunc[T]
	local.InternalImpl
}

// New constructs a linked list.
func New[T any](options ...SListOptionFunc[T]) *SList[T] {
	sl := &SList[T]{}

	for _, o := range options {
		o(sl)
	}

	if sl.copy == nil {
		sl.copy = util.DefaultDeepCopy[T]
	}

	if sl.compare == nil {
		sl.compare = util.GetDefaultComparer[T]()
	}

	return sl
}

// Option function to make the collection thread-safe. Adds overhead.
func WithThreadSafe[T any]() SListOptionFunc[T] {
	return func(sl *SList[T]) {
		sl.lock = &sync.RWMutex{}
	}
}

// Option functionto provide a comparer function for values of type T.
// Required if the element type is not numeric, bool, pointer or string.
func WithComparer[T any](comparer functions.ComparerFunc[T]) SListOptionFunc[T] {
	if comparer == nil {
		panic(messages.COMP_FN_NIL)
	}
	return func(sl *SList[T]) {
		sl.compare = comparer
	}
}

// Option function to provide a deep copy implementation for collection elements.
func WithDeepCopy[T any](copier functions.DeepCopyFunc[T]) SListOptionFunc[T] {
	// Can be nil
	return func(sl *SList[T]) {
		sl.copy = copier
	}
}

// AddItemFirst adds the given value at the head of the list.
func (l *SList[T]) AddItemFirst(value T) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	newNode := &SListNode[T]{
		list: l,
		item: value,
	}

	l.prependNode(newNode)
	l.version++
}

// AddItemLast adds the given value at the end of the list.
func (l *SList[T]) AddItemLast(value T) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	newNode := &SListNode[T]{
		list: l,
		item: value,
	}

	l.appendNode(newNode)
	l.version++
}

// Add adds a value to the end of the list.
//
// Always returns true.
func (l *SList[T]) Add(value T) bool {

	l.AddItemLast(value)
	return true
}

// Count returns the number of values in the list.
func (l *SList[T]) Count() int {

	return l.count
}

// IsEmpty returns true if the collection has no elements.
func (l *SList[T]) IsEmpty() bool {
	return l.count == 0
}

// AddRange adds a slice of values to the end of the list.
func (l *SList[T]) AddRange(values []T) {

	if len(values) == 0 {
		return
	}

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	for _, v := range values {
		l.appendNode(&SListNode[T]{
			list: l,
			item: v,
		})
	}

	l.version++
}

// AddCollection adds the values of the given collection to the end of this list.
// Values are added in the order defined by the other collection.
func (l *SList[T]) AddCollection(collection collections.Collection[T]) {

	l.AddRange(collection.ToSliceDeep())
}

// First returns the node at the head of the list.
// Will be nil if the list is empty.
func (l *SList[T]) First() *SListNode[T] {

	return l.head
}

// Last returns the node at the end of the list.
// Will be nil if the list is empty.
func (l *SList[T]) Last() *SListNode[T] {

	return l.tail
}

// AddItemAfter inserts a value after the given node in the list and returns the newly inserted node.
//
// Panics if the node argument is nil or belongs to another list.
func (l *SList[T]) AddItemAfter(node *SListNode[T], value T) *SListNode[T] {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	l.validateNode(node)
	newNode := &SListNode[T]{
		list: l,
		item: value,
	}

	if node.next == nil {
		// node is the tail, so append
		l.appendNode(newNode)
	} else {
		l.insertNodeAfter(node, newNode)
	}

	l.version++
	return newNode
}

// AddNodeAfter inserts newNode after the given node in the list.
//
// Panics if either node argument is nil or belongs to another list.
func (l *SList[T]) AddNodeAfter(node, newNode *SListNode[T]) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	l.validateNode(node)
	l.validateNewNode(newNode)

	if node.next == nil {
		// node is the tail, so append
		l.appendNode(newNode)
	} else {
		l.insertNodeAfter(node, newNode)
	}

	l.version++
}

// AddNodeFirst inserts node at the head of the list.
//
// Panics if node argument is nil or belongs to another list.
func (l *SList[T]) AddNodeFirst(node *SListNode[T]) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	l.validateNewNode(node)
	l.prependNode(node)
	l.version++
}

// AddNodeLast inserts node at the end of the list.
//
// Panics if node argument is nil or belongs to another list.
func (l *SList[T]) AddNodeLast(node *SListNode[T]) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	l.validateNewNode(node)

	l.appendNode(node)
	l.version++
}

// Clear empties the list, detaching and invalidating all nodes.
func (l *SList[T]) Clear() {

	if l.count == 0 {
		return
	}

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	var empty T
	current := l.head
	for current != nil {
		temp := current
		current = current.next
		temp.invalidate()
		temp.item = empty
	}

	l.head = nil
	l.count = 0
	l.version++
}

// Contains returns true if the given value is in the list; else false. Up to O(n).
func (l *SList[T]) Contains(value T) bool {

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	if l.head == nil {
		return false
	}

	for current := l.head; current != nil; current = current.next {
		if l.compare(value, current.item) == 0 {
			return true
		}
	}

	return false
}

// Remove is an alias for [SList.RemoveItem].
func (l *SList[T]) Remove(value T) bool {

	return l.RemoveItem(value)
}

// RemoveItem searches the list for the first occurrence of value
// and removes the node containing that value. Up to O(n).
//
// Returns true if a node was removed; else false.
func (l *SList[T]) RemoveItem(value T) bool {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	node := l.findNode(value)
	if node != nil {
		l.removeNode(node)

		var empty T
		node.item = empty
		return true
	}

	return false
}

// RemoveNode removes the given node from the list.
//
// Panics if node argument is nil or belongs to another list.
func (l *SList[T]) RemoveNode(node *SListNode[T]) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	l.validateNode(node)
	l.removeNode(node)
}

// RemoveFirst removes the node at the head of the list and returns the value that was stored
//
// Panics if list is empty.
func (l *SList[T]) RemoveFirst() T {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	if l.head == nil {
		panic(messages.COLLECTION_EMPTY)
	}

	item := l.head.item
	l.removeNode(l.head)
	return item
}

// RemoveLast removes the node at the end of the list and returns the value that was stored.
//
// Panics if list is empty.
func (l *SList[T]) RemoveLast() T {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	if l.head == nil {
		panic(messages.COLLECTION_EMPTY)
	}

	item := l.tail.item
	l.removeNode(l.tail)
	return item
}

// TryRemoveFirst removes the node at the head of the list and returns the value that was stored and true,
// or the zero value of T and false if the list is empty.
func (l *SList[T]) TryRemoveFirst() (T, bool) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	if l.head == nil {
		var v T
		return v, false
	}

	item := l.head.item
	l.removeNode(l.head)
	return item, true
}

// TryRemoveLast removes the node at the end of the list and returns the value that was stored and true,
// or the zero value of T and false if the list is empty.
func (l *SList[T]) TryRemoveLast() (T, bool) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	if l.head == nil {
		var v T
		return v, false
	}

	item := l.tail.item
	l.removeNode(l.tail)
	return item, true
}

// ToSlice returns a copy of the list content as a slice.
func (l *SList[T]) ToSlice() []T {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	return l.toSlice(false)
}

// ToSliceDeep returns the content of the collection as a slice using the provided [functions.DeepCopyFunc] if any.
//
// Elements are deep copied using the provided [functions.DeepCopyFunc] if any.
func (l *SList[T]) ToSliceDeep() []T {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	return l.toSlice(true)
}

func (l *SList[T]) toSlice(deepCopy bool) []T {
	slc := make([]T, l.count)

	for i, node := 0, l.head; node != nil; i++ {
		if deepCopy {
			slc[i] = util.DeepCopy(node.item, l.copy)
		} else {
			slc[i] = node.item
		}
		node = node.next
	}

	return slc
}

// Type returns the type of this collection.
func (*SList[T]) Type() collections.CollectionType {
	return collections.COLLECTION_SLIST
}

// String returns a string representation of container.
func (l *SList[T]) String() string {

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	var values []string
	for _, value := range l.toSlice(false) {
		values = append(values, fmt.Sprintf("%v", value))
	}

	return "SList\n" + strings.Join(values, ", ")
}

func (l *SList[T]) appendNode(newNode *SListNode[T]) {
	newNode.list = l
	newNode.next = nil
	if l.count == 0 {
		l.head = newNode
		l.tail = newNode
	} else {
		last := l.tail
		last.next = newNode
		l.tail = newNode
	}

	l.count++
}

func (l *SList[T]) prependNode(newNode *SListNode[T]) {
	if l.count == 0 {
		l.head = newNode
		l.tail = newNode
		newNode.next = nil
	} else {
		newNode.next = l.head
		l.head = newNode
	}

	newNode.list = l
	l.count++
}

func (l *SList[T]) insertNodeAfter(node, newNode *SListNode[T]) {

	newNode.next = node.next
	node.next = newNode
	newNode.list = l
	l.count++
}

func (l *SList[T]) validateNode(node *SListNode[T]) {
	if node == nil {
		panic(messages.NIL_NODE)
	}

	if node.list != l {
		panic(messages.FOREIGN_NODE)
	}
}

func (*SList[T]) validateNewNode(node *SListNode[T]) {
	if node == nil {
		panic(messages.NIL_NODE)
	}

	if node.list != nil {
		panic(messages.FOREIGN_NODE)
	}
}

func (l *SList[T]) removeNode(node *SListNode[T]) {
	// validateNode should be called before it gets here.

	if node == l.head {
		if l.count > 1 {
			l.head = node.next
		}
	} else {
		// Requires a search to find the node pointing to this one
		var n *SListNode[T]
		for n = l.head; n != nil; n = n.next {
			if n.next == node {
				if node == l.tail {
					l.tail = n
				}
				n.next = node.next
				break
			}
		}

		if n == nil {
			// Shouldn't get here. Perhaps panic
			return
		}
	}

	node.invalidate()
	if l.count == 1 {
		l.head = nil
		l.tail = nil
	}

	l.count--
	l.version++
}

func (l *SList[T]) addItemLast(value T) *SListNode[T] {

	newNode := &SListNode[T]{
		list: l,
		item: value,
	}

	l.appendNode(newNode)
	l.version++
	return newNode
}

// Make a new empty list with the same attributes as this.
func (l *SList[T]) makeCopy() *SList[T] {
	ll1 := &SList[T]{
		copy:    l.copy,
		compare: l.compare,
	}

	if l.lock != nil {
		ll1.lock = &sync.RWMutex{}
	}

	return ll1
}
