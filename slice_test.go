package rflutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SliceLen(t *testing.T) {
	t.Run("#1: nil slice", func(t *testing.T) {
		var s []int
		v, err := SliceLen(valOf(s))
		assert.Nil(t, err)
		assert.Equal(t, 0, v)
	})

	t.Run("#2: empty slice", func(t *testing.T) {
		v, err := SliceLen(valOf([]int{}))
		assert.Nil(t, err)
		assert.Equal(t, 0, v)
	})

	t.Run("#3: success", func(t *testing.T) {
		v, err := SliceLen(valOf([]int{1, 2, 3}))
		assert.Nil(t, err)
		assert.Equal(t, 3, v)
	})
}

func Test_SliceLen_failure(t *testing.T) {
	t.Run("#1: input is not slice", func(t *testing.T) {
		_, err := SliceLen(valOf("abc"))
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})
}

func Test_SliceGet(t *testing.T) {
	t.Run("#1: same type", func(t *testing.T) {
		v, err := SliceGet[int](valOf([]int{1, 2, 3}), 1)
		assert.Nil(t, err)
		assert.Equal(t, 2, v)
	})

	t.Run("#2: output is interface", func(t *testing.T) {
		v, err := SliceGet[any](valOf([]int{1, 2, 3}), 1)
		assert.Nil(t, err)
		assert.Equal(t, 2, v)
	})

	t.Run("#3: output is interface and item is nil", func(t *testing.T) {
		v, err := SliceGet[any](valOf([]any{1, nil, 3}), 1)
		assert.Nil(t, err)
		assert.Nil(t, v)
	})

	t.Run("#4: output is map type", func(t *testing.T) {
		v, err := SliceGet[map[int]int](valOf([]map[int]int{{1: 1}, {2: 2}}), 1)
		assert.Nil(t, err)
		assert.Equal(t, map[int]int{2: 2}, v)
	})

	t.Run("#5: output is map type", func(t *testing.T) {
		v, err := SliceGet[map[int]int](valOf([]map[int]int{{1: 1}, nil}), 1)
		assert.Nil(t, err)
		assert.Nil(t, v)
	})
}

func Test_SliceGet_failure(t *testing.T) {
	t.Run("#1: input is not slice", func(t *testing.T) {
		_, err := SliceGet[int](valOf("abc"), 1)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#2: output type not match", func(t *testing.T) {
		_, err := SliceGet[uint](valOf([]int{1, 2, 3}), 1)
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})

	t.Run("#3: index out of range", func(t *testing.T) {
		_, err := SliceGet[int](valOf([]int{1, 2, 3}), 3)
		assert.ErrorIs(t, err, ErrIndexOutOfRange)
	})
}

func Test_SliceGetAllAs_failure(t *testing.T) {
	t.Run("#1: input is not slice", func(t *testing.T) {
		s, err := SliceAs[string](valOf("abc"))
		assert.Nil(t, s)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#2: type not convertible", func(t *testing.T) {
		s, err := SliceAs[string](valOf([]any{"abc", 123.123}))
		assert.Nil(t, s)
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})
}

func Test_SliceSet(t *testing.T) {
	t.Run("#1: success", func(t *testing.T) {
		s := []int{1, 2, 3}
		err := SliceSet(valOf(s), 1, 22)
		assert.Nil(t, err)
		assert.Equal(t, 22, s[1])
	})

	t.Run("#2: success with pass pointer", func(t *testing.T) {
		s := []int{1, 2, 3}
		err := SliceSet(valOf(&s), 1, 22)
		assert.Nil(t, err)
		assert.Equal(t, 22, s[1])
	})

	t.Run("#3: type is interface", func(t *testing.T) {
		s := []any{1, 2, 3}
		err := SliceSet(valOf(&s), 1, 22)
		assert.Nil(t, err)
		assert.Equal(t, 22, s[1])
	})

	t.Run("#4: type is interface and set nil", func(t *testing.T) {
		s := []any{1, 2, 3}
		err := SliceSet[any](valOf(&s), 1, nil)
		assert.Nil(t, err)
		assert.Equal(t, nil, s[1])
	})
}

func Test_SliceSet_failure(t *testing.T) {
	t.Run("#1: input is not slice", func(t *testing.T) {
		err := SliceSet(valOf("abc"), 1, 11)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#2: type not match", func(t *testing.T) {
		err := SliceSet(valOf([]int{1, 2, 3}), 1, uint(22))
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})

	t.Run("#3: index out of range", func(t *testing.T) {
		err := SliceSet(valOf([]int{1, 2, 3}), 3, 33)
		assert.ErrorIs(t, err, ErrIndexOutOfRange)
	})
}

func Test_SliceAppend(t *testing.T) {
	t.Run("#1: success", func(t *testing.T) {
		s := []int{1, 2, 3}
		s2, err := SliceAppend(valOf(s), 4)
		assert.Nil(t, err)
		assert.Equal(t, []int{1, 2, 3, 4}, s2)
	})

	t.Run("#2: success with pass pointer", func(t *testing.T) {
		s := []int{1, 2, 3}
		s2, err := SliceAppend(valOf(&s), 4)
		assert.Nil(t, err)
		assert.Equal(t, []int{1, 2, 3, 4}, s2)
	})

	t.Run("#3: type is interface", func(t *testing.T) {
		s := []any{1, 2, 3}
		s2, err := SliceAppend(valOf(&s), any(4))
		assert.Nil(t, err)
		assert.Equal(t, []any{1, 2, 3, 4}, s2)
	})

	t.Run("#4: type is interface and add nil", func(t *testing.T) {
		s := []any{1, 2, 3}
		s2, err := SliceAppend(valOf(&s), any(nil))
		assert.Nil(t, err)
		assert.Equal(t, []any{1, 2, 3, nil}, s2)
	})
}

func Test_SliceAppend_failure(t *testing.T) {
	t.Run("#1: input is not slice", func(t *testing.T) {
		_, err := SliceAppend(valOf("abc"), 1)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#2: type not match", func(t *testing.T) {
		_, err := SliceAppend(valOf([]int{1, 2, 3}), uint(1))
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})
}

func Test_SliceGetAll(t *testing.T) {
	t.Run("#1: nil slice", func(t *testing.T) {
		var s []int
		v, err := SliceGetAll(valOf(s))
		assert.Nil(t, err)
		assert.Equal(t, []reflect.Value{}, v)
	})

	t.Run("#2: empty slice", func(t *testing.T) {
		s, err := SliceGetAll(valOf([]int{}))
		assert.Nil(t, err)
		assert.Equal(t, []reflect.Value{}, s)
	})

	t.Run("#3: success", func(t *testing.T) {
		s, err := SliceGetAll(valOf([]int{1, 2, 3}))
		assert.Nil(t, err)
		assert.Equal(t, 3, len(s))
		assert.Equal(t, 1, s[0].Interface())
		assert.Equal(t, 2, s[1].Interface())
		assert.Equal(t, 3, s[2].Interface())
	})
}

func Test_SliceGetAll_failure(t *testing.T) {
	t.Run("#1: input is not slice", func(t *testing.T) {
		s, err := SliceGetAll(valOf("abc"))
		assert.Nil(t, s)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})
}

func Test_SliceAs(t *testing.T) {
	t.Run("#1: input is not slice", func(t *testing.T) {
		s, err := SliceAs[string](valOf([]any{"abc", "123", ""}))
		assert.Nil(t, err)
		assert.Equal(t, []string{"abc", "123", ""}, s)
	})

	t.Run("#2: type not convertible", func(t *testing.T) {
		s, err := SliceAs[string](valOf([]any{"abc", 97}))
		assert.Nil(t, err)
		assert.Equal(t, []string{"abc", "a"}, s)
	})
}

func Test_SliceAs_failure(t *testing.T) {
	t.Run("#1: input is not slice", func(t *testing.T) {
		s, err := SliceAs[string](valOf("abc"))
		assert.Nil(t, s)
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#2: type not convertible", func(t *testing.T) {
		s, err := SliceAs[string](valOf([]any{"abc", 123.123}))
		assert.Nil(t, s)
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})

	t.Run("#3: type not convertible", func(t *testing.T) {
		s, err := SliceAs[string](valOf([]any{"abc", nil}))
		assert.Nil(t, s)
		assert.ErrorIs(t, err, ErrTypeUnmatched)
	})
}
