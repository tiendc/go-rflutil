package rflutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isKindIn(t *testing.T) {
	assert.True(t, isKindIn(reflect.Int8, reflect.Int8))
	assert.True(t, isKindIn(reflect.Int, reflect.Int8, reflect.Int, reflect.Float32))
	assert.False(t, isKindIn(reflect.Float64, reflect.Int8, reflect.Int, reflect.Float32))
}

func Test_indirectValueTilRoot(t *testing.T) {
	v := indirectValueTilRoot(valOf(1))
	assert.Equal(t, reflect.Int, v.Kind())
	assert.Equal(t, 1, v.Interface())

	v = indirectValueTilRoot(valOf(ptrOf(ptrOf("a"))))
	assert.Equal(t, reflect.String, v.Kind())
	assert.Equal(t, "a", v.Interface())

	// Nil pointer
	var ptr1 *string
	v = indirectValueTilRoot(valOf(ptr1))
	assert.False(t, v.IsValid())
}
