package util

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/messages"
)

// Default capacity for collections (where capacity matters).
const DefaultCapacity = 16

// Point (slice length) at which some slice operations switch to concurrent.
const concurrentThreshold = 65536

// Default predicate function for anywhere a predicate is required
// but all elements should be included.
func DefaultPredicate[T any](T) bool { return true }

// Reverse order of slice elements
// As per .NET Array.Reverse(Array).
func Reverse[T any](slc []T) []T {

	return ReverseSubset(slc, 0, len(slc))
}

// Reverse a subset of slice elements
// As per .NET Array.Reverse(Array, Int32, Int32).
func ReverseSubset[T any](slc []T, start, length int) []T {
	end := start + length - 1

	for start < end {
		slc[start], slc[end] = slc[end], slc[start]
		start++
		end--
	}

	return slc
}

// As per .NET Array.Copy(Array, Int32, Array, Int32, Int32).
func PartialCopy[T any](source []T, sourceIndex int, dest []T, destIndex int, length int) {
	for i := sourceIndex; i < sourceIndex+length; i++ {
		dest[destIndex] = source[i]
		destIndex++
	}
}

// As per .NET Array.Copy(Array, Int32, Array, Int32, Int32)
// but call provided DeepCopy function on every element.
func PartialCopyDeep[T any](source []T, sourceIndex int, dest []T, destIndex int, length int, f functions.DeepCopyFunc[T]) {
	for index := sourceIndex; index < sourceIndex+length; index++ {
		if f != nil {
			dest[destIndex] = f(source[index])
		} else {
			dest[destIndex] = source[index]
		}
		destIndex++
	}
}

// Get index in slice of given value. Return -1 if no match.
func IndexOf[T any](slc []T, value T, compare functions.ComparerFunc[T], concurrent bool) (index int) {
	l := len(slc)
	if l == 0 {
		return -1
	}
	if !concurrent || l < concurrentThreshold {
		return indexOf[T](slc, value, compare)
	} else {
		return indexOfConcurrent[T](slc, value, compare, runtime.NumCPU())
	}

}

func indexOf[T any](slc []T, value T, compare functions.ComparerFunc[T]) (index int) {
	index = -1

	for i := 0; i < len(slc); i++ {
		if compare(slc[i], value) == 0 {
			index = i
			break
		}
	}

	return
}

func indexOfConcurrent[T any](data []T, value T, compare functions.ComparerFunc[T], numChunks int) (index int) {

	chunkSize := (len(data) / numChunks) + 1
	indexChan := make(chan int)

	// Process
	var wg sync.WaitGroup
	for i := 0; i < numChunks; i++ {
		wg.Add(1)
		go func(chunkNum, chunkSize int, val T, indexChan chan int) {
			defer wg.Done()
			startIndex := chunkNum * chunkSize
			endIndex := startIndex + chunkSize

			if endIndex > len(data) {
				endIndex = len(data)
			}

			for i, e := range data[startIndex:endIndex] {

				if compare(e, val) == 0 {
					indexChan <- i + startIndex
					return
				}
			}

		}(i, chunkSize, value, indexChan)
	}

	resultsChan := make(chan int)
	defer close(resultsChan)

	go func(indexChan, resultsChan chan int) {

		resultIndex := -1
		for num := range indexChan {
			if num != -1 && resultIndex < num {
				resultIndex = num
			}
		}

		resultsChan <- resultIndex
	}(indexChan, resultsChan)

	wg.Wait()
	close(indexChan)
	return <-resultsChan
}

// Get last index in slice of given value. Return -1 if no match.
func LastIndexOf[T any](slc []T, value T, compare functions.ComparerFunc[T], concurrent bool) (index int) {
	l := len(slc)
	if l == 0 {
		return -1
	}
	if !concurrent || l < concurrentThreshold {
		return lastIndexOf[T](slc, value, compare)
	} else {
		return lastIndexOfConcurrent[T](slc, value, compare, runtime.NumCPU())
	}

}

func lastIndexOf[T any](slc []T, value T, compare functions.ComparerFunc[T]) (index int) {
	index = -1

	for i := len(slc) - 1; i >= 0; i-- {
		if compare(slc[i], value) == 0 {
			index = i
			break
		}
	}

	return
}

func lastIndexOfConcurrent[T any](data []T, value T, compare functions.ComparerFunc[T], numChunks int) (index int) {

	chunkSize := (len(data) / numChunks) + 1
	indexChan := make(chan int)

	// Process
	var wg sync.WaitGroup
	for i := 0; i < numChunks; i++ {
		wg.Add(1)
		go func(chunkNum, chunkSize int, val T, indexChan chan int) {
			defer wg.Done()
			startIndex := chunkNum * chunkSize
			endIndex := startIndex + chunkSize

			if endIndex > len(data) {
				endIndex = len(data)
			}

			for i := endIndex - 1; i >= startIndex; i-- {

				if compare(data[i], val) == 0 {
					indexChan <- i
					return
				}
			}

		}(i, chunkSize, value, indexChan)
	}

	resultsChan := make(chan int)
	defer close(resultsChan)

	go func(indexChan, resultsChan chan int) {

		resultIndex := -1
		for num := range indexChan {
			if num != -1 && resultIndex < num {
				resultIndex = num
			}
		}

		resultsChan <- resultIndex
	}(indexChan, resultsChan)

	wg.Wait()
	close(indexChan)
	return <-resultsChan
}

func Iif[T any](pred bool, trueVal T, falseVal T) T {
	// Cannot use function calls as args to this function,
	// because both calls are evaluated first
	// which will likely have undesirable side effects.
	if pred {
		return trueVal
	}

	return falseVal
}

func DefaultDeepCopy[T any](value T) T {
	return value
}

func DeepCopy[T any](value T, f functions.DeepCopyFunc[T]) T {
	if f == nil {
		return value
	}

	return f(value)
}

func DeepCopySlice[T any](dest, src []T, f functions.DeepCopyFunc[T]) {
	// Same semantics as built-in copy function
	for i := 0; i < len(dest) && i < len(src); i++ {
		dest[i] = DeepCopy(src[i], f)
	}
}

func Min[T any](slc []T, compare functions.ComparerFunc[T], concurrent bool) T {
	l := len(slc)
	if l == 0 {
		panic(messages.AGG_SLICE_EMPTY)
	}
	if !concurrent || l < concurrentThreshold {
		return min(slc, compare)
	} else {
		return getMinOrMaxConcurrent[T](slc, compare, min[T], runtime.NumCPU())
	}
}

func Max[T any](slc []T, compare functions.ComparerFunc[T], concurrent bool) T {
	l := len(slc)
	if l == 0 {
		panic(messages.AGG_SLICE_EMPTY)
	}
	if !concurrent || l < concurrentThreshold {
		return max(slc, compare)
	} else {
		return getMinOrMaxConcurrent[T](slc, compare, max[T], runtime.NumCPU())
	}
}

// min value of slice O(n).
func min[T any](slc []T, compare functions.ComparerFunc[T]) T {
	var m T
	for i, e := range slc {
		if i == 0 || compare(e, m) < 0 {
			m = e
		}
	}

	return m
}

// max value of slice O(n).
func max[T any](slc []T, compare functions.ComparerFunc[T]) T {
	var m T
	for i, e := range slc {
		if i == 0 || compare(e, m) > 0 {
			m = e
		}
	}

	return m
}

func getMinOrMaxConcurrent[T any](data []T, compare functions.ComparerFunc[T], boundfn func([]T, functions.ComparerFunc[T]) T, numChunks int) T {

	numChan := make(chan T)
	chunkSize := (len(data) / numChunks) + 1

	// Process
	var wg sync.WaitGroup
	for i := 0; i < numChunks; i++ {
		wg.Add(1)
		go func(i, chunkSize int, numChan chan T) {
			startIndex := i * chunkSize
			endIndex := startIndex + chunkSize

			if endIndex > len(data) {
				endIndex = len(data)
			}

			numChan <- boundfn(data[startIndex:endIndex], compare)
			wg.Done()
		}(i, chunkSize, numChan)
	}

	// Collect results
	resultsChan := make(chan T)
	defer close(resultsChan)

	go func(numChan, resultsChan chan T) {
		arr := make([]T, 0)

		for num := range numChan {
			arr = append(arr, num)
		}

		resultsChan <- boundfn(arr, compare)
	}(numChan, resultsChan)

	wg.Wait()
	close(numChan) // needed so results routine can return

	return <-resultsChan
}

/*
The following use pointer arithmetic to access memebers of collection
types from other collections without having to expose public methods
for things that should remain private between packages in this module.

Hopefully there should be no weirdness around struct member alignment
on pretty much most processor architectures, since both bits of data
being accessed are the size of a CPU word.

Tests on each collection verify the action of these functions.
*/

// Copy of internal representation of interface{}.
type eface struct {
	typ, val unsafe.Pointer
}

func GetVersion[T any](c collections.Collection[T]) int {
	// Expects version to be the first member of the collection struct
	return *(*int)((*eface)(unsafe.Pointer(&c)).val)
}

func GetLock[T any](c collections.Collection[T]) *sync.RWMutex {
	// Expects lock to be the second member of the version struct, following version
	return *(**sync.RWMutex)(unsafe.Add((*eface)(unsafe.Pointer(&c)).val, intSizeBytes))
}
