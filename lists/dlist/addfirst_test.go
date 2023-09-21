package dlist

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestAddFirst_Item(t *testing.T) {

	var tempItems, headItems, headItemsReverse, tailItems, tailItemsReverse []int
	var linkedList *DList[int]
	arraySize := 16
	seed := int64(21543)
	headItems, tailItems, headItemsReverse, tailItemsReverse = util.CreateIntListData(arraySize, &seed)

	t.Run("Verify value is default(T)", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(defaultT[int]())
		initialItems_Tests(t, linkedList, []int{defaultT[int]()})
	})

	t.Run("Call AddFirst(T) several times", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemFirst(headItems[i])
		}

		initialItems_Tests(t, linkedList, headItemsReverse)
	})

	t.Run("Call AddFirst(T) several times remove some of the items", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemFirst(headItems[i])
		}

		initialItems_Tests(t, linkedList, headItemsReverse)

		linkedList.RemoveItem(headItems[2])
		linkedList.RemoveItem(headItems[len(headItems)-3])
		linkedList.RemoveItem(headItems[1])
		linkedList.RemoveItem(headItems[len(headItems)-2])
		linkedList.RemoveFirst()
		linkedList.RemoveLast()
		//With the above remove we should have removed the first and last 3 items
		// expected items are headItems in reverse order, or a subset of them.
		tempItems = make([]int, len(headItemsReverse)-6)
		util.PartialCopy(headItemsReverse, 3, tempItems, 0, len(headItemsReverse)-6)
		initialItems_Tests(t, linkedList, tempItems)

		for i := 0; i < arraySize; i++ {
			linkedList.AddItemFirst(tailItems[i])
		}

		tempItems2 := make([]int, len(tempItems)+len(tailItemsReverse))
		copy(tempItems2, tailItemsReverse)
		util.PartialCopy(tempItems, 0, tempItems2, len(tailItemsReverse), len(tempItems))
		initialItems_Tests(t, linkedList, tempItems2)
	})

	t.Run("Call AddFirst(T) several times remove all of the items", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemFirst(headItems[i])
		}

		for i := 0; i < arraySize; i++ {
			linkedList.RemoveFirst()
		}

		for i := 0; i < arraySize; i++ {
			linkedList.AddItemFirst(tailItems[i])
		}

		initialItems_Tests(t, linkedList, tailItemsReverse)
	})

	t.Run("Call AddFirst(T) several times then call Clear", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemFirst(headItems[i])
		}

		linkedList.Clear()

		for i := 0; i < arraySize; i++ {
			linkedList.AddItemFirst(tailItems[i])
		}

		initialItems_Tests(t, linkedList, tailItemsReverse)
	})

	t.Run("Mix AddHead and AddTail calls", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemFirst(headItems[i])
			linkedList.AddItemLast(tailItems[i])
		}

		tempItems = make([]int, len(headItemsReverse)+len(tailItems))
		copy(tempItems, headItemsReverse)
		util.PartialCopy(tailItems, 0, tempItems, len(headItemsReverse), len(tailItems))

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Adding nil node panics", func(t *testing.T) {
		linkedList = New[int]()
		require.Panics(t, func() { linkedList.AddNodeFirst(nil) })
	})

	t.Run("Adding foreign node panics", func(t *testing.T) {
		linkedList = New[int]()
		linkedList2 := New[int]()
		linkedList2.AddItemFirst(1)
		require.Panics(t, func() { linkedList.AddNodeFirst(linkedList2.First()) })
	})

	t.Run("Set value on a node updates list content", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(1)
		linkedList.First().SetValue(2)
		initialItems_Tests(t, linkedList, []int{2})
	})
}

func TestAddFirst_LLNode(t *testing.T) {

	var tempItems, headItems, headItemsReverse, tailItems, tailItemsReverse []int
	var linkedList *DList[int]
	arraySize := 16
	seed := int64(21543)
	headItems, tailItems, headItemsReverse, tailItemsReverse = util.CreateIntListData(arraySize, &seed)

	t.Run("Verify value is default(T)", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddNodeFirst(NewNode(defaultT[int]()))
		initialItems_Tests(t, linkedList, []int{defaultT[int]()})
	})

	t.Run("Call AddFirst(DListNode<T>) several times", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddNodeFirst(NewNode(headItems[i]))
		}

		linkedList.RemoveItem(headItems[2])
		linkedList.RemoveItem(headItems[len(headItems)-3])
		linkedList.RemoveItem(headItems[1])
		linkedList.RemoveItem(headItems[len(headItems)-2])
		linkedList.RemoveFirst()
		linkedList.RemoveLast()
		//With the above remove we should have removed the first and last 3 items
		tempItems = make([]int, len(headItemsReverse)-6)
		util.PartialCopy(headItemsReverse, 3, tempItems, 0, len(headItemsReverse)-6)
		initialItems_Tests(t, linkedList, tempItems)

		for i := 0; i < arraySize; i++ {
			linkedList.AddNodeFirst(NewNode(tailItems[i]))
		}

		tempItems2 := make([]int, len(tailItemsReverse)+len(tempItems))
		copy(tempItems2, tailItemsReverse)
		util.PartialCopy(tempItems, 0, tempItems2, len(tailItemsReverse), len(tempItems))
		initialItems_Tests(t, linkedList, tempItems2)
	})

	t.Run("Call AddFirst(DListNode<T>) several times remove all of the items", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddNodeFirst(NewNode(headItems[i]))
		}

		for i := 0; i < arraySize; i++ {
			linkedList.RemoveFirst()
		}

		for i := 0; i < arraySize; i++ {
			linkedList.AddNodeFirst(NewNode(tailItems[i]))
		}

		initialItems_Tests(t, linkedList, tailItemsReverse)
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

func TestAddFirst_LLNode_Negative(t *testing.T) {
	var (
		linkedList, tempDList *DList[int]
		items                 []int
	)

	seed := int64(21543)

	t.Run("Verify Null node", func(t *testing.T) {
		linkedList = New[int]()
		require.Panics(t, func() { linkedList.AddNodeFirst(nil) })
		initialItems_Tests(t, linkedList, []int{})
	})

	t.Run("Verify Node that already exists in this collection that is the Head", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		require.Panics(t, func() { linkedList.AddNodeFirst(linkedList.First()) })
		initialItems_Tests(t, linkedList, items)
	})

	t.Run("Verify Node that already exists in this collection that is the Tail", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed), util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		linkedList.AddItemLast(items[1])
		require.Panics(t, func() { linkedList.AddNodeFirst(linkedList.Last()) })
		initialItems_Tests(t, linkedList, items)
	})

	t.Run("Verify Node that already exists in another collection", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed), util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		linkedList.AddItemLast(items[1])

		tempDList = New[int]()
		tempDList.Clear()
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		require.Panics(t, func() { linkedList.AddNodeFirst(tempDList.Last()) })
		initialItems_Tests(t, linkedList, items)
	})
}
