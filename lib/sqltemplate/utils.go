package sqltemplate

import (
	"fmt"
	"github.com/hdget/hdutils/convert"
	"strings"
)

// ERROR:
// when select returns null value then there is not child table will be created
//
// 构造查询条件不同的批量查询SQL语句
// 目的是将for循环的多条查询语句合并成一个查询语句，一次取出所有结果，
// 减少与数据库的连接建立，来回的时间消耗。
// 原理：每一个查询语句的结构AS为一个表，然后实现多表关联，
// 将所有不同子表的查询字段的值都用","连接起来重新定义为一个新的查询字段,
// 将结果取到本地后，进行字段中值字符串的分割处理
// e,g:
//	SELECT
//		CONCAT_WS(",",table1.user_id, table2.user_id) AS user_id,
//  	CONCAT_WS(",",table1.amount, table2.amount) as amount
// 	FROM (
//			(
//				SELECT user_id, amount FROM fission_reward
//				WHERE user_id=123 AND created_at BETWEEN 1575376982 and (1575376982+24*3600*7)
// 			) AS table1,
//			(
//				SELECT user_id, amount FROM fission_reward
//				WHERE user_id=456 AND created_at BETWEEN 1575376982 and (1575376982+24*3600*7)
//			) AS table2
//	)

// GetUnionSQL 通过Union关键字实现多SQL批量查询
//
//			(
//			SELECT user_id, amount FROM fission_reward
//			WHERE user_id=123 AND created_at BETWEEN 1575376982 and (1575376982+24*3600*7)
//	     )
//		    UNION ALL
//			(
//		     	SELECT user_id, amount FROM fission_reward
//				WHERE user_id=456 AND created_at BETWEEN 1575376982 and (1575376982+24*3600*7)
//			)
func GetUnionSQL(subSQLs []string) string {
	return strings.Join(subSQLs, " UNION ALL ")
}

// GetBatchUpdateSQL 根据给定的主键列表对相应的值进行批量更新
// UPDATE <tableName> SET
//
//	<fieldname1> = CASE pkName
//	    WHEN <pk1> THEN <interface{}>
//	    WHEN <pk2> THEN <interface{}>
//	    WHEN <pk3> THEN <interface{}>
//	END,
//	<fieldname1> = CASE pkName
//	    WHEN <pk1> THEN <interface{}>
//	    WHEN <pk1> THEN <interface{}>
//	    WHEN <pk1> THEN <interface{}>
//	END
//
// WHERE pkName IN (<pks>)
//
//	values: <fieldName1> => {
//	     pk1 => interface{}
//	     pk2 => interface{}
//	     pk3 => interface{}
//	}
func GetBatchUpdateSQL(tableName, pkName string, pks []int64, values map[string]map[int64]interface{}) string {
	if len(pks) == 0 {
		return ""
	}

	// when SQL的template
	tmplWhenSQL := "WHEN %d THEN %v"
	tmplConditionSQL := "%s = CASE %s %s END"
	conditionSQLs := make([]string, 0)
	// 根据不同ID对应的
	for fieldName, dataMap := range values {
		whenSQLs := make([]string, 0)
		for pkValue, setValue := range dataMap {
			whenSQLs = append(whenSQLs, fmt.Sprintf(tmplWhenSQL, pkValue, setValue))
		}

		if len(whenSQLs) > 0 {
			conditionSQLs = append(conditionSQLs, fmt.Sprintf(tmplConditionSQL, fieldName, pkName, strings.Join(whenSQLs, " ")))
		}
	}

	// 如果没有要更新的值
	if len(conditionSQLs) == 0 {
		return ""
	}

	// convert []int64 to id1,id2,... string
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s IN (%s)", tableName, strings.Join(conditionSQLs, ","), pkName, convert.Int64sToCsv(pks))
	return sql
}
