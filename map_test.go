package rflutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MapGet(t *testing.T) {
	t.Run("#1: success", func(t *testing.T) {
		v, err := MapGet[uint](valOf(map[int]uint{1: 1, 2: 2, 3: 3}), 1)
		assert.Nil(t, err)
		assert.Equal(t, uint(1), v)
	})
}

func Test_MapGet_failure(t *testing.T) {
	t.Run("#1: input is not a map", func(t *testing.T) {
		_, err := MapGet[int](valOf("abc"), 1)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#2: value type unmatched", func(t *testing.T) {
		_, err := MapGet[int](valOf(map[int]uint{1: 1, 2: 2, 3: 3}), 1)
		assert.ErrorIs(t, err, ErrTypeUnmatched)

		type T uint
		_, err = MapGet[T](valOf(map[int]uint{1: 1, 2: 2, 3: 3}), 2)
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})

	t.Run("#3: key type unmatched", func(t *testing.T) {
		_, err := MapGet[int](valOf(map[int]uint{1: 1, 2: 2, 3: 3}), int64(1))
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})

	t.Run("#4: key not exist", func(t *testing.T) {
		_, err := MapGet[uint](valOf(map[int]uint{1: 1, 2: 2, 3: 3}), 4)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func Test_MapSet(t *testing.T) {
	t.Run("#1: success", func(t *testing.T) {
		m := map[int]uint{1: 1, 2: 2, 3: 3}
		err := MapSet(valOf(m), 1, uint(11))
		assert.Nil(t, err)
		assert.Equal(t, uint(11), m[1])
	})

	t.Run("#2: success with new key", func(t *testing.T) {
		m := map[int]uint{1: 1, 2: 2, 3: 3}
		err := MapSet(valOf(m), 4, uint(44))
		assert.Nil(t, err)
		assert.Equal(t, uint(44), m[4])
	})
}

func Test_MapSet_failure(t *testing.T) {
	t.Run("#1: input is not a map", func(t *testing.T) {
		err := MapSet(valOf("abc"), 1, 11)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#2: value type unmatched", func(t *testing.T) {
		err := MapSet(valOf(map[int]uint{1: 1, 2: 2, 3: 3}), 1, 11)
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})

	t.Run("#3: key type unmatched", func(t *testing.T) {
		err := MapSet(valOf(map[int]uint{1: 1, 2: 2, 3: 3}), int64(1), uint(11))
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})
}

func Test_MapDelete(t *testing.T) {
	t.Run("#1: success", func(t *testing.T) {
		m := map[int]uint{1: 1, 2: 2, 3: 3}
		err := MapDelete(valOf(m), 1)
		assert.Nil(t, err)
		assert.Equal(t, map[int]uint{2: 2, 3: 3}, m)
	})

	t.Run("#2: success with key not exist, nothing change", func(t *testing.T) {
		m := map[int]uint{1: 1, 2: 2, 3: 3}
		err := MapDelete(valOf(m), 4)
		assert.Nil(t, err)
		assert.Equal(t, map[int]uint{1: 1, 2: 2, 3: 3}, m)
	})
}

func Test_MapDelete_failure(t *testing.T) {
	t.Run("#1: input is not a map", func(t *testing.T) {
		err := MapDelete(valOf("abc"), 1)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#2: key type unmatched", func(t *testing.T) {
		err := MapDelete(valOf(map[uint]uint{1: 1, 2: 2, 3: 3}), 1)
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})
}

func Test_MapAs(t *testing.T) {
	t.Run("#1: same map type", func(t *testing.T) {
		m1 := map[int]uint{1: 1, 2: 2, 3: 3}
		m2, err := MapAs[int, uint](valOf(m1))
		assert.Nil(t, err)
		assert.Equal(t, m1, m2)
	})

	t.Run("#2: success", func(t *testing.T) {
		m2, err := MapAs[int64, int64](valOf(map[int]uint{1: 1, 2: 2, 3: 3}))
		assert.Nil(t, err)
		assert.Equal(t, map[int64]int64{1: 1, 2: 2, 3: 3}, m2)
	})
}

func Test_MapAs_failure(t *testing.T) {
	t.Run("#1: input is not a map", func(t *testing.T) {
		_, err := MapAs[int, uint](valOf("abc"))
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#2: value type unmatched", func(t *testing.T) {
		_, err := MapAs[int, string](valOf(map[int]uint{1: 1, 2: 2, 3: 3}))
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})

	t.Run("#3: key type unmatched", func(t *testing.T) {
		_, err := MapAs[string, uint](valOf(map[int]uint{1: 1, 2: 2, 3: 3}))
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})
}
