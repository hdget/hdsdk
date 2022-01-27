package tx

import (
	"github.com/hdget/hdsdk"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Transaction struct {
	Tx *sqlx.Tx
}

func NewTransaction(db *sqlx.DB) (*Transaction, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction")
	}

	return &Transaction{Tx: tx}, nil
}

func (t *Transaction) CommitOrRollback(err error) {
	if err != nil {
		t.Rollback()
		return
	}
	t.Commit()
}

func (t *Transaction) Rollback() {
	err := t.Tx.Rollback()
	if err != nil {
		hdsdk.Logger.Error("transaction rollback", "err", err)
	}
}

func (t *Transaction) Commit() {
	err := t.Tx.Commit()
	if err != nil {
		hdsdk.Logger.Error("transaction commit", "err", err)
		t.Rollback()
	}
}
