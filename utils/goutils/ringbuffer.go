package goutils

import (
	"sync"
)

type RingBuffer[T any] struct {
	data     []T
	length   int
	ptrStart int
	ptrEnd   int
	isFull   bool

	mu sync.RWMutex
}

func NewRingBuffer[T any](length int) *RingBuffer[T] {
	return &RingBuffer[T]{
		data:   make([]T, length),
		length: length,
	}
}

func (rb *RingBuffer[T]) Size() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	return rb.size()
}

func (rb *RingBuffer[T]) size() int {
	sz := rb.ptrEnd - rb.ptrStart
	if sz < 0 || rb.isFull {
		sz += rb.length
	}
	return sz
}

func (rb *RingBuffer[T]) Peak() (res T, ok bool) {
	arr, ok := rb.PeakN(1)
	if !ok {
		return
	}
	return arr[0], true
}

func (rb *RingBuffer[T]) PeakN(n int) (res []T, ok bool) {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	return rb.peakN(n)
}

func (rb *RingBuffer[T]) peakN(n int) (res []T, ok bool) {
	if rb.size() < n {
		return nil, false
	}

	ok = true
	res = make([]T, n)
	for i := 0; i < n; i++ {
		realIndex := (i + rb.ptrStart) % rb.length
		res[i] = rb.data[realIndex]
	}
	return
}

func (rb *RingBuffer[T]) Pop() (res T, ok bool) {
	arr, ok := rb.PopN(1)
	if !ok {
		return
	}
	return arr[0], true
}

func (rb *RingBuffer[T]) PopN(n int) (res []T, ok bool) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	return rb.popN(n)
}

func (rb *RingBuffer[T]) popN(n int) (res []T, ok bool) {
	res, ok = rb.peakN(n)
	if !ok {
		return res, ok
	}

	rb.ptrStart = (rb.ptrStart + n) % rb.length
	if rb.isFull {
		rb.isFull = false
	}
	return
}

func (rb *RingBuffer[T]) Write(data []T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	for i := max(0, len(data)-rb.length); i < len(data); i++ {
		rb.data[rb.ptrEnd] = data[i]

		if rb.ptrEnd == rb.ptrStart && rb.isFull {
			rb.ptrStart++
			rb.ptrStart %= rb.length
		}
		rb.ptrEnd++
		rb.ptrEnd %= rb.length
		if rb.ptrEnd == rb.ptrStart && !rb.isFull {
			rb.isFull = true
		}
	}
}
