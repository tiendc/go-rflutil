package rflutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValueAs(t *testing.T) {
	t.Run("#1: nil ptr", func(t *testing.T) {
		var s *string
		v, err := ValueAs[*string](valOf(s))
		assert.Nil(t, err)
		assert.Nil(t, v)
	})

	t.Run("#2: empty string", func(t *testing.T) {
		v, err := ValueAs[string](valOf(""))
		assert.Nil(t, err)
		assert.Equal(t, "", v)
	})

	t.Run("#3: str convertible", func(t *testing.T) {
		type Str string
		v, err := ValueAs[Str](valOf("abc"))
		assert.Nil(t, err)
		assert.Equal(t, Str("abc"), v)
	})

	t.Run("#4: number convertible", func(t *testing.T) {
		v, err := ValueAs[int](valOf(123.123))
		assert.Nil(t, err)
		assert.Equal(t, 123, v)
	})

	t.Run("#5: number assignable", func(t *testing.T) {
		v, err := ValueAs[int64](valOf(int64(123)))
		assert.Nil(t, err)
		assert.Equal(t, int64(123), v)
	})

	t.Run("#6: slice assignable", func(t *testing.T) {
		v, err := ValueAs[[]int](valOf([]int{1, 2, 3}))
		assert.Nil(t, err)
		assert.Equal(t, []int{1, 2, 3}, v)
	})

	t.Run("#7: any-type str to str", func(t *testing.T) {
		v, err := ValueAs[string](valOf(any("abc")))
		assert.Nil(t, err)
		assert.Equal(t, "abc", v)
	})

	t.Run("#8: any-type int to int64", func(t *testing.T) {
		v, err := ValueAs[int64](valOf(any(123)))
		assert.Nil(t, err)
		assert.Equal(t, int64(123), v)
	})
}

func Test_ValueAs_failure(t *testing.T) {
	t.Run("#1: str to int", func(t *testing.T) {
		_, err := ValueAs[int](valOf("123"))
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})

	t.Run("#2: ptr to value", func(t *testing.T) {
		s := "123"
		_, err := ValueAs[string](valOf(&s))
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})

	t.Run("#3: different slice element type", func(t *testing.T) {
		_, err := ValueAs[[]int64](valOf([]int{1, 2, 3}))
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})
}
