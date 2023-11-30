package dlist

import (
	"fmt"
	"sort"
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestSort(t *testing.T) {

	var tempItems, headItems []int
	var linkedList *DList[int]
	seed := int64(8293)

	t.Run("Empty list", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.Sort()
		verifyLLState(t, linkedList, []int{})
	})

	for arraySize := 16; arraySize <= 32768; arraySize *= 4 {
		t.Run(fmt.Sprintf("%d elements", arraySize), func(t *testing.T) {
			headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)
			linkedList = New[int]()
			linkedList.AddRange(headItems)
			linkedList.Sort()
			tempItems = make([]int, len(headItems), len(headItems)+1)
			copy(tempItems, headItems)
			sort.Ints(tempItems)
			verifyLLState(t, linkedList, tempItems)

			// Assert tail is correct
			i := util.CreateRandInt(&seed)
			linkedList.AddItemLast(i)
			tempItems = append(tempItems, i)
			verifyLLState(t, linkedList, tempItems)
		})
	}
}

func TestSorted(t *testing.T) {

	var tempItems, headItems []int
	var linkedList *DList[int]
	seed := int64(8293)

	t.Run("Empty list", func(t *testing.T) {
		linkedList = New[int]()
		ll1, ok := linkedList.Sorted().(*DList[int])
		require.True(t, ok, "type assertion failed")
		verifyLLState(t, ll1, []int{})
	})

	for arraySize := 16; arraySize <= 32768; arraySize *= 4 {
		t.Run(fmt.Sprintf("%d elements", arraySize), func(t *testing.T) {
			headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)
			linkedList = New[int]()
			linkedList.AddRange(headItems)
			ll1, ok := linkedList.Sorted().(*DList[int])
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
	var linkedList *DList[int]
	seed := int64(8293)

	t.Run("Empty list", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.SortDescending()
		verifyLLState(t, linkedList, []int{})
	})

	for arraySize := 16; arraySize <= 32768; arraySize *= 4 {
		t.Run(fmt.Sprintf("%d elements", arraySize), func(t *testing.T) {
			headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)
			linkedList = New[int]()
			linkedList.AddRange(headItems)
			linkedList.SortDescending()
			tempItems = make([]int, len(headItems), len(headItems)+1)
			copy(tempItems, headItems)
			sort.Ints(tempItems)
			util.Reverse(tempItems)
			verifyLLState(t, linkedList, tempItems)

			// Assert tail is correct
			i := util.CreateRandInt(&seed)
			linkedList.AddItemLast(i)
			tempItems = append(tempItems, i)
			verifyLLState(t, linkedList, tempItems)
		})
	}
}

func TestSortedDescending(t *testing.T) {

	var tempItems, headItems []int
	var linkedList *DList[int]
	seed := int64(8293)

	t.Run("Empty list", func(t *testing.T) {
		linkedList = New[int]()
		ll1, ok := linkedList.SortedDescending().(*DList[int])
		require.True(t, ok, "type assertion failed")
		verifyLLState(t, ll1, []int{})
	})

	for arraySize := 16; arraySize <= 32768; arraySize *= 4 {
		t.Run(fmt.Sprintf("%d elements", arraySize), func(t *testing.T) {
			headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)
			linkedList = New[int]()
			linkedList.AddRange(headItems)
			ll1, ok := linkedList.SortedDescending().(*DList[int])
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
