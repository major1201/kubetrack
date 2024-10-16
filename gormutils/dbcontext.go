package gormutils

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// DBContextImpl is a DB operation package
type DBContextImpl struct {
	db *gorm.DB

	context context.Context
	txLevel int
}

// NewDBContext return a new DBContextImpl
func NewDBContext() DBContext {
	dbCtx := new(DBContextImpl)
	dbCtx.db = GetDB()
	return dbCtx
}

// SetContext the context of the DBContextImpl
func (dbCtx *DBContextImpl) SetContext(ctx context.Context) {
	dbCtx.context = ctx
}

// Begin a new transaction
func (dbCtx *DBContextImpl) Begin() (err error) {
	if dbCtx.txLevel == 0 {
		dbCtx.db = dbCtx.db.Begin()
		err = errors.WithStack(dbCtx.db.Error)
	}
	dbCtx.txLevel++
	return err
}

// Commit a transaction
func (dbCtx *DBContextImpl) Commit() (err error) {
	if dbCtx.txLevel <= 0 {
		return errors.Errorf("no transaction to commit, txLevel: %d", dbCtx.txLevel)
	}

	dbCtx.txLevel--
	if dbCtx.txLevel == 0 {
		dbCtx.db = dbCtx.db.Commit()
		err = errors.WithStack(dbCtx.db.Error)
	}
	return err
}

// Rollback a transaction
func (dbCtx *DBContextImpl) Rollback() error {
	dbCtx.db = dbCtx.db.Rollback()
	dbCtx.txLevel = 0
	return errors.WithStack(dbCtx.db.Error)
}

// NewQuery return a new QueryImpl
func (dbCtx *DBContextImpl) NewQuery() Query {
	query := new(QueryImpl)
	query.dbCtx = dbCtx
	return query
}

// Create insert the value into database
func (dbCtx *DBContextImpl) Create(v any) error {
	dbCtx.db = dbCtx.db.Create(v)
	return errors.WithStack(dbCtx.db.Error)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (dbCtx *DBContextImpl) Save(v any) error {
	dbCtx.db = dbCtx.db.Save(v)
	return errors.WithStack(dbCtx.db.Error)
}

// First find first record that match given conditions, order by primary key
func (dbCtx *DBContextImpl) First(out any, where ...any) (exist bool, err error) {
	dbCtx.db = dbCtx.db.First(out, where...)
	exist = !errors.Is(dbCtx.db.Error, gorm.ErrRecordNotFound)
	if !exist {
		return
	}
	err = errors.WithStack(dbCtx.db.Error)
	return
}

// Find find records that match given conditions
func (dbCtx *DBContextImpl) Find(out any, where ...any) (err error) {
	dbCtx.db = dbCtx.db.Find(out, where...)
	return errors.WithStack(dbCtx.db.Error)
}

// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
func (dbCtx *DBContextImpl) Delete(value any, where ...any) (err error) {
	dbCtx.db = dbCtx.db.Delete(value, where...)
	return errors.WithStack(dbCtx.db.Error)
}

// Exec execute a raw SQL expression
func (dbCtx *DBContextImpl) Exec(sql string, values ...any) error {
	dbCtx.db = dbCtx.db.Exec(sql, values...)
	return errors.WithStack(dbCtx.db.Error)
}

// Tx starts a transaction
func (dbCtx *DBContextImpl) Tx(tx func(dbCtx DBContext) error) (err error) {
	if err = dbCtx.Begin(); err != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			if errRb := dbCtx.Rollback(); errRb != nil {
				logger.Error(errRb, "rollback transaction error")
			}

			if re, ok := r.(error); ok {
				err = errors.WithMessage(re, "[Recover from Panic]")
			} else {
				err = errors.Errorf("[Recover from Panic]%v", r)
			}
		}
	}()

	if err = tx(dbCtx); err != nil {
		if errRb := dbCtx.Rollback(); errRb != nil {
			return errRb
		}
		return
	}

	if err = dbCtx.Commit(); err != nil {
		return
	}

	return
}

// DefaultDBContext returns a new DBContextImpl if the input dbCtx is nil
func DefaultDBContext(dbCtx DBContext) DBContext {
	if dbCtx == nil {
		return NewDBContext()
	}

	return dbCtx
}
