package set

import (
	mapset "github.com/deckarep/golang-set"
)

func newSetFromIntSlice(i []int) mapset.Set {
	set := mapset.NewSet()
	for _, item := range i {
		set.Add(item)
	}
	return set
}

func setToIntSlice(set mapset.Set) []int {
	sli := set.ToSlice()
	result := make([]int, len(sli))
	for index, value := range sli {
		result[index] = value.(int)
	}
	return result
}

// 交集
func IntSliceIntersect(a, b []int) []int {
	aSet := newSetFromIntSlice(a)
	bSet := newSetFromIntSlice(b)
	return setToIntSlice(aSet.Intersect(bSet))
}

// 差集
func IntSliceDifference(a, b []int) []int {
	aSet := newSetFromIntSlice(a)
	bSet := newSetFromIntSlice(b)
	return setToIntSlice(aSet.Difference(bSet))
}
