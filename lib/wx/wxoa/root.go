package wxoa

import (
	"hdsdk/lib/wx/typwx"
	"hdsdk/lib/wx/wxoa/cache"
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
