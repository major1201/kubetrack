package goutils

func Map[From any, To any](src []From, fn func(obj From) To) []To {
	res := make([]To, len(src))
	for i := 0; i < len(src); i++ {
		res[i] = fn(src[i])
	}
	return res
}

func MapInPlace[T any](src []T, fn func(obj T) T) {
	for i, element := range src {
		src[i] = fn(element)
	}
}

// Reduce a alice with
func Reduce[T any, R any](slice []T, fn func(T, R) R, zeroValue R) R {
	for _, element := range slice {
		zeroValue = fn(element, zeroValue)
	}
	return zeroValue
}

// Filter a slice with fn
func Filter[T any](slice []T, fn func(T) bool) []T {
	res := make([]T, 0)
	for _, element := range slice {
		if fn(element) {
			res = append(res, element)
		}
	}
	return res
}

func FilterInPlace[T any](slice *[]T, fn func(T) bool) {
	var nextIndex int
	for _, element := range *slice {
		if fn(element) {
			(*slice)[nextIndex] = element
			nextIndex++
		}
	}
	*slice = (*slice)[:nextIndex:nextIndex]
}

func PartialSwap[P1 any, P2 any, R any](fn func(P1, P2) R) func(P2, P1) R {
	return func(p2 P2, p1 P1) R {
		return fn(p1, p2)
	}
}

func IgnoreErr[R any](fn func() (R, error)) func() R {
	return func() R {
		r, _ := fn()
		return r
	}
}

func IgnoreErr1[P any, R any](fn func(P) (R, error)) func(P) R {
	return func(p P) R {
		r, _ := fn(p)
		return r
	}
}

func IgnoreErr2[P1 any, P2 any, R any](fn func(P1, P2) (R, error)) func(P1, P2) R {
	return func(p1 P1, p2 P2) R {
		r, _ := fn(p1, p2)
		return r
	}
}
