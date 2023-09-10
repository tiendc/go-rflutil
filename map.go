package rflutil

import (
	"fmt"
	"reflect"
)

// MapGet get value from a map by key
func MapGet[V any, K comparable](m reflect.Value, k K) (V, error) {
	var ret V
	val := indirectValueTilRoot(m)
	if !val.IsValid() || val.Kind() != reflect.Map {
		return ret, fmt.Errorf("%w: require map type (got %v)", ErrTypeInvalid, m.Type())
	}

	mapType := val.Type()
	keyVal := reflect.ValueOf(k)
	if !mapType.Key().AssignableTo(keyVal.Type()) {
		return ret, fmt.Errorf("%w: key type is %v (expect %v)", ErrTypeUnmatched,
			mapType.Key(), keyVal.Type())
	}

	valueVal := val.MapIndex(keyVal)
	if !valueVal.IsValid() {
		return ret, ErrNotFound
	}
	ret, ok := valueVal.Interface().(V)
	if !ok {
		return ret, fmt.Errorf("%w: value type is %v (expect %v)", ErrTypeUnmatched,
			mapType.Elem(), reflect.TypeOf(ret))
	}
	return ret, nil
}

// MapSet set value for a key of a map
func MapSet[K comparable, V any](m reflect.Value, k K, v V) error {
	val := indirectValueTilRoot(m)
	if !val.IsValid() || val.Kind() != reflect.Map {
		return fmt.Errorf("%w: require map type (got %v)", ErrTypeInvalid, m.Type())
	}

	mapType := val.Type()
	keyVal := reflect.ValueOf(k)
	if !mapType.Key().AssignableTo(keyVal.Type()) {
		return fmt.Errorf("%w: key type is %v (expect %v)", ErrTypeUnmatched,
			mapType.Key(), keyVal.Type())
	}
	valVal := reflect.ValueOf(v)
	if !valVal.IsValid() {
		if mapType.Elem().Kind() == reflect.Interface {
			val.SetMapIndex(keyVal, reflect.Zero(mapType.Elem()))
			return nil
		}
	}
	if !mapType.Elem().AssignableTo(valVal.Type()) {
		return fmt.Errorf("%w: value type is %v (expect %v)", ErrTypeUnmatched,
			mapType.Elem(), valVal.Type())
	}

	val.SetMapIndex(keyVal, valVal)
	return nil
}

// MapDelete delete the given key from a map
func MapDelete[K comparable](m reflect.Value, k K) error {
	val := indirectValueTilRoot(m)
	if !val.IsValid() || val.Kind() != reflect.Map {
		return fmt.Errorf("%w: require map type (got %v)", ErrTypeInvalid, m.Type())
	}

	mapType := val.Type()
	keyVal := reflect.ValueOf(k)
	if !mapType.Key().AssignableTo(keyVal.Type()) {
		return fmt.Errorf("%w: key type is %v (expect %v)", ErrTypeUnmatched,
			mapType.Key(), keyVal.Type())
	}
	// Set zero value means delete the key from the map
	val.SetMapIndex(keyVal, reflect.Value{})
	return nil
}

// MapAs convert a map to the expected map type.
// Key and Value types must be assignable or convertible to the equivalent input Key and Value types.
func MapAs[K comparable, V any](v reflect.Value) (map[K]V, error) {
	src := indirectValueTilRoot(v)
	if !src.IsValid() || src.Kind() != reflect.Map {
		return nil, fmt.Errorf("%w: require map type (got %v)", ErrTypeInvalid, v.Type())
	}

	srcType := src.Type()
	dstType := reflect.TypeOf(map[K]V{})
	if dstType == srcType {
		return src.Interface().(map[K]V), nil // nolint: forcetypeassert
	}

	dstKeyType := dstType.Key()
	keyAssignable := dstKeyType.AssignableTo(srcType.Key())
	keyConvertible := !keyAssignable && dstKeyType.ConvertibleTo(srcType.Key())

	dstValType := dstType.Elem()
	valAssignable := dstValType.AssignableTo(srcType.Elem())
	valConvertible := !valAssignable && dstValType.ConvertibleTo(srcType.Elem())

	if !keyAssignable && !keyConvertible {
		return nil, fmt.Errorf("%w: unable to convert key type '%v' -> '%v'",
			ErrTypeUnmatched, srcType.Key(), dstKeyType)
	}

	if !valAssignable && !valConvertible {
		return nil, fmt.Errorf("%w: unable to convert value type '%v' -> '%v'",
			ErrTypeUnmatched, dstType.Elem(), dstValType)
	}

	dst := reflect.MakeMapWithSize(dstType, src.Len())
	iter := src.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		if keyConvertible {
			k = k.Convert(dstKeyType)
		}
		if valConvertible {
			v = v.Convert(dstValType)
		}
		dst.SetMapIndex(k, v)
	}
	return dst.Interface().(map[K]V), nil // nolint: forcetypeassert
}
