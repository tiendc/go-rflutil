package rflutil

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// StructGetField get struct field value by field name as T type
// Input should be a struct, a ptr to a struct, or an interface containing a struct
func StructGetField[T any](v reflect.Value, name string, caseSensitive bool) (T, error) {
	var zeroT T
	val := indirectValueTilRoot(v)
	if !val.IsValid() || val.Kind() != reflect.Struct {
		return zeroT, fmt.Errorf("%w: require struct type (got %v)", ErrTypeInvalid, v.Type())
	}

	sf := structGetField(val, name, caseSensitive)
	if sf == nil {
		return zeroT, fmt.Errorf("%w: field '%s' not found", ErrNotFound, name)
	}
	field := val.Field(sf.Index[0])
	if !sf.IsExported() {
		if !field.CanAddr() {
			return zeroT, fmt.Errorf("%w: accessing unexported field requires it to be addressable",
				ErrValueUnaddressable)
		}
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	}

	t, ok := field.Interface().(T)
	if !ok {
		return zeroT, fmt.Errorf("%w: field type is %v (expect %v)",
			ErrTypeUnmatched, field.Type(), reflect.TypeOf(t))
	}
	return t, nil
}

// StructSetField set struct field value by field name as T type
func StructSetField[T any](v reflect.Value, name string, value T, caseSensitive bool) error {
	val := indirectValueTilRoot(v)
	if !val.IsValid() || val.Kind() != reflect.Struct {
		return fmt.Errorf("%w: require struct type (got %v)", ErrTypeInvalid, v.Type())
	}

	sf := structGetField(val, name, caseSensitive)
	if sf == nil {
		return fmt.Errorf("%w: field '%s' not found", ErrNotFound, name)
	}
	field := val.Field(sf.Index[0])
	if !sf.IsExported() {
		if !field.CanAddr() {
			return fmt.Errorf("%w: accessing unexported field requires it to be addressable",
				ErrValueUnaddressable)
		}
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	}
	if !field.CanSet() {
		return ErrValueUnsettable
	}

	dstVal := reflect.ValueOf(value)
	if !dstVal.Type().AssignableTo(field.Type()) {
		return fmt.Errorf("%w: field type is %v (expect %v)",
			ErrTypeUnmatched, field.Type(), dstVal.Type())
	}

	field.Set(dstVal)
	return nil
}

func structGetField(v reflect.Value, name string, caseSensitive bool) *reflect.StructField {
	if caseSensitive {
		f, ok := v.Type().FieldByName(name)
		if !ok {
			return nil
		}
		return &f
	}
	// NOTE: this may be faster than using strings.EqualFold
	name = strings.ToLower(name)
	f, ok := v.Type().FieldByNameFunc(func(f string) bool {
		return strings.ToLower(f) == name
	})
	if !ok {
		return nil
	}
	return &f
}

// StructToMap convert a struct to a map
// Pass the keyFunc as nil to default to use field name and ignore unexported fields.
// nolint: gocognit
func StructToMap(
	v reflect.Value,
	parseJSONTag bool,
	keyFunc func(name string, isExported bool) string,
) (map[string]any, error) {
	val := indirectValueTilRoot(v)
	if !val.IsValid() || val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: require struct type (got %v)", ErrTypeInvalid, v.Type())
	}

	parseJSONName := func(sf *reflect.StructField, v *reflect.Value) (string, error) {
		tag, err := ParseTag(sf, "json", ",")
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return "", nil
			}
			return "", err
		}
		if tag.Ignored || tag.Name == "" {
			return "", nil
		}
		if tag.HasAttr("omitempty") && (!v.IsValid() || v.IsZero()) {
			return "", nil
		}
		return tag.Name, nil
	}

	typ := val.Type()
	numFields := typ.NumField()
	result := make(map[string]any, numFields)
	for i := 0; i < numFields; i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		name := structField.Name
		if parseJSONTag {
			jsonName, err := parseJSONName(&structField, &field)
			if err != nil {
				return nil, err
			}
			if jsonName == "" {
				continue
			}
			name = jsonName
		}

		if keyFunc == nil {
			if structField.IsExported() {
				result[name] = field.Interface()
			}
			continue
		}

		name = keyFunc(name, structField.IsExported())
		if name == "" {
			continue
		}
		if !structField.IsExported() {
			if !field.CanAddr() {
				return nil, fmt.Errorf("%w: accessing unexported field requires it to be addressable",
					ErrValueUnaddressable)
			}
			field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		}
		result[name] = field.Interface()
	}
	return result, nil
}
