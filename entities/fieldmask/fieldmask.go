package fieldmask

import (
	"reflect"
)

type BitmaskType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func EnumerateFields[T BitmaskType](bitmask T, fields *[]T) {
	bits := int(reflect.TypeOf(bitmask).Size())
	for i := 0; i < bits; i++ {
		p := bitmask & T(1) << i
		if p != 0 {
			*fields = append(*fields, p)
		}
	}
}
