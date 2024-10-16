package gormutils

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Create insert the value into database
func Create(v any) error {
	db := GetDB()
	db = db.Create(v)
	return errors.WithStack(db.Error)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func Save(v any) error {
	db := GetDB()
	db = db.Save(v)
	return errors.WithStack(db.Error)
}

// First find first record that match given conditions, order by primary key
func First(out any, where ...any) (exist bool, err error) {
	db := GetDB()
	db = db.First(out, where...)
	exist = !errors.Is(db.Error, gorm.ErrRecordNotFound)
	if !exist {
		return
	}
	err = errors.WithStack(db.Error)
	return
}

// Find find records that match given conditions
func Find(out any, where ...any) (err error) {
	db := GetDB()
	db = db.Find(out, where...)
	return errors.WithStack(db.Error)
}

// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
func Delete(value any, where ...any) (err error) {
	db := GetDB()
	db = db.Delete(value, where...)
	return errors.WithStack(db.Error)
}

// Tx runs a transaction in a function
func Tx(tx func(dbCtx DBContext) error) error {
	return NewDBContext().Tx(tx)
}
