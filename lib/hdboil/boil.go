package hdboil

import (
	"fmt"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/lib/pagination"
	"github.com/hdget/hdsdk/protobuf"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
	"time"
)

func CommitOrRollback(tx boil.Transactor, err error) {
	if err != nil {
		e := tx.Rollback()
		if e != nil {
			hdsdk.Logger.Error("db rollback", "err", e)
		}
		return
	}

	e := tx.Commit()
	if e != nil {
		hdsdk.Logger.Error("db commit", "err", e)
	}
}

// GetLimitQueryMods 获取Limit相关QueryMods
func GetLimitQueryMods(list *protobuf.ListParam) []qm.QueryMod {
	p := pagination.NewWithParam(list)
	return []qm.QueryMod{qm.Offset(int(p.Offset)), qm.Limit(int(p.PageSize))}
}

// ByDesc OrderBy字段加入desc
func ByDesc(field string) string {
	return field + " DESC"
}

// JoinQueryMods 连接QueryMods
func JoinQueryMods(mods ...any) []qm.QueryMod {
	combined := make([]qm.QueryMod, 0)
	for _, mod := range mods {
		switch v := mod.(type) {
		case []qm.QueryMod:
			combined = append(combined, v...)
		case qm.QueryMod:
			combined = append(combined, v)
		}
	}
	return combined
}

func JoinAlias(thisTable, alias, thisColumn, thatTableColumn string, args ...string) string {
	clause := fmt.Sprintf("%s AS %s ON %s.%s=%s", thisTable, alias, alias, thisColumn, thatTableColumn)
	if len(args) > 0 {
		return clause + " AND " + strings.Join(args, " AND ")
	}
	return clause
}

func Join(thisTable, thisTableColumn, thatTableColumn string, args ...string) string {
	clause := fmt.Sprintf("%s ON %s=%s", thisTable, thisTableColumn, thatTableColumn)
	if len(args) > 0 {
		return clause + " AND " + strings.Join(args, " AND ")
	}
	return clause
}

// JoinTableColumn 需要指定Table和column
func JoinTableColumn(thisTable, thisColumn, thatTable, thatColumn string, args ...string) string {
	clause := fmt.Sprintf("%s ON %s.%s=%s.%s", thisTable, thisTable, thisColumn, thatTable, thatColumn)
	if len(args) > 0 {
		return clause + " AND " + strings.Join(args, " AND ")
	}
	return clause
}

// IfNullZeroString 如果传了args则用args[0]做为alias, 否则就用oldValue做为alias
func IfNullZeroString(oldValue string, args ...string) string {
	alias := oldValue
	if len(args) > 0 {
		alias = args[0]
	}
	return fmt.Sprintf("IFNULL((%s), '') AS %s", oldValue, alias)
}

// IfNullZeroNumber 如果传了args则用args[0]做为alias, 否则就用oldValue做为alias
func IfNullZeroNumber(oldValue string, args ...string) string {
	alias := oldValue
	if len(args) > 0 {
		alias = args[0]
	}
	return fmt.Sprintf("IFNULL((%s), 0) AS %s", oldValue, alias)
}

func GetUpdateCols(cols map[string]any, args ...string) map[string]any {
	updateColName := "updated_at"
	if len(args) > 0 {
		updateColName = args[0]
	}

	cols[updateColName] = time.Now().In(boil.GetLocation())
	return cols
}

func AsAliasColumn(alias, colName string) string {
	return fmt.Sprintf("`%s`.`%s` AS \"%s.%s\"", alias, colName, alias, colName)
}
