package goutils

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	list := []string{"1", "2", "3", "4", "5", "6"}
	expect := []string{"111", "222", "333", "444", "555", "666"}
	result := Map(list, func(a string) string {
		return a + a + a
	})

	if !reflect.DeepEqual(expect, result) {
		t.Fatalf("Transform failed: expect %v got %v", expect, result)
	}
}

func TestTransformInPlace(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	expect := []int{3, 6, 9, 12, 15, 18, 21, 24, 27}
	MapInPlace(list, func(a int) int {
		return a * 3
	})

	if !reflect.DeepEqual(expect, list) {
		t.Fatalf("Transform failed: expect %v got %v", expect, list)
	}
}

type Employee struct {
	Name     string
	Age      int
	Vacation int
	Salary   int
}

func TestMapEmployee(t *testing.T) {
	var list = []Employee{
		{"Hao", 44, 0, 8000},
		{"Bob", 34, 10, 5000},
		{"Alice", 23, 5, 9000},
		{"Jack", 26, 0, 4000},
		{"Tom", 48, 9, 7500},
	}
	var expect = []Employee{
		{"Hao", 45, 0, 9000},
		{"Bob", 35, 10, 6000},
		{"Alice", 24, 5, 10000},
		{"Jack", 27, 0, 5000},
		{"Tom", 49, 9, 8500},
	}
	MapInPlace(list, func(e Employee) Employee {
		e.Salary += 1000
		e.Age++
		return e
	})
	if !reflect.DeepEqual(expect, list) {
		t.Fatalf("Transform failed: expect %v got %v", expect, list)
	}
}

func TestReduce(t *testing.T) {
	ta := assert.New(t)

	// case 1
	mul := func(a, b int) int {
		return a * b
	}

	a := make([]int, 10)
	for i := range a {
		a[i] = i + 1
	}

	// Compute 10!
	out := Reduce(a, mul, 1)
	ta.Equal(1*2*3*4*5*6*7*8*9*10, out)

	// case 2
	var list = []Employee{
		{"Hao", 44, 0, 8000},
		{"Bob", 34, 10, 5000},
		{"Alice", 23, 5, 9000},
		{"Jack", 26, 0, 4000},
		{"Tom", 48, 9, 7500},
		{"Marry", 29, 0, 6000},
		{"Mike", 32, 8, 4000},
	}
	result := Reduce(list, func(a, b Employee) Employee {
		return Employee{"Total Salary", 0, 0, a.Salary + b.Salary}
	}, Employee{Salary: 0})
	expect := 43500
	if result.Salary != expect {
		t.Fatalf("expected %v got %v", expect, result)
	}
}

func TestFilter(t *testing.T) {
	isEven := func(a int) bool {
		return a%2 == 0
	}

	ta := assert.New(t)

	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	ta.Equal([]int{2, 4, 6, 8}, Filter(a, isEven))
}

func TestFilterInPlace(t *testing.T) {
	isOddString := func(s string) bool {
		i, _ := strconv.ParseInt(s, 10, 32)
		return i%2 == 1
	}

	ta := assert.New(t)

	s := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	FilterInPlace(&s, isOddString)
	ta.Equal([]string{"1", "3", "5", "7", "9"}, s)
}
