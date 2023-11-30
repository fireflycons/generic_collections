package slist

import (
	"fmt"
	"sort"
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestSort(t *testing.T) {

	var tempItems, headItems []int
	var list *SList[int]
	seed := int64(8293)

	t.Run("Empty list", func(t *testing.T) {
		list = New[int]()
		list.Sort()
		verifyLLState(t, list, []int{})
	})

	for arraySize := 16; arraySize <= 32768; arraySize *= 4 {
		t.Run(fmt.Sprintf("%d elements", arraySize), func(t *testing.T) {
			headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)
			list = New[int]()
			list.AddRange(headItems)
			list.Sort()
			tempItems = make([]int, len(headItems), len(headItems)+1)
			copy(tempItems, headItems)
			sort.Ints(tempItems)
			verifyLLState(t, list, tempItems)

			// Assert tail is correct
			i := util.CreateRandInt(&seed)
			list.AddItemLast(i)
			tempItems = append(tempItems, i)
			verifyLLState(t, list, tempItems)
		})
	}
}

func TestSorted(t *testing.T) {

	var tempItems, headItems []int
	var list *SList[int]
	seed := int64(8293)

	t.Run("Empty list", func(t *testing.T) {
		list = New[int]()
		ll1, ok := list.Sorted().(*SList[int])
		require.True(t, ok, "type assertion failed")
		verifyLLState(t, ll1, []int{})
	})

	for arraySize := 16; arraySize <= 32768; arraySize *= 4 {
		t.Run(fmt.Sprintf("%d elements", arraySize), func(t *testing.T) {
			headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)
			list = New[int]()
			list.AddRange(headItems)
			ll1, ok := list.Sorted().(*SList[int])
			require.True(t, ok, "type assertion failed")
			tempItems = make([]int, len(headItems), len(headItems)+1)
			copy(tempItems, headItems)
			sort.Ints(tempItems)
			verifyLLState(t, ll1, tempItems)

			// Assert tail is correct
			i := util.CreateRandInt(&seed)
			ll1.AddItemLast(i)
			tempItems = append(tempItems, i)
			verifyLLState(t, ll1, tempItems)
		})
	}
}

func TestSortDescending(t *testing.T) {

	var tempItems, headItems []int
	var list *SList[int]
	seed := int64(8293)

	t.Run("Empty list", func(t *testing.T) {
		list = New[int]()
		list.SortDescending()
		verifyLLState(t, list, []int{})
	})

	for arraySize := 16; arraySize <= 32768; arraySize *= 4 {
		t.Run(fmt.Sprintf("%d elements", arraySize), func(t *testing.T) {
			headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)
			list = New[int]()
			list.AddRange(headItems)
			list.SortDescending()
			tempItems = make([]int, len(headItems), len(headItems)+1)
			copy(tempItems, headItems)
			sort.Ints(tempItems)
			util.Reverse(tempItems)
			verifyLLState(t, list, tempItems)

			// Assert tail is correct
			i := util.CreateRandInt(&seed)
			list.AddItemLast(i)
			tempItems = append(tempItems, i)
			verifyLLState(t, list, tempItems)
		})
	}
}

func TestSortedDescending(t *testing.T) {

	var tempItems, headItems []int
	var list *SList[int]
	seed := int64(8293)

	t.Run("Empty list", func(t *testing.T) {
		list = New[int]()
		ll1, ok := list.SortedDescending().(*SList[int])
		require.True(t, ok, "type assertion failed")
		verifyLLState(t, ll1, []int{})
	})

	for arraySize := 16; arraySize <= 32768; arraySize *= 4 {
		t.Run(fmt.Sprintf("%d elements", arraySize), func(t *testing.T) {
			headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)
			list = New[int]()
			list.AddRange(headItems)
			ll1, ok := list.SortedDescending().(*SList[int])
			require.True(t, ok, "type assertion failed")
			tempItems = make([]int, len(headItems), len(headItems)+1)
			copy(tempItems, headItems)
			sort.Ints(tempItems)
			util.Reverse(tempItems)
			verifyLLState(t, ll1, tempItems)

			// Assert tail is correct
			i := util.CreateRandInt(&seed)
			ll1.AddItemLast(i)
			tempItems = append(tempItems, i)
			verifyLLState(t, ll1, tempItems)
		})
	}
}
