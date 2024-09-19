package rflutil

import (
	"errors"
	"fmt"
	"reflect"
)

// StructToMap converts a struct to a map.
func StructToMap(v reflect.Value, customTag string, flattenEmbeddedStructs bool) (map[string]any, error) {
	detailsMap, err := structToMapEx(v, customTag, flattenEmbeddedStructs)
	if err != nil {
		return nil, err
	}
	result := make(map[string]any, len(detailsMap))
	for _, detail := range detailsMap {
		result[detail.OutKey] = detail.Value
	}
	return result, nil
}

type structFieldDetail struct {
	Name   string
	OutKey string
	Value  any
}

//nolint:gocognit,gocyclo
func structToMapEx(
	v reflect.Value,
	customTag string,
	flattenEmbeddedStructs bool,
) (result map[string]*structFieldDetail, err error) {
	val := indirectValueTilRoot(v)
	if !val.IsValid() || val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: struct or struct pointer required, got '%v'", ErrTypeInvalid, v.Type())
	}

	parseCustomTag := func(sf *reflect.StructField, v *reflect.Value) (string, error) {
		tag, err := ParseTag(sf, customTag, ",")
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return sf.Name, nil
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
	result = make(map[string]*structFieldDetail, numFields)
	for i := 0; i < numFields; i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		if structField.Anonymous && flattenEmbeddedStructs {
			if !structField.IsExported() && !field.CanAddr() {
				continue
			}
			fieldRootVal := indirectValueTilRoot(field)
			if fieldRootVal.IsValid() && fieldRootVal.Kind() == reflect.Struct {
				embeddedFields, err := structToMapEx(fieldRootVal, customTag, flattenEmbeddedStructs)
				if err != nil {
					return nil, err
				}
				mapExtend(result, embeddedFields, true)
				continue
			}
		}

		if !structField.IsExported() {
			continue
		}
		keyName := structField.Name
		if customTag != "" {
			keyName, err = parseCustomTag(&structField, &field)
			if err != nil {
				return nil, err
			}
		}
		if keyName == "" {
			continue
		}
		result[structField.Name] = &structFieldDetail{
			Name:   structField.Name,
			OutKey: keyName,
			Value:  field.Interface(),
		}
	}
	return result, nil
}
