package util

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReverse(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []int{5, 4, 3, 2, 1}

	actual := Reverse(input)

	require.ElementsMatch(t, expected, actual)
}

func TestReverseSubset1(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []int{1, 5, 4, 3, 2}

	actual := ReverseSubset(input, 1, len(input)-1)

	require.ElementsMatch(t, expected, actual)
}

func TestReverseSubset2(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []int{1, 4, 3, 2, 5}

	actual := ReverseSubset(input, 1, len(input)-2)

	require.ElementsMatch(t, expected, actual)
}

func TestPartialCopy(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []int{0, 2, 3, 4, 0}
	dest := make([]int, 5)

	PartialCopy(input, 1, dest, 1, 3)

	require.ElementsMatch(t, expected, dest)
}

func cmp(v1, v2 int) int {
	return v1 - v2
}

func TestMinMax(t *testing.T) {
	seed := int64(2163)
	slc, expectedMin, expectedMax := CreateMinMaxTestData(DefaultCapacity, &seed)

	t.Run("Min", func(t *testing.T) {
		min := min(slc, cmp)

		require.Equal(t, expectedMin, min)
	})

	t.Run("Max", func(t *testing.T) {
		max := max(slc, cmp)

		require.Equal(t, expectedMax, max)
	})
}

func TestIndexOf(t *testing.T) {

	arr := []int{0, 1, 2, 3, 4, 5, 6, 6, 7, 8, 9, 10}
	comparer := GetDefaultComparer[int]()

	t.Run("IndexOf returns index of first matching element", func(t *testing.T) {
		expectedIndex := 6
		value := 6
		actualIndex := indexOf(arr, value, comparer)

		require.Equal(t, expectedIndex, actualIndex)
	})

	t.Run("IndexOf returns -1 when no match", func(t *testing.T) {
		expectedIndex := -1
		value := 100
		actualIndex := indexOf(arr, value, comparer)

		require.Equal(t, expectedIndex, actualIndex)
	})
}

func TestIndexOfConcurrent(t *testing.T) {
	arrSize := 1000
	expectedIndex := arrSize - 1
	value := 1
	arr := make([]int, arrSize)
	arr[expectedIndex] = value
	arr[(arrSize/2)-1] = value

	for n := 2; n <= 20; n++ {
		actualIndex := indexOfConcurrent[int](arr, value, xxCompareSignedInt[int], n)
		require.Equal(t, expectedIndex, actualIndex)
	}
}

func TestMinMaxConcurrent(t *testing.T) {
	arrSize := 1000
	lastIndex := arrSize - 1
	maxValue := 1
	minValue := -1

	arr := make([]int, arrSize)
	arr[lastIndex] = maxValue

	for n := 2; n <= 20; n++ {
		max := getMinOrMaxConcurrent[int](arr, xxCompareSignedInt[int], max[int], n)
		require.Equal(t, maxValue, max)
	}

	arr[lastIndex] = minValue

	for n := 2; n <= 20; n++ {
		min := getMinOrMaxConcurrent[int](arr, xxCompareSignedInt[int], min[int], n)
		require.Equal(t, minValue, min)
	}
}

func TestDeepCopy(t *testing.T) {

	t.Run("With int slice", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := input
		dest := make([]int, len(input))
		DeepCopySlice(dest, input, nil)
		require.Equal(t, dest, expected)
	})

	t.Run("With pointer slice", func(t *testing.T) {
		i1, i2, i3, i4 := 1, 2, 3, 4
		input := []*int{&i1, &i2, &i3, &i4}
		dest := make([]*int, len(input))

		DeepCopySlice(dest, input, func(ptr *int) *int {
			val := *ptr
			return &val
		})

		for i := 0; i < len(input); i++ {
			require.NotSame(t, input[i], dest[i])
			require.Equal(t, *input[i], *dest[i])
		}
	})

	t.Run("With struct", func(t *testing.T) {
		type myStruct struct {
			i   int
			ptr *int
		}

		i1 := 42

		s1 := myStruct{
			i:   1,
			ptr: &i1,
		}

		s2 := DeepCopy(s1, func(s myStruct) myStruct {
			i2 := *s.ptr
			return myStruct{
				i:   s.i,
				ptr: &i2,
			}
		})

		require.NotSame(t, s1.ptr, s2.ptr)
		require.Equal(t, s1.i, s2.i)
		require.Equal(t, *s1.ptr, *s2.ptr)
	})
}

func TestLastIndexOf(t *testing.T) {

	arr := []int{0, 1, 2, 3, 4, 5, 6, 6, 7, 8, 9, 10}
	comparer := GetDefaultComparer[int]()

	t.Run("LastIndexOf returns index of last matching element", func(t *testing.T) {
		expectedIndex := 7
		value := 6
		actualIndex := lastIndexOf(arr, value, comparer)

		require.Equal(t, expectedIndex, actualIndex)
	})

	t.Run("IndexOf returns -1 when no match", func(t *testing.T) {
		expectedIndex := -1
		value := 100
		actualIndex := lastIndexOf(arr, value, comparer)

		require.Equal(t, expectedIndex, actualIndex)
	})
}

func TestLastIndexOfConcurrent(t *testing.T) {
	arrSize := 1000
	expectedIndex := arrSize - 1
	value := 1
	arr := make([]int, arrSize)

	// Smatter some values in
	r := rand.New(rand.NewSource(10007))
	for i := 0; i < 25; i++ {
		arr[r.Intn(arrSize)] = value
	}

	// And make the last one the one we want
	arr[arrSize-1] = value

	for n := 2; n <= 20; n++ {
		actualIndex := lastIndexOfConcurrent[int](arr, value, xxCompareSignedInt[int], n)
		require.Equal(t, expectedIndex, actualIndex)
	}
}

func BenchmarkMinMax(b *testing.B) {

	iters := []int{10000, 100000, 1000000}
	seed := int64(2163)

	for _, v := range iters {

		data := CreateSingleIntListData(v, &seed)

		b.Run(fmt.Sprintf("Min %d", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Min(data, cmp, false)
			}
		})

		b.Run(fmt.Sprintf("Min %d concurrent", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Min(data, cmp, true)
			}
		})

		b.Run(fmt.Sprintf("Max %d", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Max(data, cmp, false)
			}
		})

		b.Run(fmt.Sprintf("Max %d concurrent", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Max(data, cmp, true)
			}
		})
	}
}

func BenchmarkIndexOf(b *testing.B) {

	iters := []int{10000, 100000, 1000000}
	seed := int64(2163)

	for _, v := range iters {

		data := CreateSingleIntListData(v, &seed)

		b.Run(fmt.Sprintf("IndexOf %d", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				IndexOf[int](data, data[len(data)-1], xxCompareSignedInt[int], false)
			}
		})

		b.Run(fmt.Sprintf("IndexOf %d concurrent", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				IndexOf[int](data, data[len(data)-1], xxCompareSignedInt[int], true)
			}
		})

	}
}
