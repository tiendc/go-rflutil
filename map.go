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
