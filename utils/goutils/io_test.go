package goutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExist(t *testing.T) {
	assert.True(t, IsExist(os.Args[0]))
	assert.True(t, IsExist(filepath.Dir(os.Args[0])))
	assert.False(t, IsExist(filepath.Join(filepath.Dir(os.Args[0]), "must_not_exist")))
}

func TestIsFile(t *testing.T) {
	assert.True(t, IsFile(os.Args[0]))
	assert.False(t, IsFile(filepath.Dir(os.Args[0])))
	assert.False(t, IsFile(filepath.Join(filepath.Dir(os.Args[0]), "must_not_exist")))
}

func TestIsDir(t *testing.T) {
	assert.False(t, IsDir(os.Args[0]))
	assert.True(t, IsDir(filepath.Dir(os.Args[0])))
	assert.False(t, IsDir(filepath.Join(filepath.Dir(os.Args[0]), "must_not_exist")))
}
