package log

import (
	"github.com/hdget/sdk/lib/log/provider"
	"github.com/hdget/sdk/lib/log/typlog"
	"github.com/hdget/sdk/types"
)

type CapImpl struct {
	types.LogProvider
}

// 缺省的日志能力提供者
const DEFAULT_PROVIDER_TYPE = types.LibLogZerolog

func (c *CapImpl) Init(configer types.Configer, logger types.LogProvider, args ...interface{}) error {
	var p types.LogProvider
	switch len(args) {
	case 0:
		p = getProvider(DEFAULT_PROVIDER_TYPE)
	// 第一个参数为指定的capability provider名字
	case 1:
		capType, ok := args[0].(types.SdkType)
		if ok {
			p = getProvider(capType)
		}
	}

	if p == nil {
		return typlog.ErrNoProvider
	}

	err := p.Init(configer, logger)
	if err != nil {
		return err
	}

	// if initialize ok, set anonymous interface to concrete provider
	c.LogProvider = p

	return nil
}

func getProvider(capType types.SdkType) types.LogProvider {
	switch capType {
	case types.LibLogZerolog:
		return &provider.ZerologProvider{}
	}
	return nil
}
