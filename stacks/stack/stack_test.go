package stack

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

var validCollectionSizes = []int{0, 1, 2, 4, 8, 13, 29, 128}

func generateIntStack(count int) *Stack[int] {

	stack := New(WithCapacity[int](count))
	if count == 0 {
		return stack
	}

	seed := int64(count * 34)

	for i := 0; i < count; i++ {
		rand.Seed(seed)
		stack.Push(rand.Int())
		seed++
	}

	return stack
}

func TestStackConstructor(t *testing.T) {

	t.Run("Create with default capacity", func(t *testing.T) {
		stack := New[int]()
		require.Equal(t, stack.capacity(), util.DefaultCapacity)
		require.Equal(t, stack.length(), util.DefaultCapacity)
	})

	t.Run("Create with capacity", func(t *testing.T) {
		for _, cap := range validCollectionSizes {
			t.Run(fmt.Sprintf("Create with capacity %d", cap), func(t *testing.T) {
				stack := New(WithCapacity[int](cap))

				require.Equal(t, stack.capacity(), cap)
				require.Equal(t, stack.length(), cap)
			})
		}
	})
}

func TestStackCount(t *testing.T) {
	stackSize := 20
	stack := generateIntStack(stackSize)
	require.Equal(t, stackSize, stack.Count())
}

func TestPeekAllElements(t *testing.T) {

	stackSize := 20
	stack := generateIntStack(stackSize)
	elements := stack.ToSlice()

	for _, expected := range elements {
		require.Equal(t, expected, stack.Peek())
		stack.Pop()
	}

	require.Equal(t, 0, stack.Count())
}

func TestContainsAllElements(t *testing.T) {

	var stackItems []int

	seed := int64(2163)
	stackItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)
	stack := New[int]()

	for _, v := range stackItems {
		stack.Push(v)
	}

	for _, v := range stackItems {
		require.True(t, stack.Contains(v))
	}
}

func TestPeekOnEmptyStackPanics(t *testing.T) {
	// Should not matter how we size it
	stack := New(WithCapacity[int](20))

	require.Panics(t, func() { stack.Peek() })
}

func TestPushedElementsArePopped(t *testing.T) {

	var stackItems, stackItemsReverse []int

	seed := int64(2163)
	stackItems, _, stackItemsReverse, _ = util.CreateIntListData(util.DefaultCapacity, &seed)
	stack := New[int]()

	for _, v := range stackItems {
		stack.Push(v)
	}

	verifyStackState(t, stack, stackItems)

	for _, v := range stackItemsReverse {
		require.Equal(t, v, stack.Pop())
	}

	require.Equal(t, 0, stack.Count())
}

func TestAddRange(t *testing.T) {

	var stackItems, stackItemsAdditional, tempItems []int

	seed := int64(2163)
	stackItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)
	stackItemsAdditional, _, _, _ = util.CreateIntListData(4, &seed)

	t.Run("Add nil slice doesn't panic", func(t *testing.T) {
		stack := New[int]()

		require.NotPanics(t, func() { stack.AddRange(nil) })
	})

	t.Run("Add single item to empty stack", func(t *testing.T) {
		stack := New[int]()
		tempItems = []int{util.CreateRandInt(&seed)}
		stack.AddRange(tempItems)
		verifyStackState(t, stack, tempItems)
	})

	t.Run("Add multiple items to empty stack", func(t *testing.T) {
		stack := New[int]()

		stack.AddRange(stackItems)
		verifyStackState(t, stack, stackItems)
	})

	t.Run("Add single item to populated stack", func(t *testing.T) {
		stack := New[int]()
		tempItems = []int{util.CreateRandInt(&seed)}

		for i := 0; i < len(stackItems); i++ {
			stack.Push(stackItems[i])
		}

		stack.AddRange(tempItems)
		tempItems2 := make([]int, len(stackItems), len(stackItems)+1)
		copy(tempItems2, stackItems)
		tempItems2 = append(tempItems2, tempItems...)

		verifyStackState(t, stack, tempItems2)
	})

	t.Run("Add multiple items to populated stack", func(t *testing.T) {
		stack := New[int]()

		for i := 0; i < len(stackItems); i++ {
			stack.Push(stackItems[i])
		}

		stack.AddRange(stackItemsAdditional)
		tempItems2 := make([]int, len(stackItems), len(stackItems)+len(stackItemsAdditional))
		copy(tempItems2, stackItems)
		tempItems2 = append(tempItems2, stackItemsAdditional...)

		verifyStackState(t, stack, tempItems2)
	})

	t.Run("AddRange preserves at least initial capacity", func(t *testing.T) {
		stack := New(WithCapacity[int](32))
		require.Equal(t, stack.initialCapacity, len(stack.buffer))
		stack.AddRange(stackItems)
		verifyStackState(t, stack, stackItems)
		require.Equal(t, stack.initialCapacity, len(stack.buffer))
	})
}

func TestRemove(t *testing.T) {

	var stackItems []int

	seed := int64(2163)

	t.Run("Full stack", func(t *testing.T) {
		stackItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)
		stack := New[int]()

		for _, v := range stackItems {
			stack.Push(v)
		}

		verifyStackState(t, stack, stackItems)

		removeItem := 4
		require.True(t, stack.Remove(stackItems[removeItem]))

		expectedItems := make([]int, len(stackItems)-1)
		index := 0
		for i := 0; i < removeItem; i++ {
			expectedItems[index] = stackItems[i]
			index++
		}
		for i := removeItem + 1; i < len(stackItems); i++ {
			expectedItems[index] = stackItems[i]
			index++
		}

		verifyStackState(t, stack, expectedItems)
	})

	t.Run("Resized stack", func(t *testing.T) {
		stackItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity+4, &seed)
		stack := New(WithCapacity[int](17))

		for _, v := range stackItems {
			stack.Push(v)
		}

		verifyStackState(t, stack, stackItems)

		removeItem := 4
		require.True(t, stack.Remove(stackItems[removeItem]))

		expectedItems := make([]int, len(stackItems)-1)
		index := 0
		for i := 0; i < removeItem; i++ {
			expectedItems[index] = stackItems[i]
			index++
		}
		for i := removeItem + 1; i < len(stackItems); i++ {
			expectedItems[index] = stackItems[i]
			index++
		}

		verifyStackState(t, stack, expectedItems)
	})
}

func TestTrimExcess(t *testing.T) {

	stack := generateIntStack(20)
	orginalLength := stack.length()
	originalSize := stack.size
	orignalCapacity := stack.capacity()
	shrinkBy := 5

	for i := 0; i < shrinkBy; i++ {
		stack.Pop()
	}

	require.Equal(t, stack.Count(), originalSize-shrinkBy)
	require.Equal(t, stack.length(), orginalLength)
	require.Equal(t, stack.capacity(), orignalCapacity)

	stack.TrimExcess()

	require.Equal(t, stack.Count(), originalSize-shrinkBy)
	require.Equal(t, stack.length(), originalSize-shrinkBy)
	require.Equal(t, stack.capacity(), originalSize-shrinkBy)
}

func TestTryStackOperations(t *testing.T) {

	seed := int64(2163)

	t.Run("TryPeek empty stack returns false", func(t *testing.T) {

		stack := New[int]()
		_, ok := stack.TryPeek()
		require.False(t, ok)
	})

	t.Run("TryPop empty stack returns false", func(t *testing.T) {
		stack := New[int]()
		_, ok := stack.TryPop()
		require.False(t, ok)
	})

	t.Run("TryPeek stack with value returns true and value", func(t *testing.T) {
		stack := New[int]()
		expected := util.CreateRandInt(&seed)
		stack.Push(expected)
		actual, ok := stack.TryPeek()
		require.True(t, ok)
		require.Equal(t, expected, actual)
		require.Equal(t, stack.Count(), 1)
	})

	t.Run("TryPop stack with value returns true and value", func(t *testing.T) {
		stack := New[int]()
		expected := util.CreateRandInt(&seed)
		stack.Push(expected)
		actual, ok := stack.TryPop()
		require.True(t, ok)
		require.Equal(t, expected, actual)
		require.Equal(t, stack.Count(), 0)
	})
}

func TestUnsafe(t *testing.T) {

	t.Run("GetVersion", func(t *testing.T) {
		s := New[int]()

		for i := 0; i < 10; i++ {
			s.Add(1)
			require.Equal(t, s.version, util.GetVersion[int](s))
		}
	})

	t.Run("GetLock", func(t *testing.T) {
		s := New[int]()

		ptrMutex := util.GetLock[int](s)

		require.Same(t, s.lock, ptrMutex)
	})
}

func TestThreadSafety(t *testing.T) {

	seed := int64(2163)
	itemsPerThread := 1024
	items1 := util.CreateSingleIntListData(itemsPerThread, &seed)
	items2 := util.CreateSingleIntListData(itemsPerThread, &seed)
	itemsCombined := make([]int, len(items1)+len(items2))
	copy(itemsCombined, items1)
	copy(itemsCombined[len(items1):], items2)

	t.Run("Parallel Push", func(t *testing.T) {
		s := New(WithThreadSafe[int]())
		wg := sync.WaitGroup{}
		wg.Add(2)

		pushFunc := func(stk *Stack[int], slc []int, w *sync.WaitGroup) {
			for _, v := range slc {
				stk.Push(v)
			}
			w.Done()
		}

		go pushFunc(s, items1, &wg)
		go pushFunc(s, items2, &wg)
		wg.Wait()
		require.ElementsMatch(t, itemsCombined, s.ToSlice())
	})

	t.Run("Parallel Pop", func(t *testing.T) {
		s := New(WithThreadSafe[int]())
		s.AddRange(itemsCombined)
		popped1 := make([]int, itemsPerThread)
		popped2 := make([]int, itemsPerThread)
		wg := sync.WaitGroup{}
		wg.Add(2)

		popFunc := func(stk *Stack[int], slc []int, w *sync.WaitGroup) {
			for i := 0; i < itemsPerThread; i++ {
				slc[i] = stk.Pop()
			}
			w.Done()
		}

		go popFunc(s, popped1, &wg)
		go popFunc(s, popped2, &wg)
		wg.Wait()

		popCombined := make([]int, len(itemsCombined))
		copy(popCombined, popped1)
		copy(popCombined[len(items1):], popped2)
		require.ElementsMatch(t, itemsCombined, popCombined)
	})
}

func verifyStackState[T any](t *testing.T, stack *Stack[T], expectedItems []T) {
	require.Equal(t, stack.Count(), len(expectedItems))
	require.ElementsMatch(t, stack.ToSlice(), expectedItems)
}

func benchmarkPush(s *Stack[int], data []int) {
	for _, v := range data {
		s.Push(v)
	}
}

func benchmarkPop(s *Stack[int], nitems int) {
	for i := 0; i < nitems; i++ {
		s.Pop()
	}
}

func BenchmarkStack(b *testing.B) {

	seed := int64(2163)
	data := make(map[int][]int, 4)
	elements := []int{100, 1000, 10000, 100000}
	list := util.CreateSingleIntListData(100000, &seed)

	for _, elem := range elements {
		data[elem] = list[:elem]
	}

	var s *Stack[int]

	for z := 0; z <= 1; z++ {
		threadsafe := z == 1

		for _, elems := range elements {
			b.Run(fmt.Sprintf("Stack-Push-%d-%s-NoPresize-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					if threadsafe {
						s = New(WithThreadSafe[int]())
					} else {
						s = New[int]()
					}
					b.StartTimer()
					benchmarkPush(s, data[elems])
				}
			})
		}

		for _, elems := range elements {
			b.Run(fmt.Sprintf("Stack-Push-%d-%s-Presize-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					if threadsafe {
						s = New[int](WithCapacity[int](elems), WithThreadSafe[int]())
					} else {
						s = New[int](WithCapacity[int](elems))
					}
					b.StartTimer()
					benchmarkPush(s, data[elems])
				}
			})
		}

		for _, elems := range elements {
			b.Run(fmt.Sprintf("Stack-Pop-%d-%s-NA-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					if threadsafe {
						s = New[int](WithCapacity[int](elems), WithThreadSafe[int]())
					} else {
						s = New[int](WithCapacity[int](elems))
					}
					s.AddRange(data[elems])
					b.StartTimer()
					benchmarkPop(s, elems)
				}
			})
		}
	}

	for _, elems := range elements {
		b.Run(fmt.Sprintf("Stack-Sort-%d-NA-NA-NA", elems), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				s = New[int](WithCapacity[int](elems))
				s.AddRange(data[elems])
				b.StartTimer()
				s.Sort()
			}
		})
	}

	for _, elems := range elements {
		b.Run(fmt.Sprintf("Stack-Contains-%d-NA-NA-NoConcurrent", elems), func(b *testing.B) {
			s = New[int](WithCapacity[int](elems))
			s.AddRange(data[elems])
			lookup := make([]int, elems)
			copy(lookup, data[elems])
			rand.Shuffle(elems, func(i, j int) {
				lookup[i], lookup[j] = lookup[j], lookup[i]
			})

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Contains(lookup[i%elems])
			}
		})
	}

	elems := 100000
	b.Run(fmt.Sprintf("Stack-Contains-%d-NA-NA-Concurrent", elems), func(b *testing.B) {
		s = New[int](WithCapacity[int](elems), WithConcurrent[int]())
		s.AddRange(data[elems])
		lookup := make([]int, elems)
		copy(lookup, data[elems])
		rand.Shuffle(elems, func(i, j int) {
			lookup[i], lookup[j] = lookup[j], lookup[i]
		})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Contains(lookup[i%elems])
		}
	})

	for _, elems := range elements {
		s = New[int](WithCapacity[int](elems))
		s.AddRange(data[elems])
		b.ResetTimer()
		b.Run(fmt.Sprintf("Stack-Min-%d-NA-NA-NoConcurrent", elems), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s.Min()
			}
		})

		b.Run(fmt.Sprintf("Stack-Max-%d-NA-NA-NoConcurrent", elems), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s.Max()
			}
		})
	}

	elems = 100000
	s = New[int](WithCapacity[int](elems), WithConcurrent[int]())
	s.AddRange(data[elems])
	b.ResetTimer()
	b.Run(fmt.Sprintf("Stack-Min-%d-NA-NA-Concurrent", elems), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.Min()
		}
	})

	b.Run(fmt.Sprintf("Stack-Max-%d-NA-NA-Concurrent", elems), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.Max()
		}
	})

}
