package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRingBuffer(t *testing.T) {
	ta := assert.New(t)

	rb := NewRingBuffer[int](5)

	ta.Equal(0, rb.Size())
	_, ok := rb.Peak()
	ta.False(ok)

	rb.Write([]int{1, 2})
	ta.Equal(2, rb.Size())
	one, ok := rb.Peak()
	ta.True(ok)
	ta.Equal(1, one)

	// once again should be the same on Peak()
	one, ok = rb.Peak()
	ta.True(ok)
	ta.Equal(1, one)

	arr, ok := rb.PeakN(2)
	ta.True(ok)
	ta.Equal([]int{1, 2}, arr)

	_, ok = rb.PeakN(3)
	ta.False(ok)

	rb.Write([]int{3, 4, 5, 6})

	ta.Equal(5, rb.Size())

	arr, ok = rb.PeakN(5)
	ta.True(ok)
	ta.Equal([]int{2, 3, 4, 5, 6}, arr)

	one, ok = rb.Pop()
	ta.True(ok)
	ta.Equal(2, one)

	one, ok = rb.Pop()
	ta.True(ok)
	ta.Equal(3, one)

	_, ok = rb.PopN(4)
	ta.False(ok)

	arr, ok = rb.PopN(3)
	ta.True(ok)
	ta.Equal([]int{4, 5, 6}, arr)
}
