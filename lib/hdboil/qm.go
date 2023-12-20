package hdboil

import "github.com/volatiletech/sqlboiler/v4/queries/qm"

type qmBuilder struct {
	mods []qm.QueryMod
}

func AppendQueryMod(mods ...qm.QueryMod) *qmBuilder {
	return &qmBuilder{
		mods: mods,
	}
}

func (q *qmBuilder) AppendQueryMod(mods ...qm.QueryMod) *qmBuilder {
	q.mods = append(q.mods, mods...)
	return q
}

func (q *qmBuilder) Output() []qm.QueryMod {
	return q.mods
}
