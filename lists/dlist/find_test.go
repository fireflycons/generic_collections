package dlist

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestFind_T(t *testing.T) {

	var tempItems, headItems, tailItems []int
	var linkedList *DList[int]
	defaultValue := defaultT[int]()
	arraySize := 16
	seed := int64(21543)
	headItems, tailItems, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Call Find an empty collection", func(t *testing.T) {
		linkedList = New[int]()
		require.Nil(t, linkedList.Find(func(v int) bool { return v == headItems[0] }))
		require.Nil(t, linkedList.Find(func(v int) bool { return v == defaultValue }))
	})

	t.Run("Call Find on a collection with one item in it", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemLast(headItems[0])
		require.Nil(t, linkedList.Find(func(v int) bool { return v == headItems[1] }))
		require.Nil(t, linkedList.Find(func(v int) bool { return v == defaultValue }))
		verifyFind(t, linkedList, []int{headItems[0]})
	})

	t.Run("Call Find on a collection with two items in it", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		require.Nil(t, linkedList.Find(func(v int) bool { return v == headItems[2] }))
		require.Nil(t, linkedList.Find(func(v int) bool { return v == defaultValue }))
		verifyFind(t, linkedList, []int{headItems[0], headItems[1]})
	})

	t.Run("Call Find on a collection with three items in it", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		linkedList.AddItemLast(headItems[2])
		require.Nil(t, linkedList.Find(func(v int) bool { return v == headItems[3] }))
		require.Nil(t, linkedList.Find(func(v int) bool { return v == defaultValue }))
		verifyFind(t, linkedList, []int{headItems[0], headItems[1], headItems[2]})
	})

	t.Run("Call Find on a collection with multiple items in it", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}
		require.Nil(t, linkedList.Find(func(v int) bool { return v == tailItems[0] }))
		require.Nil(t, linkedList.Find(func(v int) bool { return v == defaultValue }))
		verifyFind(t, linkedList, headItems)
	})

	t.Run("Call Find on a collection with duplicate items in it", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}
		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}

		require.Nil(t, linkedList.Find(func(v int) bool { return v == tailItems[0] }))
		require.Nil(t, linkedList.Find(func(v int) bool { return v == defaultValue }))
		tempItems := make([]int, len(headItems)+len(headItems))
		copy(tempItems, headItems)
		util.PartialCopy(headItems, 0, tempItems, len(headItems), len(headItems))
		verifyFindDuplicates(t, linkedList, tempItems)

	})

	t.Run("Call Find with default(T)) at the beginning", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}
		linkedList.AddItemFirst(defaultValue)

		require.Nil(t, linkedList.Find(func(v int) bool { return v == tailItems[0] }))

		tempItems := make([]int, len(headItems)+1)
		tempItems[0] = defaultValue
		util.PartialCopy(headItems, 0, tempItems, 1, len(headItems))

		verifyFind(t, linkedList, tempItems)
	})

	t.Run("Call Find with default(T)) in the middle", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}
		linkedList.AddItemLast(defaultValue)
		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(tailItems[i])
		}

		randVal := util.CreateRandInt(&seed)
		require.Nil(t, linkedList.Find(func(v int) bool { return v == randVal }))

		// prepending tempitems2 to tailitems into tempitems
		tempItems := make([]int, len(tailItems)+1)
		tempItems[0] = defaultValue
		util.PartialCopy(tailItems, 0, tempItems, 1, len(tailItems))

		tempItems2 := make([]int, len(headItems)+len(tempItems))
		copy(tempItems2, headItems)
		util.PartialCopy(tempItems, 0, tempItems2, len(headItems), len(tempItems))

		verifyFind(t, linkedList, tempItems2)
	})

	t.Run("Call Find on a collection with duplicate items in it", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}
		linkedList.AddItemLast(defaultValue)

		require.Nil(t, linkedList.Find(func(v int) bool { return v == tailItems[0] }))
		tempItems = make([]int, len(headItems)+1)
		tempItems[len(headItems)] = defaultValue
		copy(tempItems, headItems)
		verifyFind(t, linkedList, tempItems)
	})

}

func TestFindAll(t *testing.T) {
	var tempItems, headItems []int
	var linkedList *DList[int]
	arraySize := 16
	seed := int64(21543)
	headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Finds all even numbers", func(t *testing.T) {
		linkedList = New[int]()
		expected := make([]int, 0, len(headItems))
		for i := 0; i < len(headItems); i++ {
			if headItems[i]%2 == 0 {
				expected = append(expected, headItems[i])
			}
		}
		linkedList.AddRange(headItems)
		elems := linkedList.FindAll(func(v int) bool { return v%2 == 0 })
		tempItems = make([]int, len(elems))
		for i := 0; i < len(elems); i++ {
			tempItems[i] = elems[i].Value()
		}

		require.Equal(t, expected, tempItems)
	})
}
