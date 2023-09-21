package dlist

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestRemove_LLNode(t *testing.T) {
	var tempItems, headItems []int
	var linkedList *DList[int]
	var tempNode1, tempNode2, tempNode3 *DListNode[int]
	var def int
	arraySize := 16
	seed := int64(21543)
	headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Call Remove with an item that exists in the collection size=1", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemLast(headItems[0])
		tempNode1 = linkedList.First()

		linkedList.RemoveNode(linkedList.First()) //Remove when  VS Whidbey: 234648 is resolved

		verifyRemovedNode4(t, linkedList, []int{}, tempNode1, headItems[0])
		initialItems_Tests(t, linkedList, []int{})
	})

	t.Run("Call Remove with the Head collection size=2", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		tempNode1 = linkedList.First()

		linkedList.RemoveNode(linkedList.First()) //Remove when  VS Whidbey: 234648 is resolved

		initialItems_Tests(t, linkedList, []int{headItems[1]})
		verifyRemovedNode2(t, tempNode1, headItems[0])
	})

	t.Run("Call Remove with the Tail collection size=2", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		tempNode1 = linkedList.Last()

		linkedList.RemoveNode(linkedList.Last()) //Remove when  VS Whidbey: 234648 is resolved

		initialItems_Tests(t, linkedList, []int{headItems[0]})
		verifyRemovedNode2(t, tempNode1, headItems[1])
	})

	t.Run("Call Remove all the items collection size=2", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		tempNode1 = linkedList.First()
		tempNode2 = linkedList.Last()

		linkedList.RemoveNode(linkedList.First()) //Remove when  VS Whidbey: 234648 is resolved
		linkedList.RemoveNode(linkedList.Last())  //Remove when  VS Whidbey: 234648 is resolved

		initialItems_Tests(t, linkedList, []int{})
		verifyRemovedNode4(t, linkedList, []int{}, tempNode1, headItems[0])
		verifyRemovedNode4(t, linkedList, []int{}, tempNode2, headItems[1])
	})

	t.Run("Call Remove with the Head collection size=3", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		linkedList.AddItemLast(headItems[2])
		tempNode1 = linkedList.First()

		linkedList.RemoveNode(linkedList.First()) //Remove when  VS Whidbey: 234648 is resolved
		initialItems_Tests(t, linkedList, []int{headItems[1], headItems[2]})
		verifyRemovedNode4(t, linkedList, []int{headItems[1], headItems[2]}, tempNode1, headItems[0])
	})

	t.Run("Call Remove with the middle item collection size=3", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		linkedList.AddItemLast(headItems[2])
		tempNode1 = linkedList.First().Next()

		linkedList.RemoveNode(linkedList.First().Next()) //Remove when  VS Whidbey: 234648 is resolved
		initialItems_Tests(t, linkedList, []int{headItems[0], headItems[2]})
		verifyRemovedNode4(t, linkedList, []int{headItems[0], headItems[2]}, tempNode1, headItems[1])
	})

	t.Run("Call Remove with the Tail collection size=3", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		linkedList.AddItemLast(headItems[2])
		tempNode1 = linkedList.Last()

		linkedList.RemoveNode(linkedList.Last()) //Remove when  VS Whidbey: 234648 is resolved
		initialItems_Tests(t, linkedList, []int{headItems[0], headItems[1]})
		verifyRemovedNode4(t, linkedList, []int{headItems[0], headItems[1]}, tempNode1, headItems[2])
	})

	t.Run("Call Remove all the items collection size=3", func(t *testing.T) {
		linkedList = New[int]()
		linkedList.AddItemFirst(headItems[0])
		linkedList.AddItemLast(headItems[1])
		linkedList.AddItemLast(headItems[2])
		tempNode1 = linkedList.First()
		tempNode2 = linkedList.First().Next()
		tempNode3 = linkedList.Last()

		linkedList.RemoveNode(linkedList.First().Next().Next()) //Remove when  VS Whidbey: 234648 is resolved
		linkedList.RemoveNode(linkedList.First().Next())        //Remove when  VS Whidbey: 234648 is resolved
		linkedList.RemoveNode(linkedList.First())               //Remove when  VS Whidbey: 234648 is resolved

		initialItems_Tests(t, linkedList, []int{})
		verifyRemovedNode2(t, tempNode1, headItems[0])
		verifyRemovedNode2(t, tempNode2, headItems[1])
		verifyRemovedNode2(t, tempNode3, headItems[2])
	})

	t.Run("Call Remove all the items starting with the first collection size=16", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemLast(headItems[i])
		}

		for i := 0; i < arraySize; i++ {
			linkedList.RemoveNode(linkedList.First()) //Remove when  VS Whidbey: 234648 is resolved
			startIndex := i + 1
			length := arraySize - i - 1
			expectedItems := make([]int, length)
			util.PartialCopy(headItems, startIndex, expectedItems, 0, length)
			initialItems_Tests(t, linkedList, expectedItems)
		}
	})

	t.Run("Call Remove all the items starting with the last collection size=16", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemLast(headItems[i])
		}

		for i := arraySize - 1; 0 <= i; i-- {
			linkedList.RemoveNode(linkedList.Last()) //Remove when  VS Whidbey: 234648 is resolved
			expectedItems := make([]int, i)
			copy(expectedItems, headItems)
			initialItems_Tests(t, linkedList, expectedItems)
		}
	})

	t.Run("Remove some items in the middle", func(t *testing.T) {
		linkedList = New[int]()
		for i := 0; i < arraySize; i++ {
			linkedList.AddItemFirst(headItems[i])
		}

		linkedList.RemoveNode(linkedList.First().Next().Next())        //Remove when  VS Whidbey: 234648 is resolved
		linkedList.RemoveNode(linkedList.Last().Previous().Previous()) //Remove when  VS Whidbey: 234648 is resolved
		linkedList.RemoveNode(linkedList.First().Next())               //Remove when  VS Whidbey: 234648 is resolved
		linkedList.RemoveNode(linkedList.Last().Previous())
		linkedList.RemoveNode(linkedList.First()) //Remove when  VS Whidbey: 234648 is resolved
		linkedList.RemoveNode(linkedList.Last())  //Remove when  VS Whidbey: 234648 is resolved

		//With the above remove we should have removed the first and last 3 items
		headItemsReverse := make([]int, arraySize)
		copy(headItemsReverse, headItems)
		util.Reverse(headItemsReverse)

		tempItems = make([]int, len(headItemsReverse)-6)
		util.PartialCopy(headItemsReverse, 3, tempItems, 0, len(headItemsReverse)-6)

		initialItems_Tests(t, linkedList, tempItems)
	})

	t.Run("Remove an item with a value of default(T)", func(t *testing.T) {
		linkedList = New[int]()

		for i := 0; i < arraySize; i++ {
			linkedList.AddItemLast(headItems[i])
		}

		linkedList.AddItemLast(def)

		linkedList.RemoveNode(linkedList.Last()) //Remove when  VS Whidbey: 234648 is resolved

		initialItems_Tests(t, linkedList, headItems)
	})
}

func TestRemove_Duplicates_LLNode(t *testing.T) {
	linkedList := New[int]()
	arraySize := 16
	seed := int64(21543)
	var items []int
	nodes := make([]*DListNode[int], arraySize*2)
	var currentNode *DListNode[int]
	var index int

	items = make([]int, arraySize)

	for i := 0; i < arraySize; i++ {
		items[i] = util.CreateRandInt(&seed)
	}

	for i := 0; i < arraySize; i++ {
		linkedList.AddItemLast(items[i])
	}

	for i := 0; i < arraySize; i++ {
		linkedList.AddItemLast(items[i])
	}

	currentNode = linkedList.First()
	index = 0

	for currentNode != nil {
		nodes[index] = currentNode
		currentNode = currentNode.Next()
		index++
	}

	linkedList.RemoveNode(linkedList.First().Next().Next())        //Remove when  VS Whidbey: 234648 is resolved
	linkedList.RemoveNode(linkedList.Last().Previous().Previous()) //Remove when  VS Whidbey: 234648 is resolved
	linkedList.RemoveNode(linkedList.First().Next())               //Remove when  VS Whidbey: 234648 is resolved
	linkedList.RemoveNode(linkedList.Last().Previous())
	linkedList.RemoveNode(linkedList.First()) //Remove when  VS Whidbey: 234648 is resolved
	linkedList.RemoveNode(linkedList.Last())  //Remove when  VS Whidbey: 234648 is resolved

	//[] Verify that the duplicates were removed from the beginning of the collection
	currentNode = linkedList.First()

	//Verify the duplicates that should have been removed
	for i := 3; i < len(nodes)-3; i++ {
		require.NotNil(t, currentNode)                             //"Err_48588ahid CurrentNode is null index=" + i
		require.Equal(t, currentNode, nodes[i])                    //"Err_5488ahid CurrentNode is not the expected node index=" + i
		require.Equal(t, items[i%len(items)], currentNode.Value()) //"Err_16588ajide CurrentNode value index=" + i

		currentNode = currentNode.Next()
	}

	require.Nil(t, currentNode) //"Err_30878ajid Expected CurrentNode to be null after moving through entire list"
}

func TestRemove_LLNode_Negative(t *testing.T) {

	var linkedList *DList[int]
	tempDList := New[int]()
	seed := int64(21543)
	var items []int

	t.Run("Verify Null node", func(t *testing.T) {
		linkedList = New[int]()
		require.Panics(t, func() { linkedList.RemoveNode(nil) }) //"Err_858ahia Expected null node to throws ArgumentNullException\n"

		initialItems_Tests(t, linkedList, []int{})
	})

	t.Run("Verify Node that is a new Node", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		require.Panics(t, func() { linkedList.RemoveNode(NewNode(util.CreateRandInt(&seed))) }) //"Err_0568ajods Expected Node that is a new Node throws InvalidOperationException\n"

		initialItems_Tests(t, linkedList, items)
	})

	t.Run("Verify Node that already exists in another collection", func(t *testing.T) {
		linkedList = New[int]()
		items = []int{util.CreateRandInt(&seed), util.CreateRandInt(&seed)}
		linkedList.AddItemLast(items[0])
		linkedList.AddItemLast(items[1])

		tempDList.Clear()
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		tempDList.AddItemLast(util.CreateRandInt(&seed))
		require.Panics(t, func() { linkedList.RemoveNode(tempDList.Last()) }) //"Err_98809ahied Node that already exists in another collection throws InvalidOperationException\n"

		initialItems_Tests(t, linkedList, items)
	})
}
