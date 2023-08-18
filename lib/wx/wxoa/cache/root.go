package cache

import (
	"fmt"
	"github.com/hdget/hdsdk"
)

const (
	tplWxoaAccessToken = "wxoa:%s:accesstoken" // string
	keyWxoaTicket      = "wxoa:ticket"         // string
)

type Cache interface {
	GetAccessToken(appId string) (string, error)
	SetAccessToken(appId, token string, expires int) error
	GetTicket() (string, error)
	SetTicket(ticket string, expires int) error
}

type implCache struct{}

var _ Cache = (*implCache)(nil)

func New() Cache {
	return &implCache{}
}

func (c *implCache) GetAccessToken(appId string) (string, error) {
	bs, err := hdsdk.Redis.My().Get(getAccessTokenKey(appId))
	return string(bs), err
}

func (c *implCache) SetAccessToken(appId, token string, expires int) error {
	return hdsdk.Redis.My().SetEx(getAccessTokenKey(appId), token, expires)
}

func (c *implCache) GetTicket() (string, error) {
	ticket, err := hdsdk.Redis.My().GetString(keyWxoaTicket)
	if err != nil {
		return "", nil
	}
	return ticket, nil
}

func (c *implCache) SetTicket(ticket string, expires int) error {
	return hdsdk.Redis.My().SetEx(keyWxoaTicket, ticket, expires)
}

func getAccessTokenKey(appId string) string {
	return fmt.Sprintf(tplWxoaAccessToken, appId)
}
