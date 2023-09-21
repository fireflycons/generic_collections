package ringbuffer

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestForwardIterator(t *testing.T) {

	var additionalItems, bufferItems, tempItems []int
	var ringBuffer *RingBuffer[int]
	arraySize := 16
	seed := int64(21543)
	bufferItems, _, _, _ = util.CreateIntListData(arraySize, &seed)
	additionalArraySize := 4
	additionalItems, _, _, _ = util.CreateIntListData(additionalArraySize, &seed)

	_ = additionalItems

	t.Run("Iterates all values", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)

		for _, v := range bufferItems {
			ringBuffer.Enqueue(v)
		}

		iter := ringBuffer.Iterator()

		tempItems = make([]int, 0, arraySize)
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, bufferItems, tempItems)
	})

	t.Run("Remove some items, re-add them", func(t *testing.T) {

		// So that head and tail != 0
		expectedItems := removeAndReAdd(&ringBuffer, bufferItems)

		iter := ringBuffer.Iterator()
		tempItems = make([]int, 0, arraySize)
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, expectedItems, tempItems)
	})

	t.Run("Queue with gap in - some elements removed", func(t *testing.T) {
		expectedItems := createGappedBuffer(&ringBuffer, bufferItems)

		iter := ringBuffer.Iterator()
		tempItems = make([]int, 0, len(expectedItems))
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, expectedItems, tempItems)
	})

	t.Run("Start iteration on empty queue returns nil element", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)
		iter := ringBuffer.Iterator()
		require.Nil(t, iter.Start())
	})

}

func TestTakeWhile(t *testing.T) {
	var bufferItems, iteratedItems []int
	seed := int64(2163)
	bufferItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Returns even numbers", func(t *testing.T) {
		ringBuffer := New[int](util.DefaultCapacity)

		ringBuffer.AddRange(bufferItems)

		iter := ringBuffer.TakeWhile(func(val int) bool { return val%2 == 0 })

		iteratedItems = make([]int, 0, util.DefaultCapacity)

		for e := iter.Start(); e != nil; e = iter.Next() {
			iteratedItems = append(iteratedItems, e.Value())
		}

		tempItems := make([]int, 0, util.DefaultCapacity)

		for _, v := range bufferItems {
			if v%2 == 0 {
				tempItems = append(tempItems, v)
			}
		}

		require.ElementsMatch(t, tempItems, iteratedItems)

	})
}

func TestWhere(t *testing.T) {
	var bufferItems []int
	seed := int64(2163)
	bufferItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Returns even numbers", func(t *testing.T) {
		ringBuffer := New[int](util.DefaultCapacity)
		ringBuffer.AddRange(bufferItems)
		ringBuffer1 := ringBuffer.Select(func(val int) bool { return val%2 == 0 })

		tempItems := make([]int, 0, util.DefaultCapacity)

		for _, v := range bufferItems {
			if v%2 == 0 {
				tempItems = append(tempItems, v)
			}
		}

		require.ElementsMatch(t, tempItems, ringBuffer1.ToSlice())

	})
}

func TestReverseIterator(t *testing.T) {

	var additionalItems, bufferItems, tempItems, bufferItemsReverse []int
	var ringBuffer *RingBuffer[int]
	arraySize := 16
	seed := int64(21543)
	bufferItems, _, bufferItemsReverse, _ = util.CreateIntListData(arraySize, &seed)
	additionalArraySize := 4
	additionalItems, _, _, _ = util.CreateIntListData(additionalArraySize, &seed)

	_ = additionalItems

	t.Run("Iterates all values", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)

		for i := 0; i < arraySize; i++ {
			ringBuffer.Enqueue(bufferItems[i])
		}

		iter := ringBuffer.ReverseIterator()

		tempItems = make([]int, 0, arraySize)
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, bufferItemsReverse, tempItems)
	})

	t.Run("Remove some items, re-add them", func(t *testing.T) {

		// So that head and tail != 0
		expectedItems := removeAndReAdd(&ringBuffer, bufferItems)

		iter := ringBuffer.ReverseIterator()
		tempItems = make([]int, 0, arraySize)
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, util.Reverse(expectedItems), tempItems)
	})

	t.Run("Queue with gap in - some elements removed", func(t *testing.T) {
		expectedItems := util.Reverse(createGappedBuffer(&ringBuffer, bufferItems))

		iter := ringBuffer.ReverseIterator()
		tempItems = make([]int, 0, len(expectedItems))
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, expectedItems, tempItems)
	})

}

func TestForwardIterator_Negative(t *testing.T) {

	var ringBuffer *RingBuffer[int]
	seed := int64(8293)

	t.Run("Modifying collection during iteration invalidates iterator", func(t *testing.T) {

		ringBuffer = New[int](util.DefaultCapacity)

		for i := 0; i < 3; i++ {
			ringBuffer.Enqueue(util.CreateRandInt(&seed))
		}

		iter := ringBuffer.Iterator()
		require.NotPanics(t, func() { iter.Start() })
		ringBuffer.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Next() })
	})

	t.Run("Modifying collection before iteration invalidates iterator", func(t *testing.T) {

		ringBuffer = New[int](util.DefaultCapacity)

		for i := 0; i < 3; i++ {
			ringBuffer.Enqueue(util.CreateRandInt(&seed))
		}

		iter := ringBuffer.Iterator()
		ringBuffer.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Start() })
	})

	t.Run("Modifying collection after taking an element invalidates element (Value)", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)

		for i := 0; i < 3; i++ {
			ringBuffer.Enqueue(util.CreateRandInt(&seed))
		}

		iter := ringBuffer.Iterator()
		element := iter.Start()
		ringBuffer.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { element.Value() })
	})

	t.Run("Modifying collection after taking an element invalidates element (ValuePtr)", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)

		for i := 0; i < 3; i++ {
			ringBuffer.Enqueue(util.CreateRandInt(&seed))
		}

		iter := ringBuffer.Iterator()
		element := iter.Start()
		ringBuffer.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { element.ValuePtr() })
	})
}

func TestReverseIterator_Negative(t *testing.T) {

	var ringBuffer *RingBuffer[int]
	seed := int64(8293)

	t.Run("Modifying collection during iteration invalidates iterator", func(t *testing.T) {

		ringBuffer = New[int](util.DefaultCapacity)

		for i := 0; i < 3; i++ {
			ringBuffer.Enqueue(util.CreateRandInt(&seed))
		}

		iter := ringBuffer.ReverseIterator()
		require.NotPanics(t, func() { iter.Start() })
		ringBuffer.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Next() })
	})

	t.Run("Modifying collection before iteration invalidates iterator", func(t *testing.T) {

		ringBuffer = New[int](util.DefaultCapacity)

		for i := 0; i < 3; i++ {
			ringBuffer.Enqueue(util.CreateRandInt(&seed))
		}

		iter := ringBuffer.ReverseIterator()
		ringBuffer.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Start() })
	})

	t.Run("Modifying collection after taking an element invalidates element (Value)", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)

		for i := 0; i < 3; i++ {
			ringBuffer.Enqueue(util.CreateRandInt(&seed))
		}

		iter := ringBuffer.ReverseIterator()
		element := iter.Start()
		ringBuffer.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { element.Value() })
	})

	t.Run("Modifying collection after taking an element invalidates element (ValuePtr)", func(t *testing.T) {
		ringBuffer = New[int](util.DefaultCapacity)

		for i := 0; i < 3; i++ {
			ringBuffer.Enqueue(util.CreateRandInt(&seed))
		}

		iter := ringBuffer.ReverseIterator()
		element := iter.Start()
		ringBuffer.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { element.ValuePtr() })
	})
}

func createGappedBuffer[T any](queue **RingBuffer[T], queueItems []T) []T {
	*queue = New[T](util.DefaultCapacity)
	for _, v := range queueItems {
		(*queue).Enqueue(v)
	}

	// Remove four elements, add two back in leaving a gap in the middle of the buffer
	removed1 := make([]T, 4)

	for i := 0; i < 4; i++ {
		removed1[i] = (*queue).Dequeue()
	}

	for i := 0; i < 2; i++ {
		(*queue).Enqueue(removed1[i])
	}

	expectedItems := make([]T, len(queueItems)-2)
	util.PartialCopy(queueItems, 4, expectedItems, 0, len(queueItems)-4)
	util.PartialCopy(removed1, 0, expectedItems, len(queueItems)-4, 2)

	return expectedItems
}
