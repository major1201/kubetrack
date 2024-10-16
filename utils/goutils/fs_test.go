package goutils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFileInt(t *testing.T) {
	ta := assert.New(t)
	f, err := os.CreateTemp("", "")
	ta.NoError(err, "create temp failed")
	Must(err)
	f.Close()
	defer os.Remove(f.Name())

	getF := func() *os.File {
		f, err = os.OpenFile(f.Name(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		Must(err)
		return f
	}

	{
		f = getF()
		_, err = f.WriteString("123")
		Must(err)
		val, err := ReadFileInt(f.Name())
		Must(err)
		ta.Equal(123, val)
		f.Close()
	}

	{
		f = getF()
		_, err = f.WriteString("123\n")
		Must(err)
		val, err := ReadFileInt(f.Name())
		Must(err)
		ta.Equal(123, val)
		f.Close()
	}

	{
		f = getF()
		_, err = f.WriteString("aaa")
		Must(err)
		val, err := ReadFileInt(f.Name())
		ta.Error(err)
		ta.Equal(0, val)
		f.Close()
	}
}

func TestReadFileAndUnmarshal(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	ta := assert.New(t)
	f, err := os.CreateTemp("", "")
	ta.NoError(err, "create temp failed")
	Must(err)
	f.Close()
	defer os.Remove(f.Name())

	getF := func() *os.File {
		f, err = os.OpenFile(f.Name(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		Must(err)
		return f
	}

	{
		f = getF()
		_, err = f.WriteString(`{"name": "major1201", "value": 123}`)
		Must(err)
		obj := testStruct{}
		err := ReadFileAndUnmarshal(f.Name(), &obj)
		Must(err)
		ta.Equal("major1201", obj.Name)
		ta.Equal(123, obj.Value)
		f.Close()
	}
}

func TestMarshalAndWriteFile(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	ta := assert.New(t)
	f, err := os.CreateTemp("", "")
	ta.NoError(err, "create temp failed")
	Must(err)
	f.Close()
	defer os.Remove(f.Name())

	ts := testStruct{
		Name:  "major1201",
		Value: 123,
	}

	err = MarshalAndWriteFile(ts, f.Name())
	ta.NoError(err, "write file failed")
	Must(err)

	b, err := os.ReadFile(f.Name())
	Must(err)
	ta.Equal([]byte(`{"name":"major1201","value":123}`), b)
}
