package set

import (
	"testing"
)

func TestSetToIntSlice(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	if !equal(intSlice, setToIntSlice(newSetFromIntSlice(intSlice))) {
		t.Error("erro 1")
	}
}

func TestIntSliceIntersect(t *testing.T) {
	sli := IntSliceIntersect([]int{1, 2, 3, 4, 5}, []int{2, 3, 6, 7})
	if !equal([]int{2, 3}, sli) {
		t.Error("erro 1")
	}
}

func TestIntSliceDifference(t *testing.T) {
	if !equal([]int{1, 4, 5}, IntSliceDifference([]int{1, 2, 3, 4, 5}, []int{2, 3, 6, 7})) {
		t.Error("erro 1")
	}
	if !equal([]int{6, 7}, IntSliceDifference([]int{2, 3, 6, 7}, []int{1, 2, 3, 4, 5})) {
		t.Error("erro 2")
	}
}

func equal(a []int, sli []int) bool {
	if len(a) != len(sli) {
		return false
	}
	for _, value := range a {
		if !inIntSlice(value, sli) {
			return false
		}
	}
	return true
}

func inIntSlice(i int, sli []int) bool {
	for _, value := range sli {
		if value == i {
			return true
		}
	}
	return false
}
