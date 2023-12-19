package hdboil

import (
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
)

func InnerJoin(thisTable, thisTableColumn, thatTableColumn string, args ...any) qm.QueryMod {
	clause := fmt.Sprintf("%s ON %s=%s", thisTable, thisTableColumn, thatTableColumn)
	return qm.InnerJoin(clause, args...)
}

func InnerJoinAlias(thisTable, asTable, thisColumn, thatTableColumn string, args ...any) qm.QueryMod {
	clause := fmt.Sprintf("%s AS %s ON %s.%s=%s", thisTable, asTable, asTable, thisColumn, thatTableColumn)
	return qm.InnerJoin(clause, args...)
}

// InnerJoinAnd inner join a on a.x=b.x AND a.y=b.y
func InnerJoinAnd(thisTable, thisTableColumn, thatTableColumn, andClauses []string, args ...any) qm.QueryMod {
	clause := fmt.Sprintf("%s ON %s=%s", thisTable, thisTableColumn, thatTableColumn)
	if len(andClauses) > 0 {
		clause = clause + " AND " + strings.Join(andClauses, " AND ")
	}
	return qm.InnerJoin(clause, args...)
}

func LeftJoin(thisTable, thisTableColumn, thatTableColumn string, args ...any) qm.QueryMod {
	clause := fmt.Sprintf("%s ON %s=%s", thisTable, thisTableColumn, thatTableColumn)
	return qm.LeftOuterJoin(clause, args...)
}

func LeftJoinAlias(thisTable, asTable, thisColumn, thatTableColumn string, args ...any) qm.QueryMod {
	clause := fmt.Sprintf("%s AS %s ON %s.%s=%s", thisTable, asTable, asTable, thisColumn, thatTableColumn)
	return qm.LeftOuterJoin(clause, args...)
}

func LeftJoinAnd(thisTable, thisTableColumn, thatTableColumn, andClauses []string, args ...any) qm.QueryMod {
	clause := fmt.Sprintf("%s ON %s=%s", thisTable, thisTableColumn, thatTableColumn)
	if len(andClauses) > 0 {
		clause = clause + " AND " + strings.Join(andClauses, " AND ")
	}
	return qm.LeftOuterJoin(clause, args...)
}
