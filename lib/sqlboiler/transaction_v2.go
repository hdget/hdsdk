package sqlboiler

import (
	"context"
	"github.com/hdget/hdsdk/v2"
	"github.com/hdget/hdutils/logger"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Transactor interface {
	Executor() boil.Executor
	Finalize(err error)
}

type trans struct {
	tx boil.Transactor
}

func NewTransactor() (Transactor, error) {
	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &trans{tx: tx}, nil
}

func (i *trans) Executor() boil.Executor {
	if i.tx != nil {
		return i.tx
	}
	return boil.GetDB()
}

func (i *trans) Finalize(err error) {
	if i.tx == nil {
		return
	}

	// 处理transaction
	errLogger := logger.Error
	if hdsdk.HasInitialized() {
		errLogger = hdsdk.Logger().Error
	}

	// need commit
	if err != nil {
		e := i.tx.Rollback()
		errLogger("db roll back", "err", err, "rollback", e)
		return
	}

	e := i.tx.Commit()
	if e != nil {
		errLogger("db commit", "err", e)
	}
}
