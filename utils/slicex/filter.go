package slicex

// Filter iterates over a collection of T, returning an array of
// all T elements predicate returns truthy for.
func Filter[T any](s []T, cb func(T) bool) []T {
	results := []T{}
	for _, v := range s {
		if cb(v) {
			results = append(results, v)
		}
	}
	return results
}

// FilterInplace iterates over a collection of T,
// returns a slice(point to self'Data) of all T elements predicate returns truthy for.
// the operation occurs inplace, the argument s'Data would be modified.
func FilterInplace[T any](s []T, cb func(T) bool) []T {
	index, cursor, length := 0, -1, len(s)
	for index < length {
		val := s[index]
		if !cb(val) {
			index++
			continue
		}
		if cursor+1 == index {
			index++
			cursor++
			continue
		}
		// move val forward
		s[cursor+1] = val
		cursor++
		index++
	}
	return s[0 : cursor+1]
}
