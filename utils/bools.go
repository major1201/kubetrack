package utils

import (
	"fmt"
	"math"
	"math/cmplx"
	"reflect"
	"time"

	"github.com/major1201/kubetrack/utils/internal/refutil"
)

func Bool(from any) (bool, error) {
	if T, ok := from.(string); ok {
		return convStrToBool(T)
	} else if T, ok := from.(bool); ok {
		return T, nil
	}

	value := refutil.IndirectVal(reflect.ValueOf(from))
	kind := value.Kind()
	switch {
	case reflect.String == kind:
		return convStrToBool(value.String())
	case refutil.IsKindNumeric(kind):
		if parsed, ok := convNumToBool(kind, value); ok {
			return parsed, nil
		}
	case reflect.Bool == kind:
		return value.Bool(), nil
	case refutil.IsKindLength(kind):
		return value.Len() > 0, nil
	case reflect.Struct == kind && value.CanInterface():
		v := value.Interface()
		if t, ok := v.(time.Time); ok {
			return time.Time{} != t, nil
		}
	}
	return false, newConvErr(from, "bool")
}

// BoolDefault convert value to bool, it will return default value if convert failed.
func BoolDefault(from any, d bool) bool {
	if v, err := Bool(from); err == nil {
		return v
	}
	return d
}

func newConvErr(from interface{}, to string) error {
	return fmt.Errorf("cannot convert %#v (type %[1]T) to %v", from, to)
}

func convStrToBool(v string) (bool, error) {
	// @TODO Need to find a clean way to expose the truth list to be modified by
	// API to allow INTL.
	if 1 > len(v) || len(v) > 5 {
		return false, fmt.Errorf("cannot parse string with len %d as bool", len(v))
	}

	// @TODO lut
	switch v {
	case "1", "t", "T", "true", "True", "TRUE", "y", "Y", "yes", "Yes", "YES":
		return true, nil
	case "0", "f", "F", "false", "False", "FALSE", "n", "N", "no", "No", "NO":
		return false, nil
	}
	return false, fmt.Errorf("cannot parse %#v (type string) as bool", v)
}

func convNumToBool(k reflect.Kind, value reflect.Value) (bool, bool) {
	switch {
	case refutil.IsKindInt(k):
		return value.Int() != 0, true
	case refutil.IsKindUint(k):
		return value.Uint() != 0, true
	case refutil.IsKindFloat(k):
		T := value.Float()
		if math.IsNaN(T) {
			return false, true
		}
		return T != 0, true
	case refutil.IsKindComplex(k):
		T := value.Complex()
		if cmplx.IsNaN(T) {
			return false, true
		}
		return real(T) != 0, true
	}
	return false, false
}
