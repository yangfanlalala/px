package slice

import (
	"reflect"
)

func Reverse(slice interface{}) {
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		panic("not slice")
	}
	size := val.Len()
	if size == 0 || size == 1 {
		return
	}
	mid := size / 2
	swapper := reflect.Swapper(slice)
	for i := 0; i < mid; i++ {
		swapper(i, size-1-i)
	}
}
