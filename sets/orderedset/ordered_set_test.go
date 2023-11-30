package orderedset

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/fireflycons/generic_collections/lists/dlist"
	"github.com/fireflycons/generic_collections/sets/hashset"
	"github.com/stretchr/testify/require"
)

func TestConstructor(t *testing.T) {

	t.Run("With comparer", func(t *testing.T) {
		magic := 42
		comp := func(v1, v2 int) int { return magic }
		set := New(WithComparer(comp))

		require.Equal(t, magic, set.compare(1, 0))
	})

	t.Run("With nil comparer panics", func(t *testing.T) {
		var comp func(v1, v2 int) int
		require.Panics(t, func() { New(WithComparer(comp)) })
	})

}

func TestClear(t *testing.T) {

	t.Run("Clear empties the queue", func(t *testing.T) {
		var setItems []int
		var set *OrderedSet[int]
		seed := int64(2163)

		setItems = util.CreateSingleIntListData(util.DefaultCapacity, &seed)
		set = New[int]()

		for _, v := range setItems {
			set.Add(v)
		}

		require.ElementsMatch(t, set.ToSlice(), setItems)

		set.Clear()
		require.ElementsMatch(t, set.ToSlice(), []int{})
	})

	t.Run("Clear with struct type", func(t *testing.T) {
		type strct struct {
			ptr *int
		}

		set := New(
			WithComparer(func(s1, s2 strct) int {
				return 0
			}))

		i := 1
		set.Add(strct{ptr: &i})
		set.Clear()
		require.ElementsMatch(t, set.ToSlice(), []strct{})
	})
}

func TestAddItems(t *testing.T) {

	var setItems, tempItems []int
	seed := int64(2163)
	setItems = util.CreateSingleIntListData(util.DefaultCapacity, &seed)

	_ = setItems
	_ = tempItems

	t.Run("Add single item", func(t *testing.T) {
		s := New[int]()

		s.Add(0)

		require.Equal(t, 1, s.Count())
	})

	t.Run("Add duplicate item", func(t *testing.T) {
		s := New[int]()

		require.True(t, s.Add(0))
		require.False(t, s.Add(0))
		require.Equal(t, 1, s.Count())
	})

	t.Run("Add unique items", func(t *testing.T) {
		s := New[int]()

		require.True(t, s.Add(0))
		require.True(t, s.Add(1))
		require.Equal(t, 2, s.Count())
	})

	t.Run("Add many items and remove them", func(t *testing.T) {
		var manyItems []int
		var newseed = int64(1)
		manyItems, _, _, _ = util.CreateIntListData(1024, &newseed)
		s := New[int]()

		s.AddRange(manyItems)
		tempItems = make([]int, len(manyItems))
		copy(tempItems, manyItems)
		sort.Ints(tempItems)
		require.Equal(t, tempItems, s.ToSlice())

		tempItems = make([]int, len(manyItems))

		for i := 0; i < len(manyItems); i++ {
			require.True(t, s.Remove(manyItems[i]))
			tempItems[i] = manyItems[i]
		}

		require.ElementsMatch(t, manyItems, tempItems)
		require.Equal(t, 0, s.Count())
	})

	t.Run("ToSlice returns elements sorted ascending", func(t *testing.T) {
		s := New[int]()
		for _, v := range setItems {
			s.Add(v)
		}

		tempItems = make([]int, len(setItems))
		copy(tempItems, setItems)
		sort.Ints(tempItems)
		require.Equal(t, tempItems, s.ToSlice())
	})
	t.Run("AddRange", func(t *testing.T) {
		s := New[int]()
		s.AddRange(setItems)

		tempItems = make([]int, len(setItems))
		copy(tempItems, setItems)
		sort.Ints(tempItems)
		require.Equal(t, tempItems, s.ToSlice())
	})

	t.Run("AddCollection to empty set", func(t *testing.T) {
		linkedList1 := dlist.New[int]()
		linkedList1.AddRange(setItems)

		s := New[int]()
		s.AddCollection(linkedList1)
		require.Equal(t, len(setItems), s.Count())
		require.ElementsMatch(t, setItems, s.ToSlice())
	})

	t.Run("Using time.Time", func(t *testing.T) {
		// Tests that the comparer func for time.Time is valid
		arraySize := 16
		seed := int64(2163)
		items := util.CreateTimeListData(arraySize, &seed)
		set := New[time.Time]()

		set.AddRange(items)
		expected := make([]time.Time, arraySize)
		copy(expected, items)
		sort.Slice(expected, func(i, j int) bool {
			return expected[i].Before(expected[j])
		})

		require.Equal(t, expected, set.ToSlice())
	})
}

func TestContains(t *testing.T) {

	var setItems []int
	seed := int64(2163)
	setItems = util.CreateSerialIntListData(util.DefaultCapacity, &seed)

	var set = New[int]()
	//set.AddRange(setItems)
	for _, v := range setItems {
		set.Add(v)
	}
	require.ElementsMatch(t, setItems, set.ToSlice())

	require.True(t, set.Contains(setItems[4]))
}

func TestClearWithPointers(t *testing.T) {

	var setItems, tempItems []int
	seed := int64(2163)
	setItems = util.CreateSingleIntListData(util.DefaultCapacity, &seed)
	s := New[*int]()

	tempItems = make([]int, len(setItems))
	copy(tempItems, setItems)

	for i := 0; i < len(setItems); i++ {
		s.Add(&setItems[i])
	}

	s.Clear()
	require.Equal(t, 0, s.Count())
	require.Equal(t, tempItems, setItems)
}

func TestIntersection(t *testing.T) {

	t.Run("Intersect empty sets yields empty set", func(t *testing.T) {
		set1 := New[int]()
		set2 := New[int]()

		intersection := set1.Intersection(set2)

		require.ElementsMatch(t, []int{}, intersection.ToSlice())
	})

	t.Run("Intersect yields items only on both", func(t *testing.T) {
		set1 := New[int]()
		set2 := New[int]()

		set1.AddRange([]int{1, 2, 3, 4})
		set2.AddRange([]int{3, 4, 5, 6})
		intersection := set1.Intersection(set2)
		expected := []int{3, 4}

		require.ElementsMatch(t, expected, intersection.ToSlice())
	})

	t.Run("Intersect yields items only on both where 2 is hash set", func(t *testing.T) {
		set1 := New[int]()
		set2 := hashset.New[int]()

		set1.AddRange([]int{1, 2, 3, 4})
		set2.AddRange([]int{3, 4, 5, 6})
		intersection := set1.Intersection(set2)
		expected := []int{3, 4}

		require.ElementsMatch(t, expected, intersection.ToSlice())
	})
}

func BenchmarkIntersection(b *testing.B) {
	set1 := New[int]()
	set2 := New[int]()
	elems := 100000
	seed := int64(0)
	data1 := util.CreateSerialIntListData(elems, &seed)
	seed /= 2
	data2 := util.CreateSerialIntListData(elems-1, &seed)
	rand.Shuffle(elems, func(i, j int) {
		data1[i], data1[j] = data1[j], data1[i]
	})
	rand.Shuffle(elems-1, func(i, j int) {
		data2[i], data2[j] = data2[j], data2[i]
	})

	set1.AddRange(data1)
	set2.AddRange(data2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		set1.Intersection(set2)
	}
}

func TestUnion(t *testing.T) {

	t.Run("Union empty sets yields empty set", func(t *testing.T) {
		set1 := New[int]()
		set2 := New[int]()

		union := set1.Union(set2)

		require.ElementsMatch(t, []int{}, union.ToSlice())
	})

	t.Run("Union yields all unique items from both", func(t *testing.T) {
		set1 := New[int]()
		set2 := New[int]()

		set1.AddRange([]int{1, 2, 3, 4})
		set2.AddRange([]int{3, 4, 5, 6})
		union := set1.Union(set2)
		expected := []int{1, 2, 3, 4, 5, 6}

		require.ElementsMatch(t, expected, union.ToSlice())
	})

	t.Run("Union yields all unique items from both where 2 is hash set", func(t *testing.T) {
		set1 := New[int]()
		set2 := hashset.New[int]()

		set1.AddRange([]int{1, 2, 3, 4})
		set2.AddRange([]int{3, 4, 5, 6})
		union := set1.Union(set2)
		expected := []int{1, 2, 3, 4, 5, 6}

		require.ElementsMatch(t, expected, union.ToSlice())
	})
}

func TestDifference(t *testing.T) {

	t.Run("Difference empty sets yields empty set", func(t *testing.T) {
		set1 := New[int]()
		set2 := New[int]()

		difference := set1.Difference(set2)

		require.ElementsMatch(t, []int{}, difference.ToSlice())
	})

	t.Run("Difference yields all items from 1 that are not in 2", func(t *testing.T) {
		set1 := New[int]()
		set2 := New[int]()

		set1.AddRange([]int{1, 2, 3, 4})
		set2.AddRange([]int{3, 4, 5, 6})
		difference := set1.Difference(set2)
		expected := []int{1, 2}

		require.ElementsMatch(t, expected, difference.ToSlice())
	})

	t.Run("Difference yields all items from 1 that are not in 2 where 2 is hash set", func(t *testing.T) {
		set1 := New[int]()
		set2 := hashset.New[int]()

		set1.AddRange([]int{1, 2, 3, 4})
		set2.AddRange([]int{3, 4, 5, 6})
		difference := set1.Difference(set2)
		expected := []int{1, 2}

		require.ElementsMatch(t, expected, difference.ToSlice())
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

	t.Run("Parallel Add", func(t *testing.T) {
		s := New(WithThreadSafe[int]())
		wg := sync.WaitGroup{}
		wg.Add(2)

		addFunc := func(set *OrderedSet[int], slc []int, w *sync.WaitGroup) {
			for _, v := range slc {
				set.Add(v)
			}
			w.Done()
		}

		go addFunc(s, items1, &wg)
		go addFunc(s, items2, &wg)
		wg.Wait()

		require.ElementsMatch(t, itemsCombined, s.ToSlice())
	})

	t.Run("Parallel Remove", func(t *testing.T) {
		s := New(WithThreadSafe[int]())
		s.AddRange(itemsCombined)
		wg := sync.WaitGroup{}
		wg.Add(2)

		removeFunc := func(stk *OrderedSet[int], slc []int, w *sync.WaitGroup) {
			for _, v := range slc {
				stk.Remove(v)
			}
			w.Done()
		}

		go removeFunc(s, items1, &wg)
		go removeFunc(s, items2, &wg)
		wg.Wait()

		require.Equal(t, 0, s.Count())
	})
}

func benchmarkAdd(s *OrderedSet[int], data []int) {
	for _, v := range data {
		s.Add(v)
	}
}

func benchmarkRemove(s *OrderedSet[int], data []int) {
	for _, v := range data {
		s.Remove(v)
	}
}

func BenchmarkOrderedSet(b *testing.B) {

	seed := int64(2163)
	data := make(map[int][]int, 4)
	elements := []int{100, 1000, 10000, 100000}
	list := util.CreateSingleIntListData(100000, &seed)

	for _, elem := range elements {
		data[elem] = list[:elem]
	}

	var s *OrderedSet[int]

	for z := 0; z <= 1; z++ {
		threadsafe := z == 1

		for _, elems := range elements {
			b.Run(fmt.Sprintf("Set-Add-%d-%s-NA-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					if threadsafe {
						s = New(WithThreadSafe[int]())
					} else {
						s = New[int]()
					}
					b.StartTimer()
					benchmarkAdd(s, data[elems])
				}
			})
		}

		for _, elems := range elements {
			b.Run(fmt.Sprintf("Set-Remove-%d-%s-NA-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					if threadsafe {
						s = New[int](WithThreadSafe[int]())
					} else {
						s = New[int]()
					}
					s.AddRange(data[elems])
					b.StartTimer()
					benchmarkRemove(s, data[elems])
				}
			})
		}
	}

	for _, elems := range elements {
		s = New[int]()
		s.AddRange(data[elems])
		lookup := make([]int, elems)
		copy(lookup, data[elems])
		rand.Shuffle(elems, func(i, j int) {
			lookup[i], lookup[j] = lookup[j], lookup[i]
		})

		b.ResetTimer()

		b.Run(fmt.Sprintf("Set-Contains-%d-NA-NA-NA", elems), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s.Contains(lookup[i%elems])
			}
		})
	}

	for _, elems := range elements {
		s = New[int]()
		s.AddRange(data[elems])
		b.ResetTimer()
		b.Run(fmt.Sprintf("Set-Min-%d-NA-NA-NA", elems), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s.Min()
			}
		})

		b.Run(fmt.Sprintf("Set-Max-%d-NA-NA-NA", elems), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s.Max()
			}
		})
	}

}
