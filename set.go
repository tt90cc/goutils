package utils

import "reflect"

func DiffSet(a []int64, b []int64) []int64 {
	var diff []int64
	mp := make(map[int64]int)
	for _, n := range b {
		if _, ok := mp[n]; !ok {
			mp[n] = 1
		}
	}

	for _, n := range a {
		if _, ok := mp[n]; !ok {
			diff = append(diff, n)
		}
	}

	return diff
}

func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}

	return false
}
