package dlist

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestAddBefore_LLNode(t *testing.T) {

	var tempItems, headItems, headItemsReverse, tailItems, tailItemsReverse []int
	var linkedList *DList[int]
	arraySize := 16
	seed := int64(8293)
	headItems, tailItems, headItemsReverse, tailItemsReverse = util.CreateIntListData(arraySize, &seed)

	t.Run("Verify value is default(T)", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemBefore(linkedList.First(), defaultT[int]())
		initialItems_Tests(t, linkedList, []int{defaultT[int](), headItems[0]})
	})

	t.Run("Node is the Head", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddItemBefore(linkedList.First(), headItems[i])
		}

		initialItems_Tests(t, linkedList, headItemsReverse)
	})

	t.Run("Node is the Tail", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		tempItems = make([]int, len(headItems))
		util.PartialCopy(headItems, 1, tempItems, 0, len(headItems)-1)
		tempItems[len(tempItems)-1] = headItems[0]
		for i := 1; i < arraySize; i++ {
			linkedList.AddItemBefore(linkedList.Last(), headItems[i])
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Node is after the Head", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])

		tempItems = make([]int, len(headItems))
		copy(tempItems, headItems)
		util.ReverseSubset(tempItems, 1, len(headItems)-1)

		for i := 2; i < arraySize; i++ {
			linkedList.AddItemBefore(linkedList.First().Next(), headItems[i])
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Node is before the Tail", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])

		tempItems = make([]int, len(headItems))
		util.PartialCopy(headItems, 2, tempItems, 0, len(headItems)-2)
		tempItems[len(tempItems)-2] = headItems[0]
		tempItems[len(tempItems)-1] = headItems[1]

		for i := 2; i < arraySize; i++ {
			linkedList.AddItemBefore(linkedList.Last().Previous(), headItems[i])
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Node is somewhere in the middle", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		linkedList.AddItemLast(headItems[2])

		tempItems = make([]int, len(headItems))
		util.PartialCopy(headItems, 3, tempItems, 0, len(headItems)-3)
		tempItems[len(tempItems)-3] = headItems[0]
		tempItems[len(tempItems)-2] = headItems[1]
		tempItems[len(tempItems)-1] = headItems[2]

		for i := 3; i < arraySize; i++ {
			linkedList.AddItemBefore(linkedList.Last().Previous().Previous(), headItems[i])
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Call AddBefore several times remove some of the items", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddItemBefore(linkedList.First(), headItems[i])
		}

		linkedList.Remove(headItems[2])
		linkedList.Remove(headItems[len(headItems)-3])
		linkedList.Remove(headItems[1])
		linkedList.Remove(headItems[len(headItems)-2])
		linkedList.RemoveFirst()
		linkedList.RemoveLast()

		//With the above remove we should have removed the first and last 3 items
		tempItems = make([]int, len(headItemsReverse)-6)
		util.PartialCopy(headItemsReverse, 3, tempItems, 0, len(headItemsReverse)-6)

		initialItems_Tests(t, linkedList, tempItems)

		for i := 0; i < arraySize; i++ {
			linkedList.AddItemBefore(linkedList.First(), tailItems[i])
		}

		tempItems2 := make([]int, len(tailItemsReverse)+len(tempItems))
		copy(tempItems2, tailItemsReverse)
		util.PartialCopy(tempItems, 0, tempItems2, len(tailItemsReverse), len(tempItems))

		initialItems_Tests(t, linkedList, tempItems2)
	})

	t.Run("Call AddBefore several times remove all of the items", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddItemBefore(linkedList.First(), headItems[i])
		}

		for i := 0; i < arraySize; i++ {
			linkedList.RemoveFirst()
		}

		linkedList.AddItemFirst(tailItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddItemBefore(linkedList.First(), tailItems[i])
		}

		initialItems_Tests(t, linkedList, tailItemsReverse)
	})

	t.Run("Call AddBefore several times then call Clear", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddItemBefore(linkedList.First(), headItems[i])
		}

		linkedList.Clear()

		linkedList.AddItemFirst(tailItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddItemBefore(linkedList.First(), tailItems[i])
		}

		initialItems_Tests(t, linkedList, tailItemsReverse)
	})

	t.Run("Mix AddBefore and AddAfter calls", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemLast(headItems[0])
		linkedList.AddItemLast(tailItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddItemBefore(linkedList.First(), headItems[i])
			linkedList.AddItemAfter(linkedList.Last(), tailItems[i])
		}

		tempItems = make([]int, len(headItemsReverse)+len(tailItems))
		copy(tempItems, headItemsReverse)
		util.PartialCopy(tailItems, 0, tempItems, len(headItemsReverse), len(tailItems))

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Remove item that is not in list returrs false", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemLast(1)
		require.False(t, linkedList.Remove(2))
	})
}

func TestAddBefore_LLNode_Negative(t *testing.T) {
	var (
		linkedList, tempDList *DList[int]
		items                 []int
	)

	seed := int64(8293)

	t.Run("Verify Null node", func(t *testing.T) {
		linkedList = New[int]()
		require.Panics(t, func() { linkedList.AddItemBefore(nil, util.CreateRandInt(&seed)) })
		initialItems_Tests(t, linkedList, []int{})
	})

	t.Run("Verify Node that is a new Node", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		require.Panics(t, func() {
			linkedList.AddItemBefore(NewNode(util.CreateRandInt(&seed)), util.CreateRandInt(&seed))
		})

		initialItems_Tests(t, linkedList, items)
	})

	t.Run(" Verify Node that already exists in another collection", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed), util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		linkedList.AddItemLast(items[1])

		tempDList = New[int]()
		tempDList.Clear()
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		require.Panics(t, func() { linkedList.AddItemBefore(tempDList.Last(), util.CreateRandInt(&seed)) })
		initialItems_Tests(t, linkedList, items)

	})
}

func TestAddBefore_LLNode_LLNode(t *testing.T) {

	var tempItems, headItems, headItemsReverse, tailItems, tailItemsReverse []int
	var linkedList *DList[int]
	arraySize := 16
	seed := int64(8293)
	headItems, tailItems, headItemsReverse, tailItemsReverse = util.CreateIntListData(arraySize, &seed)

	t.Run("Verify value is default(T)", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddNodeBefore(linkedList.First(), NewNode(defaultT[int]()))
		initialItems_Tests(t, linkedList, []int{defaultT[int](), headItems[0]})
	})

	t.Run("Node is the Head", func(t *testing.T) {

		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeBefore(linkedList.First(), NewNode(headItems[i]))
		}

		initialItems_Tests(t, linkedList, headItemsReverse)
	})

	t.Run("Node is the Tail", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		tempItems = make([]int, len(headItems))
		util.PartialCopy(headItems, 1, tempItems, 0, len(headItems)-1)
		tempItems[len(tempItems)-1] = headItems[0]
		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeBefore(linkedList.Last(), NewNode(headItems[i]))
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Node is after the Head", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		tempItems = make([]int, len(headItems))
		copy(tempItems, headItems)
		util.ReverseSubset(tempItems, 1, len(headItems)-1)

		for i := 2; i < arraySize; i++ {
			linkedList.AddNodeBefore(linkedList.First().Next(), NewNode(headItems[i]))
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Node is before the Tail", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])

		tempItems = make([]int, len(headItems))
		util.PartialCopy(headItems, 2, tempItems, 0, len(headItems)-2)
		tempItems[len(tempItems)-2] = headItems[0]
		tempItems[len(tempItems)-1] = headItems[1]

		for i := 2; i < arraySize; i++ {
			linkedList.AddNodeBefore(linkedList.Last().Previous(), NewNode(headItems[i]))
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Node is somewhere in the middle", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		linkedList.AddItemLast(headItems[2])

		tempItems = make([]int, len(headItems))
		util.PartialCopy(headItems, 3, tempItems, 0, len(headItems)-3)
		tempItems[len(tempItems)-3] = headItems[0]
		tempItems[len(tempItems)-2] = headItems[1]
		tempItems[len(tempItems)-1] = headItems[2]

		for i := 3; i < arraySize; i++ {
			linkedList.AddNodeBefore(linkedList.Last().Previous().Previous(), NewNode(headItems[i]))
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Call AddBefore several times remove some of the items", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeBefore(linkedList.First(), NewNode(headItems[i]))
		}

		for i := 0; i < arraySize; i++ {
			linkedList.RemoveFirst()
		}

		linkedList.AddItemFirst(tailItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeBefore(linkedList.First(), NewNode(tailItems[i]))
		}

		initialItems_Tests(t, linkedList, tailItemsReverse)
	})

	t.Run("Call AddBefore several times then call Clear", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeBefore(linkedList.First(), NewNode(headItems[i]))
		}

		linkedList.Clear()

		linkedList.AddItemFirst(tailItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeBefore(linkedList.First(), NewNode(tailItems[i]))
		}

		initialItems_Tests(t, linkedList, tailItemsReverse)
	})

	t.Run("Mix AddBefore and AddAfter calls", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemLast(headItems[0])
		linkedList.AddItemLast(tailItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeBefore(linkedList.First(), NewNode(headItems[i]))
			linkedList.AddNodeAfter(linkedList.Last(), NewNode(tailItems[i]))
		}

		tempItems = make([]int, len(headItemsReverse)+len(tailItems))
		copy(tempItems, headItemsReverse)
		util.PartialCopy(tailItems, 0, tempItems, len(headItemsReverse), len(tailItems))

		initialItems_Tests(t, linkedList, tempItems)
	})
}

func TestAddBefore_LLNode_LLNode_Negative(t *testing.T) {
	var (
		linkedList, tempDList *DList[int]
		items                 []int
	)

	seed := int64(8293)

	t.Run("Verify Null node", func(t *testing.T) {
		linkedList = New[int]()
		require.Panics(t, func() { linkedList.AddNodeBefore(nil, NewNode(util.CreateRandInt(&seed))) })
		initialItems_Tests(t, linkedList, []int{})
	})

	t.Run("Verify Node that is a new Node", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		require.Panics(t, func() {
			linkedList.AddNodeBefore(NewNode(util.CreateRandInt(&seed)), NewNode(util.CreateRandInt(&seed)))
		})

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
		require.Panics(t, func() { linkedList.AddNodeBefore(tempDList.Last(), NewNode(util.CreateRandInt(&seed))) })
		initialItems_Tests(t, linkedList, items)
	})

	t.Run("Verify Null newNode", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		require.Panics(t, func() { linkedList.AddNodeBefore(linkedList.First(), nil) })
		initialItems_Tests(t, linkedList, items)
	})

	t.Run("Verify newNode that already exists in this collection", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed), util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		linkedList.AddItemLast(items[1])
		require.Panics(t, func() { linkedList.AddNodeBefore(linkedList.First(), linkedList.Last()) })

		initialItems_Tests(t, linkedList, items)
	})

	t.Run("Verify newNode that already exists in another collection", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed), util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		linkedList.AddItemLast(items[1])

		tempDList = New[int]()
		tempDList.Clear()
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		require.Panics(t, func() { linkedList.AddNodeBefore(linkedList.First(), tempDList.Last()) })

		initialItems_Tests(t, linkedList, items)
	})
}
