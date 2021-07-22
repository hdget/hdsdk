package cache

import (
	"github.com/hdget/sdk/types"
)

type BaseCacheProvider struct {
	Default types.CacheClient            //
	Items   map[string]types.CacheClient // 额外数据库
}

func (p *BaseCacheProvider) My() types.CacheClient {
	return p.Default
}

func (p *BaseCacheProvider) By(name string) types.CacheClient {
	return p.Items[name]
}
