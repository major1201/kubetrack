package funcx

// Must returns fn()'s first return if error is nil, otherwise panics with error.
func Must[T any](fn func() (T, error)) T {
	res, err := fn()
	if err != nil {
		panic(err)
	}
	return res
}
