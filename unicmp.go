package unicmp

import (
	"cmp"

	"github.com/Snawoot/maphash"
)

const maxRounds = 50

type Ordering[T comparable] struct {
	h maphash.Hasher[T]
}

func ForType[T comparable]() Ordering[T] {
	return Ordering[T]{maphash.NewHasher[T](maphash.NewSeed(0))}
}

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

func (o Ordering[T]) Less(x, y T) bool {
	return o.Cmp(x, y) < 0
}

func (o Ordering[T]) Equal(x, y T) bool {
	return o.Cmp(x, y) == 0
}
