package ringbuffer

import (
	"sort"
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestSort(t *testing.T) {

	var bufferItems, tempItems []int
	var ringBuffer *RingBuffer[int]
	arraySize := 16
	seed := int64(21543)
	bufferItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Empty buffer", func(t *testing.T) {
		ringBuffer = New[int](arraySize)
		ringBuffer.Sort()
		verifyBufferState(t, ringBuffer, []int{})
	})

	t.Run("Full buffer", func(t *testing.T) {
		ringBuffer = New[int](arraySize)
		ringBuffer.AddRange(bufferItems)
		ringBuffer.Sort()
		tempItems = make([]int, len(bufferItems))
		copy(tempItems, bufferItems)
		sort.Ints(tempItems)
		verifyBufferState(t, ringBuffer, tempItems)
		require.Equal(t, tempItems[0], ringBuffer.Peek())
		require.Equal(t, ringBuffer.head, ringBuffer.tail)
	})

	t.Run("Buffer with space at end", func(t *testing.T) {
		ringBuffer = New[int](arraySize * 2)
		ringBuffer.AddRange(bufferItems)
		ringBuffer.Sort()
		tempItems = make([]int, len(bufferItems))
		copy(tempItems, bufferItems)
		sort.Ints(tempItems)
		verifyBufferState(t, ringBuffer, tempItems)
		require.Equal(t, tempItems[0], ringBuffer.Peek())
		require.Equal(t, ringBuffer.maxSize, len(ringBuffer.buffer))
	})
}

func TestSorted(t *testing.T) {

	var bufferItems, tempItems []int
	var ringBuffer *RingBuffer[int]
	arraySize := 16
	seed := int64(21543)
	bufferItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Empty buffer", func(t *testing.T) {
		ringBuffer = New[int](arraySize)
		buf1, ok := ringBuffer.Sorted().(*RingBuffer[int])
		require.True(t, ok, "type assertion failed")
		verifyBufferState(t, buf1, []int{})
	})

	t.Run("Full buffer", func(t *testing.T) {
		ringBuffer = New[int](arraySize)
		ringBuffer.AddRange(bufferItems)
		buf1, ok := ringBuffer.Sorted().(*RingBuffer[int])
		require.True(t, ok, "type assertion failed")
		tempItems = make([]int, len(bufferItems))
		copy(tempItems, bufferItems)
		sort.Ints(tempItems)
		verifyBufferState(t, buf1, tempItems)
		require.Equal(t, tempItems[0], buf1.Peek())
		require.Equal(t, buf1.head, buf1.tail)
	})

	t.Run("Buffer with space at end", func(t *testing.T) {
		ringBuffer = New[int](arraySize * 2)
		ringBuffer.AddRange(bufferItems)
		buf1, ok := ringBuffer.Sorted().(*RingBuffer[int])
		require.True(t, ok, "type assertion failed")
		tempItems = make([]int, len(bufferItems))
		copy(tempItems, bufferItems)
		sort.Ints(tempItems)
		verifyBufferState(t, buf1, tempItems)
		require.Equal(t, tempItems[0], buf1.Peek())
		require.Equal(t, buf1.maxSize, len(buf1.buffer))
	})
}

func TestSortDescending(t *testing.T) {

	var bufferItems, tempItems []int
	var ringBuffer *RingBuffer[int]
	arraySize := 16
	seed := int64(21543)
	bufferItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Empty buffer", func(t *testing.T) {
		ringBuffer = New[int](arraySize)
		ringBuffer.SortDescending()
		verifyBufferState(t, ringBuffer, []int{})
	})

	t.Run("Full buffer", func(t *testing.T) {
		ringBuffer = New[int](arraySize)
		ringBuffer.AddRange(bufferItems)
		ringBuffer.SortDescending()
		tempItems = make([]int, len(bufferItems))
		copy(tempItems, bufferItems)
		sort.Ints(tempItems)
		util.Reverse(tempItems)
		verifyBufferState(t, ringBuffer, tempItems)
		require.Equal(t, tempItems[0], ringBuffer.Peek())
		require.Equal(t, ringBuffer.head, ringBuffer.tail)
	})

	t.Run("Buffer with space at end", func(t *testing.T) {
		ringBuffer = New[int](arraySize * 2)
		ringBuffer.AddRange(bufferItems)
		ringBuffer.SortDescending()
		tempItems = make([]int, len(bufferItems))
		copy(tempItems, bufferItems)
		sort.Ints(tempItems)
		util.Reverse(tempItems)
		verifyBufferState(t, ringBuffer, tempItems)
		require.Equal(t, tempItems[0], ringBuffer.Peek())
		require.Equal(t, ringBuffer.maxSize, len(ringBuffer.buffer))
	})
}

func TestSortedDescending(t *testing.T) {

	var bufferItems, tempItems []int
	var ringBuffer *RingBuffer[int]
	arraySize := 16
	seed := int64(21543)
	bufferItems, _, _, _ = util.CreateSmallIntListData(arraySize, &seed)

	t.Run("Empty buffer", func(t *testing.T) {
		ringBuffer = New[int](arraySize)
		buf1, ok := ringBuffer.SortedDescending().(*RingBuffer[int])
		require.True(t, ok, "type assertion failed")
		verifyBufferState(t, buf1, []int{})
	})

	t.Run("Full buffer", func(t *testing.T) {
		ringBuffer = New[int](arraySize)
		ringBuffer.AddRange(bufferItems)
		buf1, ok := ringBuffer.SortedDescending().(*RingBuffer[int])
		require.True(t, ok, "type assertion failed")
		tempItems = make([]int, len(bufferItems))
		copy(tempItems, bufferItems)
		sort.Ints(tempItems)
		util.Reverse(tempItems)
		verifyBufferState(t, buf1, tempItems)
		require.Equal(t, tempItems[0], buf1.Peek())
		require.Equal(t, buf1.head, buf1.tail)
	})

	t.Run("Buffer with space at end", func(t *testing.T) {
		ringBuffer = New[int](arraySize * 2)
		ringBuffer.AddRange(bufferItems)
		buf1, ok := ringBuffer.SortedDescending().(*RingBuffer[int])
		require.True(t, ok, "type assertion failed")
		tempItems = make([]int, len(bufferItems))
		copy(tempItems, bufferItems)
		sort.Ints(tempItems)
		util.Reverse(tempItems)
		verifyBufferState(t, buf1, tempItems)
		require.Equal(t, tempItems[0], buf1.Peek())
		require.Equal(t, buf1.maxSize, len(buf1.buffer))
	})
}
