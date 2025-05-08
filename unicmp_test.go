package unicmp

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func randBytes(n uint) []byte {
	s := make([]byte, n)
	for i := range s {
		s[i] = byte(rand.IntN(int(^byte(0)) + 1))
	}
	return s
}

func sum(s []byte) int {
	sum := 0
	for _, v := range s {
		sum += int(v)
	}
	return sum
}

func shuffle[S ~[]T, T any](s S) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

func TestSmoke(t *testing.T) {
	orig := randBytes(10000)
	s1 := slices.Clone(orig)
	s2 := slices.Clone(orig)
	if a, b := sum(s1), sum(s2); a != b {
		t.Fatalf("sum(s1) != sum(s2): %d != %d", a, b)
	}
	shuffle(s1)
	shuffle(s2)
	if a, b := sum(s1), sum(s2); a != b {
		t.Fatalf("sum(s1) != sum(s2): %d != %d", a, b)
	} else {
		if c := sum(orig); a != c {
			t.Fatalf("sum(s1) != sum(orig): %d != %d", a, c)
		}
	}
	o := ForType[byte]()
	slices.SortFunc(s1, o.Cmp)
	slices.SortFunc(s2, o.Cmp)
	if !slices.Equal(s1, s2) {
		t.Error("s1 content is not the same as s2 after sorting")
	}
	for _, x := range orig {
		if _, found := slices.BinarySearchFunc(s1, x, o.Cmp); !found {
			t.Error("value was not found with bisect in s1")
			break
		}
	}
}

func TestAny(t *testing.T) {
	numbers := randBytes(10000)
	a := make([]any, 0, len(numbers)+3)
	for _, x := range numbers {
		a = append(a, x)
	}
	a = append(a, struct {
		a int
		b string
	}{123, "abcdef"})
	a = append(a, struct {
		a int
		b string
	}{123, "abcdef"})
	a = append(a, struct {
		a int
		b string
	}{100, "!!!!!!!!!!!!"})
	shuffle(a)
	o := ForType[any]()
	slices.SortFunc(a, o.Cmp)
	idx := slices.Index[[]any, any](a, struct{
		a int
		b string
	}{123, "abcdef"})
	if idx == -1 {
		t.Errorf("struct not found in slice")
		return
	}
	if a[idx+1] != struct{
		a int
		b string
	}{123, "abcdef"} {
		t.Errorf("second same struct not found in slice")
	}
}
