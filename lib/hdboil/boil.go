package hdboil

import (
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/lib/pagination"
	"github.com/hdget/hdsdk/protobuf"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
