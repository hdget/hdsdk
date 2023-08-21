package cache

import (
	"fmt"
	"github.com/hdget/hdsdk"
)

const (
	tplWxqrAccessToken      = "wxqr:%s:access_token"  // string
	tplWxqrRfreshToken      = "wxqr:%s:refresh_token" // string
	wxqrRefreshTokenExpires = 30 * 24 * 60 * 60       // refresh token过期时间为30天
)

type Cache interface {
	GetAccessToken(appId string) (string, error)
	GetRefreshToken(appId string) (string, error)
	SetAccessToken(appId, accessToken string, expires int) error
	SetRefreshToken(appId, refreshToken string) error
}

type implCache struct{}

var _ Cache = (*implCache)(nil)

func New() Cache {
	return &implCache{}
}

func (c *implCache) GetAccessToken(appId string) (string, error) {
	bs, err := hdsdk.Redis.My().Get(getAccessToken(appId))
	return string(bs), err
}

func (c *implCache) SetAccessToken(appId, token string, expires int) error {
	return hdsdk.Redis.My().SetEx(getAccessToken(appId), token, expires)
}

func (c *implCache) GetRefreshToken(appId string) (string, error) {
	bs, err := hdsdk.Redis.My().Get(getRefreshToken(appId))
	return string(bs), err
}

func (c *implCache) SetRefreshToken(appId, token string) error {
	return hdsdk.Redis.My().SetEx(getRefreshToken(appId), token, wxqrRefreshTokenExpires)
}

func getAccessToken(appId string) string {
	return fmt.Sprintf(tplWxqrAccessToken, appId)
}

func getRefreshToken(appId string) string {
	return fmt.Sprintf(tplWxqrRfreshToken, appId)
}
