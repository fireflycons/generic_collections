package dlist

import (
	"github.com/fireflycons/generic_collections/collections"
)

// DList has its own implementation of sort.
// Where other collections use the common slice quick sort methods,
// DList sorts faster using a merge sort. While creating this
// I benchmarked qsort vs merge sort and merge sort was 100x faster
// over 16K elements.

// Sort performs an in-place sort of this collection with a time complexity of O(n*log n).
func (l *DList[T]) Sort() {

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

// Sorted returns a sorted copy of this DList as a new DList using the provided [functions.DeepCopyFunc] if any.
func (l *DList[T]) Sorted() collections.Collection[T] {

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
func (l *DList[T]) SortDescending() {

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

// Sorted returns a descending order sorted copy of this DList as a new DList using the provided [functions.DeepCopyFunc] if any.
func (l *DList[T]) SortedDescending() collections.Collection[T] {

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

func (ll *DList[T]) mergeSort(dir direction) {

	var n *DListNode[T]
	ll.head = ll.mergeSortRecursive(ll.head, dir)

	for n = ll.head; n.next != nil; n = n.next {
	}

	ll.tail = n
}

func (ll *DList[T]) mergeSortRecursive(node *DListNode[T], dir direction) *DListNode[T] {
	if node == nil || node.next == nil {
		return node
	}

	second := split(node)

	// Recur for left and right halves
	node = ll.mergeSortRecursive(node, dir)
	second = ll.mergeSortRecursive(second, dir)

	// Merge the two sorted halves
	return ll.merge(node, second, dir)
}

func split[T any](head *DListNode[T]) *DListNode[T] {
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

func (ll *DList[T]) merge(first, second *DListNode[T], dir direction) *DListNode[T] {
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
		x = ll.compare(first.item, second.item) < 0
	} else {
		// Pick the larger value
		x = ll.compare(first.item, second.item) > 0
	}

	if x {
		first.next = ll.merge(first.next, second, dir)
		first.next.prev = first
		first.prev = nil
		return first
	}

	second.next = ll.merge(first, second.next, dir)
	second.next.prev = second
	second.prev = nil
	return second

}
