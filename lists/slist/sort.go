package slist

import "github.com/fireflycons/generic_collections/collections"

type direction bool

const (
	forward, reverse direction = true, false
)

// Sort performs an in-place sort of this collection with a time complexity of O(n*log n).
func (l *SList[T]) Sort() {

	if l.head == nil || l.count < 2 {
		return
	}

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	l.mergeSort(forward)
	l.version++
}

// Sorted returns a sorted copy of this SList as a new SList using the provided [functions.DeepCopyFunc] if any.
func (l *SList[T]) Sorted() collections.Collection[T] {

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	ll1 := l.makeCopy()

	switch {
	case l.head == nil:
		return ll1

	case l.count == 1:
		ll1.addItemLast(l.head.item)
		return ll1

	default:

		ll1.AddRange(l.toSlice(true))
		ll1.mergeSort(forward)
		return ll1
	}
}

// SortDescending performs an in-place sort of this collection with a time complexity of O(n*log n).
func (l *SList[T]) SortDescending() {

	if l.head == nil || l.count < 2 {
		return
	}

	if l.lock != nil {
		l.lock.Lock()
		defer l.lock.Unlock()
	}

	l.mergeSort(reverse)
	l.version++
}

// Sorted returns a descending order sorted copy of this SList as a new SList using the provided [functions.DeepCopyFunc] if any.
func (l *SList[T]) SortedDescending() collections.Collection[T] {

	if l.lock != nil {
		l.lock.RLock()
		defer l.lock.RUnlock()
	}

	ll1 := l.makeCopy()

	switch {
	case l.head == nil:
		return ll1

	case l.count == 1:
		ll1.addItemLast(l.head.item)
		return ll1

	default:

		ll1.AddRange(l.toSlice(true))
		ll1.mergeSort(reverse)
		return ll1
	}
}

func (l *SList[T]) mergeSort(dir direction) {

	var n *SListNode[T]
	l.head = l.mergeSortRecursive(l.head, dir)

	for n = l.head; n.next != nil; n = n.next {
	}

	l.tail = n
}

func (l *SList[T]) mergeSortRecursive(node *SListNode[T], dir direction) *SListNode[T] {
	if node == nil || node.next == nil {
		return node
	}

	second := split(node)

	// Recur for left and right halves
	node = l.mergeSortRecursive(node, dir)
	second = l.mergeSortRecursive(second, dir)

	// Merge the two sorted halves
	return l.merge(node, second, dir)
}

func split[T any](head *SListNode[T]) *SListNode[T] {
	fast := head
	slow := head

	for fast.next != nil && fast.next.next != nil {
		fast = fast.next.next
		slow = slow.next
	}

	temp := slow.next
	slow.next = nil
	return temp
}

func (l *SList[T]) merge(first, second *SListNode[T], dir direction) *SListNode[T] {
	// If first linked list is empty
	if first == nil {
		return second
	}

	// If second linked list is empty
	if second == nil {
		return first
	}

	var x bool

	if dir == forward {
		// Pick the smaller value
		x = l.compare(first.item, second.item) < 0
	} else {
		// Pick the larger value
		x = l.compare(first.item, second.item) > 0
	}

	if x {
		first.next = l.merge(first.next, second, dir)
		return first
	}

	second.next = l.merge(first, second.next, dir)
	return second
}
