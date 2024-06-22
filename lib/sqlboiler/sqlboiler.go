package sqlboiler

import (
	"fmt"
	"github.com/hdget/hdsdk/v2/lib/pagination"
	"github.com/hdget/hdsdk/v2/protobuf"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
)

// GetLimitQueryMods 获取Limit相关QueryMods
func GetLimitQueryMods(list *protobuf.ListParam) []qm.QueryMod {
	p := pagination.New(list)
	return []qm.QueryMod{qm.Offset(int(p.Offset)), qm.Limit(int(p.PageSize))}
}

// IfNullZeroString 如果传了args则用args[0]做为alias, 否则就用oldValue做为alias
func IfNullZeroString(oldValue string, args ...string) string {
	alias := oldValue
	if len(args) > 0 {
		alias = args[0]
	}
	return fmt.Sprintf("IFNULL((%s), '') AS \"%s\"", oldValue, alias)
}

// IfNullZeroNumber 如果传了args则用args[0]做为alias, 否则就用oldValue做为alias
func IfNullZeroNumber(oldValue string, args ...string) string {
	alias := oldValue
	if len(args) > 0 {
		alias = args[0]
	}
	return fmt.Sprintf("IFNULL((%s), 0) AS \"%s\"", oldValue, alias)
}

// WithUpdateTime 除了cols中的会更新以外还会更新更新时间字段
func WithUpdateTime(cols map[string]any, args ...string) map[string]any {
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

func GetDB(args ...boil.Executor) boil.Executor {
	if len(args) > 0 {
		return args[0]
	}
	return boil.GetDB()
}

func SUM(col string, args ...string) string {
	return IfNullZeroNumber(fmt.Sprintf("SUM(%s)", col), args...)
}
