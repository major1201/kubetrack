package goutils

import (
	"fmt"
	"regexp"
)

// RegIPv4 is the IPv4 Regexp
const RegIPv4 = "(([2][5][0-5]|[2][0-4][0-9]|[1][0-9]{2}|[1-9][0-9]|[0-9])[.]){3}([2][5][0-5]|[2][0-4][0-9]|[1][0-9]{2}|[1-9][0-9]|[0-9])"

// IsIPv4 tells a string matches IPv4 form or not
func IsIPv4(s string) bool {
	ok, _ := regexp.MatchString(fmt.Sprintf("^%v$", RegIPv4), s)
	return ok
}
