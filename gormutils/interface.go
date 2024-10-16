package gormutils

import "context"

// DBContext indicates the db context interface
type DBContext interface {
	SetContext(ctx context.Context)

	// db
	Create(v any) error
	Save(v any) error
	First(out any, where ...any) (exist bool, err error)
	Find(out any, where ...any) (err error)
	Delete(value any, where ...any) (err error)
	Exec(sql string, values ...any) error
	Tx(tx func(dbCtx DBContext) error) (err error)

	// transaction
	Begin() (err error)
	Commit() (err error)
	Rollback() error

	// query
	NewQuery() Query
}

// Query indicates the db query interface
type Query interface {
	DBContext() DBContext
	Model(model any) Query
	Where(query any, args ...any) Query
	OrderBy(exp any) Query
	Offset(offset int) Query
	Limit(limit int) Query
	Page(page, pageSize int) Query
	Raw(sql string, values ...any) Query
	Find(out any) error
	First(out any) (exist bool, err error)
	Count() (count int64, err error)
	Scan(out any) (err error)
}
