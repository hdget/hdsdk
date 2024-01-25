package mysql

import (
	"database/sql"
	"github.com/hdget/hdsdk/errdef"
	"github.com/hdget/hdsdk/intf"
	"github.com/jmoiron/sqlx"
	"reflect"
)

type dbClient struct {
	builder intf.DbBuilder
	db      *sqlx.DB
}

func newClient(db *sqlx.DB) *dbClient {
	return &dbClient{db: db}
}

func (s dbClient) UseBuilder(builder intf.DbBuilder) {
	s.builder = builder
}

/* basic api */

func (s dbClient) Exec(query string, args ...any) (sql.Result, error) {
	if s.db == nil {
		return nil, errdef.ErrEmptyDb
	}

	return s.db.Exec(query, args...)
}

func (s dbClient) Query(query string, args ...any) (*sql.Rows, error) {
	if s.db == nil {
		return nil, errdef.ErrEmptyDb
	}

	return s.db.Query(query, args...)
}

func (s dbClient) QueryRow(query string, args ...any) *sql.Row {
	if s.db == nil {
		return nil
	}

	return s.db.QueryRow(query, args...)
}

func (s dbClient) NamedExec(query string, arg any) (sql.Result, error) {
	if s.db == nil {
		return nil, errdef.ErrEmptyDb
	}

	return s.db.NamedExec(query, arg)
}

func (s dbClient) One(v any) error {
	if s.db == nil {
		return errdef.ErrEmptyDb
	}

	if s.builder == nil {
		return errdef.ErrEmptyDbBuilder
	}

	if !reflect.ValueOf(v).CanSet() {
		return errdef.ErrValueNotSettable
	}

	sqlQuery, sqlArgs, err := s.builder.ToSql()
	if err != nil {
		return err
	}

	return s.db.Get(v, sqlQuery, sqlArgs)
}

func (s dbClient) All(v any) error {
	if s.db == nil {
		return errdef.ErrEmptyDb
	}

	if s.builder == nil {
		return errdef.ErrEmptyDbBuilder
	}

	if !reflect.ValueOf(v).CanSet() {
		return errdef.ErrValueNotSettable
	}

	sqlQuery, sqlArgs, err := s.builder.ToSql()
	if err != nil {
		return err
	}
	return s.db.Select(v, sqlQuery, sqlArgs...)
}

func (s dbClient) Count() (int64, error) {
	if s.db == nil {
		return 0, errdef.ErrEmptyDb
	}

	if s.builder == nil {
		return 0, errdef.ErrEmptyDbBuilder
	}

	sqlQuery, sqlArgs, err := s.builder.ToSql()
	if err != nil {
		return 0, err
	}
	var total int64
	err = s.db.Get(&total, sqlQuery, sqlArgs...)
	if err != nil {
		return 0, err
	}
	return total, nil
}
