package queue

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/fireflycons/generic_collections/lists/dlist"
	"github.com/stretchr/testify/require"
)

func TestEnqueuedDequeue(t *testing.T) {

	var queueItems []int
	var queue *Queue[int]
	seed := int64(2163)
	//arraySize := util.DefaultCapacity
	queueItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Enqueued items are dequeued", func(t *testing.T) {
		queue = New[int]()

		for _, v := range queueItems {
			queue.Enqueue(v)
		}

		verifyQueueState(t, queue, queueItems)

		for _, v := range queueItems {
			require.Equal(t, v, queue.Dequeue())
		}

		verifyQueueState(t, queue, []int{})
	})

	t.Run("Add is same as enqueue", func(t *testing.T) {
		queue = New[int]()

		for _, v := range queueItems {
			queue.Add(v)
		}

		verifyQueueState(t, queue, queueItems)

		for _, v := range queueItems {
			require.Equal(t, v, queue.Dequeue())
		}

		verifyQueueState(t, queue, []int{})
	})

	t.Run("Remove some items, re-add them", func(t *testing.T) {

		expectedItems := removeAndReAdd(&queue, queueItems)
		verifyQueueState(t, queue, expectedItems)
	})

	t.Run("Force queue to grow", func(t *testing.T) {
		queue = New(WithCapacity[int](10))

		for _, v := range queueItems {
			queue.Enqueue(v)
		}

		verifyQueueState(t, queue, queueItems)

		for _, v := range queueItems {
			require.Equal(t, v, queue.Dequeue())
		}

		verifyQueueState(t, queue, []int{})
	})

	t.Run("Peek queue with value returns value", func(t *testing.T) {
		queue := New[int]()
		expected := util.CreateRandInt(&seed)
		queue.Enqueue(expected)
		actual := queue.Peek()
		require.Equal(t, expected, actual)
		require.Equal(t, queue.Count(), 1)
	})
}

func TestClear(t *testing.T) {

	t.Run("Clear empties the queue", func(t *testing.T) {
		var queueItems []int
		var queue *Queue[int]
		seed := int64(2163)

		queueItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)
		queue = New[int]()

		for _, v := range queueItems {
			queue.Enqueue(v)
		}

		verifyQueueState(t, queue, queueItems)

		queue.Clear()
		verifyQueueState(t, queue, []int{})
	})

	t.Run("Clear with struct type", func(t *testing.T) {
		type strct struct {
			ptr *int
		}

		queue := New(WithComparer(func(s1, s2 strct) int {
			return 0
		}))

		i := 1
		queue.Enqueue(strct{ptr: &i})
		queue.Clear()
		verifyQueueState(t, queue, []strct{})

	})
}

func TestToSlice(t *testing.T) {

	var queueItems, tempItems, additionalItems []int
	var queue *Queue[int]
	seed := int64(2163)
	arraySize := util.DefaultCapacity
	queueItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)
	additionalArraySize := 4
	additionalItems, _, _, _ = util.CreateIntListData(additionalArraySize, &seed)

	t.Run("Empty queue", func(t *testing.T) {

		queue = New[int]()

		tempItems = queue.ToSlice()
		require.Equal(t, tempItems, []int{})
	})

	t.Run("Enqueued only", func(t *testing.T) {
		queue = New[int]()

		for _, v := range queueItems {
			queue.Enqueue(v)
		}

		tempItems = queue.ToSlice()
		tempItems2 := make([]int, 0, arraySize)

		for i := 0; i < arraySize; i++ {
			tempItems2 = append(tempItems2, queue.Dequeue())
		}

		require.Equal(t, tempItems2, tempItems)
	})

	t.Run("Move head and tail pointers", func(t *testing.T) {

		queue = New[int]()
		moveElems := 4
		for i := 0; i < arraySize; i++ {
			queue.Enqueue(queueItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := queue.Dequeue()
			tempItems = append(tempItems, v)
			queue.Enqueue(v)
		}

		slc := queue.ToSlice()
		tempItems2 := make([]int, arraySize)
		copy(tempItems2, queueItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		require.Equal(t, slc, tempItems2)
	})

	t.Run("Move head and tail pointers then add more items", func(t *testing.T) {

		queue = New[int]()
		moveElems := 4
		for i := 0; i < arraySize; i++ {
			queue.Enqueue(queueItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := queue.Dequeue()
			tempItems = append(tempItems, v)
			queue.Enqueue(v)
		}

		for i := 0; i < additionalArraySize; i++ {
			queue.Enqueue(additionalItems[i])
		}

		slc := queue.ToSlice()
		tempItems2 := make([]int, arraySize+additionalArraySize)
		copy(tempItems2, queueItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		util.PartialCopy(additionalItems, 0, tempItems2, arraySize, additionalArraySize)
		require.Equal(t, slc, tempItems2)
	})
}

func TestQueueTryOperations(t *testing.T) {

	seed := int64(2163)

	t.Run("TryDequeue empty queue returns false", func(t *testing.T) {

		queue := New[int]()
		_, ok := queue.TryDequeue()
		require.False(t, ok)
	})

	t.Run("TryPeek empty queue returns false", func(t *testing.T) {

		queue := New[int]()
		_, ok := queue.TryPeek()
		require.False(t, ok)
	})

	t.Run("TryPeek queue with value returns true and value", func(t *testing.T) {
		queue := New[int]()
		expected := util.CreateRandInt(&seed)
		queue.Enqueue(expected)
		actual, ok := queue.TryPeek()
		require.True(t, ok)
		require.Equal(t, expected, actual)
		require.Equal(t, queue.Count(), 1)
	})

	t.Run("TryDequeue queue with value returns true and value", func(t *testing.T) {
		queue := New[int]()
		expected := util.CreateRandInt(&seed)
		queue.Enqueue(expected)
		actual, ok := queue.TryDequeue()
		require.True(t, ok)
		require.Equal(t, expected, actual)
		require.Equal(t, queue.Count(), 0)
	})
}

func TestContains(t *testing.T) {

	var queueItems []int
	seed := int64(2163)
	queueItems = util.CreateSerialIntListData(util.DefaultCapacity, &seed)

	var queue = New[int]()
	queue.AddRange(queueItems)
	verifyQueueState(t, queue, queueItems)

	require.True(t, queue.Contains(queueItems[4]))
}

func TestRemove(t *testing.T) {

	var queueItems, tempItems, additionalItems []int
	var queue *Queue[int]
	seed := int64(2163)
	arraySize := util.DefaultCapacity
	queueItems = util.CreateSerialIntListData(util.DefaultCapacity, &seed)
	additionalArraySize := 4
	additionalItems = util.CreateSerialIntListData(additionalArraySize, &seed)

	_ = additionalItems
	t.Run("Empty queue", func(t *testing.T) {

		queue = New[int]()

		require.False(t, queue.Remove(queueItems[0]))
	})

	t.Run("Enqueued only", func(t *testing.T) {
		queue = New[int]()

		for _, v := range queueItems {
			queue.Enqueue(v)
		}

		verifyQueueState(t, queue, queueItems)

		removeItem := 4
		require.True(t, queue.Remove(queueItems[removeItem]))

		expectedItems := make([]int, len(queueItems)-1)
		index := 0
		for i := 0; i < removeItem; i++ {
			expectedItems[index] = queueItems[i]
			index++
		}
		for i := removeItem + 1; i < len(queueItems); i++ {
			expectedItems[index] = queueItems[i]
			index++
		}

		verifyQueueState(t, queue, expectedItems)
	})

	t.Run("Move head and tail pointers", func(t *testing.T) {

		queue = New[int]()
		moveElems := 4
		for i := 0; i < arraySize; i++ {
			queue.Enqueue(queueItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := queue.Dequeue()
			tempItems = append(tempItems, v)
			queue.Enqueue(v)
		}

		tempItems2 := make([]int, arraySize)
		copy(tempItems2, queueItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		verifyQueueState(t, queue, tempItems2)

		removeItem := 4
		require.True(t, queue.Remove(tempItems2[removeItem]))

		expectedItems := make([]int, len(tempItems2)-1)
		index := 0
		for i := 0; i < removeItem; i++ {
			expectedItems[index] = tempItems2[i]
			index++
		}
		for i := removeItem + 1; i < len(tempItems2); i++ {
			expectedItems[index] = tempItems2[i]
			index++
		}

		verifyQueueState(t, queue, expectedItems)
	})

	t.Run("Move head and tail pointers then add more items", func(t *testing.T) {

		queue = New[int]()
		moveElems := 4
		for i := 0; i < arraySize; i++ {
			queue.Enqueue(queueItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := queue.Dequeue()
			tempItems = append(tempItems, v)
			queue.Enqueue(v)
		}

		for i := 0; i < additionalArraySize; i++ {
			queue.Enqueue(additionalItems[i])
		}

		tempItems2 := make([]int, arraySize+additionalArraySize)
		copy(tempItems2, queueItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		util.PartialCopy(additionalItems, 0, tempItems2, arraySize, additionalArraySize)
		verifyQueueState(t, queue, tempItems2)

		removeItem := 4
		require.True(t, queue.Remove(tempItems2[removeItem]))

		expectedItems := make([]int, len(tempItems2)-1)
		index := 0
		for i := 0; i < removeItem; i++ {
			expectedItems[index] = tempItems2[i]
			index++
		}
		for i := removeItem + 1; i < len(tempItems2); i++ {
			expectedItems[index] = tempItems2[i]
			index++
		}

		verifyQueueState(t, queue, expectedItems)
	})

	t.Run("Head gt tail and delete ahead of head", func(t *testing.T) {

		queue = New[int]()
		moveElems := 4
		dequeueElems := 2

		for i := 0; i < arraySize; i++ {
			queue.Enqueue(queueItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := queue.Dequeue()
			tempItems = append(tempItems, v)
			queue.Enqueue(v)
		}

		tempItems2 := make([]int, arraySize)
		copy(tempItems2, queueItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		verifyQueueState(t, queue, tempItems2)

		for i := 0; i < dequeueElems; i++ {
			queue.Dequeue()
		}

		tempItems3 := make([]int, len(tempItems2)-dequeueElems)
		util.PartialCopy(tempItems2, dequeueElems, tempItems3, 0, len(tempItems2)-dequeueElems)
		verifyQueueState(t, queue, tempItems3)

		removeItem := 4
		require.True(t, queue.Remove(tempItems3[removeItem]))

		expectedItems := make([]int, len(tempItems3)-1)
		index := 0
		for i := 0; i < removeItem; i++ {
			expectedItems[index] = tempItems3[i]
			index++
		}
		for i := removeItem + 1; i < len(tempItems3); i++ {
			expectedItems[index] = tempItems3[i]
			index++
		}

		verifyQueueState(t, queue, expectedItems)
	})

	t.Run("Head gt tail and delete behind tail", func(t *testing.T) {

		queue = New[int]()
		moveElems := 4
		dequeueElems := 2

		for i := 0; i < arraySize; i++ {
			queue.Enqueue(queueItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := queue.Dequeue()
			tempItems = append(tempItems, v)
			queue.Enqueue(v)
		}

		tempItems2 := make([]int, arraySize)
		copy(tempItems2, queueItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		verifyQueueState(t, queue, tempItems2)

		for i := 0; i < dequeueElems; i++ {
			queue.Dequeue()
		}

		tempItems3 := make([]int, len(tempItems2)-dequeueElems)
		util.PartialCopy(tempItems2, dequeueElems, tempItems3, 0, len(tempItems2)-dequeueElems)
		verifyQueueState(t, queue, tempItems3)

		removeItem := 12
		require.True(t, queue.Remove(tempItems3[removeItem]))

		expectedItems := make([]int, len(tempItems3)-1)
		index := 0
		for i := 0; i < removeItem; i++ {
			expectedItems[index] = tempItems3[i]
			index++
		}
		for i := removeItem + 1; i < len(tempItems3); i++ {
			expectedItems[index] = tempItems3[i]
			index++
		}

		verifyQueueState(t, queue, expectedItems)
	})

	t.Run("Head gt tail and delete head", func(t *testing.T) {

		queue = New[int]()
		moveElems := 4
		dequeueElems := 2

		for i := 0; i < arraySize; i++ {
			queue.Enqueue(queueItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := queue.Dequeue()
			tempItems = append(tempItems, v)
			queue.Enqueue(v)
		}

		tempItems2 := make([]int, arraySize)
		copy(tempItems2, queueItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		verifyQueueState(t, queue, tempItems2)

		for i := 0; i < dequeueElems; i++ {
			queue.Dequeue()
		}

		tempItems3 := make([]int, len(tempItems2)-dequeueElems)
		util.PartialCopy(tempItems2, dequeueElems, tempItems3, 0, len(tempItems2)-dequeueElems)
		verifyQueueState(t, queue, tempItems3)

		headValue := tempItems3[0]
		require.True(t, queue.Remove(headValue))

		expectedItems := tempItems3[1:]
		verifyQueueState(t, queue, expectedItems)
	})

	t.Run("Head gt tail and delete behind tail using custom comparer", func(t *testing.T) {

		queue = New(WithComparer(func(a, b int) int {
			if a == b {
				return 0
			}

			if a > b {
				return 1
			}

			return -1
		}))

		moveElems := 4
		dequeueElems := 2

		for i := 0; i < arraySize; i++ {
			queue.Enqueue(queueItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := queue.Dequeue()
			tempItems = append(tempItems, v)
			queue.Enqueue(v)
		}

		tempItems2 := make([]int, arraySize)
		copy(tempItems2, queueItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		verifyQueueState(t, queue, tempItems2)

		for i := 0; i < dequeueElems; i++ {
			queue.Dequeue()
		}

		tempItems3 := make([]int, len(tempItems2)-dequeueElems)
		util.PartialCopy(tempItems2, dequeueElems, tempItems3, 0, len(tempItems2)-dequeueElems)
		verifyQueueState(t, queue, tempItems3)

		removeItem := 12
		require.True(t, queue.Remove(tempItems3[removeItem]))

		expectedItems := make([]int, len(tempItems3)-1)
		index := 0
		for i := 0; i < removeItem; i++ {
			expectedItems[index] = tempItems3[i]
			index++
		}
		for i := removeItem + 1; i < len(tempItems3); i++ {
			expectedItems[index] = tempItems3[i]
			index++
		}

		verifyQueueState(t, queue, expectedItems)
	})

	t.Run("Remove item not in queue returns false", func(t *testing.T) {
		queue = New[int]()

		for i := 0; i < arraySize; i++ {
			queue.Enqueue(queueItems[i])
		}

		removeItem := -6
		require.True(t, util.IndexOf(queueItems, removeItem, queue.compare, false) == -1)
		require.False(t, queue.Remove(removeItem))
	})
}

func TestQueue_Negative(t *testing.T) {

	t.Run("Dequeue on empty queue panics", func(t *testing.T) {

		queue := New[int]()
		require.Panics(t, func() { queue.Dequeue() })
	})

	t.Run("Peek on empty queue panics", func(t *testing.T) {

		queue := New[int]()
		require.Panics(t, func() { queue.Peek() })
	})
}

func TestQueueAddRange(t *testing.T) {

	var additionalItems, queueItems, tempItems []int
	var queue *Queue[int]
	arraySize := 16
	seed := int64(21543)
	queueItems, _, _, _ = util.CreateIntListData(arraySize, &seed)
	additionalArraySize := 4
	additionalItems, _, _, _ = util.CreateIntListData(additionalArraySize, &seed)

	t.Run("Empty queue", func(t *testing.T) {
		queue = New[int]()
		queue.AddRange(queueItems)
		verifyQueueState(t, queue, queueItems)
	})

	t.Run("Populated queue", func(t *testing.T) {
		queue = New[int]()
		for i := 0; i < arraySize; i++ {
			queue.Enqueue(queueItems[i])
		}

		queue.AddRange(additionalItems)
		tempItems = make([]int, arraySize+additionalArraySize)
		copy(tempItems, queueItems)
		index := len(queueItems)
		for i := 0; i < additionalArraySize; i++ {
			tempItems[index] = additionalItems[i]
			index++
		}
		verifyQueueState(t, queue, tempItems)
	})

	t.Run("Move head and tail pointers", func(t *testing.T) {

		queue = New[int]()
		moveElems := 4
		for i := 0; i < arraySize; i++ {
			queue.Enqueue(queueItems[i])
		}

		tempItems = make([]int, moveElems)

		for i := 0; i < moveElems; i++ {
			v := queue.Dequeue()
			tempItems[i] = v
			queue.Enqueue(v)
		}

		queue.AddRange(additionalItems)

		tempItems2 := make([]int, arraySize+additionalArraySize)
		copy(tempItems2, queueItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		util.PartialCopy(additionalItems, 0, tempItems2, arraySize, additionalArraySize)
		verifyQueueState(t, queue, tempItems2)
	})

	t.Run("Full Queue", func(t *testing.T) {

		queue = New[int]()
		for i := 0; i < arraySize; i++ {
			queue.Enqueue(queueItems[i])
		}

		// Will cause reallocation
		queue.AddRange(additionalItems)

		tempItems = make([]int, arraySize+additionalArraySize)
		copy(tempItems, queueItems)
		util.PartialCopy(additionalItems, 0, tempItems, arraySize, additionalArraySize)
		verifyQueueState(t, queue, tempItems)
	})

	t.Run("Queue with gap in - some elements removed", func(t *testing.T) {
		tempItems = createGappedQueue(&queue, queueItems)
		verifyQueueState(t, queue, tempItems)

		queue.AddRange(additionalItems)
		expectedItems := make([]int, len(tempItems)+len(additionalItems))
		copy(expectedItems, tempItems)
		util.PartialCopy(additionalItems, 0, expectedItems, len(tempItems), len(additionalItems))
		verifyQueueState(t, queue, expectedItems)
	})

	t.Run("From collection", func(t *testing.T) {
		linkedList := dlist.New[int]()
		linkedList.AddRange(queueItems)
		queue := New[int]()
		queue.AddCollection(linkedList)
		verifyQueueState(t, queue, queueItems)

	})
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

	t.Run("Parallel Enqueue", func(t *testing.T) {
		q := New(WithThreadSafe[int]())
		wg := sync.WaitGroup{}
		wg.Add(2)

		enqueueFunc := func(q1 *Queue[int], slc []int, w *sync.WaitGroup) {
			for _, v := range slc {
				q1.Enqueue(v)
			}
			w.Done()
		}

		go enqueueFunc(q, items1, &wg)
		go enqueueFunc(q, items2, &wg)
		wg.Wait()
		require.ElementsMatch(t, itemsCombined, q.ToSlice())
	})

	t.Run("Parallel Dequeue", func(t *testing.T) {
		q := New(WithThreadSafe[int]())
		q.AddRange(itemsCombined)
		popped1 := make([]int, itemsPerThread)
		popped2 := make([]int, itemsPerThread)
		wg := sync.WaitGroup{}
		wg.Add(2)

		dequeueFunc := func(q1 *Queue[int], slc []int, w *sync.WaitGroup) {
			for i := 0; i < itemsPerThread; i++ {
				slc[i] = q1.Dequeue()
			}
			w.Done()
		}

		go dequeueFunc(q, popped1, &wg)
		go dequeueFunc(q, popped2, &wg)
		wg.Wait()

		popCombined := make([]int, len(itemsCombined))
		copy(popCombined, popped1)
		copy(popCombined[len(items1):], popped2)
		require.ElementsMatch(t, itemsCombined, popCombined)
	})
}

func verifyQueueState[T any](t *testing.T, queue *Queue[T], expectedItems []T) {
	require.Equal(t, len(expectedItems), queue.Count())
	require.Equal(t, expectedItems, queue.ToSlice())

	items := make([]T, 0, queue.size)
	q1 := queue.makeDeepCopy()

	for i := 0; i < queue.size; i++ {
		items = append(items, q1.Dequeue())
	}

	require.Equal(t, expectedItems, items)
}

func benchmarkEnqueue(s *Queue[int], data []int) {
	for _, v := range data {
		s.Enqueue(v)
	}
}

func benchmarkDequeue(s *Queue[int], nitems int) {
	for i := 0; i < nitems; i++ {
		s.Dequeue()
	}
}

func BenchmarkQueue(b *testing.B) {

	seed := int64(2163)
	data := make(map[int][]int, 4)
	elements := []int{100, 1000, 10000, 100000}
	list := util.CreateSingleIntListData(100000, &seed)

	for _, elem := range elements {
		data[elem] = list[:elem]
	}

	var q *Queue[int]

	for z := 0; z <= 1; z++ {
		threadsafe := z == 1

		for _, elems := range elements {
			b.Run(fmt.Sprintf("Queue-Enqueue-%d-%s-NoPresize-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					if threadsafe {
						q = New(WithThreadSafe[int]())
					} else {
						q = New[int]()
					}
					b.StartTimer()
					benchmarkEnqueue(q, data[elems])
				}
			})
		}

		for _, elems := range elements {
			b.Run(fmt.Sprintf("Queue-Enqueue-%d-%s-Presize-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					if threadsafe {
						q = New[int](WithCapacity[int](elems), WithThreadSafe[int]())
					} else {
						q = New[int](WithCapacity[int](elems))
					}
					b.StartTimer()
					benchmarkEnqueue(q, data[elems])
				}
			})
		}

		for _, elems := range elements {
			b.Run(fmt.Sprintf("Queue-Dequeue-%d-%s-NA-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					if threadsafe {
						q = New[int](WithCapacity[int](elems), WithThreadSafe[int]())
					} else {
						q = New[int](WithCapacity[int](elems))
					}
					q.AddRange(data[elems])
					b.StartTimer()
					benchmarkDequeue(q, elems)
				}
			})
		}
	}

	for _, elems := range elements {
		b.Run(fmt.Sprintf("Queue-Sort-%d-NA-NA-NA", elems), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				q = New[int](WithCapacity[int](elems))
				q.AddRange(data[elems])
				b.StartTimer()
				q.Sort()
			}
		})
	}

	for _, elems := range elements {
		b.Run(fmt.Sprintf("Queue-Contains-%d-NA-NA-NA", elems), func(b *testing.B) {
			q = New[int](WithCapacity[int](elems))
			q.AddRange(data[elems])
			lookup := make([]int, elems)
			copy(lookup, data[elems])
			rand.Shuffle(elems, func(i, j int) {
				lookup[i], lookup[j] = lookup[j], lookup[i]
			})

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				q.Contains(lookup[i%elems])
			}
		})
	}

	for _, elems := range elements {
		q = New[int](WithCapacity[int](elems))
		q.AddRange(data[elems])
		b.ResetTimer()
		b.Run(fmt.Sprintf("Queue-Min-%d-NA-NA-NoConcurrent", elems), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				q.Min()
			}
		})

		b.Run(fmt.Sprintf("Queue-Max-%d-NA-NA-NoConcurrent", elems), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				q.Max()
			}
		})
	}

	elems := 100000
	q = New[int](WithCapacity[int](elems), WithConcurrent[int]())
	q.AddRange(data[elems])
	b.ResetTimer()
	b.Run(fmt.Sprintf("Queue-Min-%d-NA-NA-Concurrent", elems), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			q.Min()
		}
	})

	b.Run(fmt.Sprintf("Queue-Max-%d-NA-NA-Concurrent", elems), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			q.Max()
		}
	})
}
