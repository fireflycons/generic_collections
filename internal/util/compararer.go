package util

import (
	"fmt"
	"math/bits"
	"reflect"
	"time"
	"unsafe"

	"github.com/fireflycons/generic_collections/functions"
	"github.com/fireflycons/generic_collections/internal/messages"
	"golang.org/x/exp/constraints"
)

const intSizeBytes = bits.UintSize / 8

type ordered interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~string
}

// GetDefaultComparer returns a function to compare two values of
// types supported by this module.
//
// Return value semantics are that if first value less than second value, the result
// is a negative integer. If they are equal, the result is zero, else a positive integer.
// Magniitude of the result is unimportant. This permits a faster compare for signed int
// values by use of simple subtraction.
func GetDefaultComparer[T any]() functions.ComparerFunc[T] {
	var key T

	kind := reflect.ValueOf(&key).Elem().Type().Kind()
	switch kind {
	case reflect.Bool:
		return *(*func(T, T) int)(unsafe.Pointer(&xxCompareBool))
	case reflect.Int:
		switch intSizeBytes {
		case 2:
			return *(*func(T, T) int)(unsafe.Pointer(&xxCompareSignedWord))
		case 4:
			return *(*func(T, T) int)(unsafe.Pointer(&xxCompareSignedDword))
		case 8:
			return *(*func(T, T) int)(unsafe.Pointer(&xxCompareSignedQword))

		default:
			panic(fmt.Errorf(messages.COMPARER_INVALID_INT_FMT, intSizeBytes))
		}
	case reflect.Uint, reflect.Uintptr, reflect.Pointer:
		switch intSizeBytes {
		case 2:
			return *(*func(T, T) int)(unsafe.Pointer(&xxCompareUnsignedWord))
		case 4:
			return *(*func(T, T) int)(unsafe.Pointer(&xxCompareUnsignedDword))
		case 8:
			return *(*func(T, T) int)(unsafe.Pointer(&xxCompareUnsignedQword))

		default:
			panic(fmt.Errorf(messages.COMPARER_INVALID_INT_FMT, intSizeBytes))
		}

	case reflect.Int8:
		return *(*func(T, T) int)(unsafe.Pointer(&xxCompareSignedByte))
	case reflect.Uint8:
		return *(*func(T, T) int)(unsafe.Pointer(&xxCompareUnsignedByte))
	case reflect.Int16:
		return *(*func(T, T) int)(unsafe.Pointer(&xxCompareSignedWord))
	case reflect.Uint16:
		return *(*func(T, T) int)(unsafe.Pointer(&xxCompareUnsignedWord))
	case reflect.Int32:
		return *(*func(T, T) int)(unsafe.Pointer(&xxCompareSignedDword))
	case reflect.Uint32:
		return *(*func(T, T) int)(unsafe.Pointer(&xxCompareUnsignedDword))
	case reflect.Int64:
		return *(*func(T, T) int)(unsafe.Pointer(&xxCompareSignedQword))
	case reflect.Uint64:
		return *(*func(T, T) int)(unsafe.Pointer(&xxCompareUnsignedQword))
	case reflect.Float32:
		return *(*func(T, T) int)(unsafe.Pointer(&xxCompareFloat32))
	case reflect.Float64:
		return *(*func(T, T) int)(unsafe.Pointer(&xxCompareFloat64))
	case reflect.String:
		return *(*func(T, T) int)(unsafe.Pointer(&xxCompareString))
	case reflect.Struct:
		_, ok := reflect.ValueOf(key).Interface().(time.Time)
		if ok {
			return *(*func(T, T) int)(unsafe.Pointer(&xxCompareTime))
		}
		fallthrough

	default:
		panic(fmt.Sprintf(messages.COMPARER_INVALID_KEY_FMT, key, kind))
	}
}

var xxCompareTime = func(t1, t2 time.Time) int {
	d := t1.Sub(t2)
	if d < 0 {
		return -1
	}

	if d == 0 {
		return 0
	}

	return 1
}

var xxCompareBool = func(v1, v2 bool) int {
	// Consider true > false as these are normally represented as 1 and 0
	if v1 == v2 {
		return 0
	}

	if v1 && !v2 {
		return 1
	}

	return -1
}

var xxCompareUnsignedByte = func(v1, v2 uint8) int {
	return xxCompare(v1, v2)
}

var xxCompareUnsignedWord = func(v1, v2 uint16) int {
	return xxCompare(v1, v2)
}

var xxCompareUnsignedDword = func(v1, v2 uint32) int {
	return xxCompare(v1, v2)
}

var xxCompareUnsignedQword = func(v1, v2 uint64) int {
	return xxCompare(v1, v2)
}

var xxCompareSignedByte = func(v1, v2 int8) int {
	return xxCompareSignedInt(v1, v2)
}

var xxCompareSignedWord = func(v1, v2 int16) int {
	return xxCompareSignedInt(v1, v2)
}

var xxCompareSignedDword = func(v1, v2 int32) int {
	return xxCompareSignedInt(v1, v2)
}

var xxCompareSignedQword = func(v1, v2 int64) int {
	return xxCompareSignedInt(v1, v2)
}

var xxCompareFloat32 = func(v1, v2 float32) int {
	return xxCompareFloat(v1, v2)
}

var xxCompareFloat64 = func(v1, v2 float64) int {
	return xxCompareFloat(v1, v2)
}

var xxCompareString = func(v1, v2 string) int {
	return xxCompare(v1, v2)
}

// Compare for unsigned and string types that
// can't be done with a simple subtraction
func xxCompare[T ordered](v1, v2 T) int {
	if v1 == v2 {
		return 0
	}

	if v1 < v2 {
		return -1
	}

	return 1
}

func xxCompareSignedInt[T constraints.Signed](v1, v2 T) int {
	return int(v1 - v2)
}

// Compare for float types that handle NaN as per the
// Go 1.21 cmp package
func xxCompareFloat[T constraints.Float](v1, v2 T) int {
	v1Nan := isNaN(v1)
	v2Nan := isNaN(v2)
	if v1Nan && v2Nan {
		return 0
	}
	if v1Nan || v1 < v2 {
		return -1
	}
	if v2Nan || v1 > v2 {
		return 1
	}
	return 0
}

// isNaN reports whether x is a NaN without requiring the math package.
func isNaN[T constraints.Float](x T) bool {
	return x != x
}
