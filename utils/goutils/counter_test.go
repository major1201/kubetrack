package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	ta := assert.New(t)

	c := NewCounterFromSlice[int](nil)
	ta.Len(c.Count(), 0)
	c.Inc(3)
	ta.Equal(map[int]int{3: 1}, c.Count())

	c.Add(3, 2)
	c.Add(4, 2)
	ta.Equal(map[int]int{3: 3, 4: 2}, c.Count())

	c.Update([]int{1, 2, 3, 1})
	ta.Equal(map[int]int{1: 2, 2: 1, 3: 4, 4: 2}, c.Count())
}
