package sqlboiler

import (
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
)

type joinKind int

const (
	joinKindUnknown joinKind = iota
	joinKindInner
	joinKindLeft
	joinKindRight
)

type joinClause struct {
	kind      joinKind
	joinTable string
	asTable   string
	clauses   []string
}

func InnerJoin(joinTable string, args ...string) *joinClause {
	var asTable string
	if len(args) > 0 {
		asTable = args[0]
	}
	return &joinClause{
		kind:      joinKindInner,
		joinTable: joinTable,
		asTable:   asTable,
		clauses:   make([]string, 0),
	}
}

func LeftJoin(joinTable string, args ...string) *joinClause {
	var asTable string
	if len(args) > 0 {
		asTable = args[0]
	}
	return &joinClause{
		kind:      joinKindLeft,
		joinTable: joinTable,
		asTable:   asTable,
		clauses:   make([]string, 0),
	}
}

func (j *joinClause) On(columnOrTableColumn, thatTableColumn string) *joinClause {
	if j.asTable != "" {
		j.clauses = append(j.clauses, fmt.Sprintf("%s AS %s ON %s.%s=%s", j.joinTable, j.asTable, j.asTable, columnOrTableColumn, thatTableColumn))
	} else {
		j.clauses = append(j.clauses, fmt.Sprintf("%s ON %s=%s", j.joinTable, columnOrTableColumn, thatTableColumn))
	}
	return j
}

func (j *joinClause) And(clause string) *joinClause {
	j.clauses = append(j.clauses, clause)
	return j
}

func (j *joinClause) Output(args ...any) qm.QueryMod {
	switch j.kind {
	case joinKindInner:
		return qm.InnerJoin(strings.Join(j.clauses, " AND "), args...)
	case joinKindLeft:
		return qm.LeftOuterJoin(strings.Join(j.clauses, " AND "), args...)
	case joinKindRight:
		return qm.RightOuterJoin(strings.Join(j.clauses, " AND "), args...)
	}
	return nil
}
