package orderedset

import (
	"sort"
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

var setSize = 1024

func TestForwardIterator(t *testing.T) {

	var setItems, iteratedItems []int
	seed := int64(2163)
	setItems, _, _, _ = util.CreateIntListData(setSize, &seed)

	t.Run("Iterator receives all values added to set", func(t *testing.T) {
		set := New[int]()

		set.AddRange(setItems)

		iter := set.Iterator()

		iteratedItems = make([]int, 0, setSize)

		for e := iter.Start(); e != nil; e = iter.Next() {
			iteratedItems = append(iteratedItems, e.Value())
		}

		tempItems := make([]int, len(setItems))
		copy(tempItems, setItems)
		sort.Ints(tempItems)
		require.Equal(t, tempItems, iteratedItems)
	})

	t.Run("Using ValuePtr on an element panics", func(t *testing.T) {
		set := New[int]()
		set.AddRange(setItems)
		iter := set.Iterator()

		e := iter.Start()
		require.Panics(t, func() { e.ValuePtr() })
	})
}

func TestReverseIterator(t *testing.T) {

	var setItems, iteratedItems []int
	seed := int64(2163)
	setItems, _, _, _ = util.CreateIntListData(setSize, &seed)

	t.Run("Iterator receives all values added to set", func(t *testing.T) {
		set := New[int]()

		set.AddRange(setItems)

		iter := set.ReverseIterator()

		iteratedItems = make([]int, 0, setSize)

		for e := iter.Start(); e != nil; e = iter.Next() {
			iteratedItems = append(iteratedItems, e.Value())
		}

		tempItems := make([]int, len(setItems))
		copy(tempItems, setItems)
		sort.Ints(tempItems)
		tempItems = util.Reverse(tempItems)
		require.Equal(t, tempItems, iteratedItems)
	})
}

func TestTakeWhile(t *testing.T) {
	var setItems, iteratedItems []int
	seed := int64(2163)
	setItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Returns even numbers", func(t *testing.T) {
		set := New[int]()

		set.AddRange(setItems)

		iter := set.TakeWhile(func(val int) bool { return val%2 == 0 })

		iteratedItems = make([]int, 0, util.DefaultCapacity)

		for e := iter.Start(); e != nil; e = iter.Next() {
			iteratedItems = append(iteratedItems, e.Value())
		}

		tempItems := make([]int, 0, util.DefaultCapacity)

		for _, v := range setItems {
			if v%2 == 0 {
				tempItems = append(tempItems, v)
			}
		}

		require.ElementsMatch(t, tempItems, iteratedItems)

	})
}

func TestWhere(t *testing.T) {
	var setItems []int
	seed := int64(2163)
	setItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Returns even numbers", func(t *testing.T) {
		set := New[int]()
		set.AddRange(setItems)
		set1 := set.Select(func(val int) bool { return val%2 == 0 })

		tempItems := make([]int, 0, util.DefaultCapacity)

		for _, v := range setItems {
			if v%2 == 0 {
				tempItems = append(tempItems, v)
			}
		}

		require.ElementsMatch(t, tempItems, set1.ToSlice())

	})
}
