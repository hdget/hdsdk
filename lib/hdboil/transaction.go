package hdboil

import (
	"context"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Transaction struct {
	Tx         boil.Transactor
	needCommit bool
}

func NewTransaction(args ...boil.Transactor) (*Transaction, error) {
	var (
		err error
		tx  boil.Transactor
	)

	needCommit := true
	if len(args) > 0 && args[0] != nil {
		tx = args[0]
		// 外部传递过来的transactor我们不需要commit
		needCommit = false
	} else {
		tx, err = boil.BeginTx(context.Background(), nil)
	}
	if err != nil {
		return nil, err
	}

	return &Transaction{Tx: tx, needCommit: needCommit}, nil
}

func (t Transaction) CommitOrRollback(err error) {
	if t.needCommit {
		CommitOrRollback(t.Tx, err)
	}
}
