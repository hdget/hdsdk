package wxoa

import (
	"github.com/hdget/hdsdk/lib/wx/typwx"
	"github.com/hdget/hdsdk/lib/wx/wxoa/cache"
)

type Wxoa interface {
	GetSignature(appId, appSecret, url string) (*typwx.WxoaSignature, error)
}

type implWxoa struct{}

var (
	_      Wxoa = (*implWxoa)(nil)
	_cache      = cache.New()
)

func New() Wxoa {
	return &implWxoa{}
}
