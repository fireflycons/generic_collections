//go:build !386

package hashset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashingBool(t *testing.T) {
	m := New[bool]()
	require.Equal(t, uintptr(1231), m.hasher(true))
	require.Equal(t, uintptr(1237), m.hasher(false))
}

func TestHashingUintptr(t *testing.T) {
	m := New[uintptr]()
	require.Equal(t, uintptr(0x9f29cb17a2a49995), m.hasher(1))
	require.Equal(t, uintptr(0xeac73e4044e82db0), m.hasher(2))
}

func TestHashingUint64(t *testing.T) {
	m := New[uint64]()
	require.Equal(t, uintptr(0x9f29cb17a2a49995), m.hasher(1))
	require.Equal(t, uintptr(0xeac73e4044e82db0), m.hasher(2))
}

func TestHashingUint32(t *testing.T) {
	m := New[uint32]()
	require.Equal(t, uintptr(0xf42f94001fcb5351), m.hasher(1))
	require.Equal(t, uintptr(0x277af360cedcb29e), m.hasher(2))
}

func TestHashingUint16(t *testing.T) {
	m := New[uint16]()
	require.Equal(t, uintptr(0xdd8f621dbf7f57f1), m.hasher(1))
	require.Equal(t, uintptr(0xfc2f33e9edde6f4a), m.hasher(0x102))
}

func TestHashingUint8(t *testing.T) {
	m := New[uint8]()
	require.Equal(t, uintptr(0x8a4127811b21e730), m.hasher(1))
	require.Equal(t, uintptr(0x4b79b8c95732b0e7), m.hasher(2))
}

func TestHashingFloat32(t *testing.T) {
	m := New[float32]()
	require.Equal(t, uintptr(0xe44a4057281797c8), m.hasher(1.234e9))
	require.Equal(t, uintptr(0x3aefa6fd5cf2deb4), m.hasher(-4.543259e-15))
}

func TestHashingFloat64(t *testing.T) {
	m := New[float64]()
	require.Equal(t, uintptr(0xb12a79c149c830ec), m.hasher(1.234e201))
	require.Equal(t, uintptr(0x3aefa6fd5cf2deb4), m.hasher(6.345543e-124))
}

func TestHashingString(t *testing.T) {
	m := New[string]()
	require.Equal(t, uintptr(0xd24ec4f1a98c6e5b), m.hasher("a"))
	require.Equal(t, uintptr(0xac7a9b03a75c7b4e), m.hasher("12345678901234"))
	require.Equal(t, uintptr(0x6a1faf26e7da4cb9), m.hasher("properunittesting"))
	require.Equal(t, uintptr(0x2d4ff7e12135f1f3), m.hasher("longstringlongstringlongstringlongstring"))
}

func TestHashingPointer(t *testing.T) {
	var a, b int

	a = 1
	b = 1
	aPtr := &a
	a2Ptr := &a
	bPtr := &b

	m := New[*int]()
	require.Equal(t, m.hasher(aPtr), m.hasher(a2Ptr))
	require.NotEqual(t, m.hasher(aPtr), m.hasher(bPtr))
}

func TestHashingPointerToStruct(t *testing.T) {
	type someStruct struct {
		a int
		b string
	}

	a := someStruct{
		a: 1,
		b: "x",
	}

	b := someStruct{
		a: 1,
		b: "x",
	}
	aPtr := &a
	a2Ptr := &a
	bPtr := &b

	m := New[*someStruct]()
	require.Equal(t, m.hasher(aPtr), m.hasher(a2Ptr))
	require.NotEqual(t, m.hasher(aPtr), m.hasher(bPtr))
}
