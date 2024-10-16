package setx

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashSet(t *testing.T) {
	is := assert.New(t)
	var s, s1, s2, s3, sc *HashSet[int]

	s = NewHashSetWithCap[int](10)
	is.True(s.Add(1))
	is.False(s.Add(1))
	is.Equal(s.Len(), 1)
	is.True(s.AddAll([]int{1, 2, 3}))
	is.False(s.AddAll([]int{1, 2, 3}))
	is.Equal(s.Len(), 3)

	is.True(s.Contains(1))
	is.True(s.ContainsAll([]int{1, 2}))
	is.False(s.Contains(4))
	is.False(s.ContainsAll([]int{2, 4}))

	s1 = NewHashSetFromSlice([]int{2, 3, 4})
	is.Equal(s1.Len(), 3)
	is.False(s.Equal(s1))

	s.Clear()
	is.Equal(s.Len(), 0)
	is.True(s.AddAll([]int{0, 1, 2, 3}))
	is.False(s.Remove(4))
	is.True(s.Remove(2))
	is.False(s.RemoveAll([]int{4, 5, 6}))
	is.True(s.RemoveAll([]int{3, 4}))
	is.Equal(s.Len(), 2)

	s.AddAll([]int{0, 1, 2, 3})

	s2 = NewHashSetWithCap[int](0)
	s2.AddAll([]int{2, 3, 5, 6})

	// difference
	expectedDiffSet := NewHashSetWithCap[int](0)
	expectedDiffSet.AddAll([]int{0, 1})
	diffSet := s.Difference(s2)
	is.True(diffSet.Equal(expectedDiffSet))
	is.False(diffSet.Equal(s))
	is.False(diffSet.Equal(s2))
	is.Equal(diffSet, expectedDiffSet)

	// intersection
	expectedInterSet := NewHashSetWithCap[int](0)
	expectedInterSet.AddAll([]int{2, 3})
	interSet := s.Intersection(s2)
	is.True(interSet.Equal(expectedInterSet))
	is.False(interSet.Equal(s))
	is.False(interSet.Equal(s2))
	is.Equal(interSet, expectedInterSet)

	// union copy
	expectedUnionSet := NewHashSetWithCap[int](0)
	expectedUnionSet.AddAll([]int{0, 1, 2, 3, 5, 6})
	unionSet := s.Union(s2)
	is.True(unionSet.Equal(expectedUnionSet))
	is.False(unionSet.Equal(s))
	is.False(unionSet.Equal(s2))
	is.Equal(unionSet, expectedUnionSet)

	is.True(unionSet.IsSupersetOf(s))
	is.True(unionSet.IsSupersetOf(s2))
	is.False(s.IsSupersetOf(s2))

	is.True(s.IsSubsetOf(unionSet))
	is.True(s2.IsSubsetOf(unionSet))
	is.False(s.IsSupersetOf(s2))
	is.False(s.IsSupersetOf(unionSet))

	s3 = NewHashSetWithCap[int](0)
	s3.AddAll([]int{0, 1, 2, 3, 5, 6})
	is.True(s3.Union(s).Equal(expectedUnionSet))
	// union inplace
	s3.Clear()
	s3.AddAll([]int{0, 1, 2, 3})
	is.True(s3.UnionInplace(s2).Equal(expectedUnionSet))
	is.True(s3.Equal(expectedUnionSet))
	is.Equal(s3, expectedUnionSet)

	// copy
	sc = NewHashSetWithCap[int](0)
	s3.CopyInto(sc)
	is.True(s3.Equal(sc))
	is.Equal(s3.Len(), 6)
	is.Equal(sc.Len(), 6)

	// to slice
	sl := s3.ToSlice()
	sort.Ints(sl)
	is.Equal([]int{0, 1, 2, 3, 5, 6}, sl)
}
