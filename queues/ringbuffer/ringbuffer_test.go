package ringbuffer

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/fireflycons/generic_collections/sets/orderedset"
	"github.com/stretchr/testify/require"
)

func TestEnqueueDequeue(t *testing.T) {

	var queueItems []int
	var ringBuffer *RingBuffer[int]
	seed := int64(2163)
	//arraySize := util.DefaultCapacity
	queueItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Enqueued items are dequeued", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)

		for _, v := range queueItems {
			ringBuffer.Enqueue(v)
		}

		verifyBufferState(t, ringBuffer, queueItems)

		for _, v := range queueItems {
			require.Equal(t, v, ringBuffer.Dequeue())
		}

		verifyBufferState(t, ringBuffer, []int{})
	})

	t.Run("Add is same as enqueue", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)

		for _, v := range queueItems {
			ringBuffer.Add(v)
		}

		verifyBufferState(t, ringBuffer, queueItems)

		for _, v := range queueItems {
			require.Equal(t, v, ringBuffer.Dequeue())
		}

		verifyBufferState(t, ringBuffer, []int{})
	})

	t.Run("Remove some items, re-add them", func(t *testing.T) {

		expectedItems := removeAndReAdd(&ringBuffer, queueItems)
		verifyBufferState(t, ringBuffer, expectedItems)
	})

	t.Run("Peek buffer with value returns value", func(t *testing.T) {
		ringBuffer := New[int](util.DefaultCapacity)
		expected := util.CreateRandInt(&seed)
		ringBuffer.Enqueue(expected)
		actual := ringBuffer.Peek()
		require.Equal(t, expected, actual)
		require.Equal(t, ringBuffer.Count(), 1)
	})

	t.Run("Enqueue on full buffer displaces head item", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)

		for _, v := range queueItems {
			ringBuffer.Enqueue(v)
		}

		verifyBufferState(t, ringBuffer, queueItems)
		require.True(t, ringBuffer.Full())

		newValue := util.CreateRandInt(&seed)
		ringBuffer.Enqueue(newValue)
		expectedItems := make([]int, util.DefaultCapacity)
		util.PartialCopy(queueItems, 1, expectedItems, 0, util.DefaultCapacity-1)
		expectedItems[util.DefaultCapacity-1] = newValue
		require.Equal(t, expectedItems[0], ringBuffer.Peek())
		verifyBufferState(t, ringBuffer, expectedItems)
	})
}

func TestOffer(t *testing.T) {

	t.Run("Offer should return false when buffer full", func(t *testing.T) {
		ringBuffer := New[int](1)
		require.True(t, ringBuffer.Offer(1))
		require.False(t, ringBuffer.Offer(2))
		verifyBufferState(t, ringBuffer, []int{1})
	})
}

func TestAddRange(t *testing.T) {

	var additionalItems, queueItems, tempItems []int
	var ringBuffer *RingBuffer[int]
	arraySize := util.DefaultCapacity
	seed := int64(21543)
	queueItems, _, _, _ = util.CreateIntListData(arraySize, &seed)
	additionalArraySize := 4
	additionalItems, _, _, _ = util.CreateIntListData(additionalArraySize, &seed)

	t.Run("Empty buffer", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)
		ringBuffer.AddRange(queueItems)
		verifyBufferState(t, ringBuffer, queueItems)
	})

	t.Run("Slice bigger than buffer only adds buf length items from end of slice", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)
		tempItems = make([]int, arraySize+additionalArraySize)
		copy(tempItems, queueItems)
		util.PartialCopy(additionalItems, 0, tempItems, arraySize, additionalArraySize)
		ringBuffer.AddRange(tempItems)
		expectedItems := make([]int, arraySize)
		util.PartialCopy(tempItems, additionalArraySize, expectedItems, 0, arraySize)
		verifyBufferState(t, ringBuffer, expectedItems)
	})

	t.Run("Empty buffer and slice less than max size", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)
		ringBuffer.AddRange(additionalItems)
		verifyBufferState(t, ringBuffer, additionalItems)
	})

	t.Run("Partially filled buffer and slice will exceed capacity", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)
		for _, v := range additionalItems {
			ringBuffer.Enqueue(v)
		}
		verifyBufferState(t, ringBuffer, additionalItems)

		// Slice same size as buffer
		ringBuffer.AddRange(queueItems)
		// so existing items displaced and head pointing
		// to first added item
		verifyBufferState(t, ringBuffer, queueItems)
	})
}

func TestAddCollection(t *testing.T) {

	seed := int64(21543)

	t.Run("Add larger OrderedSet than buffer results in items at end of set being enqueued", func(t *testing.T) {
		setItems, _, _, _ := util.CreateIntListData(util.DefaultCapacity, &seed)
		sortedSetData := make([]int, util.DefaultCapacity)
		bufSize := 4
		copy(sortedSetData, setItems)
		sort.Ints(sortedSetData)
		ss := orderedset.New[int]()
		ss.AddRange(setItems)
		ringBuffer := New[int](bufSize)
		ringBuffer.AddCollection(ss)
		expected := sortedSetData[util.DefaultCapacity-bufSize:]
		verifyBufferState(t, ringBuffer, expected)
	})
}

func TestClear(t *testing.T) {

	t.Run("Clear empties the buffer", func(t *testing.T) {
		var queueItems []int
		var ringBuffer *RingBuffer[int]
		seed := int64(2163)

		queueItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)
		ringBuffer = New[int](util.DefaultCapacity)

		for _, v := range queueItems {
			ringBuffer.Enqueue(v)
		}

		verifyBufferState(t, ringBuffer, queueItems)

		ringBuffer.Clear()
		verifyBufferState(t, ringBuffer, []int{})
	})
}

func TestContains(t *testing.T) {

	var additionalItems, queueItems []int
	var ringBuffer *RingBuffer[int]
	arraySize := util.DefaultCapacity
	seed := int64(21543)
	queueItems = util.CreateSerialIntListData(arraySize, &seed)
	additionalArraySize := 4
	additionalItems = util.CreateSerialIntListData(additionalArraySize, &seed)

	t.Run("Partially filled queue", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)
		ringBuffer.AddRange(additionalItems)

		for _, v := range additionalItems {
			require.True(t, ringBuffer.Contains(v))
		}
	})

	t.Run("Full queue", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)
		ringBuffer.AddRange(queueItems)

		for _, v := range queueItems {
			require.True(t, ringBuffer.Contains(v))
		}
	})
}

func TestToSlice(t *testing.T) {

	var queueItems, tempItems, additionalItems []int
	var ringBuffer *RingBuffer[int]
	seed := int64(2163)
	arraySize := util.DefaultCapacity
	queueItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)
	additionalArraySize := 4
	additionalItems, _, _, _ = util.CreateIntListData(additionalArraySize, &seed)

	t.Run("Empty buffer", func(t *testing.T) {

		ringBuffer = New[int](util.DefaultCapacity)

		tempItems = ringBuffer.ToSlice()
		require.Equal(t, tempItems, []int{})
	})

	t.Run("Enqueued only", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)

		for _, v := range queueItems {
			ringBuffer.Enqueue(v)
		}

		tempItems = ringBuffer.ToSlice()
		tempItems2 := make([]int, 0, arraySize)

		for i := 0; i < arraySize; i++ {
			tempItems2 = append(tempItems2, ringBuffer.Dequeue())
		}

		require.Equal(t, tempItems2, tempItems)
	})

	t.Run("Move head and tail pointers", func(t *testing.T) {

		ringBuffer = New[int](util.DefaultCapacity)
		moveElems := 4
		for i := 0; i < arraySize; i++ {
			ringBuffer.Enqueue(queueItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := ringBuffer.Dequeue()
			tempItems = append(tempItems, v)
			ringBuffer.Enqueue(v)
		}

		slc := ringBuffer.ToSlice()
		tempItems2 := make([]int, arraySize)
		copy(tempItems2, queueItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		require.Equal(t, slc, tempItems2)
	})

	t.Run("Move head and tail pointers then add more items", func(t *testing.T) {

		ringBuffer = New[int](util.DefaultCapacity + additionalArraySize)
		moveElems := 4
		for i := 0; i < arraySize; i++ {
			ringBuffer.Enqueue(queueItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := ringBuffer.Dequeue()
			tempItems = append(tempItems, v)
			ringBuffer.Enqueue(v)
		}

		for i := 0; i < additionalArraySize; i++ {
			ringBuffer.Enqueue(additionalItems[i])
		}

		slc := ringBuffer.ToSlice()
		tempItems2 := make([]int, arraySize+additionalArraySize)
		copy(tempItems2, queueItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		util.PartialCopy(additionalItems, 0, tempItems2, arraySize, additionalArraySize)
		require.Equal(t, slc, tempItems2)
	})
}

func TestTryOperations(t *testing.T) {

	seed := int64(2163)

	t.Run("TryDequeue empty buffer returns false", func(t *testing.T) {

		ringBuffer := New[int](util.DefaultCapacity)
		_, ok := ringBuffer.TryDequeue()
		require.False(t, ok)
	})

	t.Run("TryPeek empty buffer returns false", func(t *testing.T) {

		ringBuffer := New[int](util.DefaultCapacity)
		_, ok := ringBuffer.TryPeek()
		require.False(t, ok)
	})

	t.Run("TryPeek buffer with value returns true and value", func(t *testing.T) {
		ringBuffer := New[int](util.DefaultCapacity)
		expected := util.CreateRandInt(&seed)
		ringBuffer.Enqueue(expected)
		actual, ok := ringBuffer.TryPeek()
		require.True(t, ok)
		require.Equal(t, expected, actual)
		require.Equal(t, ringBuffer.Count(), 1)
	})

	t.Run("TryDequeue buffer with value returns true and value", func(t *testing.T) {
		ringBuffer := New[int](util.DefaultCapacity)
		expected := util.CreateRandInt(&seed)
		ringBuffer.Enqueue(expected)
		actual, ok := ringBuffer.TryDequeue()
		require.True(t, ok)
		require.Equal(t, expected, actual)
		require.Equal(t, ringBuffer.Count(), 0)
	})
}

func TestRemove(t *testing.T) {

	var bufferItems, tempItems, additionalItems []int
	var ringBuffer *RingBuffer[int]
	seed := int64(2163)
	arraySize := util.DefaultCapacity
	bufferItems = util.CreateSerialIntListData(util.DefaultCapacity, &seed)
	additionalArraySize := 4
	additionalItems = util.CreateSerialIntListData(additionalArraySize, &seed)

	_ = additionalItems
	t.Run("Empty queue", func(t *testing.T) {

		ringBuffer = New[int](util.DefaultCapacity)

		require.False(t, ringBuffer.Remove(bufferItems[0]))
	})

	t.Run("Enqueued only", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)

		for _, v := range bufferItems {
			ringBuffer.Enqueue(v)
		}

		verifyBufferState(t, ringBuffer, bufferItems)

		removeItem := 4
		require.True(t, ringBuffer.Remove(bufferItems[removeItem]))

		expectedItems := make([]int, len(bufferItems)-1)
		index := 0
		for i := 0; i < removeItem; i++ {
			expectedItems[index] = bufferItems[i]
			index++
		}
		for i := removeItem + 1; i < len(bufferItems); i++ {
			expectedItems[index] = bufferItems[i]
			index++
		}

		verifyBufferState(t, ringBuffer, expectedItems)
	})

	t.Run("Move head and tail pointers", func(t *testing.T) {

		ringBuffer = New[int](util.DefaultCapacity)
		moveElems := 4
		for i := 0; i < arraySize; i++ {
			ringBuffer.Enqueue(bufferItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := ringBuffer.Dequeue()
			tempItems = append(tempItems, v)
			ringBuffer.Enqueue(v)
		}

		tempItems2 := make([]int, arraySize)
		copy(tempItems2, bufferItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		verifyBufferState(t, ringBuffer, tempItems2)

		removeItem := 4
		require.True(t, ringBuffer.Remove(tempItems2[removeItem]))

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

		verifyBufferState(t, ringBuffer, expectedItems)
	})

	t.Run("Head gt tail and delete ahead of head", func(t *testing.T) {

		ringBuffer = New[int](util.DefaultCapacity)
		moveElems := 4
		dequeueElems := 2

		for i := 0; i < arraySize; i++ {
			ringBuffer.Enqueue(bufferItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := ringBuffer.Dequeue()
			tempItems = append(tempItems, v)
			ringBuffer.Enqueue(v)
		}

		tempItems2 := make([]int, arraySize)
		copy(tempItems2, bufferItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		verifyBufferState(t, ringBuffer, tempItems2)

		for i := 0; i < dequeueElems; i++ {
			ringBuffer.Dequeue()
		}

		tempItems3 := make([]int, len(tempItems2)-dequeueElems)
		util.PartialCopy(tempItems2, dequeueElems, tempItems3, 0, len(tempItems2)-dequeueElems)
		verifyBufferState(t, ringBuffer, tempItems3)

		removeItem := 4
		require.True(t, ringBuffer.Remove(tempItems3[removeItem]))

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

		verifyBufferState(t, ringBuffer, expectedItems)
	})

	t.Run("Head gt tail and delete behind tail", func(t *testing.T) {

		ringBuffer = New[int](util.DefaultCapacity)
		moveElems := 4
		dequeueElems := 2

		for i := 0; i < arraySize; i++ {
			ringBuffer.Enqueue(bufferItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := ringBuffer.Dequeue()
			tempItems = append(tempItems, v)
			ringBuffer.Enqueue(v)
		}

		tempItems2 := make([]int, arraySize)
		copy(tempItems2, bufferItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		verifyBufferState(t, ringBuffer, tempItems2)

		for i := 0; i < dequeueElems; i++ {
			ringBuffer.Dequeue()
		}

		tempItems3 := make([]int, len(tempItems2)-dequeueElems)
		util.PartialCopy(tempItems2, dequeueElems, tempItems3, 0, len(tempItems2)-dequeueElems)
		verifyBufferState(t, ringBuffer, tempItems3)

		removeItem := 12
		require.True(t, ringBuffer.Remove(tempItems3[removeItem]))

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

		verifyBufferState(t, ringBuffer, expectedItems)
	})

	t.Run("Head gt tail and delete head", func(t *testing.T) {

		ringBuffer = New[int](util.DefaultCapacity)
		moveElems := 4
		dequeueElems := 2

		for i := 0; i < arraySize; i++ {
			ringBuffer.Enqueue(bufferItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := ringBuffer.Dequeue()
			tempItems = append(tempItems, v)
			ringBuffer.Enqueue(v)
		}

		tempItems2 := make([]int, arraySize)
		copy(tempItems2, bufferItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		verifyBufferState(t, ringBuffer, tempItems2)

		for i := 0; i < dequeueElems; i++ {
			ringBuffer.Dequeue()
		}

		tempItems3 := make([]int, len(tempItems2)-dequeueElems)
		util.PartialCopy(tempItems2, dequeueElems, tempItems3, 0, len(tempItems2)-dequeueElems)
		verifyBufferState(t, ringBuffer, tempItems3)

		headValue := tempItems3[0]
		require.True(t, ringBuffer.Remove(headValue))

		expectedItems := tempItems3[1:]
		verifyBufferState(t, ringBuffer, expectedItems)
	})

	t.Run("Head gt tail and delete behind tail using custom comparer", func(t *testing.T) {

		ringBuffer = New(util.DefaultCapacity, WithComparer(func(a, b int) int {
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
			ringBuffer.Enqueue(bufferItems[i])
		}

		tempItems = make([]int, 0, moveElems)

		for i := 0; i < moveElems; i++ {
			v := ringBuffer.Dequeue()
			tempItems = append(tempItems, v)
			ringBuffer.Enqueue(v)
		}

		tempItems2 := make([]int, arraySize)
		copy(tempItems2, bufferItems[moveElems:])
		util.PartialCopy(tempItems, 0, tempItems2, arraySize-moveElems, moveElems)
		verifyBufferState(t, ringBuffer, tempItems2)

		for i := 0; i < dequeueElems; i++ {
			ringBuffer.Dequeue()
		}

		tempItems3 := make([]int, len(tempItems2)-dequeueElems)
		util.PartialCopy(tempItems2, dequeueElems, tempItems3, 0, len(tempItems2)-dequeueElems)
		verifyBufferState(t, ringBuffer, tempItems3)

		removeItem := 12
		require.True(t, ringBuffer.Remove(tempItems3[removeItem]))

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

		verifyBufferState(t, ringBuffer, expectedItems)
	})

	t.Run("Remove item not in queue returns false", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)

		for i := 0; i < arraySize; i++ {
			ringBuffer.Enqueue(bufferItems[i])
		}

		removeItem := -6
		require.True(t, util.IndexOf(bufferItems, removeItem, ringBuffer.compare, false) == -1)
		require.False(t, ringBuffer.Remove(removeItem))
	})
}

func TestBuffer_Negative(t *testing.T) {

	t.Run("Dequeue on empty buffer panics", func(t *testing.T) {

		queue := New[int](util.DefaultCapacity)
		require.Panics(t, func() { queue.Dequeue() })
	})

	t.Run("Peek on empty buffer panics", func(t *testing.T) {

		queue := New[int](util.DefaultCapacity)
		require.Panics(t, func() { queue.Peek() })
	})
}

func TestUnsafe(t *testing.T) {

	t.Run("GetVersion", func(t *testing.T) {
		s := New[int](util.DefaultCapacity)

		for i := 0; i < 10; i++ {
			s.Add(1)
			require.Equal(t, s.version, util.GetVersion[int](s))
		}
	})

	t.Run("GetLock", func(t *testing.T) {
		s := New[int](util.DefaultCapacity)

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
		buf := New(len(itemsCombined), WithThreadSafe[int]())
		wg := sync.WaitGroup{}
		wg.Add(2)

		enqueueFunc := func(buf1 *RingBuffer[int], slc []int, w *sync.WaitGroup) {
			for _, v := range slc {
				buf1.Enqueue(v)
			}
			w.Done()
		}

		go enqueueFunc(buf, items1, &wg)
		go enqueueFunc(buf, items2, &wg)
		wg.Wait()
		require.ElementsMatch(t, itemsCombined, buf.ToSlice())
	})

	t.Run("Parallel Dequeue", func(t *testing.T) {
		q := New(len(itemsCombined), WithThreadSafe[int]())
		q.AddRange(itemsCombined)
		popped1 := make([]int, itemsPerThread)
		popped2 := make([]int, itemsPerThread)
		wg := sync.WaitGroup{}
		wg.Add(2)

		dequeueFunc := func(buf1 *RingBuffer[int], slc []int, w *sync.WaitGroup) {
			for i := 0; i < itemsPerThread; i++ {
				slc[i] = buf1.Dequeue()
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

func verifyBufferState[T any](t *testing.T, buf *RingBuffer[T], expectedItems []T) {
	require.Equal(t, len(expectedItems), buf.Count())
	require.Equal(t, expectedItems, buf.ToSlice())

	items := make([]T, 0, buf.size)
	q1 := buf.makeDeepCopy()

	for i := 0; i < buf.size; i++ {
		items = append(items, q1.Dequeue())
	}

	require.Equal(t, expectedItems, items)
}

func removeAndReAdd[T any](buf **RingBuffer[T], queueItems []T) []T {
	*buf = New[T](util.DefaultCapacity)
	moveElems := 4
	for _, v := range queueItems {
		(*buf).Enqueue(v)
	}

	tempItems := make([]T, 0, moveElems)

	for i := 0; i < moveElems; i++ {
		v := (*buf).Dequeue()
		tempItems = append(tempItems, v)
		(*buf).Enqueue(v)
	}

	expectedItems := make([]T, len(queueItems))
	copy(expectedItems, queueItems[moveElems:])
	util.PartialCopy(tempItems, 0, expectedItems, 12, moveElems)

	return expectedItems
}

func benchmarkEnqueue(s *RingBuffer[int], data []int) {
	for _, v := range data {
		s.Enqueue(v)
	}
}

func benchmarkDequeue(s *RingBuffer[int], nitems int) {
	for i := 0; i < nitems; i++ {
		s.Dequeue()
	}
}

func BenchmarkRingBuffer(b *testing.B) {

	seed := int64(2163)
	data := make(map[int][]int, 4)
	elements := []int{100, 1000, 10000, 100000}
	list := util.CreateSingleIntListData(100000, &seed)

	for _, elem := range elements {
		data[elem] = list[:elem]
	}

	var buf *RingBuffer[int]

	for z := 0; z <= 1; z++ {
		threadsafe := z == 1

		for _, elems := range elements {
			b.Run(fmt.Sprintf("Queue-Enqueue-%d-%s-NA-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					if threadsafe {
						buf = New(elems, WithThreadSafe[int]())
					} else {
						buf = New[int](elems)
					}
					b.StartTimer()
					benchmarkEnqueue(buf, data[elems])
				}
			})
		}

		for _, elems := range elements {
			b.Run(fmt.Sprintf("Queue-Dequeue-%d-%s-NA-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					if threadsafe {
						buf = New[int](elems, WithThreadSafe[int]())
					} else {
						buf = New[int](elems)
					}
					buf.AddRange(data[elems])
					b.StartTimer()
					benchmarkDequeue(buf, elems)
				}
			})
		}
	}

	for _, elems := range elements {
		b.Run(fmt.Sprintf("Queue-Sort-%d-NA-NA-NA", elems), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				buf = New[int](elems)
				buf.AddRange(data[elems])
				b.StartTimer()
				buf.Sort()
			}
		})
	}

	for _, elems := range elements {
		b.Run(fmt.Sprintf("Queue-Contains-%d-NA-NA-NA", elems), func(b *testing.B) {
			buf = New[int](elems)
			buf.AddRange(data[elems])
			lookup := make([]int, elems)
			copy(lookup, data[elems])
			rand.Shuffle(elems, func(i, j int) {
				lookup[i], lookup[j] = lookup[j], lookup[i]
			})

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				buf.Contains(lookup[i%elems])
			}
		})
	}

	for _, elems := range elements {
		buf = New[int](elems)
		buf.AddRange(data[elems])
		b.ResetTimer()
		b.Run(fmt.Sprintf("Queue-Min-%d-NA-NA-NA", elems), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				buf.Min()
			}
		})

		b.Run(fmt.Sprintf("Queue-Max-%d-NA-NA-NA", elems), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				buf.Max()
			}
		})
	}
}
