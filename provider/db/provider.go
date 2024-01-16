package db

//
//type BaseDbProvider struct {
//	Default intf.DbClient            // 缺省的数据库连接
//	Main    intf.DbClient            // 主数据库连接,可读写
//	Slaves  []intf.DbClient          // 只读从数据库列表，只可读
//	Items   map[string]intf.DbClient // 额外数据库
//}
//
//func (b *BaseDbProvider) My() intf.DbClient {
//	return b.Default
//}
//
//func (b *BaseDbProvider) Master() intf.DbClient {
//	return b.Main
//}
//
//func (b *BaseDbProvider) Slave(i int) intf.DbClient {
//	return b.Slaves[i]
//}
//
//func (b *BaseDbProvider) By(s string) intf.DbClient {
//	return b.Items[s]
//}
