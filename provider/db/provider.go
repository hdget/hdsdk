package db

import (
	"github.com/hdget/hdsdk/types"
)

type BaseDbProvider struct {
	Default types.DbClient            // 缺省的数据库连接
	Main    types.DbClient            // 主数据库连接,可读写
	Slaves  []types.DbClient          // 只读从数据库列表，只可读
	Items   map[string]types.DbClient // 额外数据库
}

func (b *BaseDbProvider) My() types.DbClient {
	return b.Default
}

func (b *BaseDbProvider) Master() types.DbClient {
	return b.Main
}

func (b *BaseDbProvider) Slave(i int) types.DbClient {
	return b.Slaves[i]
}

func (b *BaseDbProvider) By(s string) types.DbClient {
	return b.Items[s]
}
