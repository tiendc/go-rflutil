package rflutil

import (
	"fmt"
	"reflect"
)

// SliceGet get value from a slice type at the given index
func SliceGet[T any](s reflect.Value, i int) (T, error) {
	var ret T
	slice := indirectValueTilRoot(s)
	if !slice.IsValid() || !isKindIn(slice.Kind(), reflect.Slice, reflect.Array) {
		return ret, fmt.Errorf("%w: require slice or array type (got %v)", ErrTypeInvalid, s.Type())
	}

	if i < 0 || i >= slice.Len() {
		return ret, fmt.Errorf("%w: index %d is out of range", ErrIndexOutOfRange, i)
	}
	item := slice.Index(i).Interface()
	if item == nil {
		dstType := reflect.TypeOf(ret)
		if dstType == nil || dstType.Kind() == reflect.Interface {
			return ret, nil
		}
		return ret, fmt.Errorf("%w: item type is %v (got %v)",
			ErrTypeUnmatched, slice.Type().Elem(), dstType)
	}

	ret, ok := item.(T)
	if !ok {
		return ret, fmt.Errorf("%w: item type is %v (got %v)",
			ErrTypeUnmatched, slice.Type().Elem(), reflect.TypeOf(ret))
	}
	return ret, nil
}

// SliceSet set value in a slice type at the given index
func SliceSet[T any](s reflect.Value, i int, v T) error {
	slice := indirectValueTilRoot(s)
	if !slice.IsValid() || !isKindIn(slice.Kind(), reflect.Slice, reflect.Array) {
		return fmt.Errorf("%w: require slice or array type (got %v)", ErrTypeInvalid, s.Type())
	}

	if i < 0 || i >= slice.Len() {
		return fmt.Errorf("%w: index %d is out of range", ErrIndexOutOfRange, i)
	}

	itemType := slice.Type().Elem()
	val := reflect.ValueOf(v)
	if !val.IsValid() {
		if itemType.Kind() == reflect.Interface {
			slice.Index(i).Set(reflect.Zero(itemType))
			return nil
		}
		return fmt.Errorf("%w: item type is %v (expect %v)",
			ErrTypeUnmatched, itemType, reflect.TypeOf([]interface{}{}).Elem())
	}
	if !val.Type().AssignableTo(itemType) {
		return fmt.Errorf("%w: item type is %v (expect %v)", ErrTypeUnmatched, itemType, val.Type())
	}
	slice.Index(i).Set(val)
	return nil
}

// SliceAppend appends the given value to a slice
func SliceAppend[T any](s reflect.Value, v T) ([]T, error) {
	slice := indirectValueTilRoot(s)
	if !slice.IsValid() || slice.Kind() != reflect.Slice {
		return nil, fmt.Errorf("%w: require slice type (got %v)", ErrTypeInvalid, s.Type())
	}

	itemType := slice.Type().Elem()
	val := reflect.ValueOf(v)
	if !val.IsValid() {
		if itemType.Kind() == reflect.Interface {
			return reflect.Append(slice, reflect.Zero(itemType)).Interface().([]T), nil // nolint: forcetypeassert
		}
		return nil, fmt.Errorf("%w: item type is %v (expect %v)",
			ErrTypeUnmatched, itemType, reflect.TypeOf([]interface{}{}).Elem())
	}
	if !val.Type().AssignableTo(itemType) {
		return nil, fmt.Errorf("%w: item type is %v (expect %v)", ErrTypeUnmatched, itemType, val.Type())
	}
	return reflect.Append(slice, val).Interface().([]T), nil // nolint: forcetypeassert
}

// SliceAs convert a slice or array to the expected slice type.
// Type T must be assignable or convertible to the input slice's item type.
func SliceAs[T any](s reflect.Value) ([]T, error) {
	src := indirectValueTilRoot(s)
	if !src.IsValid() || !isKindIn(src.Kind(), reflect.Slice, reflect.Array) {
		return nil, fmt.Errorf("%w: require slice or array type (got %v)", ErrTypeInvalid, s.Type())
	}

	srcType := src.Type()
	dstType := reflect.TypeOf([]T{})
	if dstType == srcType {
		return src.Interface().([]T), nil // nolint: forcetypeassert
	}

	dstItemType := dstType.Elem()
	assignable := dstItemType.AssignableTo(srcType.Elem())
	convertible := !assignable && dstItemType.ConvertibleTo(srcType.Elem())
	if !assignable && !convertible {
		return nil, fmt.Errorf("%w: unable to convert slice to '%v'", ErrTypeInvalid, dstType)
	}

	numItems := src.Len()
	dst := reflect.MakeSlice(dstType, numItems, numItems)
	for i := 0; i < numItems; i++ {
		if assignable {
			dst.Index(i).Set(src.Index(i))
		} else {
			dst.Index(i).Set(src.Index(i).Convert(dstItemType))
		}
	}
	return dst.Interface().([]T), nil // nolint: forcetypeassert
}
