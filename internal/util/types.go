package util

import (
	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/internal/local"
	"github.com/fireflycons/generic_collections/internal/messages"
)

// Properties shared by all iterator types.
type IteratorBase[T any] struct {
	Version    int
	NilElement collections.Element[T]
}

// Concrete representation of Element interface.
type ElementType[T any] struct {
	Collection collections.Collection[T]
	Version    int
	ValueP     *T
	local.InternalImpl
}

func NewElementType[T any](collection collections.Collection[T], val *T) *ElementType[T] {
	return &ElementType[T]{
		Collection: collection,
		Version:    GetVersion[T](collection),
		ValueP:     val,
	}
}

func (e *ElementType[T]) Value() T {
	if e.Version != GetVersion[T](e.Collection) {
		panic(messages.COLLECTION_MODIFIED)
	}
	return *e.ValueP
}

func (e *ElementType[T]) ValuePtr() *T {
	collectionType := e.Collection.Type()
	if collectionType == collections.COLLECTION_HASHSET || collectionType == collections.COLLECTION_ORDEREDSET {
		panic(messages.SET_POINTER_MODIFICATION)
	}
	if e.Version != GetVersion[T](e.Collection) {
		panic(messages.COLLECTION_MODIFIED)
	}
	return e.ValueP
}
