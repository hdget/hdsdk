package cache

import (
	"fmt"
	"github.com/hdget/hdsdk/v2"
	"github.com/hdget/hdsdk/v2/lib/weixin/types"
)

type ApiWeixinCache interface {
	GetAccessToken() (string, error)
	SetAccessToken(token string, expires int) error
	GetTicket() (string, error)
	SetTicket(ticket string, expires int) error
	GetSessKey() (string, error)
	SetSessKey(sessKey string, expires int) error
}

const (
	tplAccessToken = "wxmp:%s:accesstoken"
	// todo:wxmp和wxoa共享同一个access token, 为了兼容以前的版本将weixin改成wxmp
	//tplAccessToken = "weixin:%s:accesstoken"
	tplTicket     = "%s:%s:ticket"
	tplSessionKey = `%s:%s:session`
)

type weixinCacheImpl struct {
	App   types.WeixinApp
	appId string
}

var _ ApiWeixinCache = (*weixinCacheImpl)(nil)

func New(app types.WeixinApp, appId string) ApiWeixinCache {
	return &weixinCacheImpl{App: app, appId: appId}
}

func (c *weixinCacheImpl) GetAccessToken() (string, error) {
	bs, err := hdsdk.Redis().My().Get(fmt.Sprintf(tplAccessToken, c.appId))
	return string(bs), err
}

func (c *weixinCacheImpl) SetAccessToken(token string, expires int) error {
	return hdsdk.Redis().My().SetEx(fmt.Sprintf(tplAccessToken, c.appId), token, expires)
}

func (c *weixinCacheImpl) GetTicket() (string, error) {
	ticket, err := hdsdk.Redis().My().GetString(fmt.Sprintf(tplTicket, c.App, c.appId))
	if err != nil {
		return "", nil
	}
	return ticket, nil
}

func (c *weixinCacheImpl) SetTicket(ticket string, expires int) error {
	return hdsdk.Redis().My().SetEx(fmt.Sprintf(tplTicket, c.App, c.appId), ticket, expires)
}

func (c *weixinCacheImpl) GetSessKey() (string, error) {
	return hdsdk.Redis().My().GetString(fmt.Sprintf(tplSessionKey, c.App, c.appId))
}

func (c *weixinCacheImpl) SetSessKey(sessKey string, expires int) error {
	return hdsdk.Redis().My().SetEx(fmt.Sprintf(tplSessionKey, c.App, c.appId), sessKey, expires)
}
