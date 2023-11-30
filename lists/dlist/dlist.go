/*
Package dlist implements a doubly linked list.
*/
package dlist

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

// Assert DList implements required interfaces
var _ lists.List[int] = (*DList[int])(nil)
var _ collections.ReverseIterable[int] = (*DList[int])(nil)

// DListOptionFunc is the signature of a function
// for providing options to the DList constructor.
type DListOptionFunc[T any] func(*DList[T])

// DList represents a doubly linked list of elements of type T.
type DList[T any] struct {
	version    int
	lock       *sync.RWMutex
	head       *DListNode[T]
	tail       *DListNode[T]
	count      int
	compare    functions.ComparerFunc[T]
	copy       functions.DeepCopyFunc[T]
	concurrent bool
	local.InternalImpl
}

// New constructs a linked list.
func New[T any](options ...DListOptionFunc[T]) *DList[T] {
	ll := &DList[T]{}

	for _, o := range options {
		o(ll)
	}

	if ll.copy == nil {
		ll.copy = util.DefaultDeepCopy[T]
	}

	if ll.compare == nil {
		ll.compare = util.GetDefaultComparer[T]()
	}

	return ll
}

// Option function for New to make the collection thread-safe. Adds overhead.
func WithThreadSafe[T any]() DListOptionFunc[T] {
	return func(ll *DList[T]) {
		ll.lock = &sync.RWMutex{}
	}
}

// Option function for NewDList to provide a comparer function for values of type T.
// Required if the element type is not numeric, bool, pointer or string.
func WithComparer[T any](comparer functions.ComparerFunc[T]) DListOptionFunc[T] {
	if comparer == nil {
		panic(messages.COMP_FN_NIL)
	}
	return func(ll *DList[T]) {
		ll.compare = comparer
	}
}

// Option func to provide a deep copy implementation for collection elements.
func WithDeepCopy[T any](copier functions.DeepCopyFunc[T]) DListOptionFunc[T] {
	// Can be nil
	return func(ll *DList[T]) {
		ll.copy = copier
	}
}

// AddItemFirst adds the given value at the head of the list and returns the newly inserted node.
func (l *DList[T]) AddItemFirst(value T) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	newNode := &DListNode[T]{
		list: l,
		item: value,
	}

	l.prependNode(newNode)
	l.version++
}

// AddItemLast adds the given value at the end of the list and returns the newly inserted node.
func (l *DList[T]) AddItemLast(value T) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	newNode := &DListNode[T]{
		list: l,
		item: value,
	}

	l.appendNode(newNode)
	l.version++
}

// Add adds a value to the end of the list.
//
// Always returns true
func (l *DList[T]) Add(value T) bool {

	l.AddItemLast(value)
	return true
}

// Count returns the number of values in the list.
func (l *DList[T]) Count() int {

	return l.count
}

// IsEmpty returns true if the collection has no elements
func (l *DList[T]) IsEmpty() bool {
	return l.count == 0
}

// AddRange adds a slice of values to the end of the list.
func (l *DList[T]) AddRange(values []T) {

	if len(values) == 0 {
		return
	}

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	for _, v := range values {
		l.appendNode(&DListNode[T]{
			list: l,
			item: v,
		})
	}

	l.version++
}

// AddCollection adds the values of the given collection to the end of this list.
// Values are added in the order defined by the other collection.
func (l *DList[T]) AddCollection(collection collections.Collection[T]) {

	l.AddRange(collection.ToSliceDeep())
}

// First returns the node at the head of the list.
// Will be nil if the list is empty.
func (l *DList[T]) First() *DListNode[T] {

	return l.head
}

// Last returns the node at the end of the list.
// Will be nil if the list is empty.
func (l *DList[T]) Last() *DListNode[T] {

	return l.tail
}

// AddItemAfter inserts a value after the given node in the list and returns the newly inserted node.
//
// Panics if the node argument is nil or belongs to another list.
func (l *DList[T]) AddItemAfter(node *DListNode[T], value T) { //*DListNode[T] {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	l.validateNode(node)
	newNode := &DListNode[T]{
		list: l,
		item: value,
	}

	if node.next == nil {
		// node is the tail, so append
		l.appendNode(newNode)
	} else {
		l.insertNodeBefore(node.next, newNode)
	}

	l.version++
	//return newNode
}

// AddNodeAfter inserts newNode after the given node in the list.
//
// Panics if either node argument is nil or belongs to another list.
func (l *DList[T]) AddNodeAfter(node, newNode *DListNode[T]) {

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
		l.insertNodeBefore(node.next, newNode)
	}

	l.version++
}

// AddItemBefore inserts a value before the given node in the list and returns the newly inserted node.
//
// Panics if the node argument is nil or belongs to another list.
func (l *DList[T]) AddItemBefore(node *DListNode[T], value T) *DListNode[T] {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	l.validateNode(node)
	result := &DListNode[T]{
		list: l,
		item: value,
	}

	l.insertNodeBefore(node, result)

	if node == l.head {
		l.head = result
	}

	l.version++
	return result
}

// AddNodeBefore inserts newNode before the given node in the list.
//
// Panics if either node argument is nil or belongs to another list.
func (l *DList[T]) AddNodeBefore(node, newNode *DListNode[T]) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	l.validateNode(node)
	l.validateNewNode(newNode)

	l.insertNodeBefore(node, newNode)

	if node == l.head {
		l.head = newNode
	}
	l.version++
}

// AddNodeFirst inserts node at the head of the list.
//
// Panics if node argument is nil or belongs to another list.
func (l *DList[T]) AddNodeFirst(node *DListNode[T]) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	l.validateNewNode(node)

	if l.head == nil {
		l.appendNode(node)
	} else {
		l.insertNodeBefore(l.head, node)
		l.head = node
	}

	node.list = l
	l.version++
}

func (l *DList[T]) addItemLast(value T) *DListNode[T] {

	newNode := &DListNode[T]{
		list: l,
		item: value,
	}

	l.appendNode(newNode)
	l.version++
	return newNode
}

// AddNodeLast inserts node at the end of the list.
//
// Panics if node argument is nil or belongs to another list.
func (l *DList[T]) AddNodeLast(node *DListNode[T]) {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	l.validateNewNode(node)

	l.appendNode(node)
	l.version++
}

// Clear empties the list, detaching and invalidating all nodes
func (l *DList[T]) Clear() {

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
func (l *DList[T]) Contains(value T) bool {

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

// Remove is an alias for [linkedlist.RemoveItem]
func (l *DList[T]) Remove(value T) bool {

	return l.RemoveItem(value)
}

// RemoveItem searches the list for the first occurrence of value
// and removes the node containing that value. Up to O(n).
//
// Returns true if a node was removed; else false.
func (l *DList[T]) RemoveItem(value T) bool {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	node := l.findNode(value, forward)
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
func (l *DList[T]) RemoveNode(node *DListNode[T]) {

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
func (l *DList[T]) RemoveFirst() T {

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

// RemoveList removes the node at the end of the list and returns the value that was stored.
//
// Panics if list is empty.
func (l *DList[T]) RemoveLast() T {

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
func (l *DList[T]) TryRemoveFirst() (T, bool) {

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
func (l *DList[T]) TryRemoveLast() (T, bool) {

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
func (l *DList[T]) ToSlice() []T {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	return l.toSlice(false)
}

// ToSliceDeep returns the content of the collection as a slice using the provided [functions.DeepCopyFunc] if any.
//
// Elements are deep copied using the provided [functions.DeepCopyFunc] if any.
func (l *DList[T]) ToSliceDeep() []T {

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	return l.toSlice(true)
}

func (l *DList[T]) toSlice(deepCopy bool) []T {
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

// Type returns the type of this collection
func (*DList[T]) Type() collections.CollectionType {
	return collections.COLLECTION_DLIST
}

// String returns a string representation of container
func (l *DList[T]) String() string {

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	var values []string
	for _, value := range l.toSlice(false) {
		values = append(values, fmt.Sprintf("%v", value))
	}

	return "DList\n" + strings.Join(values, ", ")
}

func (ll *DList[T]) appendNode(newNode *DListNode[T]) {
	newNode.list = ll
	newNode.next = nil

	if ll.count == 0 {
		newNode.prev = nil
		ll.head = newNode
		ll.tail = newNode
	} else {
		last := ll.tail
		last.next = newNode
		newNode.prev = last
		ll.tail = newNode
	}

	ll.count++
}

func (ll *DList[T]) prependNode(newNode *DListNode[T]) {
	newNode.list = ll
	newNode.prev = nil

	if ll.count == 0 {
		ll.head = newNode
		ll.tail = newNode
		newNode.next = nil
	} else {
		ll.head.prev = newNode
		newNode.next = ll.head
		ll.head = newNode
	}

	ll.count++
}

func (ll *DList[T]) insertNodeBefore(nextNode, newNode *DListNode[T]) {

	previousNode := nextNode.prev
	newNode.prev = previousNode
	newNode.next = nextNode
	nextNode.prev = newNode
	newNode.list = ll

	if previousNode == nil {
		ll.head = newNode
	} else {
		previousNode.next = newNode
	}

	ll.count++
}

func (ll *DList[T]) validateNode(node *DListNode[T]) {
	if node == nil {
		panic(messages.NIL_NODE)
	}

	if node.list != ll {
		panic(messages.FOREIGN_NODE)
	}
}

func (*DList[T]) validateNewNode(node *DListNode[T]) {
	if node == nil {
		panic(messages.NIL_NODE)
	}

	if node.list != nil {
		panic(messages.FOREIGN_NODE)
	}
}

func (ll *DList[T]) removeNode(node *DListNode[T]) {
	// validateNode should be called before it gets here.

	if node == ll.head {
		if ll.count > 1 {
			ll.head = node.next
			ll.head.prev = nil
		}
	} else if node == ll.tail {
		ll.tail = node.prev
		ll.tail.next = nil
	} else {
		node.next.prev = node.prev
		node.prev.next = node.next
	}

	node.invalidate()
	if ll.count == 1 {
		ll.head = nil
		ll.tail = nil
	}

	ll.count--
	ll.version++
}

// Make a new empty list with the same attributes as this
func (ll *DList[T]) makeCopy() *DList[T] {
	ll1 := &DList[T]{
		copy:    ll.copy,
		compare: ll.compare,
	}

	if ll.lock != nil {
		ll1.lock = &sync.RWMutex{}
	}

	return ll1
}
