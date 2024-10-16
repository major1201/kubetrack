package goutils

import "unsafe"

// TODO! REMOVE ME WHEN pkg/lang/conv/unsafe.go pkg/lang/v2/conv/unsafe.go supports

// UnsafeBytesToString returns the byte slice as a volatile string
// THIS IS EVIL CODE.
// YOU HAVE BEEN WARNED.
func UnsafeBytesToString(b []byte) string {
	// same as strings.Builder::String()
	return *(*string)(unsafe.Pointer(&b))
}

// UnsafeStringToBytes returns the string as a byte slice
// THIS IS EVIL CODE.
// YOU HAVE BEEN WARNED.
func UnsafeStringToBytes(s string) (b []byte) {
	return *(*[]byte)(unsafe.Pointer(&s))
}
