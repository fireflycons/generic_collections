package dlist

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestLLConstructor(t *testing.T) {

	t.Run("Simple construction", func(t *testing.T) {
		linkedList := New[int]()

		require.Nil(t, linkedList.head)
		require.Empty(t, linkedList.count)
	})

	t.Run("Construct with nil comparer panics", func(t *testing.T) {
		var comp functions.ComparerFunc[int]

		require.Panics(t, func() { New(WithComparer(comp)) })
	})
}

func TestOperationOnNilCollectionPanics(t *testing.T) {
	var linkedList *DList[int]
	require.Panics(t, func() { linkedList.Count() })
}

func TestUnsafe(t *testing.T) {

	t.Run("GetVersion", func(t *testing.T) {
		s := New[int]()

		for i := 0; i < 10; i++ {
			s.Add(1)
			require.Equal(t, s.version, util.GetVersion[int](s))
		}
	})

	t.Run("GetLock", func(t *testing.T) {
		s := New[int]()

		ptrMutex := util.GetLock[int](s)

		require.Same(t, s.lock, ptrMutex)
	})
}

func TestThreadSafety(t *testing.T) {

	seed := int64(2163)
	itemsPerThread := 1024
	items1 := util.CreateSingleIntListData(itemsPerThread, &seed)
	items2 := util.CreateSingleIntListData(itemsPerThread, &seed)
	itemsCombined := make([]int, len(items1)+len(items2))
	copy(itemsCombined, items1)
	copy(itemsCombined[len(items1):], items2)

	t.Run("Parallel Add", func(t *testing.T) {
		ll := New(WithThreadSafe[int]())
		wg := sync.WaitGroup{}
		wg.Add(2)

		addFunc := func(set *DList[int], slc []int, w *sync.WaitGroup) {
			for _, v := range slc {
				set.Add(v)
			}
			w.Done()
		}

		go addFunc(ll, items1, &wg)
		go addFunc(ll, items2, &wg)
		wg.Wait()

		require.ElementsMatch(t, itemsCombined, ll.ToSlice())
	})

	t.Run("Parallel Remove", func(t *testing.T) {
		ll := New(WithThreadSafe[int]())
		ll.AddRange(itemsCombined)
		wg := sync.WaitGroup{}
		wg.Add(2)

		removeFunc := func(stk *DList[int], slc []int, w *sync.WaitGroup) {
			for _, v := range slc {
				stk.Remove(v)
			}
			w.Done()
		}

		go removeFunc(ll, items1, &wg)
		go removeFunc(ll, items2, &wg)
		wg.Wait()

		require.Equal(t, 0, ll.Count())
	})
}

func benchmarkAdd(s *DList[int], data []int) {
	for _, v := range data {
		s.Add(v)
	}
}

func benchmarkRemove(s *DList[int], data []int) {
	for _, v := range data {
		s.Remove(v)
	}
}

func BenchmarkDList(b *testing.B) {

	seed := int64(2163)
	data := make(map[int][]int, 4)
	elements := []int{100, 1000, 10000, 100000}
	list := util.CreateSingleIntListData(100000, &seed)

	for _, elem := range elements {
		data[elem] = list[:elem]
	}

	var ll *DList[int]

	for z := 0; z <= 1; z++ {
		threadsafe := z == 1

		for _, elems := range elements {
			b.Run(fmt.Sprintf("List-Add-%d-%s-NA-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					if threadsafe {
						ll = New(WithThreadSafe[int]())
					} else {
						ll = New[int]()
					}
					b.StartTimer()
					benchmarkAdd(ll, data[elems])
				}
			})
		}

		for _, elems := range elements {
			b.Run(fmt.Sprintf("List-Remove-%d-%s-NA-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					if threadsafe {
						ll = New[int](WithThreadSafe[int]())
					} else {
						ll = New[int]()
					}
					ll.AddRange(data[elems])
					b.StartTimer()
					benchmarkRemove(ll, data[elems])
				}
			})
		}
	}
	for _, elems := range elements {
		b.Run(fmt.Sprintf("List-Sort-%d-NA-NA-NA", elems), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				ll = New[int]()
				ll.AddRange(data[elems])
				b.StartTimer()
				ll.Sort()
			}
		})
	}

	for _, elems := range elements {
		b.Run(fmt.Sprintf("List-Contains-%d-NA-NA-NA", elems), func(b *testing.B) {
			ll = New[int]()
			ll.AddRange(data[elems])
			lookup := make([]int, elems)
			copy(lookup, data[elems])
			rand.Shuffle(elems, func(i, j int) {
				lookup[i], lookup[j] = lookup[j], lookup[i]
			})

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ll.Contains(lookup[i%elems])
			}
		})
	}

	for _, elems := range elements {
		ll = New[int]()
		ll.AddRange(data[elems])
		b.ResetTimer()
		b.Run(fmt.Sprintf("List-Min-%d-NA-NA-NA", elems), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ll.Min()
			}
		})

		b.Run(fmt.Sprintf("List-Max-%d-NA-NA-NA", elems), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ll.Max()
			}
		})
	}

}

func defaultT[T any]() T {
	var result T
	return result
}

// Verify the state (integrity) of the list being tested.
func verifyLLState[T any](t *testing.T, linkedList *DList[T], expectedItems []T) {

	var previousNode, currentNode, nextNode *DListNode[T]

	expectedItemsLength := len(expectedItems)

	// Verify count
	require.Equal(t, linkedList.count, expectedItemsLength)

	// Verify head/tail
	if expectedItemsLength == 0 {
		require.Nil(t, linkedList.First())
		require.Nil(t, linkedList.Last())
	} else if expectedItemsLength == 1 {
		verifyDListNodeA(t, linkedList.First(), expectedItems[0], linkedList, nil, nil)
		verifyDListNodeA(t, linkedList.Last(), expectedItems[0], linkedList, nil, nil)
	} else {
		verifyDListNodeB(t, linkedList.First(), expectedItems[0], linkedList, true, false)
		verifyDListNodeB(t, linkedList.Last(), expectedItems[expectedItemsLength-1], linkedList, false, true)
	}

	// Moving forward through the collection starting at head
	previousNode = nil
	currentNode = linkedList.First()

	for i := 0; currentNode != nil; i++ {
		nextNode = currentNode.Next()
		verifyDListNodeA(t, currentNode, expectedItems[i], linkedList, previousNode, nextNode)
		previousNode = currentNode
		currentNode = currentNode.Next()
	}

	// Moving backward through he collection starting at Tail
	currentNode = linkedList.Last()
	nextNode = nil

	for i := 0; currentNode != nil; i++ {
		previousNode = currentNode.Previous()
		verifyDListNodeA(t, currentNode, expectedItems[expectedItemsLength-1-i], linkedList, previousNode, nextNode)
		nextNode = currentNode
		currentNode = currentNode.Previous()
	}

	// Verify Contains
	for i := 0; i < expectedItemsLength; i++ {
		require.True(t, linkedList.Contains(expectedItems[i]))
	}

	// Verify ToSlice()
	slc := linkedList.ToSlice()

	require.Equal(t, expectedItems, slc)
}

func verifyFind[T any](t *testing.T, linkedList *DList[T], expectedItems []T) {
	var previousNode, currentNode, nextNode *DListNode[T]

	currentNode = nil

	for i := 0; i < len(expectedItems); i++ {
		previousNode = currentNode
		currentNode = linkedList.findNode(expectedItems[i], forward)
		nextNode = currentNode.Next()
		verifyDListNodeA(t, currentNode, expectedItems[i], linkedList, previousNode, nextNode)
	}

	currentNode = nil

	for i := len(expectedItems) - 1; i >= 0; i-- {
		nextNode = currentNode
		currentNode = linkedList.findNode(expectedItems[i], forward)
		previousNode = currentNode.Previous()
		verifyDListNodeA(t, currentNode, expectedItems[i], linkedList, previousNode, nextNode)
	}
}

func verifyFindDuplicates[T any](t *testing.T, linkedList *DList[T], expectedItems []T) {
	var previousNode, currentNode, nextNode *DListNode[T]

	nodes := make([]*DListNode[T], len(expectedItems))
	index := 0

	currentNode = linkedList.First()

	for currentNode != nil {
		nodes[index] = currentNode
		currentNode = currentNode.Next()
		index++
	}

	for i := 0; i < len(expectedItems); i++ {
		currentNode = linkedList.findNode(expectedItems[i], forward)
		index := util.IndexOf(expectedItems, expectedItems[i], linkedList.compare, false)

		if index > 0 {
			previousNode = nodes[index-1]
		} else {
			previousNode = nil
		}

		if len(nodes)-1 > index {
			nextNode = nodes[index+1]
		} else {
			nextNode = nil
		}

		require.Equal(t, nodes[index], currentNode)
		verifyDListNodeA(t, currentNode, expectedItems[i], linkedList, previousNode, nextNode)
	}
}

func verifyFindLast[T any](t *testing.T, linkedList *DList[T], expectedItems []T) {
	var previousNode, currentNode, nextNode *DListNode[T]

	currentNode = nil

	for i := 0; i < len(expectedItems); i++ {
		previousNode = currentNode
		currentNode = linkedList.findNode(expectedItems[i], reverse)
		nextNode = currentNode.Next()
		verifyDListNodeA(t, currentNode, expectedItems[i], linkedList, previousNode, nextNode)
	}

	currentNode = nil

	for i := len(expectedItems) - 1; i >= 0; i-- {
		nextNode = currentNode
		currentNode = linkedList.findNode(expectedItems[i], reverse)
		previousNode = currentNode.Previous()
		verifyDListNodeA(t, currentNode, expectedItems[i], linkedList, previousNode, nextNode)
	}
}

func verifyFindLastDuplicates[T any](t *testing.T, linkedList *DList[T], expectedItems []T) {
	var previousNode, currentNode, nextNode *DListNode[T]

	nodes := make([]*DListNode[T], len(expectedItems))
	index := 0

	currentNode = linkedList.First()

	for currentNode != nil {
		nodes[index] = currentNode
		currentNode = currentNode.Next()
		index++
	}

	for i := 0; i < len(expectedItems); i++ {
		currentNode = linkedList.findNode(expectedItems[i], reverse)
		index = util.LastIndexOf(expectedItems, expectedItems[i], linkedList.compare, false)

		if index > 0 {
			previousNode = nodes[index-1]
		} else {
			previousNode = nil
		}

		if len(nodes)-1 > index {
			nextNode = nodes[index+1]
		} else {
			nextNode = nil
		}

		require.Equal(t, nodes[index], currentNode)
		verifyDListNodeA(t, currentNode, expectedItems[i], linkedList, previousNode, nextNode)
	}
}

func verifyDListNodeA[T any](t *testing.T, node *DListNode[T], expectedValue T, expectedList *DList[T], expectedPrevious, expectedNext *DListNode[T]) {

	require.Equal(t, expectedValue, node.Value())
	require.Equal(t, expectedList, node.List())
	require.Equal(t, expectedPrevious, node.Previous())
	require.Equal(t, expectedNext, node.Next())
}

func verifyDListNodeB[T any](t *testing.T, node *DListNode[T], expectedValue T, expectedList *DList[T], expectedPreviousNil, expectedNextNil bool) {

	require.Equal(t, expectedValue, node.Value())
	require.Equal(t, expectedList, node.List())

	if expectedPreviousNil {
		require.Nil(t, node.Previous())
	} else {
		require.NotNil(t, node.Previous())
	}

	if expectedNextNil {
		require.Nil(t, node.Next())
	} else {
		require.NotNil(t, node.Next())
	}
}

func verifyRemovedNode2[T any](t *testing.T, node *DListNode[T], expectedValue T) {
	var def T
	var headNode, tailNode *DListNode[T]
	tempDList := New[T]()

	tempDList.AddItemLast(def)
	tempDList.AddItemLast(def)
	headNode = tempDList.First()
	tailNode = tempDList.Last()

	require.Nil(t, node.List())
	require.Nil(t, node.Previous())
	require.Nil(t, node.Next())
	require.Equal(t, expectedValue, node.Value())

	tempDList.AddNodeAfter(tempDList.First(), node)

	require.Equal(t, tempDList, node.List())
	require.Equal(t, headNode, node.Previous())
	require.Equal(t, tailNode, node.Next())
	require.Equal(t, expectedValue, node.Value())

	initialItems_Tests(t, tempDList, []T{def, expectedValue, def})
}

func verifyRemovedNode3[T any](t *testing.T, linkedList *DList[T], node *DListNode[T], expectedValue T) {
	tailNode := linkedList.Last()

	require.Nil(t, node.List())
	require.Nil(t, node.Previous())
	require.Nil(t, node.Next())
	require.Equal(t, expectedValue, node.Value())

	linkedList.AddNodeLast(node)
	require.Equal(t, linkedList, node.List())
	require.Equal(t, tailNode, node.Previous())
	require.Nil(t, node.Next())

	linkedList.RemoveLast()
}

func verifyRemovedNode4[T any](t *testing.T, linkedList *DList[T], linkedListValues []T, node *DListNode[T], expectedValue T) {
	tailNode := linkedList.Last()

	require.Nil(t, node.List())
	require.Nil(t, node.Previous())
	require.Nil(t, node.Next())
	require.Equal(t, expectedValue, node.Value())

	linkedList.AddNodeLast(node)
	require.Equal(t, linkedList, node.List())
	require.Equal(t, tailNode, node.Previous())
	require.Nil(t, node.Next())
	require.Equal(t, expectedValue, node.Value())

	expected := make([]T, len(linkedListValues)+1)
	copy(expected, linkedListValues)
	expected[len(linkedListValues)] = expectedValue

	initialItems_Tests(t, linkedList, expected)
	linkedList.RemoveLast()
}

func initialItems_Tests[T any](t *testing.T, collection *DList[T], expectedItems []T) {
	verifyLLState(t, collection, expectedItems)
}
