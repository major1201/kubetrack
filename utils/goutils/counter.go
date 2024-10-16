package goutils

type Counter[T comparable] struct {
	mp map[T]int
}

func NewCounterFromSlice[T comparable](list []T) Counter[T] {
	c := Counter[T]{
		mp: make(map[T]int),
	}
	c.Update(list)
	return c
}

func (c Counter[T]) Inc(ele T) {
	c.mp[ele]++
}

func (c Counter[T]) Add(ele T, n int) {
	c.mp[ele] += n
}

func (c Counter[T]) Update(list []T) {
	for _, ele := range list {
		c.Inc(ele)
	}
}

func (c Counter[T]) Count() map[T]int {
	return c.mp
}
