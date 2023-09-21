/*
Package stack implements a slice-backed LIFO stack.
*/
package stack

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fireflycons/generic_collections/collections"
	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/local"
	"github.com/fireflycons/generic_collections/internal/messages"
	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/fireflycons/generic_collections/stacks"
)

// Assert Stack implements required interfaces
var _ stacks.Stack[int] = (*Stack[int])(nil)
var _ collections.ReverseIterable[int] = (*Stack[int])(nil)

const (
	growFactor  = 200
	minimumGrow = 4
)

type StackOptionFunc[T any] func(*Stack[T])

// Stack implements a last-in, first-out collection.
type Stack[T any] struct {
	version         int
	lock            *sync.RWMutex
	size            int
	initialCapacity int
	compare         functions.ComparerFunc[T]
	copy            functions.DeepCopyFunc[T]
	buffer          []T
	concurrent      bool

	local.InternalImpl
}

// New creates a new stack with intial capacity for 16 elements
func New[T any](options ...StackOptionFunc[T]) *Stack[T] {
	stack := &Stack[T]{
		initialCapacity: util.DefaultCapacity,
	}
	for _, o := range options {
		o(stack)
	}

	stack.buffer = make([]T, stack.initialCapacity)

	if stack.copy == nil {
		stack.copy = util.DefaultDeepCopy[T]
	}

	if stack.compare == nil {
		stack.compare = util.GetDefaultComparer[T]()
	}

	return stack
}

// Option function for New to make the collection thread-safe. Adds overhead.
func WithThreadSafe[T any]() StackOptionFunc[T] {
	return func(s *Stack[T]) {
		s.lock = &sync.RWMutex{}
	}
}

// Option function to enable concurrency feature
func WithConcurrent[T any]() StackOptionFunc[T] {
	return func(s *Stack[T]) {
		s.concurrent = true
	}
}

// Option function for New to set initial capacity to
// something other than the default 16 elements.
func WithCapacity[T any](capacity int) StackOptionFunc[T] {
	if capacity < 0 {
		panic(messages.NEGATIVE_CAPACITY)
	}
	return func(s *Stack[T]) {
		s.initialCapacity = capacity
	}
}

// Option function for New to provide a comparer function for values of type T.
// Required if the element type is not numeric, bool, pointer or string.
func WithComparer[T any](comparer functions.ComparerFunc[T]) StackOptionFunc[T] {
	if comparer == nil {
		panic(messages.COMP_FN_NIL)
	}
	return func(s *Stack[T]) {
		s.compare = comparer
	}
}

// Option func to provide a deep copy implementation for collection elements.
func WithDeepCopy[T any](copier functions.DeepCopyFunc[T]) StackOptionFunc[T] {
	// Can be nil
	return func(s *Stack[T]) {
		s.copy = copier
	}
}

// Add is an alias for [stack.Push].
//
// Always returns true.
func (s *Stack[T]) Add(value T) bool {

	s.Push(value)
	return true
}

// AddRange adds a slice of values to the set,
// effectively pushing the elements from first to last in the slice.
func (s *Stack[T]) AddRange(values []T) {

	if len(values) == 0 {
		return
	}

	lv := len(values)
	if lv == 1 {
		s.Push(values[0])
		return
	}

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}

	newSize := s.size + lv
	newCapacity := util.Iif(newSize > s.initialCapacity, newSize, s.initialCapacity)
	newBuffer := make([]T, newCapacity)
	index := s.size

	copy(newBuffer, s.buffer)

	for i := 0; i < lv; i++ {
		newBuffer[index] = values[i]
		index++
	}

	s.size += lv
	s.version++
	s.buffer = newBuffer
}

// AddCollection pushes the values of the given collection onto this stack.
// Values are pushed in the order defined by the other collection.
func (s *Stack[T]) AddCollection(collection collections.Collection[T]) {

	s.AddRange(collection.ToSliceDeep())
}

// Contains returns true if the stack contains the given value
//
// Stack is searched from most recently pushed value downwards.
func (s *Stack[T]) Contains(value T) bool {

	return s.size != 0 && util.LastIndexOf(s.buffer, value, s.compare, s.concurrent) != -1
}

// Count returns the number of values on the stack.
func (s *Stack[T]) Count() int {

	return s.size
}

// IsEmpty returns true if the collection has no elements
func (s *Stack[T]) IsEmpty() bool {
	return s.size == 0
}

// ToSlice returns a copy of the stack content as a slice.
// The slice is ordered from top to bottom of the stack
// (most recently pushed value is first in the slice).
func (s *Stack[T]) ToSlice() []T {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}
	return s.toSlice(false)
}

// ToSliceDeep returns a copy of the stack content as a slice.
// The slice is ordered from top to bottom of the stack
// (most recently pushed value is first in the slice).
//
// Elements are deep-copied using the provided DeepCopyFunc if any.
func (s *Stack[T]) ToSliceDeep() []T {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}
	return s.toSlice(true)
}

func (s *Stack[T]) toSlice(deepCopy bool) []T {
	slc := make([]T, s.size)

	if deepCopy {
		util.DeepCopySlice(slc, s.buffer[:s.size], s.copy)
	} else {
		copy(slc, s.buffer[:s.size])
	}

	return util.Reverse(slc)
}

// Clear removes all values from the stack.
func (s *Stack[T]) Clear() {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}
	s.buffer = make([]T, 0, cap(s.buffer))
	s.size = 0
	s.version++
}

// Peek returns the value at the top of the stack without adjusting the stack.
//
// Panics if the stack is empty
func (s *Stack[T]) Peek() T {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}
	if s.size == 0 {
		panic(messages.COLLECTION_EMPTY)
	}

	return s.buffer[s.size-1]
}

// TryPeek returns the value at the top of the stack and true if
// the stack is not empty; else zero value of T and false.
func (s *Stack[T]) TryPeek() (T, bool) {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}

	if s.size == 0 {
		var empty T
		return empty, false
	}

	return s.buffer[s.size-1], true
}

// Push adds a value to the top of the stack.
func (s *Stack[T]) Push(value T) {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}
	s.push(value)
}

// Pop removes and returns the value at the top of the stack.
//
// Panics if the stack is empty.
func (s *Stack[T]) Pop() T {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}
	return s.pop()
}

// TryPop pops and returns the value at the top of the stack and true if
// the stack is not empty; else zero value of T and false.
func (s *Stack[T]) TryPop() (T, bool) {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}
	if s.size == 0 {
		var empty T
		return empty, false
	}

	return s.pop(), true
}

// TrimExcess resizes the backing store's length and capacity
// to match the number of elements in the stack.
func (s *Stack[T]) TrimExcess() {

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}
	slc := make([]T, s.size)
	copy(slc, s.buffer[:s.size])
	s.buffer = slc
}

// Remove removes the first occurrence of value found, searching
// from the most recently pushed value.
//
// Returns true if the value was removed.
func (s *Stack[T]) Remove(value T) bool {

	if s.size == 0 {
		return false
	}

	if s.lock != nil {
		s.lock.Lock()
		defer s.lock.Unlock()
	}
	index := util.LastIndexOf(s.buffer, value, s.compare, s.concurrent)

	if index == -1 {
		return false
	}

	var empty T
	s.buffer[index] = empty

	buf := make([]T, len(s.buffer)-1)
	util.PartialCopy(s.buffer, 0, buf, 0, index)
	util.PartialCopy(s.buffer, index+1, buf, index, len(s.buffer)-(index+1))
	s.buffer = buf
	s.version++
	s.size--
	return true
}

// String returns a string representation of container
func (s *Stack[T]) String() string {

	if s.lock != nil {
		s.lock.RLock()
		defer s.lock.RUnlock()
	}
	var values []string
	for _, value := range s.toSlice(false) {
		values = append(values, fmt.Sprintf("%v", value))
	}

	return "Queue\n" + strings.Join(values, ", ")
}

// Type returns the type of this collection
func (*Stack[T]) Type() collections.CollectionType {
	return collections.COLLECTION_STACK
}

func (s *Stack[T]) grow(numElems int) {
	newSize := s.size + numElems
	newBufferSize := util.Iif(newSize > util.DefaultCapacity, newSize*growFactor/100, util.DefaultCapacity)
	buf := make([]T, newBufferSize)
	copy(buf, s.buffer)
	s.buffer = buf
}

func (s *Stack[T]) push(value T) {

	if s.size >= len(s.buffer) {
		s.grow(1)
	}

	s.buffer[s.size] = value
	s.version++
	s.size++
}

func (s *Stack[T]) pop() T {

	if s.size == 0 {
		panic(messages.COLLECTION_EMPTY)
	}

	var empty T
	value := s.buffer[s.size-1]
	s.buffer[s.size-1] = empty
	s.size--
	s.version++

	return value
}

func (s *Stack[T]) length() int {
	return len(s.buffer)
}

func (s *Stack[T]) capacity() int {
	return cap(s.buffer)
}

func (s *Stack[T]) makeDeepCopy() *Stack[T] {

	other := &Stack[T]{
		size:            s.size,
		initialCapacity: s.initialCapacity,
		compare:         s.compare,
		copy:            s.copy,
	}

	if s.lock != nil {
		other.lock = &sync.RWMutex{}
	}

	other.buffer = make([]T, len(s.buffer), cap(s.buffer))
	util.DeepCopySlice(other.buffer, s.buffer, s.copy)
	return other
}
