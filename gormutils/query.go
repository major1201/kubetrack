package gormutils

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// QueryImpl DB query model
type QueryImpl struct {
	dbCtx        *DBContextImpl
	whereQueries []expr
	orderBy      any
	offset       int
	limit        int
	rawExpr      *expr
}

type expr struct {
	query any
	args  []any
}

// DBContext returns the DBContext query uses
func (q *QueryImpl) DBContext() DBContext {
	return q.dbCtx
}

// Model sets the query model
func (q *QueryImpl) Model(model any) Query {
	q.dbCtx.db = q.dbCtx.db.Model(model)
	return q
}

// Where sets the query condition
func (q *QueryImpl) Where(query any, args ...any) Query {
	q.whereQueries = append(q.whereQueries, expr{query: query, args: args})
	return q
}

// OrderBy sets the query order by expr
func (q *QueryImpl) OrderBy(exp any) Query {
	q.orderBy = exp
	return q
}

// Offset sets the query offset
func (q *QueryImpl) Offset(offset int) Query {
	q.offset = offset
	return q
}

// Limit sets the query limit
func (q *QueryImpl) Limit(limit int) Query {
	q.limit = limit
	return q
}

// Page sets the query offset and limit with pagination logical
func (q *QueryImpl) Page(page, pageSize int) Query {
	q.Offset((page - 1) * pageSize)
	q.Limit(pageSize)
	return q
}

// Raw sets the query raw SQL
func (q *QueryImpl) Raw(sql string, values ...any) Query {
	q.rawExpr = &expr{
		query: sql,
		args:  values,
	}
	return q
}

func (q *QueryImpl) queried() {
	var rawArgs []string

	for _, wq := range q.whereQueries {
		q.dbCtx.db = q.dbCtx.db.Where(wq.query, wq.args...)
	}

	if q.orderBy != nil {
		if q.rawExpr == nil {
			q.dbCtx.db = q.dbCtx.db.Order(q.orderBy)
		} else {
			rawArgs = append(rawArgs, "order by ?")
			q.rawExpr.args = append(q.rawExpr.args, q.orderBy)
		}
	}

	if q.offset > 0 {
		if q.rawExpr == nil {
			q.dbCtx.db = q.dbCtx.db.Offset(q.offset)
		} else {
			rawArgs = append(rawArgs, "offset ?")
			q.rawExpr.args = append(q.rawExpr.args, q.offset)
		}
	}

	if q.limit > 0 {
		if q.rawExpr == nil {
			q.dbCtx.db = q.dbCtx.db.Limit(q.limit)
		} else {
			rawArgs = append(rawArgs, "limit ?")
			q.rawExpr.args = append(q.rawExpr.args, q.limit)
		}
	}

	if q.rawExpr != nil {
		var sql string
		if rawArgs == nil {
			sql = q.rawExpr.query.(string)
		} else {
			sql = fmt.Sprintf("select * from (%s) _ %s", q.rawExpr.query.(string), strings.Join(rawArgs, " "))
		}

		q.dbCtx.db = q.dbCtx.db.Raw(sql, q.rawExpr.args...)
	}
}

// Find fetches the query records
func (q *QueryImpl) Find(out any) error {
	q.queried()

	q.dbCtx.db = q.dbCtx.db.Find(out)
	return errors.WithStack(q.dbCtx.db.Error)
}

// First return the first item
func (q *QueryImpl) First(out any) (exist bool, err error) {
	q.queried()

	q.dbCtx.db = q.dbCtx.db.First(out)
	exist = !errors.Is(q.dbCtx.db.Error, gorm.ErrRecordNotFound)
	if !exist {
		return
	}
	err = errors.WithStack(q.dbCtx.db.Error)
	return
}

// Count returns the record count
func (q *QueryImpl) Count() (count int64, err error) {
	q.queried()
	q.dbCtx.db = q.dbCtx.db.Count(&count)
	err = errors.WithStack(q.dbCtx.db.Error)
	return
}

// Scan the raw sql results
func (q *QueryImpl) Scan(out any) (err error) {
	q.queried()
	q.dbCtx.db.Scan(out)
	return errors.WithStack(q.dbCtx.db.Error)
}
