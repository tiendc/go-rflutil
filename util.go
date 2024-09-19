package rflutil

import "reflect"

func isKindIn(k reflect.Kind, kinds ...reflect.Kind) bool {
	for _, kk := range kinds {
		if k == kk {
			return true
		}
	}
	return false
}

func indirectValueTilRoot(v reflect.Value) reflect.Value {
	for {
		k := v.Kind()
		if k == reflect.Pointer || k == reflect.Interface {
			v = v.Elem()
			if !v.IsValid() {
				return v
			}
			continue
		}
		break
	}
	return v
}

func indirectTypeTilRoot(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

// mapExtend extends a map with another map.
func mapExtend[K comparable, V any, M ~map[K]V](m1, m2 M, newKeysOnly bool) M {
	if m1 == nil {
		m1 = make(M, len(m2))
	}
	for k, v := range m2 {
		if !newKeysOnly {
			m1[k] = v
			continue
		}
		if _, exist := m1[k]; !exist {
			m1[k] = v
		}
	}
	return m1
}

// sliceIndexOf finds index of a value in a slice.
func sliceIndexOf[T comparable, S ~[]T](s S, v T) int {
	for i := range s {
		if s[i] == v {
			return i
		}
	}
	return -1
}

// sliceRemove returns a new slice with removing the given value.
func sliceRemove[T comparable, S ~[]T](s S, v T) S {
	i := sliceIndexOf(s, v)
	if i == -1 {
		return s
	}
	return append(s[:i], s[i+1:]...)
}
