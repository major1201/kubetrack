package goutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHumanReadableDuration(t *testing.T) {
	const (
		minute = 60
		hour   = 60 * minute
		day    = 24 * hour
		week   = 7 * day
		month  = 30 * day
		year   = 12 * month
	)

	ta := assert.New(t)
	ta.Equal("now", HumanReadableDuration(-100))
	ta.Equal("now", HumanReadableDuration(0))
	ta.Equal("1 second", HumanReadableDuration(1))
	ta.Equal("33 seconds", HumanReadableDuration(33))
	ta.Equal("1 minute", HumanReadableDuration(60))
	ta.Equal("1 minute, 37 seconds", HumanReadableDuration(1*minute+37))
	ta.Equal("2 minutes, 1 second", HumanReadableDuration(2*minute+1))
	ta.Equal("1 hour", HumanReadableDuration(1*hour))
	ta.Equal("1 hour, 7 minutes, 3 seconds", HumanReadableDuration(1*hour+7*minute+3))
	ta.Equal("23 hours, 7 minutes, 3 seconds", HumanReadableDuration(23*hour+7*minute+3))
	ta.Equal("7 years", HumanReadableDuration(7*year+5*month+1*week+23*hour+7*minute+3))
}
