package util

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGosort(t *testing.T) {

	compare := GetDefaultComparer[int]()
	seed := int64(4759)

	for arraySize := 100; arraySize <= 100000; arraySize *= 10 {
		t.Run(fmt.Sprintf("%d elements", arraySize), func(t *testing.T) {
			input, _, _, _ := CreateIntListData(arraySize, &seed)

			expected := make([]int, len(input))
			copy(expected, input)
			sort.Ints(expected)

			Gosort(input, len(input), compare)

			require.Equal(t, expected, input)
		})
	}

	// When buffer is larger than the number of elements in it that we want to sort,
	// zero values should remain at the end.
	for i := 1; i <= 2; i++ {
		arraySize := 192 * i
		t.Run(fmt.Sprintf("oversize buffer %d elements", arraySize), func(t *testing.T) {
			oversize := 16
			tempInput, _, _, _ := CreateIntListData(arraySize, &seed)

			input := make([]int, len(tempInput)+oversize)
			expected := make([]int, len(tempInput)+oversize)
			tempExpected := make([]int, len(tempInput))

			copy(input, tempInput)
			copy(tempExpected, input)
			sort.Ints(tempExpected)
			copy(expected, tempExpected)

			Gosort(input, len(tempInput), compare)

			require.Equal(t, expected, input)
		})
	}
}

func TestGosortDescending(t *testing.T) {

	compare := GetDefaultComparer[int]()
	seed := int64(4759)

	for arraySize := 100; arraySize <= 100000; arraySize *= 10 {
		t.Run(fmt.Sprintf("%d elements", arraySize), func(t *testing.T) {
			input, _, _, _ := CreateIntListData(arraySize, &seed)

			expected := make([]int, len(input))
			copy(expected, input)
			sort.Ints(expected)
			Reverse(expected)
			GosortDescending(input, len(input), compare)

			require.Equal(t, expected, input)
		})
	}

	// When buffer is larger than the number of elements in it that we want to sort,
	// zero values should remain at the end.
	for i := 1; i <= 2; i++ {
		arraySize := 192 * i
		t.Run(fmt.Sprintf("oversize buffer %d elements", arraySize), func(t *testing.T) {
			oversize := 16
			tempInput, _, _, _ := CreateIntListData(arraySize, &seed)

			input := make([]int, len(tempInput)+oversize)
			expected := make([]int, len(tempInput)+oversize)
			tempExpected := make([]int, len(tempInput))

			copy(input, tempInput)
			copy(tempExpected, input)
			sort.Ints(tempExpected)
			Reverse(tempExpected)
			copy(expected, tempExpected)
			GosortDescending(input, len(tempInput), compare)

			require.Equal(t, expected, input)
		})
	}
}

func BenchmarkGosort(b *testing.B) {

	compare := GetDefaultComparer[int]()
	seed := int64(4759)

	for arraySize := 100; arraySize <= 100000; arraySize *= 10 {
		b.Run(fmt.Sprintf("Ascending %d elements", arraySize), func(b *testing.B) {
			input, _, _, _ := CreateIntListData(arraySize, &seed)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Gosort(input, len(input), compare)
			}
		})
	}

	for arraySize := 100; arraySize <= 100000; arraySize *= 10 {
		b.Run(fmt.Sprintf("Descending %d elements", arraySize), func(b *testing.B) {
			input, _, _, _ := CreateIntListData(arraySize, &seed)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				GosortDescending(input, len(input), compare)
			}
		})
	}
}
