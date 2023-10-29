package rflutil

import (
	"fmt"
	"reflect"
)

func ValueAs[T any](v reflect.Value) (T, error) {
	var ret T
	targetType := reflect.TypeOf(ret)
	sourceType := v.Type()

	for {
		if sourceType.AssignableTo(targetType) {
			return v.Interface().(T), nil // nolint: forcetypeassert
		}
		if sourceType.ConvertibleTo(targetType) {
			return v.Convert(targetType).Interface().(T), nil // nolint: forcetypeassert
		}

		if v.IsValid() && v.Kind() == reflect.Interface {
			v = v.Elem()
			if !v.IsValid() {
				break
			}
			sourceType = v.Type()
		} else {
			break
		}
	}

	return ret, fmt.Errorf("%w: value type is %v (expect %v)", ErrTypeUnmatched, sourceType, targetType)
}
