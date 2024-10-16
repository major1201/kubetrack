//go:build !windows

package goutils

import (
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSafeExec(t *testing.T) {
	ta := assert.New(t)

	var (
		cmd       *exec.Cmd
		killed    bool
		err       error
		startTime time.Time
	)

	// normal task
	cmd = exec.Command("echo", "dante")
	startTime = time.Now()
	killed, err = SafeExec(cmd, time.Second, func(*exec.Cmd) error {
		outb, e := cmd.CombinedOutput()
		if e != nil {
			return e
		}
		fmt.Println(string(outb))
		return nil
	})
	ta.False(killed)
	ta.Nil(err)
	ta.WithinDuration(startTime, time.Now(), 100*time.Millisecond)

	// timeout
	cmd = exec.Command("sleep", "30")
	startTime = time.Now()
	killed, err = SafeExec(cmd, time.Second, func(*exec.Cmd) error {
		outb, e := cmd.CombinedOutput()
		if e != nil {
			return e
		}
		fmt.Println(string(outb))
		return nil
	})
	ta.True(killed)
	ta.Nil(err)
	ta.WithinDuration(startTime, time.Now(), 1100*time.Millisecond /*1.1s*/)

	// error
	cmd = exec.Command("commandnotfound")
	killed, err = SafeExec(cmd, time.Second, func(*exec.Cmd) error {
		outb, e := cmd.CombinedOutput()
		if e != nil {
			return e
		}
		fmt.Println(string(outb))
		return nil
	})
	ta.False(killed)
	ta.NotNil(err)
}

func TestSafeExecWithCombinedOutput(t *testing.T) {
	ta := assert.New(t)

	var (
		cmd       *exec.Cmd
		killed    bool
		err       error
		startTime time.Time
		outb      []byte
	)

	// normal task
	cmd = exec.Command("echo", "dante")
	startTime = time.Now()
	outb, killed, err = SafeExecWithCombinedOutput(cmd, time.Second)
	ta.False(killed)
	ta.Nil(err)
	ta.Equal([]byte("dante\n"), outb)
	ta.WithinDuration(startTime, time.Now(), 100*time.Millisecond)

	// timeout
	cmd = exec.Command("sleep", "30")
	startTime = time.Now()
	_, killed, err = SafeExecWithCombinedOutput(cmd, time.Second)
	ta.True(killed)
	ta.Nil(err)
	ta.WithinDuration(startTime, time.Now(), 1100*time.Millisecond /*1.1s*/)

	// error
	cmd = exec.Command("commandnotfound")
	outb, killed, err = SafeExecWithCombinedOutput(cmd, time.Second)
	ta.False(killed)
	ta.NotNil(err)
	ta.Equal([]byte(nil), outb)
}
