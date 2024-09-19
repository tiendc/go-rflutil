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

func Test_mapExtend(t *testing.T) {
	// nil input
	assert.Equal(t, map[int]int{1: 111, 3: 333},
		mapExtend((map[int]int)(nil), map[int]int{1: 111, 3: 333}, true))
	assert.Equal(t, map[int]int{1: 111, 3: 333},
		mapExtend(map[int]int{1: 111, 3: 333}, (map[int]int)(nil), false))

	assert.Equal(t, map[int]int{1: 11, 2: 22, 3: 333},
		mapExtend(map[int]int{1: 11, 2: 22}, map[int]int{1: 111, 3: 333}, true))
	assert.Equal(t, map[int]int{1: 111, 2: 22, 3: 333},
		mapExtend(map[int]int{1: 11, 2: 22}, map[int]int{1: 111, 3: 333}, false))
}

func Test_sliceIndexOf(t *testing.T) {
	assert.Equal(t, -1, sliceIndexOf(([]int)(nil), 3))
	assert.Equal(t, -1, sliceIndexOf([]int{}, 3))
	assert.Equal(t, -1, sliceIndexOf([]int{-1, 2, 3}, 0))
	assert.Equal(t, 0, sliceIndexOf([]int{-1, 2, 3}, -1))
	assert.Equal(t, 2, sliceIndexOf([]int{-1, 2, 3}, 3))
}

func Test_sliceRemove(t *testing.T) {
	assert.Nil(t, sliceRemove(([]int)(nil), 3))
	assert.Equal(t, []int{-1, 2, 3}, sliceRemove([]int{-1, 2, 3}, 0))
	assert.Equal(t, []int{2, 3}, sliceRemove([]int{-1, 2, 3}, -1))
	assert.Equal(t, []int{-1, 3}, sliceRemove([]int{-1, 2, 3}, 2))
}
