package rflutil

import (
	"reflect"
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

	t.Run("#1: field not found", func(t *testing.T) {
		_, err := StructGetField[int](valOf(SS{I: 1, u: 2}), "II", true)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("#2: input not struct", func(t *testing.T) {
		_, err := StructGetField[int](valOf("abc123"), "I", true)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#3: output type unmatched", func(t *testing.T) {
		_, err := StructGetField[uint](valOf(&SS{I: 1, u: 2}), "I", true)
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})

	t.Run("#4: 2 fields have same case-insensitive names", func(t *testing.T) {
		// nolint: unused
		type SS2 struct {
			I int
			i string
			u uint
		}
		_, err := StructGetField[int](valOf(SS2{I: 1, u: 2}), "i", false)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("#5: unexported but can't get address of field", func(t *testing.T) {
		_, err := StructGetField[uint](valOf(SS{I: 1, u: 2}), "u", true)
		assert.ErrorIs(t, err, ErrValueUnaddressable)
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

	t.Run("#5: unexported but can't get address of field", func(t *testing.T) {
		s := SS{I: 1, u: 2}
		err := StructSetField(valOf(s), "u", uint(11), true)
		assert.ErrorIs(t, err, ErrValueUnaddressable)
	})
}

func Test_StructListFields(t *testing.T) {
	type SS struct {
		I int  `mytag:"ii"`
		u uint `mytag:"uu"`
	}

	t.Run("#1: success", func(t *testing.T) {
		v, err := StructListFields(valOf(&SS{}), false)
		assert.Nil(t, err)
		assert.Equal(t, []string{"I"}, v)
	})

	t.Run("#2: success - flatten embedded struct", func(t *testing.T) {
		type SS2 struct {
			SS `mytag:"ss"`
			I  int `mytag:"ii"`
		}
		v, err := StructListFields(valOf(&SS2{}), true)
		assert.Nil(t, err)
		assert.Equal(t, []string{"I"}, v)
	})

	t.Run("#3: success - multi-level embedded structs", func(t *testing.T) {
		type SS2 struct {
			SS `mytag:"ss"`
			I  int `mytag:"ii"`
		}
		type SS3 struct {
			I2  int `mytag:"ii2"`
			SS2 `mytag:"ss2"`
		}
		v, err := StructListFields(valOf(SS3{}), true)
		assert.Nil(t, err)
		assert.Equal(t, []string{"I2", "I"}, v)
	})

	t.Run("#4: success - embedded structs, no flatten embedded fields", func(t *testing.T) {
		type SS2 struct {
			SS `mytag:"ss"`
			I  int `mytag:"ii"`
		}
		type SS3 struct {
			I2  int `mytag:"ii2"`
			SS2 `mytag:"ss2"`
		}
		v, err := StructListFields(valOf(SS3{}), false)
		assert.Nil(t, err)
		assert.Equal(t, []string{"I2", "SS2"}, v)
	})
}

func Test_StructListFields_failure(t *testing.T) {
	t.Run("#1: input not struct", func(t *testing.T) {
		_, err := StructListFields(valOf("abc123"), false)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})
}
