package hdboil

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
)

type orderByClause struct {
	tokens []string
}

// OrderBy OrderBy字段加入desc
func OrderBy() *orderByClause {
	return &orderByClause{tokens: make([]string, 0)}
}

func (o *orderByClause) Desc(col string) *orderByClause {
	o.tokens = append(o.tokens, col+" DESC")
	return o
}

func (o *orderByClause) Asc(col string) *orderByClause {
	o.tokens = append(o.tokens, col+" ASC")
	return o
}

func (o orderByClause) Output(args ...any) qm.QueryMod {
	return qm.OrderBy(strings.Join(o.tokens, ","), args...)
}
