package utils

import (
	"math/rand"
)

// Array extends `[]interface{}`, and it should be used to represent json-like data.
type Array []interface{}

func ShuffleSlice[T any](slice []T) {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

func Filter[T any](data []T, f func(T) bool) []T {
	fltd := make([]T, 0, len(data))

	for _, e := range data {
		if f(e) {
			fltd = append(fltd, e)
		}
	}

	return fltd
}

func (l *Array) do(do func()) {
	if len(*l) < 1 {
		*l = make(Array, 0)
	}

	if do != nil {
		do()
	}
}

func (l *Array) Add(v ...interface{}) {
	l.do(func() {
		for _, v := range v {
			*l = append(*l, v)
		}
	})
}

func (l *Array) Delete(k int) {
	l.do(func() {
		if k < len(*l) {
			copy((*l)[k:], (*l)[k+1:])
			(*l)[len(*l)-1] = ""
			*l = (*l)[:len(*l)-1]
		}
	})
}

func (l *Array) Map(fn func(k int, v interface{})) {
	l.do(func() {
		if fn != nil {
			for k, v := range *l {
				fn(k, v)
			}
		}
	})
}
