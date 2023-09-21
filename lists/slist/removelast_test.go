package slist

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestRemoveLast(t *testing.T) {

	var headItems []int
	var linkedList *SList[int]
	var tempNode1, tempNode2, tempNode3 *SListNode[int]
	arraySize := 16
	seed := int64(21543)
	headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Call RemoveHead on a collection with one item in it", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemLast(headItems[0])
		tempNode1 = linkedList.Last()

		linkedList.RemoveLast()
		initialItems_Tests(t, linkedList, []int{})
		verifyRemovedNode2(t, tempNode1, headItems[0])
	})

	t.Run("Call RemoveHead on a collection with two items in it", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		tempNode1 = linkedList.Last()
		tempNode2 = linkedList.First()

		linkedList.RemoveLast()
		initialItems_Tests(t, linkedList, []int{headItems[0]})

		linkedList.RemoveLast()
		initialItems_Tests(t, linkedList, []int{})

		verifyRemovedNode3(t, linkedList, tempNode1, headItems[1])
		verifyRemovedNode3(t, linkedList, tempNode2, headItems[0])
	})

	t.Run("Call RemoveHead on a collection with three items in it", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		linkedList.AddItemLast(headItems[2])
		tempNode1 = linkedList.Last()
		tempNode2 = linkedList.First().Next()
		tempNode3 = linkedList.First()

		linkedList.RemoveLast()
		initialItems_Tests(t, linkedList, []int{headItems[0], headItems[1]})

		linkedList.RemoveLast()
		initialItems_Tests(t, linkedList, []int{headItems[0]})

		linkedList.RemoveLast()
		initialItems_Tests(t, linkedList, []int{})
		verifyRemovedNode2(t, tempNode1, headItems[2])
		verifyRemovedNode2(t, tempNode2, headItems[1])
		verifyRemovedNode2(t, tempNode3, headItems[0])
	})

	t.Run("Call RemoveHead on a collection with 16 items in it", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemLast(headItems[i])
		}
		for i := 0; i < arraySize; i++ {
			linkedList.RemoveLast()
			length := arraySize - i - 1
			expectedItems := make([]int, length)
			copy(expectedItems, headItems)
			initialItems_Tests(t, linkedList, expectedItems)
		}
	})

	t.Run("Mix RemoveHead and RemoveTail call", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemLast(headItems[i])
		}
		for i := 0; i < arraySize; i++ {
			if (i & 1) == 0 {
				linkedList.RemoveFirst()
			} else {
				linkedList.RemoveLast()
			}
			startIndex := (i / 2) + 1
			length := arraySize - i - 1
			expectedItems := make([]int, length)
			util.PartialCopy(headItems, startIndex, expectedItems, 0, length)
			initialItems_Tests(t, linkedList, expectedItems)
		}
	})
}

func TestRemoveLast_Negative(t *testing.T) {

	t.Run("Call RemoveHead an empty collection", func(t *testing.T) {
		linkedList := New[int]()
		require.Panics(t, func() { linkedList.RemoveLast() })
		initialItems_Tests(t, linkedList, []int{})
	})
}
