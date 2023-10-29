package rflutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MapLen(t *testing.T) {
	t.Run("#1: nil map", func(t *testing.T) {
		var m map[int]int
		v, err := MapLen(valOf(m))
		assert.Nil(t, err)
		assert.Equal(t, 0, v)
	})

	t.Run("#2: empty map", func(t *testing.T) {
		v, err := MapLen(valOf(map[int]int{}))
		assert.Nil(t, err)
		assert.Equal(t, 0, v)
	})

	t.Run("#3: success", func(t *testing.T) {
		v, err := MapLen(valOf(map[int]int{1: 1, 2: 2, 3: 3}))
		assert.Nil(t, err)
		assert.Equal(t, 3, v)
	})
}

func Test_MapLen_failure(t *testing.T) {
	t.Run("#1: input is not map", func(t *testing.T) {
		_, err := MapLen(valOf("abc"))
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})
}

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

func Test_MapKeys(t *testing.T) {
	t.Run("#1: nil map", func(t *testing.T) {
		var m map[int]uint
		keys, err := MapKeys(valOf(m))
		assert.Nil(t, err)
		assert.Equal(t, []reflect.Value{}, keys)
	})

	t.Run("#2: empty map", func(t *testing.T) {
		keys, err := MapKeys(valOf(map[int]uint{}))
		assert.Nil(t, err)
		assert.Equal(t, []reflect.Value{}, keys)
	})

	t.Run("#3: success", func(t *testing.T) {
		keys, err := MapKeys(valOf(map[int]uint{1: 1, 2: 2, 3: 3}))
		assert.Nil(t, err)
		assert.Equal(t, 3, len(keys))

		mapChk := map[int]struct{}{1: {}, 2: {}, 3: {}}
		for _, k := range keys {
			_, ok := mapChk[k.Interface().(int)]
			assert.True(t, ok)
		}
	})
}

func Test_MapKeys_failure(t *testing.T) {
	t.Run("#1: input is not a map", func(t *testing.T) {
		keys, err := MapKeys(valOf("abc"))
		assert.ErrorIs(t, err, ErrTypeInvalid)
		assert.Nil(t, keys)
	})
}

func Test_MapEntries(t *testing.T) {
	t.Run("#1: nil map", func(t *testing.T) {
		var m map[int]uint
		entries, err := MapEntries(valOf(m))
		assert.Nil(t, err)
		assert.Equal(t, []MapEntry{}, entries)
	})

	t.Run("#2: empty map", func(t *testing.T) {
		entries, err := MapEntries(valOf(map[int]uint{}))
		assert.Nil(t, err)
		assert.Equal(t, []MapEntry{}, entries)
	})

	t.Run("#3: success", func(t *testing.T) {
		entries, err := MapEntries(valOf(map[int]uint{1: 1, 2: 2, 3: 3}))
		assert.Nil(t, err)
		assert.Equal(t, 3, len(entries))

		mapChk := map[int]uint{1: 1, 2: 2, 3: 3}
		for _, e := range entries {
			v, ok := mapChk[e.Key.Interface().(int)]
			assert.True(t, ok)
			assert.Equal(t, e.Value.Interface(), v)
		}
	})
}

func Test_MapEntries_failure(t *testing.T) {
	t.Run("#1: input is not a map", func(t *testing.T) {
		entries, err := MapEntries(valOf("abc"))
		assert.ErrorIs(t, err, ErrTypeInvalid)
		assert.Nil(t, entries)
	})
}
