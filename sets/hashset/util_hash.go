package hashset

// https://github.com/cornelk/hashmap/blob/36b3b9c2b7ec993f1ef12a6957d45826aca726e6/util_hash.go#L49

import (
	"encoding/binary"
	"fmt"
	"math/bits"
	"reflect"
	"time"
	"unsafe"
)

const (
	prime1 uint64 = 11400714785074694791
	prime2 uint64 = 14029467366897019727
	prime3 uint64 = 1609587929392839161
	prime4 uint64 = 9650029242287828579
	prime5 uint64 = 2870177450012600261
)

const (
	fnvOffset uint64 = 0xcbf29ce484222325
	fnvPrime  uint64 = 0x00000100000001b3
)

var prime1v = prime1

const intSizeBytes = bits.UintSize / 8

/*
Original algorithms and setDefaultHasher function Copyright (c) 2016 Caleb Spare under MIT license.
*/

// setDefaultHasher sets the default hasher depending on the key type.
// Inlines hashing as anonymous functions for performance improvements, other options like
// returning an anonymous functions from another function turned out to not be as performant.
func (s *HashSet[T]) setDefaultHasher() {
	var key T
	kind := reflect.ValueOf(&key).Elem().Type().Kind()

	switch kind {
	case reflect.Bool:
		s.hasher = *(*func(T) uintptr)(unsafe.Pointer(&HashBoolean))
	case reflect.Int, reflect.Uint, reflect.Uintptr, reflect.Pointer:
		switch intSizeBytes {
		case 2:
			s.hasher = *(*func(T) uintptr)(unsafe.Pointer(&HashWord))
		case 4:
			s.hasher = *(*func(T) uintptr)(unsafe.Pointer(&HashDword))
		case 8:
			s.hasher = *(*func(T) uintptr)(unsafe.Pointer(&HashQword))

		default:
			panic(fmt.Errorf("unsupported integer byte size %d", intSizeBytes))
		}

	case reflect.Int8, reflect.Uint8:
		s.hasher = *(*func(T) uintptr)(unsafe.Pointer(&HashByte))
	case reflect.Int16, reflect.Uint16:
		s.hasher = *(*func(T) uintptr)(unsafe.Pointer(&HashWord))
	case reflect.Int32, reflect.Uint32:
		s.hasher = *(*func(T) uintptr)(unsafe.Pointer(&HashDword))
	case reflect.Int64, reflect.Uint64:
		s.hasher = *(*func(T) uintptr)(unsafe.Pointer(&HashQword))
	case reflect.Float32:
		s.hasher = *(*func(T) uintptr)(unsafe.Pointer(&HashFloat32))
	case reflect.Float64:
		s.hasher = *(*func(T) uintptr)(unsafe.Pointer(&HashFloat64))
	case reflect.String:
		s.hasher = *(*func(T) uintptr)(unsafe.Pointer(&HashString))
	case reflect.Struct:
		_, ok := reflect.ValueOf(key).Interface().(time.Time)
		if ok {
			s.hasher = *(*func(T) uintptr)(unsafe.Pointer(&hashTime))
			return
		}
		fallthrough
	default:
		panic(fmt.Errorf("unsupported key type %T of kind %v", key, kind))
	}
}

// https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function
func hashFnv(data []byte) uintptr {
	hash := fnvOffset

	for _, b := range data {
		hash ^= uint64(b)
		hash *= fnvPrime
	}

	return uintptr(hash)
}

// HashTime is a variable containing a function that hashes a [time.Time] struct with the following signature
//
//	func (time.Time) uintptr
var HashTime = hashTime

var hashTime = func(key time.Time) uintptr {

	return hashFnv(unsafe.Slice((*byte)(unsafe.Pointer(&key)), int(unsafe.Sizeof(key))))
}

// Specialized hash functions, optimized for the bit size of the key where available.

// HashBoolean is a variable containing a function that hashes a boolean value with the following signature
//
//	func (bool) uintptr
var HashBoolean = hashBoolean

var hashBoolean = func(key bool) uintptr {
	// Copy what Java does
	if key {
		return uintptr(1231)
	}

	return uintptr(1237)
}

// HashByte is a variable containing a function that hashes an 8 bit integer value with the following signature
//
//	func (uint8) uintptr
var HashByte = hashByte

var hashByte = func(key uint8) uintptr {
	h := prime5 + 1
	h ^= uint64(key) * prime5
	h = bits.RotateLeft64(h, 11) * prime1

	h ^= h >> 33
	h *= prime2
	h ^= h >> 29
	h *= prime3
	h ^= h >> 32

	return uintptr(h)
}

// HashWord is a variable containing a function that hashes a 16 bit integer value with the following signature
//
//	func (uint16) uintptr
var HashWord = hashWord

var hashWord = func(key uint16) uintptr {
	h := prime5 + 2
	h ^= (uint64(key) & 0xff) * prime5
	h = bits.RotateLeft64(h, 11) * prime1
	h ^= ((uint64(key) >> 8) & 0xff) * prime5
	h = bits.RotateLeft64(h, 11) * prime1

	h ^= h >> 33
	h *= prime2
	h ^= h >> 29
	h *= prime3
	h ^= h >> 32

	return uintptr(h)
}

// HashDword is a variable containing a function that hashes a 32 bit integer value with the following signature
//
//	func (unit32) uintptr
var HashDword = hashDword

var hashDword = func(key uint32) uintptr {
	h := prime5 + 4
	h ^= uint64(key) * prime1
	h = bits.RotateLeft64(h, 23)*prime2 + prime3

	h ^= h >> 33
	h *= prime2
	h ^= h >> 29
	h *= prime3
	h ^= h >> 32

	return uintptr(h)
}

// HashFloat32 is a variable containing a function that hashes a 32 bit float value with the following signature
//
//	func (float32) uintptr
var HashFloat32 = hashFloat32

var hashFloat32 = func(key float32) uintptr {
	h := prime5 + 4
	h ^= uint64(key) * prime1
	h = bits.RotateLeft64(h, 23)*prime2 + prime3

	h ^= h >> 33
	h *= prime2
	h ^= h >> 29
	h *= prime3
	h ^= h >> 32

	return uintptr(h)
}

// HashFloat64 is a variable containing a function that hashes a 64 bit float value with the following signature
//
//	func (float64) uintptr
var HashFloat64 = hashFloat64

var hashFloat64 = func(key float64) uintptr {
	h := prime5 + 4
	h ^= uint64(key) * prime1
	h = bits.RotateLeft64(h, 23)*prime2 + prime3

	h ^= h >> 33
	h *= prime2
	h ^= h >> 29
	h *= prime3
	h ^= h >> 32

	return uintptr(h)
}

// HashDword is a variable containing a function that hashes a 64 bit integer value with the following signature
//
//	func (uint64) uintptr
var HashQword = hashQword

var hashQword = func(key uint64) uintptr {
	k1 := key * prime2
	k1 = bits.RotateLeft64(k1, 31)
	k1 *= prime1
	h := (prime5 + 8) ^ k1
	h = bits.RotateLeft64(h, 27)*prime1 + prime4

	h ^= h >> 33
	h *= prime2
	h ^= h >> 29
	h *= prime3
	h ^= h >> 32

	return uintptr(h)
}

// HashString is a variable containing a function that hashes a string value with the following signature
//
//	func (string) uintptr
var HashString = hashString

var hashString = func(key string) uintptr {
	slen := len(key)
	b := unsafe.Slice(unsafe.StringData(key), slen)
	var h uint64

	if slen >= 32 {
		v1 := prime1v + prime2
		v2 := prime2
		v3 := uint64(0)
		v4 := -prime1v
		for len(b) >= 32 {
			v1 = round(v1, binary.LittleEndian.Uint64(b[0:8:len(b)]))
			v2 = round(v2, binary.LittleEndian.Uint64(b[8:16:len(b)]))
			v3 = round(v3, binary.LittleEndian.Uint64(b[16:24:len(b)]))
			v4 = round(v4, binary.LittleEndian.Uint64(b[24:32:len(b)]))
			b = b[32:len(b):len(b)]
		}
		h = rol1(v1) + rol7(v2) + rol12(v3) + rol18(v4)
		h = mergeRound(h, v1)
		h = mergeRound(h, v2)
		h = mergeRound(h, v3)
		h = mergeRound(h, v4)
	} else {
		h = prime5
	}

	h += uint64(slen)

	i, end := 0, len(b)
	for ; i+8 <= end; i += 8 {
		k1 := round(0, binary.LittleEndian.Uint64(b[i:i+8:len(b)]))
		h ^= k1
		h = rol27(h)*prime1 + prime4
	}
	if i+4 <= end {
		h ^= uint64(binary.LittleEndian.Uint32(b[i:i+4:len(b)])) * prime1
		h = rol23(h)*prime2 + prime3
		i += 4
	}
	for ; i < end; i++ {
		h ^= uint64(b[i]) * prime5
		h = rol11(h) * prime1
	}

	h ^= h >> 33
	h *= prime2
	h ^= h >> 29
	h *= prime3
	h ^= h >> 32

	return uintptr(h)
}

func round(acc, input uint64) uint64 {
	acc += input * prime2
	acc = rol31(acc)
	acc *= prime1
	return acc
}

func mergeRound(acc, val uint64) uint64 {
	val = round(0, val)
	acc ^= val
	acc = acc*prime1 + prime4
	return acc
}

func rol1(x uint64) uint64  { return bits.RotateLeft64(x, 1) }
func rol7(x uint64) uint64  { return bits.RotateLeft64(x, 7) }
func rol11(x uint64) uint64 { return bits.RotateLeft64(x, 11) }
func rol12(x uint64) uint64 { return bits.RotateLeft64(x, 12) }
func rol18(x uint64) uint64 { return bits.RotateLeft64(x, 18) }
func rol23(x uint64) uint64 { return bits.RotateLeft64(x, 23) }
func rol27(x uint64) uint64 { return bits.RotateLeft64(x, 27) }
func rol31(x uint64) uint64 { return bits.RotateLeft64(x, 31) }
