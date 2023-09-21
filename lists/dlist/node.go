package dlist

// DListNode represents a node in a DList
type DListNode[T any] struct {
	prev *DListNode[T]
	next *DListNode[T]
	list *DList[T]
	item T
}

// NewNode initializes a new node with given value
func NewNode[T any](value T) *DListNode[T] {
	return &DListNode[T]{
		item: value,
	}
}

// List returns the list that this node belongs to. Can be nil.
func (n *DListNode[T]) List() *DList[T] {
	return n.list
}

// Next returns the next node in the chain.
// Will be nil if this is the last node.
func (n *DListNode[T]) Next() *DListNode[T] {
	if n.next == nil {
		return nil
	}

	if n.next == n.list.head {
		return nil
	}

	return n.next
}

// Previous retuns the previous node in the list
// Will be nil if this is the first node.
func (n *DListNode[T]) Previous() *DListNode[T] {
	if n.prev == nil {
		return nil
	}

	if n == n.list.head {
		return nil
	}

	return n.prev
}

// Value returns the value stored in this node.
func (n *DListNode[T]) Value() T {
	return n.item
}

// SetValue sets the value of this node
func (n *DListNode[T]) SetValue(value T) {
	n.item = value
}

// ValuePtr returns a pointer to this node's value
func (n *DListNode[T]) ValuePtr() *T {
	return &n.item
}

func (n *DListNode[T]) invalidate() {
	n.list = nil
	n.next = nil
	n.prev = nil
}
