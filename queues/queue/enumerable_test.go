package queue

import (
	"fmt"
	"testing"

	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestEnumerable(t *testing.T) {

	var q *Queue[int]
	evens := []int{2, 4, 6, 8, 10}
	mixed := []int{2, 4, 6, 7, 10}

	t.Run("All is true for all even numbers", func(t *testing.T) {
		q = New[int]()
		q.AddRange(evens)
		q.All(func(i int) bool { return i%2 == 0 })
		require.True(t, q.All(func(i int) bool { return i%2 == 0 }))
	})

	t.Run("All (even numbers) is false for mixed even and odd numbers", func(t *testing.T) {
		q = New[int]()
		q.AddRange(mixed)
		require.False(t, q.All(func(i int) bool { return i%2 == 0 }))
	})

	t.Run("Any is true for even number in mixed even and odd numbers", func(t *testing.T) {
		q = New[int]()
		q.AddRange(mixed)
		require.True(t, q.Any(func(i int) bool { return i%2 == 0 }))
	})

	t.Run("Any is true for odd number in mixed even and odd numbers", func(t *testing.T) {
		q = New[int]()
		q.AddRange(mixed)
		require.True(t, q.Any(func(i int) bool { return i%2 != 0 }))
	})

	t.Run("Any is false for number not in input slice", func(t *testing.T) {
		q = New[int]()
		q.AddRange(mixed)
		require.False(t, q.Any(func(i int) bool { return i > 1000 }))
	})

	t.Run("ForEach applies func to all elements", func(t *testing.T) {
		q = New[int]()
		expected := make([]int, len(evens))
		actual := make([]int, 0, len(evens))

		for i, v := range evens {
			expected[i] = v * v
		}

		q.AddRange(evens)
		q.ForEach(func(e collections.Element[int]) {
			actual = append(actual, e.Value()*e.Value())
		})

		// Output won't be in the same order as input slice
		require.ElementsMatch(t, expected, actual)
	})

	t.Run("ForEach sets all even elements to zero", func(t *testing.T) {
		q = New[int]()
		expected := make([]int, len(mixed))

		for i, v := range mixed {
			expected[i] = util.Iif(v%2 == 0, 0, v)
		}

		q.AddRange(mixed)
		q.ForEach(func(e collections.Element[int]) {
			if e.Value()%2 == 0 {
				*(e.ValuePtr()) = 0
			}
		})

		require.Equal(t, expected, q.ToSlice())
	})

	t.Run("Map applies func to all elements and retuns new collection", func(t *testing.T) {
		q = New[int]()
		expected := make([]int, len(evens))

		for i, v := range evens {
			expected[i] = v * v
		}

		q.AddRange(evens)
		q1 := q.Map(func(i int) int {
			return i * i
		})

		// Output won't be in the same order as input slice
		require.ElementsMatch(t, expected, q1.ToSlice())
	})

	t.Run("Select selects all values <= 6", func(t *testing.T) {
		q = New[int]()
		expected := []int{2, 4, 6}

		q.AddRange(evens)
		q1 := q.Select(func(i int) bool { return i <= 6 })

		// Output won't be in the same order as input slice
		require.ElementsMatch(t, expected, q1.ToSlice())
	})
}

func TestFindAll(t *testing.T) {
	var tempItems, headItems []int
	var q *Queue[int]
	arraySize := 16
	seed := int64(21543)
	headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Finds all even numbers", func(t *testing.T) {
		q = New[int]()
		expected := make([]int, 0, len(headItems))
		for i := 0; i < len(headItems); i++ {
			if headItems[i]%2 == 0 {
				expected = append(expected, headItems[i])
			}
		}
		q.AddRange(headItems)
		elems := q.FindAll(func(v int) bool { return v%2 == 0 })
		tempItems = make([]int, len(elems))
		for i := 0; i < len(elems); i++ {
			tempItems[i] = elems[i].Value()
		}

		require.Equal(t, expected, tempItems)
	})
}

func TestMinMax(t *testing.T) {

	var q *Queue[int]
	seed := int64(2163)

	for _, iters := range []int{100, 10000, 100000} {
		queueItems, min, max := util.CreateMinMaxTestData(iters, &seed)

		t.Run(fmt.Sprintf("Min %d", iters), func(t *testing.T) {
			q = New[int](WithCapacity[int](iters))
			q.AddRange(queueItems)
			require.Equal(t, min, q.Min())
		})

		t.Run(fmt.Sprintf("Min %d rotate buffer", iters), func(t *testing.T) {
			q = New[int](WithCapacity[int](iters))
			q.AddRange(queueItems)
			q.Enqueue(q.Dequeue())
			q.Enqueue(q.Dequeue())
			q.Enqueue(q.Dequeue())
			q.Enqueue(q.Dequeue())
			require.Equal(t, min, q.Min())
		})

		t.Run(fmt.Sprintf("Max %d", iters), func(t *testing.T) {
			q = New[int](WithCapacity[int](iters))
			q.AddRange(queueItems)
			require.Equal(t, max, q.Max())
		})

		t.Run(fmt.Sprintf("Max %d rotate buffer", iters), func(t *testing.T) {
			q = New[int](WithCapacity[int](iters))
			q.AddRange(queueItems)
			q.Enqueue(q.Dequeue())
			q.Enqueue(q.Dequeue())
			q.Enqueue(q.Dequeue())
			q.Enqueue(q.Dequeue())
			require.Equal(t, max, q.Max())
		})
	}

}
