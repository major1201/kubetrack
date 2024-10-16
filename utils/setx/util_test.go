package setx

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDifference(t *testing.T) {
	is := assert.New(t)

	hashSet := NewHashSetFromSlice([]int{3, 1, 2, 5, 4})
	hashSet2 := NewHashSetFromSlice([]int{4, 1, 3, 7, 9})

	differenceRet := Difference[int](hashSet, hashSet2)
	is.Equal(2, differenceRet.Len())
	diffSlice := differenceRet.ToSlice()
	sort.Ints(diffSlice)
	is.Equal([]int{2, 5}, diffSlice)
	hashSetDiffCopy := NewHashSetWithCap[int](0)
	differenceRet.CopyInto(hashSetDiffCopy)
	copySlice := hashSetDiffCopy.ToSlice()
	sort.Ints(copySlice)
	is.Equal([]int{2, 5}, copySlice)
	is.True(differenceRet.Equal(NewHashSetFromSlice([]int{5, 2})))
	is.False(differenceRet.Equal(NewHashSetFromSlice([]int{5, 2, 1})))
	is.True(differenceRet.Contains(2))
	is.False(differenceRet.Contains(3))
	is.False(differenceRet.ContainsAll([]int{2, 5, 1}))
	is.True(differenceRet.ContainsAll([]int{2, 5}))

	diffRangeSlice := make([]int, 0)
	differenceRet.Range(func(_ int, elem int) bool {
		diffRangeSlice = append(diffRangeSlice, elem)
		return true
	})
	sort.Ints(diffRangeSlice)
	is.Equal([]int{2, 5}, diffRangeSlice)

	copyIntoSet := NewHashSetWithCap[int](0)
	DifferenceIntoSet[int](hashSet, hashSet2, copyIntoSet)
	diffSlice2 := copyIntoSet.ToSlice()
	sort.Ints(diffSlice2)
	is.Equal([]int{2, 5}, diffSlice)

}

func TestIntersection(t *testing.T) {
	is := assert.New(t)

	hashSet := NewHashSetFromSlice([]int{3, 1, 2, 5, 4})
	orderedHashSet := NewHashSetFromSlice([]int{4, 1, 3, 7, 9})

	intersectionRet := Intersection[int](hashSet, orderedHashSet)
	is.Equal(3, intersectionRet.Len())
	intersectionSlice := intersectionRet.ToSlice()
	sort.Ints(intersectionSlice)
	is.Equal([]int{1, 3, 4}, intersectionSlice)
	is.True(intersectionRet.Equal(NewHashSetFromSlice([]int{3, 4, 1})))
	is.True(intersectionRet.Contains(1))
	is.False(intersectionRet.Contains(5))
	is.False(intersectionRet.ContainsAll([]int{3, 1, 5}))
	is.True(intersectionRet.ContainsAll([]int{3, 4}))

	intersectionRangeSlice := make([]int, 0)
	intersectionRet.Range(func(_ int, elem int) bool {
		intersectionRangeSlice = append(intersectionRangeSlice, elem)
		return true
	})
	sort.Ints(intersectionRangeSlice)
	is.Equal([]int{1, 3, 4}, intersectionRangeSlice)

	copyIntoSet := NewHashSetWithCap[int](0)
	IntersectionIntoSet[int](hashSet, orderedHashSet, copyIntoSet)
	intersectionSlice2 := copyIntoSet.ToSlice()
	sort.Ints(intersectionSlice2)
	is.Equal([]int{1, 3, 4}, intersectionSlice2)

}

func TestUnion(t *testing.T) {
	is := assert.New(t)

	hashSet := NewHashSetFromSlice([]int{3, 1, 2, 5, 4})
	orderedHashSet := NewHashSetFromSlice([]int{4, 1, 3, 7, 9})

	unionRet := Union[int](hashSet, orderedHashSet)
	is.Equal(7, unionRet.Len())
	unionSlice := unionRet.ToSlice()
	sort.Ints(unionSlice)
	is.Equal([]int{1, 2, 3, 4, 5, 7, 9}, unionSlice)

	copyIntoSet := NewHashSetWithCap[int](0)
	UnionIntoSet[int](hashSet, orderedHashSet, copyIntoSet)
	unionSlice2 := copyIntoSet.ToSlice()
	sort.Ints(unionSlice2)
	is.Equal([]int{1, 2, 3, 4, 5, 7, 9}, unionSlice2)

	is.True(unionRet.Equal(NewHashSetFromSlice([]int{5, 9, 7, 1, 2, 3, 4})))
	is.True(unionRet.Contains(1))
	is.True(unionRet.Contains(5))
	is.False(unionRet.Contains(10))
	is.False(unionRet.ContainsAll([]int{1, 3, 10}))
	is.True(unionRet.ContainsAll([]int{9, 7, 3, 4}))

	unionRangeSlice1 := make([]int, 0)
	unionRet.Range(func(_ int, elem int) bool {
		if elem%3 == 0 {
			unionRangeSlice1 = append(unionRangeSlice1, elem)
		}
		return true
	})
	sort.Ints(unionRangeSlice1)
	is.Equal([]int{3, 9}, unionRangeSlice1)

	unionRangeSlice2 := make([]int, 0)
	unionRet.Range(func(idx int, elem int) bool {
		if idx <= 3 {
			unionRangeSlice2 = append(unionRangeSlice2, elem)
			return true
		}
		return false
	})
	is.Len(unionRangeSlice2, 4)

}

func TestSuperOfAndSubsetOf(t *testing.T) {
	is := assert.New(t)

	hashSet := NewHashSetFromSlice([]int{3, 1, 2, 5, 4})
	hashSet2 := NewHashSetFromSlice([]int{4, 1, 3, 7, 9})

	is.False(IsSupersetOf[int](hashSet, hashSet2))
	is.False(IsSubsetOf[int](hashSet, hashSet2))

	hashSet2 = NewHashSetWithCap[int](4)
	hashSet2.AddAll([]int{3, 1, 2, 5})
	is.True(IsSupersetOf[int](hashSet, hashSet2))
	is.False(IsSupersetOf[int](hashSet2, hashSet))
	is.True(IsSubsetOf[int](hashSet2, hashSet))
	is.False(IsSubsetOf[int](hashSet, hashSet2))

	hashSet3 := NewHashSetWithCap[int](4)
	hashSet3.AddAll([]int{2, 1, 3, 4})
	is.True(IsSupersetOf[int](hashSet, hashSet3))
	is.False(IsSupersetOf[int](hashSet3, hashSet))
	is.True(IsSubsetOf[int](hashSet3, hashSet))
	is.False(IsSubsetOf[int](hashSet, hashSet3))

}
