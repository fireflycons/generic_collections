/*
Package hashset provides a hash bucket backed collection that contains an unordered collection of unique values.
*/
package hashset

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/local"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/fireflycons/generic_collections/sets"
)

// Assert HashSet implements required interfaces.
var _ sets.Set[int] = (*HashSet[int])(nil)

// Capacity of initial hash buckets.
// If the hashing algorithm to generate keys is good enough
// then there should be few collisions.
// Ideally 99.9999% of values should generate unique hash keys.
const defaultBucketCapacity = 2

// Option function signature for HashSet contructor options.
type HashSetOptionFunc[T any] func(*HashSet[T])

// HashSet stores an unordered collection of unique elements.
type HashSet[T any] struct {
	version        int
	lock           *sync.RWMutex
	bucketCapacity int
	collisionCount int
	size           int
	capacity       int
	hasher         func(T) uintptr
	compare        functions.ComparerFunc[T]
	copy           functions.DeepCopyFunc[T]
	buffer         map[uintptr][]T
	concurrent     bool
	local.InternalImpl
}

// Constructs a new HashSet[T].
func New[T any](options ...HashSetOptionFunc[T]) *HashSet[T] {
	s := &HashSet[T]{}

	for _, o := range options {
		o(s)
	}

	if s.hasher == nil {
		// Will panic if T is not comparable
		s.setDefaultHasher()
	}

	if s.copy == nil {
		s.copy = util.DefaultDeepCopy[T]
	}

	if s.compare == nil {
		s.compare = util.GetDefaultComparer[T]()
	}

	if s.buffer == nil {
		s.buffer = make(map[uintptr][]T, util.DefaultCapacity)
		s.capacity = util.DefaultCapacity
	}

	if s.bucketCapacity == 0 {
		s.bucketCapacity = defaultBucketCapacity
	}

	return s
}

// Option function for New to make the collection thread-safe. Adds overhead.
func WithThreadSafe[T any]() HashSetOptionFunc[T] {
	return func(s *HashSet[T]) {
		s.lock = &sync.RWMutex{}
	}
}

// Option function to enable concurrency feature.
func WithConcurrent[T any]() HashSetOptionFunc[T] {
	return func(s *HashSet[T]) {
		s.concurrent = true
	}
}

// Option function for NewSet to provide an alternative hash function
// for the type of values stored in the set. This is required for any type
// that is not one of the supported types. HashSet creation will panic
// if the type T is not one of the supported types and a hasher is not
// provided.
//
//	// Creates a HashSet[int] with a custom hasher.
//	// Clearly _don't_ use a hash function like this :-)
//	set := NewSet(WithHasher(func(v int) uintptr { return uintptr(v % 4) }))
func WithHasher[T any](hasher functions.HashFunc[T]) HashSetOptionFunc[T] {
	return func(s *HashSet[T]) {
		s.hasher = hasher
	}
}

// Option function for NewSet to provide a comparer function for values of type T.
// Required if the element type is not numeric, bool, pointer or string.
func WithComparer[T any](comparer functions.ComparerFunc[T]) HashSetOptionFunc[T] {
	if comparer == nil {
		panic(messages.COMP_FN_NIL)
	}
	return func(s *HashSet[T]) {
		s.compare = comparer
	}
}

// Option func to provide a deep copy implementation for collection elements.
func WithDeepCopy[T any](copier functions.DeepCopyFunc[T]) HashSetOptionFunc[T] {
	// Can be nil
	return func(s *HashSet[T]) {
		s.copy = copier
	}
}

// Option function for NewSet to set the initial hash bucket capacity associated with a new hash key.
// The default capacity is 2, which should be sufficient for the default hashing algorithms.
func WithHashBucketCapacity[T any](bucketCapacity int) HashSetOptionFunc[T] {
	if bucketCapacity < 1 {
		panic(messages.HASH_BUCKET_SIZE_INVALID)
	}
	return func(s *HashSet[T]) {
		s.bucketCapacity = bucketCapacity
	}
}

// Option function for NewSet to set initial key capacity to
// something other than the default 16 elements. If you have
// some idea of hom many elements you will be storing, addition
// of values is approx 30% faster in a preallocated collection.
func WithCapacity[T any](capacity int) HashSetOptionFunc[T] {
	if capacity < 0 {
		panic(messages.NEGATIVE_CAPACITY)
	}
	return func(s *HashSet[T]) {
		s.buffer = make(map[uintptr][]T, capacity)
		s.capacity = capacity
	}
}

// AddCollection inserts the values of the given collection into this set.
// Values are added in the order defined by the other collection.
func (s *HashSet[T]) AddCollection(collection collections.Collection[T]) {

	s.AddRange(collection.ToSliceDeep())
}

// Clear removes all values from the set, restoring it to its initial capacity.
func (s *HashSet[T]) Clear() {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	s.buffer = make(map[uintptr][]T, max(s.bucketCapacity, util.DefaultCapacity))
	s.size = 0
	s.version++
}

// Add adds a value into the set. Returns true if the value was added;
// else false if the value already exists in the set.
func (s *HashSet[T]) Add(value T) bool {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	s.version++
	return s.add(value)
}

// AddRange adds a slice of values to the set.
func (s *HashSet[T]) AddRange(values []T) {

	if len(values) == 0 {
		return
	}

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	for _, v := range values {
		s.add(v)
	}

	s.version++
}

// Count returns the number of elements stored in the set.
func (s *HashSet[T]) Count() int {

	return s.size
}

// Contains returns true if the given value exists in the set; else false
//
// O(1) average and amortized. Up to O(n) if user hasher is poor.
func (s *HashSet[T]) Contains(value T) bool {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	return s.contains(s.hasher(value), value) >= 0
}

func (s *HashSet[T]) UnlockedContains(value T) bool {
	return s.contains(s.hasher(value), value) >= 0
}

// Get returns the collection element that matches the given value, or nil if it is not found.
// Useful if the set contains struct elements you want to modify in-place.
func (s *HashSet[T]) Get(value T) collections.Element[T] {
	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	hash := s.hasher(value)
	ind := s.contains(hash, value)

	if ind == -1 {
		return nil
	}

	return util.NewElementType[T](s, &s.buffer[hash][ind])
}

// IsEmpty returns true if the collection has no elements.
func (s *HashSet[T]) IsEmpty() bool {
	return s.size == 0
}

// ToSlice returns a copy of the set content as a slice.
func (s *HashSet[T]) ToSlice() []T {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	return s.toSlice(false)
}

// ToSliceDeep returns a copy of the set content as a slice using the provided [functions.DeepCopyFunc] if any.
func (s *HashSet[T]) ToSliceDeep() []T {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	return s.toSlice(true)
}

// Remove removes a value from the set.
//
// Returns true if the value was present and was removed;
// else false.
func (s *HashSet[T]) Remove(value T) bool {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	hash := s.hasher(value)
	index := s.contains(hash, value)
	if index == -1 {
		return false
	}

	if len(s.buffer[hash]) == 1 {
		delete(s.buffer, hash)
		//s.buffer[hash] = make([]T, 0, s.bucketCapacity)
	} else {
		// More than one value for this hash
		s.collisionCount--
		tmp := s.buffer[hash]
		tmp[index] = tmp[len(tmp)-1]
		s.buffer[hash] = tmp[:len(tmp)-1]
	}

	s.version++
	s.size--
	return true
}

// Type returns the type of this collection.
func (*HashSet[T]) Type() collections.CollectionType {
	return collections.COLLECTION_HASHSET
}

// Difference returns the difference between two sets.
// The new set consists of all elements that are in this set, but not other set.
//
// The argument can be any implementation of Set[T]. The result is a new HashSet with the same properties as this one.
// Items are shallow-copied.
func (s *HashSet[T]) Difference(other sets.Set[T]) sets.Set[T] {

	ol := util.GetLock[T](other)

	if ol != nil {
		ol.RLock()
		defer ol.RUnlock()
	}

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	result := s.makeEmptyCopy(s.capacity)

	if other.Count() == 0 {
		// Resultant set is a direct copy of this one
		for key, bucket := range s.buffer {
			b := make([]T, len(bucket))
			copy(b, bucket)
			result.buffer[key] = b
		}

		result.size = len(s.buffer)
		result.collisionCount = s.collisionCount
		return result
	}

	for _, bucket := range s.buffer {
		for _, value := range bucket {
			if !other.UnlockedContains(value) {
				result.add(value)
			}
		}
	}

	return result
}

// Intersection returns the intersection between two sets.
// The new set consists of all elements that are in both this set and the other.
//
// The argument can be any implementation of Set[T]. The result is a new HashSet with the same properties as this one.
// Items are shallow-copied.
func (s *HashSet[T]) Intersection(other sets.Set[T]) sets.Set[T] {

	if s.size == 0 || other.Count() == 0 {
		// No intersection if either set empty
		return New[T]()
	}

	ol := util.GetLock[T](other)
	if ol != nil {
		ol.RLock()
		defer ol.RUnlock()
	}

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	// It's much quicker to scan the smaller collection
	// and look up values in the larger one as lookup
	// is very fast in sets.
	var smaller, larger sets.Set[T]
	var smallerHS *HashSet[T]
	var smallerIsHS bool

	if s.size <= other.Count() {
		smaller = s
		larger = other
		smallerHS, smallerIsHS = smaller.(*HashSet[T])
	} else {
		smaller = other
		larger = s
	}

	result := s.makeEmptyCopy(smaller.Count())

	if smallerIsHS {
		for _, bucket := range smallerHS.buffer {
			for _, value := range bucket {
				if larger.UnlockedContains(value) {
					result.add(value)
				}
			}
		}

		return result
	}

	for _, value := range smaller.ToSlice() {
		if larger.UnlockedContains(value) {
			result.add(value)
		}
	}

	return result
}

// Union returns the union of two sets.
// The new set consists of all elements that are in both this and the other set.
//
// The argument can be any implementation of Set[T]. The result is a new HashSet with the same properties as this one.
// Items are shallow-copied.
func (s *HashSet[T]) Union(other sets.Set[T]) sets.Set[T] {

	ol := util.GetLock[T](other)

	if ol != nil {
		ol.RLock()
		defer ol.RUnlock()
	}

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	result := s.makeEmptyCopy(s.size + other.Count())
	result.AddCollection(s)
	result.AddCollection(other)
	return result
}

// returns hasbucket index of value if found.
func (s *HashSet[T]) contains(hash uintptr, value T) int {
	bucket, ok := s.buffer[hash]

	if !ok {
		return -1
	}

	for i, v := range bucket {
		if s.compare(v, value) == 0 {
			return i
		}
	}

	return -1
}

// String returns a string representation of container.
func (s *HashSet[T]) String() string {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	var values []string
	for _, value := range s.toSlice(false) {
		values = append(values, fmt.Sprintf("%v", value))
	}

	return "HashSet\n" + strings.Join(values, ", ")
}

func (s *HashSet[T]) toSlice(deepCopy bool) []T {
	slc := make([]T, s.size)
	i := 0

	for _, bucket := range s.buffer {
		for _, v := range bucket {
			if deepCopy {
				slc[i] = util.DeepCopy(v, s.copy)
			} else {
				slc[i] = v
			}
			i++
		}
	}

	return slc
}

func (s *HashSet[T]) add(value T) bool {

	var bucket []T
	hash := s.hasher(value)
	if s.contains(hash, value) > -1 {
		return false
	}

	bucket, ok := s.buffer[hash]

	if !ok {
		// hash bucket doesn't exist
		bucket = make([]T, 0, s.bucketCapacity)
	} else if len(bucket) > 0 {
		// Bucket exists and holds another value
		s.collisionCount++
	}

	bucket = append(bucket, value)
	s.buffer[hash] = bucket
	s.size++
	return true
}

func (s *HashSet[T]) makeEmptyCopy(capacity int) *HashSet[T] {
	other := &HashSet[T]{
		bucketCapacity: s.bucketCapacity,
		capacity:       capacity,
		hasher:         s.hasher,
		compare:        s.compare,
		copy:           s.copy,
		buffer:         make(map[uintptr][]T, capacity),
		concurrent:     s.concurrent,
	}

	if s.lock != nil {
		other.lock = &sync.RWMutex{}
	}

	return other
}
