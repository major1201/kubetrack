package goutils

import (
	"fmt"
	"math"

	"github.com/major1201/kubetrack/utils/goutils/internal/constraints"
)

func humanReadableBytes[T constraints.Integer](s T, base float64, units []string) string {
	if s < 10 {
		return fmt.Sprintf("%d %s", s, units[0])
	}
	e := math.Floor(math.Log(float64(s)) / math.Log(base))
	suffix := units[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+" %s", val, suffix)
}

// FileSize translate bytes number into human-readable size
func FileSize[T constraints.Integer](s T) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanReadableBytes(s, 1024, units)
}

// MakeRange returns a slice starts from min and ends with max
func MakeRange(min, max int) []int {
	if min > max {
		return nil
	}

	list := make([]int, max-min+1)
	for i := range list {
		list[i] = i + min
	}
	return list
}

// Round round half up
// For example:
//
//	Round(0.363636, 3)  // 0.364
//	Round(0.363636, 2)  // 0.36
//	Round(0.363636, 1)  // 0.4
//	Round(32, 1)        // 30
func Round(val float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Floor(val*p+0.5) / p
}
