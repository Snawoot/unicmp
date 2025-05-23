// Package unicmp provides generic ordering function for all comparable types.
//
// The ordering is made by compares of one or more output of maphash function with
// different seeds. The resulting ordering is transitive and anticommutative, suitable
// for sorting and implementation of algorithms on top of sorted collections.
package unicmp

import (
	"cmp"

	"hash/maphash"
)

const maxRounds = 64

var seeds [maxRounds]maphash.Seed

func init() {
	for i := range seeds {
		seeds[i] = maphash.MakeSeed()
	}
}

func needsReflection[T comparable](x T) bool {
	return x != x
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
	var i int
	fxEven, fyEven := maphash.Comparable[T], maphash.Comparable[T]
	fxOdd, fyOdd := extendedMapHash[T], extendedMapHash[T]
	if needsReflection(x) {
		fxEven = extendedMapHash[T]
	}
	if needsReflection(y) {
		fyEven = extendedMapHash[T]
	}
	for ; h1 == h2 && i < maxRounds; i++ {
		if i%2 == 0 {
			h1, h2 = fxEven(seeds[i], x), fyEven(seeds[i], y)
		} else {
			h1, h2 = fxOdd(seeds[i], x), fyOdd(seeds[i], y)
		}
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
