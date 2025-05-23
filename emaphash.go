package unicmp

import (
	"encoding/binary"
	"errors"
	"hash/maphash"
	"math"
	"reflect"
)

type emaphash struct {
	maphash.Hash
}

var typeSeed = maphash.MakeSeed()

func (h *emaphash) float64(f float64) {
	if f == 0 {
		h.WriteByte(0)
		return
	}
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(f))
	h.Write(buf[:])
}

func (h *emaphash) uint64(u uint64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], u)
	h.Write(buf[:])
}

// use only in the same scope as maphash.Comparable against
// same variables, otherwise escapeForHash compiler intrinsic
// won't be applied
func extendedMapHash[T comparable](seed maphash.Seed, v T) uint64 {
	var h emaphash
	h.SetSeed(seed)
	writeComparable(&h, v)
	return h.Sum64()
}

func writeComparable[T comparable](h *emaphash, v T) {
	vv := reflect.Indirect(reflect.ValueOf(&v))
	h.appendT(vv)
}

func (h *emaphash) appendT(v reflect.Value) {
	h.uint64(maphash.Comparable(typeSeed, v.Type()))
	switch v.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		h.uint64(uint64(v.Int()))
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uintptr:
		h.uint64(v.Uint())
	case reflect.Array:
		for i := range uint64(v.Len()) {
			h.appendT(v.Index(int(i)))
		}
	case reflect.String:
		h.uint64(uint64(v.Len()))
		h.WriteString(v.String())
	case reflect.Struct:
		for i := range v.NumField() {
			h.appendT(v.Field(i))
		}
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		h.float64(real(c))
		h.float64(imag(c))
	case reflect.Float32, reflect.Float64:
		h.float64(v.Float())
	case reflect.Bool:
		h.WriteByte(btoi(v.Bool()))
	case reflect.UnsafePointer, reflect.Pointer, reflect.Chan:
		h.uint64(uint64(v.Pointer()))
	case reflect.Interface:
		h.appendT(v.Elem())
	default:
		panic(errors.New("maphash: hash of unhashable type " + v.Type().String()))
	}
}

func btoi(b bool) byte {
	if b {
		return 1
	}
	return 0
}
