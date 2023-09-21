package stack

import (
	"sort"
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestSort(t *testing.T) {

	var stackItems, tempItems []int

	seed := int64(2163)
	stackItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Empty stack", func(t *testing.T) {
		stack := New[int]()
		stack.Sort()
		verifyStackState(t, stack, []int{})
	})

	t.Run("Full stack", func(t *testing.T) {
		stack := New[int]()
		stack.AddRange(stackItems)
		stack.Sort()
		tempItems = make([]int, len(stackItems))
		copy(tempItems, stackItems)
		sort.Ints(tempItems)
		verifyStackState(t, stack, tempItems)
		require.Equal(t, tempItems[0], stack.Peek())
	})

	t.Run("Stack with space at end", func(t *testing.T) {
		stack := New(WithCapacity[int](32))
		stack.AddRange(stackItems)
		stack.Sort()
		tempItems = make([]int, len(stackItems))
		copy(tempItems, stackItems)
		sort.Ints(tempItems)
		verifyStackState(t, stack, tempItems)
		require.Equal(t, tempItems[0], stack.Peek())
	})
}

func TestSorted(t *testing.T) {

	var stackItems, tempItems []int

	seed := int64(2163)
	stackItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Empty stack", func(t *testing.T) {
		stack := New[int]()
		s1 := stack.Sorted().(*Stack[int])
		verifyStackState(t, s1, []int{})
	})

	t.Run("Full stack", func(t *testing.T) {
		stack := New[int]()
		stack.AddRange(stackItems)
		s1 := stack.Sorted().(*Stack[int])
		tempItems = make([]int, len(stackItems))
		copy(tempItems, stackItems)
		sort.Ints(tempItems)
		verifyStackState(t, s1, tempItems)
		require.Equal(t, tempItems[0], s1.Peek())
	})

	t.Run("Stack with space at end", func(t *testing.T) {
		stack := New(WithCapacity[int](32))
		stack.AddRange(stackItems)
		s1 := stack.Sorted().(*Stack[int])
		tempItems = make([]int, len(stackItems))
		copy(tempItems, stackItems)
		sort.Ints(tempItems)
		verifyStackState(t, s1, tempItems)
		require.Equal(t, tempItems[0], s1.Peek())
	})
}

func TestSortDescending(t *testing.T) {

	var stackItems, tempItems []int

	seed := int64(2163)
	stackItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Empty stack", func(t *testing.T) {
		stack := New[int]()
		stack.SortDescending()
		verifyStackState(t, stack, []int{})
	})

	t.Run("Full stack", func(t *testing.T) {
		stack := New[int]()
		stack.AddRange(stackItems)
		stack.SortDescending()
		tempItems = make([]int, len(stackItems))
		copy(tempItems, stackItems)
		sort.Ints(tempItems)
		util.Reverse(tempItems)
		verifyStackState(t, stack, tempItems)
		require.Equal(t, tempItems[0], stack.Peek())
	})

	t.Run("Stack with space at end", func(t *testing.T) {
		stack := New(WithCapacity[int](32))
		stack.AddRange(stackItems)
		stack.SortDescending()
		tempItems = make([]int, len(stackItems))
		copy(tempItems, stackItems)
		sort.Ints(tempItems)
		util.Reverse(tempItems)
		verifyStackState(t, stack, tempItems)
		require.Equal(t, tempItems[0], stack.Peek())
	})
}

func TestSortedDescending(t *testing.T) {

	var stackItems, tempItems []int

	seed := int64(2163)
	stackItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Empty stack", func(t *testing.T) {
		stack := New[int]()
		s1 := stack.SortedDescending().(*Stack[int])
		verifyStackState(t, s1, []int{})
	})

	t.Run("Full stack", func(t *testing.T) {
		stack := New[int]()
		stack.AddRange(stackItems)
		s1 := stack.SortedDescending().(*Stack[int])
		tempItems = make([]int, len(stackItems))
		copy(tempItems, stackItems)
		sort.Ints(tempItems)
		util.Reverse(tempItems)
		verifyStackState(t, s1, tempItems)
		require.Equal(t, tempItems[0], s1.Peek())
	})

	t.Run("Stack with space at end", func(t *testing.T) {
		stack := New(WithCapacity[int](32))
		stack.AddRange(stackItems)
		s1 := stack.SortedDescending().(*Stack[int])
		tempItems = make([]int, len(stackItems))
		copy(tempItems, stackItems)
		sort.Ints(tempItems)
		util.Reverse(tempItems)
		verifyStackState(t, s1, tempItems)
		require.Equal(t, tempItems[0], s1.Peek())
	})
}
