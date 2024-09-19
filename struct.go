package rflutil

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// StructGetField get struct field value by field name as T type.
// Input should be a struct, a ptr to a struct, or an interface containing a struct.
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
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem() //nolint:gosec
	}

	t, ok := field.Interface().(T)
	if !ok {
		return zeroT, fmt.Errorf("%w: field type is %v (expect %v)",
			ErrTypeUnmatched, field.Type(), reflect.TypeOf(t))
	}
	return t, nil
}

// StructSetField set struct field value by field name as T type.
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
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem() //nolint:gosec
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

// StructListFields lists all fields of a struct with flattening embedded structs option.
func StructListFields(
	v reflect.Value,
	flattenEmbeddedStructs bool,
) ([]string, error) {
	return structListFields(v.Type(), flattenEmbeddedStructs)
}

// structListFields lists all fields of a struct with flattening embedded structs option.
func structListFields(
	t reflect.Type,
	flattenEmbeddedStructs bool,
) ([]string, error) {
	typ := indirectTypeTilRoot(t)
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: struct or struct pointer required, got '%v'", ErrTypeInvalid, t)
	}

	numFields := typ.NumField()
	result := make([]string, 0, numFields)
	for i := 0; i < numFields; i++ {
		structField := typ.Field(i)
		//nolint:nestif
		if structField.Anonymous && flattenEmbeddedStructs {
			fieldType := structField.Type
			fieldRootType := indirectTypeTilRoot(fieldType)
			if fieldRootType.Kind() == reflect.Struct {
				embeddedFields, err := structListFields(fieldRootType, flattenEmbeddedStructs)
				if err != nil {
					return nil, err
				}
				for _, f := range embeddedFields {
					if sliceIndexOf(result, f) == -1 {
						result = append(result, f)
					}
				}
				continue
			}
		}
		if structField.IsExported() {
			result = sliceRemove(result, structField.Name)
			result = append(result, structField.Name)
		}
	}
	return result, nil
}
