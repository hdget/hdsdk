package hdboil

import "github.com/volatiletech/sqlboiler/v4/queries/qm"

type qmBuilder struct {
	mods []qm.QueryMod
}

func Qm() *qmBuilder {
	return &qmBuilder{
		mods: make([]qm.QueryMod, 0),
	}
}

func (q *qmBuilder) Use(mods ...qm.QueryMod) *qmBuilder {
	q.mods = append(q.mods, mods...)
	return q
}

func (q *qmBuilder) Slice(mods []qm.QueryMod) *qmBuilder {
	q.mods = append(q.mods, mods...)
	return q
}

func (q *qmBuilder) Output() qm.QueryMod {
	return qm.Expr(q.mods...)
}
