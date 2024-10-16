package goutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsIPv4(t *testing.T) {
	ta := assert.New(t)
	ta.True(IsIPv4("192.168.1.1"))
	ta.False(IsIPv4("go192.168.1.1lang"))
}
