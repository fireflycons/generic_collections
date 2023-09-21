package slist

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestAddLast_Item(t *testing.T) {

	var tempItems, headItems, tailItems []int
	var linkedList *SList[int]
	arraySize := 16
	seed := int64(21543)
	headItems, tailItems, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Verify value is default(T)", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemLast(defaultT[int]())
		initialItems_Tests(t, linkedList, []int{defaultT[int]()})
	})

	t.Run("Call AddLast(T) several times", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemLast(tailItems[i])
		}

		initialItems_Tests(t, linkedList, tailItems)
	})

	t.Run("Call Add(T) several times", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.Add(tailItems[i])
		}

		initialItems_Tests(t, linkedList, tailItems)
	})

	t.Run("Call AddLast(T) several times remove some of the items using custom comparer", func(t *testing.T) {
		linkedList = New(WithComparer(func(a, b int) int {
			if a == b {
				return 0
			}

			if a > b {
				return 1
			}

			return -1
		}))

		for i := 0; i < arraySize; i++ {
			linkedList.AddItemLast(tailItems[i])
		}

		linkedList.RemoveItem(tailItems[2])
		linkedList.RemoveItem(tailItems[len(tailItems)-3])
		linkedList.RemoveItem(tailItems[1])
		linkedList.RemoveItem(tailItems[len(tailItems)-2])
		linkedList.RemoveFirst()
		linkedList.RemoveLast()
		//With the above remove we should have removed the first and last 3 items
		tempItems = make([]int, len(tailItems)-6)
		util.PartialCopy(tailItems, 3, tempItems, 0, len(tailItems)-6)

		initialItems_Tests(t, linkedList, tempItems)

		// adding some more items to the tail of the linked list.
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemLast(headItems[i])
		}

		tempItems2 := make([]int, len(tempItems)+len(headItems))
		copy(tempItems2, tempItems)
		util.PartialCopy(headItems, 0, tempItems2, len(tempItems), len(headItems))
		initialItems_Tests(t, linkedList, tempItems2)
	})

	t.Run("Call AddLast(T) several times then call Clear", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemLast(tailItems[i])
		}

		linkedList.Clear()

		for i := 0; i < arraySize; i++ {
			linkedList.AddItemLast(headItems[i])
		}

		initialItems_Tests(t, linkedList, headItems)
	})

	t.Run("Mix AddHead and AddTail calls", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemFirst(headItems[i])
			linkedList.AddItemLast(tailItems[i])
		}

		tempItems = make([]int, len(headItems))
		// adding the headItems in reverse order.
		for i := 0; i < len(headItems); i++ {
			index := (len(headItems) - 1) - i
			tempItems[i] = headItems[index]
		}

		tempItems2 := make([]int, len(tempItems)+len(tailItems))
		copy(tempItems2, tempItems)
		util.PartialCopy(tailItems, 0, tempItems2, len(tempItems), len(tailItems))
		initialItems_Tests(t, linkedList, tempItems2)
	})
}

func TestAddLast_LLNode(t *testing.T) {
	var tempItems, headItems, headItemsReverse, tailItems []int
	var linkedList *SList[int]
	arraySize := 16
	seed := int64(21543)
	headItems, tailItems, headItemsReverse, _ = util.CreateSerialSmallIntListData(arraySize, &seed)

	t.Run("Verify value is default(T)", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddNodeLast(NewNode(defaultT[int]()))
		initialItems_Tests(t, linkedList, []int{defaultT[int]()})
	})

	t.Run("Call AddLast(SListNode<T>) several times", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddNodeLast(NewNode(tailItems[i]))
		}

		initialItems_Tests(t, linkedList, tailItems)
	})

	t.Run("Call AddLast(SListNode<T>) several times remove some of the items", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddNodeLast(NewNode(tailItems[i]))
		}

		linkedList.RemoveItem(tailItems[2])
		linkedList.RemoveItem(tailItems[len(tailItems)-3])
		linkedList.RemoveItem(tailItems[1])
		linkedList.RemoveItem(tailItems[len(tailItems)-2])
		linkedList.RemoveFirst()
		linkedList.RemoveLast()

		//With the above remove we should have removed the first and last 3 items
		tempItems = make([]int, len(tailItems)-6)
		util.PartialCopy(tailItems, 3, tempItems, 0, len(tailItems)-6)
		initialItems_Tests(t, linkedList, tempItems)

		for i := 0; i < arraySize; i++ {
			linkedList.AddNodeLast(NewNode(headItems[i]))
		}

		tempItems2 := make([]int, len(tempItems)+len(headItems))
		copy(tempItems2, tempItems)
		util.PartialCopy(headItems, 0, tempItems2, len(tempItems), len(headItems))

		initialItems_Tests(t, linkedList, tempItems2)
	})

	t.Run("Call AddLst(T) several times then call Clear", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddNodeLast(NewNode(tailItems[i]))
		}

		linkedList.Clear()

		for i := 0; i < arraySize; i++ {
			linkedList.AddNodeLast(NewNode(headItems[i]))
		}

		initialItems_Tests(t, linkedList, headItems)
	})

	t.Run("Mix AddFirst and AddLast calls", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddNodeFirst(NewNode(headItems[i]))
			linkedList.AddNodeLast(NewNode(tailItems[i]))
		}

		tempItems = make([]int, len(headItemsReverse)+len(tailItems))
		copy(tempItems, headItemsReverse)
		util.PartialCopy(tailItems, 0, tempItems, len(headItemsReverse), len(tailItems))
		initialItems_Tests(t, linkedList, tempItems)
	})
}

func TestAddLast_LLNode_Negative(t *testing.T) {
	var (
		linkedList, tempSList *SList[int]
		items                 []int
	)

	seed := int64(21543)

	t.Run("Verify Null node", func(t *testing.T) {
		linkedList = New[int]()
		require.Panics(t, func() { linkedList.AddNodeLast(nil) })
		initialItems_Tests(t, linkedList, []int{})
	})

	t.Run("Verify Node that already exists in this collection that is the Head", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		require.Panics(t, func() { linkedList.AddNodeLast(linkedList.First()) })
		initialItems_Tests(t, linkedList, items)
	})

	t.Run("Verify Node that already exists in this collection that is the Tail", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed), util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		linkedList.AddItemLast(items[1])
		require.Panics(t, func() { linkedList.AddNodeLast(linkedList.Last()) })
		initialItems_Tests(t, linkedList, items)
	})

	t.Run("Verify Node that already exists in another collection", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed), util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		linkedList.AddItemLast(items[1])

		tempSList = New[int]()
		tempSList.Clear()
		tempSList.AddItemLast(util.CreateRandInt(&seed))
		tempSList.AddItemLast(util.CreateRandInt(&seed))
		require.Panics(t, func() { linkedList.AddNodeLast(tempSList.Last()) })
		initialItems_Tests(t, linkedList, items)
	})
}
