package slist

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestFindLast_T(t *testing.T) {
	var headItems, tailItems, prependDefaultHeadItems, prependDefaultTailItems []int
	var linkedList *SList[int]
	arraySize := 16
	seed := int64(21543)
	headItems, tailItems, _, _ = util.CreateIntListData(arraySize, &seed)
	defaultVal := defaultT[int]()
	randVal := util.CreateRandInt(&seed)

	prependDefaultHeadItems = make([]int, len(headItems)+1)
	prependDefaultHeadItems[0] = defaultT[int]()
	util.PartialCopy(headItems, 0, prependDefaultHeadItems, 1, len(headItems))

	prependDefaultTailItems = make([]int, len(tailItems)+1)
	prependDefaultTailItems[0] = defaultT[int]()
	util.PartialCopy(tailItems, 0, prependDefaultTailItems, 1, len(tailItems))

	t.Run("Call FindLast an empty collection", func(t *testing.T) {
		linkedList = New[int]()
		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == headItems[0] }))
		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == defaultVal }))
	})

	t.Run("Call FindLast on a collection with one item in it", func(t *testing.T) {

		linkedList = New[int]()
		linkedList.AddItemLast(headItems[0])
		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == headItems[1] }))
		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == defaultVal }))
		verifyFindLast(t, linkedList, []int{headItems[0]})
	})

	t.Run("Call FindLast on a collection with two items in it", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == headItems[2] }))
		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == defaultVal }))
		verifyFindLast(t, linkedList, []int{headItems[0], headItems[1]})
	})

	t.Run("Call FindLast on a collection with three items in it", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		linkedList.AddItemLast(headItems[2])
		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == headItems[3] }))
		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == defaultVal }))
		verifyFindLast(t, linkedList, []int{headItems[0], headItems[1], headItems[2]})
	})

	t.Run("Call FindLast on a collection with multiple items in it", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}

		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == tailItems[0] }))
		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == defaultVal }))
		verifyFindLast(t, linkedList, headItems)
	})

	t.Run("Call FindLast on a collection with duplicate items in it", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}
		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}

		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == tailItems[0] }))
		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == defaultVal }))
		tempItems := make([]int, len(headItems)+len(headItems))
		copy(tempItems, headItems)
		util.PartialCopy(headItems, 0, tempItems, len(headItems), len(headItems))
		verifyFindLastDuplicates(t, linkedList, tempItems)
	})

	t.Run("Call FindLast with defaultVal at the beginning", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}
		linkedList.AddItemFirst(defaultVal)

		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == tailItems[0] }))
		verifyFindLast(t, linkedList, prependDefaultHeadItems)
	})

	t.Run("Call FindLast with defaultVal in the middle", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(headItems[i])
		}
		linkedList.AddItemLast(defaultVal)
		for i := 0; i < len(headItems); i++ {
			linkedList.AddItemLast(tailItems[i])
		}

		require.Nil(t, linkedList.FindLast(func(v int) bool { return v == randVal }))

		tempItems := make([]int, len(headItems)+len(prependDefaultTailItems))
		copy(tempItems, headItems)
		util.PartialCopy(prependDefaultTailItems, 0, tempItems, len(headItems), len(prependDefaultTailItems))

		verifyFindLast(t, linkedList, tempItems)
	})
}
