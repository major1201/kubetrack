package goutils

import (
	"os"
	"strings"
)

// ListAllEnvs list all envs in map type
func ListAllEnvs() map[string]string {
	envs := os.Environ()
	ret := make(map[string]string, len(os.Environ()))
	for _, env := range envs {
		v := strings.SplitN(env, "=", 2)
		ret[v[0]] = v[1]
	}
	return ret
}
