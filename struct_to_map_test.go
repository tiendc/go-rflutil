package rflutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StructToMap(t *testing.T) {
	type SS struct {
		I int    `json:"i"`
		S string `json:"s,omitempty"`
		U *uint  `json:"-"`
		b bool
	}

	t.Run("#1: failure, input is not struct and struct pointer", func(t *testing.T) {
		_, err := StructToMap(valOf("abc123"), "", true)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#2: success", func(t *testing.T) {
		s := SS{I: 1, S: "2", U: ptrOf(uint(3)), b: true}
		m, err := StructToMap(valOf(&s), "", true)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"I": 1, "S": "2", "U": ptrOf(uint(3))}, m)
	})

	t.Run("#3: success with parsing json", func(t *testing.T) {
		s := SS{I: 1, S: "2", U: ptrOf(uint(3)), b: true}
		m, err := StructToMap(valOf(&s), "json", true)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"i": 1, "s": "2"}, m)
	})

	t.Run("#4: success with parsing json with omitempty", func(t *testing.T) {
		s := SS{I: 1, S: "", U: ptrOf(uint(3)), b: true}
		m, err := StructToMap(valOf(&s), "json", true)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"i": 1}, m)
	})

	t.Run("#5: success with tag not found", func(t *testing.T) {
		type SS struct {
			I int
			S string `json:"s,omitempty"`
			b bool
		}
		s := SS{I: 1, S: "abc", b: true}
		m, err := StructToMap(valOf(&s), "json", true)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"I": 1, "s": "abc"}, m)
	})
}

func Test_StructToMap_embeddedStruct(t *testing.T) {
	type SS1 struct {
		I int    `json:"i"`
		S string `json:"s,omitempty"`
		U *uint  `json:"-"`
		b bool
	}
	type SS2 struct {
		I  int `json:"i"`
		I2 int `json:"i2"`
		SS1
	}
	type SS3 struct {
		S string `json:"s"`
		SS2
	}
	type SS4 struct {
		S string `json:"s"`
		*SS2
	}
	type SS5 struct {
		S   string `json:"s"`
		Sub *SS2
	}

	t.Run("#1: success", func(t *testing.T) {
		s := SS3{}
		m, err := StructToMap(valOf(&s), "", true)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"I": 0, "I2": 0, "S": "", "U": (*uint)(nil)}, m)
	})

	t.Run("#2: success with embedding `nil` struct pointer", func(t *testing.T) {
		s := SS4{SS2: nil}
		m, err := StructToMap(valOf(&s), "", true)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"S": "", "SS2": (*SS2)(nil)}, m)
	})

	t.Run("#3: success with embedding `non-nil` struct pointer", func(t *testing.T) {
		s := SS4{SS2: &SS2{I2: 2}}
		m, err := StructToMap(valOf(&s), "", true)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"I": 0, "I2": 2, "S": "", "U": (*uint)(nil)}, m)
	})

	t.Run("#4: success with embedding struct, no flatten", func(t *testing.T) {
		s := SS3{}
		m, err := StructToMap(valOf(&s), "", false)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"S": "", "SS2": s.SS2}, m)
	})

	t.Run("#5: success with embedding `nil` struct pointer, no flatten", func(t *testing.T) {
		s := SS5{}
		m, err := StructToMap(valOf(&s), "", false)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"S": "", "Sub": (*SS2)(nil)}, m)
	})

	t.Run("#4: success with embedding struct, no flatten", func(t *testing.T) {
		s := SS3{}
		m, err := StructToMap(valOf(&s), "", false)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"S": "", "SS2": s.SS2}, m)
	})
}
