package goutils

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

func ReadFileInt(file string) (val int, err error) {
	b, err := os.ReadFile(file)
	if err != nil {
		err = errors.Wrapf(err, "read file failed, path=%s", file)
		return
	}
	s := string(bytes.TrimSpace(b))
	val, err = strconv.Atoi(s)
	err = errors.Wrapf(err, "strconv failed=%s", s)
	return
}

func ReadFileAndUnmarshal(file string, v any) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return errors.Wrapf(err, "read file failed, path=%s", file)
	}
	return json.Unmarshal(b, v)
}

func MarshalAndWriteFile(v any, file string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "open file for write failed, path=%s", file)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return errors.Wrapf(err, "marshal object failed")
	}
	_, err = f.Write(b)
	return errors.Wrapf(err, "write to file failed, path=%s", file)
}
