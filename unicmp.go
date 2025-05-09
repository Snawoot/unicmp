// Package unicmp provides generic ordering function for all comparable types.
//
// The ordering is made by compares of one or more output of maphash function with
// different seeds. The resulting ordering is transitive and anticommutative, suitable
// for sorting and implementation of algorithms on top of sorted collections.
package unicmp

import (
	"cmp"
	"unsafe"

	"hash/maphash"
)

const maxRounds = 8

var seeds [maxRounds]maphash.Seed

func init() {
	for i := range seeds {
		seeds[i] = maphash.MakeSeed()
	}
}

// Cmp returns
//
//	-1 if x is less than y,
//	 0 if x equals y,
//	+1 if x is greater than y.
func Cmp[T comparable](x, y T) int {
	if x == y {
		return 0
	}
	var h1, h2 uint64
	for i := 0; h1 == h2 && i < maxRounds; i++ {
		h1, h2 = maphash.Comparable(seeds[i], x), maphash.Comparable(seeds[i], y)
	}
	if h1 == h2 {
		return memcmp(unsafe.Pointer(&x), unsafe.Pointer(&y), unsafe.Sizeof(x))
	}
	return cmp.Compare(h1, h2)
}

// Less returns true if x sorts before y (x < y).
func Less[T comparable](x, y T) bool {
	return Cmp(x, y) < 0
}

// Equal returns true if x == y.
func Equal[T comparable](x, y T) bool {
	return Cmp(x, y) == 0
}

func memcmp(x, y unsafe.Pointer, size uintptr) int {
	for i := uintptr(0); i < size; i++ {
		px := (*byte)(unsafe.Pointer(uintptr(x) + uintptr(i)))
		py := (*byte)(unsafe.Pointer(uintptr(y) + uintptr(i)))
		switch {
		case *px < *py:
			return -1
		case *px > *py:
			return 1
		}
	}
	return 0
}
