package sqlboiler

import (
	"fmt"
	"github.com/hdget/hdsdk/v2/lib/pagination"
	"github.com/hdget/hdsdk/v2/protobuf"
	"github.com/hdget/hdutils/convert"
	jsonUtils "github.com/hdget/hdutils/json"
	reflectUtils "github.com/hdget/hdutils/reflect"
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

	if defaultValue == nil {
		return fmt.Sprintf("IFNULL((%s), '') AS \"%s\"", column, alias)
	}

	v := reflectUtils.Indirect(defaultValue)

	switch vv := reflect.ValueOf(v); vv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("IFNULL((%s), %d) AS \"%s\"", column, v, alias)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("IFNULL((%s), %.4f) AS \"%s\"", column, v, alias)
	case reflect.Slice:
		if vv.Type().Elem().Kind() == reflect.Uint8 {
			if jsonUtils.IsEmptyJsonObject(vv.Bytes()) {
				return fmt.Sprintf("IFNULL((%s), '{}') AS \"%s\"", column, alias)
			} else if jsonUtils.IsEmptyJsonArray(vv.Bytes()) {
				return fmt.Sprintf("IFNULL((%s), '[]') AS \"%s\"", column, alias)
			} else {
				return fmt.Sprintf("IFNULL((%s), '%s') AS \"%s\"", column, convert.BytesToString(vv.Bytes()), alias)
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
	if len(args) > 0 && args[0] != nil {
		return args[0]
	}
	return boil.GetDB()
}

func SUM(col string, args ...string) string {
	return IfNull(fmt.Sprintf("SUM(%s)", col), 0, args...)
}
