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
