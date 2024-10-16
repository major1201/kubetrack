package goutils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAllEnvs(t *testing.T) {
	assert.NoError(t, os.Setenv("TEST1", "golang"))
	assert.NoError(t, os.Setenv("TEST2", "TEST=golang"))
	assert.NoError(t, os.Setenv("TEST3", ""))

	envs := ListAllEnvs()

	assert.Equal(t, "golang", envs["TEST1"])
	assert.Equal(t, "TEST=golang", envs["TEST2"])
	assert.Equal(t, "", envs["TEST3"])

	_, exist := envs["TEST4"]
	assert.False(t, exist)
}
