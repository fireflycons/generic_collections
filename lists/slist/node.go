package slist

// SListNode represents a node in an SList
type SListNode[T any] struct {
	next *SListNode[T]
	list *SList[T]
	item T
}

// NewNode initializes a new node with given value
func NewNode[T any](value T) *SListNode[T] {
	return &SListNode[T]{
		item: value,
	}
}

// List returns the list that this node belongs to. Can be nil.
func (n *SListNode[T]) List() *SList[T] {
	return n.list
}

// Next returns the next node in the chain.
// Will be nil if this is the last node.
func (n *SListNode[T]) Next() *SListNode[T] {
	if n.next == nil {
		return nil
	}

	if n.next == n.list.head {
		return nil
	}

	return n.next
}

// Value returns the value stored in this node.
func (n *SListNode[T]) Value() T {
	return n.item
}

// SetValue sets the value of this node
func (n *SListNode[T]) SetValue(value T) {
	n.item = value
}

// ValuePtr returns a pointer to this node's value
func (n *SListNode[T]) ValuePtr() *T {
	return &n.item
}

func (n *SListNode[T]) invalidate() {
	n.list = nil
	n.next = nil
}
