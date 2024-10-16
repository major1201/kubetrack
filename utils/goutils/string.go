package goutils

import (
	"cmp"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// EmptyStr is the empty string const
const EmptyStr = ""

// Trim cuts the blanks of a string, the beginning and the end
func Trim(str string) string {
	return strings.Trim(str, " ")
}

// TrimLeft cuts the left side blanks of a string
func TrimLeft(str string) string {
	return strings.TrimLeft(str, " ")
}

// TrimRight cuts the right side blanks of a string
func TrimRight(str string) string {
	return strings.TrimRight(str, " ")
}

// LeftPad pad a string with specified character to the left side
func LeftPad(s string, padStr string, length int) string {
	prefix := EmptyStr
	if len(s) < length {
		prefix = strings.Repeat(padStr, length-len(s))
	}
	return prefix + s
}

// RightPad pad a string with specified character to the right side
func RightPad(s string, padStr string, length int) string {
	postfix := EmptyStr
	if len(s) < length {
		postfix = strings.Repeat(padStr, length-len(s))
	}
	return s + postfix
}

// ZeroFill pad a string(usually a number string) with "0" to the left
func ZeroFill(s string, length int) string {
	const zeroStr = "0"
	return LeftPad(s, zeroStr, length)
}

// Index return the location of a string in another long string, if it doesn't exist, returns -1
// this function supports CJK characters
func Index(s, substr string) int {
	sRune := []rune(s)
	subRune := []rune(substr)
	if len(subRune) > len(sRune) {
		return -1
	}
	for i := 0; i < len(sRune)-len(subRune)+1; i++ {
		if reflect.DeepEqual(sRune[i:i+len(subRune)], subRune) {
			return i
		}
	}
	return -1
}

// Indent inserts prefix at the beginning of each non-empty line of s. The
// end-of-line marker is NL.
func Indent(s, prefix string) string {
	return string(IndentBytes([]byte(s), []byte(prefix)))
}

// IndentBytes inserts prefix at the beginning of each non-empty line of b.
// The end-of-line marker is NL.
func IndentBytes(b, prefix []byte) []byte {
	var res []byte
	bol := true
	for _, c := range b {
		if bol && c != '\n' {
			res = append(res, prefix...)
		}
		res = append(res, c)
		bol = c == '\n'
	}
	return res
}

// FormatMsgAndArgs format msg and args
func FormatMsgAndArgs(msgAndArgs ...any) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}

// UUID returns a random generated UUID string
func UUID() string {
	return strings.Replace(uuid.New().String(), "-", "", 4)
}

// Deprecated: DefaultStringIfEmpty returns a default string if s is empty
// Use cmp.Or instead
func DefaultStringIfEmpty(s, dv string) string {
	return cmp.Or(s, dv)
}
