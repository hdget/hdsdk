package provider

import "github.com/jmoiron/sqlx"

type BaseDbProvider struct {
	Default *sqlx.DB            // 缺省的数据库连接
	Main    *sqlx.DB            // 主数据库连接,可读写
	Slaves  []*sqlx.DB          // 只读从数据库列表，只可读
	Items   map[string]*sqlx.DB // 额外数据库
}

func (c *BaseDbProvider) My() *sqlx.DB {
	return c.Default
}

func (c *BaseDbProvider) Master() *sqlx.DB {
	return c.Main
}

func (c *BaseDbProvider) Slave(index int) *sqlx.DB {
	return c.Slaves[index]
}

func (c *BaseDbProvider) By(name string) *sqlx.DB {
	return c.Items[name]
}
