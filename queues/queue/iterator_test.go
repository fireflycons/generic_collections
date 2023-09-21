package queue

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestForwardIterator(t *testing.T) {

	var additionalItems, queueItems, tempItems []int
	var queue *Queue[int]
	arraySize := 16
	seed := int64(21543)
	queueItems, _, _, _ = util.CreateIntListData(arraySize, &seed)
	additionalArraySize := 4
	additionalItems, _, _, _ = util.CreateIntListData(additionalArraySize, &seed)

	_ = additionalItems

	t.Run("Iterates all values", func(t *testing.T) {
		queue = New[int]()

		for _, v := range queueItems {
			queue.Enqueue(v)
		}

		iter := queue.Iterator()

		tempItems = make([]int, 0, arraySize)
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, queueItems, tempItems)
	})

	t.Run("Remove some items, re-add them", func(t *testing.T) {

		// So that head and tail != 0
		expectedItems := removeAndReAdd(&queue, queueItems)

		iter := queue.Iterator()
		tempItems = make([]int, 0, arraySize)
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, expectedItems, tempItems)
	})

	t.Run("Force queue to grow", func(t *testing.T) {
		queue = New(WithCapacity[int](10))

		for _, v := range queueItems {
			queue.Enqueue(v)
		}

		iter := queue.Iterator()
		tempItems = make([]int, 0, arraySize)
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, queueItems, tempItems)
	})

	t.Run("Queue with gap in - some elements removed", func(t *testing.T) {
		expectedItems := createGappedQueue(&queue, queueItems)

		iter := queue.Iterator()
		tempItems = make([]int, 0, len(expectedItems))
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, expectedItems, tempItems)
	})

	t.Run("Start iteration on empty queue returns nil element", func(t *testing.T) {
		queue = New(WithCapacity[int](10))
		iter := queue.Iterator()
		require.Nil(t, iter.Start())
	})

}

func TestReverseIterator(t *testing.T) {

	var additionalItems, queueItems, tempItems, queueItemsReverse []int
	var queue *Queue[int]
	arraySize := 16
	seed := int64(21543)
	queueItems, _, queueItemsReverse, _ = util.CreateIntListData(arraySize, &seed)
	additionalArraySize := 4
	additionalItems, _, _, _ = util.CreateIntListData(additionalArraySize, &seed)

	_ = additionalItems

	t.Run("Iterates all values", func(t *testing.T) {
		queue = New[int]()

		for i := 0; i < arraySize; i++ {
			queue.Enqueue(queueItems[i])
		}

		iter := queue.ReverseIterator()

		tempItems = make([]int, 0, arraySize)
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, queueItemsReverse, tempItems)
	})

	t.Run("Remove some items, re-add them", func(t *testing.T) {

		// So that head and tail != 0
		expectedItems := removeAndReAdd(&queue, queueItems)

		iter := queue.ReverseIterator()
		tempItems = make([]int, 0, arraySize)
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, util.Reverse(expectedItems), tempItems)
	})

	t.Run("Force queue to grow", func(t *testing.T) {
		queue = New(WithCapacity[int](10))

		for _, v := range queueItems {
			queue.Enqueue(v)
		}

		iter := queue.ReverseIterator()
		tempItems = make([]int, 0, arraySize)
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, queueItemsReverse, tempItems)
	})

	t.Run("Queue with gap in - some elements removed", func(t *testing.T) {
		expectedItems := util.Reverse(createGappedQueue(&queue, queueItems))

		iter := queue.ReverseIterator()
		tempItems = make([]int, 0, len(expectedItems))
		for e := iter.Start(); e != nil; e = iter.Next() {
			tempItems = append(tempItems, e.Value())
		}

		require.Equal(t, expectedItems, tempItems)
	})

}

func TestTakeWhile(t *testing.T) {
	var queueItems, iteratedItems []int
	seed := int64(2163)
	queueItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Returns even numbers", func(t *testing.T) {
		queue := New[int]()

		queue.AddRange(queueItems)

		iter := queue.TakeWhile(func(val int) bool { return val%2 == 0 })

		iteratedItems = make([]int, 0, util.DefaultCapacity)

		for e := iter.Start(); e != nil; e = iter.Next() {
			iteratedItems = append(iteratedItems, e.Value())
		}

		tempItems := make([]int, 0, util.DefaultCapacity)

		for _, v := range queueItems {
			if v%2 == 0 {
				tempItems = append(tempItems, v)
			}
		}

		require.ElementsMatch(t, tempItems, iteratedItems)

	})
}

func TestWhere(t *testing.T) {
	var queueItems []int
	seed := int64(2163)
	queueItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Returns even numbers", func(t *testing.T) {
		queue := New[int]()
		queue.AddRange(queueItems)
		queue1 := queue.Select(func(val int) bool { return val%2 == 0 })

		tempItems := make([]int, 0, util.DefaultCapacity)

		for _, v := range queueItems {
			if v%2 == 0 {
				tempItems = append(tempItems, v)
			}
		}

		require.ElementsMatch(t, tempItems, queue1.ToSlice())

	})
}

func TestQueueForwardIterator_Negative(t *testing.T) {

	var queue *Queue[int]
	seed := int64(8293)

	t.Run("Modifying collection during iteration invalidates iterator", func(t *testing.T) {

		queue = New[int]()

		for i := 0; i < 3; i++ {
			queue.Enqueue(util.CreateRandInt(&seed))
		}

		iter := queue.Iterator()
		require.NotPanics(t, func() { iter.Start() })
		queue.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Next() })
	})

	t.Run("Modifying collection before iteration invalidates iterator", func(t *testing.T) {

		queue = New[int]()

		for i := 0; i < 3; i++ {
			queue.Enqueue(util.CreateRandInt(&seed))
		}

		iter := queue.Iterator()
		queue.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Start() })
	})

	t.Run("Modifying collection after taking an element invalidates element (Value)", func(t *testing.T) {
		queue = New[int]()

		for i := 0; i < 3; i++ {
			queue.Enqueue(util.CreateRandInt(&seed))
		}

		iter := queue.Iterator()
		element := iter.Start()
		queue.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { element.Value() })
	})

	t.Run("Modifying collection after taking an element invalidates element (ValuePtr)", func(t *testing.T) {
		queue = New[int]()

		for i := 0; i < 3; i++ {
			queue.Enqueue(util.CreateRandInt(&seed))
		}

		iter := queue.Iterator()
		element := iter.Start()
		queue.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { element.ValuePtr() })
	})
}

func TestQueueReverseIterator_Negative(t *testing.T) {

	var queue *Queue[int]
	seed := int64(8293)

	t.Run("Modifying collection during iteration invalidates iterator", func(t *testing.T) {

		queue = New[int]()

		for i := 0; i < 3; i++ {
			queue.Enqueue(util.CreateRandInt(&seed))
		}

		iter := queue.ReverseIterator()
		require.NotPanics(t, func() { iter.Start() })
		queue.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Next() })
	})

	t.Run("Modifying collection before iteration invalidates iterator", func(t *testing.T) {

		queue = New[int]()

		for i := 0; i < 3; i++ {
			queue.Enqueue(util.CreateRandInt(&seed))
		}

		iter := queue.ReverseIterator()
		queue.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Start() })
	})

	t.Run("Modifying collection after taking an element invalidates element (Value)", func(t *testing.T) {
		queue = New[int]()

		for i := 0; i < 3; i++ {
			queue.Enqueue(util.CreateRandInt(&seed))
		}

		iter := queue.ReverseIterator()
		element := iter.Start()
		queue.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { element.Value() })
	})

	t.Run("Modifying collection after taking an element invalidates element (ValuePtr)", func(t *testing.T) {
		queue = New[int]()

		for i := 0; i < 3; i++ {
			queue.Enqueue(util.CreateRandInt(&seed))
		}

		iter := queue.ReverseIterator()
		element := iter.Start()
		queue.Enqueue(util.CreateRandInt(&seed))
		require.Panics(t, func() { element.ValuePtr() })
	})
}
