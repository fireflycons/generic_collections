package hashset

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/fireflycons/generic_collections/lists/dlist"
	"github.com/fireflycons/generic_collections/sets/orderedset"
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

	t.Run("With key capacity", func(t *testing.T) {
		require.NotPanics(t, func() { New(WithCapacity[int](10)) })
	})

	t.Run("With negative key capacity panics", func(t *testing.T) {
		require.Panics(t, func() { New(WithCapacity[int](-10)) })
	})

	t.Run("With bucket capacity", func(t *testing.T) {
		capacity := 10
		set := New(WithHashBucketCapacity[int](capacity))
		require.Equal(t, set.bucketCapacity, capacity)
	})

	t.Run("With negative bucket capacity panics", func(t *testing.T) {
		capacity := -10
		require.Panics(t, func() { New(WithHashBucketCapacity[int](capacity)) })
	})
}

func TestClear(t *testing.T) {

	t.Run("Clear empties the set", func(t *testing.T) {
		var setItems []int
		var set *HashSet[int]
		seed := int64(2163)

		setItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)
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
			}),
			WithHasher(func(s1 strct) uintptr { return uintptr(42) }))

		i := 1
		set.Add(strct{ptr: &i})
		set.Clear()
		require.ElementsMatch(t, set.ToSlice(), []strct{})
	})
}

func TestAddItems(t *testing.T) {

	var setItems, tempItems []int
	seed := int64(2163)
	setItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

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
		s := New[int]()

		for _, v := range setItems {
			s.Add(v)
		}

		tempItems = make([]int, len(setItems))

		for i := 0; i < len(setItems); i++ {
			require.True(t, s.Remove(setItems[i]))
			tempItems[i] = setItems[i]
		}

		require.ElementsMatch(t, setItems, tempItems)
		require.Equal(t, 0, s.Count())
	})

	t.Run("Add range", func(t *testing.T) {
		s := New[int]()
		s.AddRange(setItems)
		require.Equal(t, util.DefaultCapacity, s.Count())
	})

	t.Run("ToSlice", func(t *testing.T) {
		s := New[int]()
		s.AddRange(setItems)
		require.ElementsMatch(t, setItems, s.ToSlice())
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
		// Tests that the hash algorithm for time.Time is suitable
		var set *HashSet[time.Time]
		arraySize := 16384
		seed := int64(2163)
		items := util.CreateTimeListData(arraySize, &seed)

		set = New(WithCapacity[time.Time](arraySize))

		set.AddRange(items)
		require.ElementsMatch(t, items, set.ToSlice())
		require.Equal(t, 0, set.collisionCount)
	})
}

func TestAddItemsWithHashCollisions(t *testing.T) {

	// In these tests, we use a "rubbish hash" that yields only four possible keys.

	var setItems, tempItems []int
	seed := int64(2163)
	setItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Add single item", func(t *testing.T) {
		s := New(WithHasher(func(v int) uintptr { return uintptr(v % 4) }))
		s.Add(0)
		require.Equal(t, 1, s.Count())
	})

	t.Run("Add duplicate item", func(t *testing.T) {
		s := New(WithHasher(func(v int) uintptr { return uintptr(v % 4) }))

		require.True(t, s.Add(0))
		require.False(t, s.Add(0))
		require.Equal(t, 1, s.Count())
	})

	t.Run("Add unique items", func(t *testing.T) {
		s := New(WithHasher(func(v int) uintptr { return uintptr(v % 4) }))

		require.True(t, s.Add(0))
		require.True(t, s.Add(1))
		require.Equal(t, 2, s.Count())
	})

	t.Run("Add unique items with same hash", func(t *testing.T) {
		s := New(WithHasher(func(v int) uintptr { return uintptr(v % 4) }))

		require.True(t, s.Add(1))
		require.True(t, s.Add(5))
		require.Equal(t, 2, s.Count())
		require.Equal(t, 1, s.collisionCount)
	})

	t.Run("Add many items and remove them", func(t *testing.T) {
		s := New(WithHasher(func(v int) uintptr { return uintptr(v % 4) }))

		for _, v := range setItems {
			s.Add(v)
		}

		require.Greater(t, s.collisionCount, 0)

		tempItems = make([]int, len(setItems))

		for i := 0; i < len(setItems); i++ {
			require.True(t, s.Remove(setItems[i]))
			tempItems[i] = setItems[i]
		}

		require.ElementsMatch(t, setItems, tempItems)
		require.Equal(t, 0, s.Count())
		require.Equal(t, 0, s.collisionCount)
		require.Equal(t, 0, len(s.buffer))
	})

	t.Run("Add range", func(t *testing.T) {
		s := New(WithHasher(func(v int) uintptr { return uintptr(v % 4) }))
		s.AddRange(setItems)
		require.Equal(t, util.DefaultCapacity, s.Count())
	})

	t.Run("ToSlice", func(t *testing.T) {
		s := New(WithHasher(func(v int) uintptr { return uintptr(v % 4) }))
		s.AddRange(setItems)
		require.ElementsMatch(t, setItems, s.ToSlice())
	})
}

func TestContains(t *testing.T) {

	var setItems []int
	seed := int64(2163)
	setItems = util.CreateSerialIntListData(util.DefaultCapacity, &seed)

	var set = New[int]()
	set.AddRange(setItems)
	require.ElementsMatch(t, setItems, set.ToSlice())

	require.True(t, set.Contains(setItems[4]))
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

	t.Run("Intersect yields items only on both where 2 is orderedset", func(t *testing.T) {
		set1 := New[int]()
		set2 := orderedset.New[int]()

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

	t.Run("Union yields all unique items from both where 2 is ordered set", func(t *testing.T) {
		set1 := New[int]()
		set2 := orderedset.New[int]()

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

	t.Run("Difference yields all items from 1 that are not in 2 where 2 is ordered set", func(t *testing.T) {
		set1 := New[int]()
		set2 := orderedset.New[int]()

		set1.AddRange([]int{1, 2, 3, 4})
		set2.AddRange([]int{3, 4, 5, 6})
		difference := set1.Difference(set2)
		expected := []int{1, 2}

		require.ElementsMatch(t, expected, difference.ToSlice())
	})
}

func TestTime(t *testing.T) {
	var set *HashSet[time.Time]
	arraySize := 16
	seed := int64(2163)
	items := util.CreateTimeListData(arraySize, &seed)

	set = New[time.Time]()

	set.AddRange(items)
	require.ElementsMatch(t, items, set.ToSlice())
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

		addFunc := func(set *HashSet[int], slc []int, w *sync.WaitGroup) {
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

		removeFunc := func(stk *HashSet[int], slc []int, w *sync.WaitGroup) {
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

func TestUnsafe(t *testing.T) {

	t.Run("GetVersion", func(t *testing.T) {
		s := New(WithThreadSafe[int]())

		for i := 0; i < 10; i++ {
			s.Add(1)
			require.Equal(t, s.version, util.GetVersion[int](s))
		}
	})

	t.Run("GetLock", func(t *testing.T) {
		s := New(WithThreadSafe[int]())

		ptrMutex := util.GetLock[int](s)

		require.Same(t, s.lock, ptrMutex)
	})
}

func benchmarkAdd(s *HashSet[int], data []int) {
	for _, v := range data {
		s.Add(v)
	}
}

func benchmarkRemove(s *HashSet[int], data []int) {
	for _, v := range data {
		s.Remove(v)
	}
}

func BenchmarkHashSet(b *testing.B) {

	seed := int64(2163)
	data := make(map[int][]int, 4)
	elements := []int{100, 1000, 10000, 100000}
	list := util.CreateSingleIntListData(100000, &seed)

	for _, elem := range elements {
		data[elem] = list[:elem]
	}

	var s *HashSet[int]

	for z := 0; z <= 1; z++ {
		threadsafe := z == 1

		for _, elems := range elements {
			b.Run(fmt.Sprintf("Set-Add-%d-%s-NoPresize-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
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
			b.Run(fmt.Sprintf("Set-Add-%d-%s-Presize-NA", elems, util.Iif(threadsafe, "ThreadSafe", "NoThreadSafe")), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					if threadsafe {
						s = New[int](WithCapacity[int](elems), WithThreadSafe[int]())
					} else {
						s = New[int](WithCapacity[int](elems))
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
						s = New[int](WithCapacity[int](elems), WithThreadSafe[int]())
					} else {
						s = New[int](WithCapacity[int](elems))
					}
					s.AddRange(data[elems])
					b.StartTimer()
					benchmarkRemove(s, data[elems])
				}
			})
		}
	}

	for _, elems := range elements {
		b.Run(fmt.Sprintf("Set-Contains-%d-NA-NA-NA", elems), func(b *testing.B) {
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

	for _, elems := range elements {
		s = New[int](WithCapacity[int](elems))
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
