package setx

var emptyStruct = struct{}{}

// DifferenceIntoSet calc from Difference and copy into set.
func DifferenceIntoSet[T any](s1, s2 SetView[T], dstSet Set[T]) {
	Difference(s1, s2).CopyInto(dstSet)
}

// Difference returns a difference view from set s1 and s2 that contains all elements of s1 that are not also elements of s2.
func Difference[T any](s1, s2 SetView[T]) SetView[T] {
	return &differenceView[T]{
		s1: s1,
		s2: s2,
	}
}

// IntersectionIntoSet calc from Difference and copy into set.
func IntersectionIntoSet[T any](s1, s2 SetView[T], dstSet Set[T]) {
	Intersection(s1, s2).CopyInto(dstSet)
}

// Intersection returns an intersection view of two sets.
func Intersection[T any](s1, s2 SetView[T]) SetView[T] {
	return &intersectionView[T]{
		s1: s1,
		s2: s2,
	}
}

// UnionIntoSet calc from Difference and copy into set.
func UnionIntoSet[T any](s1, s2 SetView[T], dstSet Set[T]) {
	Union(s1, s2).CopyInto(dstSet)
}

// Union unions two sets and returns a view.
func Union[T any](s1, s2 SetView[T]) SetView[T] {
	return &unionView[T]{
		s1: s1,
		s2: s2,
	}
}

// IsSubsetOf Determines if every element in child is in parent.
func IsSubsetOf[T any](child, parent SetView[T]) bool {
	if child.Len() > parent.Len() {
		return false
	}

	count := 0
	child.Range(func(_ int, elem T) bool {
		if parent.Contains(elem) {
			count++
			return true
		}
		return false
	})

	return count == child.Len()
}

// IsSupersetOf Determines if every element in child is in parent.
func IsSupersetOf[T any](parent, child SetView[T]) bool {
	return IsSubsetOf(child, parent)
}

func differenceLength[T any](s1, s2 SetView[T]) int {
	length := 0
	s1.Range(func(_ int, elem T) bool {
		if !s2.Contains(elem) {
			length++
		}
		return true
	})
	return length
}

func differenceElems[T any](s1, s2 SetView[T]) []T {
	elems := make([]T, 0)
	s1.Range(func(_ int, elem T) bool {
		if !s2.Contains(elem) {
			elems = append(elems, elem)
		}
		return true
	})
	return elems
}

func intersectionLength[T any](s1, s2 SetView[T]) int {
	length := 0
	s1.Range(func(_ int, elem T) bool {
		if s2.Contains(elem) {
			length++
		}
		return true
	})
	return length
}

func intersectionElems[T any](s1, s2 SetView[T]) []T {
	elems := make([]T, 0)
	s1.Range(func(_ int, elem T) bool {
		if s2.Contains(elem) {
			elems = append(elems, elem)
		}
		return true
	})
	return elems
}

func unionLength[T any](s1, s2 SetView[T]) int {
	length := 0
	s1.Range(func(_ int, elem T) bool {
		length++
		return true
	})
	s2.Range(func(_ int, elem T) bool {
		if !s1.Contains(elem) {
			length++
		}
		return true
	})
	return length
}

func unionElems[T any](s1, s2 SetView[T]) []T {
	elems := make([]T, 0)
	s1.Range(func(_ int, elem T) bool {
		elems = append(elems, elem)
		return true
	})
	s2.Range(func(_ int, elem T) bool {
		if !s1.Contains(elem) {
			elems = append(elems, elem)
		}
		return true
	})
	return elems
}
