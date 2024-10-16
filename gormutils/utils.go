package gormutils

import (
	"encoding/json"

	"github.com/pkg/errors"
	"gorm.io/datatypes"
)

// ToJsonb convert object to jsonb
func ToJsonb(v any) (jsonb datatypes.JSON, err error) {
	b, err := json.Marshal(v)
	if err != nil {
		err = errors.Wrap(err, "parse json error")
		return
	}
	jsonb = b
	return
}

// MustToJsonb convert object to jsonb and ignore errors
func MustToJsonb(v any) (jsonb datatypes.JSON) {
	jsonb, _ = ToJsonb(v)
	return
}

// FromJsonb convert jsonb object to entity
func FromJsonb(jsonb datatypes.JSON, v any) (err error) {
	b, err := jsonb.MarshalJSON()
	if err != nil {
		err = errors.Wrapf(err, "marshal jsonb failed")
		return
	}

	if err = json.Unmarshal(b, v); err != nil {
		err = errors.Wrapf(err, "failed to unmarshal jsonb")
		return
	}
	return
}
