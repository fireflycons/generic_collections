package util

import (
	"sort"

	"github.com/fireflycons/generic_collections/functions"
)

type sortable[T any] struct {
	values  []T
	compare functions.ComparerFunc[T]
}

type ascendingSortable[T any] struct {
	sortable[T]
}

type SortFunc[T any] func([]T, int, functions.ComparerFunc[T])

type descendingSortable[T any] struct {
	sortable[T]
}

// Gosort sorts values (in-place) with respect to the given ComparerFunc using Go's built in sort.
func Gosort[T any](values []T, length int, compare functions.ComparerFunc[T]) {
	sort.Sort(ascendingSortable[T]{sortable[T]{values[:length], compare}})
}

// GosortDescending sorts values (in-place) in descending order with respect to the given ComparerFunc using Go's built in sort.
func GosortDescending[T any](values []T, length int, compare functions.ComparerFunc[T]) {
	sort.Sort(descendingSortable[T]{sortable[T]{values[:length], compare}})
}

func (s sortable[T]) Len() int {
	return len(s.values)
}

func (s sortable[T]) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

func (s ascendingSortable[T]) Less(i, j int) bool {
	return s.compare(s.values[i], s.values[j]) < 0
}

func (s descendingSortable[T]) Less(i, j int) bool {
	return s.compare(s.values[i], s.values[j]) > 0
}
