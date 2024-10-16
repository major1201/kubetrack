package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	ta := assert.New(t)
	ta.Equal("", Trim(""))
	ta.Equal("", Trim(" "))
	ta.Equal("", Trim("  "))
	ta.Equal("golang", Trim("golang "))
	ta.Equal("golang", Trim(" golang"))
	ta.Equal("golang", Trim(" golang "))
	ta.Equal("golang", Trim("  golang  "))
	ta.NotEqual("golang", Trim("  golang \n"))
}

func TestTrimLeft(t *testing.T) {
	ta := assert.New(t)
	ta.Equal("", TrimLeft(""))
	ta.Equal("", TrimLeft(" "))
	ta.Equal("", TrimLeft("  "))
	ta.Equal("golang ", TrimLeft("golang "))
	ta.Equal("golang", TrimLeft(" golang"))
	ta.Equal("golang ", TrimLeft(" golang "))
	ta.Equal("golang  ", TrimLeft("  golang  "))
}

func TestTrimRight(t *testing.T) {
	ta := assert.New(t)
	ta.Equal("", TrimRight(""))
	ta.Equal("", TrimRight(" "))
	ta.Equal("", TrimRight("  "))
	ta.Equal("golang", TrimRight("golang "))
	ta.Equal(" golang", TrimRight(" golang"))
	ta.Equal(" golang", TrimRight(" golang "))
	ta.Equal("  golang", TrimRight("  golang  "))
}

func TestLeftPad(t *testing.T) {
	ta := assert.New(t)
	ta.Equal("golang", LeftPad("golang", "-", 0))
	ta.Equal("golang", LeftPad("golang", "-", 6))
	ta.Equal("----golang", LeftPad("golang", "-", 10))
	ta.Equal("", LeftPad("", "-", 0))
	ta.Equal("-", LeftPad("", "-", 1))
}

func TestRightPad(t *testing.T) {
	ta := assert.New(t)
	ta.Equal("golang", RightPad("golang", "-", 0))
	ta.Equal("golang", RightPad("golang", "-", 6))
	ta.Equal("golang----", RightPad("golang", "-", 10))
	ta.Equal("", RightPad("", "-", 0))
	ta.Equal("-", RightPad("", "-", 1))
}

func TestZeroFill(t *testing.T) {
	ta := assert.New(t)
	ta.Equal("00", ZeroFill("", 2))
	ta.Equal("00", ZeroFill("0", 2))
	ta.Equal("01", ZeroFill("1", 2))
	ta.Equal("11", ZeroFill("11", 2))
	ta.Equal("111", ZeroFill("111", 2))
}

func TestIndex(t *testing.T) {
	ta := assert.New(t)
	ta.Equal(0, Index("Hello, world!", ""))
	ta.Equal(0, Index("Hello, world!", "He"))
	ta.Equal(2, Index("Hello, world!", "llo"))
	ta.Equal(-1, Index("Hello, world!", "not exist"))
	ta.Equal(0, Index("Hello, 世界!", "He"))
	ta.Equal(7, Index("Hello, 世界!", "世"))
	ta.Equal(7, Index("Hello, 世界!", "世界"))
	ta.Equal(-1, Index("Hello, 世界!", "不存在"))
}

func TestIndent(t *testing.T) {
	ta := assert.New(t)
	ta.Equal(`  a:

    b:
      c: bb`, Indent(`a:

  b:
    c: bb`, "  "))
}

func TestFormatMsgAndArgs(t *testing.T) {
	ta := assert.New(t)
	ta.Equal("", FormatMsgAndArgs())
	ta.Equal("hello world", FormatMsgAndArgs("hello world"))
	ta.Equal("12345", FormatMsgAndArgs(12345))
	ta.Equal("hello world", FormatMsgAndArgs("hello %s", "world"))
	ta.Equal("hello world 12345", FormatMsgAndArgs("hello world %d", 12345))
}

func TestDefaultStringIfEmpty(t *testing.T) {
	ta := assert.New(t)
	ta.Equal("hello world", DefaultStringIfEmpty("hello world", "foo"))
	ta.Equal("foo", DefaultStringIfEmpty("", "foo"))
	ta.Equal("", DefaultStringIfEmpty("", ""))
}
