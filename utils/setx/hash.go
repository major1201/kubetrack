package setx

var _ Set[int] = (*HashSet[int])(nil)

type HashSet[T comparable] struct {
	data map[T]struct{}
}

func NewHashSetFromSlice[T comparable](s []T) *HashSet[T] {
	ss := NewHashSetWithCap[T](len(s))
	for _, elem := range s {
		ss.data[elem] = emptyStruct
	}
	return ss
}

func NewHashSetWithCap[T comparable](capacity int) *HashSet[T] {
	data := make(map[T]struct{}, capacity)
	return &HashSet[T]{data: data}
}

// region implements SetView

func (s *HashSet[T]) Contains(elem T) bool {
	_, exists := s.data[elem]
	return exists
}

func (s *HashSet[T]) ContainsAll(elems []T) bool {
	var exists bool
	for _, elem := range elems {
		_, exists = s.data[elem]
		if !exists {
			return false
		}
	}
	return true
}

func (s *HashSet[T]) CopyInto(toSet Set[T]) {
	for elem := range s.data {
		toSet.Add(elem)
	}
}

func (s *HashSet[T]) Equal(s2 SetView[T]) bool {
	if len(s.data) != s2.Len() {
		return false
	}
	for elem := range s.data {
		if !s2.Contains(elem) {
			return false
		}
	}
	return true
}

func (s *HashSet[T]) Len() int {
	return len(s.data)
}

func (s *HashSet[T]) Range(fn func(idx int, elem T) bool) {
	idx := 0
	for elem := range s.data {
		if !fn(idx, elem) {
			return
		}
		idx++
	}
}

func (s *HashSet[T]) ToSlice() []T {
	sl := make([]T, 0, len(s.data))
	for elem := range s.data {
		sl = append(sl, elem)
	}
	return sl
}

// region implements Set

func (s *HashSet[T]) Add(elem T) bool {
	_, exists := s.data[elem]
	if exists {
		return false
	}
	s.data[elem] = emptyStruct
	return true
}

func (s *HashSet[T]) AddAll(elems []T) bool {
	var exists bool
	added := false
	for _, elem := range elems {
		if _, exists = s.data[elem]; !exists {
			s.data[elem] = emptyStruct
			added = true
		}
	}
	return added
}

func (s *HashSet[T]) Remove(elem T) bool {
	_, exists := s.data[elem]
	if exists {
		delete(s.data, elem)
		return true
	}
	return false
}

func (s *HashSet[T]) RemoveAll(elems []T) bool {
	modified := false
	exists := false
	for _, elem := range elems {
		if _, exists = s.data[elem]; exists {
			delete(s.data, elem)
			modified = true
		}
	}
	return modified
}

func (s *HashSet[T]) Clear() {
	s.data = make(map[T]struct{})
}

func (s *HashSet[T]) Difference(other SetView[T]) *HashSet[T] {
	result := make(map[T]struct{})
	for elem := range s.data {
		if !other.Contains(elem) {
			result[elem] = emptyStruct
		}
	}
	return &HashSet[T]{data: result}
}

// Intersection returns an intersection of two sets as HashSet.
func (s *HashSet[T]) Intersection(other SetView[T]) *HashSet[T] {
	result := make(map[T]struct{})
	for elem := range s.data {
		if other.Contains(elem) {
			result[elem] = emptyStruct
		}
	}
	return &HashSet[T]{data: result}
}

// UnionInplace union two sets and returns self.
func (s *HashSet[T]) UnionInplace(other SetView[T]) *HashSet[T] {
	unionHashSet(s, other)
	return s
}

func (s *HashSet[T]) Union(other SetView[T]) *HashSet[T] {
	copyData := make(map[T]struct{}, len(s.data))
	sc := &HashSet[T]{data: copyData}
	s.CopyInto(sc)
	unionHashSet(sc, other)
	return sc
}

func (s *HashSet[T]) IsSubsetOf(other SetView[T]) bool {
	return IsSubsetOf[T](s, other)
}

func (s *HashSet[T]) IsSupersetOf(other SetView[T]) bool {
	return IsSubsetOf[T](other, s)
}

func unionHashSet[T comparable](toMs *HashSet[T], fromSet SetView[T]) {
	fromSet.Range(func(_ int, elem T) bool {
		toMs.Add(elem)
		return true
	})
}
