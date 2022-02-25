package sql

import (
	"strings"
)

// Builder sql builder
type Builder interface {
	GetConditions() []string
	GetWhereClause() string
}

type BaseSqlBuilder struct {
	Filters map[string]string
}

// BuildWhereClause 创建where子句
func (b BaseSqlBuilder) BuildWhereClause(builder Builder) string {
	conditions := builder.GetConditions()
	if len(conditions) == 0 {
		return "1=1"
	}
	return strings.Join(conditions, " AND ")
}
