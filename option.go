package hdsdk

import (
	"github.com/hdget/hdsdk/v2/intf"
)

type Option func(*SdkInstance)

// WithConfigProvider config provider option
func WithConfigProvider(configProvider intf.ConfigProvider) Option {
	return func(i *SdkInstance) {
		i.configProvider = configProvider
	}
}

func WithDebug() Option {
	return func(i *SdkInstance) {
		i.debug = true
	}
}
