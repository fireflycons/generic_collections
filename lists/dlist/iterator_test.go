package dlist

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestValidateIteratorNotNil(t *testing.T) {
	var iter *DListIterator[int]

	require.Panics(t, func() { iter.Start() })
}

func TestLLForwardIterator(t *testing.T) {

	var headItems []int
	var linkedList *DList[int]
	arraySize := 16
	seed := int64(8293)
	headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Calling start on emply list returns nil element", func(t *testing.T) {
		linkedList = New[int]()
		iter := linkedList.Iterator()
		e := iter.Start()
		require.Nil(t, e)
	})

	t.Run("Iterates all values", func(t *testing.T) {
		linkedList = New[int]()

		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}

		index := 0
		iter := linkedList.Iterator()

		for e := iter.Start(); e != nil; e = iter.Next() {
			require.Equal(t, e.Value(), headItems[index])
			index++
		}

		initialItems_Tests(t, linkedList, headItems)
	})

	t.Run("Calling start mid-iteration restarts iteration", func(t *testing.T) {
		linkedList = New[int]()

		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}

		iter := linkedList.Iterator()

		e := iter.Start()
		require.Equal(t, e.Value(), headItems[0])
		e = iter.Next()
		require.Equal(t, e.Value(), headItems[1])
		e = iter.Start()
		require.Equal(t, e.Value(), headItems[0])
	})

	t.Run("Using ValuePtr on an element does not panic", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemLast(util.CreateRandInt(&seed))
		iter := linkedList.Iterator()

		e := iter.Start()
		require.NotPanics(t, func() { e.ValuePtr() })
	})
}

func TestLLForwardIterator_Negative(t *testing.T) {

	var linkedList *DList[int]
	seed := int64(8293)

	t.Run("Modifying collection invalidates iterator", func(t *testing.T) {

		linkedList = New[int]()

		for i := 0; i < 3; i++ {
			linkedList.AddItemLast(util.CreateRandInt(&seed))
		}

		iter := linkedList.Iterator()
		require.NotPanics(t, func() { iter.Start() })
		linkedList.AddItemLast(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Next() })
	})
}

func TestTakeWhile(t *testing.T) {
	var listItems, iteratedItems []int
	seed := int64(2163)
	listItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Returns even numbers", func(t *testing.T) {
		ll := New[int]()

		ll.AddRange(listItems)

		iter := ll.TakeWhile(func(val int) bool { return val%2 == 0 })

		iteratedItems = make([]int, 0, util.DefaultCapacity)

		for e := iter.Start(); e != nil; e = iter.Next() {
			iteratedItems = append(iteratedItems, e.Value())
		}

		tempItems := make([]int, 0, util.DefaultCapacity)

		for _, v := range listItems {
			if v%2 == 0 {
				tempItems = append(tempItems, v)
			}
		}

		require.ElementsMatch(t, tempItems, iteratedItems)

	})
}

func TestWhere(t *testing.T) {
	var listItems []int
	seed := int64(2163)
	listItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Returns even numbers", func(t *testing.T) {
		ll := New[int]()
		ll.AddRange(listItems)
		ll1 := ll.Select(func(val int) bool { return val%2 == 0 })

		tempItems := make([]int, 0, util.DefaultCapacity)

		for _, v := range listItems {
			if v%2 == 0 {
				tempItems = append(tempItems, v)
			}
		}

		require.ElementsMatch(t, tempItems, ll1.ToSlice())

	})
}

func TestLLReverseIterator(t *testing.T) {

	var headItems, headItemsReverse []int
	var linkedList *DList[int]
	arraySize := 16
	seed := int64(8293)
	headItems, _, headItemsReverse, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Iterates all values", func(t *testing.T) {
		linkedList = New[int]()

		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}

		index := 0
		iter := linkedList.ReverseIterator()

		for e := iter.Start(); e != nil; e = iter.Next() {
			require.Equal(t, e.Value(), headItemsReverse[index])
			index++
		}

		initialItems_Tests(t, linkedList, headItems)
	})
}

func TestLLReverseIterator_Negative(t *testing.T) {

	var linkedList *DList[int]
	seed := int64(8293)

	t.Run("Modifying collection invalidates iterator", func(t *testing.T) {

		linkedList = New[int]()

		for i := 0; i < 3; i++ {
			linkedList.AddItemLast(util.CreateRandInt(&seed))
		}

		iter := linkedList.ReverseIterator()
		require.NotPanics(t, func() { iter.Start() })
		linkedList.AddItemLast(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Next() })
	})
}
