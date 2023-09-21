package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComparerBool(t *testing.T) {
	f := GetDefaultComparer[bool]()

	require.Equal(t, 0, f(true, true))
	require.Greater(t, f(true, false), 0)
	require.Less(t, f(false, true), 0)
}

func TestComparerUintptr(t *testing.T) {
	f := GetDefaultComparer[uintptr]()

	require.Equal(t, 0, f(uintptr(0x9f29cb17a2a49995), uintptr(0x9f29cb17a2a49995)))
	require.Greater(t, f(uintptr(0x9f29cb17a2a49995), uintptr(0x9f29cb17a2a40000)), 0)
}

func TestComparerUint64(t *testing.T) {
	f := GetDefaultComparer[uint64]()

	require.Equal(t, 0, f(0, 0))
	require.Greater(t, f(1, 0), 0)
}

func TestComparerInt64(t *testing.T) {
	f := GetDefaultComparer[int64]()

	require.Equal(t, 0, f(-10, -10))
	require.Less(t, f(-1, 0), 0)
}

func TestComparerUint32(t *testing.T) {
	f := GetDefaultComparer[uint32]()

	require.Equal(t, 0, f(0, 0))
	require.Greater(t, f(1, 0), 0)
}

func TestComparerInt32(t *testing.T) {
	f := GetDefaultComparer[int32]()

	require.Equal(t, 0, f(-10, -10))
	require.Less(t, f(-1, 0), 0)
}

func TestComparerUint16(t *testing.T) {
	f := GetDefaultComparer[uint16]()

	require.Equal(t, 0, f(0, 0))
	require.Greater(t, f(1, 0), 0)
}

func TestComparerInt16(t *testing.T) {
	f := GetDefaultComparer[int16]()

	require.Equal(t, 0, f(-10, -10))
	require.Less(t, f(-1, 0), 0)
}

func TestComparerUint8(t *testing.T) {
	f := GetDefaultComparer[uint8]()

	require.Equal(t, 0, f(0, 0))
	require.Greater(t, f(1, 0), 0)
}

func TestComparerInt8(t *testing.T) {
	f := GetDefaultComparer[int8]()

	require.Equal(t, 0, f(-10, -10))
	require.Less(t, f(-1, 0), 0)
}

func TestComparerInt(t *testing.T) {
	f := GetDefaultComparer[int]()

	require.Equal(t, 0, f(-5, -5))
	require.Less(t, f(-5, 5), 0)
}

func TestCompareString(t *testing.T) {
	f := GetDefaultComparer[string]()
	require.Equal(t, 0, f("properunittesting", "properunittesting"))
	require.Greater(t, f("properunittesting", "longstringlongstringlongstringlongstring"), 0)
}

func TestCompareFloat32(t *testing.T) {
	f := GetDefaultComparer[float32]()
	require.Equal(t, 0, f(1.234e+9, 1.234e+9))
	require.Greater(t, f(1.234e+9, 1.234e-9), 0)
}

func TestCompareFloat64(t *testing.T) {
	f := GetDefaultComparer[float64]()
	require.Equal(t, 0, f(1.234e+9, 1.234e+9))
	require.Greater(t, f(1.234e+9, 1.234e-9), 0)
}

func TestGetDefaultComparerPanicsWithUnsupportedType(t *testing.T) {
	type testType struct{}

	require.Panics(t, func() { GetDefaultComparer[testType]() })
}

func TestComparePointer(t *testing.T) {
	var a, b int

	a = 1
	b = 1
	aPtr := &a
	a2Ptr := &a
	bPtr := &b

	f := GetDefaultComparer[*int]()
	require.Equal(t, 0, f(aPtr, a2Ptr))
	require.NotEqual(t, 0, f(aPtr, bPtr))
}

func TestComparePointerToStruct(t *testing.T) {

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

	f := GetDefaultComparer[*someStruct]()
	require.Equal(t, 0, f(aPtr, a2Ptr))
	require.NotEqual(t, 0, f(aPtr, bPtr))
}
