package dlist

import (
	"testing"

	"github.com/fireflycons/generic_collections/internal/util"
)

func TestNode(t *testing.T) {

	var (
		node  *DListNode[int]
		value int
	)

	seed := int64(21543)

	t.Run("Verify passing default(T) into the constructor", func(t *testing.T) {
		node = NewNode(defaultT[int]())
		verifyDListNodeA(t, node, defaultT[int](), nil, nil, nil)
	})

	t.Run("Verify passing something other then default(T) into the constructor", func(t *testing.T) {
		value = util.CreateRandInt(&seed)
		node = NewNode(value)
		verifyDListNodeA(t, node, value, nil, nil, nil)
	})

	t.Run("Verify passing something other then default(T) into the constructor and set the value to something other then default(T)", func(t *testing.T) {
		value = util.CreateRandInt(&seed)
		node = NewNode(value)
		value = util.CreateRandInt(&seed)
		node.item = value
		verifyDListNodeA(t, node, value, nil, nil, nil)
	})

	t.Run("Verify passing default(T) into the constructor and set the value to default(T)", func(t *testing.T) {
		node = NewNode(defaultT[int]())
		value = util.CreateRandInt(&seed)
		node.item = defaultT[int]()
		verifyDListNodeA(t, node, defaultT[int](), nil, nil, nil)
	})

	t.Run("Verify passing something other then default(T) into the constructor and set the value to something other then default(T)", func(t *testing.T) {
		value = util.CreateRandInt(&seed)
		node = NewNode(value)
		value = util.CreateRandInt(&seed)
		ref := node.ValuePtr()
		*ref = value
		verifyDListNodeA(t, node, value, nil, nil, nil)
	})

	t.Run("Verify passing something other then default(T) into the constructor and set the value to default(T)", func(t *testing.T) {
		value = util.CreateRandInt(&seed)
		node = NewNode(value)
		ref := node.ValuePtr()
		*ref = defaultT[int]()
		verifyDListNodeA(t, node, defaultT[int](), nil, nil, nil)
	})

	t.Run("Verify passing default(T) into the constructor and set the value to something other then default(T)", func(t *testing.T) {
		node = NewNode(defaultT[int]())
		value = util.CreateRandInt(&seed)
		ref := node.ValuePtr()
		*ref = value
		verifyDListNodeA(t, node, value, nil, nil, nil)
	})

	t.Run("Verify passing default(T) into the constructor and set the value to default(T)", func(t *testing.T) {
		node = NewNode(defaultT[int]())
		value = util.CreateRandInt(&seed)
		ref := node.ValuePtr()
		*ref = defaultT[int]()
		verifyDListNodeA(t, node, defaultT[int](), nil, nil, nil)
	})
}
