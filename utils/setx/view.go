package setx

type differenceView[T any] struct {
	s1 SetView[T]
	s2 SetView[T]
}

func (v *differenceView[T]) Contains(elem T) bool {
	return v.s1.Contains(elem) && !v.s2.Contains(elem)
}

func (v *differenceView[T]) ContainsAll(elems []T) bool {
	for _, elem := range elems {
		if !v.Contains(elem) {
			return false
		}
	}
	return true
}

func (v *differenceView[T]) CopyInto(s Set[T]) {
	s.AddAll(v.calcElems())
}

func (v *differenceView[T]) Len() int {
	return differenceLength(v.s1, v.s2)
}

func (v *differenceView[T]) ToSlice() []T {
	return differenceElems(v.s1, v.s2)
}

func (v *differenceView[T]) Equal(other SetView[T]) bool {
	return equalFromCalcElems(v.calcElems, other)
}

func (v *differenceView[T]) Range(fn func(idx int, elem T) bool) {
	rangeFromCalcElems(v.calcElems, fn)
}

func (v *differenceView[T]) calcElems() []T {
	return differenceElems(v.s1, v.s2)
}

type intersectionView[T any] struct {
	s1 SetView[T]
	s2 SetView[T]
}

func (v *intersectionView[T]) Contains(elem T) bool {
	return v.s1.Contains(elem) && v.s2.Contains(elem)
}

func (v *intersectionView[T]) ContainsAll(elems []T) bool {
	for _, elem := range elems {
		if !v.Contains(elem) {
			return false
		}
	}
	return true
}

func (v *intersectionView[T]) CopyInto(s Set[T]) {
	s.AddAll(v.calcElems())
}

func (v *intersectionView[T]) Len() int {
	return intersectionLength(v.s1, v.s2)
}

func (v *intersectionView[T]) ToSlice() []T {
	return v.calcElems()
}

func (v *intersectionView[T]) Equal(other SetView[T]) bool {
	return equalFromCalcElems(v.calcElems, other)
}

func (v *intersectionView[T]) Range(fn func(idx int, elem T) bool) {
	rangeFromCalcElems(v.calcElems, fn)
}

func (v *intersectionView[T]) calcElems() []T {
	return intersectionElems(v.s1, v.s2)
}

type unionView[T any] struct {
	s1 SetView[T]
	s2 SetView[T]
}

func (v *unionView[T]) Contains(elem T) bool {
	return v.contains(elem)
}

func (v *unionView[T]) contains(elem T) bool {
	return v.s1.Contains(elem) || v.s2.Contains(elem)
}

func (v *unionView[T]) ContainsAll(elems []T) bool {
	for _, elem := range elems {
		if !v.Contains(elem) {
			return false
		}
	}
	return true
}

func (v *unionView[T]) CopyInto(s Set[T]) {
	s.AddAll(v.calcElems())
}

func (v *unionView[T]) Len() int {
	return unionLength(v.s1, v.s2)
}

func (v *unionView[T]) ToSlice() []T {
	return unionElems(v.s1, v.s2)
}

func (v *unionView[T]) Equal(other SetView[T]) bool {
	return equalFromCalcElems(v.calcElems, other)
}

func (v *unionView[T]) Range(fn func(idx int, elem T) bool) {
	rangeFromCalcElems(v.calcElems, fn)
}

func (v *unionView[T]) calcElems() []T {
	return unionElems(v.s1, v.s2)
}

func equalFromCalcElems[T any](calcElems func() []T, other SetView[T]) bool {
	elems := calcElems()
	if other.Len() != len(elems) {
		return false
	}
	return other.ContainsAll(elems)
}

func rangeFromCalcElems[T any](calcElems func() []T, fn func(idx int, elem T) bool) {
	elems := calcElems()
	for idx, elem := range elems {
		if !fn(idx, elem) {
			return
		}
	}
	return
}
