package unicmp

import (
	"cmp"
	"fmt"
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
	slices.SortFunc(s1, Cmp)
	slices.SortFunc(s2, Cmp)
	if !slices.Equal(s1, s2) {
		t.Error("s1 content is not the same as s2 after sorting")
	}
	for _, x := range orig {
		if _, found := slices.BinarySearchFunc(s1, x, Cmp); !found {
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
	slices.SortFunc(a, Cmp)
	idx := slices.Index[[]any, any](a, struct {
		a int
		b string
	}{123, "abcdef"})
	if idx == -1 {
		t.Errorf("struct not found in slice")
		return
	}
	if a[idx+1] != struct {
		a int
		b string
	}{123, "abcdef"} {
		t.Errorf("second same struct not found in slice")
	}
}

func TestAny2(t *testing.T) {
	var x *int
	var y *byte
	if Cmp[any](x, y) == 0 {
		t.Error("any(*int(nil)) == any(*byte(nil))")
	}
}

var result int

func BenchmarkWorstCase(b *testing.B) {
	var (
		r int
		x *int
		y *float32
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = Cmp[any](x, y)
	}
	result = r
}

func BenchmarkBestCase(b *testing.B) {
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = Cmp(1, 2)
	}
	result = r
}

func BenchmarkSimpleCmp(b *testing.B) {
	var r int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = cmp.Compare(1, 2)
	}
	result = r
}

type Bar struct {
	msg string
}

func (b *Bar) Xf() {
}

type Foo struct {
	X, Y, Z float64
}

func (f *Foo) Xf() {
	fmt.Println("X called")
}

type Xer interface {
	Xf()
}

func TestAny3(t *testing.T) {
	var f1 *Foo
	x1 := Xer(f1)
	var b1 *Bar
	x2 := Xer(b1)
	if Cmp(x1, x2) == 0 {
		t.Error("type-specific nil-values of interface value are equal!")
	}
}
