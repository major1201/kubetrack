package setx

type SetView[T any] interface {
	// Contains returns true if this set contains the specified element.
	Contains(elem T) bool
	// ContainsAll returns true if this set contains all elements of the specified collection.
	ContainsAll(elems []T) bool
	// CopyInto Copies the current contents of this set view into an existing set.
	CopyInto(Set[T])
	Equal(SetView[T]) bool
	Len() int
	// Range iterates elements, when fn returns false, iterates stop.
	// it does not provide any ordering guarantees(it depends on implements).
	Range(fn func(idx int, elem T) bool)
	// ToSlice transfer to slice.
	ToSlice() []T
}

// Set A collection that contains no duplicate elements.
// More formally, sets contain no pair of elements e1 and e2 such that e1.equals(e2), and at most one null element.
type Set[T any] interface {
	SetView[T]

	// Add adds the specified element to this set if it is not already present (optional operation).
	Add(T) bool
	// AddAll adds all of the elements in the specified collection to this set
	// if they're not already present (optional operation).
	AddAll(elems []T) bool
	// Remove removes the specified element from this set if it is present (optional operation).
	Remove(elem T) bool
	// RemoveAll Removes from this set all of its elements that are contained
	// in the specified collection (optional operation).
	RemoveAll(elems []T) bool
	// Clear removes all elements from this set.
	Clear()
}

type SortedSet[T any] interface {
	Set[T]

	// First returns a view of the portion of this set whose elements are strictly less than toElement.
	First() T
	// Last returns the last (highest) element currently in this set.
	Last() T
	// Sub returns a view of the portion of this set whose elements range from fromElement, inclusive, to toElement, exclusive.
	Sub(fromElem, to T)
	// HeadSet returns a view of the portion of this set whose elements are strictly less than toElement.
	HeadSet(toElem T) SortedSet[T]
	// TailSet returns a view of the portion of this set whose elements are greater than or equal to fromElement.
	TailSet(fromElem T) SortedSet[T]
}
