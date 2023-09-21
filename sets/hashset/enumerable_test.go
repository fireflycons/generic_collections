package hashset

import (
	"testing"

	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestEnumerable(t *testing.T) {

	var s *HashSet[int]
	evens := []int{2, 4, 6, 8, 10}
	mixed := []int{2, 4, 6, 7, 10}

	t.Run("All is true for all even numbers", func(t *testing.T) {
		s = New[int]()
		s.AddRange(evens)
		s.All(func(i int) bool { return i%2 == 0 })
		require.True(t, s.All(func(i int) bool { return i%2 == 0 }))
	})

	t.Run("All (even numbers) is false for mixed even and odd numbers", func(t *testing.T) {
		s = New[int]()
		s.AddRange(mixed)
		require.False(t, s.All(func(i int) bool { return i%2 == 0 }))
	})

	t.Run("Any is true for even number in mixed even and odd numbers", func(t *testing.T) {
		s = New[int]()
		s.AddRange(mixed)
		require.True(t, s.Any(func(i int) bool { return i%2 == 0 }))
	})

	t.Run("Any is true for odd number in mixed even and odd numbers", func(t *testing.T) {
		s = New[int]()
		s.AddRange(mixed)
		require.True(t, s.Any(func(i int) bool { return i%2 != 0 }))
	})

	t.Run("Any is false for number not in input slice", func(t *testing.T) {
		s = New[int]()
		s.AddRange(mixed)
		require.False(t, s.Any(func(i int) bool { return i > 1000 }))
	})

	t.Run("ForEach applies func to all elements", func(t *testing.T) {
		s = New[int]()
		expected := make([]int, len(evens))
		actual := make([]int, 0, len(evens))

		for i, v := range evens {
			expected[i] = v * v
		}

		s.AddRange(evens)
		s.ForEach(func(e collections.Element[int]) {
			actual = append(actual, e.Value()*e.Value())
		})

		// Output won't be in the same order as input slice
		require.ElementsMatch(t, expected, actual)
	})

	t.Run("ForEach panics when attempting to set value", func(t *testing.T) {
		s = New[int]()
		expected := make([]int, len(evens))

		for i, v := range evens {
			expected[i] = v * v
		}

		s.AddRange(evens)
		require.Panics(t, func() {
			s.ForEach(func(e collections.Element[int]) {
				*(e.ValuePtr()) = 0
			})
		})
	})

	t.Run("Map applies func to all elements and retuns new collection", func(t *testing.T) {
		s = New[int]()
		expected := make([]int, len(evens))

		for i, v := range evens {
			expected[i] = v * v
		}

		s.AddRange(evens)
		s1 := s.Map(func(i int) int {
			return i * i
		})

		// Output won't be in the same order as input slice
		require.ElementsMatch(t, expected, s1.ToSlice())
	})

	t.Run("Select selects all values <= 6", func(t *testing.T) {
		s = New[int]()
		expected := []int{2, 4, 6}

		s.AddRange(evens)
		s1 := s.Select(func(i int) bool { return i <= 6 })

		// Output won't be in the same order as input slice
		require.ElementsMatch(t, expected, s1.ToSlice())
	})
}

func TestFindAll(t *testing.T) {
	var tempItems, headItems []int
	var s *HashSet[int]
	arraySize := 16
	seed := int64(21543)
	headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Finds all even numbers", func(t *testing.T) {
		s = New[int]()
		expected := make([]int, 0, len(headItems))
		for i := 0; i < len(headItems); i++ {
			if headItems[i]%2 == 0 {
				expected = append(expected, headItems[i])
			}
		}
		s.AddRange(headItems)
		elems := s.FindAll(func(v int) bool { return v%2 == 0 })
		tempItems = make([]int, len(elems))
		for i := 0; i < len(elems); i++ {
			tempItems[i] = elems[i].Value()
		}

		// Set is not ordered
		require.ElementsMatch(t, expected, tempItems)
	})
}

func TestMinMax(t *testing.T) {

	var s *HashSet[int]
	seed := int64(2163)
	setItems, min, max := util.CreateMinMaxTestData(util.DefaultCapacity, &seed)

	t.Run("Min", func(t *testing.T) {
		s = New[int]()
		s.AddRange(setItems)
		require.Equal(t, min, s.Min())
	})

	t.Run("Max", func(t *testing.T) {
		s = New[int]()
		s.AddRange(setItems)
		require.Equal(t, max, s.Max())
	})
}
