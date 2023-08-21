package log

import (
	"github.com/hdget/hdsdk/provider/log/typlog"
	"github.com/hdget/hdsdk/provider/log/zerolog"
	"github.com/hdget/hdsdk/types"
)

type LoggerImpl struct {
	types.LogProvider
}

const defaultProviderType = types.SdkTypeLogZerolog // 缺省的日志能力提供者

func (impl *LoggerImpl) Init(configer types.Configer, logger types.LogProvider, args ...interface{}) error {
	var p types.LogProvider
	switch len(args) {
	case 0:
		p = getProvider(defaultProviderType)
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
	impl.LogProvider = p

	return nil
}

func getProvider(capType types.SdkType) types.LogProvider {
	switch capType {
	case types.SdkTypeLogZerolog:
		return &zerolog.ZerologProvider{}
	}
	return nil
}
