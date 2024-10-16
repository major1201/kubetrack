package funcx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testPartial1(a *int16) string {
	*a++
	return fmt.Sprintf("a=%d", *a)
}

func testPartial1E(a *int16) (string, error) {
	return testPartial1(a), errors.New("1")
}

func testPartial1Of2(a *int16, b int64) string {
	*a++
	return fmt.Sprintf("a=%d a+b=%d", *a, int64(*a)+b)
}

func testPartial1Of2E(a *int16, b int64) (string, error) {
	return testPartial1Of2(a, b), errors.New("1Of2")
}

func testPartial1Of3(a *int16, b int64, c string) string {
	*a++
	b *= 2
	return fmt.Sprintf("a=%d b=%d a+b=%d - ", *a, b, int64(*a)+b) + c
}

func testPartial1Of3E(a *int16, b int64, c string) (string, error) {
	return testPartial1Of3(a, b, c), errors.New("1Of3")
}

func testPartial2Of3(a *int16, b int64, c string) string {
	*a++
	b *= 2
	return fmt.Sprintf("a=%d b=%d a+b=%d - ", *a, b, int64(*a)+b) + c
}

func testPartial2Of3E(a *int16, b int64, c string) (string, error) {
	return testPartial2Of3(a, b, c), errors.New("2Of3")
}

func TestPartial1(t *testing.T) {
	is := assert.New(t)

	var err error
	var result string

	a1 := int16(100)
	p1 := Partial1(testPartial1, &a1)
	is.Equal("a=101", p1())
	is.Equal("a=102", p1())

	a1 = int16(100)
	p1e := Partial1E(testPartial1E, &a1)
	result, err = p1e()
	is.Equal("a=101", result)
	is.EqualError(err, "1")

}

func TestPartial1Of2(t *testing.T) {
	is := assert.New(t)

	var err error
	var result string

	a1Of2 := int16(1)
	p1Of2 := Partial1Of2(testPartial1Of2, &a1Of2)
	is.Equal("a=2 a+b=102", p1Of2(100))
	is.Equal("a=3 a+b=103", p1Of2(100))

	a1Of2 = int16(1)
	p1Of2e := Partial1Of2E(testPartial1Of2E, &a1Of2)
	result, err = p1Of2e(100)
	is.Equal("a=2 a+b=102", result)
	is.EqualError(err, "1Of2")

	// tail
	p1Of2FromTail := Partial1Of2FromTail(testPartial1Of2, int64(100))
	a1Of2 = int16(1)
	is.Equal("a=2 a+b=102", p1Of2FromTail(&a1Of2))

	p1Of2FromTailE := Partial1Of2FromTailE(testPartial1Of2E, int64(100))
	a1Of2 = int16(1)
	result, err = p1Of2FromTailE(&a1Of2)
	is.Equal("a=2 a+b=102", result)
	is.EqualError(err, "1Of2")

}

func TestPartial1Of3(t *testing.T) {
	is := assert.New(t)

	var err error
	var result string

	a1Of3 := int16(1)
	p1Of3 := Partial1Of3(testPartial1Of3, &a1Of3)
	is.Equal("a=2 b=16 a+b=18 - abc", p1Of3(8, "abc"))
	is.Equal("a=3 b=16 a+b=19 - abc", p1Of3(8, "abc"))

	a1Of3 = int16(1)
	p1Of3e := Partial1Of3E(testPartial1Of3E, &a1Of3)
	result, err = p1Of3e(8, "abc")
	is.Equal("a=2 b=16 a+b=18 - abc", result)
	is.EqualError(err, "1Of3")

	// tail
	p1Of3FromTail := Partial1Of3FromTail(testPartial1Of3, "abc")
	a1Of3 = int16(1)
	is.Equal("a=2 b=16 a+b=18 - abc", p1Of3FromTail(&a1Of3, 8))

	p1Of3FromTailE := Partial1Of3FromTailE(testPartial1Of3E, "abc")
	a1Of3 = int16(1)
	result, err = p1Of3FromTailE(&a1Of3, 8)
	is.Equal("a=2 b=16 a+b=18 - abc", result)
	is.EqualError(err, "1Of3")

}

func TestPartial2Of3(t *testing.T) {
	is := assert.New(t)

	var err error
	var result string

	a2Of3 := int16(1)
	p2Of3 := Partial2Of3(testPartial2Of3, &a2Of3, 2)
	is.Equal("a=2 b=4 a+b=6 - abc", p2Of3("abc"))
	is.Equal("a=3 b=4 a+b=7 - abc", p2Of3("abc"))

	a2Of3 = int16(1)
	p2Of3e := Partial2Of3E(testPartial2Of3E, &a2Of3, 2)
	result, err = p2Of3e("abc")
	is.Equal("a=2 b=4 a+b=6 - abc", result)
	is.EqualError(err, "2Of3")

	// tail
	p2Of3FromTail := Partial2Of3FromTail(testPartial2Of3, 2, "abc")
	a2Of3 = int16(1)
	is.Equal("a=2 b=4 a+b=6 - abc", p2Of3FromTail(&a2Of3))
	p2Of3FromTailE := Partial2Of3FromTailE(testPartial2Of3E, 2, "abc")
	a2Of3 = int16(1)
	result, err = p2Of3FromTailE(&a2Of3)
	is.Equal("a=2 b=4 a+b=6 - abc", result)
	is.EqualError(err, "2Of3")

}
