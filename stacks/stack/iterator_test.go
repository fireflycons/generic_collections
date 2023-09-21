package stack

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
	"github.com/stretchr/testify/require"
)

func TestStackForwardIterator(t *testing.T) {

	var headItems, headItemsReverse []int
	var stack *Stack[int]
	arraySize := 16
	seed := int64(8293)
	headItems, _, headItemsReverse, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Iterates all values", func(t *testing.T) {
		stack = New[int]()

		for i := 0; i < len(headItems); i++ {
			stack.Push(headItems[i])
		}

		index := 0
		iter := stack.Iterator()

		for e := iter.Start(); e != nil; e = iter.Next() {
			require.Equal(t, e.Value(), headItemsReverse[index])
			index++
		}
	})

	t.Run("Calling start mid-iteration restarts iteration", func(t *testing.T) {
		stack = New[int]()

		for i := 0; i < len(headItems); i++ {
			stack.Push(headItems[i])
		}

		iter := stack.Iterator()

		e := iter.Start()
		require.Equal(t, e.Value(), headItemsReverse[0])
		e = iter.Next()
		require.Equal(t, e.Value(), headItemsReverse[1])
		e = iter.Start()
		require.Equal(t, e.Value(), headItemsReverse[0])
	})

	t.Run("Start returns nil if stack is empty -forward iterator", func(t *testing.T) {
		stack = New[int]()
		iter := stack.Iterator()
		require.Nil(t, iter.Start())
	})

	t.Run("Start returns nil if stack is empty - reverse iterator", func(t *testing.T) {
		stack = New[int]()
		iter := stack.ReverseIterator()
		require.Nil(t, iter.Start())
	})

	t.Run("Using ValuePtr on an element does not panic", func(t *testing.T) {
		stack = New[int]()
		stack.Push(util.CreateRandInt(&seed))
		iter := stack.Iterator()

		e := iter.Start()
		require.NotPanics(t, func() { e.ValuePtr() })
	})
}

func TestTakeWhile(t *testing.T) {
	var stackItems, iteratedItems []int
	seed := int64(2163)
	stackItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Returns even numbers", func(t *testing.T) {
		stack := New[int]()

		stack.AddRange(stackItems)

		iter := stack.TakeWhile(func(val int) bool { return val%2 == 0 })

		iteratedItems = make([]int, 0, util.DefaultCapacity)

		for e := iter.Start(); e != nil; e = iter.Next() {
			iteratedItems = append(iteratedItems, e.Value())
		}

		tempItems := make([]int, 0, util.DefaultCapacity)

		for _, v := range stackItems {
			if v%2 == 0 {
				tempItems = append(tempItems, v)
			}
		}

		require.ElementsMatch(t, tempItems, iteratedItems)

	})
}

func TestWhere(t *testing.T) {
	var stackItems []int
	seed := int64(2163)
	stackItems, _, _, _ = util.CreateIntListData(util.DefaultCapacity, &seed)

	t.Run("Returns even numbers", func(t *testing.T) {
		stack := New[int]()
		stack.AddRange(stackItems)
		stack1 := stack.Select(func(val int) bool { return val%2 == 0 })

		tempItems := make([]int, 0, util.DefaultCapacity)

		for _, v := range stackItems {
			if v%2 == 0 {
				tempItems = append(tempItems, v)
			}
		}

		require.ElementsMatch(t, tempItems, stack1.ToSlice())

	})
}

func TestStackForwardIterator_Negative(t *testing.T) {

	var stack *Stack[int]
	seed := int64(8293)

	t.Run("Modifying collection during iteration invalidates iterator", func(t *testing.T) {

		stack = New[int]()

		for i := 0; i < 3; i++ {
			stack.Push(util.CreateRandInt(&seed))
		}

		iter := stack.Iterator()
		require.NotPanics(t, func() { iter.Start() })
		stack.Push(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Next() })
	})

	t.Run("Modifying collection before iteration invalidates iterator", func(t *testing.T) {

		stack = New[int]()

		for i := 0; i < 3; i++ {
			stack.Push(util.CreateRandInt(&seed))
		}

		iter := stack.Iterator()
		stack.Push(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Start() })
	})

	t.Run("Modifying collection after taking an element invalidates element (Value)", func(t *testing.T) {
		stack = New[int]()

		for i := 0; i < 3; i++ {
			stack.Push(util.CreateRandInt(&seed))
		}

		iter := stack.Iterator()
		element := iter.Start()
		stack.Push(util.CreateRandInt(&seed))
		require.Panics(t, func() { element.Value() })
	})

	t.Run("Modifying collection after taking an element invalidates element (ValuePtr)", func(t *testing.T) {
		stack = New[int]()

		for i := 0; i < 3; i++ {
			stack.Push(util.CreateRandInt(&seed))
		}

		iter := stack.Iterator()
		element := iter.Start()
		stack.Push(util.CreateRandInt(&seed))
		require.Panics(t, func() { element.ValuePtr() })
	})

}

func TestStackReverseIterator(t *testing.T) {

	var headItems []int
	var stack *Stack[int]
	arraySize := 16
	seed := int64(8293)
	headItems, _, _, _ = util.CreateIntListData(arraySize, &seed)

	t.Run("Iterates all values", func(t *testing.T) {
		stack = New[int]()

		for i := 0; i < len(headItems); i++ {
			stack.Push(headItems[i])
		}

		index := 0
		iter := stack.ReverseIterator()

		for e := iter.Start(); e != nil; e = iter.Next() {
			require.Equal(t, e.Value(), headItems[index])
			index++
		}
	})
}

func TestStackReverseIterator_Negative(t *testing.T) {

	var stack *Stack[int]
	seed := int64(8293)

	t.Run("Modifying collection during iteration invalidates iterator", func(t *testing.T) {

		stack = New[int]()

		for i := 0; i < 3; i++ {
			stack.Push(util.CreateRandInt(&seed))
		}

		iter := stack.ReverseIterator()
		require.NotPanics(t, func() { iter.Start() })
		stack.Push(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Next() })
	})

	t.Run("Modifying collection before iteration invalidates iterator", func(t *testing.T) {

		stack = New[int]()

		for i := 0; i < 3; i++ {
			stack.Push(util.CreateRandInt(&seed))
		}

		iter := stack.ReverseIterator()
		stack.Push(util.CreateRandInt(&seed))
		require.Panics(t, func() { iter.Start() })
	})

	t.Run("Modifying collection after taking an element invalidates element (Value)", func(t *testing.T) {
		stack = New[int]()

		for i := 0; i < 3; i++ {
			stack.Push(util.CreateRandInt(&seed))
		}

		iter := stack.ReverseIterator()
		element := iter.Start()
		stack.Push(util.CreateRandInt(&seed))
		require.Panics(t, func() { element.Value() })
	})

	t.Run("Modifying collection after taking an element invalidates element (ValuePtr)", func(t *testing.T) {
		stack = New[int]()

		for i := 0; i < 3; i++ {
			stack.Push(util.CreateRandInt(&seed))
		}

		iter := stack.ReverseIterator()
		element := iter.Start()
		stack.Push(util.CreateRandInt(&seed))
		require.Panics(t, func() { element.ValuePtr() })
	})

}
