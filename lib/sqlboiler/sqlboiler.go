package sqlboiler

import (
	"fmt"
	"github.com/hdget/hdsdk/v2/lib/pagination"
	"github.com/hdget/hdsdk/v2/protobuf"
	jsonUtils "github.com/hdget/hdutils/json"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"reflect"
	"time"
)

// GetLimitQueryMods 获取Limit相关QueryMods
func GetLimitQueryMods(list *protobuf.ListParam) []qm.QueryMod {
	p := pagination.New(list)
	return []qm.QueryMod{qm.Offset(int(p.Offset)), qm.Limit(int(p.PageSize))}
}

func IfNull(column string, defaultValue any, args ...string) string {
	alias := column
	if len(args) > 0 {
		alias = args[0]
	}

	if reflect.ValueOf(defaultValue).IsZero() {
		switch v := defaultValue.(type) {
		case string:
			if v == "" {
				return fmt.Sprintf("IFNULL((%s), '') AS \"%s\"", column, alias)
			}
		case int8, int, int32, int64, uint8, uint, uint32, uint64:
			return fmt.Sprintf("IFNULL((%s), 0) AS \"%s\"", column, alias)
		case float32, float64:
			return fmt.Sprintf("IFNULL((%s), 0.0) AS \"%s\"", column, alias)
		case []byte:
			if jsonUtils.IsEmptyJsonObject(v) {
				return fmt.Sprintf("IFNULL((%s), '{}') AS \"%s\"", column, alias)
			} else if jsonUtils.IsEmptyJsonArray(v) {
				return fmt.Sprintf("IFNULL((%s), '[]') AS \"%s\"", column, alias)
			}
		}
	}

	return fmt.Sprintf("IFNULL((%s), '%v') AS \"%s\"", column, defaultValue, alias)
}

func IfNullWithColumn(column string, anotherColumn string, args ...string) string {
	alias := column
	if len(args) > 0 {
		alias = args[0]
	}
	return fmt.Sprintf("IFNULL((%s), %s) AS \"%s\"", column, anotherColumn, alias)
}

//// IfNullZeroString 如果传了args则用args[0]做为alias, 否则就用column做为alias
//func IfNullZeroString(column string, args ...string) string {
//	alias := column
//	if len(args) > 0 {
//		alias = args[0]
//	}
//	return fmt.Sprintf("IFNULL((%s), '') AS \"%s\"", column, alias)
//}
//
//// IfNullZeroNumber 如果传了args则用args[0]做为alias, 否则就用oldValue做为alias
//func IfNullZeroNumber(column string, args ...string) string {
//	alias := column
//	if len(args) > 0 {
//		alias = args[0]
//	}
//	return fmt.Sprintf("IFNULL((%s), 0) AS \"%s\"", column, alias)
//}
//
//// IfNullJsonObject 如果传了args则用args[0]做为alias, 否则就用column做为alias
//func IfNullJsonObject(column string, args ...string) string {
//	alias := column
//	if len(args) > 0 {
//		alias = args[0]
//	}
//	return fmt.Sprintf("IFNULL((%s), 0) AS \"%s\"", column, alias)
//}

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
	return IfNull(fmt.Sprintf("SUM(%s)", col), 0, args...)
}
