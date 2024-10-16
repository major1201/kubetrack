package slicex

// Contains returns true if an T is present in a iteratee.
func Contains[T comparable](s []T, v T) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}
