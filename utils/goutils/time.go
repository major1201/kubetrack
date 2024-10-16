package goutils

import (
	"fmt"
	"strings"
)

// HumanReadableDuration translates seconds into human-readable text
func HumanReadableDuration(diffSeconds int64) string {
	diffArr := make([]string, 0, 7)

	const (
		minute = 60
		hour   = 60 * minute
		day    = 24 * hour
		week   = 7 * day
		month  = 30 * day
		year   = 12 * month
	)

	for {
		switch {
		case diffSeconds <= 0:
			diffSeconds = 0
			diffArr = append(diffArr, "now")
		case diffSeconds < 2:
			diffSeconds = 0
			diffArr = append(diffArr, "1 second")
		case diffSeconds < 1*minute:
			diffArr = append(diffArr, fmt.Sprintf("%d seconds", diffSeconds))
			diffSeconds = 0
		case diffSeconds < 2*minute:
			diffSeconds -= 1 * minute
			diffArr = append(diffArr, "1 minute")
		case diffSeconds < 1*hour:
			diffArr = append(diffArr, fmt.Sprintf("%d minutes", diffSeconds/minute))
			diffSeconds -= diffSeconds / minute * minute
		case diffSeconds < 2*hour:
			diffSeconds -= 1 * hour
			diffArr = append(diffArr, "1 hour")
		case diffSeconds < 1*day:
			diffArr = append(diffArr, fmt.Sprintf("%d hours", diffSeconds/hour))
			diffSeconds -= diffSeconds / hour * hour
		case diffSeconds < 2*day:
			diffSeconds -= 1 * day
			diffArr = append(diffArr, "1 day")
		case diffSeconds < 1*week:
			diffArr = append(diffArr, fmt.Sprintf("%d days", diffSeconds/day))
			diffSeconds -= diffSeconds / day * day
		case diffSeconds < 2*week:
			diffSeconds -= 1 * week
			diffArr = append(diffArr, "1 week")
		case diffSeconds < 1*month:
			diffArr = append(diffArr, fmt.Sprintf("%d weeks", diffSeconds/week))
			diffSeconds -= diffSeconds / week * week
		case diffSeconds < 2*month:
			diffSeconds -= 1 * month
			diffArr = append(diffArr, "1 month")
		case diffSeconds < 1*year:
			diffArr = append(diffArr, fmt.Sprintf("%d months", diffSeconds/month))
			diffSeconds -= diffSeconds / month * month
		case diffSeconds < 2*year:
			diffSeconds -= 1 * year
			diffArr = append(diffArr, "1 year")
		default:
			diffArr = append(diffArr, fmt.Sprintf("%d years", diffSeconds/year))
			diffSeconds = 0
		}
		if diffSeconds == 0 {
			break
		}
	}
	return strings.Join(diffArr, ", ")
}
