package funcx

// Partial1 will apply 1 param to the passed function (first parameter).
// The function returns another function which can be called with the rest of the parameters.
func Partial1[P any, R any](fn func(P) R, p P) func() R {
	return func() R {
		return fn(p)
	}
}

// Partial1E will apply 1 param to the passed function (first parameter).
// The function returns another function which can be called with the rest of the parameters.
func Partial1E[P any, R any](fn func(P) (R, error), p P) func() (R, error) {
	return func() (R, error) {
		return fn(p)
	}
}

// Partial1Of2 will apply 1 param from head to the passed function.
// The function returns another function which can be called with the rest of the parameters.
func Partial1Of2[P1 any, P2 any, R any](fn func(P1, P2) R, p1 P1) func(P2) R {
	return func(p2 P2) R {
		return fn(p1, p2)
	}
}

// Partial1Of2E will apply 1 param from head to the passed function.
// The function returns another function which can be called with the rest of the parameters.
func Partial1Of2E[P1 any, P2 any, R any](fn func(P1, P2) (R, error), p1 P1) func(P2) (R, error) {
	return func(p2 P2) (R, error) {
		return fn(p1, p2)
	}
}

// Partial1Of2FromTail will apply 1 param from tail to the passed function.
// The function returns another function which can be called with the rest of the parameters.
func Partial1Of2FromTail[P1 any, P2 any, R any](fn func(P1, P2) R, p2 P2) func(P1) R {
	return func(p1 P1) R {
		return fn(p1, p2)
	}
}

// Partial1Of2FromTailE will apply 1 param from tail to the passed function.
// The function returns another function which can be called with the rest of the parameters.
func Partial1Of2FromTailE[P1 any, P2 any, R any](fn func(P1, P2) (R, error), p2 P2) func(P1) (R, error) {
	return func(p1 P1) (R, error) {
		return fn(p1, p2)
	}
}

// Partial1Of3 will apply 1 param from head to the passed function.
// The function returns another function which can be called with the rest of the parameters.
func Partial1Of3[P1 any, P2 any, P3 any, R any](fn func(P1, P2, P3) R, p1 P1) func(P2, P3) R {
	return func(p2 P2, p3 P3) R {
		return fn(p1, p2, p3)
	}
}

// Partial1Of3E will apply 1 param from head to the passed function.
// The function returns another function which can be called with the rest of the parameters.
func Partial1Of3E[P1 any, P2 any, P3 any, R any](fn func(P1, P2, P3) (R, error), p1 P1) func(P2, P3) (R, error) {
	return func(p2 P2, p3 P3) (R, error) {
		return fn(p1, p2, p3)
	}
}

// Partial1Of3FromTail will apply 1 param from tail to the passed function.
// The function returns another function which can be called with the rest of the parameters.
func Partial1Of3FromTail[P1 any, P2 any, P3 any, R any](fn func(P1, P2, P3) R, p3 P3) func(P1, P2) R {
	return func(p1 P1, p2 P2) R {
		return fn(p1, p2, p3)
	}
}

// Partial1Of3FromTailE will apply 1 param from tail to the passed function.
// The function returns another function which can be called with the rest of the parameters.
func Partial1Of3FromTailE[P1 any, P2 any, P3 any, R any](fn func(P1, P2, P3) (R, error), p3 P3) func(P1, P2) (R, error) {
	return func(p1 P1, p2 P2) (R, error) {
		return fn(p1, p2, p3)
	}
}

// Partial2Of3 will apply 2 param from head to the passed function.
// The function returns another function which can be called with the rest of the parameters.
func Partial2Of3[P1 any, P2 any, P3 any, R any](fn func(P1, P2, P3) R, p1 P1, p2 P2) func(P3) R {
	return func(p3 P3) R {
		return fn(p1, p2, p3)
	}
}

// Partial2Of3E will apply 2 param from head to the passed function.
// The function returns another function which can be called with the rest of the parameters.
func Partial2Of3E[P1 any, P2 any, P3 any, R any](fn func(P1, P2, P3) (R, error), p1 P1, p2 P2) func(P3) (R, error) {
	return func(p3 P3) (R, error) {
		return fn(p1, p2, p3)
	}
}

// Partial2Of3FromTail will apply 2 param from tail to the passed function.
// The function returns another function which can be called with the rest of the parameters.
func Partial2Of3FromTail[P1 any, P2 any, P3 any, R any](fn func(P1, P2, P3) R, p2 P2, p3 P3) func(P1) R {
	return func(p1 P1) R {
		return fn(p1, p2, p3)
	}
}

// Partial2Of3FromTailE will apply 2 param from tail to the passed function.
// The function returns another function which can be called with the rest of the parameters.
func Partial2Of3FromTailE[P1 any, P2 any, P3 any, R any](fn func(P1, P2, P3) (R, error), p2 P2, p3 P3) func(P1) (R, error) {
	return func(p1 P1) (R, error) {
		return fn(p1, p2, p3)
	}
}
