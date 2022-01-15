package utils

import "reflect"

// ReverseSlice reverses the original slice
func ReverseSlice(slice interface{}) {
	length := reflect.ValueOf(slice).Len()
	swap := reflect.Swapper(slice)
	for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
