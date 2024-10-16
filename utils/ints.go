package utils

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/major1201/kubetrack/utils/internal/refutil"
)

// IntDefault convert value to int, it will return default value if convert failed.
func IntDefault(from any, d int) int {
	if v, err := Int(from); err == nil {
		return v
	}
	return d
}

func Int(from any) (int, error) {
	if T, ok := from.(int); ok {
		return T, nil
	}

	to64, err := Int64(from)
	if err != nil {
		return 0, newConvErr(from, "int")
	}
	return int(to64), nil
}

func Int64(from any) (int64, error) {
	if T, ok := from.(string); ok {
		return convStrToInt64(T)
	} else if T, ok := from.(int64); ok {
		return T, nil
	}

	value := refutil.IndirectVal(reflect.ValueOf(from))
	kind := value.Kind()
	switch {
	case reflect.String == kind:
		return convStrToInt64(value.String())
	case refutil.IsKindInt(kind):
		return value.Int(), nil
	case refutil.IsKindUint(kind):
		val := value.Uint()
		if val > math.MaxInt64 {
			val = math.MaxInt64
		}
		return int64(val), nil
	case refutil.IsKindFloat(kind):
		return int64(value.Float()), nil
	case refutil.IsKindComplex(kind):
		return int64(real(value.Complex())), nil
	case reflect.Bool == kind:
		if value.Bool() {
			return 1, nil
		}
		return 0, nil
	case refutil.IsKindLength(kind):
		return int64(value.Len()), nil
	}
	return 0, newConvErr(from, "int64")
}

func convStrToInt64(v string) (int64, error) {
	if parsed, err := strconv.ParseInt(v, 10, 0); err == nil {
		return parsed, nil
	}
	if parsed, err := strconv.ParseFloat(v, 64); err == nil {
		return int64(parsed), nil
	}
	if parsed, err := convStrToBool(v); err == nil {
		if parsed {
			return 1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("cannot convert %#v (type string) to int64", v)
}
