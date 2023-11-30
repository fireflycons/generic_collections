package dlist

import (
	"testing"

	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestEnumerable(t *testing.T) {

	var ll *DList[int]
	evens := []int{2, 4, 6, 8, 10}
	mixed := []int{2, 4, 6, 7, 10}

	t.Run("All is true for all even numbers", func(t *testing.T) {
		ll = New[int]()
		ll.AddRange(evens)
		ll.All(func(i int) bool { return i%2 == 0 })
		require.True(t, ll.All(func(i int) bool { return i%2 == 0 }))
	})

	t.Run("All (even numbers) is false for mixed even and odd numbers", func(t *testing.T) {
		ll = New[int]()
		ll.AddRange(mixed)
		require.False(t, ll.All(func(i int) bool { return i%2 == 0 }))
	})

	t.Run("Any is true for even number in mixed even and odd numbers", func(t *testing.T) {
		ll = New[int]()
		ll.AddRange(mixed)
		require.True(t, ll.Any(func(i int) bool { return i%2 == 0 }))
	})

	t.Run("Any is true for odd number in mixed even and odd numbers", func(t *testing.T) {
		ll = New[int]()
		ll.AddRange(mixed)
		require.True(t, ll.Any(func(i int) bool { return i%2 != 0 }))
	})

	t.Run("Any is false for number not in input slice", func(t *testing.T) {
		ll = New[int]()
		ll.AddRange(mixed)
		require.False(t, ll.Any(func(i int) bool { return i > 1000 }))
	})

	t.Run("ForEach applies func to all elements", func(t *testing.T) {
		ll = New[int]()
		expected := make([]int, len(evens))
		actual := make([]int, 0, len(evens))

		for i, v := range evens {
			expected[i] = v * v
		}

		ll.AddRange(evens)
		ll.ForEach(func(e collections.Element[int]) {
			actual = append(actual, e.Value()*e.Value())
		})

		require.Equal(t, expected, actual)
	})

	t.Run("ForEach sets all even elements to zero", func(t *testing.T) {
		ll = New[int]()
		expected := make([]int, len(mixed))

		for i, v := range mixed {
			expected[i] = util.Iif(v%2 == 0, 0, v)
		}

		ll.AddRange(mixed)
		ll.ForEach(func(e collections.Element[int]) {
			if e.Value()%2 == 0 {
				*(e.ValuePtr()) = 0
			}
		})

		require.Equal(t, expected, ll.ToSlice())
	})

	t.Run("Map applies func to all elements and returns new collection", func(t *testing.T) {
		ll = New[int]()
		expected := make([]int, len(evens))

		for i, v := range evens {
			expected[i] = v * v
		}

		ll.AddRange(evens)
		ll1 := ll.Map(func(i int) int {
			return i * i
		})

		// Output won't be in the same order as input slice
		require.ElementsMatch(t, expected, ll1.ToSlice())
	})

	t.Run("Select selects all values <= 6", func(t *testing.T) {
		ll = New[int]()
		expected := []int{2, 4, 6}

		ll.AddRange(evens)
		s1 := ll.Select(func(i int) bool { return i <= 6 })

		// Output won't be in the same order as input slice
		require.ElementsMatch(t, expected, s1.ToSlice())
	})
}

func TestMinMax(t *testing.T) {

	var ll *DList[int]
	seed := int64(2163)
	setItems, min, max := util.CreateMinMaxTestData(util.DefaultCapacity, &seed)

	t.Run("Min", func(t *testing.T) {
		ll = New[int]()
		ll.AddRange(setItems)
		require.Equal(t, min, ll.Min())
	})

	t.Run("Max", func(t *testing.T) {
		ll = New[int]()
		ll.AddRange(setItems)
		require.Equal(t, max, ll.Max())
	})
}
