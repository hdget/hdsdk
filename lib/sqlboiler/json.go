package sqlboiler

import (
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func JsonValue(jsonColumn string, jsonKey string, defaultValue any) qm.QueryMod {
	var template string
	switch v := defaultValue.(type) {
	case string:
		template = fmt.Sprintf("IFNULL(JSON_UNQUOTE(JSON_EXTRACT(%s, '$.%s')), '%s') AS %s", jsonColumn, jsonKey, v, jsonKey)
	case int8, int, int32, int64:
		template = fmt.Sprintf("IFNULL(JSON_EXTRACT(%s, '$.%s'), %d) AS %s", jsonColumn, jsonKey, v, jsonKey)
	case float32, float64:
		template = fmt.Sprintf("IFNULL(JSON_EXTRACT(%s, '$.%s'), %f) AS %s", jsonColumn, jsonKey, v, jsonKey)
	default:
		return nil
	}
	return qm.Select(template)
}

func JsonValueCompare(jsonColumn string, jsonKey string, operator string, compareValue any) qm.QueryMod {
	var template string
	switch v := compareValue.(type) {
	case string:
		template = fmt.Sprintf("JSON_UNQUOTE(JSON_EXTRACT(%s, '$.%s')) %s '%s'", jsonColumn, jsonKey, operator, v)
	case int8, int, int32, int64:
		template = fmt.Sprintf("JSON_EXTRACT(%s, '$.%s') %s %d", jsonColumn, jsonKey, operator, v)
	case float32, float64:
		template = fmt.Sprintf("JSON_EXTRACT(%s, '$.%s') %s %f", jsonColumn, jsonKey, operator, v)
	default:
		return nil
	}
	return qm.Where(template)
}
