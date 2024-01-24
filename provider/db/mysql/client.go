package mysql

import (
	"github.com/hdget/hdsdk/v1/errdef"
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/jmoiron/sqlx"
	"reflect"
)

type dbClient struct {
	*sqlx.DB
	builder intf.DbBuilder
}

func newClient(db *sqlx.DB) *dbClient {
	return &dbClient{DB: db}
}

func (s *dbClient) UseBuilder(builder intf.DbBuilder) intf.DbClient {
	s.builder = builder
	return s
}

func (s *dbClient) One(v any) error {
	if s.DB == nil {
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

	return s.DB.Get(v, sqlQuery, sqlArgs)
}

func (s *dbClient) All(v any) error {
	if s.DB == nil {
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
	return s.DB.Select(v, sqlQuery, sqlArgs...)
}

func (s *dbClient) Count() (int64, error) {
	if s.DB == nil {
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
	err = s.DB.Get(&total, sqlQuery, sqlArgs...)
	if err != nil {
		return 0, err
	}
	return total, nil
}
