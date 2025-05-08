// Package unicmp provides generic ordering function for all comparable types.
//
// The ordering is made by compares of one or more output of maphash function with
// different seeds. The resulting ordering is transitive and anticommutative, suitable
// for sorting and implementation of algorithms on top of sorted collections.
package unicmp

import (
	"cmp"

	"github.com/Snawoot/maphash"
)

const maxRounds = 10

// Ordering is an instance of ordering hasher. All Ordering methods are safe for
// concurrent use.
type Ordering[T comparable] struct {
	h maphash.Hasher[T]
}

// ForType returns a new Ordering for specified type.
func ForType[T comparable]() Ordering[T] {
	return Ordering[T]{maphash.NewHasher[T](maphash.NewSeed(0))}
}

// Cmp returns
//
//	-1 if x is less than y,
//	 0 if x equals y,
//	+1 if x is greater than y.
func (o Ordering[T]) Cmp(x, y T) int {
	if x == y {
		return 0
	}
	h := o.h
	h1, h2 := h.Hash2(x), h.Hash2(y)
	for i := uintptr(1); h1 == h2 && i < maxRounds; h1, h2, i = h.Hash2(x), h.Hash2(y), i+1 {
		h = h.WithSeed(maphash.NewSeed(i))
	}
	return cmp.Compare(h1, h2)
}

// Less returns true if x sorts before y (x < y).
func (o Ordering[T]) Less(x, y T) bool {
	return o.Cmp(x, y) < 0
}

// Equal returns true if x == y.
func (o Ordering[T]) Equal(x, y T) bool {
	return o.Cmp(x, y) == 0
}
