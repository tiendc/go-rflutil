package rflutil

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Tag struct {
	Name      string
	FieldName string
	Ignored   bool // when name is "-"
	Attrs     map[string]string
}

func (tag *Tag) GetAttr(key string) (string, bool) {
	val, ok := tag.Attrs[key]
	return val, ok
}

func (tag *Tag) GetAttrDefault(key string, defVal string) string {
	if val, ok := tag.Attrs[key]; ok {
		return val
	}
	return defVal
}

func (tag *Tag) HasAttr(key string) bool {
	_, ok := tag.Attrs[key]
	return ok
}

// ParseTag parse tag for the given struct field
func ParseTag(field *reflect.StructField, tagName, delim string) (*Tag, error) {
	tagValue, ok := field.Tag.Lookup(tagName)
	if !ok {
		return nil, fmt.Errorf("%w: struct tag '%s'", ErrNotFound, tagName)
	}

	tag := &Tag{
		FieldName: field.Name,
		Attrs:     map[string]string{},
	}
	tags := strings.Split(tagValue, delim)
	if len(tags) == 0 {
		return tag, nil
	}

	tag.Name = tags[0]
	tag.Ignored = tag.Name == "-"

	for _, tagOpt := range tags[1:] {
		kv := strings.SplitN(tagOpt, "=", 2) // nolint: gomnd
		if len(kv) == 1 {
			tag.Attrs[kv[0]] = ""
		} else {
			tag.Attrs[kv[0]] = kv[1]
		}
	}
	return tag, nil
}

// ParseTagOf parse tag for the struct and field name
func ParseTagOf(v reflect.Value, fieldName, tagName, delim string) (*Tag, error) {
	val := indirectValueTilRoot(v)
	if !val.IsValid() || val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: require struct type (got %v)", ErrTypeInvalid, v.Type())
	}

	field, ok := val.Type().FieldByName(fieldName)
	if !ok {
		return nil, fmt.Errorf("%w: struct field '%s'", ErrNotFound, fieldName)
	}

	return ParseTag(&field, tagName, delim)
}

// ParseTagsOf parse tags of all struct fields
func ParseTagsOf(v reflect.Value, tagName, delim string) ([]*Tag, error) {
	val := indirectValueTilRoot(v)
	if !val.IsValid() || val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: require struct type (got %v)", ErrTypeInvalid, v.Type())
	}

	numFields := v.NumField()
	tags := make([]*Tag, 0, numFields)
	for i := 0; i < numFields; i++ {
		field := v.Type().Field(i)
		tag, err := ParseTag(&field, tagName, delim)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				continue // This is not error, just ignore
			}
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
