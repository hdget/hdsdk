package sqlboiler

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
)

type orderByBuilder struct {
	tokens []string
}

// OrderBy OrderBy字段加入desc
func OrderBy() *orderByBuilder {
	return &orderByBuilder{tokens: make([]string, 0)}
}

func (o *orderByBuilder) Desc(col string) *orderByBuilder {
	o.tokens = append(o.tokens, col+" DESC")
	return o
}

func (o *orderByBuilder) Asc(col string) *orderByBuilder {
	o.tokens = append(o.tokens, col+" ASC")
	return o
}

func (o orderByBuilder) Output(args ...any) qm.QueryMod {
	return qm.OrderBy(strings.Join(o.tokens, ","), args...)
}
