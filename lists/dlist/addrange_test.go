package dlist

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestAddRange(t *testing.T) {

	var additionalItems, headItems, tempItems []int
	var linkedList *DList[int]
	arraySize := 16
	additionalArraySize := 4
	seed := int64(21543)
	headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)
	additionalItems, _, _, _ = util.CreateIntListData(additionalArraySize, &seed)

	t.Run("AddRange to empty list", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddRange(headItems)
		initialItems_Tests(t, linkedList, headItems)
	})

	t.Run("AddRange to populated list", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemLast(headItems[i])
		}

		linkedList.AddRange(additionalItems)
		tempItems = make([]int, arraySize+additionalArraySize)
		copy(tempItems, headItems)
		index := len(headItems)
		for i := 0; i < additionalArraySize; i++ {
			tempItems[index] = additionalItems[i]
			index++
		}
		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("AddCollection to empty list", func(t *testing.T) {
		linkedList1 := New[int]()
		linkedList1.AddRange(headItems)

		linkedList = New[int]()
		linkedList.AddCollection(linkedList1)
		require.Equal(t, len(headItems), linkedList.Count())
		initialItems_Tests(t, linkedList, headItems)
	})

}
