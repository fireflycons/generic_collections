package queue

import (
	"sort"
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestSort(t *testing.T) {

	var queueItems, tempItems []int
	var queue *Queue[int]
	arraySize := 16
	seed := int64(21543)
	queueItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Empty queue", func(t *testing.T) {
		queue = New[int]()
		queue.Sort()
		verifyQueueState(t, queue, []int{})
	})

	t.Run("Full queue", func(t *testing.T) {
		queue = New[int]()
		queue.AddRange(queueItems)
		queue.Sort()
		tempItems = make([]int, len(queueItems))
		copy(tempItems, queueItems)
		sort.Ints(tempItems)
		verifyQueueState(t, queue, tempItems)
		require.Equal(t, tempItems[0], queue.Peek())
	})

	t.Run("Queue with space at end", func(t *testing.T) {
		queue = New(WithCapacity[int](32))
		queue.AddRange(queueItems)
		queue.Sort()
		tempItems = make([]int, len(queueItems))
		copy(tempItems, queueItems)
		sort.Ints(tempItems)
		verifyQueueState(t, queue, tempItems)
		require.Equal(t, tempItems[0], queue.Peek())
	})

	t.Run("Queue with gap in - some elements removed", func(t *testing.T) {
		tempItems := createGappedQueue(&queue, queueItems)
		verifyQueueState(t, queue, tempItems)
		queue.Sort()
		expectedItems := make([]int, len(tempItems))
		copy(expectedItems, tempItems)
		sort.Ints(expectedItems)
		verifyQueueState(t, queue, expectedItems)
		require.Equal(t, expectedItems[0], queue.Peek())
	})
}

func TestSorted(t *testing.T) {

	var queueItems, tempItems []int
	var queue *Queue[int]
	arraySize := 16
	seed := int64(21543)
	queueItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Empty queue", func(t *testing.T) {
		queue = New[int]()
		q1, ok := queue.Sorted().(*Queue[int])
		require.True(t, ok, "type assertion failed")
		verifyQueueState(t, q1, []int{})
	})

	t.Run("Full queue", func(t *testing.T) {
		queue = New[int]()
		queue.AddRange(queueItems)
		q1, ok := queue.Sorted().(*Queue[int])
		require.True(t, ok, "type assertion failed")
		tempItems = make([]int, len(queueItems))
		copy(tempItems, queueItems)
		sort.Ints(tempItems)
		verifyQueueState(t, q1, tempItems)
		require.Equal(t, tempItems[0], q1.Peek())
	})

	t.Run("Queue with space at end", func(t *testing.T) {
		queue = New(WithCapacity[int](32))
		queue.AddRange(queueItems)
		q1, ok := queue.Sorted().(*Queue[int])
		require.True(t, ok, "type assertion failed")
		tempItems = make([]int, len(queueItems))
		copy(tempItems, queueItems)
		sort.Ints(tempItems)
		verifyQueueState(t, q1, tempItems)
		require.Equal(t, tempItems[0], q1.Peek())
	})

	t.Run("Queue with gap in - some elements removed", func(t *testing.T) {
		tempItems := createGappedQueue(&queue, queueItems)
		verifyQueueState(t, queue, tempItems)
		q1, ok := queue.Sorted().(*Queue[int])
		require.True(t, ok, "type assertion failed")
		expectedItems := make([]int, len(tempItems))
		copy(expectedItems, tempItems)
		sort.Ints(expectedItems)
		verifyQueueState(t, q1, expectedItems)
		require.Equal(t, expectedItems[0], q1.Peek())
	})
}

func TestSortDescending(t *testing.T) {

	var queueItems, tempItems []int
	var queue *Queue[int]
	arraySize := 16
	seed := int64(21543)
	queueItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Empty queue", func(t *testing.T) {
		queue = New[int]()
		queue.SortDescending()
		verifyQueueState(t, queue, []int{})
	})

	t.Run("Full queue", func(t *testing.T) {
		queue = New[int]()
		queue.AddRange(queueItems)
		queue.SortDescending()
		tempItems = make([]int, len(queueItems))
		copy(tempItems, queueItems)
		sort.Ints(tempItems)
		util.Reverse(tempItems)
		verifyQueueState(t, queue, tempItems)
		require.Equal(t, tempItems[0], queue.Peek())
	})

	t.Run("Queue with space at end", func(t *testing.T) {
		queue = New(WithCapacity[int](32))
		queue.AddRange(queueItems)
		queue.SortDescending()
		tempItems = make([]int, len(queueItems))
		copy(tempItems, queueItems)
		sort.Ints(tempItems)
		util.Reverse(tempItems)
		verifyQueueState(t, queue, tempItems)
		require.Equal(t, tempItems[0], queue.Peek())
	})

	t.Run("Queue with gap in - some elements removed", func(t *testing.T) {
		tempItems := createGappedQueue(&queue, queueItems)
		verifyQueueState(t, queue, tempItems)
		queue.SortDescending()
		expectedItems := make([]int, len(tempItems))
		copy(expectedItems, tempItems)
		sort.Ints(expectedItems)
		util.Reverse(expectedItems)
		verifyQueueState(t, queue, expectedItems)
		require.Equal(t, expectedItems[0], queue.Peek())
	})
}

func TestSortedDescending(t *testing.T) {

	var queueItems, tempItems []int
	var queue *Queue[int]
	arraySize := 16
	seed := int64(21543)
	queueItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Empty queue", func(t *testing.T) {
		queue = New[int]()
		q1, ok := queue.SortedDescending().(*Queue[int])
		require.True(t, ok, "type assertion failed")
		verifyQueueState(t, q1, []int{})
	})

	t.Run("Full queue", func(t *testing.T) {
		queue = New[int]()
		queue.AddRange(queueItems)
		q1, ok := queue.SortedDescending().(*Queue[int])
		require.True(t, ok, "type assertion failed")
		tempItems = make([]int, len(queueItems))
		copy(tempItems, queueItems)
		sort.Ints(tempItems)
		util.Reverse(tempItems)
		verifyQueueState(t, q1, tempItems)
		require.Equal(t, tempItems[0], q1.Peek())
	})

	t.Run("Queue with space at end", func(t *testing.T) {
		queue = New(WithCapacity[int](32))
		queue.AddRange(queueItems)
		q1, ok := queue.SortedDescending().(*Queue[int])
		require.True(t, ok, "type assertion failed")
		tempItems = make([]int, len(queueItems))
		copy(tempItems, queueItems)
		sort.Ints(tempItems)
		util.Reverse(tempItems)
		verifyQueueState(t, q1, tempItems)
		require.Equal(t, tempItems[0], q1.Peek())
	})

	t.Run("Queue with gap in - some elements removed", func(t *testing.T) {
		tempItems := createGappedQueue(&queue, queueItems)
		verifyQueueState(t, queue, tempItems)
		q1, ok := queue.SortedDescending().(*Queue[int])
		require.True(t, ok, "type assertion failed")
		expectedItems := make([]int, len(tempItems))
		copy(expectedItems, tempItems)
		sort.Ints(expectedItems)
		util.Reverse(expectedItems)
		verifyQueueState(t, q1, expectedItems)
		require.Equal(t, expectedItems[0], q1.Peek())
	})
}
