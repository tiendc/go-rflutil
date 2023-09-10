package rflutil

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func valOf(v any) reflect.Value {
	return reflect.ValueOf(v)
}

func ptrOf[T any](v T) *T {
	return &v
}

func Test_StructGetField(t *testing.T) {
	type SS struct {
		I int
		u uint
	}

	t.Run("#1: success", func(t *testing.T) {
		v, err := StructGetField[int](valOf(&SS{I: 1, u: 2}), "I", true)
		assert.Nil(t, err)
		assert.Equal(t, 1, v)
	})

	t.Run("#2: case insensitive", func(t *testing.T) {
		v, err := StructGetField[int](valOf(SS{I: 1, u: 2}), "i", false)
		assert.Nil(t, err)
		assert.Equal(t, 1, v)
	})

	t.Run("#3: unexported", func(t *testing.T) {
		v, err := StructGetField[uint](valOf(&SS{I: 1, u: 2}), "u", true)
		assert.Nil(t, err)
		assert.Equal(t, uint(2), v)
	})

	t.Run("#4: unexported - case insensitive", func(t *testing.T) {
		v, err := StructGetField[uint](valOf(&SS{I: 1, u: 2}), "U", false)
		assert.Nil(t, err)
		assert.Equal(t, uint(2), v)
	})
}

func Test_StructGetField_failure(t *testing.T) {
	type SS struct {
		I int
		u uint
	}

	t.Run("#3: field not found", func(t *testing.T) {
		_, err := StructGetField[int](valOf(SS{I: 1, u: 2}), "II", true)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("#4: input not struct", func(t *testing.T) {
		_, err := StructGetField[int](valOf("abc123"), "I", true)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#5: output type unmatched", func(t *testing.T) {
		_, err := StructGetField[uint](valOf(&SS{I: 1, u: 2}), "I", true)
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})

	t.Run("#6: 2 fields have same case-insensitive names", func(t *testing.T) {
		// nolint: unused
		type SS2 struct {
			I int
			i string
			u uint
		}
		_, err := StructGetField[int](valOf(SS2{I: 1, u: 2}), "i", false)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func Test_StructSetField(t *testing.T) {
	type SS struct {
		I int
		u uint
	}

	t.Run("#1: success", func(t *testing.T) {
		s := SS{I: 1, u: 2}
		err := StructSetField(valOf(&s), "I", 11, true)
		assert.Nil(t, err)
		assert.Equal(t, 11, s.I)
	})

	t.Run("#2: case insensitive", func(t *testing.T) {
		s := SS{I: 1, u: 2}
		err := StructSetField(valOf(&s), "i", 11, false)
		assert.Nil(t, err)
		assert.Equal(t, 11, s.I)
	})

	t.Run("#3: unexported field", func(t *testing.T) {
		s := SS{I: 1, u: 2}
		err := StructSetField(valOf(&s), "u", uint(22), true)
		assert.Nil(t, err)
		assert.Equal(t, uint(22), s.u)
	})

	t.Run("#4: unexported field and case insensitive", func(t *testing.T) {
		s := SS{I: 1, u: 2}
		err := StructSetField(valOf(&s), "U", uint(22), false)
		assert.Nil(t, err)
		assert.Equal(t, uint(22), s.u)
	})
}

func Test_StructSetField_failure(t *testing.T) {
	type SS struct {
		I int
		u uint
	}

	t.Run("#1: field not found", func(t *testing.T) {
		s := SS{I: 1, u: 2}
		err := StructSetField(valOf(&s), "II", 1, true)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("#2: input not struct", func(t *testing.T) {
		s := 123456
		err := StructSetField(valOf(&s), "I", 1, true)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#3: field unsettable", func(t *testing.T) {
		s := SS{I: 1, u: 2}
		err := StructSetField(valOf(s), "I", 11, true)
		assert.ErrorIs(t, err, ErrValueUnsettable)
	})

	t.Run("#4: field type unmatched", func(t *testing.T) {
		s := SS{I: 1, u: 2}
		err := StructSetField(valOf(&s), "I", int32(11), true)
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})
}

func Test_StructToMap(t *testing.T) {
	type SS struct {
		I int    `json:"i"`
		S string `json:"s,omitempty"`
		U *uint  `json:"-"`
		b bool
	}

	t.Run("#1: success", func(t *testing.T) {
		s := SS{I: 1, S: "2", U: ptrOf(uint(3)), b: true}
		m, err := StructToMap(valOf(&s), false, nil)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"I": 1, "S": "2", "U": ptrOf(uint(3))}, m)
	})

	t.Run("#2: success with keyFunc", func(t *testing.T) {
		s := SS{I: 1, S: "2", U: ptrOf(uint(3)), b: true}
		m, err := StructToMap(valOf(&s), false, func(name string, isExported bool) string {
			return strings.ToLower(name)
		})
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"i": 1, "s": "2", "u": ptrOf(uint(3)), "b": true}, m)
	})

	t.Run("#3: keyFunc returns empty str", func(t *testing.T) {
		s := SS{I: 1, S: "2", U: ptrOf(uint(3)), b: true}
		m, err := StructToMap(valOf(&s), false, func(name string, isExported bool) string {
			if name == "I" {
				return ""
			}
			return name
		})
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"S": "2", "U": ptrOf(uint(3)), "b": true}, m)
	})

	t.Run("#4: success with parsing json", func(t *testing.T) {
		s := SS{I: 1, S: "2", U: ptrOf(uint(3)), b: true}
		m, err := StructToMap(valOf(&s), true, nil)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"i": 1, "s": "2"}, m)
	})

	t.Run("#5: success with parsing json with omitempty", func(t *testing.T) {
		s := SS{I: 1, S: "", U: ptrOf(uint(3)), b: true}
		m, err := StructToMap(valOf(&s), true, func(name string, isExported bool) string {
			return name
		})
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"i": 1}, m)
	})
}

func Test_StructToMap_failure(t *testing.T) {
	type SS struct {
		I int
		S string
		U *uint
		b bool
	}

	t.Run("#1: unaddressable error", func(t *testing.T) {
		s := SS{I: 1, S: "2", U: ptrOf(uint(3)), b: true}
		_, err := StructToMap(valOf(s), false, func(name string, isExported bool) string {
			return strings.ToLower(name)
		})
		assert.ErrorIs(t, err, ErrValueUnaddressable)
	})

	t.Run("#2: input not struct", func(t *testing.T) {
		_, err := StructToMap(valOf("abc123"), false, nil)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})
}
