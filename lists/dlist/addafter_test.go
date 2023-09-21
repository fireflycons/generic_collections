package dlist

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestAddAfter_LLNode(t *testing.T) {

	var tempItems, headItems, headItemsReverse, tailItems []int
	var linkedList *DList[int]
	arraySize := 16
	seed := int64(8293)
	headItems, tailItems, headItemsReverse, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Verify value is default(T)", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemAfter(linkedList.First(), defaultT[int]())
		initialItems_Tests(t, linkedList, []int{headItems[0], defaultT[int]()})

	})

	t.Run("Node is the head", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])

		tempItems = make([]int, len(headItems))
		copy(tempItems, headItems)
		util.ReverseSubset(tempItems, 1, len(headItems)-1)

		for i := 1; i < arraySize; i++ {
			linkedList.AddItemAfter(linkedList.First(), headItems[i])
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Node is the tail", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])

		for i := 1; i < arraySize; i++ {
			linkedList.AddItemAfter(linkedList.Last(), headItems[i])
		}

		initialItems_Tests(t, linkedList, headItems)
	})

	t.Run("Node is after the head", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])

		tempItems = make([]int, len(headItems))
		copy(tempItems, headItems)
		util.ReverseSubset(tempItems, 2, len(headItems)-2)

		for i := 2; i < arraySize; i++ {
			linkedList.AddItemAfter(linkedList.First().Next(), headItems[i])
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Node is before the tail", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])

		tempItems = make([]int, len(headItems))
		util.PartialCopy(headItems, 2, tempItems, 1, len(headItems)-2)
		tempItems[0] = headItems[0]
		tempItems[len(tempItems)-1] = headItems[1]

		for i := 2; i < arraySize; i++ {
			linkedList.AddItemAfter(linkedList.Last().Previous(), headItems[i])
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Node is somewhere in the middle", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		linkedList.AddItemLast(headItems[2])

		tempItems = make([]int, len(headItems))
		util.PartialCopy(headItems, 3, tempItems, 1, len(headItems)-3)

		tempItems[0] = headItems[0]
		tempItems[len(tempItems)-2] = headItems[1]
		tempItems[len(tempItems)-1] = headItems[2]

		for i := 3; i < arraySize; i++ {
			linkedList.AddItemAfter(linkedList.Last().Previous().Previous(), headItems[i])
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Call AddAfter several times remove some of the items", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddItemAfter(linkedList.Last(), headItems[i])
		}

		linkedList.RemoveItem(headItems[2])
		linkedList.RemoveItem(headItems[len(headItems)-3])
		linkedList.RemoveItem(headItems[1])
		linkedList.RemoveItem(headItems[len(headItems)-2])
		linkedList.RemoveFirst()
		linkedList.RemoveLast()

		//With the above remove we should have removed the first and last 3 items
		tempItems = make([]int, len(headItems)-6)
		util.PartialCopy(headItems, 3, tempItems, 0, len(headItems)-6)

		initialItems_Tests(t, linkedList, tempItems)

		for i := 0; i < arraySize; i++ {
			linkedList.AddItemAfter(linkedList.Last(), tailItems[i])
		}

		tempItems2 := make([]int, len(tempItems)+len(tailItems))
		copy(tempItems2, tempItems)
		util.PartialCopy(tailItems, 0, tempItems2, len(tempItems), len(tailItems))

		initialItems_Tests(t, linkedList, tempItems2)
	})

	t.Run("Call AddAfter several times remove all of the items", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])

		for i := 1; i < arraySize; i++ {
			linkedList.AddItemAfter(linkedList.Last(), headItems[i])
		}

		for i := 0; i < arraySize; i++ {
			linkedList.RemoveFirst()
		}

		linkedList.AddItemFirst(tailItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddItemAfter(linkedList.Last(), tailItems[i])
		}

		initialItems_Tests(t, linkedList, tailItems)
	})

	t.Run("Call AddAfter several times then call Clear", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])

		for i := 1; i < arraySize; i++ {
			linkedList.AddItemAfter(linkedList.Last(), headItems[i])
		}

		linkedList.Clear()

		linkedList.AddItemFirst(tailItems[0])

		for i := 1; i < arraySize; i++ {
			linkedList.AddItemAfter(linkedList.Last(), tailItems[i])
		}

		initialItems_Tests(t, linkedList, tailItems)
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
}

func TestAddAfter_LLNode_Negative(t *testing.T) {
	var (
		linkedList *DList[int]
		items      []int
	)

	seed := int64(8293)

	t.Run("Verify Null node", func(t *testing.T) {
		linkedList = New[int]()
		require.Panics(t, func() { linkedList.AddItemAfter(nil, util.CreateRandInt(&seed)) })
		initialItems_Tests(t, linkedList, []int{})
	})

	t.Run("Verify Node that is a new Node", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		require.Panics(t, func() {
			linkedList.AddItemAfter(NewNode(util.CreateRandInt(&seed)), util.CreateRandInt(&seed))
		})
		initialItems_Tests(t, linkedList, items)
	})

	t.Run("Verify Node that already exists in another collection", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed), util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		linkedList.AddItemLast(items[1])

		tempDList := New[int]()
		tempDList.Clear()
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		require.Panics(t, func() { linkedList.AddItemAfter(tempDList.Last(), util.CreateRandInt(&seed)) })
		initialItems_Tests(t, linkedList, items)

	})
}

func TestAddAfter_LLNode_LLNode(t *testing.T) {

	var (
		linkedList                                                          *DList[int]
		tempItems, headItems, headItemsReverse, tailItems, tailItemsReverse []int
	)

	arraySize := 16
	seed := int64(8293)

	headItems = make([]int, arraySize)
	tailItems = make([]int, arraySize)
	headItemsReverse = make([]int, arraySize)
	tailItemsReverse = make([]int, arraySize)

	for i := 0; i < arraySize; i++ {
		index := (arraySize - 1) - i
		head := util.CreateRandInt(&seed)
		tail := util.CreateRandInt(&seed)
		headItems[i] = head
		headItemsReverse[index] = head
		tailItems[i] = tail
		tailItemsReverse[index] = tail
	}

	t.Run("Verify value is default(T)", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		defaultValue := defaultT[int]()
		linkedList.AddNodeAfter(linkedList.First(), NewNode(defaultValue))
		initialItems_Tests(t, linkedList, []int{headItems[0], defaultValue})
	})

	t.Run("Node is the Head", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])

		tempItems = make([]int, len(headItems))
		copy(tempItems, headItems)
		util.ReverseSubset(tempItems, 1, len(headItems)-1)

		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeAfter(linkedList.First(), NewNode(headItems[i]))
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Node is the Tail", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeAfter(linkedList.Last(), NewNode(headItems[i]))
		}

		initialItems_Tests(t, linkedList, headItems)
	})

	t.Run("Node is after the Head", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		tempItems = make([]int, len(headItems))
		copy(tempItems, headItems)
		util.ReverseSubset(tempItems, 2, len(headItems)-2)
		for i := 2; i < arraySize; i++ {
			linkedList.AddNodeAfter(linkedList.First().Next(), NewNode(headItems[i]))
		}
		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Node is before the Tail", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		tempItems = make([]int, len(headItems))
		util.PartialCopy(headItems, 2, tempItems, 1, len(headItems)-2)
		tempItems[0] = headItems[0]
		tempItems[len(tempItems)-1] = headItems[1]

		for i := 2; i < arraySize; i++ {
			linkedList.AddNodeAfter(linkedList.Last().Previous(), NewNode(headItems[i]))
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Node is somewhere in the middle", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		linkedList.AddItemLast(headItems[2])

		tempItems = make([]int, len(headItems))
		util.PartialCopy(headItems, 3, tempItems, 1, len(headItems)-3)
		tempItems[0] = headItems[0]
		tempItems[len(tempItems)-2] = headItems[1]
		tempItems[len(tempItems)-1] = headItems[2]

		for i := 3; i < arraySize; i++ {
			linkedList.AddNodeAfter(linkedList.Last().Previous().Previous(), NewNode(headItems[i]))
		}

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Call AddAfter several times remove some of the items", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeAfter(linkedList.Last(), NewNode(headItems[i]))
		}

		linkedList.RemoveItem(headItems[2])
		linkedList.RemoveItem(headItems[len(headItems)-3])
		linkedList.RemoveItem(headItems[1])
		linkedList.RemoveItem(headItems[len(headItems)-2])
		linkedList.RemoveFirst()
		linkedList.RemoveLast()
		//With the above remove we should have removed the first and last 3 items
		tempItems = make([]int, len(headItems)-6)
		util.PartialCopy(headItems, 3, tempItems, 0, len(headItems)-6)

		initialItems_Tests(t, linkedList, tempItems)

		for i := 0; i < arraySize; i++ {
			linkedList.AddNodeAfter(linkedList.Last(), NewNode(tailItems[i]))
		}

		tempItems2 := make([]int, len(tempItems)+len(tailItems))
		copy(tempItems2, tempItems)
		util.PartialCopy(tailItems, 0, tempItems2, len(tempItems), len(tailItems))

		initialItems_Tests(t, linkedList, tempItems2)
	})

	t.Run("Call AddAfter several times remove all of the items", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeAfter(linkedList.Last(), NewNode(headItems[i]))
		}

		for i := 0; i < arraySize; i++ {
			linkedList.RemoveFirst()
		}

		linkedList.AddItemFirst(tailItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeAfter(linkedList.Last(), NewNode(tailItems[i]))
		}

		initialItems_Tests(t, linkedList, tailItems)
	})

	t.Run("Call AddAfter several times then call Clear", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeAfter(linkedList.Last(), NewNode(headItems[i]))
		}

		linkedList.Clear()

		linkedList.AddItemFirst(tailItems[0])
		for i := 1; i < arraySize; i++ {
			linkedList.AddNodeAfter(linkedList.Last(), NewNode(tailItems[i]))
		}

		initialItems_Tests(t, linkedList, tailItems)
	})

	t.Run("Mix AddBefore and AddAfter calls", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemLast(headItems[0])
		linkedList.AddItemLast(tailItems[0])
		for i := 1; i < arraySize; i++ {
			if i&1 == 0 {
				linkedList.AddItemBefore(linkedList.First(), headItems[i])
				linkedList.AddItemAfter(linkedList.Last(), tailItems[i])
			} else {
				linkedList.AddNodeBefore(linkedList.First(), NewNode(headItems[i]))
				linkedList.AddNodeAfter(linkedList.Last(), NewNode(tailItems[i]))
			}
		}

		tempItems = make([]int, len(headItemsReverse)+len(tailItems))
		copy(tempItems, headItemsReverse)
		util.PartialCopy(tailItems, 0, tempItems, len(headItemsReverse), len(tailItems))
		initialItems_Tests(t, linkedList, tempItems)
	})
}

func TestAddAfter_LLNode_LLNode_Negative(t *testing.T) {
	var (
		linkedList, tempDList *DList[int]
		items                 []int
	)

	seed := int64(8293)

	t.Run("Verify Null node", func(t *testing.T) {
		linkedList = New[int]()
		require.Panics(t, func() { linkedList.AddNodeAfter(nil, NewNode(util.CreateRandInt(&seed))) })
		initialItems_Tests(t, linkedList, []int{})
	})

	t.Run("Verify Node that is a new Node", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		require.Panics(t, func() {
			linkedList.AddNodeAfter(NewNode(util.CreateRandInt(&seed)), NewNode(util.CreateRandInt(&seed)))
		})

		initialItems_Tests(t, linkedList, items)
	})

	t.Run("Verify Node that already exists in another collection", func(t *testing.T) {
		linkedList = New[int]()
		tempDList = New[int]()
		items = []int{util.CreateRandInt(&seed), util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		linkedList.AddItemLast(items[1])

		tempDList.Clear()
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		require.Panics(t, func() { linkedList.AddNodeAfter(tempDList.Last(), NewNode(util.CreateRandInt(&seed))) })

		initialItems_Tests(t, linkedList, items)
	})

	t.Run("Verify Null newNode", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		require.Panics(t, func() { linkedList.AddNodeAfter(linkedList.First(), nil) })

		initialItems_Tests(t, linkedList, items)
	})

	t.Run("Verify newNode that already exists in this collection", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed), util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		linkedList.AddItemLast(items[1])
		require.Panics(t, func() { linkedList.AddNodeAfter(linkedList.First(), linkedList.Last()) })

		initialItems_Tests(t, linkedList, items)
	})

	t.Run("Verify newNode that already exists in another collection", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed), util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		linkedList.AddItemLast(items[1])

		tempDList.Clear()
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		require.Panics(t, func() { linkedList.AddNodeAfter(linkedList.First(), tempDList.Last()) })

		initialItems_Tests(t, linkedList, items)
	})
}
