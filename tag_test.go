package rflutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseTag(t *testing.T) {
	type SS struct {
		I int    `mytag:"i,optional,k1=v1,omitempty"`
		S string `mytag:"s,optional,k1=v1,k2=v2,omitempty"`
		U uint   `mytag:"-,optional"`
	}
	s := SS{I: 1, S: "hello", U: 2}
	v := valOf(s)

	t.Run("#1: tag not exist", func(t *testing.T) {
		field, ok := v.Type().FieldByName("I")
		assert.True(t, ok)
		_, err := ParseTag(&field, "abc", ",")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("#2: inexact delimiter", func(t *testing.T) {
		field, ok := v.Type().FieldByName("I")
		assert.True(t, ok)
		tag, err := ParseTag(&field, "mytag", " ")
		assert.Nil(t, err)
		assert.Equal(t, "i,optional,k1=v1,omitempty", tag.Name)
	})

	t.Run("#3: ignored tag", func(t *testing.T) {
		field, ok := v.Type().FieldByName("U")
		assert.True(t, ok)
		tag, err := ParseTag(&field, "mytag", ",")
		assert.Nil(t, err)
		assert.Equal(t, "-", tag.Name)
		assert.True(t, tag.Ignored)
	})

	t.Run("#4: success", func(t *testing.T) {
		field, ok := v.Type().FieldByName("S")
		assert.True(t, ok)
		tag, err := ParseTag(&field, "mytag", ",")
		assert.Nil(t, err)
		assert.Equal(t, "s", tag.Name)
		assert.False(t, tag.Ignored)
		assert.Equal(t, map[string]string{
			"optional":  "",
			"k1":        "v1",
			"k2":        "v2",
			"omitempty": "",
		}, tag.Attrs)
		val, _ := tag.GetAttr("k2")
		assert.Equal(t, "v2", val)
		assert.Equal(t, "", tag.GetAttrDefault("optional", "x"))
		assert.Equal(t, "x", tag.GetAttrDefault("not-exist", "x"))
	})
}

func Test_ParseTagOf(t *testing.T) {
	type SS struct {
		I int    `mytag:"i,optional,k1=v1,omitempty"`
		S string `mytag:"s,optional,k1=v1,k2=v2,omitempty"`
		U uint   `mytag:"-,optional"`
	}
	s := SS{I: 1, S: "hello", U: 2}
	v := valOf(s)

	t.Run("#1: invalid input type", func(t *testing.T) {
		_, err := ParseTagOf(valOf("123"), "xField", "abc", ",")
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#2: field not exist", func(t *testing.T) {
		_, err := ParseTagOf(v, "xField", "mytag", ",")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("#3: ignored tag", func(t *testing.T) {
		tag, err := ParseTagOf(v, "U", "mytag", ",")
		assert.Nil(t, err)
		assert.Equal(t, "-", tag.Name)
		assert.True(t, tag.Ignored)
	})

	t.Run("#4: success", func(t *testing.T) {
		tag, err := ParseTagOf(v, "S", "mytag", ",")
		assert.Nil(t, err)
		assert.Equal(t, "s", tag.Name)
		assert.False(t, tag.Ignored)
		assert.Equal(t, map[string]string{
			"optional":  "",
			"k1":        "v1",
			"k2":        "v2",
			"omitempty": "",
		}, tag.Attrs)
	})
}

func Test_ParseTagsOf(t *testing.T) {
	type SS struct {
		I int    `mytag:"i,optional,k1=v1,omitempty"`
		S string `mytag:"s,optional,k1=v1,k2=v2,omitempty"`
		U uint   `mytag:"-,optional"`
		B bool   `tagx:"b"`
	}
	s := SS{I: 1, S: "hello", U: 2}
	v := valOf(s)

	t.Run("#1: invalid input type", func(t *testing.T) {
		_, err := ParseTagsOf(valOf("123"), "abc", ",")
		assert.ErrorIs(t, err, ErrTypeInvalid)
	})

	t.Run("#2: tag not exist", func(t *testing.T) {
		tags, err := ParseTagsOf(v, "abc", ",")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(tags))
	})

	t.Run("#3: success", func(t *testing.T) {
		tags, err := ParseTagsOf(v, "mytag", ",")
		assert.Nil(t, err)
		assert.Equal(t, 3, len(tags))

		assert.Equal(t, "i", tags[0].Name)
		assert.False(t, tags[0].Ignored)
		assert.Equal(t, map[string]string{
			"optional":  "",
			"k1":        "v1",
			"omitempty": "",
		}, tags[0].Attrs)

		assert.Equal(t, "s", tags[1].Name)
		assert.False(t, tags[1].Ignored)
		assert.Equal(t, map[string]string{
			"optional":  "",
			"k1":        "v1",
			"k2":        "v2",
			"omitempty": "",
		}, tags[1].Attrs)

		assert.Equal(t, "-", tags[2].Name)
		assert.True(t, tags[2].Ignored)
		assert.Equal(t, map[string]string{
			"optional": "",
		}, tags[2].Attrs)
	})
}
