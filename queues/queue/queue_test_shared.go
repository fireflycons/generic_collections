package queue

import "github.com/fireflycons/generic_collections/internal/util"

// Queue setup used in tests for both queue and queue_iterator

func removeAndReAdd[T any](queue **Queue[T], queueItems []T) []T {
	*queue = New[T]()
	moveElems := 4
	for _, v := range queueItems {
		(*queue).Enqueue(v)
	}

	tempItems := make([]T, 0, moveElems)

	for i := 0; i < moveElems; i++ {
		v := (*queue).Dequeue()
		tempItems = append(tempItems, v)
		(*queue).Enqueue(v)
	}

	expectedItems := make([]T, len(queueItems))
	copy(expectedItems, queueItems[moveElems:])
	util.PartialCopy(tempItems, 0, expectedItems, 12, moveElems)

	return expectedItems
}

func createGappedQueue[T any](queue **Queue[T], queueItems []T) []T {
	*queue = New[T]()
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
